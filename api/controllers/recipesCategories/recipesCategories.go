package recipescategories

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type RecipesCategoriesController struct {
	recipesCategoriesModel models.RecipeCategoriesModel
	recipeModel            models.RecipeModel
	categoryModel          models.CategoryModel
}

func NewRecipesCategoriesController(
	recipesCategoriesModel models.RecipeCategoriesModel,
	recipeModel models.RecipeModel,
	categoryModel models.CategoryModel) *RecipesCategoriesController {
	return &RecipesCategoriesController{
		recipesCategoriesModel,
		recipeModel,
		categoryModel,
	}
}

func (controller *RecipesCategoriesController) AddRecipeCategoriesController(c echo.Context) error {
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	var recipeCategories models.RecipeCategories
	c.Bind(&recipeCategories)

	categoryItem := models.RecipeCategories{
		RecipeId:   recipeCategories.RecipeId,
		CategoryId: recipeCategories.CategoryId,
	}

	if categoryItem.RecipeId == 0 || categoryItem.CategoryId == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	output, err := controller.recipesCategoriesModel.AddRecipeCategories(categoryItem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    output,
		"code":    200,
		"message": "Success Add Recipe Category",
	})

}

func (controller *RecipesCategoriesController) GetRecipeByCategoryIdController(c echo.Context) error {
	categoryId := strings.Split(c.Param("categoryId"), ",")
	var categoryName []int
	for _, v := range categoryId {
		value, err := strconv.Atoi(v)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"code":    400,
				"message": "Bad Request",
			})
		}
		categoryName = append(categoryName, value)
	}
	recipe, err := controller.recipesCategoriesModel.GetRecipeByCategoryId(categoryName)
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
		"data":    recipe,
		"message": "Success Get Recipe By Category ID",
	})
}
