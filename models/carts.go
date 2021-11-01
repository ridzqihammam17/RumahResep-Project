package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	RecipeId      int
	TotalQuantity int `json:"total_quantity" form:"total_quantity"`
	TotalPrice    int `json:"total_price" form:"total_price"`
	// 1 to 1
	UserID int
	// Role   string
}

type GormCartModel struct {
	db *gorm.DB
}

func NewCartModel(db *gorm.DB) *GormCartModel {
	return &GormCartModel{db: db}
}

type CartModel interface {
	CreateCart(cart Cart) (Cart, error)
	GetCart(cartId int) (Cart, error)
	// GetTotalPrice(cartId int) (int, error)
	// GetTotalQty(cartId int) (int, error)
	// UpdateTotalCart(cartId int, newTotalPrice int, newTotalQty int) (Cart, error)
	CheckCartId(cartId int) (interface{}, error)
	GetCartById(id int) (Cart, error)
	// DeleteCart(cartId int) (cart Cart, err error)
}

func (m *GormCartModel) CreateCart(cart Cart) (Cart, error) {

	if err := m.db.Select(&cart, "NOT EXISTS (SELECT * FROM users WHERE user_id=?", cart.UserID).Save(&cart).Error; err != nil {
		return cart, err
	}
	return cart, nil
}

// get cart by id
func (m *GormCartModel) GetCart(cartId int) (Cart, error) {
	var cart Cart
	if err := m.db.Find(&cart, "id=?", cartId).Error; err != nil {
		return cart, err
	}
	return cart, nil
}

// get total price
func (m *GormCartModel) GetTotalPrice(cartId int) (int, error) {
	var totalPrice int
	if err := m.db.Select("sum(cart_details.price*cart_details.quantity)").Joins("JOIN carts ON carts.id = cart_details.carts_id").Where("carts_id=?", cartId).First(&totalPrice).Error; err == nil {
		return totalPrice, err
	}
	return totalPrice, nil
}

//get total quantity
func (m *GormCartModel) GetTotalQty(cartId int) (int, error) {
	var cartDetails CartDetails
	var totalQty int
	if err := m.db.Model(&cartDetails).Select("SUM(cart_details.quantity)").Joins("JOIN carts ON carts.id = cart_details.carts_id").Where("carts_id=?", cartId).First(&totalQty).Error; err == nil {
		return totalQty, err
	}
	return totalQty, nil
}

//update total cart
func (m *GormCartModel) UpdateTotalCart(cartId int, newTotalPrice int, newTotalQty int) (Cart, error) {
	var cart Cart

	if err := m.db.Find(&cart, "id=?", cartId).Error; err != nil {
		return cart, err
	}

	cart.TotalQuantity += newTotalQty
	cart.TotalPrice += (newTotalPrice * newTotalQty)

	if err := m.db.Save(&cart).Error; err != nil {
		return cart, err
	}
	return cart, nil
}

//check is cart id exist on table cart
func (m *GormCartModel) CheckCartId(cartId int) (interface{}, error) {
	var cart []Cart
	if err := m.db.Where("id=?", cartId).First(&cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

// get cart by id
func (m *GormCartModel) GetCartById(id int) (Cart, error) {
	var cart Cart
	if err := m.db.Find(&cart, "id=?", id).Error; err != nil {
		return cart, err
	}
	return cart, nil
}

//delete cart
func (m *GormCartModel) DeleteCart(cartId int) (cart Cart, err error) {

	if err := m.db.Find(&cart, "id=?", cartId).Unscoped().Delete(&cart).Error; err != nil {
		return cart, err
	}
	return cart, nil
}
