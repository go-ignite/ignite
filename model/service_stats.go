package model

import (
	"time"
)

type ServiceStats struct {
	ID          int64 `gorm:"primary_key"`
	UserID      string
	NodeID      string
	ServiceID   int64
	StartTime   time.Time
	EndTime     time.Time
	TrafficUsed uint64 // unit: byte
	StatsResult uint64 // unit: byte
	CreatedAt   time.Time
}

func (h *Handler) GetLastServiceStats(serviceID int64) (*ServiceStats, error) {
	ls := new(ServiceStats)
	r := h.db.Last(ls, "service_id = ?", serviceID)
	if r.RecordNotFound() {
		return nil, nil
	}

	return ls, r.Error
}
