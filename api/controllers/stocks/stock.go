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
			"message": "Unauthorized Error",
		})
	}

	var stock models.Stock
	c.Bind(&stock)

	ingredientId, err := strconv.Atoi(c.Param("ingredientId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	newStock := models.Stock{
		IngredientId: uint(ingredientId),
		Stock:        stock.Stock,
		UserId:       userId,
	}
	if newStock.Stock == 0 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
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
		"message": "Success Create Stock Ingredient",
		"data":    output,
	})

}

func (controller *StockController) GetRestockDateController(c echo.Context) error {
	//check role seller or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "seller" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	daterange := c.Param("range")
	if daterange != "daily" && daterange != "weekly" && daterange != "monthly" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	stock, err := controller.stockModel.GetRestockDate(int(userId), daterange)
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

func (controller *StockController) GetRestockAllController(c echo.Context) error {
	//check role seller or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "seller" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	stock, _ := controller.stockModel.GetRestockAll(int(userId))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Get All Restock",
		"data":    stock,
	})
}
