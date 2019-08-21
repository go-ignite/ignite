package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid"

	"github.com/go-ignite/ignite/api"
)

type User struct {
	ID              string `gorm:"primary_key"`
	Name            string
	HashedPwd       []byte
	InviteCode      string
	ServicePassword string
	PackageLimit    int // unit: GB/month
	ExpiredAt       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `sql:"index"`
}

func NewUser(name string, hashedPwd []byte, inviteCode string) *User {
	return &User{
		ID:              shortuuid.New(),
		ServicePassword: shortuuid.New(),
		Name:            name,
		HashedPwd:       hashedPwd,
		InviteCode:      inviteCode,
	}
}

func (u User) Output() *api.User {
	return &api.User{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}

func (h *Handler) GetUsers() ([]*User, error) {
	var users []*User
	return users, h.db.Find(&users).Error
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

func (h *Handler) SaveUser(u *User) error {
	return h.db.Save(u).Error
}

func (h *Handler) CreateUser(u *User) error {
	return h.runTX(func(h *Handler) error {
		ic := new(InviteCode)
		if err := h.db.First(ic, "code = ? AND available = 1", u.InviteCode).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return api.ErrInviteCodeNotExistOrUnavailable
			}

			return err
		}

		if ic.ExpiredAt.Before(time.Now()) {
			return api.ErrInviteCodeExpired
		}

		if err := h.db.Model(InviteCode{ID: ic.ID}).UpdateColumn("available", false).Error; err != nil {
			return err
		}

		user, err := h.GetUserByName(u.Name)
		if err != nil {
			return err
		}

		if user != nil {
			return api.ErrUserNameExists
		}

		u.ExpiredAt = ic.ExpiredAt
		u.PackageLimit = ic.Limit
		return h.db.Create(u).Error
	})
}

func (h *Handler) DestroyUser(userID string) error {
	return h.runTX(func(h *Handler) error {
		user, err := h.GetUserByID(userID)
		if err != nil {
			return err
		}

		if user == nil {
			return nil
		}

		if err := h.db.Delete(&User{ID: userID}).Error; err != nil {
			return err
		}

		// FIXME depends on the agent response to remove services
		if err := h.db.Delete(Service{}, "user_id = ?", userID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (h *Handler) ChangeUserPassword(id string, hashedPwd []byte) error {
	return h.db.Model(User{ID: id}).UpdateColumn("hashed_pwd", hashedPwd).Error
}
