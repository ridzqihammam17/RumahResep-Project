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
	// GetIngredientsByRecipeId(recipeId int) ([]RecipeIngredients, error)
}

func (m *GormRecipeIngredientsModel) AddIngredientsRecipe(recipeIngredients RecipeIngredients) (RecipeIngredients, error) {
	if err := m.db.Save(&recipeIngredients).Error; err != nil {
		return recipeIngredients, err
	}
	return recipeIngredients, nil
}

// func (m *GormIngredientModel) GetIngredientsByRecipeId(recipeId int) (string, error) {
// 	// var recipeIngredients []RecipeIngredients
// 	var recipeIngredients string
// 	// if err := m.db.Select("name").Find(&ingredients).Where()
// 	if err := m.db.Raw("SELECT ingredients.name FROM `ingredients` left join recipe_ingredients ON ingredients.id=recipe_ingredients.ingredient_id WHERE recipe_ingredients.recipe_id = ?", "recipeId", recipeId).Error; err != nil {
// 		return recipeIngredients, err
// 	}
// 	fmt.Println(recipeIngredients)

// 	// if err := m.db.Find(&recipeIngredients, "recipe_id", recipeId).Error; err != nil {

// 	// }
// 	if len(recipeIngredients) == 0 {

// 		return "", errors.New("Data Not Found")
// 	}
// 	return recipeIngredients, nil
// }
