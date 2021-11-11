package stock

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
	db.Migrator().DropTable(&models.Stock{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Stock{})

	// preparate dummy data
	// dummy data for seller
	var newUser models.User
	newUser.Name = "Yudi"
	newUser.Address = "Tangerang"
	newUser.Email = "yudi@test.com"
	newUser.Password = "yudi77"
	newUser.Role = "seller"
	userModel := models.NewUserModel(db)
	_, err := userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCreateStockUpdateController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.POST("/api/stocks/:ingredientId", stockController.CreateStockUpdateController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller
	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "yudi@test.com",
		"password": "yudi77",
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

	// Get All Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"stock": 30,
	})

	req := httptest.NewRequest(echo.POST, "/api/stocks/8", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool         `json:"success"`
		Code    int          `json:"code"`
		Message string       `json:"message"`
		Data    models.Stock `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)
	fmt.Println(response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 30, response.Data.Stock)
}

func TestGetRestockDateController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/stocks/:range", stockController.GetRestockDateController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller
	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "yudi@test.com",
		"password": "yudi77",
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

	// Get All Controller
	req := httptest.NewRequest(echo.GET, "/api/stocks/daily", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool           `json:"success"`
		Code    int            `json:"code"`
		Message string         `json:"message"`
		Data    []models.Stock `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Get Restock Date", response.Message)
	// assert.NotEmpty(t, response.Data)
}
