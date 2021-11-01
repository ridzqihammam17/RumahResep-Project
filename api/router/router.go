package router

import (
	"rumah_resep/api/controllers/auth"
	"rumah_resep/constants"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Route(
	e *echo.Echo,
	authController *auth.AuthController,

) {
	// ------------------------------------------------------------------
	// Auth Login & Register
	// ------------------------------------------------------------------
	e.POST("/api/register", authController.RegisterUserController)
	e.POST("/api/login", authController.LoginUserController)

	// ------------------------------------------------------------------
	// Admin Role
	// ------------------------------------------------------------------
	eAdmin := e.Group("/api/admin")
	eAdmin.Use(middleware.JWT([]byte(constants.SECRET_JWT)))
}
