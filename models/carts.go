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
	CreateCart(cart Cart, userId int) (Cart, error)
	GetCartIdByUserId(userId int) (int, error)
}

func (m *GormCartModel) CreateCart(cart Cart, userId int) (Cart, error) {
	// Select(&cart, "NOT EXISTS (SELECT * FROM carts WHERE user_id=?", cart.UserID).
	// var cartUserId Cart
	// SELECT user_id FROM carts WHERE EXISTS (SELECT * FROM t2);
	// if err := m.db.Raw("SELECT id, user_id FROM carts WHERE user_id = ", userId).Scan(&cart).Error; err != nil {
	// 	return cart, err
	// }

	if err := m.db.Save(&cart).Error; err != nil {
		return cart, err
	}
	return cart, nil
}

func (m *GormCartModel) GetCartIdByUserId(userId int) (int, error) {
	var cartId int
	if err := m.db.Raw("SELECT id FROM carts WHERE user_id = ?", userId).Scan(&cartId).Error; err != nil {
		return cartId, err
	}
	return cartId, nil
}
