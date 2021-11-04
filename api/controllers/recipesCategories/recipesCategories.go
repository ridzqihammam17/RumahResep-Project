package recipescategories

import (
	"fmt"
	"net/http"
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

// type RecipesCategories struct {
// 	// gorm.Model
// 	RecipeId   int `json:"recipes_id" form:"recipes_id"`
// 	CategoryId int `json:"categories_id" form:"categories_id"`
// 	// CreatedAt  time.Time
// 	// DeletedAt  gorm.DeletedAt
// }

func (controller *RecipesCategoriesController) AddRecipeCategoriesController(c echo.Context) error {
	var recipeCategories models.RecipeCategories
	c.Bind(&recipeCategories)
	// fmt.Println(c.Param("categoryId"))
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, "Bad Request")
	// }
	categoryItem := models.RecipeCategories{
		RecipeId:   recipeCategories.RecipeId,
		CategoryId: recipeCategories.CategoryId,
	}
	addCategory, err := controller.recipesCategoriesModel.AddRecipeCategories(categoryItem)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    addCategory,
		"message": "Success Add Recipe Category",
	})

}

func (controller *RecipesCategoriesController) GetRecipeByCategoryIdController(c echo.Context) error {
	categoryId := strings.Split(c.Param("categoryId"), ",")
	var categoryName []int
	for _, v := range categoryId {
		value, _ := strconv.Atoi(v)
		categoryName = append(categoryName, value)
	}
	fmt.Println(categoryName)
	recipe, err := controller.recipesCategoriesModel.GetRecipeByCategoryId(categoryName)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    recipe,
		"message": "Success Get Recipe By Category ID",
	})
}
