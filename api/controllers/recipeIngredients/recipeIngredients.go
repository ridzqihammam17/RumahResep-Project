package recipeingredients

import (
	"net/http"
	"rumah_resep/models"
	"strconv"

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

func (controller *RecipeIngredientsController) AddIngredientsRecipeController(c echo.Context) error {
	var recipeIngredients models.RecipeIngredients
	if err := c.Bind(&recipeIngredients); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	categoryItem := models.RecipeIngredients{
		RecipeId:      recipeIngredients.RecipeId,
		IngredientId:  recipeIngredients.IngredientId,
		QtyIngredient: recipeIngredients.QtyIngredient,
	}
	addIngredient, err := controller.recipeIngredientsModel.AddIngredientsRecipe(categoryItem)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"data":    addIngredient,
		"message": "Success Add Recipe Ingredient",
	})

}

func (controller *RecipeIngredientsController) GetIngredientsByRecipeIdController(c echo.Context) error {
	recipeId, _ := strconv.Atoi(c.Param("recipeId"))
	ingredient, err := controller.ingredientModel.GetIngredientsByRecipeId(recipeId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"succuess": false,
			"code":     400,
			"message":  "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"data":    ingredient,
		"message": "Success Get Ingredients By Recipe ID",
	})
}
