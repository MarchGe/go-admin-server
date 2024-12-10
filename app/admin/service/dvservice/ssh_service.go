package dvservice

import (
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/MarchGe/go-admin-server/config"
	"golang.org/x/crypto/ssh"
	"time"
)

var _sshService = &SshService{}

type SshService struct {
}

func GetSshService() *SshService {
	return _sshService
}

func (s *SshService) CreateSshClient(host *dvmodel.Host) (*ssh.Client, error) {
	encryptKey := config.GetConfig().EncryptKey
	decryptPasswd, err := utils.DecryptString(encryptKey, host.Password, "")
	if err != nil {
		return nil, fmt.Errorf("decrypt password error, %w", err)
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
	client, err := ssh.Dial("tcp", addr, &clientConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh connect failed, %w", err)
	}
	return client, nil
}
