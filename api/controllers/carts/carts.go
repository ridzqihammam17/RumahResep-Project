package carts

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"

	"github.com/labstack/echo/v4"
)

type CartController struct {
	cartModel models.CartModel
}

func NewCartController(cartModel models.CartModel) *CartController {
	return &CartController{
		cartModel,
	}
}

func (controller *CartController) CreateCartController(c echo.Context) error {
	var cart models.Cart
	c.Bind(&cart)

	// get id user & role login
	userId, role := middlewares.ExtractTokenUser(c)

	// check role is customer
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	//set data cart and create new cart
	cart = models.Cart{
		UserID: userId,
	}

	_, err := controller.cartModel.CreateCart(cart, int(userId))
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
		"message": "Create cart success",
	})

}
