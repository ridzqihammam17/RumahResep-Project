package ingredient

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"
	"strconv"
)

type IngredientController struct {
	IngredientModel models.IngredientModel
	StockModel      models.StockModel
}

func NewIngredientController(IngredientModel models.IngredientModel, StockModel models.StockModel) *IngredientController {
	// cartDetailModel models.CartDetailModel, recipeModel models.RecipeModel
	return &IngredientController{
		IngredientModel,
		StockModel,
	}
}

func (controller *IngredientController) CreateIngredientController(c echo.Context) error {
	var ingredient models.Ingredient
	if err := c.Bind(&ingredient); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	//check role admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Role is not Admin",
		})
	}

	// input request body to newRecipe
	newIngredient := models.Ingredient{
		Name:  ingredient.Name,
		Price: ingredient.Price,
		Stock: ingredient.Stock,
	}

	// create recipe
	output, err := controller.IngredientModel.CreateIngredient(newIngredient)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"recipe":  output,
		"message": "Create Ingredient Success",
	})
}

func (controller *IngredientController) GetAllIngredientController(c echo.Context) error {
	ingredient, err := controller.IngredientModel.GetAllIngredient()
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, ingredient)
}

func (controller *IngredientController) GetIngredientByIdController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	ingredient, err := controller.IngredientModel.GetIngredientById(id)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.JSON(http.StatusOK, ingredient)
}

func (controller *IngredientController) UpdateIngredientController(c echo.Context) error {
	// check admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Role is not Admin",
		})
	}

	// check id recipe
	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	//bind recipe from request body
	var ingredientRequest models.Ingredient
	if err := c.Bind(&ingredientRequest); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	ingredient := models.Ingredient{
		Name:  ingredientRequest.Name,
		Price: ingredientRequest.Price,
	}

	output, err := controller.IngredientModel.UpdateIngredient(ingredient, id)
	if err != nil {
		return c.String(http.StatusNotFound, "Data Not Found")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    output,
		"message": "Success Update Recipe",
	})
}

func (controller *IngredientController) UpdateIngredientStockController(c echo.Context) error {
	// check admin or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Role is not Admin",
		})
	}

	// check id recipe
	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	//bind recipe from request body
	var stockRequest models.Stock
	var ingredientRequest models.Ingredient

	if err := c.Bind(&ingredientRequest); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	ingredient := models.Ingredient{
		Stock: ingredientRequest.Stock,
	}

	if err := c.Bind(&stockRequest); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	stock := models.Stock{
		UserId: userId,
		Stock:  ingredientRequest.Stock,
	}

	_, err2 := controller.IngredientModel.UpdateStock(ingredient, id)
	if err2 != nil {
		return c.String(http.StatusNotFound, "Data Not Found")
	}

	output2, err3 := controller.StockModel.CreateStockUpdate(stock)
	if err3 != nil {
		return c.String(http.StatusNotFound, "Data Not Found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    output2,
		"message": "Success Update Recipe",
	})
}
func (controller *IngredientController) DeleteIngredientController(c echo.Context) error {
	// check admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Role is not Admin",
		})
	}

	// check id recipe
	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if _, err := controller.IngredientModel.DeleteIngredient(id); err != nil {
		return c.String(http.StatusNotFound, "Data Not Found")
	}
	return c.String(http.StatusOK, "Success Delete Recipe")
}
