package models

import (
	"rumah_resep/api/middlewares"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" form:"name"`
	Address  string `json:"address" form:"address"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Gender   string `json:"gender" form:"gender"`
	Role     string `json:"role" form:"role"`
	Token    string `json:"token" form:"token"`
	Carts    Cart   `gorm:"foreignKey:UserID"`
}

type GormUserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *GormUserModel {
	return &GormUserModel{db: db}
}

type UserModel interface {
	Register(User) (User, error)
	Login(email, password string) (User, error)
}

func (m *GormUserModel) Register(user User) (User, error) {
	if err := m.db.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (m *GormUserModel) Login(email, password string) (User, error) {
	var user User
	var err error

	if err = m.db.Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		return user, err
	}

	user.Token, err = middlewares.CreateToken(int(user.ID), user.Role)


	if err != nil {
		return user, err
	}

	if err := m.db.Save(user).Error; err != nil {
		return user, err
	}

	return user, nil
}
