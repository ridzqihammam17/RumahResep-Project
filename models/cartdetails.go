package models

import "gorm.io/gorm"

type CartDetails struct {
	gorm.Model
	CartID   int `gorm:"primaryKey" json:"cart_id" form:"cart_id"`
	RecipeID int `gorm:"primaryKey" json:"recipe_id" form:"recipe_id"`
	Quantity int `json:"quantity" form:"quantity"`
	Price    int `json:"price" form:"price"`
}

type RecipePriceResponse struct {
	RecipeID      int `gorm:"primaryKey" json:"recipe_id" form:"recipe_id"`
	TotalQuantity int `json:"total_qty" form:"total_qty"`
	TotalPrice    int `json:"total_price" form:"total_price"`
}

type GormCartDetailsModel struct {
	db *gorm.DB
}

func NewCartDetailsModel(db *gorm.DB) *GormCartDetailsModel {
	return &GormCartDetailsModel{db: db}
}

type CartDetailsModel interface {
	GetAllRecipeByCartId(cartId int) ([]CartDetails, error)
	// GetRecipePriceByRecipeId(recipeId int) (RecipePriceResponse, error)
	AddRecipeToCart(cartDetails CartDetails) (CartDetails, error)
	UpdateRecipePortion(cartDetails CartDetails, recipeId int) (CartDetails, error)
	DeleteRecipeFromCart(cartId int) (CartDetails, error)
	CountQtyRecipeOnCart(cartId int) (int, error)
	CountTotalPriceOnCart(cartId int) (int, int, error)
}

func (m *GormCartDetailsModel) GetAllRecipeByCartId(cartId int) ([]CartDetails, error) {
	var cartDetails []CartDetails
	if err := m.db.Find(&cartDetails, cartId).Error; err != nil {
		return nil, err
	}
	return cartDetails, nil
}

// func (m *GormCartDetailsModel) GetRecipePriceByRecipeId(recipeId int) ([]CartDetails, error) {

// }

func (m *GormCartDetailsModel) AddRecipeToCart(cartDetails CartDetails) (CartDetails, error) {
	if err := m.db.Save(&cartDetails).Error; err != nil {
		return cartDetails, err
	}
	return cartDetails, nil
}

func (m *GormCartDetailsModel) UpdateRecipePortion(newCartDetails CartDetails, recipeId int) (CartDetails, error) {
	var cartDetails CartDetails
	if err := m.db.Find(&cartDetails, "id=?", recipeId).Error; err != nil {
		return cartDetails, err
	}

	cartDetails.Quantity = newCartDetails.Quantity

	if err := m.db.Save(&cartDetails).Error; err != nil {
		return cartDetails, err
	}
	return cartDetails, nil
}

func (m *GormCartDetailsModel) DeleteRecipeFromCart(cartId int) (CartDetails, error) {
	var cartDetails CartDetails
	if err := m.db.Find(&cartDetails, "id=?", cartId).Error; err != nil {
		return cartDetails, err
	}
	if err := m.db.Delete(&cartDetails).Error; err != nil {
		return cartDetails, err
	}
	return cartDetails, nil
}

func (m *GormCartDetailsModel) CountQtyRecipeOnCart(cartId int) (int, error) {
	var countProduct int
	if err := m.db.Select("COUNT(carts_id)").Where("carts_id=?", cartId).First(&countProduct).Error; err == nil {
		return countProduct, err
	}
	return countProduct, nil
}

func (m *GormCartDetailsModel) CountTotalPriceOnCart(cartId int) (int, int, error) {
	var countProduct, Price int
	if err := m.db.Select("COUNT(carts_id), Price").Where("carts_id=?", cartId).First(&countProduct).Error; err == nil {
		return countProduct, Price, err
	}
	return countProduct, Price, nil
}
