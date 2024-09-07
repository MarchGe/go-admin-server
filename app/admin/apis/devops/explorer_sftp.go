package devops

import (
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

var _explorerSftpApi = &ExplorerSftpApi{
	hostService:         dvservice.GetHostService(),
	explorerSftpService: dvservice.GetExplorerSftpService(),
}

type ExplorerSftpApi struct {
	hostService         *dvservice.HostService
	explorerSftpService *dvservice.ExplorerSftpService
}

func GetExplorerSftpApi() *ExplorerSftpApi {
	return _explorerSftpApi
}

// GetEntries godoc
//
//	@Summary	查询entry列表
//	@Tags		资源管理器（SFTP）
//	@Produce	application/json
//	@Param		dir		query		string	true	"目录路径"
//	@Param		hostId	query		int64	true	"主机的主键ID"
//	@Success	200		{object}	R.Result{value=[]res.ExplorerEntry}
//	@Router		/devops/explorer/sftp/entries [get]
func (a *ExplorerSftpApi) GetEntries(c *gin.Context) {
	dir := ginUtils.GetStringQuery(c, "dir", "")
	if dir == "" {
		R.Fail(c, "目录路径不能为空", http.StatusBadRequest)
		return
	}
	hostId, err := ginUtils.GetInt64Query(c, "hostId", 0)
	if err != nil {
		E.PanicErr(err)
	}
	host := a.getHost(hostId)
	entries, err := a.explorerSftpService.ListEntries(dir, host)
	if err != nil {
		E.PanicErr(err)
	}
	sortEntries(entries)
	R.Success(c, entries)
}

// DeleteEntry godoc
//
//	@Summary	删除文件或文件夹
//	@Tags		资源管理器（SFTP）
//	@Produce	application/json
//	@Param		path	query		string	true	"文件或文件夹的路径"
//	@Param		hostId	query		int64	true	"主机的主键ID"
//	@Success	200		{object}	R.Result
//	@Router		/devops/explorer/sftp/entry [delete]
func (a *ExplorerSftpApi) DeleteEntry(c *gin.Context) {
	deletePath := ginUtils.GetStringQuery(c, "path", "")
	if deletePath == "" {
		R.Fail(c, "操作的资源路径不能为空", http.StatusBadRequest)
		return
	}
	hostId, err := ginUtils.GetInt64Query(c, "hostId", 0)
	if err != nil {
		E.PanicErr(err)
	}
	host := a.getHost(hostId)
	if err = a.explorerSftpService.DeleteEntry(deletePath, host); err != nil {
		if errors.Is(err, os.ErrPermission) {
			R.Fail(c, "文件系统：permission denied", http.StatusBadRequest)
			return
		}
		if errors.Is(err, os.ErrNotExist) {
			R.Fail(c, "文件不存在", http.StatusBadRequest)
			return
		}
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// Upload godoc
//
//	@Summary	上传文件
//	@Tags		资源管理器（SFTP）
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		dir		formData	string	true	"文件目录"
//	@Param		file	formData	file	true	"文件信息"
//	@Param		hostId	formData	int64	true	"主机的主键ID"
//	@Success	200		{object}	R.Result
//	@Router		/devops/explorer/sftp/upload [post]
func (a *ExplorerSftpApi) Upload(c *gin.Context) {
	dir := c.PostForm("dir")
	if dir == "" {
		R.Fail(c, "目录参数不能为空", http.StatusBadRequest)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		E.PanicErr(err)
	}
	sHostId := c.PostForm("hostId")
	hostId, err := strconv.Atoi(sHostId)
	if err != nil {
		E.PanicErr(err)
	}
	host := a.getHost(int64(hostId))
	filePath := path.Clean(dir) + "/" + file.Filename
	f, err := file.Open()
	if err != nil {
		E.PanicErr(err)
	}
	defer func() { _ = f.Close() }()
	if err = a.explorerSftpService.UploadFile(filePath, f, host); err != nil {
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
//	@Tags		资源管理器（SFTP）
//	@Accept		multipart/form-data
//	@Produce	application/json
//	@Param		path	query		string	true	"文件完整路径"
//	@Param		hostId	query		int64	true	"主机的主键ID"
//	@Success	200		{object}	R.Result
//	@Router		/devops/explorer/sftp/download [get]
func (a *ExplorerSftpApi) Download(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		R.Fail(c, "文件路径不能为空", http.StatusBadRequest)
		return
	}
	hostId, err := ginUtils.GetInt64Query(c, "hostId", 0)
	if err != nil {
		E.PanicErr(err)
	}
	host := a.getHost(hostId)
	parts := strings.Split(filePath, "/")
	fileName := parts[len(parts)-1]
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	file, err := a.explorerSftpService.DownloadFile(filePath, host)
	if err != nil {
		E.PanicErr(err)
	}
	defer func() { _ = file.Close() }()
	if _, err = io.Copy(c.Writer, file); err != nil {
		E.PanicErr(err)
	}
	c.Writer.Flush()
}

// CreateDir godoc
//
//	@Summary	创建目录
//	@Tags		资源管理器（SFTP）
//	@Accept		application/json
//	@Produce	application/json
//	@Param		hostId	body		int64	true	"主机主键ID"
//	@Param		dir		body		string	true	"当前目录"
//	@Param		name	body		string	true	"创建目录的名称"
//	@Success	200		{object}	R.Result
//	@Router		/devops/explorer/sftp/create [post]
func (a *ExplorerSftpApi) CreateDir(c *gin.Context) {
	var body req.SftpCreateDirReq
	if err := c.ShouldBindJSON(&body); err != nil {
		E.PanicErr(err)
	}

	host := a.getHost(body.HostId)
	if err := a.explorerSftpService.CreateDir(&body, host); err != nil {
		if errors.Is(err, os.ErrPermission) {
			R.Fail(c, "文件系统：permission denied", http.StatusBadRequest)
			return
		}
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// Rename godoc
//
//	@Summary	重命名
//	@Tags		资源管理器（SFTP）
//	@Accept		application/json
//	@Produce	application/json
//	@Param		hostId	body		int64	true	"主机主键ID"
//	@Param		dir		body		string	true	"当前目录"
//	@Param		oldName	body		string	true	"旧名称"
//	@Param		newName	body		string	true	"新名称"
//	@Success	200		{object}	R.Result
//	@Router		/devops/explorer/sftp/rename [post]
func (a *ExplorerSftpApi) Rename(c *gin.Context) {
	var body req.SftpRenameReq
	if err := c.ShouldBindJSON(&body); err != nil {
		E.PanicErr(err)
	}
	host := a.getHost(body.HostId)
	if err := a.explorerSftpService.Rename(&body, host); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

func (a *ExplorerSftpApi) getHost(hostId int64) *dvmodel.Host {
	if hostId == 0 {
		E.PanicErr(E.Message("主机ID参数不能为空"))
	}
	host, err := a.hostService.FindOneById(hostId)
	if err != nil {
		E.PanicErr(err)
	}
	return host
}
