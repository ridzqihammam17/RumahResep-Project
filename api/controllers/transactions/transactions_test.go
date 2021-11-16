package transactions

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
	db.Migrator().DropTable(&models.RecipeIngredients{})
	db.Migrator().DropTable(&models.Ingredient{})
	db.Migrator().DropTable(&models.Recipe{})
	db.Migrator().DropTable(&models.User{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Recipe{})
	db.AutoMigrate(&models.Ingredient{})
	db.AutoMigrate(&models.RecipeIngredients{})
	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.CartDetails{})
	db.AutoMigrate(&models.Transaction{})

	// -- Dummy Data with Model
	// ------ Start User ------
	var newUser models.User
	newUser.Name = "Rudi"
	newUser.Address = "Bekasi"
	newUser.Email = "rudi@test.com"
	newUser.Password = "rudi99"
	newUser.Gender = "laki"
	newUser.Role = "customer"

	var newUser1 models.User
	newUser1.Name = "Budi"
	newUser1.Email = "budi@mail.com"
	newUser1.Password = "passAdmin"
	newUser1.Address = "jl. barat laut no 1"
	newUser1.Gender = "laki"
	newUser1.Role = "admin"

	var newUser2 models.User
	newUser2.Name = "Rudi2"
	newUser2.Address = "Bekasi2"
	newUser2.Email = "rudi2@test.com"
	newUser2.Password = "rudi99"
	newUser2.Gender = "laki"
	newUser2.Role = "customer"

	userModel := models.NewUserModel(db)
	_, userModelErr := userModel.Register(newUser)
	if userModelErr != nil {
		fmt.Println(userModelErr)
	}

	userModel1 := models.NewUserModel(db)
	_, userModelErr1 := userModel1.Register(newUser1)
	if userModelErr1 != nil {
		fmt.Println(userModelErr1)
	}

	userModel2 := models.NewUserModel(db)
	_, userModelErr2 := userModel2.Register(newUser2)
	if userModelErr2 != nil {
		fmt.Println(userModelErr2)
	}
	// ------ End User ------

	// ------ Start Recipe ------
	var newRecipe models.Recipe
	newRecipe.Name = "Rujak Cingur"

	recipeModel := models.NewRecipeModel(db)
	_, recipeModelErr := recipeModel.CreateRecipe(newRecipe)
	if recipeModelErr != nil {
		fmt.Println(recipeModelErr)
	}
	// ------ End Recipe ------

	// ------ Start Ingredient
	var newIngredient models.Ingredient
	newIngredient.Name = "Bawah merah"
	newIngredient.Price = 2500

	var newIngredient1 models.Ingredient
	newIngredient1.Name = "Bawah putih"
	newIngredient1.Price = 3000

	ingredientModel := models.NewIngredientModel(db)
	_, ingredientModelErr := ingredientModel.CreateIngredient(newIngredient)
	if ingredientModelErr != nil {
		fmt.Println(ingredientModelErr)
	}

	ingredientModel1 := models.NewIngredientModel(db)
	_, ingredientModelErr1 := ingredientModel1.CreateIngredient(newIngredient1)
	if ingredientModelErr1 != nil {
		fmt.Println(ingredientModelErr1)
	}
	// ------ End Ingredient ------

	// ------ Start RecipeIngredient ------
	var newRecipeIngredient models.RecipeIngredients
	newRecipeIngredient.IngredientId = 1
	newRecipeIngredient.RecipeId = 1
	newRecipeIngredient.QtyIngredient = 3

	var newRecipeIngredient1 models.RecipeIngredients
	newRecipeIngredient1.IngredientId = 2
	newRecipeIngredient1.RecipeId = 1
	newRecipeIngredient1.QtyIngredient = 4

	recipeIngredientModel := models.NewRecipeIngredientsModel(db)
	_, recipeIngredientModelErr := recipeIngredientModel.AddIngredientsRecipe(newRecipeIngredient)
	if recipeIngredientModelErr != nil {
		fmt.Println(recipeIngredientModelErr)
	}

	recipeIngredientModel1 := models.NewRecipeIngredientsModel(db)
	_, recipeIngredientModelErr1 := recipeIngredientModel1.AddIngredientsRecipe(newRecipeIngredient1)
	if recipeIngredientModelErr1 != nil {
		fmt.Println(recipeIngredientModelErr1)
	}
	// ------ End RecipeIngredient ------

	// ------ Start Cart ------
	var newCart models.Cart
	newCart.UserID = 1

	cartModel := models.NewCartModel(db)
	_, cartModelErr := cartModel.CreateCart(newCart, 1)
	if cartModelErr != nil {
		fmt.Println(cartModelErr)
	}
	// ------ End Cart ------

	// ------ Start Cart Detail ------
	var newCartDetails models.CartDetails
	newCartDetails.CartID = int(1)
	newCartDetails.RecipeID = int(1)
	newCartDetails.Quantity = 2
	newCartDetails.Price = 50000
	newCartDetails.CheckoutID = 1

	cartDetailModel := models.NewCartDetailsModel(db)
	_, cartDetailModelErr := cartDetailModel.AddRecipeToCart(newCartDetails)
	if cartDetailModelErr != nil {
		fmt.Println(cartDetailModelErr)
	}
	// ------ End Cart Detail ------
}

func TestValidGetAllTransactionAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	cartModel := models.NewCartModel(db)
	transactionController := NewTransactionController(transactionModel, cartModel, userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/transactions", transactionController.GetAllTransactionAdmin, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "budi@mail.com",
		"password": "passAdmin",
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
	req := httptest.NewRequest(echo.GET, "/api/transactions", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool               `json:"success"`
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    models.Transaction `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Get All Transaction", response.Message)
}

func TestInvalidGetAllTransactionAdminNotAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	cartModel := models.NewCartModel(db)
	transactionController := NewTransactionController(transactionModel, cartModel, userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/transactions", transactionController.GetAllTransactionAdmin, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "rudi@test.com",
		"password": "rudi99",
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
	req := httptest.NewRequest(echo.GET, "/api/transactions", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool               `json:"success"`
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    models.Transaction `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestValidGetAllTransaction(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	cartModel := models.NewCartModel(db)
	transactionController := NewTransactionController(transactionModel, cartModel, userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/transactions/list", transactionController.GetAllTransactionCustomer, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "rudi@test.com",
		"password": "rudi99",
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
	req := httptest.NewRequest(echo.GET, "/api/transactions/list", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool               `json:"success"`
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    models.Transaction `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Get All Transaction", response.Message)
}

func TestInvalidGetAllTransactionNotCustomer(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	cartModel := models.NewCartModel(db)
	transactionController := NewTransactionController(transactionModel, cartModel, userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/transactions/list", transactionController.GetAllTransactionCustomer, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "budi@mail.com",
		"password": "passAdmin",
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
	req := httptest.NewRequest(echo.GET, "/api/transactions/list", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool               `json:"success"`
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    models.Transaction `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestValidCreateTransaction(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	cartModel := models.NewCartModel(db)
	transactionController := NewTransactionController(transactionModel, cartModel, userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/transactions", transactionController.CreateTransaction, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "rudi@test.com",
		"password": "rudi99",
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
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"shipping_method": "delivery",
	})

	req := httptest.NewRequest(echo.POST, "/api/transactions", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	type Response struct {
		Success bool               `json:"success"`
		Code    int                `json:"code"`
		Message string             `json:"message"`
		Data    models.Transaction `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Create Transaction", response.Message)
	assert.NotEmpty(t, response.Data)
}

func TestInvalidCreateTransactionNotCustomer(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	cartModel := models.NewCartModel(db)
	transactionController := NewTransactionController(transactionModel, cartModel, userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/transactions", transactionController.CreateTransaction, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "budi@mail.com",
		"password": "passAdmin",
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
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"shipping_method": "delivery",
	})

	req := httptest.NewRequest(echo.POST, "/api/transactions", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

func TestInvalidCreateTransactionNoCartId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	cartModel := models.NewCartModel(db)
	transactionController := NewTransactionController(transactionModel, cartModel, userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/transactions", transactionController.CreateTransaction, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "rudi2@test.com",
		"password": "rudi99",
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
	reqBodyPost, _ := json.Marshal(map[string]interface{}{})

	req := httptest.NewRequest(echo.POST, "/api/transactions", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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
	assert.Equal(t, 500, response.Code)
	assert.Equal(t, "Internal Server Error", response.Message)
}

func TestInvalidCreateTransactionNoShippingMethod(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	transactionModel := models.NewTransactionModel(db)
	cartModel := models.NewCartModel(db)
	transactionController := NewTransactionController(transactionModel, cartModel, userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/transactions", transactionController.CreateTransaction, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "rudi@test.com",
		"password": "rudi99",
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
	reqBodyPost, _ := json.Marshal(map[string]interface{}{})

	req := httptest.NewRequest(echo.POST, "/api/transactions", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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
