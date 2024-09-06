package devops

import (
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
)

var _explorerApi = &ExplorerApi{
	explorerService: dvservice.GetExplorerService(),
}

type ExplorerApi struct {
	explorerService *dvservice.ExplorerService
}

func GetExplorerApi() *ExplorerApi {
	return _explorerApi
}

// GetEntries godoc
//
//	@Summary	查询entry列表
//	@Tags		资源管理器
//	@Produce	application/json
//	@Param		dir	query		string	true	"目录路径"
//	@Success	200	{object}	R.Result{value=[]res.ExplorerEntry}
//	@Router		/devops/explorer/entries [get]
func (a *ExplorerApi) GetEntries(c *gin.Context) {
	dir := ginUtils.GetStringQuery(c, "dir", "")
	if dir == "" {
		R.Fail(c, "目录路径不能为空", http.StatusBadRequest)
		return
	}
	entries, err := a.explorerService.ListEntries(dir)
	if err != nil {
		E.PanicErr(err)
	}
	sortEntries(entries)
	R.Success(c, entries)
}

// DeleteEntry godoc
//
//	@Summary	删除文件或文件夹
//	@Tags		资源管理器
//	@Produce	application/json
//	@Param		path	query		string	true	"文件或文件夹的路径"
//	@Success	200		{object}	R.Result
//	@Router		/devops/explorer/entry [delete]
func (a *ExplorerApi) DeleteEntry(c *gin.Context) {
	deletePath := ginUtils.GetStringQuery(c, "path", "")
	if deletePath == "" {
		R.Fail(c, "操作的资源路径不能为空", http.StatusBadRequest)
		return
	}
	if err := os.RemoveAll(deletePath); err != nil {
		if errors.Is(err, os.ErrPermission) {
			R.Fail(c, "文件系统：permission denied", http.StatusBadRequest)
			return
		}
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// Upload godoc
//
//	@Summary	上传文件
//	@Tags		资源管理器
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		dir		formData	string	true	"文件目录"
//	@Param		file	formData	file	true	"文件信息"
//	@Success	200		{object}	R.Result
//	@Router		/devops/explorer/upload [post]
func (a *ExplorerApi) Upload(c *gin.Context) {
	dir := c.PostForm("dir")
	if dir == "" {
		R.Fail(c, "目录参数不能为空", http.StatusBadRequest)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		E.PanicErr(err)
	}
	filePath := path.Clean(dir) + "/" + file.Filename
	if err = c.SaveUploadedFile(file, filePath); err != nil {
		if errors.Is(err, os.ErrPermission) {
			R.Fail(c, "文件系统：permission denied", http.StatusBadRequest)
			return
		}
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// Download godoc
//
//	@Summary	下载文件
//	@Tags		资源管理器
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		path	query		string	true	"文件完整路径"
//	@Success	200		{object}	R.Result
//	@Router		/devops/explorer/download [get]
func (a *ExplorerApi) Download(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		R.Fail(c, "文件路径不能为空", http.StatusBadRequest)
		return
	}
	info, err := os.Stat(filePath)
	if err != nil {
		R.Fail(c, "获取文件信息失败", http.StatusBadRequest)
		return
	}
	if info.IsDir() {
		R.Fail(c, "不支持下载文件夹", http.StatusBadRequest)
		return
	}
	parts := strings.Split(filePath, "/")
	fileName := parts[len(parts)-1]
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	c.File(filePath)
}

// CreateDir godoc
//
//	@Summary	创建目录
//	@Tags		资源管理器（SFTP）
//	@Accept		application/json
//	@Produce	application/json
//	@Param		dir		body		string	true	"当前目录"
//	@Param		name	body		string	true	"创建目录的名称"
//	@Success	200		{object}	R.Result
//	@Router		/devops/explorer/create [post]
func (a *ExplorerApi) CreateDir(c *gin.Context) {
	var body req.ExplorerCreateDirReq
	if err := c.ShouldBindJSON(&body); err != nil {
		E.PanicErr(err)
	}
	dir := path.Clean(body.Dir + "/" + body.Name)
	_, err := os.Stat(dir)
	if !os.IsNotExist(err) {
		R.Fail(c, "目录已存在", http.StatusBadRequest)
		return
	}
	if err = os.MkdirAll(dir, 0750); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// sortEntries 排序规则：文件夹在前，然后按字母自然顺序排序（忽略大小写）
func sortEntries(entries []*res.ExplorerEntry) {
	length := len(entries)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Type == res.EntryTypeDir {
			return true
		}
		return false
	})
	dirNum := 0
	for i := 0; i < length; i++ {
		if entries[i].Type != res.EntryTypeDir {
			dirNum = i
			break
		}
	}
	dirEntries := entries[0:dirNum]
	sort.Slice(dirEntries, func(i, j int) bool {
		if strings.Compare(strings.ToLower(dirEntries[i].Name), strings.ToLower(dirEntries[j].Name)) <= 0 {
			return true
		}
		return false
	})
	nonDirEntries := entries[dirNum:]
	sort.Slice(nonDirEntries, func(i, j int) bool {
		if strings.Compare(strings.ToLower(nonDirEntries[i].Name), strings.ToLower(nonDirEntries[j].Name)) <= 0 {
			return true
		}
		return false
	})
}
