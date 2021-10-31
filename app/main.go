package main

import (
	"fmt"
	"rumah_resep/config"

	// "rumah_resep/util"
	"rumah_resep/api/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	//load config if available or set to default
	config := config.GetConfig()

	//initialize database connection based on given config
	// db := util.MysqlDatabaseConnection(config)

	//initiate model
	// userModel := models.NewUserModel(db)

	//initiate controller
	// newUserController := controllers.NewUserController(userModel)

	// newCheckoutController := controllers.NewCheckoutController(checkoutModel)
	//create echo http
	e := echo.New()

	//register API path and controller
	router.Route(e)

	// run server
	address := fmt.Sprintf("localhost:%d", config.Port)

	if err := e.Start(address); err != nil {
		log.Info("shutting down the server")
	}
}