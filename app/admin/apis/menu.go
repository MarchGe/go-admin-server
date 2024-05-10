package apis

import (
	"github.com/MarchGe/go-admin-server/app/admin/service"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var _menuApi = &MenuApi{
	menuService: service.GetMenuService(),
}

type MenuApi struct {
	menuService *service.MenuService
}

func GetMenuApi() *MenuApi {
	return _menuApi
}

// AddMenu godoc
//
//	@Summary	添加菜单
//	@Tags		菜单管理
//	@Accept		application/json
//	@Produce	application/json
//
//	@Param		[body]	body		req.MenuUpsertReq	true	"菜单信息"
//
//	@Success	200		{object}	R.Result
//	@Router		/menu [post]
func (a *MenuApi) AddMenu(c *gin.Context) {
	m := &req.MenuUpsertReq{}
	if err := c.ShouldBindJSON(m); err != nil {
		E.PanicErr(err)
	}

	err := a.menuService.CreateMenu(m)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateMenu godoc
//
//	@Summary	修改菜单
//	@Tags		菜单管理
//	@Accept		application/json
//	@Produce	application/json
//
//	@Param		id		path		int64				true	"菜单ID"
//	@Param		[body]	body		req.MenuUpsertReq	true	"菜单信息"
//
//	@Success	200		{object}	R.Result
//	@Router		/menu/:id [put]
func (a *MenuApi) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	m := &req.MenuUpsertReq{}
	if err = c.ShouldBindJSON(m); err != nil {
		E.PanicErr(err)
	}

	err = a.menuService.UpdateMenu(id, m)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteMenu godoc
//
//	@Summary	删除菜单
//	@Tags		菜单管理
//	@Produce	application/json
//
//	@Param		id	path		int64	true	"菜单ID"
//
//	@Success	200	{object}	R.Result
//	@Router		/menu/:id [delete]
func (a *MenuApi) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	exists, err := a.menuService.ExistsByParentId(id)
	if err != nil {
		E.PanicErr(err)
	}
	if exists {
		R.Fail(c, "该菜单下存在子菜单，不允许删除", http.StatusBadRequest)
		return
	}

	err = a.menuService.DeleteMenu(id)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetMenu godoc
//
//	@Summary	查询菜单信息
//	@Tags		菜单管理
//	@Produce	application/json
//
//	@Param		id	path		int64	true	"菜单ID"
//
//	@Success	200	{object}	R.Result{value=model.Menu}
//	@Router		/menu/:id [get]
func (a *MenuApi) GetMenu(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	menu, _ := a.menuService.FindOneById(id)
	if menu == nil {
		R.Fail(c, "查询的菜单不存在", http.StatusBadRequest)
		return
	}
	R.Success(c, menu)
}

// GetMenuTree godoc
//
//	@Summary	查询菜单树
//	@Tags		菜单管理
//	@Produce	application/json
//	@Success	200	{object}	R.Result{value=res.MenuTree}
//	@Router		/menu/tree [get]
func (a *MenuApi) GetMenuTree(c *gin.Context) {
	trees, err := a.menuService.FindMenuTree()
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, trees)
}

// GetMyMenuTree godoc
//
//	@Summary	查询当前用户有权限访问的菜单树（仅查询显示的菜单）
//	@Tags		菜单管理
//	@Produce	application/json
//	@Success	200	{object}	R.Result{value=res.MenuTree}
//	@Router		/menu/my/tree [get]
func (a *MenuApi) GetMyMenuTree(c *gin.Context) {
	isRoot := c.GetBool(constant.IsRootUser)
	userId := c.GetInt64(constant.SessionUserId)
	trees, err := a.menuService.FindUserMenuTree(userId, isRoot)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, trees)
}

// GetMyPermissions godoc
//
//	@Summary	查询当前用户拥有的所有权限列表
//	@Tags		菜单管理
//	@Produce	application/json
//	@Success	200	{object}	R.Result{value=[]string}
//	@Router		/menu/my/permissions [get]
func (a *MenuApi) GetMyPermissions(c *gin.Context) {
	isRoot := c.GetBool(constant.IsRootUser)
	userId := c.GetInt64(constant.SessionUserId)
	permissions, err := a.menuService.FindUserPermissions(userId, isRoot)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, permissions)
}
