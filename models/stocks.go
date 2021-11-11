package models

import (
	"gorm.io/gorm"
)

type Stock struct {
	gorm.Model
	IngredientId uint
	Stock        int `json:"stock" form:"stock"`
	UserId       uint
}

type GormStockModel struct {
	db *gorm.DB
}

func NewStockModel(db *gorm.DB) *GormStockModel {
	return &GormStockModel{db: db}
}

type StockModel interface {
	CreateStockUpdate(stock Stock, ingredientId int) (Stock, error)
	Restock(stock Stock, ingredientId int, useId int) (Stock, error)
	StockDecrease(stock Stock, ingredientId int) (Stock, error)
	GetRestockDate(daterange string) (Stock, error)
}

func (m *GormStockModel) CreateStockUpdate(stock Stock, ingredientId int) (Stock, error) {
	if err := m.db.Save(&stock).Error; err != nil {
		return stock, err
	}
	return stock, nil
}

func (m *GormStockModel) Restock(stock Stock, ingredientId int, userId int) (Stock, error) {
	var newStock = stock.Stock

	if err := m.db.Raw("UPDATE stocks SET stock = stock + ? WHERE ingredient_id = ? AND user_id = ?", newStock, ingredientId, userId).Scan(&stock).Error; err != nil {
		return stock, err
	}
	return stock, nil
}

func (m *GormStockModel) StockDecrease(stock Stock, ingredientId int) (Stock, error) {
	var newStock = stock.Stock

	if err := m.db.Raw("UPDATE stocks SET stock = stock - ? WHERE ingredient_id = ?", newStock, ingredientId).Scan(&stock).Error; err != nil {
		return stock, err
	}
	return stock, nil
}

func (m *GormStockModel) GetRestockDate(daterange string) (Stock, error) {
	var stock Stock

	if daterange == "daily" {
		if err := m.db.Raw("SELECT * FROM stocks WHERE created_at >= DATE_ADD(CURDATE(), INTERVAL -1 DAY)").Scan(&stock).Error; err != nil {
			return stock, err
		}
		return stock, nil
	} else if daterange == "weekly" {
		if err := m.db.Where("created_at >= ?", "DATE_ADD(CURDATE(), INTERVAL -7 DAY)").Find(&stock).Error; err != nil {
			return stock, err
		}
		return stock, nil
	} else if daterange == "monthly" {
		if err := m.db.Where("created_at >= ?", "DATE_ADD(CURDATE(), INTERVAL -30 DAY)").Find(&stock).Error; err != nil {
			return stock, err
		}
		return stock, nil
	}

	return stock, nil
}
