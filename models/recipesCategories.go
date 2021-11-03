package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type RecipesCategories struct {
	// gorm.Model
	RecipeId   int `gorm:"primaryKey" json:"recipes_id" form:"recipes_id"`
	CategoryId int `gorm:"primaryKey" json:"categories_id" form:"categories_id"`
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

type GormRecipesCategoriesModel struct {
	db *gorm.DB
}

func NewRecipesCategoriesModel(db *gorm.DB) *GormRecipesCategoriesModel {
	return &GormRecipesCategoriesModel{db: db}
}

type RecipesCategoriesModel interface {
	GetRecipeByCategoryId(categoryId []int) ([]Recipe, error)
}

func (m *GormRecipesCategoriesModel) GetRecipeByCategoryId(categoryId []int) ([]Recipe, error) {
	var recipe []Recipe

	if err := m.db.Find(&recipe, "category in ?", categoryId).Error; err != nil {
		return recipe, err
	}
	if len(recipe) == 0 {

		return nil, errors.New("Data Not Found")
	}
	return recipe, nil
}
