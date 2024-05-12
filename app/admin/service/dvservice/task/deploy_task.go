package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/message"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/google/uuid"
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
	"sync"
	"time"
)

var _deployTaskService = &DeployTaskService{
	appService:       dvservice.GetAppService(),
	scriptService:    dvservice.GetScriptService(),
	hostGroupService: dvservice.GetGroupService(),
	taskLogService:   GetTaskLogService(),
}
var _ ConcreteTask = (*DeployTaskService)(nil)

type DeployTaskService struct {
	appService       *dvservice.AppService
	scriptService    *dvservice.ScriptService
	hostGroupService *dvservice.GroupService
	taskLogService   *LogService
}

var runningMap = sync.Map{}

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

func (s *DeployTaskService) FindOneById(id int64) (*dvmodel.DeployTask, error) {
	t := &dvmodel.DeployTask{}
	if err := database.GetMysql().Where("id = ?", id).First(t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, E.Message("具体任务不存在")
		} else {
			return nil, err
		}
	}
	return t, nil
}

func (s *DeployTaskService) Delete(tx *gorm.DB, t *dvmodel.Task) error {
	if err := tx.Delete(&dvmodel.DeployTask{}, t.AssociationID).Error; err != nil {
		return err
	}
	deploymentLogDir := s.taskLogService.getTaskLogDir(t.Id)
	go func() {
		_ = os.RemoveAll(deploymentLogDir)
	}()
	return nil
}

func (s *DeployTaskService) Start(ctx context.Context, t *dvmodel.Task) error {
	if err := s.setTaskStatus(t.Id, dvmodel.TaskStatusRunning); err != nil {
		return err
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("batch deploy error", slog.Any("err", r), slog.String("err stack", string(debug.Stack())))
			}
			runningMap.Delete(t.Id)
		}()
		uId := ctx.Value("uId").(int64)
		if err := s.Run(ctx, t); err != nil {
			s.handleDeployError(t.Id, uId, err)
		}
		if err := s.setTaskStatus(t.Id, dvmodel.TaskStatusComplete); err != nil {
			slog.Error("change task status error", slog.Any("err", err))
		}
		_ = message.GetSseService().PushEventMessage(uId, constant.SseTaskExecuteEndEvent, "部署结束")
	}()
	return nil
}

func (s *DeployTaskService) Run(ctx context.Context, t *dvmodel.Task) error {
	if _, ok := runningMap.Load(t.Id); ok {
		return errors.New("任务已经启动")
	}
	runningCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	runningMap.Store(t.Id, cancel)

	deployTask, err := s.FindOneById(t.AssociationID)
	if err != nil {
		return err
	}
	app, script, group, err := s.getRelationData(deployTask)
	if err != nil {
		return err
	}
	return s.batchDeploy(runningCtx, t, path.Clean(deployTask.UploadPath), app, group.HostList, script)
}

func (s *DeployTaskService) Stop(ctx context.Context, t *dvmodel.Task) error {
	value, loaded := runningMap.LoadAndDelete(t.Id)
	if loaded {
		cancel := value.(context.CancelFunc)
		cancel()
	}
	return s.setTaskStatus(t.Id, dvmodel.TaskStatusStopped)
}

func (s *DeployTaskService) toModel(info *req.DeployTaskUpsertReq) *dvmodel.DeployTask {
	return &dvmodel.DeployTask{
		UploadPath:  info.UploadPath,
		AppId:       info.AppId,
		ScriptId:    info.ScriptId,
		HostGroupId: info.HostGroupId,
	}
}

func (s *DeployTaskService) copyProperties(info *req.DeployTaskUpsertReq, task *dvmodel.DeployTask) {
	task.UploadPath = info.UploadPath
	task.AppId = info.AppId
	task.ScriptId = info.ScriptId
	task.HostGroupId = info.HostGroupId
}

