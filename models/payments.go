package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	PaymentMethod string `json:"payment_method" form:"payment_method"`
	PaymentStatus string `json:"payment_status" form:"payment_status"`
	// One to One with Transaction
	TransactionID uint
}

type GormPaymentModel struct {
	db *gorm.DB
}

func NewPaymentModel(db *gorm.DB) *GormPaymentModel {
	return &GormPaymentModel{db: db}
}
