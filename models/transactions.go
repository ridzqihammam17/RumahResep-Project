package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID         uint
	CustomerName   string `json:"customer_name" form:"customer_name"`
	Address        string `json:"address" form:"address"`
	ShippingMethod string `json:"shipping_method" form:"shipping_method"`
	PaymentStatus  string `json:"payment_status" form:"payment_status"`
	TotalPayment   int    `json:"total_payment" form:"total_payment"`
	// One to Many with Checkout
	CheckoutID uint
	// One to One with Payment
	Payment Payment
}

type GormTransactionModel struct {
	db *gorm.DB
}

func NewTransactionModel(db *gorm.DB) *GormTransactionModel {
	return &GormTransactionModel{db: db}
}

type TransactionModel interface {
	CreateTransaction() (Transaction, error)
}
