package models

import (
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
