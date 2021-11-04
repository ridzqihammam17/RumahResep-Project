package recipeingredients

import (
	"net/http"
	"rumah_resep/models"

	"github.com/labstack/echo/v4"
)

type RecipeIngredientsController struct {
	recipeIngredientsModel models.RecipeIngredientsModel
	recipeModel            models.RecipeModel
	ingredientModel        models.IngredientModel
}

func NewRecipeIngredientsController(recipeIngredientsModel models.RecipeIngredientsModel, recipeModel models.RecipeModel, ingredientModel models.IngredientModel) *RecipeIngredientsController {
	return &RecipeIngredientsController{
		recipeIngredientsModel,
		recipeModel,
		ingredientModel,
	}
}

type Response struct {
	Data    interface{} `json:"menu"`
	Message string      `json:"message"`
}

func (controller *RecipeIngredientsController) AddIngredientsRecipeController(c echo.Context) error {
	var recipeIngredients models.RecipeIngredients
	c.Bind(&recipeIngredients)
	// fmt.Println(c.Param("categoryId"))
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, "Bad Request")
	// }
	categoryItem := models.RecipeIngredients{
		RecipeId:     recipeIngredients.RecipeId,
		IngredientId: recipeIngredients.IngredientId,
	}
	addIngredient, err := controller.recipeIngredientsModel.AddIngredientsRecipe(categoryItem)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    addIngredient,
		"message": "Success Add Recipe Ingredient",
	})

}
