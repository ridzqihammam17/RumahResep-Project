package stock

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type StockController struct {
	stockModel models.StockModel
	// ingredientModel models.IngredientModel
}

func NewStockController(stockModel models.StockModel) *StockController {
	return &StockController{
		stockModel,
		// ingredientModel,
	}
}

func (controller *StockController) CreateStockUpdateController(c echo.Context) error {

	//check role seller or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "seller" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	var stock models.Stock
	if err := c.Bind(&stock); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}
	ingredientId, _ := strconv.Atoi(c.Param("ingredientId"))
	newStock := models.Stock{
		IngredientId: uint(ingredientId),
		Stock:        stock.Stock,
		UserId:       userId,
	}

	output, err := controller.stockModel.CreateStockUpdate(newStock, ingredientId)
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
		"message": "Success Restock Ingredient",
		"data":    output,
	})

}

func (controller *StockController) UpdateStockController(c echo.Context) error {
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "seller" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	var stock models.Stock
	if err := c.Bind(&stock); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}
	ingredientId, _ := strconv.Atoi(c.Param("ingredientId"))
	newStock := models.Stock{
		IngredientId: uint(ingredientId),
		Stock:        stock.Stock,
		UserId:       userId,
	}

	output, err := controller.stockModel.UpdateStock(newStock, ingredientId, int(userId))
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
		"message": "Success Restock Ingredient",
		"data":    output,
	})
}

func (controller *StockController) GetRestockDateController(c echo.Context) error {
	daterange := c.Param("range")
	if daterange != "daily" && daterange != "weekly" && daterange != "monthly" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid range")
	}
	// var user models.User
	//check role admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "seller" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	stock, err := controller.stockModel.GetRestockDate(daterange)
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
		"message": "Success Get Restock Date",
		"data":    stock,
	})
}
