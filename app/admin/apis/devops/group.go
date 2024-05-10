package devops

import (
	"errors"
	_ "github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	_ "github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var _groupApi = &GroupApi{
	groupService: dvservice.GetGroupService(),
}

type GroupApi struct {
	groupService *dvservice.GroupService
}

func GetGroupApi() *GroupApi {
	return _groupApi
}

// AddGroup godoc
//
//	@Summary	添加分组
//	@Tags		服务器分组管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.GroupUpsertReq	true	"分组信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/group [post]
func (a *GroupApi) AddGroup(c *gin.Context) {
	groupUpsertReq := &req.GroupUpsertReq{}
	if err := c.ShouldBindJSON(groupUpsertReq); err != nil {
		E.PanicErr(err)
	}
	err := a.groupService.CreateGroup(groupUpsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateGroup godoc
//
//	@Summary	更新分组
//	@Tags		服务器分组管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		path		int64				true	"分组ID"
//	@Param		[body]	body		req.GroupUpsertReq	true	"分组信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/group/:id [put]
func (a *GroupApi) UpdateGroup(c *gin.Context) {
	upsertReq := &req.GroupUpsertReq{}
	if e := c.ShouldBindJSON(upsertReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.groupService.UpdateGroup(id, upsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteGroup godoc
//
//	@Summary	删除分组
//	@Tags		服务器分组管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"分组ID"
//	@Success	200	{object}	R.Result
//	@Router		/devops/group/:id [delete]
func (a *GroupApi) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.groupService.DeleteGroup(id)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetList godoc
//
//	@Summary	查询分组列表
//	@Tags		服务器分组管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照名称模糊搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[dvmodel.Group]}
//	@Router		/devops/group/list [get]
func (a *GroupApi) GetList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableGroups, err := a.groupService.PageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableGroups)
}
