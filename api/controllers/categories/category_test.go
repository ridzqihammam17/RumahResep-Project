package categories

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
	db := util.MysqlDatabaseConnection(config)

	// -- Clean DB Data
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Category{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Category{})

	// -- Dummy Data
	var newUser models.User
	newUser.Name = "Budi"
	newUser.Email = "budi@mail.com"
	newUser.Password = "generate111"
	newUser.Address = "jl. barat laut no 1"
	newUser.Gender = "laki"
	newUser.Role = "admin"

	var newCategory models.Category
	newCategory.Name = "Buah"

	// -- Dummy Data with Model
	userModel := models.NewUserModel(db)
	_, userModelErr := userModel.Register(newUser)
	if userModelErr != nil {
		fmt.Println(userModelErr)
	}

	categoryModel := models.NewCategoryModel(db)
	_, categoryModelErr := categoryModel.Insert(newCategory)
	if categoryModelErr != nil {
		fmt.Println(categoryModelErr)
	}
}

func TestValidGetAllCategoryController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	categoryModel := models.NewCategoryModel(db)
	categoryController := NewCategoryController(categoryModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/categories", categoryController.GetAllCategoryController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]string{
		"email":    "budi@mail.com",
		"password": "generate111",
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
	req := httptest.NewRequest(echo.GET, "/api/categories", nil)
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool              `json:"success"`
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    []models.Category `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Get All Category", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, "Buah", response.Data[0].Name)
}

func TestValidGetCategoryController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	categoryModel := models.NewCategoryModel(db)
	categoryController := NewCategoryController(categoryModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/categories/:id", categoryController.GetCategoryController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]string{
		"email":    "budi@mail.com",
		"password": "generate111",
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
	req := httptest.NewRequest(echo.GET, "/api/categories/1", nil)
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool            `json:"success"`
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    models.Category `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Get Category", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, "Buah", response.Data.Name)
}

func TestValidInsertCategoryController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	categoryModel := models.NewCategoryModel(db)
	categoryController := NewCategoryController(categoryModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/categories", categoryController.InsertCategoryController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]string{
		"email":    "budi@mail.com",
		"password": "generate111",
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

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]string{
		"name": "Software",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/categories", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Insert Category", response.Message)
}

func TestValidEditCategoryController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	categoryModel := models.NewCategoryModel(db)
	categoryController := NewCategoryController(categoryModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/categories/:id", categoryController.EditCategoryController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]string{
		"email":    "budi@mail.com",
		"password": "generate111",
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

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]string{
		"name": "Software Edit",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/categories/1", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Edit Category", response.Message)
}

func TestValidDeleteCategoryController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	categoryModel := models.NewCategoryModel(db)
	categoryController := NewCategoryController(categoryModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.DELETE("/api/categories/:id", categoryController.DeleteCategoryController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]string{
		"email":    "budi@mail.com",
		"password": "generate111",
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
	req := httptest.NewRequest(echo.DELETE, "/api/categories/1", nil)
	req.Header.Set("Authorization", fmt.Sprint("Bearer ", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Delete Category", response.Message)
}
