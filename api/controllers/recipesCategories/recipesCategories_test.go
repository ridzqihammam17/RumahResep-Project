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
	db := util.MysqlDatabaseConnTest(config)

	// cleaning data before testing
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
	db.AutoMigrate(&models.Recipe{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.RecipeCategories{})

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

	newUser.Name = "Budi"
	newUser.Email = "budi@mail.com"
	newUser.Password = "generate999"
	newUser.Address = "jakarta"
	newUser.Gender = "laki-laki"
	newUser.Role = "admin"

	_, userModelErr = userModel.Register(newUser)
	if userModelErr != nil {
		fmt.Println(userModelErr)
	}

	// dummy data for recipe
	var newRecipe models.Recipe
	newRecipe.Name = "Rujak Cingur"

	recipeModel := models.NewRecipeModel(db)
	_, recipeModelErr := recipeModel.CreateRecipe(newRecipe)
	if recipeModelErr != nil {
		fmt.Println(recipeModelErr)
	}

	var newCategory models.Category
	newCategory.Name = "Makanan"

	categoryModel := models.NewCategoryModel(db)
	_, categoryModelErr := categoryModel.Insert(newCategory)
	if categoryModelErr != nil {
		fmt.Println(categoryModelErr)
	}

	var newRecipeCategory models.RecipeCategories
	newRecipeCategory.RecipeId = 1
	newRecipeCategory.CategoryId = 1

	recipeCategoryModel := models.NewRecipesCategoriesModel(db)
	_, err := recipeCategoryModel.AddRecipeCategories(newRecipeCategory)
	if err != nil {
		fmt.Println(err)
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
		"email":    "budi@mail.com",
		"password": "generate999",
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

func TestAddRecipeCategoriesAuthInvalidController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipeCategoriesController := NewRecipesCategoriesController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.POST("/api/recipe/categories", recipeCategoriesController.AddRecipeCategoriesController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Categories Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"recipes_id":    2,
		"categories_id": 1,
	})
	req := httptest.NewRequest(echo.POST, "/api/recipe/categories", bytes.NewBuffer(reqBodyPost))
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
	assert.Equal(t, "Unauthorized", response.Message)
}

func TestAddRecipeCategoriesBadRequestController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipeCategoriesController := NewRecipesCategoriesController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.POST("/api/recipe/categories", recipeCategoriesController.AddRecipeCategoriesController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Categories Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"recipes_id":    0,
		"categories_id": 1,
	})
	req := httptest.NewRequest(echo.POST, "/api/recipe/categories", bytes.NewBuffer(reqBodyPost))
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
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestAddRecipeCategoriesController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipeCategoriesController := NewRecipesCategoriesController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.POST("/api/recipe/categories", recipeCategoriesController.AddRecipeCategoriesController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Categories Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"recipes_id":    2,
		"categories_id": 1,
	})
	req := httptest.NewRequest(echo.POST, "/api/recipe/categories", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool                    `json:"success"`
		Code    int                     `json:"code"`
		Message string                  `json:"message"`
		Data    models.RecipeCategories `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Add Recipe Category", response.Message)
	assert.NotEmpty(t, response.Data)
}

func TestGetRecipeByCategoryIdBadRequestController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipeCategoriesController := NewRecipesCategoriesController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.GET("/api/recipe/categories/:categoryId", recipeCategoriesController.GetRecipeByCategoryIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Categories Controller
	req := httptest.NewRequest(echo.GET, "/api/recipe/categories/satu", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestGetRecipeByCategoryIdController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipeCategoriesController := NewRecipesCategoriesController(recipeCategoriesModel, recipesModel, categoriesModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.GET("/api/recipe/categories/:categoryId", recipeCategoriesController.GetRecipeByCategoryIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Categories Controller
	req := httptest.NewRequest(echo.GET, "/api/recipe/categories/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
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
