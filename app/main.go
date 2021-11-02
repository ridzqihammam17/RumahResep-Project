package main

import (
	"fmt"
	"rumah_resep/api/middlewares"
	"rumah_resep/api/router"
	"rumah_resep/config"
	"rumah_resep/models"
	"rumah_resep/util"

	authControllers "rumah_resep/api/controllers/auth"
	cartControllers "rumah_resep/api/controllers/carts"
	recipeControllers "rumah_resep/api/controllers/recipes"
	categoryControllers "rumah_resep/api/controllers/categories"

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

	//initiate controller
	newCartController := cartControllers.NewCartController(cartModel)
	newAuthController := authControllers.NewAuthController(userModel)
	newRecipeController := recipeControllers.NewRecipeController(recipeModel)
	newCategoryController := categoryControllers.NewCategoryController(categoryModel)

	//create echo http with log
	e := echo.New()
	middlewares.LoggerMiddlewares(e)

	//register API path and controller
	router.Route(e, newAuthController, newCartController, newRecipeController, newCategoryController)

	// run server
	address := fmt.Sprintf("localhost:%d", config.Port)

	if err := e.Start(address); err != nil {
		log.Info("shutting down the server")
	}
}
