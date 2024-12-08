package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/task/log"
	"github.com/MarchGe/go-admin-server/app/admin/service/message"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/constant/sse"
	"github.com/MarchGe/go-admin-server/app/common/database"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var _deployTaskService = &DeployTaskService{
	appService:       dvservice.GetAppService(),
	scriptService:    dvservice.GetScriptService(),
	hostGroupService: dvservice.GetGroupService(),
	sshService:       dvservice.GetSshService(),
	logManager:       log.NewManager("data/log/task"),
}
var _ ConcreteTask = (*DeployTaskService)(nil)

type DeployTaskService struct {
	appService       *dvservice.AppService
	scriptService    *dvservice.ScriptService
	hostGroupService *dvservice.GroupService
	sshService       *dvservice.SshService
	logManager       *log.Manager
}

func (s *DeployTaskService) Create(tx *gorm.DB, data map[string]any) (int64, error) {
	info, err := s.validate(data)
	if err != nil {
		return 0, err
	}
	t := s.toModel(info)
	if err = tx.Omit("Task").Save(t).Error; err != nil {
		return 0, err
	}
	return t.Id, nil
}

func (s *DeployTaskService) Update(tx *gorm.DB, data map[string]any, id int64) error {
	info, err := s.validate(data)
	if err != nil {
		return err
	}
	deployTask, err := s.FindOneById(id)
	if err != nil {
		return err
	}
	s.copyProperties(info, deployTask)
	return tx.Omit("Task").Save(deployTask).Error
}

func (s *DeployTaskService) validate(concrete map[string]any) (*req.DeployTaskUpsertReq, error) {
	upsertReq := &req.DeployTaskUpsertReq{}
	if _, ok := concrete["uploadPath"].(string); ok {
		upsertReq.UploadPath = concrete["uploadPath"].(string)
	} else {
		upsertReq.UploadPath = ""
	}
	if _, ok := concrete["appId"].(float64); ok {
		upsertReq.AppId = int64(concrete["appId"].(float64))
	} else {
		upsertReq.AppId = 0
	}
	if _, ok := concrete["scriptId"].(float64); ok {
		upsertReq.ScriptId = int64(concrete["scriptId"].(float64))
	} else {
		upsertReq.ScriptId = 0
	}
	if _, ok := concrete["hostGroupId"].(float64); ok {
		upsertReq.HostGroupId = int64(concrete["hostGroupId"].(float64))
	} else {
		upsertReq.HostGroupId = 0
	}
	validate := ginUtils.GetValidator()
	if err := validate.Struct(upsertReq); err != nil {
		return nil, err
	}
	return upsertReq, nil
}

func (s *DeployTaskService) FindOneById(id int64) (*task.DeployTask, error) {
	t := &task.DeployTask{}
	if err := database.GetMysql().Where("id = ?", id).First(t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("具体任务不存在")
		} else {
			return nil, err
		}
	}
	return t, nil
}

func (s *DeployTaskService) Delete(tx *gorm.DB, t *task.Task) error {
	if err := tx.Delete(&task.DeployTask{}, t.AssociationID).Error; err != nil {
		return err
	}
	go s.logManager.RemoveLogs(t.Id)
	return nil
}

func (s *DeployTaskService) Start(ctx context.Context, t *task.Task) error {
	if err := s.setTaskStatus(t.Id, task.StatusRunning); err != nil {
		return err
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("batch deploy error", slog.Any("err", r), slog.String("err stack", string(debug.Stack())))
			}
		}()
		uId := ctx.Value("uId").(int64)
		if err := s.Run(ctx, t); err != nil {
			s.handleDeployError(t.Id, uId, err)
			return
		}
		if err := s.setTaskStatus(t.Id, task.StatusComplete); err != nil {
			slog.Error("change task status error", slog.Any("err", err))
		}
		_ = message.GetSseService().PushEventMessage(uId, sse.TaskExecuteEndEvent, "部署结束")
	}()
	return nil
}

func (s *DeployTaskService) Run(ctx context.Context, t *task.Task) error {
	deployTask, err := s.FindOneById(t.AssociationID)
	if err != nil {
		return err
	}
	app, script, group, err := s.getRelationData(deployTask)
	if err != nil {
		return err
	}
	return s.batchDeploy(ctx, t, path.Clean(deployTask.UploadPath), app, group.HostList, script)
}

func (s *DeployTaskService) Stop(ctx context.Context, t *task.Task) error {
	return s.setTaskStatus(t.Id, task.StatusStopped)
}

func (s *DeployTaskService) toModel(info *req.DeployTaskUpsertReq) *task.DeployTask {
	return &task.DeployTask{
		UploadPath:  info.UploadPath,
		AppId:       info.AppId,
		ScriptId:    info.ScriptId,
		HostGroupId: info.HostGroupId,
	}
}

func (s *DeployTaskService) copyProperties(info *req.DeployTaskUpsertReq, task *task.DeployTask) {
	task.UploadPath = info.UploadPath
	task.AppId = info.AppId
	task.ScriptId = info.ScriptId
	task.HostGroupId = info.HostGroupId
}

func (s *DeployTaskService) getRelationData(deployTask *task.DeployTask) (*dvmodel.App, *dvmodel.Script, *dvmodel.Group, error) {
	app, _ := s.appService.FindOneById(deployTask.AppId)
	if app == nil {
		return nil, nil, nil, E.Message("关联的应用不存在")
	}
	script, _ := s.scriptService.FindOneById(deployTask.ScriptId)
	if script == nil {
		return app, nil, nil, E.Message("关联的脚本不存在")
	}
	group, _ := s.hostGroupService.FindOneById(deployTask.HostGroupId, "HostList")
	if group == nil {
		return app, script, nil, E.Message("关联的分组不存在")
	}
	if len(group.HostList) == 0 {
		return app, script, group, E.Message("关联分组下的服务器列表为空")
	}
	return app, script, group, nil
}

