package recipes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	auth "rumah_resep/api/controllers/auth"
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
	db.Migrator().DropTable(&models.Recipe{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Recipe{})

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
	// dummy data for recipe
	var newRecipe models.Recipe
	newRecipe.Name = "Recipe A"
	// newRecipe.Categories = "1"

	recipeModel := models.NewRecipeModel(db)
	_, err = recipeModel.CreateRecipe(newRecipe)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetAllRecipeController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/recipes", recipesController.GetAllRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.GET, "/api/recipes", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool            `json:"success"`
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    []models.Recipe `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, "Recipe A", response.Data[0].Name)
}

func TestGetRecipeByIdController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/recipes/:recipeId", recipesController.GetRecipeByIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.GET, "/api/recipes/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool          `json:"success"`
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    models.Recipe `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Get Recipe By Id", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, "Recipe A", response.Data.Name)
}

func TestCreateRecipeController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.POST("/api/recipes", recipesController.CreateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Get All Controller
	reqBodyPost, _ := json.Marshal(map[string]string{
		"name": "Recipe B",
	})
	req := httptest.NewRequest(echo.POST, "/api/recipes", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool          `json:"success"`
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    models.Recipe `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)
	fmt.Println(response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, "Recipe B", response.Data.Name)
}

func TestUpdateRecipeController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.PUT("/api/recipes/:recipeId", recipesController.UpdateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Get All Controller
	reqBodyPut, _ := json.Marshal(map[string]string{
		"name": "Recipe B Updated",
	})
	req := httptest.NewRequest(echo.PUT, "/api/recipes/2", bytes.NewBuffer(reqBodyPut))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool          `json:"success"`
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    models.Recipe `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)
	fmt.Println(response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Update Recipe", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, uint(2), response.Data.ID)
	assert.Equal(t, "Recipe B Updated", response.Data.Name)
}

func TestDeleteRecipeController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.DELETE("/api/recipes/:recipeId", recipesController.DeleteRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Delete Recipe
	req := httptest.NewRequest(echo.DELETE, "/api/recipes/1", nil)
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
	fmt.Println(response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Delete Recipe", response.Message)
}
