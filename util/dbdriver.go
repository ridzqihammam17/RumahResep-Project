package util

import (
	"fmt"
	"rumah_resep/config"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MysqlDatabaseConnection(config *config.AppConfig) *gorm.DB {
	uri := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		config.Database.Username,
		config.Database.Password,
		config.Database.Address,
		config.Database.Port,
		config.Database.Name)

	db, err := gorm.Open(mysql.Open(uri), &gorm.Config{})

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

}
