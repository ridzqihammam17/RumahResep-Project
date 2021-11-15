package midtrans

import (
	"net/http"
	"rumah_resep/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type MidtransController struct {
	transactionModel models.TransactionModel
}

func NewMidtransController(transactionModel models.TransactionModel) *MidtransController {
	return &MidtransController{
		transactionModel,
	}
}

func (controller *MidtransController) RequestPayment(c echo.Context) error {
	idSplit := strings.Split(c.Param("id"), "-")
	ids, _ := strconv.Atoi(idSplit[1])

	totalPayment, _ := controller.transactionModel.GetTotalPayment(ids)

	data, err := models.RequestPayment(ids, totalPayment)
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
		"message": "Success Request Payment",
		"data":    data,
	})
}

func (controller *MidtransController) StatusPayment(c echo.Context) error {
	ids := c.Param("id")
	data, err := models.StatusPayment(ids)
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
		"message": "Success Get Status Payment",
		"data":    data,
	})
}
