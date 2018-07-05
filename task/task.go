package task

import (
	"sync"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type Task struct {
	*logrus.Logger
	*cron.Cron
	sync.Mutex
}

func New(l *logrus.Logger) *Task {
	return &Task{
		Logger: l,
		Cron:   cron.New(),
	}
}

func (t *Task) Init() {
	t.Info("init task")
	t.AddFunc("0 */5 * * * *", t.InstantStats)
	t.AddFunc("0 0 0 * * *", t.DailyStats)
	t.AddFunc("0 0 0 1 * *", t.MonthlyStats)
}
