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

var _scriptApi = &ScriptApi{
	scriptService: dvservice.GetScriptService(),
}

type ScriptApi struct {
	scriptService *dvservice.ScriptService
}

func GetScriptApi() *ScriptApi {
	return _scriptApi
}

// AddScript godoc
//
//	@Summary	添加脚本
//	@Tags		脚本管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.ScriptUpsertReq	true	"脚本信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/script [post]
func (a *ScriptApi) AddScript(c *gin.Context) {
	scriptUpsertReq := &req.ScriptUpsertReq{}
	if err := c.ShouldBindJSON(scriptUpsertReq); err != nil {
		E.PanicErr(err)
	}
	err := a.scriptService.CreateScript(scriptUpsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateScript godoc
//
//	@Summary	更新脚本
//	@Tags		脚本管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		path		int64				true	"脚本ID"
//	@Param		[body]	body		req.ScriptUpsertReq	true	"脚本信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/script/:id [put]
func (a *ScriptApi) UpdateScript(c *gin.Context) {
	upsertReq := &req.ScriptUpsertReq{}
	if e := c.ShouldBindJSON(upsertReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.scriptService.UpdateScript(id, upsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteScript godoc
//
//	@Summary	删除脚本
//	@Tags		脚本管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"脚本ID"
//	@Success	200	{object}	R.Result
//	@Router		/devops/script/:id [delete]
func (a *ScriptApi) DeleteScript(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.scriptService.DeleteScript(id)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetList godoc
//
//	@Summary	查询脚本列表
//	@Tags		脚本管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照名称模糊搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[dvmodel.Script]}
//	@Router		/devops/script/list [get]
func (a *ScriptApi) GetList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableScripts, err := a.scriptService.PageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableScripts)
}
