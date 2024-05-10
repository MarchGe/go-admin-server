package service

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"github.com/MarchGe/go-admin-server/app/common/middleware/authz"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/MarchGe/go-admin-server/app/common/utils/email"
	"github.com/MarchGe/go-admin-server/config"
	"gorm.io/gorm"
	"log/slog"
	"strings"
	"time"
)

var _userService = &UserService{}

type UserService struct {
}

func GetUserService() *UserService {
	return _userService
}

func (s *UserService) PageList(keyword string, sex int8, status int8, start *time.Time, end *time.Time, page int, pageSize int) *res.PageableData[*model.User] {
	db := database.GetMysql()
	db = db.Model(&model.User{})
	if keyword != "" {
		db = db.Where("nickname like ? or name like ? or email like ? or cellphone like ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sex != -1 {
		db = db.Where("sex = ?", sex)
	}
	if status != -1 {
		db = db.Where("status = ?", status)
	}
	if start != nil {
		db = db.Where("create_time >= ?", *start)
	}
	if end != nil {
		db = db.Where("create_time <= ?", *end)
	}
	db = db.Order("root desc, id desc").Preload("Dept").Preload("JobList").Preload("RoleList").Preload("MenuList")

	pageableUsers := &res.PageableData[*model.User]{}
	users := make([]*model.User, 0)
	var count int64
	err := db.Count(&count).Limit(pageSize).Offset(pageSize * (page - 1)).Find(&users).Error
	if err != nil {
		E.PanicErr(err)
	}
	pageableUsers.List = users
	pageableUsers.Total = count
	return pageableUsers
}

func (s *UserService) ChangeUserStatus(id int64, status int8) error {
	db := database.GetMysql()
	err := db.Model(&model.User{}).Where("id = ?", id).UpdateColumn("status", status).Error
	return err
}

func (s *UserService) IsRootUser(id int64) (result bool) {
	err := database.GetMysql().Model(&model.User{}).Select("root").Where("id = ?", id).Scan(&result).Error
	if err != nil {
		E.PanicErr(err)
	}
	return
}

func (s *UserService) FindOneByEmail(email string) *model.User {
	db := database.GetMysql()
	user := &model.User{}
	err := db.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			E.PanicErr(err)
		}
	}
	return user
}

func (s *UserService) CreateUser(info *req.UserUpsertReq) error {
	existUser := s.FindOneByEmail(info.Email)
	if existUser != nil {
		return E.Message("邮箱已存在")
	}
	user := s.toModel(info)
	user.CreateTime = time.Now()
	user.UpdateTime = time.Now()
	if user.Nickname == "" {
		user.Nickname = strings.Split(user.Email, "@")[0]
	}
	user.Status = model.UserStatusNormal
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Root", "Dept", "RoleList", "MenuList", "JobList").Save(user).Error; err != nil {
			return err
		}
		generatedPasswd := utils.RandomString(8)
		up := &model.UserPassword{
			UserId:   user.Id,
			Password: GetAuthService().PasswordHash(generatedPasswd),
		}
		if err := tx.Save(up).Error; err != nil {
			return err
		}
		if err := s.updateUserRoles(tx, user.Id, info.RoleIds); err != nil {
			return err
		}
		if err := s.updateUserJobs(tx, user.Id, info.JobIds); err != nil {
			return err
		}
		if info.ShouldSendPassword {
			s.sendUserAddedEmailAsync(user.Email, generatedPasswd)
		}
		return nil
	})
	return err
}

func (s *UserService) UpdateUser(id int64, info *req.UserUpsertReq) error {
	user := s.FindOneById(id)
	if user == nil {
		return E.Message("操作的用户不存在")
	}
	existUser := s.FindOneByEmail(info.Email)
	if existUser != nil && user.Id != existUser.Id {
		return E.Message("邮箱已存在")
	}
	s.copyProperties(info, user)
	user.UpdateTime = time.Now()
	tx := database.GetMysql().Begin()
	if err := tx.Omit("Root", "Dept", "RoleList", "MenuList", "JobList").Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := s.updateUserRoles(tx, user.Id, info.RoleIds); err != nil {
		tx.Rollback()
		return err
	}
	if err := s.updateUserJobs(tx, user.Id, info.JobIds); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *UserService) toModel(info *req.UserUpsertReq) *model.User {
	return &model.User{
		Name:      info.Name,
		Email:     info.Email,
		Cellphone: info.Cellphone,
		Nickname:  info.Nickname,
		Birthday:  info.Birthday,
		Sex:       info.Sex,
		DeptId:    info.DeptId,
	}
}

func (s *UserService) copyProperties(info *req.UserUpsertReq, user *model.User) {
	user.Sex = info.Sex
	user.Email = info.Email
	user.Birthday = info.Birthday
	user.Name = info.Name
	user.Nickname = info.Nickname
	user.Cellphone = info.Cellphone
	user.DeptId = info.DeptId
}

func (s *UserService) FindOneById(id int64, preloads ...string) *model.User {
	user := &model.User{}
	db := database.GetMysql()
	if len(preloads) > 0 {
		for _, preload := range preloads {
			db = db.Preload(preload)
		}
	}
	err := db.First(user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			E.PanicErr(err)
		}
	}
	return user
}

func (s *UserService) DeleteUser(id int64) error {
	user := s.FindOneById(id)
	if user == nil {
		return E.Message("操作的用户不存在")
	}
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(user).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.UserPassword{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.UserMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.RoleUser{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.UserJob{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&model.Settings{}).Error; err != nil {
			return err
		}
		return authz.DeleteUser(authz.UserSub(id))
	})
	return err
}

