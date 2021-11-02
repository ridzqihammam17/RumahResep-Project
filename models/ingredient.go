package models

import "gorm.io/gorm"

type Ingredient struct {
	gorm.Model
	Name  string `json:"name" form:"name"`
	Price int    `json:"price" form:"price"`
	Stock int    `json:"stock" form:"stock"`
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
