package transactions

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"

	"github.com/labstack/echo/v4"
)

// ------------------------------------------------------------------
// Start Request
// ------------------------------------------------------------------

type TransactionRequest struct {
	UserID         uint
	CustomerName   string `json:"customer_name" form:"customer_name"`
	Address        string `json:"address" form:"address"`
	ShippingMethod string `json:"shipping_method" form:"shipping_method"`
	PaymentMethod  string `json:"payment_method" form:"payment_method"`
	PaymentStatus  string `json:"payment_status" form:"payment_status"`
	TotalPayment   int    `json:"total_payment" form:"total_payment"`
	CheckoutID     uint
}

// ------------------------------------------------------------------
// End Request
// ------------------------------------------------------------------

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

func (controller *TransactionController) GetAllTransactionAdmin(c echo.Context) error {
	//check role admin or not
	_, role := middlewares.ExtractTokenUser(c)
	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	data, err := controller.transactionModel.GetAllTransactionAdmin()
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
		"message": "Success Get All Transaction",
		"data":    data,
	})
}

func (controller *TransactionController) GetAllTransactionCustomer(c echo.Context) error {
	//check role customer or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	data, err := controller.transactionModel.GetAllTransaction(int(userId))
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
		"message": "Success Get All Transaction",
		"data":    data,
	})
}

func (controller *TransactionController) CreateTransaction(c echo.Context) error {
	//check role customer or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	cartId, err := controller.cartModel.GetCartIdByUserId(int(userId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	checkoutId, err := controller.transactionModel.GetCheckoutId(cartId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	user, err := controller.userModel.GetUserData(int(userId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	totalPayment, err := controller.transactionModel.CountTotalPayment(cartId, checkoutId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	var transactionRequest TransactionRequest

	c.Bind(&transactionRequest)

	transaction := models.Transaction{
		UserID:         userId,
		CustomerName:   user.Name,
		Address:        user.Address,
		ShippingMethod: transactionRequest.ShippingMethod,
		PaymentMethod:  "Belum Dipilih",
		PaymentStatus:  "Belum Terbayar",
		TotalPayment:   totalPayment,
		CheckoutID:     uint(checkoutId),
	}
	if transaction.ShippingMethod == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	data, _ := controller.transactionModel.CreateTransaction(transaction)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Create Transaction",
		"data":    data,
	})
}
