package apis

import (
	"errors"
	_ "github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	_ "github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var _roleApi = &RoleApi{
	roleService: service.GetRoleService(),
}

type RoleApi struct {
	roleService *service.RoleService
}

func GetRoleApi() *RoleApi {
	return _roleApi
}

// AddRole godoc
//
//	@Summary	添加角色
//	@Tags		角色管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.RoleUpsertReq	true	"角色信息"
//	@Success	200		{object}	R.Result
//	@Router		/role [post]
func (a *RoleApi) AddRole(c *gin.Context) {
	roleUpsertReq := &req.RoleUpsertReq{}
	if err := c.ShouldBindJSON(roleUpsertReq); err != nil {
		E.PanicErr(err)
	}
	err := a.roleService.CreateRole(roleUpsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateRole godoc
//
//	@Summary	更新角色
//	@Tags		角色管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		path		int64				true	"角色ID"
//	@Param		[body]	body		req.RoleUpsertReq	true	"角色信息"
//	@Success	200		{object}	R.Result
//	@Router		/role/:id [put]
func (a *RoleApi) UpdateRole(c *gin.Context) {
	upsertReq := &req.RoleUpsertReq{}
	if e := c.ShouldBindJSON(upsertReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.roleService.UpdateRole(id, upsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteRole godoc
//
//	@Summary	删除角色
//	@Tags		角色管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"角色ID"
//	@Success	200	{object}	R.Result
//	@Router		/role/:id [delete]
func (a *RoleApi) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.roleService.DeleteRole(id)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetRole godoc
//
//	@Summary	查询角色信息
//	@Tags		角色管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"角色ID"
//	@Success	200	{object}	R.Result{value=model.Role}
//	@Router		/role/:id [get]
func (a *RoleApi) GetRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	role, _ := a.roleService.FindOneById(id)
	if role == nil {
		R.Fail(c, "查询的用户不存在")
		return
	}
	R.Success(c, role)
}

// GetList godoc
//
//	@Summary	查询角色列表
//	@Tags		角色管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照名称模糊搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[model.Role]}
//	@Router		/role/list [get]
func (a *RoleApi) GetList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableRoles, err := a.roleService.PageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableRoles)
}

// UpdateRoleMenus godoc
//
//	@Summary	角色绑定菜单
//	@Tags		角色管理
//	@Accept		application/json
//	@Produce	application/json
//
//	@Param		[body]	body		req.RoleMenusUpdateReq	true	"菜单列表"
//
//	@Success	200		{object}	R.Result
//	@Router		/role/:id/menus [put]
func (a *RoleApi) UpdateRoleMenus(c *gin.Context) {
	updateReq := &req.RoleMenusUpdateReq{}
	if e := c.ShouldBindJSON(updateReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err = a.roleService.UpdateRoleMenus(id, updateReq); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}
