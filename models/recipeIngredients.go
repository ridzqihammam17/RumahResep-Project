package models

import (
	"gorm.io/gorm"
)

type RecipeIngredients struct {
	gorm.Model
	QtyIngredient int  `json:"qty_ingredient" form:"qty_ingredient"`
	RecipeId      uint `gorm:"primaryKey" json:"recipe_id" form:"recipe_id"`
	IngredientId  uint `gorm:"primaryKey" json:"ingredient_id" form:"ingredient_id"`

	// ;foreignKey:ingredient_id;joinForeignKey:ingredientsID;References:recipe_id;joinReferences:recipesID"
}

type GormRecipeIngredientsModel struct {
	db *gorm.DB
}

func NewRecipeIngredientsModel(db *gorm.DB) *GormRecipeIngredientsModel {
	db.Migrator().AddColumn(&RecipeIngredients{}, "Harga")
	return &GormRecipeIngredientsModel{db: db}
}

type RecipeIngredientsModel interface {
	AddIngredientsRecipe(recipeIngredients RecipeIngredients) (RecipeIngredients, error)
	GetIdIngredientQtyIngredient(recipeId int) ([]RecipeIngredients, error)
}

func (m *GormRecipeIngredientsModel) AddIngredientsRecipe(recipeIngredients RecipeIngredients) (RecipeIngredients, error) {
	if err := m.db.Save(&recipeIngredients).Error; err != nil {
		return recipeIngredients, err
	}
	return recipeIngredients, nil
}

func (m *GormRecipeIngredientsModel) GetIdIngredientQtyIngredient(recipeId int) ([]RecipeIngredients, error) {
	var mapIngredientIdQty []RecipeIngredients
	if err := m.db.Raw("SELECT ingredient_id, qty_ingredient FROM recipe_ingredients WHERE recipe_id = ?", recipeId).Scan(&mapIngredientIdQty).Error; err != nil {
		return mapIngredientIdQty, err
	}

	return mapIngredientIdQty, nil
}
