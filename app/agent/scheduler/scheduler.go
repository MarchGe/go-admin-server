package scheduler

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log/slog"
	"strconv"
)

var c *cron.Cron

// AddTask dynamic add task after scheduler started
func AddTask(t Task) (runtimeId string, err error) {
	taskId, e := c.AddFunc(t.CRON(), t.Execute)
	if e != nil {
		err = fmt.Errorf("add task '%s' error, %w", t.NAME(), e)
		return
	}
	runtimeId = strconv.Itoa(int(taskId))
	return
}

// RemoveTask dynamic remove task after scheduler started
func RemoveTask(runtimeId string) error {
	id, err := strconv.Atoi(runtimeId)
	if err != nil {
		return fmt.Errorf("remove task error, taskId: %s, %w", runtimeId, err)
	}
	c.Remove(cron.EntryID(id))
	return nil
}

func Run(ctx context.Context, tasks []Task) {
	if c != nil {
		slog.Error("scheduler has been started already.")
		return
	}
	option := cron.WithParser(cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor))
	c = cron.New(option)
	for i := range tasks {
		runtimeId, err := c.AddFunc(tasks[i].CRON(), tasks[i].Execute)
		if err != nil {
			slog.Error("add task '"+tasks[i].NAME()+"' failed.", slog.Any("err", err))
			continue
		}
		slog.Info("add task '" + tasks[i].NAME() + "' success.")
		tasks[i].SetRuntimeId(strconv.Itoa(int(runtimeId)))
	}
	c.Start()
	slog.Info("scheduler started success.")
	defer func() {
		c.Stop()
		slog.Debug("scheduler stopped success.")
	}()
	<-ctx.Done()
}
