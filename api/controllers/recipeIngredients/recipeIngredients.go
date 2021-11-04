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

	// type Response struct {
	// 	Data    models.Ingredient `json:"data"`
	// 	Message string            `json:"message"`
	// }

	// response := &Response{
	// 	Data:    ingredient,
	// 	Message: "Success Get Ingredient By Recipe Id",
	// }

	// var response Response
	// fmt.Println(c.Param("recipeId"))
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"data":    ingredient,
		"message": "Success Get Recipe By Category ID",
	})
}
