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
	GetRestockDate(daterange string) ([]Stock, error)
}

func (m *GormStockModel) CreateStockUpdate(stock Stock) (Stock, error) {
	if err := m.db.Save(&stock).Error; err != nil {
		return stock, err
	}
	return stock, nil
}

func (m *GormStockModel) GetRestockDate(daterange string) ([]Stock, error) {
	var stock []Stock

	if daterange == "daily" {
		if err := m.db.Where("created_at >= ?", "DATE_ADD(CURDATE(), INTERVAL -1 DAY)").Find(&stock).Error; err != nil {
			return stock, err
		}
	} else if daterange == "weekly" {
		if err := m.db.Where("created_at >= ?", "DATE_ADD(CURDATE(), INTERVAL -7 DAY)").Find(&stock).Error; err != nil {
			return stock, err
		}
	} else if daterange == "monthly" {
		if err := m.db.Where("created_at >= ?", "DATE_ADD(CURDATE(), INTERVAL -30 DAY)").Find(&stock).Error; err != nil {
			return stock, err
		}
	}

	return stock, nil
}
