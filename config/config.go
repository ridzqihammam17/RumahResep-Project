package config

import (
	"os"
	"strconv"
	"sync"
)

//AppConfig Application configuration
type AppConfig struct {
	Port     int
	Database struct {
		Driver     string
		Connection string
	}
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}

	return appConfig
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func initConfig() *AppConfig {
	var defaultConfig AppConfig

	httpPort, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		return &defaultConfig
	}

	defaultConfig.Port = httpPort
	defaultConfig.Database.Driver = "mysql"
	defaultConfig.Database.Connection = getEnv("CONNECTION_STRING", "root:@tcp(localhost:3306)/db-rumah-resep?charset=utf8&parseTime=True&loc=Local")

	return &defaultConfig
}
