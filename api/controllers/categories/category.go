package categories

import (
	"net/http"
	"rumah_resep/api/middlewares"
	"rumah_resep/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ------------------------------------------------------------------
// Start Request
// ------------------------------------------------------------------

type InsertCategoryRequest struct {
	Name string `json:"name" form:"name"`
}

type EditCategoryRequest struct {
	Name string `json:"name" form:"name"`
}

// ------------------------------------------------------------------
// End Request
// ------------------------------------------------------------------

type CategoryController struct {
	categoryModel models.CategoryModel
}

func NewCategoryController(categoryModel models.CategoryModel) *CategoryController {
	return &CategoryController{
		categoryModel,
	}
}

// ------------------------------------------------------------------
// Admin Authorize Check
// ------------------------------------------------------------------
func AuthorizeAdmin(c echo.Context) bool {
	_, role := middlewares.ExtractTokenUser(c)
	return role == "admin"
}

func (controller *CategoryController) GetAllCategoryController(c echo.Context) error {
	checkAuthorize := AuthorizeAdmin(c)
	if !checkAuthorize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	category, err := controller.categoryModel.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Get All Category",
		"data":    category,
	})
}

func (controller *CategoryController) GetCategoryController(c echo.Context) error {
	checkAuthorize := AuthorizeAdmin(c)
	if !checkAuthorize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	category, err := controller.categoryModel.Get(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Get Category",
		"data":    category,
	})
}

func (controller *CategoryController) InsertCategoryController(c echo.Context) error {
	checkAuthorize := AuthorizeAdmin(c)
	if !checkAuthorize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	var categoryRequest InsertCategoryRequest

	if err := c.Bind(&categoryRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	category := models.Category{
		Name: categoryRequest.Name,
	}

	_, err := controller.categoryModel.Insert(category)

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
		"message": "Success Insert Category",
	})
}

func (controller *CategoryController) EditCategoryController(c echo.Context) error {
	checkAuthorize := AuthorizeAdmin(c)
	if !checkAuthorize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	var categoryRequest EditCategoryRequest
	if err := c.Bind(&categoryRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	category := models.Category{
		Name: categoryRequest.Name,
	}

	if _, err := controller.categoryModel.Edit(category, id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Edit Category",
	})
}

func (controller *CategoryController) DeleteCategoryController(c echo.Context) error {
	checkAuthorize := AuthorizeAdmin(c)
	if !checkAuthorize {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"code":    401,
			"message": "Unauthorized Error",
		})
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	if _, err := controller.categoryModel.Delete(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"code":    400,
			"message": "Bad Request",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"code":    200,
		"message": "Success Delete Category",
	})
}
