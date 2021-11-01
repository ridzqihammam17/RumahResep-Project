package carts

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CartController struct {
	cartModel models.CartModel
	// cartDetailModel models.CartDetailModel
	// recipeModel     models.RecipeModel
}

func NewCartController(cartModel models.CartModel) *CartController {
	// cartDetailModel models.CartDetailModel, recipeModel models.RecipeModel
	return &CartController{
		cartModel,
		// cartDetailModel,
		// recipeModel,
	}
}

func (controller *CartController) CreateCartController(c echo.Context) error {
	var cart models.Cart
	c.Bind(&cart)

	// get id user & role login
	userId, role := middlewares.ExtractTokenUser(c)

	// check role is customer
	if role != "customer" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Role not Customer",
		})
	}
	//check recipe id on table recipe
	// paymentId := cart.PaymentMethodsID

	//set data cart and create new cart
	cart = models.Cart{
		// StatusTransactions: "ordered",
		RecipeId:      0,
		TotalQuantity: 0,
		TotalPrice:    0,
		UserID:        userId,
		// Role:          role,
		// PaymentMethodsID:   paymentId,
	}

	// check if userId doesn't have cart yet

	newCart, _ := controller.cartModel.CreateCart(cart)

	//------------ cart detail -------------//
	// convert recipe id
	// fmt.Println(c.Param("recipeId"))
	// recipeId, err := strconv.Atoi(c.Param("recipeId"))

	// fmt.Println(recipeId, err)

	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]interface{}{
	// 		"message": "Recipe id is invalid",
	// 	})
	// }

	// check recipe id on table recipe

	// checkRecipeId, err := controller.recipeModel.CheckRecipeId(recipeId)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]interface{}{
	// 		"message":       "Recipe isn't found with recipe id #" + strconv.Itoa(recipeId),
	// 		"checkRecipeId": checkRecipeId,
	// 	})
	// }

	//get price
	// getRecipe, _ := controller.recipeModel.Get(recipeId)

	//convert qty
	// fmt.Println(c.Param("cnt"))
	// qty, _ := strconv.Atoi(c.Param("cnt"))
	// fmt.Println(qty)
	// recipePrice := getRecipe.Price
	// //set data cart details
	// cartDetails := models.CartDetails{
	// 	RecipesID: recipeId,
	// 	CartID:    newCart.ID,
	// 	Quantity:  qty,
	// 	Price:     recipePrice,
	// }

	//create cart detail
	// newCartDetail, _ := controller.cartDetailModel.AddToCart(cartDetails)

	//update total quantity and total price on table carts
	// controller.cartModel.UpdateTotalCart(newCart.ID, recipePrice, qty)

	//get cart updated (total qty&total price)
	// updatedCart, _ := controller.cartModel.GetCart(int(newCart.ID))

	// //custom data cart for body response
	// outputCart := map[string]interface{}{
	// 	// "ID":                  updatedCart.ID,
	// 	"customers_id": updatedCart.UserID,
	// 	// "payment_methods_id":  updatedCart.PaymentMethodsID,
	// 	// "status_transactions": updatedCart.StatusTransactions,
	// 	"total_quantity": updatedCart.TotalQuantity,
	// 	"total_price":    updatedCart.TotalPrice,
	// 	// "CreatedAt":           updatedCart.CreatedAt,
	// 	// "UpdatedAt":           updatedCart.UpdatedAt,
	// 	// "DeletedAt":           updatedCart.DeletedAt,
	// }

	return c.JSON(http.StatusOK, map[string]interface{}{
		"cart": newCart,
		// "cartDetails": newCartDetail,
		"status": "Create cart success",
	})

}

//func for update total quantity and total price on table carts
// func (controller *CartController) UpdateTotalCart(cartId int) (int, int) {
// 	newTotalPrice, _ := controller.cartModel.GetTotalPrice(cartId)
// 	newTotalQty, _ := controller.cartModel.GetTotalQty(cartId)
// 	newCart, _ := controller.cartModel.UpdateTotalCart(cartId, newTotalPrice, newTotalQty)

// 	return newCart.TotalQuantity, newCart.TotalPrice
// }

func (controller *CartController) GetCartController(c echo.Context) error {
	//convert cart_id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid id cart",
		})
	}
	//is cart id exist
	// var cart models.Cart
	checkCartId, err := controller.cartModel.CheckCartId(int(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":       "Cant find cart id",
			"checkRecipeId": checkCartId,
		})
	}

	listCart, _ := controller.cartModel.GetCartById(id) //get cart by id
	// recipes, _ := controller.cartDetailModel.GetListRecipeCart(id) //get all recipes based on cart id

	//custom data cart for body response
	outputCart := map[string]interface{}{
		"ID": listCart.ID,
		// "customers_id":        listCart.CustomersID,
		// "payment_methods_id":  listCart.PaymentMethodsID,
		// "status_transactions": listCart.StatusTransactions,
		"total_quantity": listCart.TotalQuantity,
		"total_price":    listCart.TotalPrice,
		"customer_id":    listCart.UserID,
		// "CreatedAt":           listCart.CreatedAt,
		// "UpdatedAt":           listCart.UpdatedAt,
		// "DeletedAt":           listCart.DeletedAt,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"cart": outputCart,
		// "recipes": recipes,
		"status": "Success get all recipes by cart id",
	})
}

// func (controller *CartController) DeleteCartController(c echo.Context) error {
// 	//convert cart id
// 	cartId, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]interface{}{
// 			"message": "Invalid cart id",
// 		})
// 	}

// 	//check is cart id exist on table cart
// 	//var cart models.Cart
// 	checkCartId, err := controller.cartModel.CheckCartId(cartId)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]interface{}{
// 			"message":     "Can't find cart",
// 			"checkCartId": checkCartId,
// 		})
// 	}

// 	//delete cart and all recipes included on it
// 	deletedCart, _ := controller.cartModel.DeleteCart(cartId)

// 	//custom output data cart for body response
// 	outputCart := map[string]interface{}{
// 		"ID": deletedCart.ID,
// 		// "customers_id":        deletedCart.CustomersID,
// 		// "payment_methods_id":  deletedCart.PaymentMethodsID,
// 		"status_transactions": deletedCart.StatusTransactions,
// 		"total_quantity":      deletedCart.TotalQuantity,
// 		"total_price":         deletedCart.TotalPrice,
// 		"CreatedAt":           deletedCart.CreatedAt,
// 		"UpdatedAt":           deletedCart.UpdatedAt,
// 		"DeletedAt":           deletedCart.DeletedAt,
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"status":       "Delete cart success",
// 		"Deleted Cart": outputCart,
// 	})
// }
