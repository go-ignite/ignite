package task

import (
	"sync"

	"github.com/go-ignite/ignite/logger"
	"github.com/robfig/cron"
)

type Task struct {
	logger *logger.Logger
	cron   *cron.Cron
	sync.Mutex
}

func New() *Task {
	return &Task{
		cron:   cron.New(),
		logger: logger.GetTaskLogger(),
	}
}

func (t *Task) SetLogger(l *logger.Logger) *Task {
	t.logger = l
	return t
}

func (t *Task) AsyncRun() {
	t.logger.Info("init task")
	t.cron.AddFunc("0 */5 * * * *", t.InstantStats)
	t.cron.AddFunc("0 0 0 * * *", t.DailyStats)
	t.cron.AddFunc("0 0 0 1 * *", t.MonthlyStats)

	t.cron.Start()
}
