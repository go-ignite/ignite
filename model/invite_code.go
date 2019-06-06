package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/go-ignite/ignite/utils"
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
		Code:      utils.RandString(16),
		Limit:     limit,
		ExpiredAt: expiredAt,
		Available: true,
	}
}

func (h *Handler) CreateInviteCodes(inviteCodes []*InviteCode) error {
	return h.runTX(func(tx *gorm.DB) error {
		for _, inviteCode := range inviteCodes {
			if err := tx.Create(inviteCode).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (h *Handler) GetInviteCode(code string) (*InviteCode, error) {
	iv := new(InviteCode)
	r := h.db.First(iv, "code = ? AND available = 1", code)
	if r.RecordNotFound() {
		return nil, nil
	}

	return iv, r.Error
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
