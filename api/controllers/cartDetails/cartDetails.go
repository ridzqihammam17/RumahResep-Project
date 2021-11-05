package cartdetails

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"

	"github.com/labstack/echo/v4"
)

type CartDetailsController struct {
	cartDetailsModel       models.CartDetailsModel
	recipeModel            models.RecipeModel
	ingredientsModel       models.IngredientModel
	recipeIngredientsModel models.RecipeIngredientsModel
	cartModel              models.CartModel
}

func NewCartDetailsController(cartDetailsModel models.CartDetailsModel, recipeModel models.RecipeModel, ingredientModel models.IngredientModel, recipeIngredientsModel models.RecipeIngredientsModel, cartModel models.CartModel) *CartDetailsController {
	return &CartDetailsController{
		cartDetailsModel,
		recipeModel,
		ingredientModel,
		recipeIngredientsModel,
		cartModel,
	}
}

func (controller *CartDetailsController) GetAllRecipeByCartIdController(c echo.Context) error {
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	cartDetails, err := controller.cartDetailsModel.GetAllRecipeByCartId(int(userId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Get All Recipe In Cart",
		"data":    cartDetails,
	})
}

func (controller *CartDetailsController) AddRecipeToCartController(c echo.Context) error {
	//check role admin or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	// Get cart id
	cartId, _ := controller.cartModel.GetCartIdByUserId(int(userId))

	// Bind body request
	var cartDetails models.CartDetails
	if err := c.Bind(&cartDetails); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Success Bad Request",
		})
	}

	// Get Total Price
	var idQtyingredient []models.RecipeIngredients
	var totalPrice int
	idQtyingredient, _ = controller.recipeIngredientsModel.GetIdIngredientQtyIngredient(cartDetails.RecipeID)

	for i := 0; i < len(idQtyingredient); i++ {
		price, _ := controller.ingredientsModel.GetIngredientPrice(idQtyingredient[i].IngredientId)
		totalPrice += idQtyingredient[i].QtyIngredient * price
	}

	// Prepare Post Body
	newCartDetails := models.CartDetails{
		CartID:   cartId,
		RecipeID: cartDetails.RecipeID,
		Quantity: cartDetails.Quantity,
		Price:    totalPrice,
	}

	output, err := controller.cartDetailsModel.AddRecipeToCart(newCartDetails)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Add Recipe To Cart",
		"data":    output,
	})
}

func (controller *CartDetailsController) UpdateRecipePortionController(c echo.Context) error {
	//check role customer or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	var cartDetails models.CartDetails
	if err := c.Bind(&cartDetails); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Success Bad Request",
		})
	}

	newCartDetails := models.CartDetails{
		RecipeID: cartDetails.RecipeID,
		Quantity: cartDetails.Quantity,
	}

	output, err := controller.cartDetailsModel.UpdateRecipePortion(newCartDetails, int(userId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Update Recipe Portion From Cart",
		"data":    output,
	})
}

func (controller *CartDetailsController) DeleteRecipeFromCartController(c echo.Context) error {
	//check role customer or not
	userId, role := middlewares.ExtractTokenUser(c)
	if role != "customer" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized",
		})
	}

	var cartDetails models.CartDetails
	if err := c.Bind(&cartDetails); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Success Bad Request",
		})
	}

	output, err := controller.cartDetailsModel.DeleteRecipeFromCart(int(userId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Update Recipe Portion From Cart",
		"data":    output,
	})
}
