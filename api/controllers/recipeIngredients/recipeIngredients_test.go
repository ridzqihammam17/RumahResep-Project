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
	db := util.MysqlDatabaseConnection(config)

	// cleaning data before testing
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.RecipeIngredients{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.RecipeIngredients{})

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
	var newRecipeIngredient models.RecipeIngredients
	newRecipeIngredient.RecipeId = 1
	newRecipeIngredient.IngredientId = 1
	newRecipeIngredient.QtyIngredient = 10
	// newRecipe.Ingredients = "1"

	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	_, err = recipeIngredientModel.AddIngredientsRecipe(newRecipeIngredient)
	if err != nil {
		fmt.Println(err)
	}
}

func TestAddIngredientsRecipeController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	recipeIngredientsModel := models.NewRecipeIngredientsModel(db)
	recipesModel := models.NewRecipeModel(db)
	ingredientModel := models.NewIngredientModel(db)
	recipeIngredientsController := NewRecipeIngredientsController(recipeIngredientsModel, recipesModel, ingredientModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.POST("/api/ingredients/recipe", recipeIngredientsController.AddIngredientsRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Add Recipe Ingredients Controller
	reqBodyPost, _ := json.Marshal(map[string]int{
		"ingredient_id":  2,
		"recipe_id":      1,
		"qty_ingredient": 15,
	})
	req := httptest.NewRequest(echo.POST, "/api/ingredients/recipe", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
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

func TestGetIngredientsByRecipeIdController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	recipeIngredientsModel := models.NewRecipeIngredientsModel(db)
	recipesModel := models.NewRecipeModel(db)
	ingredientsModel := models.NewIngredientModel(db)
	recipeIngredientsController := NewRecipeIngredientsController(recipeIngredientsModel, recipesModel, ingredientsModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/ingredients/recipe/:recipeId", recipeIngredientsController.GetIngredientsByRecipeIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Add Recipe Ingredients Controller
	req := httptest.NewRequest(echo.GET, "/api/ingredients/recipe/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
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
