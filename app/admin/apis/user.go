package apis

import (
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	_ "github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var _userApi = &UserApi{
	userService: service.GetUserService(),
}

type UserApi struct {
	userService *service.UserService
}

func GetUserApi() *UserApi {
	return _userApi
}

// AddUser godoc
//
//	@Summary		添加用户
//	@Description	添加用户，添加成功后，密码会通过邮箱发给用户
//	@Tags			用户管理
//	@Accept			application/json
//	@Produce		application/json
//	@Param			[body]	body		req.UserUpsertReq	true	"用户信息"
//	@Success		200		{object}	R.Result
//	@Router			/user [post]
func (a *UserApi) AddUser(c *gin.Context) {
	upsertReq := &req.UserUpsertReq{}
	if err := c.ShouldBindJSON(upsertReq); err != nil {
		E.PanicErr(err)
	}
	err := a.userService.CreateUser(upsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetList godoc
//
//	@Summary	查询用户列表
//	@Tags		用户管理
//	@Produce	application/json
//	@Param		keyword		query		string		false	"搜索关键字（支持按名字、昵称、邮箱、手机号模糊搜索）"
//	@Param		status		query		int64		false	"用户状态"
//	@Param		page		query		int64		false	"页码"
//	@Param		pageSize	query		int64		false	"每页查询条数"
//	@Param		dates		query		[]string	false	"开始时间和结束时间"
//	@Success	200			{object}	R.Result{value=res.PageableData[model.User]}
//	@Router		/user/list [get]
func (a *UserApi) GetList(c *gin.Context) {
	keyword := c.Query("keyword")
	sex, err1 := ginUtils.GetInt8Query(c, "sex", -1)
	status, err2 := ginUtils.GetInt8Query(c, "status", -1)
	page, err3 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err4 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2, err3, err4)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	dates := c.QueryArray("dates[]")
	var start, end *time.Time
	if dates != nil && len(dates) == 2 {
		time1, e1 := time.Parse(constant.DateFormat, dates[0])
		time2, e2 := time.Parse(constant.DateFormat, dates[1])
		if err = errors.Join(e1, e2); err != nil {
			R.Fail(c, "参数中的时间解析出错！")
			return
		}
		if time1.Before(time2) {
			start = &time1
			end = &time2
		} else {
			start = &time2
			end = &time1
		}
	}

	pageableUsers := a.userService.PageList(keyword, sex, status, start, end, page, pageSize)
	R.Success(c, pageableUsers)
}

// EnableAccount godoc
//
//	@Summary	启用账号
//	@Tags		用户管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"用户ID"
//	@Success	200	{object}	R.Result
//	@Router		/user/enable/:id [put]
func (a *UserApi) EnableAccount(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.userService.ChangeUserStatus(id, model.UserStatusNormal)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DisableAccount godoc
//
//	@Summary	禁用账号
//	@Tags		用户管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"用户ID"
//	@Success	200	{object}	R.Result
//	@Router		/user/disable/:id [put]
func (a *UserApi) DisableAccount(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if a.userService.IsRootUser(id) {
		R.Fail(c, "不允许对超级用户执行该操作", http.StatusBadRequest)
		return
	}
	err = a.userService.ChangeUserStatus(id, model.UserStatusDisable)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateUser godoc
//
//	@Summary	更新用户信息
//	@Tags		用户管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		path		int64				true	"用户ID"
//	@Param		[body]	body		req.UserUpsertReq	true	"用户信息"
//	@Success	200		{object}	R.Result
//	@Router		/user/:id [put]
func (a *UserApi) UpdateUser(c *gin.Context) {
	upsertReq := &req.UserUpsertReq{}
	if e := c.ShouldBindJSON(upsertReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	root := c.GetBool(constant.IsRootUser)
	if !root && a.userService.IsRootUser(id) {
		R.Fail(c, "不允许对超级用户执行该操作", http.StatusBadRequest)
		return
	}
	err = a.userService.UpdateUser(id, upsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteUser godoc
//
//	@Summary	删除用户
//	@Tags		用户管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"用户ID"
//	@Success	200	{object}	R.Result
//	@Router		/user/:id [delete]
func (a *UserApi) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if a.userService.IsRootUser(id) {
		R.Fail(c, "超级用户不允许删除", http.StatusBadRequest)
		return
	}
	err = a.userService.DeleteUser(id)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetUser godoc
//
//	@Summary	查询用户信息
//	@Tags		用户管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"用户ID"
//	@Success	200	{object}	R.Result{value=model.User}
//	@Router		/user/:id [get]
func (a *UserApi) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	user := a.userService.FindOneById(id)
	if user == nil {
		R.Fail(c, "没有查到用户信息", http.StatusBadRequest)
		return
	}
	R.Success(c, user)
}

// ChangePasswd godoc
//
//	@Summary	修改密码
//	@Tags		用户管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.ChangePasswdReq	true	"修改密码信息"
//	@Success	200		{object}	R.Result
//	@Router		/user/passwd [put]
func (a *UserApi) ChangePasswd(c *gin.Context) {
	passwdReq := &req.ChangePasswdReq{}
	if err := c.ShouldBindJSON(passwdReq); err != nil {
		E.PanicErr(err)
	}
	root := c.GetBool(constant.IsRootUser)
	if !root && a.userService.IsRootUser(passwdReq.Id) {
		R.Fail(c, "不允许对超级用户执行该操作", http.StatusBadRequest)
		return
	}
	err := a.userService.ChangePassword(passwdReq.Id, passwdReq.NewPassword)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// ResetPasswd godoc
//
//	@Summary	重置密码
//	@Tags		用户管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.ChangePasswdReq	true	"修改密码信息"
//	@Success	200		{object}	R.Result
//	@Router		/user/passwd-reset/:id [put]
func (a *UserApi) ResetPasswd(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		E.PanicErr(err)
	}
	root := c.GetBool(constant.IsRootUser)
	if !root && a.userService.IsRootUser(id) {
		R.Fail(c, "不允许对超级用户执行该操作", http.StatusBadRequest)
		return
	}
	if err = a.userService.ResetPassword(id); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateUserMenus godoc
//
//	@Summary	用户绑定菜单
//	@Tags		用户管理
//	@Accept		application/json
//	@Produce	application/json
//
//	@Param		id		path		int64					true	"用户ID"
//
//	@Param		[body]	body		req.UserMenusUpdateReq	true	"菜单列表"
//	@Success	200		{object}	R.Result
//	@Router		/user/:id/menus [put]
func (a *UserApi) UpdateUserMenus(c *gin.Context) {
	updateReq := &req.UserMenusUpdateReq{}
	if e := c.ShouldBindJSON(updateReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if a.userService.IsRootUser(id) {
		R.Fail(c, "无需对超级用户执行该操作", http.StatusBadRequest)
		return
	}
	if err = a.userService.UpdateUserMenus(id, updateReq); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}
