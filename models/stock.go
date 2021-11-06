package models

import "gorm.io/gorm"

type Stock struct {
	gorm.Model
	Stock        int  `json:"stock" form:"stock"`
	IngredientId uint `json:"ingredient_id" form:"ingredient_id"`
	UserId       uint `json:"user_id" form:"user_id"`
}

type GormStockModel struct {
	db *gorm.DB
}

func NewStockModel(db *gorm.DB) *GormStockModel {
	return &GormStockModel{db: db}
}

type StockModel interface {
	CreateStockUpdate(stock Stock) (Stock, error)
}

func (m *GormStockModel) CreateStockUpdate(stock Stock) (Stock, error) {
	if err := m.db.Save(&stock).Error; err != nil {
		return stock, err
	}
	return stock, nil
}
