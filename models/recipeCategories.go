package models

import (
	"errors"

	"gorm.io/gorm"
)

type RecipeCategories struct {
	gorm.Model
	// ID         int `gorm:"primaryKey" sql:"AUTO_INCREMENT"`
	RecipeId   int `gorm:"primaryKey" json:"recipes_id" form:"recipes_id"`
	CategoryId int `gorm:"primaryKey" json:"categories_id" form:"categories_id"`
	// CreatedAt  time.Time
	// UpdatedAt  time.Time
	// DeletedAt  gorm.DeletedAt
}

type GormRecipesCategoriesModel struct {
	db *gorm.DB
}

func NewRecipesCategoriesModel(db *gorm.DB) *GormRecipesCategoriesModel {
	return &GormRecipesCategoriesModel{db: db}
}

type RecipeCategoriesModel interface {
	AddRecipeCategories(recipeCategories RecipeCategories) (RecipeCategories, error)
	GetRecipeByCategoryId(categoryId []int) ([]RecipeCategories, error)
}

func (m *GormRecipesCategoriesModel) AddRecipeCategories(recipeCategories RecipeCategories) (RecipeCategories, error) {
	if err := m.db.Save(&recipeCategories).Error; err != nil {
		return recipeCategories, err
	}
	return recipeCategories, nil
}

func (m *GormRecipesCategoriesModel) GetRecipeByCategoryId(categoryId []int) ([]RecipeCategories, error) {
	var recipe_categories []RecipeCategories

	if err := m.db.Find(&recipe_categories, "category_id", categoryId).Error; err != nil {
		return recipe_categories, err
	}
	if len(recipe_categories) == 0 {

		return nil, errors.New("Data Not Found")
	}
	return recipe_categories, nil
}
