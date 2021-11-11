package midtrans

import (
	"net/http"
	"rumah_resep/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MidtransController struct {
	// midtransModel models.MidtransMo
	transactionModel models.TransactionModel
}

func NewMidtransController(transactionModel models.TransactionModel) *MidtransController {
	return &MidtransController{
		transactionModel,
	}
}

func (controller *MidtransController) RequestPayment(c echo.Context) error {
	ids, _ := strconv.Atoi(c.Param("id"))

	totalPayment, _ := controller.transactionModel.GetTotalPayment(ids)

	redirectURL, err := models.RequestPayment(ids, totalPayment)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error Bosqu")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Request Payment",
		"data":    redirectURL,
	})
}

func (controller *MidtransController) StatusPayment(c echo.Context) error {
	ids := c.Param("id")
	// idStr, err := strconv.Atoi(id)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "ID Invalid")
	// }
	redirectURL, err := models.StatusPayment(ids)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error Bosqu")
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Get Status Payment",
		"data":    redirectURL,
	})
}
