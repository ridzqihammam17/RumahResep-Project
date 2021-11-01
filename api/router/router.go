package router

import (
	"rumah_resep/api/controllers/auth"
	"rumah_resep/api/controllers/carts"
	"rumah_resep/constants"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Route(
	e *echo.Echo,
	authController *auth.AuthController,
	cartController *carts.CartController,

) {
	// ------------------------------------------------------------------
	// Auth Login & Register
	// ------------------------------------------------------------------
	e.POST("/api/register", authController.RegisterUserController)
	e.POST("/api/login", authController.LoginUserController)
	// Auth JWT
	jwtMiddleware := middleware.JWT([]byte(constants.SECRET_JWT))

	e.POST("/api/carts", cartController.CreateCartController, jwtMiddleware)
	e.GET("/api/carts/:id", cartController.GetCartController, jwtMiddleware)
	e.PUT("/api/carts/:id", cartController.UpdateCartController, jwtMiddleware)
}
