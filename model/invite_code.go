package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"

	"github.com/go-ignite/ignite/api"
)

type InviteCode struct {
	ID        int64 `gorm:"primary_key"`
	Code      string
	Available bool
	ExpiredAt time.Time
	Limit     int //unit: GB/month
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func NewInviteCode(limit int, expiredAt time.Time) *InviteCode {
	return &InviteCode{
		Code:      shortuuid.New(),
		Limit:     limit,
		ExpiredAt: expiredAt,
		Available: true,
	}
}

func (ic InviteCode) Output() *api.InviteCode {
	return &api.InviteCode{
		ID:        ic.ID,
		Code:      ic.Code,
		Limit:     ic.Limit,
		ExpiredAt: ic.ExpiredAt,
		CreatedAt: ic.CreatedAt,
	}
}

func (h *Handler) CreateInviteCodes(inviteCodes []*InviteCode) error {
	return h.runTX(func(h *Handler) error {
		for _, inviteCode := range inviteCodes {
			if err := h.db.Create(inviteCode).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (h *Handler) DeleteInviteCode(id int64) error {
	return h.db.Delete(&InviteCode{ID: id}).Error
}

func (h *Handler) GetAvailableInviteCodeList(pageIndex, pageSize int) ([]*InviteCode, int, error) {
	where := func() *gorm.DB {
		return h.db.Model(InviteCode{}).Where("available = 1")
	}

	var total int
	if err := where().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var inviteCodes []*InviteCode
	if err := where().Offset(pageSize * (pageIndex - 1)).Limit(pageSize).Find(&inviteCodes).Error; err != nil {
		return nil, 0, err
	}

	return inviteCodes, total, nil
}
