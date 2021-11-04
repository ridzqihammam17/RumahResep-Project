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
	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

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
	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

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
	userModel := models.NewUserModel(db)
	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	userController := auth.NewAuthController(userModel)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

	e := echo.New()

	e.POST("/api/login", userController.LoginUserController)
	e.POST("/api/recipes", recipesController.CreateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "admin@test.com",
		"password": "passAdmin",
	})

	req := httptest.NewRequest(echo.POST, "/api/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	type Response struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}
	var response Response
	resBody := rec.Body.String()
	json.Unmarshal([]byte(resBody), &response)

	assert.Equal(t, http.StatusOK, rec.Code)

	req2 := httptest.NewRequest(echo.POST, "/api/recipes", nil)
	rec2 := httptest.NewRecorder()
	req2.Header.Set("Authorization", fmt.Sprint("Bearer ", response.Data.(map[string]interface{})["token"]))
	e.ServeHTTP(rec2, req2)

	var response2 Response
	resBody2 := rec2.Body.String()
	json.Unmarshal([]byte(resBody2), &response2)

	assert.Equal(t, 200, rec2.Code)
	assert.Equal(t, "Success Create Recipe ", response2.Message)
}

func TestUpdateRecipeController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	userController := auth.NewAuthController(userModel)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

	e := echo.New()

	e.POST("/api/login", userController.LoginUserController)
	e.PUT("/api/recipes/:recipeId", recipesController.UpdateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "admin@test.com",
		"password": "passAdmin",
	})

	req := httptest.NewRequest(echo.POST, "/api/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	type Response struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}
	var response Response
	resBody := rec.Body.String()
	json.Unmarshal([]byte(resBody), &response)

	assert.Equal(t, http.StatusOK, rec.Code)

	req2 := httptest.NewRequest(echo.PUT, "/api/recipes/2", nil)
	rec2 := httptest.NewRecorder()
	req2.Header.Set("Authorization", fmt.Sprint("Bearer ", response.Data.(map[string]interface{})["token"]))
	e.ServeHTTP(rec2, req2)

	var response2 Response
	resBody2 := rec2.Body.String()
	json.Unmarshal([]byte(resBody2), &response2)

	assert.Equal(t, 200, rec2.Code)
	assert.Equal(t, "Success Update Recipe", response2.Message)
}

func TestDeleteRecipeController(t *testing.T) {
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	userModel := models.NewUserModel(db)
	recipeCategoriesModel := models.NewRecipesCategoriesModel(db)
	recipesModel := models.NewRecipeModel(db)
	categoriesModel := models.NewCategoryModel(db)
	userController := auth.NewAuthController(userModel)
	recipesController := NewRecipeController(recipeCategoriesModel, recipesModel, categoriesModel)

	e := echo.New()

	e.POST("/api/login", userController.LoginUserController)
	e.DELETE("/api/recipes/:recipeId", recipesController.CreateRecipeController, middleware.JWT([]byte(constants.SECRET_JWT)))

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "admin@test.com",
		"password": "passAdmin",
	})

	req := httptest.NewRequest(echo.POST, "/api/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	type Response struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}

	var response Response
	resBody := rec.Body.String()
	json.Unmarshal([]byte(resBody), &response)

	assert.Equal(t, http.StatusOK, rec.Code)

	req2 := httptest.NewRequest(echo.DELETE, "/api/recipes/2", nil)
	rec2 := httptest.NewRecorder()
	req2.Header.Set("Authorization", fmt.Sprint("Bearer ", response.Data.(map[string]interface{})["token"]))
	e.ServeHTTP(rec2, req2)

	var response2 Response
	resBody2 := rec2.Body.String()
	json.Unmarshal([]byte(resBody2), &response2)

	assert.Equal(t, 200, rec2.Code)
	assert.Equal(t, "Success Delete Recipe", response2.Message)
}
