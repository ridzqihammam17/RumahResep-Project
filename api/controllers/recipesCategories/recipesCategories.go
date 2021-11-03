package recipescategories

import (
	"net/http"
	"rumah_resep/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type RecipesCategoriesController struct {
	recipesCategoriesModel models.RecipesCategoriesModel
}

func NewRecipesCategoriesController(recipesCategoriesModel models.RecipesCategoriesModel) *RecipesCategoriesController {
	return &RecipesCategoriesController{
		recipesCategoriesModel,
	}
}

func (controller *RecipesCategoriesController) GetRecipeByCategoryIdController(c echo.Context) error {
	categoryId := strings.Split(c.Param("categoryId"), ",")
	var categoryName []int
	for _, v := range categoryId {
		value, _ := strconv.Atoi(v)
		categoryName = append(categoryName, value)
	}

	// fmt.Println(c.Param("categoryId"))
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, "Bad Request")
	// }

	recipe, err := controller.recipesCategoriesModel.GetRecipeByCategoryId(categoryName)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    recipe,
		"message": "Success Get Recipe By Category ID",
	})
}
