package stock

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"

	"github.com/labstack/echo/v4"
)

type StockController struct {
	stockModel models.StockModel
}

func NewStockController(stockModel models.StockModel) *StockController {
	return &StockController{
		stockModel,
	}
}

func (controller *StockController) GetRestockDate(c echo.Context) error {
	daterange := c.Param("range")
	var user models.User
	//check role admin or not
	_, user.Role = middlewares.ExtractTokenUser(c)
	if user.Role != "admin" {
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
