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

	// -- Dummy Data
	var newUser models.User
	newUser.Name = "Yudi"
	newUser.Email = "yudi@mail.com"
	newUser.Password = "generate222"
	newUser.Address = "surabaya"
	newUser.Gender = "laki-laki"
	newUser.Role = "customer"

	userModel := models.NewUserModel(db)
	_, userModelErr := userModel.Register(newUser)
	if userModelErr != nil {
		fmt.Println(userModelErr)
	}

	newUser.Name = "Rudi"
	newUser.Email = "rudi@mail.com"
	newUser.Password = "generate4455"
	newUser.Address = "tangerang"
	newUser.Gender = "laki-laki"
	newUser.Role = "seller"

	_, userModelErr = userModel.Register(newUser)
	if userModelErr != nil {
		fmt.Println(userModelErr)
	}
}

func AuthValid(t *testing.T) string {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)

	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "rudi@mail.com",
		"password": "generate4455",
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

	token := loginResponse.Token
	return token
}

func AuthInvalid(t *testing.T) string {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)

	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "yudi@mail.com",
		"password": "generate222",
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

	token := loginResponse.Token
	return token
}

func TestCreateStockUpdateAuthInvalidController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.POST("/api/stocks/:ingredientId", stockController.CreateStockUpdateController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Get All Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"stock": 30,
	})

	req := httptest.NewRequest(echo.POST, "/api/stocks/1", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)
	fmt.Println(response)

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 401, res.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestCreateStockUpdateBadRequestAController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.POST("/api/stocks/:ingredientId", stockController.CreateStockUpdateController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Get All Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"stock": 30,
	})

	req := httptest.NewRequest(echo.POST, "/api/stocks/satu", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)
	fmt.Println(response)

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, res.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestCreateStockUpdateBadRequestBController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.POST("/api/stocks/:ingredientId", stockController.CreateStockUpdateController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Get All Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"stock": 0,
	})

	req := httptest.NewRequest(echo.POST, "/api/stocks/1", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)
	fmt.Println(response)

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, res.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestCreateStockUpdateController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.POST("/api/stocks/:ingredientId", stockController.CreateStockUpdateController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Get All Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"stock": 30,
	})

	req := httptest.NewRequest(echo.POST, "/api/stocks/1", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
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
	assert.Equal(t, "Success Create Stock Ingredient", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 30, response.Data.Stock)
}

func TestGetRestockDateAuthInvalidController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.GET("/api/restocks/:range", stockController.GetRestockDateController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Get All Controller
	req := httptest.NewRequest(echo.GET, "/api/restocks/daily", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestGetRestockDateParamBadRequestController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.GET("/api/restocks/:range", stockController.GetRestockDateController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Get All Controller
	req := httptest.NewRequest(echo.GET, "/api/restocks/yearly", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, res.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestGetRestockDateController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.GET("/api/restocks/:range", stockController.GetRestockDateController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Get All Controller
	req := httptest.NewRequest(echo.GET, "/api/restocks/daily", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
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
}

func TestGetAllRestockController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.GET("/api/restocks", stockController.GetRestockAllController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Get All Controller
	req := httptest.NewRequest(echo.GET, "/api/restocks", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
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
	assert.Equal(t, "Success Get All Restock", response.Message)
}

func TestInvalidGetAllRestockControllerNotSeller(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	stockModel := models.NewStockModel(db)
	stockController := NewStockController(stockModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.GET("/api/restocks", stockController.GetRestockAllController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Get All Controller
	req := httptest.NewRequest(echo.GET, "/api/restocks", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}
