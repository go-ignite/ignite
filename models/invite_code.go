package models

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/go-ignite/ignite/utils"
)

type InviteCode struct {
	gorm.Model
	Code      string
	Available bool
	ExpiredAt time.Time
	Limit     uint //unit: GB/month
}

func NewInviteCode(limit uint, expiredAt time.Time) *InviteCode {
	return &InviteCode{
		Code:      utils.RandString(16),
		Limit:     limit,
		ExpiredAt: expiredAt,
		Available: true,
	}
}

func CreateInviteCodes(inviteCodes []*InviteCode) error {
	tx := db.Begin()
	if err := func() error {
		for _, inviteCode := range inviteCodes {
			if err := inviteCode.create(tx); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (ic *InviteCode) create(session *gorm.DB) error {
	return session.Create(ic).Error
}

func (ic *InviteCode) used(session *gorm.DB) error {
	// TODO
	return nil
}

func GetInviteCodeByCode(code string) (*InviteCode, error) {
	iv := new(InviteCode)
	r := db.First(iv, "code = ? AND available = 1", code)
	if r.RecordNotFound() {
		return nil, nil
	}
	return iv, r.Error
}

func DeleteInviteCodeByID(id int64) error {
	return db.Delete(new(InviteCode), "id = ?", id).Error
}

func GetAvailableInviteCodeList(pageIndex, pageSize int) ([]*InviteCode, int, error) {
	where := func() *gorm.DB {
		return db.Model(new(InviteCode)).Where("available = 1")
	}

	var total int
	if err := where().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var inviteCodes []*InviteCode
	return inviteCodes, total, where().Offset(pageSize * (pageIndex - 1)).Limit(pageSize).Find(&inviteCodes).Error
}
