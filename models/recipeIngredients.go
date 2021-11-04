package models

import (
	"gorm.io/gorm"
)

type RecipeIngredients struct {
	gorm.Model
	RecipeId     int `gorm:"primaryKey" json:"recipes_id" form:"recipes_id"`
	IngredientId int `gorm:"primaryKey" json:"ingredients_id" form:"ingredients_id"`
}

type GormRecipeIngredientsModel struct {
	db *gorm.DB
}

func NewRecipeIngredientsModel(db *gorm.DB) *GormRecipeIngredientsModel {
	return &GormRecipeIngredientsModel{db: db}
}

type RecipeIngredientsModel interface {
	AddIngredientsRecipe(recipeIngredients RecipeIngredients) (RecipeIngredients, error)
	// GetIngredientsByRecipeId(recipeId int) (RecipeIngredients, error)
}

func (m *GormRecipeIngredientsModel) AddIngredientsRecipe(recipeIngredients RecipeIngredients) (RecipeIngredients, error) {
	if err := m.db.Save(&recipeIngredients).Error; err != nil {
		return recipeIngredients, err
	}
	return recipeIngredients, nil
}
