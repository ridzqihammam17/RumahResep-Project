package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"rumah_resep/config"
	"rumah_resep/models"
	"rumah_resep/util"
	"testing"

	"github.com/labstack/echo/v4"
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

	// -- Dummy Data
	var newUser models.User
	newUser.Name = "Testing Name"
	newUser.Email = "testingmailt@mail.com"
	newUser.Password = "generate111"
	newUser.Address = "jl. barat lau no 1"
	newUser.Gender = "laki"
	newUser.Role = "admin"

	// -- Dummy Data with Model
	userModel := models.NewUserModel(db)
	_, err := userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}
}

func TestValidRegisterUserController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := NewAuthController(userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/register", authController.RegisterUserController)

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"name":     "Buddy",
		"email":    "buddy@gmail.com",
		"password": "generate111",
		"address":  "jl barat daya no 5",
		"gender":   "laki",
		"role":     "admin",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/register", bytes.NewBuffer(reqBodyPost))
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
	assert.Equal(t, "Success Register Account", response.Message)
}

func TestInvalidRegisterUserControllerEmailNull(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := NewAuthController(userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/register", authController.RegisterUserController)

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"name":     "Buddy",
		"password": "generate111",
		"address":  "jl barat daya no 5",
		"gender":   "laki",
		"role":     "customer",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/register", bytes.NewBuffer(reqBodyPost))
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

func TestInvalidRegisterUserControllerAllNull(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := NewAuthController(userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/register", authController.RegisterUserController)

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/register", bytes.NewBuffer(reqBodyPost))
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
	assert.Equal(t, 500, res.Code)
	assert.Equal(t, "Internal Server Error", response.Message)
}

func TestValidLoginUserController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := NewAuthController(userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"email":    "buddy@gmail.com",
		"password": "generate111",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/login", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	e.ServeHTTP(res, req)

	// -- Declare Response and Convert to JSON
	type Response struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}
	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Login", response.Message)
	assert.NotEqual(t, "", response.Token)
}

func TestInvalidLoginUserController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	authController := NewAuthController(userModel)

	// -- Declare Route
	e := echo.New()
	e.POST("/api/login", authController.LoginUserController)

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]interface{}{
		"email":    "buddy@gmail.com",
		"password": "generate112",
	})

	// -- Setting Controller
	req := httptest.NewRequest(echo.POST, "/api/login", bytes.NewBuffer(reqBodyPost))
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
