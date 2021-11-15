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
	UpdateCheckoutIdOnCartDetails(recipeId, checkoutId int) (CartDetails, error)
	// GetSeller(sellerCity string, ingredientId int) (Stock, error)
}

func (m *GormCheckoutModel) CreateCheckout(checkout Checkout) (Checkout, error) {
	if err := m.db.Save(&checkout).Error; err != nil {
		return checkout, err
	}
	return checkout, nil
}

func (m *GormCheckoutModel) UpdateCheckoutIdOnCartDetails(recipeId, checkoutId int) (CartDetails, error) {
	var cartDetails CartDetails
	var newCartDetails CartDetails
	if err := m.db.Find(&newCartDetails, "recipe_id = ?", recipeId).Error; err != nil {
		return cartDetails, err
	}

	newCartDetails.CheckoutID = checkoutId
	// newRecipe.Category = recipe.Category

	if err := m.db.Save(&newCartDetails).Error; err != nil {
		return newCartDetails, err
	}
	return newCartDetails, nil
}

// func (m *GormCheckoutModel) GetSeller(sellerCity string, recipeId int) (User, error) {
// 	if err := m.db.Raw("SELECT s.stock FROM ")
// }
