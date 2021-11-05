package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	// 1 to 1
	UserID uint
}

type GormCartModel struct {
	db *gorm.DB
}

func NewCartModel(db *gorm.DB) *GormCartModel {
	return &GormCartModel{db: db}
}

type CartModel interface {
	CreateCart(cart Cart) (Cart, error)
}

func (m *GormCartModel) CreateCart(cart Cart) (Cart, error) {
	// Select(&cart, "NOT EXISTS (SELECT * FROM carts WHERE user_id=?", cart.UserID).
	if err := m.db.Save(&cart).Error; err != nil {
		return cart, err
	}
	return cart, nil
}
