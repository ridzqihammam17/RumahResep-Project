package models

import (
	"rumah_resep/api/middlewares"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string  `json:"name" form:"name"`
	Address   string  `json:"address" form:"address"`
	Latitude  float64 `json:"latitude" form:"latitude"`
	Longitude float64 `json:"laongitude" form:"longitude"`
	Email     string  `json:"email" form:"email"`
	Password  string  `json:"password" form:"password"`
	Gender    string  `json:"gender" form:"gender"`
	Role      string  `json:"role" form:"role"`
	Token     string  `json:"token" form:"token"`
	Cart      Cart    `gorm:"foreignKey:UserID"`
	// Ingredients  []Ingredient
	Stocks       []Stock
	Transactions []Transaction
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
	GetUserData(userId int) (User, error)
}

func (m *GormUserModel) Register(user User) (User, error) {
	// Encrypt Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.Password = string(hashedPassword)

	if err := m.db.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (m *GormUserModel) Login(email, password string) (User, error) {
	var user User
	var err error

	if err = m.db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	// Start Checking Hash Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
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

func (m *GormUserModel) GetUserData(userId int) (User, error) {
	var user User
	if err := m.db.Find(&user, userId).Error; err != nil {
		return user, err
	}
	return user, nil
}
