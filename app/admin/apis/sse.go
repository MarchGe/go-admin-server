package apis

import (
	"bufio"
	"context"
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice/task"
	"github.com/MarchGe/go-admin-server/app/admin/service/message"
	"github.com/MarchGe/go-admin-server/app/common/constant"
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
	sseService: &message.SseService{},
}

type SseApi struct {
	sseService *message.SseService
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
//	@Param		id		path		int64	true	"任务ID"
//	@Router		/sse/task/:id/manifest-log [get]
func (a *SseApi) PushManifestLogEvent(c *gin.Context) {
	mtx := &sync.Mutex{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		a.sendEvent(c, mtx, constant.SseErrorEvent, err.Error())
		return
	}
	w := c.Writer
	a.setSSEHeader(c)

	taskLogService := task.GetTaskLogService()
	manifestFile, err := taskLogService.OpenLatestManifestLogFile(id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			a.sendEvent(c, mtx, constant.SseErrorEvent, "日志文件不存在")
			return
		}
		a.sendEvent(c, mtx, constant.SseErrorEvent, err.Error())
		return
	}
	defer func() { _ = manifestFile.Close() }()
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
	reader := bufio.NewReader(manifestFile)
	for {
		if connClosed {
			break
		}
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			a.sendEvent(c, mtx, constant.SseErrorEvent, err.Error())
			return
		}
		if line != "" {
			entry, err := taskLogService.ParseManifestEntry(line)
			if err != nil {
				if errors.Is(err, task.ErrReachedLogEnd) {
					a.sendEvent(c, mtx, constant.SseCloseEvent, "")
					time.Sleep(2 * time.Second)
					return
				} else {
					a.sendEvent(c, mtx, constant.SseErrorEvent, err.Error())
					return
				}
			}
			a.sendEvent(c, mtx, constant.SseManifestEntryEvent, entry)
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
//	@Param		id			path		int64	true	"任务ID"
//	@Param		host		query		string	true	"主机IP"
//	@Param		hostLogName	query		string	true	"主机日志的文件名"
//	@Router		/sse/task/:id/host-log [get]
func (a *SseApi) PushHostLogEvent(c *gin.Context) {
	mtx := &sync.Mutex{}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		a.sendEvent(c, mtx, constant.SseErrorEvent, err.Error())
		return
	}
	host := c.Query("host")
	hostLogName := c.Query("hostLogName")
	if host == "" || hostLogName == "" {
		a.sendEvent(c, mtx, constant.SseErrorEvent, "缺少参数")
		return
	}
	w := c.Writer
	a.setSSEHeader(c)

	taskLogService := task.GetTaskLogService()
	hostLogFile, err := taskLogService.OpenHostLogFile(id, host, hostLogName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			a.sendEvent(c, mtx, constant.SseErrorEvent, "日志文件不存在")
			return
		}
		a.sendEvent(c, mtx, constant.SseErrorEvent, err.Error())
		return
	}
	defer func() { _ = hostLogFile.Close() }()
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
	reader := bufio.NewReader(hostLogFile)
	for {
		if connClosed {
			break
		}
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			a.sendEvent(c, mtx, constant.SseErrorEvent, err.Error())
			return
		}
		if line != "" {
			if taskLogService.IsLogEnd(line) {
				a.sendEvent(c, mtx, constant.SseCloseEvent, "")
				time.Sleep(2 * time.Second)
				return
			}
			a.sendEvent(c, mtx, constant.SseHostLogEvent, line)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (a *SseApi) sendEvent(c *gin.Context, mtx *sync.Mutex, event string, msg any) {
	mtx.Lock()
	defer mtx.Unlock()
	c.SSEvent(event, msg)
	c.Writer.Flush()
}

func (a *SseApi) startHeartbeat(ctx context.Context, c *gin.Context, mtx *sync.Mutex) {
	mtx.Lock()
	c.SSEvent(constant.SseTickEvent, "")
	c.Writer.Flush()
	mtx.Unlock()
	ticker := time.NewTicker(constant.SseTickDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mtx.Lock()
			c.SSEvent(constant.SseTickEvent, "")
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
