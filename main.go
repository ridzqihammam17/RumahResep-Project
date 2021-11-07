package main

import (
	"fmt"
	"rumah_resep/api/middlewares"
	"rumah_resep/api/router"
	"rumah_resep/config"
	"rumah_resep/models"
	"rumah_resep/util"

	authControllers "rumah_resep/api/controllers/auth"
	cartDetailsControllers "rumah_resep/api/controllers/cartDetails"
	cartControllers "rumah_resep/api/controllers/carts"
	categoryControllers "rumah_resep/api/controllers/categories"
	checkoutControllers "rumah_resep/api/controllers/checkouts"
	ingredientControllers "rumah_resep/api/controllers/ingredient"
	midtransControllers "rumah_resep/api/controllers/midtrans"
	recipeIngredientsControllers "rumah_resep/api/controllers/recipeIngredients"
	recipeControllers "rumah_resep/api/controllers/recipes"
	recipeCategoriesControllers "rumah_resep/api/controllers/recipesCategories"
	stockControllers "rumah_resep/api/controllers/stocks"
	transactionControllers "rumah_resep/api/controllers/transactions"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	//load config if available or set to default
	config := config.GetConfig()

	//initialize database connection based on given config
	db := util.MysqlDatabaseConnection(config)

	//initiate model
	cartModel := models.NewCartModel(db)
	userModel := models.NewUserModel(db)
	recipeModel := models.NewRecipeModel(db)
	categoryModel := models.NewCategoryModel(db)
	recipesCategoriesModel := models.NewRecipesCategoriesModel(db)
	ingredientModel := models.NewIngredientModel(db)
	recipeIngredientsModel := models.NewRecipeIngredientsModel(db)
	stockModel := models.NewStockModel(db)
	cartDetailsModel := models.NewCartDetailsModel(db)
	checkoutModel := models.NewCheckoutModel(db)
	transactionModel := models.NewTransactionModel(db)
	// midtransModel := models.NewMidtransModel(db)

	//initiate controller
	newCartController := cartControllers.NewCartController(cartModel)
	newAuthController := authControllers.NewAuthController(userModel)
	newRecipeController := recipeControllers.NewRecipeController(recipesCategoriesModel, recipeModel, categoryModel)
	newCategoryController := categoryControllers.NewCategoryController(categoryModel)
	newRecipesCategoriesController := recipeCategoriesControllers.NewRecipesCategoriesController(recipesCategoriesModel, recipeModel, categoryModel)
	// newStockController := stockControllers.NewStrockController(stockModel)
	newIngredientController := ingredientControllers.NewIngredientController(ingredientModel, stockModel)
	newRecipeIngredientsController := recipeIngredientsControllers.NewRecipeIngredientsController(recipeIngredientsModel, recipeModel, ingredientModel)
	newCartDetailsController := cartDetailsControllers.NewCartDetailsController(cartDetailsModel, recipeModel, ingredientModel, recipeIngredientsModel, cartModel)
	newCheckoutController := checkoutControllers.NewCheckoutController(checkoutModel)
	newTransactionController := transactionControllers.NewTransactionController(transactionModel, cartModel, userModel)
	newMidtransController := midtransControllers.NewMidtransController(transactionModel)
	newStockController := stockControllers.NewStockController(stockModel)

	//create echo http with log
	e := echo.New()
	middlewares.LoggerMiddlewares(e)

	//register API path and controller
	router.Route(e, newAuthController, newCartController, newRecipeController, newCategoryController, newRecipesCategoriesController, newIngredientController, newRecipeIngredientsController, newCartDetailsController, newCheckoutController, newTransactionController, newMidtransController, newStockController)

	// run server
	address := fmt.Sprintf(":%d", config.Port)

	if err := e.Start(address); err != nil {
		log.Info("shutting down the server")
	}
}
