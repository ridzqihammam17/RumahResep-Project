package recipes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	auth "rumah_resep/api/controllers/auth"
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
	// create database connection
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)

	// cleaning data before testing
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Recipe{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Recipe{})

	// preparate dummy data

	// dummy data for admin

	// dummy data for recipe
	var newRecipe models.Recipe
	newRecipe.Name = "Recipe A"
	// newRecipe.Categories = "1"

	recipeModel := models.NewRecipeModel(db)
	_, err := recipeModel.CreateRecipe(newRecipe)
	if err != nil {
		fmt.Println(err)
	}
}

func LoginForAllRole(t *testing.T) (token string) {
	// create database connection and create controller
	setup()
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)

	// dummy data for customer
	var newUser models.User
	newUser.Name = "Customer A"
	newUser.Email = "customer@test.com"
	newUser.Password = "passCust"
	newUser.Role = "customer"
	userModel := models.NewUserModel(db)
	_, err := userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}
	userController := auth.NewAuthController(userModel)

	// setting controller
	e := echo.New()
	reqBodyLogin, _ := json.Marshal(auth.LoginUserRequest{Email: "customer@test.com", Password: "passCust"})
	loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBodyLogin))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	loginContext := e.NewContext(loginReq, loginRes)
	loginContext.SetPath("/api/login")

	if err := userController.LoginUserController(loginContext); err != nil {
		t.Errorf("Shouldn't get error, get error: %s", err)
	}

	var user models.User
	json.Unmarshal(loginRes.Body.Bytes(), &user)

	assert.Equal(t, 200, loginRes.Code)
	assert.NotEqual(t, "", user.Token)

	return user.Token
}

func LoginForAdmin(t *testing.T) (c echo.Context, token string) {
	// create database connection and create controller
	setup()
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)

	// dummy data for admin
	var newUser models.User
	newUser.Name = "Admin A"
	newUser.Email = "admin@test.com"
	newUser.Password = "passAdmin"
	newUser.Role = "admin"
	userModel := models.NewUserModel(db)
	_, err := userModel.Register(newUser)
	if err != nil {
		fmt.Println(err)
	}
	userController := auth.NewAuthController(userModel)

	// setting controller
	e := echo.New()
	reqBodyLogin, _ := json.Marshal(auth.LoginUserRequest{Email: "admin@test.com", Password: "passAdmin"})
	loginReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBodyLogin))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	loginContext := e.NewContext(loginReq, loginRes)
	loginContext.SetPath("/api/login")

	if err := userController.LoginUserController(loginContext); err != nil {
		t.Errorf("Shouldn't get error, get error: %s", err)
	}

	var user models.User
	json.Unmarshal(loginRes.Body.Bytes(), &user)

	assert.Equal(t, 200, loginRes.Code)
	assert.NotEqual(t, "", user.Token)

	return loginContext, user.Token
}

func TestGetAllRecipeController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	token := LoginForAllRole(t)

	// setting controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/recipes")

	recipesController.GetAllRecipeController(context)

	var recipes []models.Recipe
	json.Unmarshal(res.Body.Bytes(), &recipes)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Recipe A", recipes[0].Name)
	// assert.Equal(t, 1, b[0].Category)
}

func TestGetRecipeByIdController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	token := LoginForAllRole(t)

	// setting controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/recipes/:recipeId")
	context.SetParamNames("recipeId")
	context.SetParamValues("1")
	// context.
	recipesController.GetRecipeByIdController(context)
	fmt.Println(context)

	var recipe models.Recipe
	json.Unmarshal(res.Body.Bytes(), &recipe)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Recipe A", recipe.Name)

}

func TestCreateRecipeController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	ctx, token := LoginForAdmin(t)
	fmt.Println(token)

	reqBodyPost, _ := json.Marshal(map[string]string{
		"name": "Recipe B",
	})

	// setting controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBodyPost))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.Set("user", ctx.Get("user"))
	context.SetPath("/api/recipes")

	// fmt.Println(context)

	if err := recipesController.CreateRecipeController(context); err != nil {
		t.Errorf("Should'nt get error, get error: %s", err)
	}

	// fmt.Println()

	type Response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var response Response
	fmt.Println(response)
	json.Unmarshal(res.Body.Bytes(), &response)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Create Recipe Success", response.Message)
}

// func TestUpdateRecipeController(t *testing.T) {
// 	config := config.GetConfig()
// 	db := util.MysqlDatabaseConnection(config)
// 	recipesModel := models.NewRecipeModel(db)
// 	recipesController := NewRecipeController(recipesModel)

// 	e, token := LoginForAdmin(t)
// 	// fmt.Println(token)

// 	reqBodyPost, _ := json.Marshal(map[string]string{
// 		"name": "Recipe B Updated",
// 	})

// 	// setting controller
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBodyPost))
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
// 	req.Header.Set("Content-Type", "application/json")
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath("/api/recipes/:recipeId")
// 	context.SetParamNames("recipeId")
// 	context.SetParamValues("1")

// 	// fmt.Println(context)

// 	if err := recipesController.UpdateRecipeController(context); err != nil {
// 		t.Errorf("Should'nt get error, get error: %s", err)
// 	}

// 	// fmt.Println()

// 	type Response struct {
// 		Code    int    `json:"code"`
// 		Message string `json:"message"`
// 	}

// 	var response Response
// 	fmt.Println(response)
// 	json.Unmarshal(res.Body.Bytes(), &response)

// 	assert.Equal(t, 200, res.Code)
// 	assert.Equal(t, "Success Update Recipe", response.Message)
// }

// func TestDeleteRecipeController(t *testing.T) {
// 	config := config.GetConfig()
// 	db := util.MysqlDatabaseConnection(config)
// 	recipesModel := models.NewRecipeModel(db)
// 	recipesController := NewRecipeController(recipesModel)

// 	token := LoginForAdmin(t)
// 	// fmt.Println(token)

// 	// setting controller
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodDelete, "/", nil)
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
// 	req.Header.Set("Content-Type", "application/json")
// 	res := httptest.NewRecorder()
// 	context := e.NewContext(req, res)
// 	context.SetPath("/api/recipes/:recipeId")
// 	context.SetParamNames("recipeId")
// 	context.SetParamValues("1")

// 	// fmt.Println(context)

// 	if err := recipesController.DeleteRecipeController(context); err != nil {
// 		t.Errorf("Should'nt get error, get error: %s", err)
// 	}

// 	// fmt.Println()

// 	type Response struct {
// 		Code    int    `json:"code"`
// 		Message string `json:"message"`
// 	}

// 	var response Response
// 	fmt.Println(response)
// 	json.Unmarshal(res.Body.Bytes(), &response)

// 	assert.Equal(t, 200, res.Code)
// 	assert.Equal(t, "Success Delete Recipe", response.Message)
// }
