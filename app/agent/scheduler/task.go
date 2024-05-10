package scheduler

type Task interface {
	ID() string
	NAME() string
	CRON() string
	SetRuntimeId(id string)
	Execute()
}

type TaskI struct {
	Id        string
	Name      string
	Cron      string
	RuntimeId string // 无需设置，由cron框架自动生成
}

func (s *TaskI) ID() string {
	return s.Id
}

func (s *TaskI) NAME() string {
	return s.Name
}

func (s *TaskI) CRON() string {
	return s.Cron
}

func (s *TaskI) SetRuntimeId(id string) {
	s.RuntimeId = id
}
