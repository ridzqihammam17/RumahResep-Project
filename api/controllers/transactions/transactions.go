package transactions

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"

	// "rumah_resep/models"

	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	transactionModel models.TransactionModel
	cartModel        models.CartModel
	userModel        models.UserModel
}

func NewTransactionController(transactionModel models.TransactionModel, cartModel models.CartModel, userModel models.UserModel) *TransactionController {
	return &TransactionController{
		transactionModel,
		cartModel,
		userModel,
	}
}

func (controller *TransactionController) CreateTransaction(c echo.Context) error {
	// Check role is Customer
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	// Get Cart Id
	cartId, err := controller.cartModel.GetCartIdByUserId(int(userId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	// Get Checkout Id
	checkoutId, err := controller.transactionModel.GetCheckoutId(cartId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	// Get User Data
	user, err := controller.userModel.GetUserData(int(userId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	// Get Total Payment
	totalPayment, err := controller.transactionModel.CountTotalPayment(cartId, checkoutId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	// Bind Request Body
	var transaction models.Transaction
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	// Prepare Response Body
	newTransaction := models.Transaction{
		UserID:         userId,
		CustomerName:   user.Name,
		Address:        user.Address,
		ShippingMethod: transaction.ShippingMethod,
		PaymentMethod:  transaction.PaymentMethod,
		PaymentStatus:  "Belum Terbayar",
		TotalPayment:   totalPayment,
		CheckoutID:     uint(checkoutId),
	}

	output, _ := controller.transactionModel.CreateTransaction(newTransaction)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Create Transaction",
		"data":    output,
	})
}
