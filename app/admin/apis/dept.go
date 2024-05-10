package apis

import (
	"github.com/MarchGe/go-admin-server/app/admin/service"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var _deptApi = &DeptApi{
	deptService: service.GetDeptService(),
}

type DeptApi struct {
	deptService *service.DeptService
}

func GetDeptApi() *DeptApi {
	return _deptApi
}

// AddDept godoc
//
//	@Summary	添加部门
//	@Tags		部门管理
//	@Accept		application/json
//	@Produce	application/json
//
//	@Param		[body]	body		req.DeptUpsertReq	true	"部门信息"
//
//	@Success	200		{object}	R.Result
//	@Router		/dept [post]
func (a *DeptApi) AddDept(c *gin.Context) {
	m := &req.DeptUpsertReq{}
	if err := c.ShouldBindJSON(m); err != nil {
		E.PanicErr(err)
	}
	if err := a.deptService.CreateDept(m); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateDept godoc
//
//	@Summary	修改部门
//	@Tags		部门管理
//	@Accept		application/json
//	@Produce	application/json
//
//	@Param		id		path		int64				true	"部门ID"
//	@Param		[body]	body		req.DeptUpsertReq	true	"部门信息"
//
//	@Success	200		{object}	R.Result
//	@Router		/dept/:id [put]
func (a *DeptApi) UpdateDept(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	m := &req.DeptUpsertReq{}
	if err = c.ShouldBindJSON(m); err != nil {
		E.PanicErr(err)
	}

	err = a.deptService.UpdateDept(id, m)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteDept godoc
//
//	@Summary	删除部门
//	@Tags		部门管理
//	@Produce	application/json
//
//	@Param		id	path		int64	true	"部门ID"
//
//	@Success	200	{object}	R.Result
//	@Router		/dept/:id [delete]
func (a *DeptApi) DeleteDept(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	exists, err := a.deptService.ExistsByParentId(id)
	if err != nil {
		E.PanicErr(err)
	}
	if exists {
		R.Fail(c, "该部门下存在子部门，不允许删除", http.StatusBadRequest)
		return
	}

	err = a.deptService.DeleteDept(id)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetDept godoc
//
//	@Summary	查询部门信息
//	@Tags		部门管理
//	@Produce	application/json
//
//	@Param		id	path		int64	true	"部门ID"
//
//	@Success	200	{object}	R.Result{value=model.Dept}
//	@Router		/dept/:id [get]
func (a *DeptApi) GetDept(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	dept, _ := a.deptService.FindOneById(id)
	if dept == nil {
		R.Fail(c, "查询的部门不存在", http.StatusBadRequest)
		return
	}
	R.Success(c, dept)
}

// GetDeptTree godoc
//
//	@Summary	查询部门树
//	@Tags		部门管理
//	@Produce	application/json
//	@Param		keyword	query		string	false	"按照名称模糊搜索"
//	@Success	200		{object}	R.Result{value=res.DeptTree}
//	@Router		/dept/tree [get]
func (a *DeptApi) GetDeptTree(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	trees, err := a.deptService.FindDeptTree(keyword)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, trees)
}
