package ingredients

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

type CreateIngredientRequest struct {
	Name  string `json:"name" form:"name"`
	Price int    `json:"price" form:"price"`
}

type UpdateIngredientRequest struct {
	Name  string `json:"name" form:"name"`
	Price int    `json:"price" form:"price"`
}

// ------------------------------------------------------------------
// End Request
// ------------------------------------------------------------------

type IngredientController struct {
	ingredientModel models.IngredientModel
}

func NewIngredientController(ingredientModel models.IngredientModel) *IngredientController {
	return &IngredientController{
		ingredientModel,
	}
}

func (controller *IngredientController) GetAllIngredientController(c echo.Context) error {
	//check role admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	data, err := controller.ingredientModel.GetAllIngredient()
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
		"message": "Success Get All Ingredient",
		"data":    data,
	})
}

func (controller *IngredientController) GetIngredientByIdController(c echo.Context) error {
	//check role admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	data, err := controller.ingredientModel.GetIngredientById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Get Ingredient",
		"data":    data,
	})
}

func (controller *IngredientController) CreateIngredientController(c echo.Context) error {
	//check role admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	var ingredientRequest CreateIngredientRequest

	c.Bind(&ingredientRequest)

	ingredient := models.Ingredient{
		Name:  ingredientRequest.Name,
		Price: ingredientRequest.Price,
	}
	if ingredient.Name == "" || ingredient.Price < 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	data, err := controller.ingredientModel.CreateIngredient(ingredient)
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
		"message": "Success Create Ingredient",
		"data":    data,
	})
}

func (controller *IngredientController) UpdateIngredientController(c echo.Context) error {
	// check admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	var ingredientRequest UpdateIngredientRequest

	c.Bind(&ingredientRequest)

	ingredient := models.Ingredient{
		Name:  ingredientRequest.Name,
		Price: ingredientRequest.Price,
	}
	if ingredient.Name == "" || ingredient.Price < 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	data, err := controller.ingredientModel.UpdateIngredient(ingredient, id)
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
		"message": "Success Update Ingredient",
		"data":    data,
	})
}

func (controller *IngredientController) DeleteIngredientController(c echo.Context) error {
	// check admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	if _, err := controller.ingredientModel.DeleteIngredient(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Delete Ingredient",
	})
}
