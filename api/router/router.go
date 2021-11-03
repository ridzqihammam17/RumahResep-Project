package router

import (
	"rumah_resep/api/controllers/auth"
	"rumah_resep/api/controllers/carts"
	"rumah_resep/api/controllers/categories"
	"rumah_resep/api/controllers/recipes"
	"rumah_resep/constants"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Route(
	e *echo.Echo,
	authController *auth.AuthController,
	cartController *carts.CartController,
	recipeController *recipes.RecipeController,
	categoryController *categories.CategoryController,
) {
	// ------------------------------------------------------------------
	// Auth Login & Register
	// ------------------------------------------------------------------
	e.POST("/api/register", authController.RegisterUserController)
	e.POST("/api/login", authController.LoginUserController)

	// Auth JWT
	jwtMiddleware := middleware.JWT([]byte(constants.SECRET_JWT))

	// ------------------------------------------------------------------
	// Carts
	// ------------------------------------------------------------------
	e.POST("/api/carts", cartController.CreateCartController, jwtMiddleware)
	e.GET("/api/carts/:id", cartController.GetCartController, jwtMiddleware)
	e.PUT("/api/carts/:id", cartController.UpdateCartController, jwtMiddleware)
	e.DELETE("/api/carts/:id", cartController.DeleteCartController, jwtMiddleware)

	// Recipe
	e.GET("/api/recipes", recipeController.GetAllRecipeController, jwtMiddleware)
	e.GET("/api/recipes/:recipeId", recipeController.GetRecipeByIdController, jwtMiddleware)
	e.POST("/api/recipes", recipeController.CreateRecipeController, jwtMiddleware)
	e.PUT("/api/recipes/:recipeId", recipeController.UpdateRecipeController, jwtMiddleware)
	e.DELETE("/api/recipes/:recipeId", recipeController.DeleteRecipeController, jwtMiddleware)
	e.GET("/api/recipes/category/:categoryId", recipeController.GetRecipeByCategoryIdController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Categories
	// ------------------------------------------------------------------
	e.GET("/api/categories", categoryController.GetAllCategoryController, jwtMiddleware)
	e.POST("/api/categories", categoryController.InsertCategoryController, jwtMiddleware)
	e.GET("/api/categories/:id", categoryController.GetCategoryController, jwtMiddleware)
	e.PUT("/api/categories/:id", categoryController.EditCategoryController, jwtMiddleware)
	e.DELETE("/api/categories/:id", categoryController.DeleteCategoryController, jwtMiddleware)

}
