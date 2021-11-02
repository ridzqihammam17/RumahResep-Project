package main

import (
	"fmt"
	"rumah_resep/api/middlewares"
	"rumah_resep/api/router"
	"rumah_resep/config"
	"rumah_resep/models"
	"rumah_resep/util"

	authControllers "rumah_resep/api/controllers/auth"
	cartDetailControllers "rumah_resep/api/controllers/cartdetails"
	cartControllers "rumah_resep/api/controllers/carts"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	//load config if available or set to default
	config := config.GetConfig()

	//initialize database connection based on given config
	db := util.MysqlDatabaseConnection(config)

	//initiate model
	userModel := models.NewUserModel(db)
	cartModel := models.NewCartModel(db)
	cartDetailModel := models.NewCartDetailModel(db)

	//initiate controller
	newAuthController := authControllers.NewAuthController(userModel)
	newCartController := cartControllers.NewCartController(cartModel)
	newCartDetailController := cartDetailControllers.NewCartDetailController(cartModel, cartDetailModel)

	//create echo http with log
	e := echo.New()
	middlewares.LoggerMiddlewares(e)

	//register API path and controller
	router.Route(e, newAuthController, newCartController, newCartDetailController)

	// run server
	address := fmt.Sprintf("localhost:%d", config.Port)

	if err := e.Start(address); err != nil {
		log.Info("shutting down the server")
	}
}
