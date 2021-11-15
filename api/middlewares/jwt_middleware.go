package middlewares

import (
	"rumah_resep/constants"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateToken(userId int, role, city string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = int(userId)
	claims["role"] = role
	claims["city"] = city
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(constants.SECRET_JWT))
}

func LoggerMiddlewares(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` + "\n",
	}))
}

func ExtractTokenUser(c echo.Context) (uint, string) {

	// fmt.Println(c.Get("user"))
	// if temp := c.Get("user"); temp != nil {

	// 	u := temp.(*jwt.Token)
	// 	claims := u.Claims.(jwt.MapClaims)
	// 	userId := int(claims["userId"].(float64))
	// 	role := claims["role"].(string)
	// 	return userId, role
	// }
	// return 0, ""

	token := c.Get("user").(*jwt.Token)
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		userId := uint(claims["userId"].(float64))
		role := claims["role"].(string)
		return userId, role
	}
	return 0, ""
}

func EctractCity(c echo.Context) string {

	token := c.Get("user").(*jwt.Token)
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		city := claims["city"].(string)
		return city
	}
	return ""
}
