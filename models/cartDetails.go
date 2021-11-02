package models

import "gorm.io/gorm"

type CartDetail struct {
	gorm.Model
	CartID       int `gorm:"primaryKey" json:"carts_id" form:"carts_id"`
	IngredientID int `gorm:"primaryKey" json:"ingredients_id" form:"ingredients_id"` //gorm:"primaryKey"
	Quantity     int `json:"quantity" form:"quantity"`
	Price        int `json:"price" form:"price"`
}

type GormCartDetailModel struct {
	db *gorm.DB
}

func NewCartDetailModel(db *gorm.DB) *GormCartDetailModel {
	return &GormCartDetailModel{db: db}
}

type CartDetailModel interface {
	CheckProductAndCartId(productId, cartId int, cartDetails CartDetail) (interface{}, error)
	GetCartDetailByCartId(cartId int) (CartDetail, error)
	AddToCart(cartDetails CartDetail) (CartDetail, error)
	DeleteIngredientFromCart(ingredientId int) (interface{}, error)
	GetListProductCart(cartId int) (interface{}, error)
	CountProductOnCart(cartId int) (int, error)
	CountProductandPriceOnCart(cartId int) (int, int, error)
}

func (m *GormCartDetailModel) CheckProductAndCartId(productId, cartId int, cartDetails CartDetail) (interface{}, error) {
	if err := m.db.Where("carts_id=? AND products_id=?", cartId, productId).First(&cartDetails).Error; err != nil {
		return nil, err
	}
	return cartDetails, nil
}

// get product by id
// func (m *GormCartDetailModel) GetProduct(productId int) (Product, error) {
// 	var product Product
// 	if err := m.db.Find(&product, "id=?", productId).Error; err != nil {
// 		return product, err
// 	}
// 	return product, nil
// }

//Get cart details by Cart ID
func (m *GormCartDetailModel) GetCartDetailByCartId(cartId int) (CartDetail, error) {
	var cartDetails CartDetail
	if err := m.db.Find(&cartDetails, cartId).Error; err != nil {
		return cartDetails, err
	}
	return cartDetails, nil
}

//add product to cart
func (m *GormCartDetailModel) AddToCart(cartDetail CartDetail) (CartDetail, error) {
	if err := m.db.Save(&cartDetail).Error; err != nil {
		return cartDetail, err
	}
	return cartDetail, nil
}

//delete product from cart detail
func (m *GormCartDetailModel) DeleteIngredientFromCart(ingredientId int) (interface{}, error) {
	var cartDetails CartDetail
	if err := m.db.Find(&cartDetails, "ingredient_id=?", ingredientId).Unscoped().Delete(&cartDetails).Error; err != nil {
		return nil, err
	}
	// if err := m.db.Find(&cartDetails, "cart_id=? AND products_id=?", cartId, productId).Unscoped().Delete(&cartDetails).Error; err != nil {
	// 	return nil, err
	// }
	return cartDetails, nil
}

//get all products from cart detail
func (m *GormCartDetailModel) GetListProductCart(cartId int) (interface{}, error) {
	var cartDetail []CartDetail

	if err := m.db.Find(&cartDetail, "cart_id=?", cartId).Error; err != nil {
		return nil, err
	}
	// if err := m.db.Table("products").Joins("JOIN cart_details ON products.id = cart_details.products_id").Joins("JOIN carts ON cart_details.carts_id = carts.id").Where("carts.id=?", cartId).Find(&cartDetail).Error; err != nil {
	// 	return cartDetail, nil
	// }
	return cartDetail, nil
}

func (m *GormCartDetailModel) CountProductOnCart(cartId int) (int, error) {
	var countProduct int
	if err := m.db.Select("COUNT(carts_id)").Where("carts_id=?", cartId).First(&countProduct).Error; err == nil {
		return countProduct, err
	}
	return countProduct, nil
}

func (m *GormCartDetailModel) CountProductandPriceOnCart(cartId int) (int, int, error) {
	var countProduct, Price int
	if err := m.db.Select("COUNT(carts_id), Price").Where("carts_id=?", cartId).First(&countProduct).Error; err == nil {
		return countProduct, Price, err
	}
	return countProduct, Price, nil
}
