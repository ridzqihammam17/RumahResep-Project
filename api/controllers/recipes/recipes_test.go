package recipes

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

	// -- Dummy Data
	var newUser models.User
	newUser.Name = "Customer A"
	newUser.Email = "customer@test.com"
	newUser.Password = "passCust"
	newUser.Address = "Gg. Gewart Dago No.3"
	newUser.Gender = "laki"
	newUser.Role = "customer"

	var newUser1 models.User
	newUser1.Name = "Admin A"
	newUser1.Email = "admin@test.com"
	newUser1.Password = "passAdmin"
	newUser1.Address = "jl. barat laut no 1"
	newUser1.Gender = "laki"
	newUser1.Role = "admin"

	var newRecipe models.Recipe
	newRecipe.Name = "Recipe A"

	// -- Dummy Data with Model
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

	recipeModel := models.NewRecipeModel(db)
	_, recipeModelErr := recipeModel.CreateRecipe(newRecipe)
	if recipeModelErr != nil {
		fmt.Println(recipeModelErr)
	}
}

func TestValidGetAllRecipeController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/recipes", recipesController.GetAllRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "customer@test.com",
		"password": "passCust",
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
	req := httptest.NewRequest(echo.GET, "/api/recipes", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
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
	assert.Equal(t, "Success Get All Recipe", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, "Recipe A", response.Data[0].Name)
}

func TestValidGetRecipeByIdController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/recipes/:recipeId", recipesController.GetRecipeByIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "customer@test.com",
		"password": "passCust",
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
	req := httptest.NewRequest(echo.GET, "/api/recipes/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
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
	assert.Equal(t, "Success Get Recipe", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, "Recipe A", response.Data.Name)
}

func TestInvalidGetRecipeByIdControllerNoId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/recipes/:recipeId", recipesController.GetRecipeByIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "customer@test.com",
		"password": "passCust",
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
	req := httptest.NewRequest(echo.GET, "/api/recipes/id", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestInvalidGetRecipeByIdControllerFalseId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/recipes/:recipeId", recipesController.GetRecipeByIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "customer@test.com",
		"password": "passCust",
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
	req := httptest.NewRequest(echo.GET, "/api/recipes/100", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, response.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestValidCreateRecipeController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/recipes", recipesController.CreateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "admin@test.com",
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

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"name": "Recipe B",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/recipes", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool          `json:"success"`
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    models.Recipe `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Create Recipe", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, "Recipe B", response.Data.Name)
}

func TestInvalidCreateRecipeControllerNotAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/recipes", recipesController.CreateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "customer@test.com",
		"password": "passCust",
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
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"name": "Recipe B",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/recipes", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 401, res.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestInvalidCreateRecipeControllerNameNull(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/recipes", recipesController.CreateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "admin@test.com",
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

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/recipes", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, res.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestValidUpdateRecipeController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/recipes/:recipeId", recipesController.UpdateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "admin@test.com",
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

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"name": "Recipe B Updated",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/recipes/2", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool          `json:"success"`
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    models.Recipe `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Update Recipe", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, uint(2), response.Data.ID)
	assert.Equal(t, "Recipe B Updated", response.Data.Name)
}

func TestInvalidUpdateRecipeControllerNotAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/recipes/:recipeId", recipesController.UpdateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "customer@test.com",
		"password": "passCust",
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
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"name": "Recipe B Updated",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/recipes/2", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 401, res.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestInvalidUpdateRecipeControllerNoId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/recipes/:recipeId", recipesController.UpdateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "admin@test.com",
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

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"name": "Recipe B Updated",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/recipes/id", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, res.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestInvalidUpdateRecipeControllerFalseId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/recipes/:recipeId", recipesController.UpdateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "admin@test.com",
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

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"name": "Recipe B Updated",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/recipes/200", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, res.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestInvalidUpdateRecipeControllerNameNull(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/recipes/:recipeId", recipesController.UpdateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "admin@test.com",
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

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/recipes/2", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, res.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestValidDeleteRecipeController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.DELETE("/api/recipes/:recipeId", recipesController.DeleteRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "admin@test.com",
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
	req := httptest.NewRequest(echo.DELETE, "/api/recipes/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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
	assert.Equal(t, "Success Delete Recipe", response.Message)
}

func TestInvalidDeleteRecipeControllerNotAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.DELETE("/api/recipes/:recipeId", recipesController.DeleteRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "customer@test.com",
		"password": "passCust",
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
	req := httptest.NewRequest(echo.DELETE, "/api/recipes/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 401, res.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestInvalidDeleteRecipeControllerNoId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.DELETE("/api/recipes/:recipeId", recipesController.DeleteRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "admin@test.com",
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
	req := httptest.NewRequest(echo.DELETE, "/api/recipes/id", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, res.Code)
	assert.Equal(t, "Bad Request", response.Message)
}

func TestInvalidDeleteRecipeControllerFalseId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.DELETE("/api/recipes/:recipeId", recipesController.DeleteRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	// ------ Start Login ------
	// -- Input
	reqBodyPostLogin, _ := json.Marshal(map[string]interface{}{
		"email":    "admin@test.com",
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
	req := httptest.NewRequest(echo.DELETE, "/api/recipes/300", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
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

	assert.Equal(t, false, response.Success)
	assert.Equal(t, 400, res.Code)
	assert.Equal(t, "Bad Request", response.Message)
}
