package cartdetails

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
	db.Migrator().DropTable(&models.CartDetails{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Category{})

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
}

func TestAddRecipeToCartController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	cartDetailsModel := models.NewCartDetailsModel(db)
	recipeModel := models.NewRecipeModel(db)
	ingredientModel := models.NewIngredientModel(db)
	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	cartModel := models.NewCartModel(db)
	cartDetailsController := NewCartDetailsController(cartDetailsModel, recipeModel, ingredientModel, recipeIngredientModel, cartModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.POST("/api/cartdetails", cartDetailsController.AddRecipeToCartController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller
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

	// Add To Cart Controller
	reqBodyPost, _ := json.Marshal((map[string]int{
		"recipe_id": 1,
		"quantity":  1,
	}))

	req := httptest.NewRequest(echo.POST, "/api/cartdetails", bytes.NewReader(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool               `json:"success"`
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    models.CartDetails `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 1, response.Data.CartID)
	assert.Equal(t, 1, response.Data.RecipeID)
	assert.Equal(t, 1, response.Data.Quantity)
	assert.Equal(t, 55000, response.Data.Price)
}

func TestGetAllRecipeByCartIdController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	cartDetailsModel := models.NewCartDetailsModel(db)
	recipeModel := models.NewRecipeModel(db)
	ingredientModel := models.NewIngredientModel(db)
	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	cartModel := models.NewCartModel(db)
	cartDetailsController := NewCartDetailsController(cartDetailsModel, recipeModel, ingredientModel, recipeIngredientModel, cartModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.GET("/api/cartdetails", cartDetailsController.GetAllRecipeByCartIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller
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

	// Get All Recipe By Cart Id Controller
	req := httptest.NewRequest(echo.GET, "/api/cartdetails", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool                 `json:"success"`
		Code    int                  `json:"code"`
		Message string               `json:"message"`
		Data    []models.CartDetails `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 1, response.Data[0].CartID)
	assert.Equal(t, 1, response.Data[0].RecipeID)
	assert.Equal(t, 1, response.Data[0].Quantity)
	assert.Equal(t, 55000, response.Data[0].Price)
}

func TestUpdateRecipePortionController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	cartDetailsModel := models.NewCartDetailsModel(db)
	recipeModel := models.NewRecipeModel(db)
	ingredientModel := models.NewIngredientModel(db)
	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	cartModel := models.NewCartModel(db)
	cartDetailsController := NewCartDetailsController(cartDetailsModel, recipeModel, ingredientModel, recipeIngredientModel, cartModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.PUT("/api/cartdetails/:recipeId", cartDetailsController.UpdateRecipePortionController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller
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

	// Add To Cart Controller
	reqBodyPost, _ := json.Marshal((map[string]int{
		"quantity": 2,
	}))

	req := httptest.NewRequest(echo.PUT, "/api/cartdetails/1", bytes.NewReader(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool               `json:"success"`
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    models.CartDetails `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 1, response.Data.CartID)
	assert.Equal(t, 1, response.Data.RecipeID)
	assert.Equal(t, 2, response.Data.Quantity)
}

func TestDeleteRecipeFromCartController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)

	cartDetailsModel := models.NewCartDetailsModel(db)
	recipeModel := models.NewRecipeModel(db)
	ingredientModel := models.NewIngredientModel(db)
	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	cartModel := models.NewCartModel(db)
	cartDetailsController := NewCartDetailsController(cartDetailsModel, recipeModel, ingredientModel, recipeIngredientModel, cartModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.DELETE("/api/cartdetails/:recipeId", cartDetailsController.DeleteRecipeFromCartController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// Login Controller
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

	// Get All Recipe By Cart Id Controller
	req := httptest.NewRequest(echo.DELETE, "/api/cartdetails/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool               `json:"success"`
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    models.CartDetails `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Delete Recipe From Cart", response.Message)
}
