package ingredient

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
	db := util.MysqlDatabaseConnection(config)

	// cleaning data before testing
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Ingredient{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Ingredient{})

	// preparate dummy data
	// dummy data for customer
	var newUser models.User
	newUser.Name = "Customer A"
	newUser.Email = "customer@test.com"
	newUser.Password = "passCust"
	newUser.Role = "customer"
	userModel := models.NewUserModel(db)
	_, err := userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}
	// dummy data for admin
	newUser.Name = "Admin A"
	newUser.Email = "admin@test.com"
	newUser.Password = "passAdmin"
	newUser.Role = "admin"
	userModel = models.NewUserModel(db)
	_, err = userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}

	// dummy data for ingredient
	var newIngredient models.Ingredient
	newIngredient.Name = "Daging ayam dada"
	newIngredient.Stock = 10
	newIngredient.Price = 30000

	// recipe dummy data with model
	ingredientModel := models.NewIngredientModel(db)
	_, err = ingredientModel.CreateIngredient(newIngredient)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCreateIngredientController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	ingredientModel := models.NewIngredientModel(db)
	stockModel := models.NewStockModel(db)
	ingredientController := NewIngredientController(ingredientModel, stockModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.POST("/api/ingredients", ingredientController.CreateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller Success
	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "admin@test.com",
		"password": "passAdmin",
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

	// Create Controller
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"name":  "Santan Kelapa",
		"price": 3000,
		"stock": 20,
	})
	req := httptest.NewRequest(echo.POST, "/api/ingredients", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool              `json:"success"`
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    models.Ingredient `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Create Ingredient Success", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, uint(2), response.Data.ID)
	assert.Equal(t, "Santan Kelapa", response.Data.Name)
	assert.Equal(t, 20, response.Data.Stock)
	assert.Equal(t, 3000, response.Data.Price)
}

func TestGetAllIngredientController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	ingredientModel := models.NewIngredientModel(db)
	stockModel := models.NewStockModel(db)
	ingredientController := NewIngredientController(ingredientModel, stockModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/ingredients", ingredientController.GetAllIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller Success
	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "customer@test.com",
		"password": "passCust",
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
	req := httptest.NewRequest(echo.GET, "/api/ingredients", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool                `json:"success"`
		Code    int                 `json:"code"`
		Message string              `json:"message"`
		Data    []models.Ingredient `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, uint(1), response.Data[0].ID)
	assert.Equal(t, "Daging ayam dada", response.Data[0].Name)
	assert.Equal(t, 10, response.Data[0].Stock)
	assert.Equal(t, 30000, response.Data[0].Price)
	assert.Equal(t, uint(2), response.Data[1].ID)
	assert.Equal(t, "Santan Kelapa", response.Data[1].Name)
	assert.Equal(t, 20, response.Data[1].Stock)
	assert.Equal(t, 3000, response.Data[1].Price)
}

