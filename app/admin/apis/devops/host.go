package devops

import (
	"errors"
	_ "github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	_ "github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	_ "github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var _hostApi = &HostApi{
	hostService: dvservice.GetHostService(),
}

type HostApi struct {
	hostService *dvservice.HostService
}

func GetHostApi() *HostApi {
	return _hostApi
}

const passwordLenLimit = 50

var passwordLenLimitHint = "密码长度不能超过" + strconv.Itoa(passwordLenLimit) + "字符"

// AddHost godoc
//
//	@Summary	添加服务器
//	@Tags		服务器管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.HostUpsertReq	true	"服务器信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/host [post]
func (a *HostApi) AddHost(c *gin.Context) {
	hostUpsertReq := &req.HostUpsertReq{}
	if err := c.ShouldBindJSON(hostUpsertReq); err != nil {
		E.PanicErr(err)
	}
	if len([]rune(hostUpsertReq.Password)) > passwordLenLimit {
		R.Fail(c, passwordLenLimitHint, http.StatusBadRequest)
		return
	}
	err := a.hostService.CreateHost(hostUpsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateHost godoc
//
//	@Summary	更新服务器
//	@Tags		服务器管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		path		int64				true	"服务器ID"
//	@Param		[body]	body		req.HostUpsertReq	true	"服务器信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/host/:id [put]
func (a *HostApi) UpdateHost(c *gin.Context) {
	upsertReq := &req.HostUpsertReq{}
	if e := c.ShouldBindJSON(upsertReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if upsertReq.PasswordChanged && len([]rune(upsertReq.Password)) > passwordLenLimit {
		R.Fail(c, passwordLenLimitHint, http.StatusBadRequest)
		return
	}
	err = a.hostService.UpdateHost(id, upsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteHost godoc
//
//	@Summary	删除服务器
//	@Tags		服务器管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"服务器ID"
//	@Success	200	{object}	R.Result
//	@Router		/devops/host/:id [delete]
func (a *HostApi) DeleteHost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.hostService.DeleteHost(id)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetList godoc
//
//	@Summary	查询服务器列表
//	@Tags		服务器管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照名称模糊搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[dvmodel.Host]}
//	@Router		/devops/host/list [get]
func (a *HostApi) GetList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableHosts, err := a.hostService.PageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableHosts)
}

// GetAll godoc
//
//	@Summary	查询全部服务器列表
//	@Tags		服务器管理
//	@Produce	application/json
//	@Success	200	{object}	R.Result{value=res.PageableData[res.HostBasicRes]}
//	@Router		/devops/host/all [get]
func (a *HostApi) GetAll(c *gin.Context) {
	results, err := a.hostService.FindAll()
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, results)
}

// SshConnectTest godoc
//
//	@Summary	SSH连接测试
//	@Tags		服务器管理
//	@Produce	application/json
//	@Param		ip				query		string	true	"IP地址"
//	@Param		port			query		int16	true	"端口"
//	@Param		user			query		string	true	"账号"
//	@Param		password		query		string	true	"密码"
//	@Param		mode			query		string	true	"模式：0，新增；1，编辑"
//	@Param		passwordChanged	query		bool	true	"是否修改了密码（仅编辑时有效）"
//	@Success	200				{object}	R.Result
//	@Router		/devops/host/connect-test [get]
func (a *HostApi) SshConnectTest(c *gin.Context) {
	params := &req.SshConnectTestParams{}
	if err := c.ShouldBindQuery(params); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, map[string]bool{
		"result": a.hostService.SshConnectTest(params),
	})
}
