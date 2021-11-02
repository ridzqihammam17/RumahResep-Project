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

// Create Migration Here
func DatabaseMigration(db *gorm.DB) {
	db.AutoMigrate(models.User{})
	db.AutoMigrate(models.Cart{})
	db.AutoMigrate(models.Recipe{})
	db.AutoMigrate(models.Category{})
}
