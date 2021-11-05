package models

import "gorm.io/gorm"

type CartDetails struct {
	gorm.Model
	CartID   int `gorm:"primaryKey" json:"cart_id" form:"cart_id"`
	RecipeID int `gorm:"primaryKey" json:"recipe_id" form:"recipe_id"`
	Quantity int `json:"quantity" form:"quantity"`
	Price    int `json:"price" form:"price"`
}

type GormCartDetailsModel struct {
	db *gorm.DB
}

func NewCartDetailsModel(db *gorm.DB) *GormCartDetailsModel {
	return &GormCartDetailsModel{db: db}
}

type CartDetailsModel interface {
	GetAllRecipeByCartId(cartId int) ([]CartDetails, error)
	AddRecipeToCart(cartDetails CartDetails, cartId int) (CartDetails, error)
	UpdateRecipePortion(cartDetails CartDetails, recipeId int) (CartDetails, error)
	DeleteRecipeFromCart(cartId int) (CartDetails, error)
	CountQtyRecipeOnCart(cartId int) (int, error)
	CountTotalPriceOnCart(cartId int) (int, error)
}

func (m *GormCartDetailsModel) GetAllRecipeByCartId(cartId int) ([]CartDetails, error) {
	var cartDetails []CartDetails
	if err := m.db.Find(&cartDetails, cartId).Error; err != nil {
		return nil, err
	}
	return cartDetails, nil
}

func (m *GormCartDetailsModel) AddRecipeToCart(cartDetails CartDetails, cartId int) (CartDetails, error) {
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

	cartDetails.RecipeID = newCartDetails.RecipeID

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

func (m *GormCartDetailsModel) CountQtyRecipeOnCart(cartId int) (CartDetails, error) {
	var cartDetails CartDetails
	if err := m.db.Find(&cartDetails, "id=?", cartId).Error; err != nil {
		return cartDetails, err
	}
	if err := m.db.Delete(&cartDetails).Error; err != nil {
		return cartDetails, err
	}
	return cartDetails, nil
}

func (m *GormCartDetailsModel) CountTotalPriceOnCart(cartId int) (CartDetails, error) {
	var cartDetails CartDetails
	if err := m.db.Find(&cartDetails, "id=?", cartId).Error; err != nil {
		return cartDetails, err
	}
	if err := m.db.Delete(&cartDetails).Error; err != nil {
		return cartDetails, err
	}
	return cartDetails, nil
}
