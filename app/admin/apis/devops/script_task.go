package devops

import (
	"errors"
	_ "github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	_ "github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	_ "github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/task"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var _scriptTaskApi = &ScriptTaskApi{
	scriptTaskService: task.GetScriptTaskService(),
}

type ScriptTaskApi struct {
	scriptTaskService *task.ScriptTaskService
}

func GetScriptTaskApi() *ScriptTaskApi {
	return _scriptTaskApi
}

// Add godoc
//
//	@Summary	添加任务
//	@Tags		脚本任务管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.ScriptTaskUpsertReq	true	"任务信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/script-task [post]
func (a *ScriptTaskApi) Add(c *gin.Context) {
	upsertReq := &req.ScriptTaskUpsertReq{}
	if err := c.ShouldBindJSON(upsertReq); err != nil {
		E.PanicErr(err)
	}
	if err := upsertReq.Verify(); err != nil {
		E.PanicErr(err)
	}
	if err := a.scriptTaskService.Create(upsertReq); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// Update godoc
//
//	@Summary	更新任务
//	@Tags		脚本任务管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		path		int64				true	"任务ID"
//	@Param		[body]	body		req.TaskUpsertReq	true	"任务信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/script-task/:id [put]
func (a *ScriptTaskApi) Update(c *gin.Context) {
	upsertReq := &req.ScriptTaskUpsertReq{}
	if e := c.ShouldBindJSON(upsertReq); e != nil {
		E.PanicErr(e)
	}
	if err := upsertReq.Verify(); err != nil {
		E.PanicErr(err)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err = a.scriptTaskService.Update(id, upsertReq); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// Delete godoc
//
//	@Summary	删除任务
//	@Tags		脚本任务管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"任务ID"
//	@Success	200	{object}	R.Result
//	@Router		/devops/script-task/:id [delete]
func (a *ScriptTaskApi) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err = a.scriptTaskService.Delete(id); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetList godoc
//
//	@Summary	查询任务列表
//	@Tags		脚本任务管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照名称模糊搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[task.ScriptTask]}
//	@Router		/devops/script-task/list [get]
func (a *ScriptTaskApi) GetList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableTasks, err := a.scriptTaskService.PageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableTasks)
}

// Start godoc
//
//	@Summary	启动任务
//	@Tags		脚本任务管理
//	@Accept		application/json
//	@Produce	application/json
//	@Success	200	{object}	R.Result
//	@Router		/devops/script-task/start/:id [post]
func (a *ScriptTaskApi) Start(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	userId := c.GetInt64(constant.SessionUserId)
	if err = a.scriptTaskService.Start(userId, id); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// Stop godoc
//
//	@Summary	停止任务
//	@Tags		脚本任务管理
//	@Accept		application/json
//	@Produce	application/json
//	@Success	200	{object}	R.Result
//	@Router		/devops/script-task/stop/:id [post]
func (a *ScriptTaskApi) Stop(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err = a.scriptTaskService.Stop(id); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}