func (s *DeployTaskService) batchDeploy(ctx context.Context, t *task.Task, remoteRoot string, app *dvmodel.App, hostList []*dvmodel.Host, script *dvmodel.Script) error {
	localRoot := path.Clean(config.GetConfig().UploadPkgPath)
	localPath := localRoot + "/" + app.Key
	manifestLogger, err := s.logManager.CreateManifestLogger(t.Id)
	if err != nil {
		return err
	}
	defer func() {
		manifestLogger.WriteEnd()
		manifestLogger.Close()
		s.logManager.RemoveOldLogs(t.Id, 1) // 仅保留最新一次执行日志
	}()
	for i := range hostList {
		ip := hostList[i].Ip
		hostLogger, e := s.logManager.CreateHostLogger(t.Id, ip)
		if e != nil {
			return e
		}
		manifestLogger.Append(log.NewEntry(i, ip, hostLogger.GetName(), "正在部署..."))
		if err = s.deploy(hostLogger, localPath, remoteRoot, app, hostList[i], script); err != nil {
			manifestLogger.Append(log.NewEntry(i, ip, hostLogger.GetName(), "失败"))
			hostLogger.Append(err.Error())
			slog.Error("deploy to host error", slog.String("host", ip), slog.Any("err", err))
		} else {
			manifestLogger.Append(log.NewEntry(i, ip, hostLogger.GetName(), "完成"))
		}
		hostLogger.WriteEnd()
		hostLogger.Close()
	}
	return nil
}

// 这些key会作为变量传到部署脚本中
const (
	envAppName    = "appName"    // 应用名称
	envAppVersion = "appVersion" // 应用版本
	envAppPort    = "appPort"    // 应用端口
	envPkgPath    = "pkgPath"    // 部署包的完整路径
	envPkgName    = "pkgName"    // 部署包的文件名
)

func (s *DeployTaskService) deploy(hostLogger *log.HostLogger, localPath string, remoteRoot string, app *dvmodel.App, host *dvmodel.Host, script *dvmodel.Script) error {
	hostLogger.Append(fmt.Sprintf("Ssh connecting to host %s:%d", host.Ip, host.Port))
	sshClient, err := s.sshService.CreateSshClient(host)
	if err != nil {
		return err
	}
	defer func() { _ = sshClient.Close() }()
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("new sftp client error, %w", err)
	}
	defer func() { _ = sftpClient.Close() }()
	remotePath := remoteRoot + "/" + app.FileName
	hostLogger.Append("Uploading deployment package...")
	if err = s.uploadFile(sftpClient, localPath, remotePath); err != nil {
		return err
	}
	hostLogger.Append("Deployment package upload completed.")

	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("new ssh session error, %w", err)
	}
	defer func() { _ = session.Close() }()

	commands := []string{
		"cd " + remoteRoot,
		envAppName + "=" + app.Name,
		envAppVersion + "=" + app.Version,
		envAppPort + "=" + strconv.Itoa(int(app.Port)),
		envPkgPath + "=" + remotePath,
		envPkgName + "=" + app.FileName,
		script.Content,
	}
	cmdContent := strings.Join(commands, "\n")
	executeTimeout := config.GetConfig().ScriptExecuteTimeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(executeTimeout)*time.Second)
	defer cancel()
	go func() {
		<-ctx.Done()
		if err = ctx.Err(); err != nil && errors.Is(err, context.DeadlineExceeded) {
			_ = session.Close()
		}
	}()
	hostLogger.Append("Executing script: " + constant.NewLine + "```" + constant.NewLine + script.Content + constant.NewLine + "```")
	session.Stdout = hostLogger.Original()
	session.Stderr = hostLogger.Original()
	if err = session.Run(cmdContent); err != nil {
		if exitError, ok := err.(*ssh.ExitError); ok && exitError.Signal() == string(ssh.SIGPIPE) {
			return fmt.Errorf("ssh command execution failed, maybe caused by a timeout, %w", err)
		}
		return fmt.Errorf("ssh execute command error, %w", err)
	}
	hostLogger.Append("Deploy completed!")
	return nil
}

func (s *DeployTaskService) uploadFile(sftpClient *sftp.Client, localPath, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("app package %s not found", localPath)
		} else {
			return fmt.Errorf("open local file %s error, %w", localPath, err)
		}
	}
	defer func() { _ = localFile.Close() }()
	remoteDir := path.Dir(remotePath)
	if err = sftpClient.MkdirAll(remoteDir); err != nil {
		return fmt.Errorf("sftp create remote dir %s error, %w", remoteDir, err)
	}
	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("sftp create remote file %s error, %w", remotePath, err)
	}
	defer func() { _ = remoteFile.Close() }()
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("sftp upload file error, %w", err)
	}
	return nil
}

func (s *DeployTaskService) handleDeployError(taskId, uId int64, err error) {
	slog.Error("deploy error", slog.Any("err", err))
	if e := s.setTaskStatus(taskId, task.StatusStopped); e != nil {
		slog.Error("change task status error", slog.Any("err", e))
	}
	_ = message.GetSseService().PushEventMessage(uId, sse.TaskExecuteFailEvent, err.Error())
}

func (s *DeployTaskService) setTaskStatus(taskId int64, status task.Status) error {
	return database.GetMysql().Model(&task.Task{}).Where("id = ?", taskId).Update("status", status).Error
}
