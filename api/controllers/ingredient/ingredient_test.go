package ingredient

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
	db.Migrator().DropTable(&models.Ingredient{})
	db.AutoMigrate(&models.Ingredient{})

	// preparate dummy data
	var newIngredient models.Ingredient
	newIngredient.Name = "Daging ayam dada"
	newIngredient.Stock = 10
	newIngredient.Price = 30000

	// recipe dummy data with model
	ingredientModel := models.NewIngredientModel(db)
	_, err := ingredientModel.CreateIngredient(newIngredient)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetAllRecipeController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	ingredientModel := models.NewIngredientModel(db)
	stockModel := models.NewStockModel(db)
	ingredientController := NewIngredientController(ingredientModel, stockModel)

	// setting controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/ingredient")

	ingredientController.GetAllIngredientController(context)

	var b []models.Ingredient
	json.Unmarshal(res.Body.Bytes(), &b)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "Daging ayam dada", b[0].Name)
	assert.Equal(t, 10, b[0].Stock)
	assert.Equal(t, 30000, b[0].Price)
}

func TestGetRecipeByIdController(t *testing.T) {
	// create database connection and create controller
	config := config.GetConfig()
	db := util.MysqlDatabaseConnection(config)
	ingredientModel := models.NewIngredientModel(db)
	stockModel := models.NewStockModel(db)
	ingredientController := NewIngredientController(ingredientModel, stockModel)

	// setting controller
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	context := e.NewContext(req, res)
	context.SetPath("/api/ingredient/:ingredientId")
	context.SetParamNames("recipeId")
	context.SetParamValues("1")

	ingredientController.GetIngredientByIdController(context)

	var bList []models.Ingredient
	json.Unmarshal(res.Body.Bytes(), &bList)

	assert.Equal(t, 200, res.Code)
}
