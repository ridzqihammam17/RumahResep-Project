package checkouts

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
	db := util.MysqlDatabaseConnTest(config)

	// -- Clean DB Data
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
	db.AutoMigrate(&models.Ingredient{})
	db.AutoMigrate(&models.RecipeIngredients{})
	db.AutoMigrate(&models.Stock{})
	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.CartDetails{})
	db.AutoMigrate(&models.Checkout{})

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

	var newRecipe models.Recipe
	newRecipe.Name = "Soto Ayam"

	recipeModel := models.NewRecipeModel(db)
	_, recipeModelErr := recipeModel.CreateRecipe(newRecipe)
	if recipeModelErr != nil {
		fmt.Println(recipeModelErr)
	}

	var newIngredient models.Ingredient
	newIngredient.Name = "Bawang Merah"
	newIngredient.Price = 500

	ingredientModel := models.NewIngredientModel(db)
	_, ingredientModelErr := ingredientModel.CreateIngredient(newIngredient)
	if ingredientModelErr != nil {
		fmt.Println(ingredientModelErr)
	}

	newIngredient.Name = "Daging Ayam"
	newIngredient.Price = 4000

	_, ingredientModelErr = ingredientModel.CreateIngredient(newIngredient)
	if ingredientModelErr != nil {
		fmt.Println(ingredientModelErr)
	}

	var newRecipeIngredient models.RecipeIngredients
	newRecipeIngredient.RecipeId = 1
	newRecipeIngredient.IngredientId = 1
	newRecipeIngredient.QtyIngredient = 5

	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	_, recipeIngredientErr := recipeIngredientModel.AddIngredientsRecipe(newRecipeIngredient)
	if recipeIngredientErr != nil {
		fmt.Println(recipeIngredientErr)
	}

	newRecipeIngredient.RecipeId = 1
	newRecipeIngredient.IngredientId = 2
	newRecipeIngredient.QtyIngredient = 1

	_, recipeIngredientErr = recipeIngredientModel.AddIngredientsRecipe(newRecipeIngredient)
	if recipeIngredientErr != nil {
		fmt.Println(recipeIngredientErr)
	}

	var newCart models.Cart
	newCart.UserID = 1
	cartModel := models.NewCartModel(db)
	_, cartErr := cartModel.CreateCart(newCart, int(newCart.UserID))
	if cartErr != nil {
		fmt.Println(cartErr)
	}

	var newCartDetail models.CartDetails
	newCartDetail.CartID = 1
	newCartDetail.RecipeID = 1
	newCartDetail.Quantity = 2
	cartDetailModel := models.NewCartDetailsModel(db)
	_, cartDetailErr := cartDetailModel.AddRecipeToCart(newCartDetail)
	if cartDetailErr != nil {
		fmt.Println(cartDetailErr)
	}

	var newStock models.Stock
	newStock.IngredientId = 1
	newStock.Stock = 10
	stockModel := models.NewStockModel(db)
	_, stockErr := stockModel.CreateStockUpdate(newStock, int(newStock.IngredientId))
	if stockErr != nil {
		fmt.Println(stockErr)
	}

	newStock.IngredientId = 2
	newStock.Stock = 10
	_, stockErr = stockModel.CreateStockUpdate(newStock, int(newStock.IngredientId))
	if stockErr != nil {
		fmt.Println(stockErr)
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

func AuthInvalid(t *testing.T) string {
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

func TestCreateCheckoutAuthInvalidController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	checkoutModel := models.NewCheckoutModel(db)
	stockModel := models.NewStockModel(db)
	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	cartModel := models.NewCartModel(db)
	checkoutController := NewCheckoutController(checkoutModel, stockModel, recipeIngredientModel, cartModel)

	// Setting Route
	token := AuthInvalid(t)
	e := echo.New()
	e.POST("/api/checkouts/:recipeId", checkoutController.CreateCheckoutController, middleware.JWT([]byte(constants.SECRET_JWT)))
	// Add To Cart Controller
	req := httptest.NewRequest(echo.POST, "/api/checkouts/1", nil)
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

func TestCreateCheckoutBadRequestController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	checkoutModel := models.NewCheckoutModel(db)
	stockModel := models.NewStockModel(db)
	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	cartModel := models.NewCartModel(db)
	checkoutController := NewCheckoutController(checkoutModel, stockModel, recipeIngredientModel, cartModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.POST("/api/checkouts/:recipeId", checkoutController.CreateCheckoutController, middleware.JWT([]byte(constants.SECRET_JWT)))
	// Add To Cart Controller
	req := httptest.NewRequest(echo.POST, "/api/checkouts/satu", nil)
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

func TestCreateCheckoutController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)

	checkoutModel := models.NewCheckoutModel(db)
	stockModel := models.NewStockModel(db)
	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	cartModel := models.NewCartModel(db)
	checkoutController := NewCheckoutController(checkoutModel, stockModel, recipeIngredientModel, cartModel)

	// Setting Route
	token := AuthValid(t)
	e := echo.New()
	e.POST("/api/checkouts/:recipeId", checkoutController.CreateCheckoutController, middleware.JWT([]byte(constants.SECRET_JWT)))
	// Add To Cart Controller
	req := httptest.NewRequest(echo.POST, "/api/checkouts/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool            `json:"success"`
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    models.Checkout `json:"data"`
	}

	var response Response
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Create Checkout", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, uint(1), response.Data.ID)
}
