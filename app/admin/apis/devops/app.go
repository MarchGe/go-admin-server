package devops

import (
	"errors"
	"fmt"
	_ "github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	_ "github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

var _appApi = &AppApi{
	appService: dvservice.GetAppService(),
}

type AppApi struct {
	appService *dvservice.AppService
}

func GetAppApi() *AppApi {
	return _appApi
}

// AddApp godoc
//
//	@Summary	添加应用
//	@Tags		应用管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		[body]	body		req.AppUpsertReq	true	"应用信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/app [post]
func (a *AppApi) AddApp(c *gin.Context) {
	appUpsertReq := &req.AppUpsertReq{}
	if err := c.ShouldBindJSON(appUpsertReq); err != nil {
		E.PanicErr(err)
	}
	err := a.appService.CreateApp(appUpsertReq)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// UpdateApp godoc
//
//	@Summary	更新应用
//	@Tags		应用管理
//	@Accept		application/json
//	@Produce	application/json
//	@Param		id		path		int64				true	"应用ID"
//	@Param		[body]	body		req.AppUpsertReq	true	"应用信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/app/:id [put]
func (a *AppApi) UpdateApp(c *gin.Context) {
	upsertReq := &req.AppUpsertReq{}
	if e := c.ShouldBindJSON(upsertReq); e != nil {
		E.PanicErr(e)
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err = a.appService.UpdateApp(id, upsertReq); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// DeleteApp godoc
//
//	@Summary	删除应用
//	@Tags		应用管理
//	@Produce	application/json
//	@Param		id	path		int64	true	"应用ID"
//	@Success	200	{object}	R.Result
//	@Router		/devops/app/:id [delete]
func (a *AppApi) DeleteApp(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if err = a.appService.DeleteApp(id); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetList godoc
//
//	@Summary	查询应用列表
//	@Tags		应用管理
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照名称模糊搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[dvmodel.App]}
//	@Router		/devops/app/list [get]
func (a *AppApi) GetList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	if err := errors.Join(err1, err2); err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	pageableApps, err := a.appService.PageList(keyword, page, pageSize)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, pageableApps)
}

// UploadPkg godoc
//
//	@Summary	上传部署包
//	@Tags		应用管理
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		file	formData		file	true	"文件信息"
//	@Success	200		{object}	R.Result{value=res.UploadRes}
//	@Router		/devops/app/upload [post]
func (a *AppApi) UploadPkg(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		E.PanicErr(err)
	}
	if len([]rune(file.Filename)) > req.AppPkgFileNameMaxLength {
		R.Fail(c, "文件名太长", http.StatusBadRequest)
		return
	}
	cfg := config.GetConfig()
	slog.Debug("upload file", slog.String("size", fmt.Sprintf("%d bytes", file.Size)))
	if file.Size > int64(cfg.UploadPkgSizeLimit)*1024*1024 {
		R.Fail(c, fmt.Sprintf("文件大小不能超过%dMB", cfg.UploadPkgSizeLimit), http.StatusBadRequest)
		return
	}
	uploadTmpDir := a.appService.GetUploadTmpDir(cfg.WorkDir)
	userId := c.GetInt64(constant.SessionUserId)
	tmpKey := a.getPkgTmpKey(userId, file.Filename)
	if err = c.SaveUploadedFile(file, uploadTmpDir+"/"+tmpKey); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, &res.UploadRes{
		TmpKey:   tmpKey,
		FileName: file.Filename,
		FileSize: file.Size,
	})
}

func (a *AppApi) getPkgTmpKey(userId int64, fileName string) string {
	return fmt.Sprintf("tmp-%s%d-%s", time.Now().Format("20060102150405.999"), userId, fileName)
}

// DownloadPkg godoc
//
//	@Summary	下载部署包
//	@Tags		应用管理
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		key			query		string	true	"文件Key路径"
//	@Param		fileName	query		string	false	"文件名"
//	@Success	200		{object}	R.Result
//	@Router		/devops/app/download [get]
func (a *AppApi) DownloadPkg(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		R.Fail(c, "key不能为空", http.StatusBadRequest)
		return
	}
	if strings.Contains(key, "..") {
		R.Fail(c, "key值有误", http.StatusBadRequest)
		return
	}
	fileName := c.Query("fileName")
	if fileName == "" {
		index := strings.LastIndex(key, "/")
		fileName = key[index+1:]
	}
	uploadRoot := path.Clean(config.GetConfig().UploadPkgPath)
	filePath := uploadRoot + "/" + key
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.File(filePath)
}
