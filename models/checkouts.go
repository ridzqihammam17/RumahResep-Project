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
	UpdateCheckoutIdOnCartDetails(recipeId, checkoutId, cartId int) (CartDetails, error)
}

func (m *GormCheckoutModel) CreateCheckout(checkout Checkout) (Checkout, error) {
	if err := m.db.Save(&checkout).Error; err != nil {
		return checkout, err
	}
	return checkout, nil
}

func (m *GormCheckoutModel) UpdateCheckoutIdOnCartDetails(recipeId, checkoutId, cartId int) (CartDetails, error) {
	var cartDetails CartDetails
	// var newCartDetails CartDetails

	if err := m.db.Raw("UPDATE cart_details SET checkout_id = ? WHERE recipe_id = ? AND cart_id = ?", checkoutId, recipeId, cartId).Scan(&cartDetails).Error; err != nil {
		return cartDetails, err
	}

	return cartDetails, nil
	// if err := m.db.Find(&newCartDetails, "recipe_id = ?", recipeId).Error; err != nil {
	// 	return cartDetails, err
	// }

	// newCartDetails.CheckoutID = checkoutId
	// // newRecipe.Category = recipe.Category

	// if err := m.db.Save(&newCartDetails).Error; err != nil {
	// 	return newCartDetails, err
	// }
	// return newCartDetails, nil
}
