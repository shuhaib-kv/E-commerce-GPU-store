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

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "discount added successfully",
		"data":    discount,
	})
}
func DeleteDiscount(c *gin.Context) {
	var request struct {
		Id uint `json:"id"`
	}
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"error":   err,
			"message": "Invalid request body",
		})
		return
	}

	var discount models.Discount
	if err := database.Db.First(&discount, request.Id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Discount not found",
			"error":   err.Error(),
		})
		return
	}

	if err := database.Db.Delete(&discount).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"error":   err.Error(),
			"message": "Failed to delete discount",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Discount deleted successfully",
		"data":    discount,
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
