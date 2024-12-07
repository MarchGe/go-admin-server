package apis

import (
	"context"
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/task/log"
	"github.com/MarchGe/go-admin-server/app/admin/service/message"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/constant/sse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"
)

var _sseApi = &SseApi{
	sseService:           &message.SseService{},
	taskLogManager:       log.NewManager("data/log/task"),
	scriptTaskLogManager: log.NewManager("data/log/script-task"),
}

type SseApi struct {
	sseService           *message.SseService
	taskLogManager       *log.Manager
	scriptTaskLogManager *log.Manager
}

type LogManager interface {
	OpenLatestManifestLogger(taskId int64) (*log.ManifestLogger, error)
	OpenHostLogger(taskId int64, host string, hostLogName string) (*log.HostLogger, error)
}

func GetSseApi() *SseApi {
	return _sseApi
}

// MessagePush godoc
//
//	@Summary	服务端事件消息推送
//	@Tags		消息服务
//	@Accept		application/json
//	@Produce	text/event-stream
//	@Router		/sse/message-push [get]
func (a *SseApi) MessagePush(c *gin.Context) {
	w := c.Writer
	a.setSSEHeader(c)
	userId := c.GetInt64(constant.SessionUserId)

	sId := uuid.NewString()
	mtx := &sync.Mutex{}
	a.sseService.AddSseSession(sId, userId, w, mtx)
	defer a.sseService.RemoveSseSession(sId)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		notify := w.CloseNotify()
		<-notify
		cancel()
	}()
	a.startHeartbeat(ctx, c, mtx)
}

// PushManifestLogEvent godoc
//
//	@Summary	推送manifest日志
//	@Tags		任务管理
//	@Accept		application/json
//	@Produce	text/event-stream
//	@Param		id	path	int64	true	"任务ID"
//	@Router		/sse/task/:id/manifest-log [get]
func (a *SseApi) PushManifestLogEvent(c *gin.Context) {
	a.pushManifestLog(c, a.taskLogManager)
}

// PushScriptTaskManifestLog godoc
//
//	@Summary	推送manifest日志
//	@Tags		任务管理
//	@Accept		application/json
//	@Produce	text/event-stream
//	@Param		id	path	int64	true	"任务ID"
//	@Router		/sse/script-task/:id/manifest-log [get]
func (a *SseApi) PushScriptTaskManifestLog(c *gin.Context) {
	a.pushManifestLog(c, a.scriptTaskLogManager)
}

func (a *SseApi) pushManifestLog(c *gin.Context, logManager LogManager) {
	mtx := &sync.Mutex{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		a.sendEvent(c, mtx, sse.Error, err.Error())
		return
	}
	w := c.Writer
	a.setSSEHeader(c)

	manifestLogger, err := logManager.OpenLatestManifestLogger(id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			a.sendEvent(c, mtx, sse.Error, "所请求的日志已被清理")
			return
		}
		a.sendEvent(c, mtx, sse.Error, err.Error())
		return
	}
	defer manifestLogger.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var connClosed bool
	go func() {
		notify := w.CloseNotify()
		<-notify
		connClosed = true
		cancel()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("heartbeat error", slog.Any("err", r))
			}
		}()
		a.startHeartbeat(ctx, c, mtx)
	}()
	reader := manifestLogger.GetReader()
	for {
		if connClosed {
			break
		}
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			a.sendEvent(c, mtx, sse.Error, err.Error())
			return
		}
		if line != "" {
			if log.IsEnd(line) {
				a.sendEvent(c, mtx, sse.Close, "")
				time.Sleep(2 * time.Second)
				return
			}
			entry, err := log.ParseEntry(line)
			if err != nil {
				a.sendEvent(c, mtx, sse.Error, err.Error())
				return
			}
			a.sendEvent(c, mtx, sse.ManifestEntryEvent, entry)
		} else {
			time.Sleep(time.Second)
		}
	}
}

// PushHostLogEvent godoc
//
//	@Summary	推送远程主机日志
//	@Tags		任务管理
//	@Accept		application/json
//	@Produce	text/event-stream
//	@Param		id			path	int64	true	"任务ID"
//	@Param		host		query	string	true	"主机IP"
//	@Param		hostLogName	query	string	true	"主机日志的文件名"
//	@Router		/sse/task/:id/host-log [get]
func (a *SseApi) PushHostLogEvent(c *gin.Context) {
	a.pushHostLog(c, a.taskLogManager)
}

// PushScriptTaskHostLogEvent godoc
//
//	@Summary	推送远程主机日志
//	@Tags		任务管理
//	@Accept		application/json
//	@Produce	text/event-stream
//	@Param		id			path	int64	true	"任务ID"
//	@Param		host		query	string	true	"主机IP"
//	@Param		hostLogName	query	string	true	"主机日志的文件名"
//	@Router		/sse/script-task/:id/host-log [get]
func (a *SseApi) PushScriptTaskHostLogEvent(c *gin.Context) {
	a.pushHostLog(c, a.scriptTaskLogManager)
}
func (a *SseApi) pushHostLog(c *gin.Context, logManager LogManager) {
	mtx := &sync.Mutex{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		a.sendEvent(c, mtx, sse.Error, err.Error())
		return
	}
	host := c.Query("host")
	hostLogName := c.Query("hostLogName")
	if host == "" || hostLogName == "" {
		a.sendEvent(c, mtx, sse.Error, "缺少参数")
		return
	}
	w := c.Writer
	a.setSSEHeader(c)

	hostLogger, err := logManager.OpenHostLogger(id, host, hostLogName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			a.sendEvent(c, mtx, sse.Error, "所请求的日志已被清理")
			return
		}
		a.sendEvent(c, mtx, sse.Error, err.Error())
		return
	}
	defer hostLogger.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var connClosed bool
	go func() {
		notify := w.CloseNotify()
		<-notify
		connClosed = true
		cancel()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("heartbeat error", slog.Any("err", r))
			}
		}()
		a.startHeartbeat(ctx, c, mtx)
	}()
	reader := hostLogger.GetReader()
	for {
		if connClosed {
			break
		}
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			a.sendEvent(c, mtx, sse.Error, err.Error())
			return
		}
		if line != "" {
			if log.IsEnd(line) {
				a.sendEvent(c, mtx, sse.Close, "")
				time.Sleep(2 * time.Second)
				return
			}
			a.sendEvent(c, mtx, sse.HostLogEvent, line)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (a *SseApi) sendEvent(c *gin.Context, mtx *sync.Mutex, event sse.Event, msg any) {
	mtx.Lock()
	defer mtx.Unlock()
	c.SSEvent(event.String(), msg)
	c.Writer.Flush()
}

func (a *SseApi) startHeartbeat(ctx context.Context, c *gin.Context, mtx *sync.Mutex) {
	mtx.Lock()
	c.SSEvent(sse.Tick.String(), "")
	c.Writer.Flush()
	mtx.Unlock()
	ticker := time.NewTicker(sse.TickDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mtx.Lock()
			c.SSEvent(sse.Tick.String(), "")
			c.Writer.Flush()
			mtx.Unlock()
		}
	}
}

func (a *SseApi) setSSEHeader(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
}
