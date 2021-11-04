package models

import "gorm.io/gorm"

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

type RecipeIngredientsModel struct {
	GetIngredientsByRecipeId(recipeId int) ([]RecipeIngredients, error) 
}

func (m *GormRecipeIngredientsModel) GetIngredientsByRecipeId(recipeId int) ([]RecipeIngredients, error) {
	var recipeIngredients []RecipeIngredients

	if err := m.db.Find(&recipeIngredients, "recipe_id", recipeId).Error; err != nil {
		return recipe_ingredients, err
	}
	if len(recipe_ingredients) == 0 {

		return nil, errors.New("Data Not Found")
	}
	return recipe_ingredients, nil
}