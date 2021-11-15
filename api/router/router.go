package router

import (
	"rumah_resep/api/controllers/auth"
	cartdetails "rumah_resep/api/controllers/cartDetails"
	"rumah_resep/api/controllers/carts"
	"rumah_resep/api/controllers/categories"
	"rumah_resep/api/controllers/checkouts"
	"rumah_resep/api/controllers/ingredients"
	"rumah_resep/api/controllers/midtrans"
	recipeingredients "rumah_resep/api/controllers/recipeIngredients"
	"rumah_resep/api/controllers/recipes"
	recipescategories "rumah_resep/api/controllers/recipesCategories"
	stock "rumah_resep/api/controllers/stocks"
	"rumah_resep/api/controllers/transactions"

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
	ingredientController *ingredients.IngredientController,
	recipeIngredientsController *recipeingredients.RecipeIngredientsController,
	cartDetailsController *cartdetails.CartDetailsController,
	checkoutController *checkouts.CheckoutController,
	transactionController *transactions.TransactionController,
	midtransController *midtrans.MidtransController,
	stockController *stock.StockController,
) {
	// ------------------------------------------------------------------
	// Auth Login & Register
	// ------------------------------------------------------------------
	e.POST("/api/register", authController.RegisterUserController)
	e.POST("/api/login", authController.LoginUserController)

	// ------------------------------------------------------------------
	// Auth JWT
	// ------------------------------------------------------------------
	jwtMiddleware := middleware.JWT([]byte(constants.SECRET_JWT))

	// ------------------------------------------------------------------
	// Recipe
	// ------------------------------------------------------------------
	e.GET("/api/recipes", recipeController.GetAllRecipeController, jwtMiddleware)
	e.GET("/api/recipes/:recipeId", recipeController.GetRecipeByIdController, jwtMiddleware)
	e.POST("/api/recipes", recipeController.CreateRecipeController, jwtMiddleware)
	e.PUT("/api/recipes/:recipeId", recipeController.UpdateRecipeController, jwtMiddleware)
	e.DELETE("/api/recipes/:recipeId", recipeController.DeleteRecipeController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Categories
	// ------------------------------------------------------------------
	e.GET("/api/categories", categoryController.GetAllCategoryController, jwtMiddleware)
	e.POST("/api/categories", categoryController.InsertCategoryController, jwtMiddleware)
	e.GET("/api/categories/:id", categoryController.GetCategoryController, jwtMiddleware)
	e.PUT("/api/categories/:id", categoryController.EditCategoryController, jwtMiddleware)
	e.DELETE("/api/categories/:id", categoryController.DeleteCategoryController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Recipe Categories
	// ------------------------------------------------------------------
	e.POST("/api/recipe/categories", recipesCategroriesController.AddRecipeCategoriesController, jwtMiddleware)
	e.GET("/api/recipe/categories/:categoryId", recipesCategroriesController.GetRecipeByCategoryIdController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Ingredients
	// ------------------------------------------------------------------
	e.GET("/api/ingredients", ingredientController.GetAllIngredientController, jwtMiddleware)
	e.GET("/api/ingredients/:ingredientId", ingredientController.GetIngredientByIdController, jwtMiddleware)
	e.POST("/api/ingredients", ingredientController.CreateIngredientController, jwtMiddleware)
	e.PUT("/api/ingredients/:ingredientId", ingredientController.UpdateIngredientController, jwtMiddleware)
	e.DELETE("/api/ingredients/:ingredientId", ingredientController.DeleteIngredientController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Stock
	// ------------------------------------------------------------------
	e.POST("/api/stocks/:ingredientId", stockController.CreateStockUpdateController, jwtMiddleware)
	e.PUT("/api/stocks/:ingredientId", stockController.UpdateStockController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Restock
	// ------------------------------------------------------------------
	e.GET("/api/restocks", stockController.GetRestockAllController, jwtMiddleware)
	e.GET("/api/restocks/:range", stockController.GetRestockDateController, jwtMiddleware)

	// Recipe Ingredients
	// ------------------------------------------------------------------
	e.POST("/api/ingredients/recipe", recipeIngredientsController.AddIngredientsRecipeController, jwtMiddleware)
	e.GET("/api/ingredients/recipe/:recipeId", recipeIngredientsController.GetIngredientsByRecipeIdController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Carts
	// ------------------------------------------------------------------
	e.POST("/api/carts", cartController.CreateCartController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Cart Details
	// ------------------------------------------------------------------
	e.GET("/api/cartdetails", cartDetailsController.GetAllRecipeByCartIdController, jwtMiddleware)
	e.POST("/api/cartdetails", cartDetailsController.AddRecipeToCartController, jwtMiddleware)
	e.PUT("/api/cartdetails/:recipeId", cartDetailsController.UpdateRecipePortionController, jwtMiddleware)
	e.DELETE("/api/cartdetails/:recipeId", cartDetailsController.DeleteRecipeFromCartController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Checkouts
	// ------------------------------------------------------------------
	e.POST("/api/checkouts/:recipeId", checkoutController.CreateCheckoutController, jwtMiddleware)

	// ------------------------------------------------------------------
	// Transactions
	// ------------------------------------------------------------------
	e.GET("/api/transactions", transactionController.GetAllTransaction, jwtMiddleware)
	e.POST("/api/transactions", transactionController.CreateTransaction, jwtMiddleware)

	// ------------------------------------------------------------------
	// Payments
	// ------------------------------------------------------------------
	e.GET("/api/payments/request/:id", midtransController.RequestPayment, jwtMiddleware)
	e.GET("/api/payments/status/:id", midtransController.StatusPayment, jwtMiddleware)
}
