package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model

	// 1 to 1
	UserID uint
}
