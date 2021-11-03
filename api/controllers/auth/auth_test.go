package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	db := util.MysqlDatabaseConnection(config)

	// -- Clean DB Data
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

func TestRegisterUserController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := NewAuthController(userModel)

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]string{
		"name":     "Buddy",
		"email":    "buddy@gmail.com",
		"password": "generate111",
		"address":  "jl barat daya no 5",
		"gender":   "laki",
		"role":     "admin",
	})

	// -- Setting Controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/register")

	// -- Declare Controller
	userController.RegisterUserController(context)

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

func TestValidLoginUserController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := NewAuthController(userModel)

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]string{
		"email":    "buddy@gmail.com",
		"password": "generate111",
	})

	// -- Setting Controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/login")

	// -- Declare Controller
	userController.LoginUserController(context)

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
	assert.Equal(t, "Success Login", response.Message)
}

func TestInvalidLoginUserController(t *testing.T) {
	// -- Create Connection and Controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := NewAuthController(userModel)

	// -- Input
	reqBodyPost, _ := json.Marshal(map[string]string{
		"email":    "buddy@gmail.com",
		"password": "generate112",
	})

	// -- Setting Controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/login")

	// -- Declare Controller
	userController.LoginUserController(context)

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
