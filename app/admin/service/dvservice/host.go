package dvservice

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	dvRes "github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/MarchGe/go-admin-server/config"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

var _hostService = &HostService{}

type HostService struct {
}

func GetHostService() *HostService {
	return _hostService
}

func (s *HostService) CreateHost(info *req.HostUpsertReq) error {
	existHost, _ := s.FindOneByIp(info.Ip)
	if existHost != nil {
		return E.Message("主机'" + info.Ip + "'已存在")
	}
	host, err := s.toModel(info)
	if err != nil {
		return err
	}
	host.CreateTime = time.Now()
	host.UpdateTime = time.Now()
	err = database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(host).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *HostService) toModel(info *req.HostUpsertReq) (*dvmodel.Host, error) {
	cfg := config.GetConfig()
	encryptPasswd, err := utils.EncryptString(cfg.EncryptKey, info.Password, "")
	if err != nil {
		return nil, err
	}
	return &dvmodel.Host{
		Name:     info.Name,
		Ip:       info.Ip,
		Port:     info.Port,
		User:     info.User,
		Password: encryptPasswd,
		SortNum:  info.SortNum,
	}, nil
}

func (s *HostService) UpdateHost(id int64, info *req.HostUpsertReq) error {
	host, _ := s.FindOneById(id)
	if host == nil {
		return E.Message("操作的主机不存在")
	}
	existHost, _ := s.FindOneByIp(info.Ip)
	if existHost != nil && existHost.Id != id {
		return E.Message("主机'" + info.Ip + "'已存在")
	}
	if err := s.copyProperties(info, host); err != nil {
		return err
	}
	host.UpdateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(host).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *HostService) FindOneById(id int64) (*dvmodel.Host, error) {
	m := &dvmodel.Host{}
	err := database.GetMysql().First(m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}

func (s *HostService) copyProperties(info *req.HostUpsertReq, host *dvmodel.Host) error {
	host.Name = info.Name
	host.Ip = info.Ip
	host.Port = info.Port
	host.User = info.User
	host.SortNum = info.SortNum
	if info.PasswordChanged {
		cfg := config.GetConfig()
		encryptPasswd, err := utils.EncryptString(cfg.EncryptKey, info.Password, "")
		if err != nil {
			return err
		}
		host.Password = encryptPasswd
	}
	return nil
}

func (s *HostService) DeleteHost(id int64) error {
	host, _ := s.FindOneById(id)
	if host == nil {
		return E.Message("操作的主机不存在")
	}
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&dvmodel.Host{}, id).Error; err != nil {
			return err
		}
		if err := tx.Where("host_id = ?", id).Delete(&dvmodel.HostGroup{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *HostService) PageList(keyword string, page, pageSize int) (*res.PageableData[*dvmodel.Host], error) {
	hosts := make([]*dvmodel.Host, 0)
	pageableData := &res.PageableData[*dvmodel.Host]{}
	db := database.GetMysql().Model(&dvmodel.Host{})
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("sort_num").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&hosts).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = hosts
	pageableData.Total = count
	return pageableData, nil
}

func (s *HostService) FindOneByIp(ip string) (*dvmodel.Host, error) {
	m := &dvmodel.Host{}
	err := database.GetMysql().Where("ip = ?", ip).First(m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}

func (s *HostService) FindAll() ([]*dvRes.HostBasicRes, error) {
	results := make([]*dvRes.HostBasicRes, 0)
	if err := database.GetMysql().Model(&dvmodel.Host{}).Select("id, name, ip").Order("sort_num").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (s *HostService) SshConnectTest(params *req.SshConnectTestParams) (connectSuccess bool) {
	addr := fmt.Sprintf("%s:%d", params.Ip, params.Port)
	password := params.Password
	if params.Mode == req.HostUpdateMode && !params.PasswordChanged {
		decryptPasswd, err := utils.DecryptString(config.GetConfig().EncryptKey, password, "")
		if err != nil {
			return
		}
		password = decryptPasswd
	}
	clientConfig := ssh.ClientConfig{
		User: params.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout:         constant.SshEstablishTimeoutInSeconds * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", addr, &clientConfig)
	if err != nil {
		slog.Error("ssh connect failed", slog.Any("err", err))
		return
	}
	defer func() { _ = client.Close() }()
	connectSuccess = true
	return
}