func (s *UserService) sendUserAddedEmailAsync(userEmail, password string) {
	go func() {
		mailConfig := config.GetConfig().Email
		htmlMailer := email.CreateHtmlMailer(&mailConfig.MailConfig)
		data := make(map[string]any)
		data["SystemName"] = mailConfig.SystemName
		data["AccessUrl"] = mailConfig.AccessUrl
		data["Email"] = userEmail
		data["Password"] = password
		msg, err := htmlMailer.MsgBuilder().
			To([]string{userEmail}).
			Subject("账号开通成功通知").
			HtmlTemplate(common.UserCreateEmailTemplate).
			Data(data).
			Build()
		if err != nil {
			slog.Error("build email msg error", slog.Any("err", err))
		}
		if err = htmlMailer.Send(msg); err != nil {
			slog.Error("send email error", slog.Any("err", err))
		}
	}()
}

func (s *UserService) sendResetPasswordEmailAsync(userEmail, password string) {
	go func() {
		mailConfig := config.GetConfig().Email
		htmlMailer := email.CreateHtmlMailer(&mailConfig.MailConfig)
		data := make(map[string]any)
		data["SystemName"] = mailConfig.SystemName
		data["AccessUrl"] = mailConfig.AccessUrl
		data["Email"] = userEmail
		data["Password"] = password
		msg, err := htmlMailer.MsgBuilder().
			To([]string{userEmail}).
			Subject("密码已重置成功").
			HtmlTemplate(common.PasswordResetEmailTemplate).
			Data(data).
			Build()
		if err != nil {
			slog.Error("build email msg error", slog.Any("err", err))
		}
		if err = htmlMailer.Send(msg); err != nil {
			slog.Error("send email error", slog.Any("err", err))
		}
	}()
}

func (s *UserService) UpdateUserMenus(id int64, updateReq *req.UserMenusUpdateReq) error {
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		// 删除user和menu的关联关系
		if err := tx.Where("user_id = ?", id).Delete(&model.UserMenu{}).Error; err != nil {
			return err
		}
		if len(updateReq.Ids) > 0 {
			// 添加user和menu的关联关系
			menus := make([]*model.UserMenu, 0, len(updateReq.Ids))
			for _, menuId := range updateReq.Ids {
				menus = append(menus, &model.UserMenu{
					UserId: id,
					MenuId: menuId,
				})
			}
			if err := tx.Save(menus).Error; err != nil {
				return err
			}
		}
		symbols, err := GetMenuService().FindSymbolsByIds(updateReq.Ids)
		if err != nil {
			return fmt.Errorf("find menu symbols error, %w", err)
		}
		return authz.UpdateSubPolicies(authz.UserSub(id), symbols)
	})
	return err
}

func (s *UserService) updateUserRoles(tx *gorm.DB, userId int64, roleIds []int64) error {
	if err := tx.Where("user_id = ?", userId).Delete(&model.RoleUser{}).Error; err != nil {
		return fmt.Errorf("delete user and role relations error, %w", err)
	}
	if len(roleIds) > 0 {
		roleUsers := make([]*model.RoleUser, len(roleIds))
		for i, roleId := range roleIds {
			roleUsers[i] = &model.RoleUser{
				RoleId: roleId,
				UserId: userId,
			}
		}
		if err := tx.Save(&roleUsers).Error; err != nil {
			return fmt.Errorf("save user and role relations error, %w", err)
		}
	}
	// 更新组策略
	return authz.UpdateGroupingPolicies(authz.UserSub(userId), s.getRoleSubs(roleIds))
}

func (s *UserService) updateUserJobs(tx *gorm.DB, userId int64, jobIds []int64) error {
	if err := tx.Where("user_id = ?", userId).Delete(&model.UserJob{}).Error; err != nil {
		return fmt.Errorf("delete user and job relations error, %w", err)
	}
	if len(jobIds) > 0 {
		userJobs := make([]*model.UserJob, len(jobIds))
		for i, jobId := range jobIds {
			userJobs[i] = &model.UserJob{
				UserId: userId,
				JobId:  jobId,
			}
		}
		if err := tx.Save(&userJobs).Error; err != nil {
			return fmt.Errorf("save user and job relations error, %w", err)
		}
	}
	return nil
}

func (s *UserService) getRoleSubs(roleIds []int64) []string {
	groups := make([]string, len(roleIds))
	for i, roleId := range roleIds {
		groups[i] = authz.RoleSub(roleId)
	}
	return groups
}

func (s *UserService) ChangePassword(id int64, password string) error {
	passwdModel, _ := GetAuthService().FindPasswordById(id)
	if passwdModel == nil {
		return E.Message("操作的用户密码记录不存在")
	}
	passwdModel.Password = GetAuthService().PasswordHash(password)
	return database.GetMysql().Save(passwdModel).Error
}

func (s *UserService) ResetPassword(id int64) error {
	user := s.FindOneById(id)
	if user == nil {
		return E.Message("操作的用户记录不存在")
	}
	passwdModel, _ := GetAuthService().FindPasswordById(id)
	if passwdModel == nil {
		return E.Message("操作的用户密码记录不存在")
	}
	generatedPasswd := utils.RandomString(8)
	passwdModel.Password = GetAuthService().PasswordHash(generatedPasswd)
	err := database.GetMysql().Save(passwdModel).Error
	if err != nil {
		return err
	}
	s.sendResetPasswordEmailAsync(user.Email, generatedPasswd)
	return nil
}