func (s *DeployTaskService) getRelationData(deployTask *dvmodel.DeployTask) (*dvmodel.App, *dvmodel.Script, *dvmodel.Group, error) {
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

func (s *DeployTaskService) batchDeploy(ctx context.Context, t *dvmodel.Task, remoteRoot string, app *dvmodel.App, hostList []*dvmodel.Host, script *dvmodel.Script) error {
	localRoot := path.Clean(config.GetConfig().UploadPkgPath)
	localPath := localRoot + "/" + app.Key
	manifestFile, err := s.taskLogService.createManifestLogFile(t.Id)
	if err != nil {
		return err
	}
	defer func() {
		_ = manifestFile.Close()
		s.taskLogService.removeOldLogs(t.Id, 1) // 仅保留最新一次执行日志
	}()
	for i := range hostList {
		select {
		case <-ctx.Done():
			s.taskLogService.writeEnd(manifestFile)
			return errors.New("任务已停止")
		default:
			ip := hostList[i].Ip
			hostFileName := strings.ReplaceAll(uuid.NewString(), "-", "")
			hostFile, e := s.taskLogService.createHostLogFile(t.Id, ip, hostFileName)
			if e != nil {
				s.handleHostDeployError(manifestFile, hostFile, i, ip, hostFileName, e)
				continue
			}
			s.taskLogService.appendManifestLog(manifestFile, i, ip, hostFileName, "正在部署...")
			if err = s.deploy(hostFile, localPath, remoteRoot, app, hostList[i], script); err != nil {
				s.handleHostDeployError(manifestFile, hostFile, i, ip, hostFileName, err)
			} else {
				s.taskLogService.appendManifestLog(manifestFile, i, ip, hostFileName, "完成")
			}
			s.taskLogService.writeEnd(hostFile)
			_ = hostFile.Close()
		}
	}
	s.taskLogService.writeEnd(manifestFile)
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

func (s *DeployTaskService) deploy(hostFile *os.File, localPath string, remoteRoot string, app *dvmodel.App, host *dvmodel.Host, script *dvmodel.Script) error {
	s.taskLogService.appendHostLog(hostFile, fmt.Sprintf("ssh connecting to host %s:%d", host.Ip, host.Port))
	sshClient, err := s.createSshClient(host)
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
	s.taskLogService.appendHostLog(hostFile, "uploading deployment package...")
	if err = s.uploadFile(sftpClient, localPath, remotePath); err != nil {
		return err
	}
	s.taskLogService.appendHostLog(hostFile, "deployment package upload completed.")

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
	s.taskLogService.appendHostLog(hostFile, "executing script: "+constant.NewLine+"```"+constant.NewLine+script.Content+constant.NewLine+"```")
	session.Stdout = hostFile
	session.Stderr = hostFile
	if err = session.Run(cmdContent); err != nil {
		if _, ok := err.(*ssh.ExitError); ok {
			if exitError := err.(*ssh.ExitError); exitError.Signal() == string(ssh.SIGPIPE) {
				return fmt.Errorf("ssh execute remote command error, maybe execute timeout, %w", err)
			}
		}
		return fmt.Errorf("ssh execute remote command error, %w", err)
	}
	s.taskLogService.appendHostLog(hostFile, "deploy completed.")
	return nil
}

func (s *DeployTaskService) createSshClient(host *dvmodel.Host) (*ssh.Client, error) {
	encryptKey := config.GetConfig().EncryptKey
	decryptPasswd, err := utils.DecryptString(encryptKey, host.Password, "")
	if err != nil {
		return nil, fmt.Errorf("decrypt password error, %w", err)
	}
	clientConfig := ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(decryptPasswd),
		},
		Timeout:         constant.SshEstablishTimeoutInSeconds * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", host.Ip, host.Port)
	client, err := ssh.Dial("tcp", addr, &clientConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh connect failed, %w", err)
	}
	return client, nil
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
	if err = s.setTaskStatus(taskId, dvmodel.TaskStatusStopped); err != nil {
		slog.Error("change task status error", slog.Any("err", err))
	}
	_ = message.GetSseService().PushEventMessage(uId, constant.SseTaskExecuteFailEvent, err.Error())
}

func (s *DeployTaskService) handleHostDeployError(manifestFile *os.File, hostFile *os.File, index int, host, hostFileName string, err error) {
	s.taskLogService.appendManifestLog(manifestFile, index, host, hostFileName, "失败")
	s.taskLogService.appendHostLog(hostFile, err.Error()+constant.NewLine)
	slog.Error("deploy to host error", slog.String("host", host), slog.Any("err", err))
}

func (s *DeployTaskService) setTaskStatus(taskId int64, status dvmodel.TaskStatus) error {
	return database.GetMysql().Model(&dvmodel.Task{}).Where("id = ?", taskId).Update("status", status).Error
}
