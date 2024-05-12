package dvservice

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"github.com/MarchGe/go-admin-server/config"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"
)

var _appService = &AppService{}

type AppService struct {
}

func GetAppService() *AppService {
	return _appService
}

func (s *AppService) CreateApp(info *req.AppUpsertReq) error {
	existApp, _ := s.FindOneByNameAndVersion(info.Name, info.Version)
	if existApp != nil {
		return E.Message(fmt.Sprintf("应用'%s:%s'已存在", info.Name, info.Version))
	}
	app := s.toModel(info)
	app.CreateTime = time.Now()
	app.UpdateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		app.Key = s.getPkgKey(info.Name, info.Version, info.FileName)
		if err := tx.Save(app).Error; err != nil {
			return err
		}
		return s.moveTmpPkg(info.Key, app.Key)
	})
	return err
}

func (s *AppService) getPkgKey(appName, appVersion, fileName string) string {
	return appName + "/" + appVersion + "/" + fileName
}

func (s *AppService) GetUploadTmpDir(workDir string) string {
	return path.Clean(workDir) + "/.tmp"
}

func (s *AppService) moveTmpPkg(tmpKey, key string) error {
	cfg := config.GetConfig()
	tmpFile := s.GetUploadTmpDir(cfg.WorkDir) + "/" + tmpKey
	pkgFile := path.Clean(cfg.UploadPkgPath) + "/" + key
	if err := os.MkdirAll(path.Dir(pkgFile), 0755); err != nil {
		return fmt.Errorf("create directory %s error, %w", path.Dir(pkgFile), err)
	}
	if err := os.Rename(tmpFile, pkgFile); err != nil {
		return fmt.Errorf("move file from %s to %s error, %w", tmpFile, pkgFile, err)
	}
	return nil
}

func (s *AppService) movePkg(oldKey, newKey string) error {
	cfg := config.GetConfig()
	oldFile := path.Clean(cfg.UploadPkgPath) + "/" + oldKey
	newFile := path.Clean(cfg.UploadPkgPath) + "/" + newKey
	if err := os.MkdirAll(path.Dir(newFile), 0755); err != nil {
		return fmt.Errorf("create directory %s error, %w", path.Dir(newFile), err)
	}
	if err := os.Rename(oldFile, newFile); err != nil {
		return fmt.Errorf("move file from %s to %s error, %w", oldFile, newFile, err)
	}
	s.removeOldPkg(oldKey)
	return nil
}

func (s *AppService) toModel(info *req.AppUpsertReq) *dvmodel.App {
	return &dvmodel.App{
		Name:     info.Name,
		Version:  info.Version,
		Port:     info.Port,
		Key:      info.Key,
		FileName: info.FileName,
	}
}

func (s *AppService) UpdateApp(id int64, info *req.AppUpsertReq) error {
	app, _ := s.FindOneById(id)
	if app == nil {
		return E.Message("操作的应用不存在")
	}
	oldKey := app.Key
	existApp, _ := s.FindOneByNameAndVersion(info.Name, info.Version)
	if existApp != nil && existApp.Id != id {
		return E.Message(fmt.Sprintf("应用'%s:%s'已存在", info.Name, info.Version))
	}
	s.copyProperties(info, app)
	app.UpdateTime = time.Now()
	app.Key = s.getPkgKey(app.Name, app.Version, app.FileName)
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(app).Error; err != nil {
			return err
		}
		if info.Key != oldKey { // 上传了新的部署包
			if err := s.moveTmpPkg(info.Key, app.Key); err != nil {
				return err
			}
			if oldKey != app.Key {
				s.removeOldPkg(oldKey)
			}
		} else if app.Key != oldKey { // 没有上传新的部署包，但包路径发生变化
			if err := s.movePkg(oldKey, app.Key); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *AppService) FindOneById(id int64) (*dvmodel.App, error) {
	m := &dvmodel.App{}
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

func (s *AppService) copyProperties(info *req.AppUpsertReq, app *dvmodel.App) {
	app.Name = info.Name
	app.Version = info.Version
	app.Port = info.Port
	app.Key = info.Key
	app.FileName = info.FileName
}

func (s *AppService) DeleteApp(id int64) error {
	app, _ := s.FindOneById(id)
	if app == nil {
		return E.Message("操作的应用不存在")
	}
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&dvmodel.App{}, id).Error; err != nil {
			return err
		}
		s.removeOldPkg(app.Key)
		return nil
	})
	return err
}

func (s *AppService) PageList(keyword string, page, pageSize int) (*res.PageableData[*dvmodel.App], error) {
	apps := make([]*dvmodel.App, 0)
	pageableData := &res.PageableData[*dvmodel.App]{}
	db := database.GetMysql().Model(&dvmodel.App{})
	if keyword != "" {
		db.Where("name like ?", "%"+keyword+"%")
	}
	var count int64
	err := db.Count(&count).Order("create_time desc").Offset(pageSize * (page - 1)).Limit(pageSize).Find(&apps).Error
	if err != nil {
		return nil, err
	}
	pageableData.List = apps
	pageableData.Total = count
	return pageableData, nil
}

func (s *AppService) FindOneByNameAndVersion(name, version string) (*dvmodel.App, error) {
	m := &dvmodel.App{}
	err := database.GetMysql().Where("name = ? and version = ?", name, version).First(m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		} else {
			return nil, err
		}
	}
	return m, nil
}

func (s *AppService) removeOldPkg(oldKey string) {
	if oldKey == "" {
		return
	}
	uploadRoot := path.Clean(config.GetConfig().UploadPkgPath)
	filePath := uploadRoot + "/" + oldKey
	if err := os.Remove(filePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("delete file error", slog.String("file", filePath), slog.Any("err", err))
	}
	keyLevels := len(strings.Split(oldKey, "/")) - 1
	loopPath := filePath
	for i := 0; i < keyLevels; i++ { // 循环，判断是否需要删除上层目录
		loopPath = path.Dir(loopPath)
		entries, e := os.ReadDir(loopPath)
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			slog.Error("read dir error", slog.Any("err", e))
		}
		if len(entries) == 0 {
			if err := os.Remove(loopPath); err != nil && !errors.Is(e, os.ErrNotExist) {
				slog.Error("delete dir error", slog.String("dir", loopPath), slog.Any("err", err))
			}
		}
	}
}
