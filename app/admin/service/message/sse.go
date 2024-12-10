package message

import (
	"fmt"
	cSse "github.com/MarchGe/go-admin-server/app/common/constant/sse"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"log/slog"
	"sync"
)

var _sseService = &SseService{}

type SseService struct {
}

func GetSseService() *SseService {
	return _sseService
}

type Event string

var sseSessions = make([]*sseSession, 0)
var mtx = sync.RWMutex{}

func (s *SseService) PushEventMessage(uId int64, event cSse.Event, data any) error {
	sessions := s.getSessions(uId)
	for i := range sessions {
		resMtx := sessions[i].Mtx
		resMtx.Lock()
		err := sse.Encode(sessions[i].W, sse.Event{
			Event: event.String(),
			Data:  data,
		})
		if err != nil {
			resMtx.Unlock()
			return fmt.Errorf("encode sse event data error, %w", err)
		}
		sessions[i].W.Flush()
		resMtx.Unlock()
	}
	return nil
}

type sseSession struct {
	Id  string
	Uid int64
	W   gin.ResponseWriter
	Mtx *sync.Mutex
}

func (s *SseService) AddSseSession(id string, uid int64, w gin.ResponseWriter, mtx *sync.Mutex) {
	mtx.Lock()
	defer mtx.Unlock()
	sseSessions = append(sseSessions, &sseSession{
		Id:  id,
		Uid: uid,
		W:   w,
		Mtx: mtx,
	})
}
func (s *SseService) RemoveSseSession(id string) {
	mtx.Lock()
	defer mtx.Unlock()
	var index int
	length := len(sseSessions)
	for i := 0; i < length; i++ {
		if sseSessions[i].Id == id {
			index = i
			break
		}
	}
	if index == length-1 {
		sseSessions = sseSessions[0:index]
	} else {
		sseSessions = append(sseSessions[0:index], sseSessions[index+1:]...)
	}
	slog.Debug("remove sse session", slog.String("id", id))
}

func (s *SseService) getSessions(uId int64) []*sseSession {
	mtx.RLock()
	defer mtx.RUnlock()
	sessions := make([]*sseSession, 0, 1)
	for i := range sseSessions {
		if sseSessions[i].Uid == uId {
			sessions = append(sessions, sseSessions[i])
		}
	}
	return sessions
}
