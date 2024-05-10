package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"gorm.io/gorm"
	"strings"
)

var _authService = &AuthService{}

type AuthService struct {
}

func GetAuthService() *AuthService {
	return _authService
}

func (a *AuthService) Login(req *req.LoginReq) (*model.User, error) {
	user := GetUserService().FindOneByEmail(req.Email)
	if user == nil {
		return nil, E.Message("账号不存在")
	}
	if user.Status != model.UserStatusNormal {
		return nil, E.Message("账号已禁用")
	}
	var db = database.GetMysql()
	userPassword := &model.UserPassword{}
	err := db.Model(&model.UserPassword{}).Where("user_id = ?", user.Id).First(userPassword).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message(fmt.Sprintf("账号数据异常，密码记录不存在，账号ID：%d", user.Id))
		}
		return nil, err
	}
	sha1Passwd := a.PasswordHash(req.Password)
	if strings.ToLower(userPassword.Password) != strings.ToLower(sha1Passwd) {
		return nil, E.Message("密码错误")
	}
	return user, nil
}

func (a *AuthService) ChangeMyPassword(userId int64, oldPassword, newPassword string) error {
	passwdModel, _ := a.FindPasswordById(userId)
	if passwdModel == nil {
		return E.Message("操作的用户密码记录不存在")
	}
	oldPasswd := a.PasswordHash(oldPassword)
	if strings.ToLower(oldPasswd) != strings.ToLower(passwdModel.Password) {
		return E.Message("旧密码错误")
	}
	passwdModel.Password = a.PasswordHash(newPassword)
	return database.GetMysql().Save(passwdModel).Error
}

func (a *AuthService) PasswordHash(password string) string {
	passwdBytes := sha1.Sum([]byte(password))
	return hex.EncodeToString(passwdBytes[:])
}

func (a *AuthService) FindPasswordById(id int64) (*model.UserPassword, error) {
	up := &model.UserPassword{}
	err := database.GetMysql().Where("user_id = ?", id).First(up).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("查询的记录不存在")
		} else {
			return nil, err
		}
	}
	return up, nil
}

func (a *AuthService) UpdateMyInfo(id int64, upsertReq *req.UpsertMyInfoReq) error {
	user := GetUserService().FindOneById(id)
	if user == nil {
		return E.Message("操作的用户不存在")
	}
	user.Nickname = upsertReq.Nickname
	user.Cellphone = upsertReq.Cellphone
	return database.GetMysql().Save(user).Error
}
