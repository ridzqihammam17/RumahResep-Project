package models

import "gorm.io/gorm"

type CartDetails struct {
	gorm.Model
	RecipeID int `gorm:"primaryKey" json:"recipes_id" form:"recipes_id"`
	CartID   int `gorm:"primaryKey" json:"carts_id" form:"carts_id"`
	Quantity int `json:"quantity" form:"quantity"`
	Price    int `json:"price" form:"price"`
}
