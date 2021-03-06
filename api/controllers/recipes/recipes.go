package recipes

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ------------------------------------------------------------------
// Start Request
// ------------------------------------------------------------------

type CreateRecipeRequest struct {
	Name string `json:"name" form:"name"`
}

type UpdateRecipeRequest struct {
	Name string `json:"name" form:"name"`
}

// ------------------------------------------------------------------
// End Request
// ------------------------------------------------------------------

type RecipeController struct {
	recipeModel models.RecipeModel
}

func NewRecipeController(recipeModel models.RecipeModel) *RecipeController {
	return &RecipeController{
		recipeModel,
	}
}

func (controller *RecipeController) GetAllRecipeController(c echo.Context) error {
	data, err := controller.recipeModel.GetAllRecipe()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Get All Recipe",
		"data":    data,
	})
}

func (controller *RecipeController) GetRecipeByIdController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("recipeId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	data, err := controller.recipeModel.GetRecipeById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Get Recipe",
		"data":    data,
	})
}

func (controller *RecipeController) CreateRecipeController(c echo.Context) error {
	//check role admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	var recipeRequest CreateRecipeRequest

	c.Bind(&recipeRequest)

	recipe := models.Recipe{
		Name: recipeRequest.Name,
	}
	if recipe.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	data, err := controller.recipeModel.CreateRecipe(recipe)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Create Recipe",
		"data":    data,
	})
}

func (controller *RecipeController) UpdateRecipeController(c echo.Context) error {
	// check admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	id, err := strconv.Atoi(c.Param("recipeId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	var recipeRequest UpdateRecipeRequest

	c.Bind(&recipeRequest)

	recipe := models.Recipe{
		Name: recipeRequest.Name,
	}
	if recipe.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	data, err := controller.recipeModel.UpdateRecipe(recipe, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Update Recipe",
		"data":    data,
	})
}

func (controller *RecipeController) DeleteRecipeController(c echo.Context) error {
	// check admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	id, err := strconv.Atoi(c.Param("recipeId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	if _, err := controller.recipeModel.DeleteRecipe(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Delete Recipe",
	})
}
