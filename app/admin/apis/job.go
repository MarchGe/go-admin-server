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

var _jobApi = &JobApi{
	jobService: service.GetJobService(),
}

type JobApi struct {
	jobService *service.JobService
}

func GetJobApi() *JobApi {
	return _jobApi
}

// AddJob godoc
//
//	@Summary	添加岗位
//	@Tags		岗位管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.JobUpsertReq	true	"岗位信息"
//	@Success	200		{object}	R.Result
//	@Router		/job [post]
func (a *JobApi) AddJob(c *gin.Context) {
	jobUpsertReq := &req.JobUpsertReq{}
	if err := c.ShouldBindJSON(jobUpsertReq); err != nil {
		E.PanicErr(err)
	}
	err := a.jobService.CreateJob(jobUpsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateJob godoc
//
//	@Summary	更新岗位
//	@Tags		岗位管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		path		int64				true	"岗位ID"
//	@Param		[body]	body		req.JobUpsertReq	true	"岗位信息"
//	@Success	200		{object}	R.Result
//	@Router		/job/:id [put]
func (a *JobApi) UpdateJob(c *gin.Context) {
	upsertReq := &req.JobUpsertReq{}
	if e := c.ShouldBindJSON(upsertReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.jobService.UpdateJob(id, upsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteJob godoc
//
//	@Summary	删除岗位
//	@Tags		岗位管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"岗位ID"
//	@Success	200	{object}	R.Result
//	@Router		/job/:id [delete]
func (a *JobApi) DeleteJob(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.jobService.DeleteJob(id)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetList godoc
//
//	@Summary	查询岗位列表
//	@Tags		岗位管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照名称模糊搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[model.Job]}
//	@Router		/job/list [get]
func (a *JobApi) GetList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableJobs, err := a.jobService.PageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableJobs)
}