func TestGetRecipeByIdController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	ingredientModel := models.NewIngredientModel(db)
	stockModel := models.NewStockModel(db)
	ingredientController := NewIngredientController(ingredientModel, stockModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/ingredients/:ingredientId", ingredientController.GetIngredientByIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller
	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "customer@test.com",
		"password": "passCust",
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
	req := httptest.NewRequest(echo.GET, "/api/ingredients/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool              `json:"success"`
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    models.Ingredient `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, uint(1), response.Data.ID)
	assert.Equal(t, "Daging ayam dada", response.Data.Name)
	assert.Equal(t, 10, response.Data.Stock)
	assert.Equal(t, 30000, response.Data.Price)
}

func TestUpdateIngredientController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	ingredientModel := models.NewIngredientModel(db)
	stockModel := models.NewStockModel(db)
	ingredientController := NewIngredientController(ingredientModel, stockModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.PUT("/api/ingredients/:ingredientId", ingredientController.UpdateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller Success
	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "admin@test.com",
		"password": "passAdmin",
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

	// Create Controller
	reqBodyPut, _ := json.Marshal(map[string]interface{}{
		"name":  "Bawang Merah",
		"price": 1000,
		"stock": 15,
	})
	req := httptest.NewRequest(echo.PUT, "/api/ingredients/1", bytes.NewBuffer(reqBodyPut))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool              `json:"success"`
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    models.Ingredient `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Update Ingredient Success", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, uint(1), response.Data.ID)
	assert.Equal(t, "Bawang Merah", response.Data.Name)
	assert.Equal(t, 15, response.Data.Stock)
	assert.Equal(t, 1000, response.Data.Price)
}

// func TestUpdateIngredientStockController(t *testing.T) {
// 	// create database connection and create controller
// 	config := config.GetConfig()
// 	db := util.MysqlDatabaseConnection(config)
// 	userModel := models.NewUserModel(db)
// 	userController := auth.NewAuthController(userModel)

// 	ingredientModel := models.NewIngredientModel(db)
// 	stockModel := models.NewStockModel(db)
// 	ingredientController := NewIngredientController(ingredientModel, stockModel)

// 	// Setting Route
// 	e := echo.New()
// 	e.POST("/api/login", userController.LoginUserController)
// 	e.PUT("/api/ingredients/stock/:ingredientId", ingredientController.UpdateIngredientStockController, middleware.JWT([]byte(constants.SECRET_JWT)))

// 	// Login Controller Success
// 	reqBodyLogin, _ := json.Marshal(map[string]string{
// 		"email":    "admin@test.com",
// 		"password": "passAdmin",
// 	})

// 	loginReq := httptest.NewRequest(echo.POST, "/api/login", bytes.NewBuffer(reqBodyLogin))
// 	loginReq.Header.Set("Content-Type", "application/json")
// 	loginRes := httptest.NewRecorder()
// 	e.ServeHTTP(loginRes, loginReq)

// 	type LoginResponse struct {
// 		Success bool   `json:"success"`
// 		Code    int    `json:"code"`
// 		Message string `json:"message"`
// 		Token   string `json:"token"`
// 	}

// 	var loginResponse LoginResponse
// 	json.Unmarshal(loginRes.Body.Bytes(), &loginResponse)

// 	assert.Equal(t, true, loginResponse.Success)
// 	assert.Equal(t, 200, loginResponse.Code)
// 	assert.Equal(t, "Success Login", loginResponse.Message)
// 	assert.NotEqual(t, "", loginResponse.Token)

// 	// Create Controller
// 	reqBodyPut, _ := json.Marshal(map[string]interface{}{
// 		"stock": 50,
// 	})
// 	req := httptest.NewRequest(echo.PUT, "/api/ingredients/stock/2", bytes.NewBuffer(reqBodyPut))
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
// 	req.Header.Set("Content-Type", "application/json")
// 	res := httptest.NewRecorder()
// 	e.ServeHTTP(res, req)

// 	type Response struct {
// 		Success bool              `json:"success"`
// 		Code    int               `json:"code"`
// 		Message string            `json:"message"`
// 		Data    models.Ingredient `json:"data"`
// 	}

// 	var response Response
// 	json.Unmarshal(res.Body.Bytes(), &response)

// 	assert.Equal(t, true, response.Success)
// 	assert.Equal(t, 200, response.Code)
// 	assert.Equal(t, "Update Ingredient Stock Success", response.Message)
// 	assert.NotEmpty(t, response.Data)
// 	assert.Equal(t, uint(2), response.Data.ID)
// 	assert.Equal(t, "Santan Kelapa", response.Data.Name)
// 	assert.Equal(t, 50, response.Data.Stock)
// 	assert.Equal(t, 3000, response.Data.Price)
// }

func TestDeleteIngredientController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	ingredientModel := models.NewIngredientModel(db)
	stockModel := models.NewStockModel(db)
	ingredientController := NewIngredientController(ingredientModel, stockModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.DELETE("/api/ingredients/:ingredientId", ingredientController.DeleteIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller
	reqBodyLogin, _ := json.Marshal(map[string]string{
		"email":    "admin@test.com",
		"password": "passAdmin",
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

	// Delete Controller
	req := httptest.NewRequest(echo.DELETE, "/api/ingredients/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
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

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Delete Ingredient", response.Message)
}
