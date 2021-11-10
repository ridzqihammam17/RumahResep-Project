package carts

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
	db.Migrator().DropTable(&models.Cart{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Cart{})

	// -- Dummy Data
	var newUser models.User
	newUser.Name = "Yudi"
	newUser.Email = "yudi@mail.com"
	newUser.Password = "generate222"
	newUser.Address = "surabaya"
	newUser.Gender = "laki-laki"
	newUser.Role = "customer"

	// -- Dummy Data with Model
	userModel := models.NewUserModel(db)
	_, err := userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCreateCartController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnTest(config)
	userModel := models.NewUserModel(db)
	userController := auth.NewAuthController(userModel)
	cartModel := models.NewCartModel(db)
	cartController := NewCartController(cartModel)

	// Setting Route
	e := echo.New()
	e.POST("/api/login", userController.LoginUserController)
	e.POST("/api/carts", cartController.CreateCartController, middleware.JWT([]byte(constants.SECRET_JWT)))

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

	// Create Cart Controller
	req := httptest.NewRequest(echo.POST, "/api/carts", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", loginResponse.Token))
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

	assert.Equal(t, true, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Create cart success", response.Message)
}
