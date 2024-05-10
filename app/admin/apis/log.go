package apis

import (
	"errors"
	"fmt"
	_ "github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service"
	_ "github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var _logApi = &LogApi{
	logService: service.GetLogService(),
}

type LogApi struct {
	logService *service.LogService
}

func GetLogApi() *LogApi {
	return _logApi
}

// GetLoginLogList godoc
//
//	@Summary	查询登录日志列表
//	@Tags		日志管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照用户名搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[model.LoginLog]}
//	@Router		/log/login [get]
func (a *LogApi) GetLoginLogList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableData, err := a.logService.LoginLogPageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableData)
}

// GetOpLogList godoc
//
//	@Summary	查询操作日志列表
//	@Tags		日志管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照用户名搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[model.OpLog]}
//	@Router		/log/op [get]
func (a *LogApi) GetOpLogList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableData, err := a.logService.OpLogPageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableData)
}

// GetExceptionLogList godoc
//
//	@Summary	查询异常日志列表
//	@Tags		日志管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照用户名搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[model.ExceptionLog]}
//	@Router		/log/exception [get]
func (a *LogApi) GetExceptionLogList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableData, err := a.logService.ExceptionLogPageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableData)
}

// DeleteLoginLog godoc
//
//	@Summary	清除登录日志
//	@Tags		日志管理
//	@Produce	application/json
//	@Param		type	query		int16	true	"清除方式"
//	@Success	200		{object}	R.Result
//	@Router		/log/login [delete]
func (a *LogApi) DeleteLoginLog(c *gin.Context) {
	clearType, err := ginUtils.GetInt16Query(c, "type", 0)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	beforeTime, err := getTime(clearType)
	if err != nil {
		E.PanicErr(err)
	}
	if err = a.logService.DeleteLoginLog(beforeTime); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteOpLog godoc
//
//	@Summary	清除操作日志
//	@Tags		日志管理
//	@Produce	application/json
//	@Param		type	query		int16	true	"清除方式"
//	@Success	200		{object}	R.Result
//	@Router		/log/op [delete]
func (a *LogApi) DeleteOpLog(c *gin.Context) {
	clearType, err := ginUtils.GetInt16Query(c, "type", 0)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	beforeTime, err := getTime(clearType)
	if err != nil {
		E.PanicErr(err)
	}
	if err = a.logService.DeleteOpLog(beforeTime); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteExceptionLog godoc
//
//	@Summary	清除异常日志
//	@Tags		日志管理
//	@Produce	application/json
//	@Param		type	query		int16	true	"清除方式"
//	@Success	200		{object}	R.Result
//	@Router		/log/exception [delete]
func (a *LogApi) DeleteExceptionLog(c *gin.Context) {
	clearType, err := ginUtils.GetInt16Query(c, "type", 0)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	beforeTime, err := getTime(clearType)
	if err != nil {
		E.PanicErr(err)
	}
	if err = a.logService.DeleteExceptionLog(beforeTime); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

func getTime(clearType int16) (time.Time, error) {
	switch clearType {
	case 1:
		return time.Now(), nil
	case 2:
		return time.Now().AddDate(0, 0, -7), nil
	case 3:
		return time.Now().AddDate(0, -1, 0), nil
	case 4:
		return time.Now().AddDate(0, -3, 0), nil
	case 5:
		return time.Now().AddDate(-1, 0, 0), nil
	default:
		return time.Unix(0, 0), E.Message(fmt.Sprintf("unrecognized parameter type '%d'", clearType))
	}
}
