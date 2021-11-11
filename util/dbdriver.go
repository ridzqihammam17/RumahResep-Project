package util

import (
	"rumah_resep/config"
	"rumah_resep/models"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlDatabaseConnection(config *config.AppConfig) *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Database.Connection), &gorm.Config{})

	if err != nil {
		log.Info("failed to connect database: ", err)
		panic(err)
	}

	// Uncommand For Migration
	DatabaseMigration(db)

	return db
}

func MysqlDatabaseConnTest(config *config.AppConfig) *gorm.DB {
	db_test, err := gorm.Open(mysql.Open(config.Database.ConnTest), &gorm.Config{})
	if err != nil {
		log.Info("failed to connect database: ", err)
		panic(err)
	}
	// DBMigrationTest(db_test)
	return db_test
}

// Create Migration Here
func DatabaseMigration(db *gorm.DB) {
	db.AutoMigrate(models.User{})
	db.AutoMigrate(models.Recipe{})
	db.AutoMigrate(models.Category{})
	db.AutoMigrate(models.Ingredient{})
	db.AutoMigrate(models.Stock{})
	db.AutoMigrate(models.Cart{})
	db.AutoMigrate(models.CartDetails{})
	db.AutoMigrate(models.RecipeCategories{})
	db.AutoMigrate(models.RecipeIngredients{})
	db.AutoMigrate(models.Checkout{})
	db.AutoMigrate(models.Transaction{})
}

// func DBMigrationTest(db *gorm.DB) {
// 	db.AutoMigrate(models.User{})
// 	db.AutoMigrate(models.Recipe{})
// 	db.AutoMigrate(models.Category{})
// 	db.AutoMigrate(models.Ingredient{})
// 	db.AutoMigrate(models.Stock{})
// 	db.AutoMigrate(models.Cart{})
// 	db.AutoMigrate(models.CartDetails{})
// 	db.AutoMigrate(models.RecipeCategories{})
// 	db.AutoMigrate(models.RecipeIngredients{})
// 	db.AutoMigrate(models.Checkout{})
// 	db.AutoMigrate(models.Transaction{})
// }
