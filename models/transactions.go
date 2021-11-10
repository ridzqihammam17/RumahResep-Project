package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID         uint
	CustomerName   string `json:"customer_name" form:"customer_name"`
	Address        string `json:"address" form:"address"`
	ShippingMethod string `json:"shipping_method" form:"shipping_method"`
	PaymentMethod  string `json:"payment_method" form:"payment_method"`
	PaymentStatus  string `json:"payment_status" form:"payment_status"`
	TotalPayment   int    `json:"total_payment" form:"total_payment"`
	// One to Many with Checkout
	CheckoutID uint
}

type GormTransactionModel struct {
	db *gorm.DB
}

func NewTransactionModel(db *gorm.DB) *GormTransactionModel {
	return &GormTransactionModel{db: db}
}

type TransactionModel interface {
	CreateTransaction(Transaction) (Transaction, error)
	GetCheckoutId(cartId int) (int, error)
	GetUserData(userId int) (Transaction, error)
	CountTotalPayment(cartId, checkoutId int) (int, error)
	Get(transactionId int) (Transaction, error)
	Add(Transaction) (Transaction, error)
	GetTotalPayment(transactionId int) (int, error)
	// ChooseShippingPaymentMethod(checkoutId int, transaction Transaction) (Transaction, error)
}

func (m *GormTransactionModel) CreateTransaction(transactions Transaction) (Transaction, error) {
	if err := m.db.Save(&transactions).Error; err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (m *GormTransactionModel) GetCheckoutId(cartId int) (int, error) {
	// var cartDetails CartDetails
	var checkoutId int

	if err := m.db.Raw("SELECT checkout_id FROM cart_details WHERE checkout_id IS NOT NULL AND checkout_id != 0 AND cart_id = ?", cartId).Scan(&checkoutId).Error; err != nil {
		return checkoutId, err
	}
	return checkoutId, nil
	// if err := m.db.Find(&cartDetails, cartId).Error; err != nil {
	// 	return cartDetails.CheckoutID, err
	// }
	// return cartDetails.CheckoutID, nil
}

func (m *GormTransactionModel) GetUserData(userId int) (Transaction, error) {
	var user User
	var transaction Transaction
	if err := m.db.Find(&user, userId).Error; err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (m *GormTransactionModel) CountTotalPayment(cartId, checkoutId int) (int, error) {
	var countTotalPayment int
	// if err := m.db.Select("COUNT(carts_id)").Where("checkout_id=?", checkoutId).First(&countProduct).Error; err == nil {
	if err := m.db.Raw("SELECT SUM(price) FROM cart_details WHERE cart_id = ? AND checkout_id = ?", cartId, checkoutId).Scan(&countTotalPayment).Error; err != nil {
		return countTotalPayment, err
	}
	return countTotalPayment, nil
}

func (m *GormTransactionModel) GetTotalPayment(transactionId int) (int, error) {
	var transaction Transaction
	// var totalPayment int
	if err := m.db.Find(&transaction, transactionId).Error; err != nil {
		return transaction.TotalPayment, err
	}
	return transaction.TotalPayment, nil
}

func (m *GormTransactionModel) Get(transactionId int) (Transaction, error) {
	var transaction Transaction
	if err := m.db.Where("id=?", transactionId).First(&transaction).Error; err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (m *GormTransactionModel) Add(transaction Transaction) (Transaction, error) {
	return transaction, nil
}
