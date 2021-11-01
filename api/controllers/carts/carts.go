package carts

import (
	"fmt"
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
func (controller *CartController) UpdateCartController(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// get id user & role login
	userId, role := middlewares.ExtractTokenUser(c)
	fmt.Println(userId, role)
	// check role is customer
	if role != "customer" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Role not Customer",
		})
	}
	// listCart, _ := controller.cartModel.GetCartById(userId) //get cart by id
	// recipes, _ := controller.cartDetailModel.GetListRecipeCart(id) //get all recipes based on cart id

	// custom data cart for body response
	var newCart models.Cart
	c.Bind(&newCart)
	outputCart := models.Cart{
		// "ID": listCart.ID,
		// "customers_id":        listCart.CustomersID,
		// "payment_methods_id":  listCart.PaymentMethodsID,
		// "status_transactions": listCart.StatusTransactions,
		TotalQuantity: newCart.TotalQuantity,
		TotalPrice:    newCart.TotalPrice,
		// "customer_id":    listCart.UserID,
		// "CreatedAt":           listCart.CreatedAt,
		// "UpdatedAt":           listCart.UpdatedAt,
		// "DeletedAt":           listCart.DeletedAt,
	}
	// newTotalPrice, _ := controller.cartModel.GetTotalPrice(userId)
	// newTotalQty, _ := controller.cartModel.GetTotalQty(userId)
	// fmt.Println(newTotalPrice, newTotalQty)
	newCartq, _ := controller.cartModel.UpdateTotalCart(outputCart, id)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"cart": newCartq,
		// "cartDetails": newCartDetail,
		"status": "Update cart success",
	})
}

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

func (controller *CartController) DeleteCartController(c echo.Context) error {
	//convert cart id
	id, _ := strconv.Atoi(c.Param("id"))

	// get id user & role login
	userId, role := middlewares.ExtractTokenUser(c)
	fmt.Println(userId, role)
	// check role is customer
	if role != "customer" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Role not Customer",
		})
	}
	// listCart, _ := controller.cartModel.GetCartById(userId) //get cart by id
	// recipes, _ := controller.cartDetailModel.GetListRecipeCart(id) //get all recipes based on cart id

	// custom data cart for body response
	var newCart models.Cart
	c.Bind(&newCart)
	outputCart := models.Cart{
		// "ID": listCart.ID,
		// "customers_id":        listCart.CustomersID,
		// "payment_methods_id":  listCart.PaymentMethodsID,
		// "status_transactions": listCart.StatusTransactions,
		TotalQuantity: 0,
		TotalPrice:    0,
		// "customer_id":    listCart.UserID,
		// "CreatedAt":           listCart.CreatedAt,
		// "UpdatedAt":           listCart.UpdatedAt,
		// "DeletedAt":           listCart.DeletedAt,
	}
	// newTotalPrice, _ := controller.cartModel.GetTotalPrice(userId)
	// newTotalQty, _ := controller.cartModel.GetTotalQty(userId)
	// fmt.Println(newTotalPrice, newTotalQty)
	newCartq, _ := controller.cartModel.UpdateTotalCart(outputCart, id)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"cart": newCartq,
		// "cartDetails": newCartDetail,
		"status": "Delete cart success",
	})
}
