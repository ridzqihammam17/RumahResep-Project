package midtrans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"rumah_resep/api/controllers/auth"
	"rumah_resep/config"
	"rumah_resep/constants"
	"rumah_resep/models"
	"rumah_resep/util"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	// -- Create Connection
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	// -- Clean DB Data
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.CartDetails{})
	db.Migrator().DropTable(&models.Cart{})
	db.Migrator().DropTable(&models.Checkout{})
	db.Migrator().DropTable(&models.Stock{})
	db.Migrator().DropTable(&models.RecipeIngredients{})
	db.Migrator().DropTable(&models.RecipeCategories{})
	db.Migrator().DropTable(&models.Ingredient{})
	db.Migrator().DropTable(&models.Category{})
	db.Migrator().DropTable(&models.Recipe{})
	db.Migrator().DropTable(&models.User{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Checkout{})
	db.AutoMigrate(&models.Transaction{})

	// -- Dummy Data with Model
	// ------ Start User ------
	var newUser models.User
	newUser.Name = "Rudi"
	newUser.Address = "Bekasi"
	newUser.Email = "rudi@test.com"
	newUser.Password = "rudi99"
	newUser.Role = "customer"
	userModel := models.NewUserModel(db)
	_, userModelErr := userModel.Register(newUser)
	if userModelErr != nil {
		fmt.Println(userModelErr)
	}
	// ------ End User ------

	// ------ Start Checkout ------
	var newCheckout models.Checkout
	checkoutModel := models.NewCheckoutModel(db)
	_, checkoutModelErr := checkoutModel.CreateCheckout(newCheckout)
	if checkoutModelErr != nil {
		fmt.Println(checkoutModelErr)
	}
	// ------ End Checkout ------

	// ------ Start Transaction ------
	var newTransaction models.Transaction
	newTransaction.ID = 120
	newTransaction.UserID = 1
	newTransaction.CustomerName = "Rudi"
	newTransaction.Address = "Depok"
	newTransaction.ShippingMethod = "Kirim"
	newTransaction.PaymentMethod = "Belum Dipilih"
	newTransaction.PaymentStatus = "Belum Dibayar"
	newTransaction.TotalPayment = 51000
	newTransaction.CheckoutID = 1
	transactionModel := models.NewTransactionModel(db)
	_, transactionModelErr := transactionModel.CreateTransaction(newTransaction)
	if transactionModelErr != nil {
		fmt.Println(transactionModelErr)
	}
	// ------ End Transaction ------
}

func TestRequestPayment(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	midtransController := NewMidtransController(transactionModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/payments/request/:id", midtransController.RequestPayment, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]string{
		"email":    "rudi@test.com",
		"password": "rudi99",
	})

	// -- Setting Controller
	reqLogin := httptest.NewRequest(echo.POST, "/api/login", bytes.NewBuffer(reqBodyPostLogin))
	reqLogin.Header.Set("Content-Type", "application/json")
	resLogin := httptest.NewRecorder()
	e.ServeHTTP(resLogin, reqLogin)

	// -- Declare Response and Convert to JSON
	type ResponseLogin struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}

	var responseLogin ResponseLogin

	json.Unmarshal(resLogin.Body.Bytes(), &responseLogin)

	assert.Equal(t, true, responseLogin.Success)
	assert.Equal(t, 200, resLogin.Code)
	assert.Equal(t, "Success Login", responseLogin.Message)
	assert.NotEqual(t, "", responseLogin.Token)
	// ------ End Login ------

	// -- Setting Controller
	req := httptest.NewRequest(echo.GET, "/api/payments/request/120", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Request Payment", response.Message)
	assert.NotEmpty(t, response.Data)
}

func TestStatusPayment(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	midtransController := NewMidtransController(transactionModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/payments/status/:id", midtransController.StatusPayment, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]string{
		"email":    "rudi@test.com",
		"password": "rudi99",
	})

	// -- Setting Controller
	reqLogin := httptest.NewRequest(echo.POST, "/api/login", bytes.NewBuffer(reqBodyPostLogin))
	reqLogin.Header.Set("Content-Type", "application/json")
	resLogin := httptest.NewRecorder()
	e.ServeHTTP(resLogin, reqLogin)

	// -- Declare Response and Convert to JSON
	type ResponseLogin struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}

	var responseLogin ResponseLogin

	json.Unmarshal(resLogin.Body.Bytes(), &responseLogin)

	assert.Equal(t, true, responseLogin.Success)
	assert.Equal(t, 200, resLogin.Code)
	assert.Equal(t, "Success Login", responseLogin.Message)
	assert.NotEqual(t, "", responseLogin.Token)
	// ------ End Login ------

	// -- Setting Controller
	req := httptest.NewRequest(echo.GET, "/api/payments/status/111", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Get Status Payment", response.Message)
	assert.NotEmpty(t, response.Data)
}
