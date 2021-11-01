package main

import (
	"fmt"
	"rumah_resep/config"

	authController "rumah_resep/api/controllers/auth"
	"rumah_resep/models"
	"rumah_resep/util"

	// "rumah_resep/util"
	"rumah_resep/api/middlewares"
	"rumah_resep/api/router"

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

	//initiate controller
	newAuthController := authController.NewAuthController(userModel)

	//create echo http with log
	e := echo.New()
	middlewares.LoggerMiddlewares(e)

	//register API path and controller
	router.Route(
		e,
		newAuthController,
	)

	// run server
	address := fmt.Sprintf("localhost:%d", config.Port)

	if err := e.Start(address); err != nil {
		log.Info("shutting down the server")
	}
}
