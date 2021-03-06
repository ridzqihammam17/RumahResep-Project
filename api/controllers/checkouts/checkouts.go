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
	checkoutModel         models.CheckoutModel
	stockModel            models.StockModel
	reciepingredientModel models.RecipeIngredientsModel
	cartModel             models.CartModel
}

func NewCheckoutController(checkoutModel models.CheckoutModel, stockModel models.StockModel, recipeIngredientModel models.RecipeIngredientsModel, cartModel models.CartModel) *CheckoutController {
	return &CheckoutController{
		checkoutModel,
		stockModel,
		recipeIngredientModel,
		cartModel,
	}
}

func (controller *CheckoutController) CreateCheckoutController(c echo.Context) error {
	recipeIdArr := strings.Split(c.Param("recipeId"), ",")
	var recipeId []int
	for _, v := range recipeIdArr {
		value, err := strconv.Atoi(v)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"code":    400,
				"message": "Bad Request",
			})
		}
		recipeId = append(recipeId, value)
	}

	userId, role := middlewares.ExtractTokenUser(c)
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	var checkout models.Checkout
	output, err := controller.checkoutModel.CreateCheckout(checkout)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	cartId, _ := controller.cartModel.GetCartIdByUserId(int(userId))

	for _, v := range recipeId {
		if _, err := controller.checkoutModel.UpdateCheckoutIdOnCartDetails(v, int(output.ID), cartId); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"code":    500,
				"message": "Internal Server Error",
			})
		}
	}

	for i, v := range recipeId {
		mapIngredientIdQty, err := controller.reciepingredientModel.GetIdIngredientQtyIngredient(v)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"code":    500,
				"message": "Internal Server Error",
			})
		}
		_, err = controller.stockModel.StockDecrease(mapIngredientIdQty[i].QtyIngredient, int(mapIngredientIdQty[i].IngredientId))
		if err != nil {
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
