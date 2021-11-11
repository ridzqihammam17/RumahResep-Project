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
	db := util.MysqlDatabaseConnTest(config)

	// -- Clean DB Data
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Recipe{})
	db.Migrator().DropTable(&models.Ingredient{})
	db.Migrator().DropTable(&models.RecipeIngredients{})
	db.Migrator().DropTable(&models.Cart{})
	db.Migrator().DropTable(&models.CartDetails{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Recipe{})
	db.AutoMigrate(&models.Ingredient{})
	db.AutoMigrate(&models.RecipeIngredients{})
	db.AutoMigrate(&models.Cart{})
	db.AutoMigrate(&models.CartDetails{})

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

}

func TestAddRecipeToCartController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
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
	assert.Equal(t, 6500, response.Data.Price)
}

func TestGetAllRecipeByCartIdController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
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
	assert.Equal(t, 6500, response.Data[0].Price)
}

func TestUpdateRecipePortionController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
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
	assert.Equal(t, 13000, response.Data.Price)
}

func TestDeleteRecipeFromCartController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
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
