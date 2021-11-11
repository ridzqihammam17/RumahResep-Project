package config

import (
	"os"
	"strconv"
	"sync"
	"time"
)

//AppConfig Application configuration
type AppConfig struct {
	Port     int
	Database struct {
		Driver     string
		Connection string
		ConnTest   string
	}
}

type ThirdPartyConfig struct {
	GoogleMapsAPIKey        string
	GoogleMapsAPIUrl        string
	GoogleMapsGeoCodeAPIUrl string
}

type HTTPServerConfig struct {
	Addr            string
	ShutdownTimeout time.Duration
	WriteTimeout    time.Duration
	ReadTimeout     time.Duration
	IdleTimeout     time.Duration
}

//HTTPServer httpServer config
var HTTPServer HTTPServerConfig

var lock = &sync.Mutex{}
var appConfig *AppConfig

// var ThirdParty ThirdPartyConfig

// -- GeoCoding and MapsConfig
var ThirdParty ThirdPartyConfig

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = InitConfig()
	}

	return appConfig
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InitConfig() *AppConfig {
	var defaultConfig AppConfig

	httpPort, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		return &defaultConfig
	}

	defaultConfig.Port = httpPort
	defaultConfig.Database.Driver = "mysql"
	defaultConfig.Database.Connection = getEnv("CONNECTION_STRING", "root:@tcp(localhost:3306)/dbrumahresep?charset=utf8&parseTime=True&loc=Local")
	defaultConfig.Database.ConnTest = getEnv("CONNECTION_STRING_TEST", "root:@tcp(localhost:3306)/dbrumahreseptest?charset=utf8&parseTime=True&loc=Local")

	ThirdParty = ThirdPartyConfig{
		GoogleMapsAPIKey:        getEnv("GOOGLE_MAPS_API_KEY", "AIzaSyAfF0h3oFhZS23os2XgPF8OIxTxKtkD8qI"),
		GoogleMapsAPIUrl:        getEnv("GOOGLE_MAPS_API_URL", "https://maps.googleapis.com/maps/api/distancematrix/json?units=metric&origins=%s,%s&destinations=%s,%s&key=%s"),
		GoogleMapsGeoCodeAPIUrl: getEnv("GOOGLE_MAPS_GEOCODE_API_URL", "https://maps.googleapis.com/maps/api/geocode/json?"),
	}

	return &defaultConfig
}
