package models

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Name string `json:"name" form:"name"`

	// One to Many with Recipe Category
	RecipeCategories []RecipeCategories
	// One to Many with Recipe Ingredient
	RecipeIngredients []RecipeIngredients
}

type GormRecipeModel struct {
	db *gorm.DB
}

func NewRecipeModel(db *gorm.DB) *GormRecipeModel {
	return &GormRecipeModel{db: db}
}

type RecipeModel interface {
	CreateRecipe(recipe Recipe) (Recipe, error)
	GetAllRecipe() ([]Recipe, error)
	GetRecipeById(recipeId int) (Recipe, error)
	UpdateRecipe(recipe Recipe, recipeId int) (Recipe, error)
	DeleteRecipe(recipeId int) (Recipe, error)
}

func (m *GormRecipeModel) CreateRecipe(recipe Recipe) (Recipe, error) {
	if err := m.db.Save(&recipe).Error; err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (m *GormRecipeModel) GetAllRecipe() ([]Recipe, error) {
	var recipe []Recipe
	if err := m.db.Find(&recipe).Error; err != nil {
		return nil, err
	}
	return recipe, nil
}

func (m *GormRecipeModel) GetRecipeById(recipeId int) (Recipe, error) {
	var recipe Recipe
	if err := m.db.First(&recipe, recipeId).Error; err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (m *GormRecipeModel) UpdateRecipe(recipe Recipe, recipeId int) (Recipe, error) {
	var newRecipe Recipe
	if err := m.db.First(&newRecipe, recipeId).Error; err != nil {
		return recipe, err
	}

	newRecipe.Name = recipe.Name

	if err := m.db.Save(&newRecipe).Error; err != nil {
		return newRecipe, err
	}
	return newRecipe, nil
}

func (m *GormRecipeModel) DeleteRecipe(recipeId int) (Recipe, error) {
	var recipe Recipe
	if err := m.db.First(&recipe, recipeId).Error; err != nil {
		return recipe, err
	}

	if err := m.db.Delete(&recipe).Error; err != nil {
		return recipe, err
	}
	return recipe, nil
}
