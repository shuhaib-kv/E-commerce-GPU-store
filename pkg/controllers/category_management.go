package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddCategory(c *gin.Context) {
	var reqBody struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err,
		})
		return
	}
	var count int64
	database.Db.Model(&models.Category{}).Where("name = ?", reqBody.Name).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Category already exists",
		})
		return
	}
	category := models.Category{Name: reqBody.Name}
	if err := database.Db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to create category",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Category added",
		"data":    category,
	})
}

func ViewCategory(c *gin.Context) {
	var categories []models.Category
	result := database.Db.Find(&categories)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to retrieve categories",
			"error":   result.Error,
		})
		return
	}

	var categoriesJSON []gin.H
	for _, category := range categories {
		categoryJSON := gin.H{
			"id":   category.ID,
			"name": category.Name,
		}
		categoriesJSON = append(categoriesJSON, categoryJSON)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "fetched all categories",
		"data":    categoriesJSON,
	})
}
func ViewProductByCategory(c *gin.Context) {
	var reqBody struct {
		Id int `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err,
		})
		return
	}
	var products []models.Product
	database.Db.Find(&products, "id = ?", reqBody.Id)

	var productJSON []gin.H
	for _, p := range products {
		productJSON = append(productJSON, gin.H{
			"id":    p.ID,
			"name":  p.Name,
			"price": p.Price,
			"image": p.Image1 + p.Image2 + p.Image3,
			"brand": p.Brand,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"message":    "fetched all products on given category",
		"categories": productJSON,
	})
}

func EditCategory(c *gin.Context) {
	var reqBody struct {
		Id   int    `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err,
		})
		return
	}
	var category models.Category
	result := database.Db.First(&category, reqBody.Id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Category not found",
			"error":   result.Error,
		})
		return
	}
	category.Name = reqBody.Name
	result = database.Db.Save(&category)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed to update category",
			"error":   result.Error,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Category updated",
		"data":    category,
	})
}

func DeletECategory(c *gin.Context) {
	var reqBody struct {
		ID int `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request body",
			"error":   err,
		})
		return
	}

	database.Db.Delete(&models.Category{}, reqBody.ID)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Category deleted successfully",
	})
}
