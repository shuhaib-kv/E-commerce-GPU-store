package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddCategory(c *gin.Context) {
	Name_c := c.PostForm("name")

	catogory := models.Category{
		Name: Name_c,
	}
	var check []models.Category
	database.Db.Find(&check)
	for _, i := range check {
		if i.Name == catogory.Name {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "CATEGORY Already Exist",
			})
			return
		}

	}
	result := database.Db.Create(&catogory)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
	}
	// Return it
	c.JSON(200, gin.H{
		"message": "Category added",
	})
}

func ViewCategory(c *gin.Context) {
	var category []models.Category
	database.Db.Find(&category)
	for _, i := range category {
		c.JSON(200, gin.H{
			"id":   i.ID,
			"Name": i.Name,
		})
	}
}
func ViewProductByCategory(c *gin.Context) {
	id := c.Param("id")
	var product []models.Product
	database.Db.Find(&product, "categoryid = ?", id)
	for _, i := range product {
		c.JSON(200, gin.H{
			"id":    i.ID,
			"Name":  i.Name,
			"price": i.Price,
			"image": i.Image1 + i.Image2 + i.Image3,
			"brand": i.Brand,
		})
	}
}
func EditCategory(c *gin.Context) {
	id := c.Param("id")
	idc, _ := strconv.Atoi(id)
	Name := c.PostForm("name")
	var category []models.Category
	database.Db.First(&category, idc)
	database.Db.Model(&category).Updates(models.Category{
		Name: Name,
	})
	c.JSON(200, gin.H{
		"message": "category updated",
	})
}
func DeletECategory(c *gin.Context) {
	id := c.Param("id")
	idc, _ := strconv.Atoi(id)
	database.Db.Delete(&models.Category{}, idc)
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Deleted succesfully",
	})
}
