package devops

import (
	"errors"
	_ "github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
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

var _taskApi = &TaskApi{
	taskService: task.GetTaskService(),
}

type TaskApi struct {
	taskService *task.TaskService
}

func GetTaskApi() *TaskApi {
	return _taskApi
}

// AddTask godoc
//
//	@Summary	添加任务
//	@Tags		任务管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.TaskUpsertReq	true	"任务信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/task [post]
func (a *TaskApi) AddTask(c *gin.Context) {
	taskUpsertReq := &req.TaskUpsertReq{}
	if err := c.ShouldBindJSON(taskUpsertReq); err != nil {
		E.PanicErr(err)
	}
	if err := a.taskService.CreateTask(taskUpsertReq); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateTask godoc
//
//	@Summary	更新任务
//	@Tags		任务管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		path		int64				true	"任务ID"
//	@Param		[body]	body		req.TaskUpsertReq	true	"任务信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/task/:id [put]
func (a *TaskApi) UpdateTask(c *gin.Context) {
	upsertReq := &req.TaskUpsertReq{}
	if e := c.ShouldBindJSON(upsertReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err = a.taskService.UpdateTask(id, upsertReq); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteTask godoc
//
//	@Summary	删除任务
//	@Tags		任务管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"任务ID"
//	@Success	200	{object}	R.Result
//	@Router		/devops/task/:id [delete]
func (a *TaskApi) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err = a.taskService.DeleteTask(id); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetList godoc
//
//	@Summary	查询任务列表
//	@Tags		任务管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照名称模糊搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[dvmodel.Task]}
//	@Router		/devops/task/list [get]
func (a *TaskApi) GetList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableTasks, err := a.taskService.PageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableTasks)
}

// StartTask godoc
//
//	@Summary	启动任务
//	@Tags		任务管理
//	@Accept		application/json
//	@Produce	application/json
//	@Success	200		{object}	R.Result
//	@Router		/devops/task/start/:id [post]
func (a *TaskApi) StartTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	userId := c.GetInt64(constant.SessionUserId)
	if err = a.taskService.StartTask(userId, id); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// StopTask godoc
//
//	@Summary	停止任务
//	@Tags		任务管理
//	@Accept		application/json
//	@Produce	application/json
//	@Success	200		{object}	R.Result
//	@Router		/devops/task/stop/:id [post]
func (a *TaskApi) StopTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err = a.taskService.StopTask(id); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}
