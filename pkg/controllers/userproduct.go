package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"

	"github.com/gin-gonic/gin"
)

func ViewProductsUser(c *gin.Context) {
	var products []models.Product
	database.Db.Find(&products)

	for _, i := range products {
		c.JSON(200, gin.H{
			"id":             i.ID,
			"Name":           i.Name,
			"Actual price":   i.Price,
			"Discount_price": i.Discount_Price,
			"image":          i.Image1 + i.Image2 + i.Image3,
			"brand":          i.Brand,
			"chipset brand":  i.Chipset_brand,
			"modelgpu":       i.Model_gpu,
		})
	}

}

func ViewProductsUserbyid(c *gin.Context) {
	id := c.Param("id")
	var products []models.Product
	database.Db.Find(&products).Where("products.id=?", id).Scan(&products)

	for _, i := range products {
		c.JSON(200, gin.H{
			"id":                    i.ID,
			"Name":                  i.Name,
			"Actual price":          i.Price,
			"Discount_price":        i.Discount_Price,
			"image":                 i.Image1 + i.Image2 + i.Image3,
			"brand":                 i.Brand,
			"Description":           i.Description,
			"Chipset_brand":         i.Chipset_brand,
			"Model_gpu":             i.Model_gpu,
			"Series":                i.Series,
			"Generation":            i.Generation,
			"Memmory_type":          i.Memmory_type,
			"Thermal_design_power":  i.Thermal_design_power,
			"Released":              i.Released,
			"Architecture":          i.Architecture,
			"Memmory_size":          i.Memmory_size,
			"Recomented_resolution": i.Recomented_resolution,
			"DirectX":               i.DirectX,
			"Memmory_bus_width":     i.Memmory_bus_width,
			"Production_status":     i.Production_status,
			"Text_mapping_unit":     i.Text_mapping_unit,
			"Slots":                 i.Slots,
			"Rops":                  i.Rops,
			"Power_Connecters":      i.Power_Connecters,
		})
	}

}
