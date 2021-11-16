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
	StockDecrease(stock int, ingredientId int) (Stock, error)
	GetRestockDate(userId int, daterange string) (Stock, error)
	GetRestockAll(userId int) ([]Stock, error)
}

func (m *GormStockModel) CreateStockUpdate(stock Stock, ingredientId int) (Stock, error) {
	if err := m.db.Save(&stock).Error; err != nil {
		return stock, err
	}
	return stock, nil
}

func (m *GormStockModel) StockDecrease(stock int, ingredientId int) (Stock, error) {
	var newStock Stock

	if err := m.db.Raw("UPDATE stocks SET stock = stock - ? WHERE ingredient_id = ?", stock, ingredientId).Scan(&newStock).Error; err != nil {
		return newStock, err
	}
	return newStock, nil
}

func (m *GormStockModel) GetRestockDate(userId int, daterange string) (Stock, error) {
	var stock Stock

	if daterange == "daily" {
		if err := m.db.Where("created_at >= ? AND user_id = ?", "DATE_ADD(CURDATE(), INTERVAL -1 DAY)", userId).Find(&stock).Error; err != nil {
			return stock, err
		}
		return stock, nil
	} else if daterange == "weekly" {
		if err := m.db.Where("created_at >= ? AND user_id = ?", "DATE_ADD(CURDATE(), INTERVAL -7 DAY)", userId).Find(&stock).Error; err != nil {
			return stock, err
		}
		return stock, nil
	} else if daterange == "monthly" {
		if err := m.db.Where("created_at >= ? AND user_id = ?", "DATE_ADD(CURDATE(), INTERVAL -30 DAY)", userId).Find(&stock).Error; err != nil {
			return stock, err
		}
		return stock, nil
	}

	return stock, nil
}

func (m *GormStockModel) GetRestockAll(userId int) ([]Stock, error) {
	var stock []Stock

	if err := m.db.Where("user_id=?", userId).Find(&stock).Error; err != nil {
		return stock, err
	}

	return stock, nil
}
