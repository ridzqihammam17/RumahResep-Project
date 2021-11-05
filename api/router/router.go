package router

import (
	"rumah_resep/api/controllers/auth"
	cartdetails "rumah_resep/api/controllers/cartDetails"
	"rumah_resep/api/controllers/carts"
	"rumah_resep/api/controllers/categories"
	"rumah_resep/api/controllers/ingredient"
	recipeingredients "rumah_resep/api/controllers/recipeIngredients"
	"rumah_resep/api/controllers/recipes"
	recipescategories "rumah_resep/api/controllers/recipesCategories"

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
	recipesCategroriesController *recipescategories.RecipesCategoriesController,
	ingredientController *ingredient.IngredientController,
	recipeIngredientsController *recipeingredients.RecipeIngredientsController,
	cartDetailsController *cartdetails.CartDetailsController,
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
	// e.GET("/api/carts/:id", cartController.GetCartController, jwtMiddleware)
	// e.PUT("/api/carts/:id", cartController.UpdateCartController, jwtMiddleware)
	// e.DELETE("/api/carts/:id", cartController.DeleteCartController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Recipe
	// ------------------------------------------------------------------
	e.GET("/api/recipes", recipeController.GetAllRecipeController, jwtMiddleware)
	e.GET("/api/recipes/:recipeId", recipeController.GetRecipeByIdController, jwtMiddleware)
	e.POST("/api/recipes", recipeController.CreateRecipeController, jwtMiddleware)
	e.PUT("/api/recipes/:recipeId", recipeController.UpdateRecipeController, jwtMiddleware)
	e.DELETE("/api/recipes/:recipeId", recipeController.DeleteRecipeController, jwtMiddleware)
	// e.GET("/api/recipes/category/:categoryId", recipesCategroriesController.GetRecipeByCategoryIdController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Categories
	// ------------------------------------------------------------------
	e.GET("/api/categories", categoryController.GetAllCategoryController, jwtMiddleware)
	e.POST("/api/categories", categoryController.InsertCategoryController, jwtMiddleware)
	e.GET("/api/categories/:id", categoryController.GetCategoryController, jwtMiddleware)
	e.PUT("/api/categories/:id", categoryController.EditCategoryController, jwtMiddleware)
	e.DELETE("/api/categories/:id", categoryController.DeleteCategoryController, jwtMiddleware)

	// Recipe Categories
	e.POST("/api/recipe/categories", recipesCategroriesController.AddRecipeCategoriesController, jwtMiddleware)
	e.GET("/api/recipe/categories/:categoryId", recipesCategroriesController.GetRecipeByCategoryIdController, jwtMiddleware)

	// Ingredients
	e.GET("/api/ingredients", ingredientController.GetAllIngredientController, jwtMiddleware)
	e.GET("/api/ingredients/:ingredientId", ingredientController.GetIngredientByIdController, jwtMiddleware)
	e.POST("/api/ingredients", ingredientController.CreateIngredientController, jwtMiddleware)
	e.PUT("/api/ingredients/:ingredientId", ingredientController.UpdateIngredientController, jwtMiddleware)
	e.PUT("/api/ingredients/stock/:ingredientId", ingredientController.UpdateIngredientStockController, jwtMiddleware)
	e.DELETE("/api/ingredients/:ingredientId", ingredientController.DeleteIngredientController, jwtMiddleware)

	// Recipe Ingredients
	e.POST("/api/ingredients/recipe", recipeIngredientsController.AddIngredientsRecipeController, jwtMiddleware)
	e.GET("/api/ingredients/recipe/:recipeId", recipeIngredientsController.GetIngredientsByRecipeIdController, jwtMiddleware)

	// Cart Details
	e.GET("/api/carts", cartDetailsController.GetAllRecipeByCartIdController, jwtMiddleware)
	e.POST("/api/carts/recipe", cartDetailsController.AddRecipeToCartController, jwtMiddleware)
	// e.PUT("/api/carts/:recipeId", cartDetailsController.UpdateRecipePortionController, jwtMiddleware)
	// e.DELETE("/api/carts/:recipeId", cartDetailsController.DeleteRecipeFromCartController, jwtMiddleware)
}
