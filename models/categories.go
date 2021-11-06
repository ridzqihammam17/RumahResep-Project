package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name" form:"name"`
	// One to Many with Recipe Category
	RecipeCategories []RecipeCategories
}

type GormCategoryModel struct {
	db *gorm.DB
}

func NewCategoryModel(db *gorm.DB) *GormCategoryModel {
	return &GormCategoryModel{db: db}
}

type CategoryModel interface {
	GetAll() ([]Category, error)
	Get(categoryId int) (Category, error)
	Insert(Category) (Category, error)
	Edit(category Category, categoryId int) (Category, error)
	Delete(categoryId int) (Category, error)
}

func (m *GormCategoryModel) GetAll() ([]Category, error) {
	var category []Category
	if err := m.db.Find(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (m *GormCategoryModel) Get(categoryId int) (Category, error) {
	var category Category
	if err := m.db.First(&category, categoryId).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (m *GormCategoryModel) Insert(category Category) (Category, error) {
	if err := m.db.Save(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (m *GormCategoryModel) Edit(newCategory Category, categoryId int) (Category, error) {
	var category Category
	if err := m.db.Find(&category, "id=?", categoryId).Error; err != nil {
		return category, err
	}

	category.Name = newCategory.Name

	if err := m.db.Save(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (m *GormCategoryModel) Delete(categoryId int) (Category, error) {
	var category Category
	if err := m.db.Find(&category, "id=?", categoryId).Error; err != nil {
		return category, err
	}
	if err := m.db.Delete(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}
