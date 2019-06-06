package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"
)

type User struct {
	ID           string `gorm:"primary_key"`
	Name         string
	HashedPwd    []byte
	InviteCodeID int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
}

func NewUser(name string, hashedPwd []byte, inviteCodeID int64) *User {
	return &User{
		ID:           shortuuid.New(),
		Name:         name,
		HashedPwd:    hashedPwd,
		InviteCodeID: inviteCodeID,
	}
}

func (h *Handler) GetUserList(keyword string, pageIndex, pageSize int) ([]*User, int, error) {
	where := func() *gorm.DB {
		return h.db.Model(new(User)).Where("name like ?", "%"+keyword+"%")
	}

	var total int
	if err := where().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []*User
	if err := where().Offset(pageSize * (pageIndex - 1)).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (h *Handler) GetUserByNameAndPassword(name string, hashedPassword []byte) (*User, error) {
	user := new(User)
	r := h.db.First(user, "name = ? AND hashed_pwd = ?", name, hashedPassword)
	if r.RecordNotFound() {
		return nil, nil
	}

	if r.Error != nil {
		return nil, r.Error
	}

	return user, nil
}

func (h *Handler) GetUserByID(id string) (*User, error) {
	user := new(User)
	r := h.db.First(user, id)

	if r.RecordNotFound() {
		return nil, nil
	}

	if r.Error != nil {
		return nil, r.Error
	}

	return user, nil
}

func (h *Handler) SaveUser(u *User) error {
	return h.db.Save(u).Error
}

func (h *Handler) CreateUser(u *User) error {
	return h.runTX(func(tx *gorm.DB) error {
		r := tx.Model(InviteCode{ID: u.InviteCodeID}).UpdateColumn("available", true)
		if r.Error != nil {
			return r.Error
		}
		if r.RowsAffected == 0 {
			return fmt.Errorf("model: invite code has been used, id: %d", u.InviteCodeID)
		}

		return tx.Create(u).Error
	})
}

func (h *Handler) DestroyUser(userID string) error {
	return h.runTX(func(tx *gorm.DB) error {
		if err := tx.Delete(&User{ID: userID}).Error; err != nil {
			return err
		}

		return tx.Delete(Service{}, "user_id = ?", userID).Error
	})
}
