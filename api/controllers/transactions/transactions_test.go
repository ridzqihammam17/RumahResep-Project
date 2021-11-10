package transactions

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
	// create database connection
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	// cleaning data before testing
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Transaction{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Transaction{})

	// preparate dummy data
	// dummy data for customer
	var newUser models.User
	newUser.Name = "Rudi"
	newUser.Address = "Bekasi"
	newUser.Email = "rudi@test.com"
	newUser.Password = "rudi99"
	newUser.Role = "customer"
	userModel := models.NewUserModel(db)
	_, err := userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCreateTransaction(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	transactionModel := models.NewTransactionModel(db)
	cartModel := models.NewCartModel(db)
	transactionController := NewTransactionController(transactionModel, cartModel, userModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.POST("/api/transactions", transactionController.CreateTransaction, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller
	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "rudi@test.com",
		"password": "rudi99",
	})

	loginReq := httptest.NewRequest(echo.POST, "/api/login", bytes.NewBuffer(reqBodyLogin))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	e.ServeHTTP(loginRes, loginReq)

	type LoginResponse struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}

	var loginResponse LoginResponse
	json.Unmarshal(loginRes.Body.Bytes(), &loginResponse)

	assert.Equal(t, true, loginResponse.Success)
	assert.Equal(t, 200, loginResponse.Code)
	assert.Equal(t, "Success Login", loginResponse.Message)
	assert.NotEqual(t, "", loginResponse.Token)

	// Create Transaction Controller
	reqBodyPost, _ := json.Marshal(map[string]string{
		"shipping_method": "delivery",
	})

	req := httptest.NewRequest(echo.POST, "/api/transactions", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool               `json:"success"`
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    models.Transaction `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)
	fmt.Println(response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Create Transaction", response.Message)
	assert.NotEmpty(t, response.Data)
}
