package state

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel"
	"github.com/MarchGe/go-admin-server/app/admin/model/dvmodel/task"
	"github.com/MarchGe/go-admin-server/app/admin/scheduler"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/req"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/task/log"
	"github.com/MarchGe/go-admin-server/app/admin/service/message"
	"github.com/MarchGe/go-admin-server/app/common"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/constant/sse"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"github.com/MarchGe/go-admin-server/config"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"log/slog"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"
)

var UnsupportedOperationErr = E.Message("当前状态不支持该操作")
var _taskStateCommon = &TaskStateCommon{
	sshService: dvservice.GetSshService(),
	logManager: log.NewManager("data/log/script-task"),
}

type TaskStateCommon struct {
	sshService *dvservice.SshService
	logManager *log.Manager
}

func GetTaskStateCommon() *TaskStateCommon {
	return _taskStateCommon
}

func (s *TaskStateCommon) Start(ctx context.Context, t *task.ScriptTask) error {
	if t.ExecuteType == task.ExecuteTypeAuto {
		err := scheduler.AddTask(t.Id, t.Cron, func() {
			if e := s.Run(ctx, t); e != nil {
				slog.Error("scheduler execute task error", slog.Int64("taskId", t.Id), slog.Any("err", e))
			}
		})
		if err != nil {
			s.updateStatus(t.Id, task.StatusStopped)
			return err
		}
		if err = s.updateStatus(t.Id, task.StatusActive); err != nil {
			return err
		}
	} else {
		if err := s.updateStatus(t.Id, task.StatusRunning); err != nil {
			return err
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("execute task error", slog.Int64("taskId", t.Id), slog.Any("err", r), slog.String("err stack", string(debug.Stack())))
				}
			}()
			uId := ctx.Value("uId").(int64)
			if e := s.Run(ctx, t); e != nil {
				slog.Error("Execute script task error", slog.Int64("taskId", t.Id), slog.Any("err", e))
				_ = message.GetSseService().PushEventMessage(uId, sse.ScriptExecuteFailEvent, e.Error())
			} else {
				_ = message.GetSseService().PushEventMessage(uId, sse.ScriptExecuteEndEvent, "任务执行结束")
			}
		}()
	}
	return nil
}

func (s *TaskStateCommon) Run(ctx context.Context, t *task.ScriptTask) (err error) {
	slog.Info("Executing script task", slog.Int64("id", t.Id), slog.String("name", t.Name))
	defer func() {
		if t.ExecuteType == task.ExecuteTypeAuto {
			_ = s.updateStatusIfNotStopped(t.Id, task.StatusActive)
		} else {
			if err != nil {
				_ = s.updateStatus(t.Id, task.StatusStopped)
			} else {
				_ = s.updateStatus(t.Id, task.StatusComplete)
			}
		}
	}()
	if t.Kind == task.KindLocal {
		err = s.executeLocalTask(ctx, t)
	} else {
		err = s.executeRemoteTask(ctx, t)
	}
	return err
}

func (s *TaskStateCommon) Stop(ctx context.Context, t *task.ScriptTask) error {
	if err := s.updateStatus(t.Id, task.StatusStopped); err != nil {
		return err
	}
	if t.ExecuteType == task.ExecuteTypeAuto {
		scheduler.RemoveTask(t.Id)
	}
	return nil
}

func (s *TaskStateCommon) Update(ctx context.Context, info *req.ScriptTaskUpsertReq, t *task.ScriptTask) error {
	s.copyProperties(info, t)
	t.UpdateTime = time.Now()
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if e := tx.Save(t).Error; e != nil {
			return e
		}
		if e := tx.Where("task_id = ?", t.Id).Delete(&task.ScriptTaskScript{}).Error; e != nil {
			return e
		}
		scripts := make([]*task.ScriptTaskScript, len(info.ScriptIds))
		for i := range info.ScriptIds {
			scripts[i] = &task.ScriptTaskScript{
				TaskId:   t.Id,
				ScriptId: info.ScriptIds[i],
			}
		}
		return tx.Save(scripts).Error
	})
	return err
}

func (s *TaskStateCommon) copyProperties(info *req.ScriptTaskUpsertReq, t *task.ScriptTask) {
	t.Name = info.Name
	t.Kind = info.Kind
	t.Cron = info.Cron
	if strings.TrimSpace(info.Cron) == "" {
		t.ExecuteType = task.ExecuteTypeManual
	} else {
		t.ExecuteType = task.ExecuteTypeAuto
	}
	if info.Kind == task.KindRemote {
		t.HostGroupId = info.HostGroupId
	} else {
		t.HostGroupId = 0
	}
}

func (s *TaskStateCommon) Delete(ctx context.Context, t *task.ScriptTask) error {
	err := database.GetMysql().Transaction(func(tx *gorm.DB) error {
		if e := tx.Where("task_id = ?", t.Id).Delete(&task.ScriptTaskScript{}).Error; e != nil {
			return e
		}
		return tx.Delete(&task.ScriptTask{}, t.Id).Error
	})
	return err
}

