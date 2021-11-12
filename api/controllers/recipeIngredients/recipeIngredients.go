package recipeingredients

import (
	"net/http"
	"rumah_resep/api/middlewares"
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
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	var recipeIngredients models.RecipeIngredients
	c.Bind(&recipeIngredients)

	categoryItem := models.RecipeIngredients{
		RecipeId:      recipeIngredients.RecipeId,
		IngredientId:  recipeIngredients.IngredientId,
		QtyIngredient: recipeIngredients.QtyIngredient,
	}

	if categoryItem.RecipeId == 0 || categoryItem.IngredientId == 0 || categoryItem.QtyIngredient == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	addIngredient, err := controller.recipeIngredientsModel.AddIngredientsRecipe(categoryItem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"succuess": false,
			"code":     500,
			"message":  "Internal Server Error",
		})
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
