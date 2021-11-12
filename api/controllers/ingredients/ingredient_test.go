package ingredients

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
	db.Migrator().DropTable(&models.Ingredient{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Ingredient{})

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

	var newIngredient models.Ingredient
	newIngredient.Name = "Daging ayam dada"
	newIngredient.Price = 30000

	var newIngredient1 models.Ingredient
	newIngredient1.Name = "Santan Kelapa"
	newIngredient1.Price = 3000

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

	ingredientModel := models.NewIngredientModel(db)
	_, ingredientModelErr := ingredientModel.CreateIngredient(newIngredient)
	if ingredientModelErr != nil {
		fmt.Println(ingredientModelErr)
	}

	ingredientModel1 := models.NewIngredientModel(db)
	_, ingredientModel1Err := ingredientModel1.CreateIngredient(newIngredient1)
	if ingredientModel1Err != nil {
		fmt.Println(ingredientModel1Err)
	}
}

func TestValidGetAllIngredientController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/ingredients", ingredientController.GetAllIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.GET, "/api/ingredients", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool                `json:"success"`
		Code    int                 `json:"code"`
		Message string              `json:"message"`
		Data    []models.Ingredient `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Get All Ingredient", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, uint(1), response.Data[0].ID)
	assert.Equal(t, "Daging ayam dada", response.Data[0].Name)
	assert.Equal(t, 30000, response.Data[0].Price)
	assert.Equal(t, uint(2), response.Data[1].ID)
	assert.Equal(t, "Santan Kelapa", response.Data[1].Name)
	assert.Equal(t, 3000, response.Data[1].Price)
}

func TestInvalidGetAllIngredientControllerNotAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/ingredients", ingredientController.GetAllIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.GET, "/api/ingredients", nil)
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
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestValidGetIngredientByIdController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/ingredients/:ingredientId", ingredientController.GetIngredientByIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.GET, "/api/ingredients/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool              `json:"success"`
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    models.Ingredient `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Get Ingredient", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, uint(1), response.Data.ID)
	assert.Equal(t, "Daging ayam dada", response.Data.Name)
	assert.Equal(t, 30000, response.Data.Price)
}

func TestInvalidGetIngredientByIdControllerNotAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/ingredients/:ingredientId", ingredientController.GetIngredientByIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.GET, "/api/ingredients/1", nil)
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
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestInvalidGetIngredientByIdControllerNoId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/ingredients/:ingredientId", ingredientController.GetIngredientByIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.GET, "/api/ingredients/id", nil)
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

func TestInvalidGetIngredientByIdControllerFalseId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.GET("/api/ingredients/:ingredientId", ingredientController.GetIngredientByIdController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.GET, "/api/ingredients/200", nil)
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

func TestValidCreateIngredientController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/ingredients", ingredientController.CreateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
		"name":  "Buah Apel",
		"price": 1000,
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/ingredients", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool              `json:"success"`
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    models.Ingredient `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Create Ingredient", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, "Buah Apel", response.Data.Name)
	assert.Equal(t, 1000, response.Data.Price)
}

func TestInvalidCreateIngredientControllerNotAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/ingredients", ingredientController.CreateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
		"name":  "Buah Apel",
		"price": 1000,
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/ingredients", bytes.NewBuffer(reqBodyPost))
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
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestInvalidCreateIngredientControllerAllNull(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.POST("/api/ingredients", ingredientController.CreateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.POST, "/api/ingredients", bytes.NewBuffer(reqBodyPost))
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

func TestValidUpdateIngredientController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/ingredients/:ingredientId", ingredientController.UpdateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
		"name":  "Bawang Merah",
		"price": 1000,
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/ingredients/1", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", responseLogin.Token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool              `json:"success"`
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    models.Ingredient `json:"data"`
	}

	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Update Ingredient", response.Message)
	assert.NotEmpty(t, response.Data)
	assert.Equal(t, uint(1), response.Data.ID)
	assert.Equal(t, "Bawang Merah", response.Data.Name)
	assert.Equal(t, 1000, response.Data.Price)
}

func TestInvalidUpdateIngredientControllerNotAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/ingredients/:ingredientId", ingredientController.UpdateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
		"name":  "Bawang Merah",
		"price": 1000,
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/ingredients/1", bytes.NewBuffer(reqBodyPost))
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
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestInvalidUpdateIngredientControllerNoId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/ingredients/:ingredientId", ingredientController.UpdateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
		"name":  "Bawang Merah",
		"price": 1000,
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/ingredients/id", bytes.NewBuffer(reqBodyPost))
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

func TestInvalidUpdateIngredientControllerFalseId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/ingredients/:ingredientId", ingredientController.UpdateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
		"name":  "Bawang Merah",
		"price": 1000,
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.PUT, "/api/ingredients/100", bytes.NewBuffer(reqBodyPost))
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

func TestInvalidUpdateIngredientControllerAllNull(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.PUT("/api/ingredients/:ingredientId", ingredientController.UpdateIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.PUT, "/api/ingredients/1", bytes.NewBuffer(reqBodyPost))
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

func TestValidDeleteIngredientController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.DELETE("/api/ingredients/:ingredientId", ingredientController.DeleteIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.DELETE, "/api/ingredients/1", nil)
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
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success Delete Ingredient", response.Message)
}

func TestInvalidDeleteIngredientControllerNotAdmin(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.DELETE("/api/ingredients/:ingredientId", ingredientController.DeleteIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.DELETE, "/api/ingredients/1", nil)
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
	assert.Equal(t, 401, response.Code)
	assert.Equal(t, "Unauthorized Error", response.Message)
}

func TestInvalidDeleteIngredientControllerNoId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.DELETE("/api/ingredients/:ingredientId", ingredientController.DeleteIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.DELETE, "/api/ingredients/id", nil)
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

func TestInvalidDeleteIngredientControllerFalseId(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := auth.NewAuthController(userModel)
	ingredientModel := models.NewIngredientModel(db)
	ingredientController := NewIngredientController(ingredientModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)
	e.DELETE("/api/ingredients/:ingredientId", ingredientController.DeleteIngredientController, middleware.JWT([]byte(constants.SECRET_JWT)))

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
	req := httptest.NewRequest(echo.DELETE, "/api/ingredients/100", nil)
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
