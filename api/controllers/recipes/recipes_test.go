package recipes

import (
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
	// create database connection
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)

	// cleaning data before testing
	db.Migrator().DropTable(&models.Recipe{})
	db.AutoMigrate(&models.Recipe{})

	// preparate dummy data
	var newRecipe models.Recipe
	newRecipe.Name = "Soto"
	newRecipe.Category = 1

	// recipe dummy data with model
	recipeModel := models.NewRecipeModel(db)
	_, err := recipeModel.CreateRecipe(newRecipe)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetAllRecipeController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// setting controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/recipes")

	recipesController.GetAllRecipeController(context)

	var b []models.Recipe
	json.Unmarshal(res.Body.Bytes(), &b)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Soto", b[0].Name)
	assert.Equal(t, 1, b[0].Category)
}

func TestGetRecipeByIdController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	recipesModel := models.NewRecipeModel(db)
	recipesController := NewRecipeController(recipesModel)

	// setting controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/recipes/:recipeId")
	context.SetParamNames("recipeId")
	context.SetParamValues("1")

	recipesController.GetRecipeByIdController(context)

	var bList []models.Recipe
	json.Unmarshal(res.Body.Bytes(), &bList)

	assert.Equal(t, 200, res.Code)
}
