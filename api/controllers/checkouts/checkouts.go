package checkouts

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type CheckoutController struct {
	checkoutModel models.CheckoutModel
}

func NewCheckoutController(checkoutModel models.CheckoutModel) *CheckoutController {
	return &CheckoutController{
		checkoutModel,
	}
}

func (controller *CheckoutController) CreateCheckoutController(c echo.Context) error {
	_, role := middlewares.ExtractTokenUser(c)
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	var checkout models.Checkout
	if err := c.Bind(&checkout); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Success Bad Request",
		})
	}

	output, err := controller.checkoutModel.CreateCheckout(checkout)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	recipeIdArr := strings.Split(c.Param("recipeId"), ",")
	var recipeId []int
	for _, v := range recipeIdArr {
		value, _ := strconv.Atoi(v)
		recipeId = append(recipeId, value)
	}

	for _, v := range recipeId {
		if _, err := controller.checkoutModel.UpdateCheckoutIdOnCartDetails(v, int(output.ID)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"code":    500,
				"message": "Internal Server Error",
			})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Create Checkout",
		"data":    output,
	})
}
