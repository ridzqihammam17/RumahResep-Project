package recipescategories

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
	db.Migrator().DropTable(&models.RecipeCategories{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.RecipeCategories{})

	// preparate dummy data

	// dummy data for admin
	var newUser models.User
	newUser.Name = "Admin A"
	newUser.Email = "admin@test.com"
	newUser.Password = "passAdmin"
	newUser.Role = "admin"
	userModel := models.NewUserModel(db)
	_, err := userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}
	// dummy data for recipe
	var newRecipeCategory models.RecipeCategories
	newRecipeCategory.RecipeId = 1
	newRecipeCategory.CategoryId = 1

	recipeCategoryModel := models.NewRecipesCategoriesModel(db)
	_, err = recipeCategoryModel.AddRecipeCategories(newRecipeCategory)
	if err != nil {
		fmt.Println(err)
	}
}

func TestAddRecipeCategoriesController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipeCategoriesController := NewRecipesCategoriesController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.POST("/api/recipe/categories", recipeCategoriesController.AddRecipeCategoriesController, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Add Recipe Categories Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"recipes_id":    2,
		"categories_id": 1,
	})
	req := httptest.NewRequest(echo.POST, "/api/recipe/categories", bytes.NewBuffer(reqBodyPost))
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
	assert.Equal(t, "Success Add Recipe Category", response.Message)
	// assert.NotEmpty(t, response.Data)
	// assert.Equal(t, 2, len(response.Data))
}

func TestGetRecipeByCategoryId(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipeCategoriesController := NewRecipesCategoriesController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/recipe/categories/:categoryId", recipeCategoriesController.GetRecipeByCategoryIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Add Recipe Categories Controller
	req := httptest.NewRequest(echo.GET, "/api/recipe/categories/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool                      `json:"success"`
		Code    int                       `json:"code"`
		Message string                    `json:"message"`
		Data    []models.RecipeCategories `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Get Recipe By Category ID", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, 1, response.Data[0].RecipeId)
	assert.Equal(t, 2, response.Data[1].RecipeId)
}
