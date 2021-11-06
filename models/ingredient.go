package models

import (
	"errors"

	"gorm.io/gorm"
)

type Ingredient struct {
	gorm.Model
	Name  string `json:"name" form:"name"`
	Price int    `json:"price" form:"price"`
	Stock int    `json:"stock" form:"stock"`

	Stocks []Stock

	// One to Many with Recipe Ingredient
	RecipeIngredients []RecipeIngredients
}

type GormIngredientModel struct {
	db *gorm.DB
}

func NewIngredientModel(db *gorm.DB) *GormIngredientModel {
	return &GormIngredientModel{db: db}
}

type IngredientModel interface {
	CreateIngredient(ingredient Ingredient) (Ingredient, error)
	GetAllIngredient() ([]Ingredient, error)
	GetIngredientById(ingredientId int) (Ingredient, error)
	UpdateIngredient(ingredient Ingredient, ingredientId int) (Ingredient, error)
	DeleteIngredient(ingredientId int) (Ingredient, error)
	UpdateStock(ingredient Ingredient, ingredientId int) (Ingredient, error)
	GetIngredientsByRecipeId(recipeId int) ([]Ingredient, error)
	GetIngredientPrice(ingredientId int) (int, error)
}

func (m *GormIngredientModel) CreateIngredient(ingredient Ingredient) (Ingredient, error) {
	if err := m.db.Save(&ingredient).Error; err != nil {
		return ingredient, err
	}
	return ingredient, nil
}

func (m *GormIngredientModel) GetAllIngredient() ([]Ingredient, error) {
	var ingredient []Ingredient
	if err := m.db.Find(&ingredient).Error; err != nil {
		return nil, err
	}
	return ingredient, nil
}

func (m *GormIngredientModel) GetIngredientById(ingredientId int) (Ingredient, error) {
	var ingredient Ingredient
	if err := m.db.Find(&ingredient, ingredientId).Error; err != nil {
		return ingredient, err
	}
	return ingredient, nil
}

func (m *GormIngredientModel) UpdateStock(ingredient Ingredient, ingredientId int) (Ingredient, error) {
	var newIngredient Ingredient
	if err := m.db.Find(&newIngredient, ingredientId).Error; err != nil {
		return ingredient, err
	}

	newIngredient.Stock = ingredient.Stock

	if err := m.db.Save(&newIngredient).Error; err != nil {
		return newIngredient, err
	}
	return newIngredient, nil
}

func (m *GormIngredientModel) UpdateIngredient(ingredient Ingredient, ingredientId int) (Ingredient, error) {
	var newIngredient Ingredient
	if err := m.db.Find(&newIngredient, ingredientId).Error; err != nil {
		return ingredient, err
	}

	newIngredient.Name = ingredient.Name
	newIngredient.Price = ingredient.Price

	if err := m.db.Save(&newIngredient).Error; err != nil {
		return newIngredient, err
	}
	return newIngredient, nil
}

func (m *GormIngredientModel) DeleteIngredient(ingredientId int) (Ingredient, error) {
	var ingredient Ingredient
	if err := m.db.Find(&ingredient, ingredientId).Error; err != nil {
		return ingredient, err
	}

	if err := m.db.Delete(&ingredient).Error; err != nil {
		return ingredient, err
	}
	return ingredient, nil
}

func (m *GormIngredientModel) GetIngredientsByRecipeId(recipeId int) ([]Ingredient, error) {
	var recipeIngredients []Ingredient
	if err := m.db.Raw("SELECT i.id, i.name FROM recipe_ingredients rc join ingredients i ON i.id = rc.ingredient_id WHERE rc.recipe_id = ?", recipeId).Scan(&recipeIngredients).Error; err != nil {
		return recipeIngredients, err
	}

	if len(recipeIngredients) == 0 {

		return nil, errors.New("Data Not Found")
	}
	return recipeIngredients, nil
}

func (m *GormIngredientModel) GetIngredientPrice(ingredientId int) (int, error) {
	var ingredientPrice int
	if err := m.db.Raw("SELECT i.price FROM ingredients i  WHERE i.id = ?", ingredientId).Scan(&ingredientPrice).Error; err != nil {
		return ingredientPrice, err
	}
	return ingredientPrice, nil
}
