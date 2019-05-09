package models

import (
	"github.com/jinzhu/gorm"
)

type Users []*User

type User struct {
	gorm.Model
	Name         string
	HashedPwd    []byte
	InviteCodeID uint
}

func NewUser(name string, hashedPwd []byte, inviteCodeID uint) *User {
	return &User{
		Name:         name,
		HashedPwd:    hashedPwd,
		InviteCodeID: inviteCodeID,
	}
}

func GetUserList(keyword string, pageIndex, pageSize int) (Users, int, error) {
	where := func() *gorm.DB {
		return db.Model(new(User)).Where("name like ?", "%"+keyword+"%")
	}

	var total int
	if err := where().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users Users
	return users, total, where().Offset(pageSize * (pageIndex - 1)).Limit(pageSize).Find(&users).Error
}

func GetUserByNameAndPassword(name string, hashedPassword []byte) (*User, error) {
	user := new(User)
	r := db.First(user, "name = ? AND hashed_pwd = ?", name, hashedPassword)
	if r.RecordNotFound() {
		return nil, nil
	}
	return user, r.Error
}

func GetUserByID(id uint) (*User, error) {
	user := new(User)
	r := db.First(user, id)
	if r.RecordNotFound() {
		return nil, nil
	}
	return user, r.Error
}

func (u *User) Save() error {
	return db.Save(u).Error
}

func (u *User) Create(iv *InviteCode) error {
	tx := db.Begin()
	if err := func() error {
		if err := u.create(tx); err != nil {
			return err
		}
		if err := iv.used(tx); err != nil {
			return err
		}

		return nil
	}(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (u *User) create(session *gorm.DB) error {
	return session.Create(u).Error
}

func (u *User) Destroy(serviceIDs []uint) error {
	tx := db.Begin()
	if err := tx.Delete(u).Error; err != nil {
		return err
	}
	return deleteServicesByIDs(tx, serviceIDs)
}

func (u *User) GetServices() (Services, error) {
	var services Services
	return services, db.Find(&services, "user_id = ?", u.ID).Error
}
