package models

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Name string `json:"name" form:"name"`

	// many2many with category
	Categories []*Category `gorm:"many2many:recipe_categories" json:"categories"`
	// many2many with ingredient
	Ingredients []Ingredient `gorm:"many2many:recipe_ingredients" json:"ingredients"`

	// Category int `json:"category" form:"category"`
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
	// GetRecipeByCategoryId(categoryId []int) ([]Recipe, error)
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
	if err := m.db.Find(&recipe, recipeId).Error; err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (m *GormRecipeModel) UpdateRecipe(recipe Recipe, recipeId int) (Recipe, error) {
	var newRecipe Recipe
	if err := m.db.Find(&newRecipe, recipeId).Error; err != nil {
		return recipe, err
	}

	newRecipe.Name = recipe.Name
	// newRecipe.Category = recipe.Category

	if err := m.db.Save(&newRecipe).Error; err != nil {
		return newRecipe, err
	}
	return newRecipe, nil
}

func (m *GormRecipeModel) DeleteRecipe(recipeId int) (Recipe, error) {
	var recipe Recipe
	if err := m.db.Find(&recipe, recipeId).Error; err != nil {
		return recipe, err
	}

	if err := m.db.Delete(&recipe).Error; err != nil {
		return recipe, err
	}
	return recipe, nil
}

// func (m *GormRecipeModel) GetRecipeByCategoryId(categoryId []int) ([]Recipe, error) {
// 	var recipe []Recipe

// 	if err := m.db.Find(&recipe, "category in ?", categoryId).Error; err != nil {
// 		return recipe, err
// 	}
// 	if len(recipe) == 0 {

// 		return nil, errors.New("Data Not Found")
// 	}
// 	return recipe, nil
// }
