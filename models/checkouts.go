package models

import "gorm.io/gorm"

type Checkout struct {
	gorm.Model
	// One to Many with Cart Details
	CartDetails []CartDetails
	// One to Many with Transactions
	Transactions []Transaction
}

type GormCheckoutModel struct {
	db *gorm.DB
}

func NewCheckoutModel(db *gorm.DB) *GormCheckoutModel {
	return &GormCheckoutModel{db: db}
}

type CheckoutModel interface {
	CreateCheckout(checkout Checkout) (Checkout, error)
}

func (m *GormCheckoutModel) CreateCheckout(checkout Checkout) (Checkout, error) {
	if err := m.db.Save(&checkout).Error; err != nil {
		return checkout, err
	}
	return checkout, nil
}
