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
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)

	db.Migrator().DropTable(&models.User{})
	db.AutoMigrate(&models.User{})

	var newUser models.User
	newUser.Name = "Testing Name"
	newUser.Email = "testingmailt@mail.com"
	newUser.Password = "generate111"

	userModel := models.NewUserModel(db)
	_, err := userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}
}

func TestRegisterUserController(t *testing.T) {
	setup()
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := NewAuthController(userModel)

	e := echo.New()

	reqBodyPost, _ := json.Marshal(map[string]string{
		"name":     "Buddy",
		"email":    "buddy@gmail.com",
		"password": "generate111",
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/register")

	if err := userController.RegisterUserController(context); err != nil {
		t.Errorf("Should'nt get error, get error: %s", err)
	}

	type Response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	var response Response

	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Success Register Account", response.Message)
}

func TestLoginUserController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	userController := NewAuthController(userModel)

	// // setting controller
	e := echo.New()
	reqBodyLogin, _ := json.Marshal(models.User{Email: "testingmailt@mail.com", Password: "generate111"})
	loginreq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBodyLogin))
	loginreq.Header.Set("Content-Type", "application/json")
	loginres := httptest.NewRecorder()
	logincontext := e.NewContext(loginreq, loginres)
	logincontext.SetPath("/api/login")

	if err := userController.LoginUserController(logincontext); err != nil {
		t.Errorf("Should'nt get error, get error: %s", err)
	}

	var c models.User
	json.Unmarshal(loginres.Body.Bytes(), &c)

	assert.Equal(t, 200, loginres.Code)
}
