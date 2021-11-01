package auth

import (
	"net/http"
	"rumah_resep/models"

	echo "github.com/labstack/echo/v4"
)

type AuthController struct {
	userModel models.UserModel
}

func NewAuthController(userModel models.UserModel) *AuthController {
	return &AuthController{
		userModel,
	}
}

func (controller *AuthController) RegisterUserController(c echo.Context) error {
	var userRequest models.User

	if err := c.Bind(&userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	user := models.User{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Gender:   userRequest.Gender,
		Password: userRequest.Password,
		Role:     userRequest.Role,
	}

	_, err := controller.userModel.Register(user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"code":    500,
			"message": "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Register Account",
	})
}

func (controller *AuthController) LoginUserController(c echo.Context) error {
	var userRequest models.User

	if err := c.Bind(&userRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	user, err := controller.userModel.Login(userRequest.Email, userRequest.Password)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": user.Token,
	})
}
