package apis

import (
	"github.com/MarchGe/go-admin-server/app/admin/service"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

var _authApi = &AuthApi{
	authService: service.GetAuthService(),
	userService: service.GetUserService(),
}

type AuthApi struct {
	authService *service.AuthService
	userService *service.UserService
}

func GetAuthApi() *AuthApi {
	return _authApi
}

// Login godoc
//
//	@Summary	登录
//	@Tags		账号相关
//	@Accept		application/json
//	@Produce	application/json
//
//	@Param		[body]	body		req.LoginReq	true	"登录信息"
//
//	@Success	200		{object}	R.Result{value=model.User}
//	@Router		/auth/login [post]
func (a *AuthApi) Login(c *gin.Context) {
	loginReq := &req.LoginReq{}
	if err := c.ShouldBindJSON(loginReq); err != nil {
		E.PanicErr(err)
	}

	loginRes, err := a.authService.Login(loginReq)
	if err != nil {
		E.PanicErr(err)
	}

	session := sessions.DefaultMany(c, constant.LoginSession)
	session.Set(constant.SessionUserId, loginRes.Id)
	session.Set(constant.IsRootUser, loginRes.Root)
	if err = session.Save(); err != nil {
		E.PanicErr(err)
	}
	recorder.RecordLoginLog(c, loginRes.Id)
	R.Success(c, loginRes)
}

// ChangeMyPasswd godoc
//
//	@Summary	修改密码
//	@Tags		账号相关
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.ChangePasswdReq	true	"修改密码信息"
//	@Success	200		{object}	R.Result
//	@Router		/auth/passwd [put]
func (a *AuthApi) ChangeMyPasswd(c *gin.Context) {
	passwdReq := &req.ChangeMyPasswdReq{}
	if err := c.ShouldBindJSON(passwdReq); err != nil {
		E.PanicErr(err)
	}
	myId := c.GetInt64(constant.SessionUserId)
	err := a.authService.ChangeMyPassword(myId, passwdReq.OldPassword, passwdReq.NewPassword)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetMyInfo godoc
//
//	@Summary	查询我的信息
//	@Tags		账号相关
//	@Produce	application/json
//	@Success	200	{object}	R.Result{value=model.User}
//	@Router		/auth/my-info [get]
func (a *AuthApi) GetMyInfo(c *gin.Context) {
	myId := c.GetInt64(constant.SessionUserId)
	user := a.userService.FindOneById(myId, "Dept", "RoleList", "JobList")
	R.Success(c, user)
}

// UpdateMyInfo godoc
//
//	@Summary	更新我的信息
//	@Tags		账号相关
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.UpsertMyInfoReq	true	"用户信息"
//	@Success	200		{object}	R.Result
//	@Router		/auth/my-info [put]
func (a *AuthApi) UpdateMyInfo(c *gin.Context) {
	upsertReq := &req.UpsertMyInfoReq{}
	if err := c.ShouldBindJSON(upsertReq); err != nil {
		E.PanicErr(err)
	}

	myId := c.GetInt64(constant.SessionUserId)
	err := a.authService.UpdateMyInfo(myId, upsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetWebShellToken godoc
//
//	@Summary	获取WebShell连接需要的token
//	@Tags		账号相关
//	@Produce	application/json
//	@Success	200	{object}	R.Result
//	@Router		/auth/web-shell-token [get]
func (a *AuthApi) GetWebShellToken(c *gin.Context) {
	loginToken, err := c.Cookie(constant.LoginSession)
	if err != nil {
		E.PanicErr(err)
	}
	timestamp := strconv.Itoa(int(time.Now().UnixMilli()))
	shellToken, err := utils.EncryptString(config.GetConfig().Cookie.SecretKey, loginToken+"$"+timestamp, "")
	if err != nil {
		E.PanicErr(err)
	}
	result := map[string]string{
		"token": shellToken,
	}
	R.Success(c, result)
}
