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
	// create database connection
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	// cleaning data before testing
	db.Migrator().DropTable(&models.User{})
	db.AutoMigrate(&models.User{})

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

func TestRequestPayment(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	transactionModel := models.NewTransactionModel(db)
	midtransController := NewMidtransController(transactionModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/payments/request/:id", midtransController.RequestPayment, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Get Request Payment Controller

	req := httptest.NewRequest(echo.GET, "/api/payments/request/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
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
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	transactionModel := models.NewTransactionModel(db)
	midtransController := NewMidtransController(transactionModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/payments/status/:id", midtransController.RequestPayment, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Get Payment Status Controller

	req := httptest.NewRequest(echo.GET, "/api/payments/status/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
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
	assert.Equal(t, "Success Get Status Payment", response.Message)
	assert.NotEmpty(t, response.Data)
}
