package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ViewProductsUser(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	var products []models.Product
	query := database.Db.Model(&models.Product{})
	if name := c.Query("name"); name != "" {
		query = query.Where("name iLIKE ?", "%"+name+"%")
	}
	if brand := c.Query("brand"); brand != "" {
		query = query.Where("brand iLIKE ?", "%"+brand+"%")
	}
	if minPrice, err := strconv.Atoi(c.Query("minPrice")); err == nil {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice, err := strconv.Atoi(c.Query("maxPrice")); err == nil {
		query = query.Where("price <= ?", maxPrice)
	}
	if id, err := strconv.Atoi(c.Query("id")); err == nil {
		query = query.Where("id = ?", id)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := make([]gin.H, 0, len(products))
	for _, product := range products {
		price := product.Price
		discountAmount := uint(0)
		if product.Discount > 0 {
			var discount models.Discount
			if err := database.Db.First(&discount, product.Discount).Error; err == nil {
				discountAmount = uint(float64(price) * float64(discount.DiscountPercentage) / 100.0)
				price = price - discountAmount
			}
		}
		response = append(response, gin.H{
			"id":              product.ID,
			"name":            product.Name,
			"price":           price,
			"image":           product.Image1 + product.Image2 + product.Image3,
			"brand":           product.Brand,
			"chipset_brand":   product.Chipset_brand,
			"model_gpu":       product.Model_gpu,
			"discount_amount": discountAmount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
		"data":     response,
	})
}
