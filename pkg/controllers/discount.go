package controllers

import (
	"fmt"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddDiscount(c *gin.Context) {
	var discount models.Discount
	if err := c.BindJSON(&discount); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.Db.Create(&discount).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, discount)
}

func DeleteDiscount(c *gin.Context) {
	discountID := c.Param("id")
	var discount models.Discount
	if err := database.Db.First(&discount, discountID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Discount not found"})
		return
	}
	if err := database.Db.Delete(&discount).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Discount deleted"})
}
func ListDiscount(c *gin.Context) {
	var discount []models.Discount
	result := database.Db.Find(&discount)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "No discount found",
		})
		return
	}
	c.JSON(200, gin.H{
		"status": true,
		"data":   discount,
	})

}
func GetDiscountsWithFilters(c *gin.Context) {

	nameFilter := c.Query("name")
	idFilter := c.Query("id")
	percentageFilter := c.Query("percentage")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var totalCount int64

	var discounts []models.Discount

	query := database.Db.Model(&models.Discount{})
	if nameFilter != "" {
		query = query.Where("discount_name LIKE ?", fmt.Sprintf("%%%s%%", nameFilter))
	}
	if idFilter != "" {
		query = query.Where("id = ?", idFilter)
	}
	if percentageFilter != "" {
		query = query.Where("discount_percentage = ?", percentageFilter)
	}

	query.Count(&totalCount)
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&discounts).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{"count": totalCount,
			"page":      page,
			"pageSize":  pageSize,
			"discounts": discounts},
		"message": "Discounts",
	})
}
