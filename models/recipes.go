package models

import "gorm.io/gorm"

type Recipe struct {
	gorm.Model
	Name  string `json:"name" form:"name"`
	Price int    `json:"price" form:"price"`
	Stock int    `json:"stock" form:"stock"`

	// many2many with category
	Categories []Category `gorm:"many2many:recipe_categories"`
}
