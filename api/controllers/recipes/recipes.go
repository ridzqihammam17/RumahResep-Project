package recipes

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RecipeController struct {
	recipeModel models.RecipeModel
}

func NewRecipeController(recipeModel models.RecipeModel) *RecipeController {
	return &RecipeController{
		recipeModel,
	}
}

func (controller *RecipeController) CreateRecipeController(c echo.Context) error {
	//bind recipe from request body
	var recipe models.Recipe
	if err := c.Bind(&recipe); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	//check role admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Role not Admin",
		})
	}

	// input request body to newRecipe
	newRecipe := models.Recipe{
		Name:     recipe.Name,
		Category: recipe.Category,
	}

	// create recipe
	output, err := controller.recipeModel.CreateRecipe(newRecipe)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"recipe":  output,
		"message": "Create Recipe Success",
	})
}

func (controller *RecipeController) GetAllRecipeController(c echo.Context) error {
	recipes, err := controller.recipeModel.GetAllRecipe()
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, recipes)
}

func (controller *RecipeController) GetRecipeByIdController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("recipeId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	recipe, err := controller.recipeModel.GetRecipeById(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, recipe)
}

func (controller *RecipeController) UpdateRecipeController(c echo.Context) error {
	// check admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Role not Admin",
		})
	}

	// check id recipe
	id, err := strconv.Atoi(c.Param("recipeId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	//bind recipe from request body
	var recipeRequest models.Recipe
	if err := c.Bind(&recipeRequest); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	recipe := models.Recipe{
		Name:     recipeRequest.Name,
		Category: recipeRequest.Category,
	}

	output, err := controller.recipeModel.UpdateRecipe(recipe, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Data Not Found")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    output,
		"message": "Success Update Recipe",
	})
}

func (controller *RecipeController) DeleteRecipeController(c echo.Context) error {
	// check admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Role not Admin",
		})
	}

	// check id recipe
	id, err := strconv.Atoi(c.Param("recipeId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if _, err := controller.recipeModel.DeleteRecipe(id); err != nil {
		return c.String(http.StatusNotFound, "Data Not Found")
	}
	return c.String(http.StatusOK, "Success Delete Recipe")
}

func (controller *RecipeController) GetRecipeByCategoryIdController(c echo.Context) error {
	categoryId, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	recipe, err := controller.recipeModel.GetRecipeByCategoryId(categoryId)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    recipe,
		"message": "Success Get Recipe By Category ID",
	})
}
