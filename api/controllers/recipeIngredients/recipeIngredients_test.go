package recipeingredients

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
	db.Migrator().DropTable(&models.RecipeIngredients{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.RecipeIngredients{})

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
	var newRecipeIngredient models.RecipeIngredients
	newRecipeIngredient.RecipeId = 1
	newRecipeIngredient.IngredientId = 1
	newRecipeIngredient.QtyIngredient = 10
	// newRecipe.Ingredients = "1"

	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	_, err := recipeIngredientModel.AddIngredientsRecipe(newRecipeIngredient)
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

func TestAddIngredientsRecipeAuthInvalidController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeIngredientsModel := models.NewRecipeIngredientsModel(db)
	recipesModel := models.NewRecipeModel(db)
	ingredientModel := models.NewIngredientModel(db)
	recipeIngredientsController := NewRecipeIngredientsController(recipeIngredientsModel, recipesModel, ingredientModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.POST("/api/ingredients/recipe", recipeIngredientsController.AddIngredientsRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Ingredients Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"ingredient_id":  2,
		"recipe_id":      1,
		"qty_ingredient": 15,
	})
	req := httptest.NewRequest(echo.POST, "/api/ingredients/recipe", bytes.NewBuffer(reqBodyPost))
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

func TestAddIngredientsRecipeBadRequestController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeIngredientsModel := models.NewRecipeIngredientsModel(db)
	recipesModel := models.NewRecipeModel(db)
	ingredientModel := models.NewIngredientModel(db)
	recipeIngredientsController := NewRecipeIngredientsController(recipeIngredientsModel, recipesModel, ingredientModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.POST("/api/ingredients/recipe", recipeIngredientsController.AddIngredientsRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Ingredients Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"recipe_id":      1,
		"ingredient_id":  0,
		"qty_ingredient": 15,
	})
	req := httptest.NewRequest(echo.POST, "/api/ingredients/recipe", bytes.NewBuffer(reqBodyPost))
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

func TestAddIngredientsRecipeController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeIngredientsModel := models.NewRecipeIngredientsModel(db)
	recipesModel := models.NewRecipeModel(db)
	ingredientModel := models.NewIngredientModel(db)
	recipeIngredientsController := NewRecipeIngredientsController(recipeIngredientsModel, recipesModel, ingredientModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.POST("/api/ingredients/recipe", recipeIngredientsController.AddIngredientsRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Ingredients Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"ingredient_id":  2,
		"recipe_id":      1,
		"qty_ingredient": 15,
	})
	req := httptest.NewRequest(echo.POST, "/api/ingredients/recipe", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool                     `json:"success"`
		Code    int                      `json:"code"`
		Message string                   `json:"message"`
		Data    models.RecipeIngredients `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Add Recipe Ingredient", response.Message)
	assert.NotEmpty(t, response.Data)
}

func TestGetIngredientsByRecipeIdBadRequestController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeIngredientsModel := models.NewRecipeIngredientsModel(db)
	recipesModel := models.NewRecipeModel(db)
	ingredientsModel := models.NewIngredientModel(db)
	recipeIngredientsController := NewRecipeIngredientsController(recipeIngredientsModel, recipesModel, ingredientsModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.GET("/api/ingredients/recipe/:recipeId", recipeIngredientsController.GetIngredientsByRecipeIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Ingredients Controller
	req := httptest.NewRequest(echo.GET, "/api/ingredients/recipe/satu", nil)
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

func TestGetIngredientsByRecipeIdController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	recipeIngredientsModel := models.NewRecipeIngredientsModel(db)
	recipesModel := models.NewRecipeModel(db)
	ingredientsModel := models.NewIngredientModel(db)
	recipeIngredientsController := NewRecipeIngredientsController(recipeIngredientsModel, recipesModel, ingredientsModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.GET("/api/ingredients/recipe/:recipeId", recipeIngredientsController.GetIngredientsByRecipeIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Add Recipe Ingredients Controller
	req := httptest.NewRequest(echo.GET, "/api/ingredients/recipe/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool                       `json:"success"`
		Code    int                        `json:"code"`
		Message string                     `json:"message"`
		Data    []models.RecipeIngredients `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Get Ingredients By Recipe ID", response.Message)
	assert.NotEmpty(t, response.Data)
}
