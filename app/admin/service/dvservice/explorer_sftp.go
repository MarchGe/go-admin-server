package dvservice

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	dvRes "github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/fs"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"
)

var _explorerSftpService = &ExplorerSftpService{}

type ExplorerSftpService struct {
}

func GetExplorerSftpService() *ExplorerSftpService {
	return _explorerSftpService
}

func (s *ExplorerSftpService) ListEntries(parentDir string, host *dvmodel.Host) ([]*dvRes.ExplorerEntry, error) {
	client, err := s.getSftpClient(host)
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.Close() }()

	info, err := client.Stat(parentDir)
	if err != nil {
		pathError := &fs.PathError{}
		if errors.As(err, &pathError) {
			return nil, E.Message("父目录参数有误")
		}
		if errors.Is(err, os.ErrNotExist) {
			return nil, E.Message("访问的目录不存在")
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, E.Message("父目录参数有误")
	}
	dirEntries, err := client.ReadDir(parentDir)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return nil, E.Message("没有权限访问该目录")
		}
		return nil, err
	}
	var length = len(dirEntries)
	entries := make([]*dvRes.ExplorerEntry, length)
	for i, item := range dirEntries {
		entry := &dvRes.ExplorerEntry{
			Name: item.Name(),
			Type: parseType(item.Mode()),
		}
		entries[i] = entry
	}
	return entries, nil
}

func (s *ExplorerSftpService) getSftpClient(host *dvmodel.Host) (*sftp.Client, error) {
	sshClient, err := s.sshConnect(host)
	if err != nil {
		slog.Error("ssh connect error", slog.Any("err", err))
		return nil, E.Message("SSH连接失败")
	}
	return sftp.NewClient(sshClient)
}

func (s *ExplorerSftpService) sshConnect(host *dvmodel.Host) (*ssh.Client, error) {
	decryptPasswd, err := utils.DecryptString(config.GetConfig().EncryptKey, host.Password, "")
	if err != nil {
		return nil, err
	}
	clientConfig := ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(decryptPasswd),
		},
		Timeout:         constant.SshEstablishTimeoutInSeconds * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", host.Ip, host.Port)
	return ssh.Dial("tcp", addr, &clientConfig)
}

func (s *ExplorerSftpService) DeleteEntry(path string, host *dvmodel.Host) error {
	client, err := s.getSftpClient(host)
	if err != nil {
		return fmt.Errorf("get sftp client error, %w", err)
	}
	defer func() { _ = client.Close() }()
	return client.RemoveAll(path)
}

func (s *ExplorerSftpService) UploadFile(filePath string, file multipart.File, host *dvmodel.Host) error {
	client, err := s.getSftpClient(host)
	if err != nil {
		return fmt.Errorf("get sftp client error, %w", err)
	}
	defer func() { _ = client.Close() }()
	parentDir := filepath.ToSlash(filepath.Dir(filePath))
	if err = client.MkdirAll(parentDir); err != nil {
		return fmt.Errorf("sftp mkdir error, %w", err)
	}
	f, err := client.Create(filePath)
	if err != nil {
		return fmt.Errorf("sftp create file error, %w", err)
	}
	defer func() { _ = f.Close() }()
	_, err = io.Copy(f, file)
	return err
}

func (s *ExplorerSftpService) DownloadFile(filePath string, host *dvmodel.Host) (*sftp.File, error) {
	client, err := s.getSftpClient(host)
	if err != nil {
		return nil, fmt.Errorf("get sftp client error, %w", err)
	}
	defer func() { _ = client.Close() }()

	info, err := client.Stat(filePath)
	if err != nil {
		return nil, E.Message("获取文件信息失败")
	}
	if info.IsDir() {
		return nil, E.Message("不支持下载文件夹")
	}

	file, err := client.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("sftp open file error, %w", err)
	}
	return file, nil
}

func (s *ExplorerSftpService) CreateDir(r *req.SftpCreateDirReq, host *dvmodel.Host) error {
	client, err := s.getSftpClient(host)
	if err != nil {
		return fmt.Errorf("get sftp client error, %w", err)
	}
	defer func() { _ = client.Close() }()
	dir := path.Clean(r.Dir + "/" + r.Name)
	_, err = client.Stat(dir)
	if !os.IsNotExist(err) {
		return E.Message("目录已存在")
	}
	return client.MkdirAll(dir)
}

func (s *ExplorerSftpService) Rename(r *req.SftpRenameReq, host *dvmodel.Host) error {
	client, err := s.getSftpClient(host)
	if err != nil {
		return fmt.Errorf("get sftp client error, %w", err)
	}
	defer func() { _ = client.Close() }()
	oldPath := path.Clean(r.Dir + "/" + r.OldName)
	newPath := path.Clean(r.Dir + "/" + r.NewName)
	if oldPath == newPath {
		return nil
	}
	_, err = client.Stat(newPath)
	if !os.IsNotExist(err) {
		return E.Message("目标文件（或目录）已存在")
	}
	if err = client.Rename(oldPath, newPath); err != nil {
		if os.IsNotExist(err) {
			return E.Message("旧文件（或目录）不存在")
		}
		if errors.Is(err, os.ErrPermission) {
			return E.Message("文件系统：permission denied")
		}
		return err
	}
	return nil
}
