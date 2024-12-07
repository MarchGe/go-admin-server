package scheduler

import (
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/robfig/cron/v3"
	"sync"
)

var c = cron.New(cron.WithParser(cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)))

type Task func()

var taskMap = make(map[int64]cron.EntryID, 0)

func Start() {
	c.Start()
}

var taskMapMutex sync.Mutex

func AddTask(buzId int64, cronStr string, t Task) error {
	taskMapMutex.Lock()
	defer taskMapMutex.Unlock()
	if _, ok := taskMap[buzId]; ok {
		return nil
	}
	entryId, err := c.AddFunc(cronStr, t)
	if err != nil {
		return fmt.Errorf("add task to scheduler error, %w", err)
	}
	taskMap[buzId] = entryId
	return nil
}

func RemoveTask(buzId int64) {
	taskMapMutex.Lock()
	defer taskMapMutex.Unlock()
	if entryId, ok := taskMap[buzId]; ok {
		c.Remove(entryId)
		delete(taskMap, buzId)
	}
}

func VerifyCron(cronStr string) error {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := parser.Parse(cronStr); err != nil {
		return E.Message("Cron表达式解析失败")
	}
	return nil
}

func Stop() {
	c.Stop()
}