func (s *TaskStateCommon) updateStatus(id int64, status task.Status) error {
	err := database.GetMysql().Model(&task.ScriptTask{}).Where("id = ?", id).
		UpdateColumn("status", status).Error
	return err
}

func (s *TaskStateCommon) updateStatusIfNotStopped(id int64, status task.Status) error {
	err := database.GetMysql().Model(&task.ScriptTask{}).Where("id = ? and status != ?", id, task.StatusStopped).
		UpdateColumn("status", status).Error
	return err
}

func (s *TaskStateCommon) executeLocalTask(ctx context.Context, t *task.ScriptTask) error {
	manifestLogger, e := s.logManager.CreateManifestLogger(t.Id)
	if e != nil {
		return e
	}
	defer func() {
		manifestLogger.WriteEnd()
		manifestLogger.Close()
		s.logManager.RemoveOldLogs(t.Id, 1)
	}()
	host := "127.0.0.1"
	hostLogger, err := s.logManager.CreateHostLogger(t.Id, host)
	if err != nil {
		return err
	}
	defer func() {
		hostLogger.WriteEnd()
		hostLogger.Close()
	}()
	manifestLogger.Append(log.NewEntry(0, host, hostLogger.GetName(), "正在执行..."))
	if err = s.executeLocalScript(ctx, t, hostLogger); err != nil {
		manifestLogger.Append(log.NewEntry(0, host, hostLogger.GetName(), "失败"))
		hostLogger.Append(err.Error())
		return err
	} else {
		manifestLogger.Append(log.NewEntry(0, host, hostLogger.GetName(), "完成"))
	}
	return nil
}

func (s *TaskStateCommon) executeRemoteTask(ctx context.Context, t *task.ScriptTask) error {
	manifestLogger, e := s.logManager.CreateManifestLogger(t.Id)
	if e != nil {
		return e
	}
	defer func() {
		manifestLogger.WriteEnd()
		manifestLogger.Close()
		s.logManager.RemoveOldLogs(t.Id, 1)
	}()
	hostList := t.HostGroup.HostList
	for i, host := range hostList {
		hostLogger, err := s.logManager.CreateHostLogger(t.Id, host.Ip)
		if err != nil {
			return err
		}
		manifestLogger.Append(log.NewEntry(i, host.Ip, hostLogger.GetName(), "正在执行..."))
		if err = s.executeRemoteScript(ctx, host, t, hostLogger); err != nil {
			manifestLogger.Append(log.NewEntry(i, host.Ip, hostLogger.GetName(), "失败"))
			hostLogger.Append(err.Error())
			slog.Error("execute remote script error", slog.String("host", host.Ip), slog.Any("err", err))
		} else {
			manifestLogger.Append(log.NewEntry(i, host.Ip, hostLogger.GetName(), "完成"))
		}
		hostLogger.WriteEnd()
		hostLogger.Close()
	}
	return nil
}

func (s *TaskStateCommon) executeLocalScript(ctx context.Context, t *task.ScriptTask, hostLogger *log.HostLogger) error {
	bash, err := common.GetBash()
	if err != nil {
		return err
	}
	scripts := t.Scripts
	for _, script := range scripts {
		cmd := strings.Builder{}
		cmd.WriteString("set -e\n")
		cmd.WriteString(fmt.Sprintf("echo '>>> Executing script: %s-%s'\n", script.Name, script.Version))
		cmd.WriteString(script.Content + "\n")
		command := exec.Command(bash, "-c", cmd.String())
		command.Stdout = hostLogger.Original()
		command.Stderr = hostLogger.Original()
		if e := command.Run(); e != nil {
			return fmt.Errorf("execute script error, %w", e)
		}
	}
	return nil
}
func (s *TaskStateCommon) executeRemoteScript(ctx context.Context, host *dvmodel.Host, t *task.ScriptTask, hostLogger *log.HostLogger) error {
	sshClient, err := s.sshService.CreateSshClient(host)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	scripts := t.Scripts
	for _, script := range scripts {
		session, e := sshClient.NewSession()
		if e != nil {
			return fmt.Errorf("new ssh session error, %w", e)
		}
		session.Stdout = hostLogger.Original()
		session.Stderr = hostLogger.Original()
		cmd := strings.Builder{}
		cmd.WriteString("set -e\n")
		cmd.WriteString(fmt.Sprintf("echo '>>> Executing script: %s-%s'\n", script.Name, script.Version))
		cmd.WriteString(script.Content + "\n")

		scriptExecuteTimeout := config.GetConfig().ScriptExecuteTimeout
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(scriptExecuteTimeout)*time.Second)
		go func() {
			<-timeoutCtx.Done()
			if err = timeoutCtx.Err(); err != nil && errors.Is(err, context.DeadlineExceeded) {
				session.Close()
			}
		}()
		if e = session.Run(cmd.String()); e != nil {
			cancel()
			if exitError, ok := e.(*ssh.ExitError); ok && exitError.Signal() == string(ssh.SIGPIPE) {
				return fmt.Errorf("ssh command execution failed, maybe caused by a timeout, %w", err)
			}
			session.Close()
			return fmt.Errorf("ssh execute command error, %w", e)
		}
		cancel()
		session.Close()
	}
	return nil
}
