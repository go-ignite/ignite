package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"

	"github.com/go-ignite/ignite/api"
)

var (
	ErrInviteCodeNotExistOrUnavailable = errors.New("model: invite code does not exist or is unavailable")
	ErrInviteCodeExpired               = errors.New("model: invite code is expired")
	ErrUserNameExists                  = errors.New("model: user name already exists")
	ErrUserDeleted                     = errors.New("model: user has been deleted")
)

type User struct {
	ID              string `gorm:"primary_key"`
	Name            string
	HashedPwd       []byte
	InviteCodeID    int64
	ServicePassword string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `sql:"index"`
}

func NewUser(name string, hashedPwd []byte) *User {
	return &User{
		ID:              shortuuid.New(),
		ServicePassword: shortuuid.New(),
		Name:            name,
		HashedPwd:       hashedPwd,
	}
}

func (u User) Output() *api.User {
	return &api.User{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}

func (h *Handler) GetUserList(keyword string, pageIndex, pageSize int) ([]*User, int, error) {
	where := func() *gorm.DB {
		return h.db.Model(User{}).Where("name like ?", "%"+keyword+"%")
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

func (h *Handler) GetUserByName(name string) (*User, error) {
	user := new(User)
	r := h.db.First(user, "name = ?", name)

	if r.RecordNotFound() {
		return nil, nil
	}

	return user, r.Error
}

func (h *Handler) GetUserByID(id string) (*User, error) {
	user := new(User)
	r := h.db.First(user, "id = ?", id)

	if r.RecordNotFound() {
		return nil, nil
	}

	return user, r.Error
}

func (h *Handler) mustGetUserByID(id string) (*User, error) {
	u, err := h.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, ErrUserDeleted
	}

	return u, nil
}

func (h *Handler) SaveUser(u *User) error {
	return h.db.Save(u).Error
}

func (h *Handler) CreateUser(u *User, inviteCode string) error {
	return h.runTX(func(tx *gorm.DB) error {
		ic := new(InviteCode)
		if err := h.db.First(ic, "code = ? AND available = 1", inviteCode).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return ErrInviteCodeNotExistOrUnavailable
			}

			return err
		}

		if ic.ExpiredAt.Before(time.Now()) {
			return ErrInviteCodeExpired
		}

		if err := tx.Model(InviteCode{ID: ic.ID}).UpdateColumn("available", false).Error; err != nil {
			return err
		}

		user, err := newHandler(tx).GetUserByName(u.Name)
		if err != nil {
			return err
		}

		if user != nil {
			return ErrUserNameExists
		}

		u.InviteCodeID = ic.ID
		return tx.Create(u).Error
	})
}

func (h *Handler) DestroyUser(userID string, f func() error) error {
	return h.runTX(func(tx *gorm.DB) error {
		user, err := newHandler(tx).GetUserByID(userID)
		if err != nil {
			return err
		}

		if user == nil {
			return nil
		}

		if err := tx.Delete(&User{ID: userID}).Error; err != nil {
			return err
		}

		// FIXME depends on the agent response to remove services
		if err := tx.Delete(Service{}, "user_id = ?", userID).Error; err != nil {
			return err
		}

		return f()
	})
}
