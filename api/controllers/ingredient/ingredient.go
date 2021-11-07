package ingredient

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"
	"strconv"

	"github.com/labstack/echo/v4"
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
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "seller" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
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
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	newStock := models.Stock{
		IngredientId: ingredient.ID,
		Stock:        ingredient.Stock,
		UserId:       userId,
	}

	if _, err := controller.StockModel.CreateStockUpdate(newStock, int(ingredient.ID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"data":    output,
		"message": "Create Ingredient Success",
	})
}

func (controller *IngredientController) GetAllIngredientController(c echo.Context) error {
	ingredient, err := controller.IngredientModel.GetAllIngredient()
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
		"data":    ingredient,
		"message": "Get All Ingredient Success",
	})
}

func (controller *IngredientController) GetIngredientByIdController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	ingredient, err := controller.IngredientModel.GetIngredientById(id)
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
		"data":    ingredient,
		"message": "Get Ingredient By Id Success",
	})
}

func (controller *IngredientController) UpdateIngredientController(c echo.Context) error {
	// check admin or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "seller" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	// check id recipe
	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	//bind recipe from request body
	var ingredientRequest models.Ingredient
	if err := c.Bind(&ingredientRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}
	ingredient := models.Ingredient{
		Name:   ingredientRequest.Name,
		Price:  ingredientRequest.Price,
		Stock:  ingredientRequest.Stock,
		UserID: userId,
	}

	output, err := controller.IngredientModel.UpdateIngredient(ingredient, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"code":    404,
			"message": "Not Found",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"data":    output,
		"message": "Update Ingredient Success",
	})
}

// func (controller *IngredientController) UpdateIngredientStockController(c echo.Context) error {
// 	// check admin or not
// 	userId, role := middlewares.ExtractTokenUser(c)
// 	if role != "seller" {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"success": false,
// 			"code":    401,
// 			"message": "Unauthorized",
// 		})
// 	}

// 	// check id recipe
// 	ingredientId, err := strconv.Atoi(c.Param("ingredientId"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]interface{}{
// 			"success": false,
// 			"code":    400,
// 			"message": "Bad Request",
// 		})
// 	}

// 	//bind recipe from request body
// 	var stockRequest models.Stock
// 	if err := c.Bind(&stockRequest); err != nil {
// 		return c.String(http.StatusBadRequest, "Bad Request")
// 	}
// 	// var ingredientRequest models.Ingredient
// 	// if err := c.Bind(&ingredientRequest); err != nil {
// 	// 	return c.JSON(http.StatusBadRequest, map[string]interface{}{
// 	// 		"success": false,
// 	// 		"code":    400,
// 	// 		"message": "Bad Request",
// 	// 	})
// 	// }
// 	// ingredient := models.Ingredient{
// 	// 	Stock: ingredientRequest.Stock,
// 	// }

// 	// if err := c.Bind(&ingredientRequest); err != nil {
// 	// 	return c.String(http.StatusBadRequest, "Bad Request")
// 	// }
// 	stock := models.Stock{
// 		UserId:       userId,
// 		IngredientId: uint(ingredientId),
// 		Stock:        stockRequest.Stock,
// 	}

// 	// _, err2 := controller.IngredientModel.UpdateStock(ingredient, ingredientId)
// 	// if err2 != nil {
// 	// 	return c.JSON(http.StatusNotFound, map[string]interface{}{
// 	// 		"success": false,
// 	// 		"code":    404,
// 	// 		"message": "Not Found",
// 	// 	})
// 	// }

// 	output2, err3 := controller.StockModel.CreateStockUpdate(stock, ingredientId)
// 	if err3 != nil {
// 		return c.JSON(http.StatusNotFound, map[string]interface{}{
// 			"success": false,
// 			"code":    404,
// 			"message": "Not Found",
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"success": true,
// 		"code":    200,
// 		"data":    output2,
// 		"message": "Update Ingredient Stock Success",
// 	})
// }

func (controller *IngredientController) DeleteIngredientController(c echo.Context) error {
	// check admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "seller" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	// check id recipe
	id, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	if _, err := controller.IngredientModel.DeleteIngredient(id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"code":    404,
			"message": "Not Found",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Delete Ingredient",
	})
}
