package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
	"math"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AdminAddProduct(c *gin.Context) {

	P_Name := c.PostForm("name")
	sprice := c.PostForm("price")
	P_Price, _ := strconv.Atoi(sprice)
	mno := c.PostForm("modelno")
	Modelno, _ := strconv.Atoi(mno)
	Image1, _ := c.FormFile("image1")

	extension := filepath.Ext(Image1.Filename)
	img1 := uuid.New().String() + extension
	c.SaveUploadedFile(Image1, "./public/images"+img1)

	Image2, _ := c.FormFile("image2")
	extension = filepath.Ext(Image2.Filename)
	img2 := uuid.New().String() + extension
	c.SaveUploadedFile(Image2, "./public/images"+img2)
	Image3, _ := c.FormFile("image3")
	extension = filepath.Ext(Image3.Filename)
	img3 := uuid.New().String() + extension
	c.SaveUploadedFile(Image3, "./public/images"+img3)
	sstock := c.PostForm("stock")
	P_Stock, _ := strconv.Atoi(sstock)
	scategory := c.PostForm("CategoryID")
	P_CategoryID, _ := strconv.Atoi(scategory)

	P_description := c.PostForm("discription")
	P_Brand := c.PostForm("brand")
	P_Chipset_brand := c.PostForm("chipset_brand")
	P_Model_gpu := c.PostForm("model_gpu")
	P_Series := c.PostForm("series")
	P_Generation := c.PostForm("generation")
	P_Memmory_type := c.PostForm("memmorytype")
	P_Thermal_design_power := c.PostForm("tdp")
	P_Released := c.PostForm("released_year")
	P_Architecture := c.PostForm("architecture")
	sms := c.PostForm("memmorysize")
	P_Memmory_size, _ := strconv.Atoi(sms)
	P_Recomented_resolution := c.PostForm("recomentedresolution")
	P_DirectX := c.PostForm("directx")
	P_Memmory_bus_width := c.PostForm("memmorybuswidth")
	P_Production_status := c.PostForm("productionstatus")
	P_Text_mapping_unit := c.PostForm("tpu")
	P_Slots := c.PostForm("slot")
	P_Rops := c.PostForm("rops")
	P_Power_Connecters := c.PostForm("powerconnecters")
	p_Discount := c.PostForm("discount")
	discount, _ := strconv.Atoi(p_Discount)

	product := models.Product{
		Name:                  P_Name,
		Price:                 uint(P_Price),
		ModelNo:               uint(Modelno),
		Image1:                img1,
		Image2:                img2,
		Image3:                img3,
		Stock:                 uint(P_Stock),
		CategoryID:            uint(P_CategoryID),
		Description:           P_description,
		Brand:                 P_Brand,
		Chipset_brand:         P_Chipset_brand,
		Model_gpu:             P_Model_gpu,
		Series:                P_Series,
		Generation:            P_Generation,
		Memmory_type:          P_Memmory_type,
		Thermal_design_power:  P_Thermal_design_power,
		Released:              P_Released,
		Architecture:          P_Architecture,
		Memmory_size:          uint(P_Memmory_size),
		Recomented_resolution: P_Recomented_resolution,
		DirectX:               P_DirectX,
		Memmory_bus_width:     P_Memmory_bus_width,
		Production_status:     P_Production_status,
		Text_mapping_unit:     P_Text_mapping_unit,
		Slots:                 P_Slots,
		Rops:                  P_Rops,
		Power_Connecters:      P_Power_Connecters,
		Discount:              uint(discount),
	}
	var check []models.Product
	database.Db.Find(&check)
	for _, i := range check {
		if i.ModelNo == product.ModelNo {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Product Exist",
			})
			return
		}
	}
	createproduct := database.Db.Create(&product)
	if createproduct.Error != nil {
		c.JSON(400, gin.H{
			"message": "Error",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": "Product added",
	})
}

func EditProduct(c *gin.Context) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	var body struct {
		Name                  string
		Price                 uint
		ModelNo               uint
		Image1                string
		Image2                string
		Image3                string
		Stock                 uint
		CategoryID            uint
		SubCategoryID         uint
		Description           string
		Brand                 string
		Chipset_brand         string
		Model_gpu             string
		Series                string
		Generation            string
		Memmory_type          string
		Thermal_design_power  string
		Released              string
		Architecture          string
		Memmory_size          uint
		Recomented_resolution string
		DirectX               string
		Memmory_bus_width     string
		Production_status     string
		Text_mapping_unit     string
		Slots                 string
		Rops                  string
		Power_Connecters      string
		Discount              uint
		Discount_Price        uint
	}
	c.Bind(&body)
	var products []models.Product
	var count uint
	database.Db.Raw("select count(id) from products where id=?", ID).Scan(&count)
	if count <= 0 {
		c.JSON(404, gin.H{
			"msg": "product doesnot exist",
		})
		c.Abort()
		return
	} else {
		database.Db.First(&products, id)
		database.Db.Model(&products).Updates(models.Product{
			Name:                  body.Name,
			Price:                 body.Price,
			ModelNo:               body.ModelNo,
			Stock:                 body.Stock,
			CategoryID:            body.CategoryID,
			Description:           body.Description,
			Brand:                 body.Brand,
			Chipset_brand:         body.Chipset_brand,
			Model_gpu:             body.Model_gpu,
			Series:                body.Series,
			Generation:            body.Memmory_type,
			Memmory_type:          body.Memmory_type,
			Thermal_design_power:  body.Thermal_design_power,
			Released:              body.Released,
			Architecture:          body.Architecture,
			Memmory_size:          body.Memmory_size,
			Recomented_resolution: body.Recomented_resolution,
			DirectX:               body.DirectX,
			Memmory_bus_width:     body.Memmory_bus_width,
			Production_status:     body.Production_status,
			Text_mapping_unit:     body.Text_mapping_unit,
			Slots:                 body.Slots,
			Rops:                  body.Rops,
			Power_Connecters:      body.Power_Connecters,
			Discount:              body.Discount,
		})
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Product updated",
		})

	}

}

type ProductFilter struct {
	Name      string `form:"name"`
	Brand     string `form:"brand"`
	MinPrice  uint   `form:"minPrice"`
	MaxPrice  uint   `form:"maxPrice"`
	Category  string `form:"category"`
	PageSize  int    `form:"pageSize"`
	PageIndex int    `form:"pageIndex"`
}

func ViewProducts(c *gin.Context) {
	// Parse the filter parameters from the request query string
	var filter ProductFilter
	if err := c.BindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid filter parameters",
			"error":   err.Error(),
		})
		return
	}

	// Set default values for the filter parameters if they are not provided
	if filter.PageSize == 0 {
		filter.PageSize = 10
	}
	if filter.PageIndex == 0 {
		filter.PageIndex = 1
	}

	// Build the query to retrieve the products based on the filter parameters
	query := database.Db.Model(&models.Product{})
	if filter.Name != "" {
		query = query.Where("name iLIKE ?", "%"+filter.Name+"%")
	}
	if filter.Brand != "" {
		query = query.Where("brand = ?", filter.Brand)
	}
	if filter.MinPrice != 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice != 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}
	if filter.Category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.name = ?", filter.Category)
	}

	// Paginate the results based on the filter parameters
	var totalProducts int64
	query.Count(&totalProducts)
	totalPages := int(math.Ceil(float64(totalProducts) / float64(filter.PageSize)))
	offset := (filter.PageIndex - 1) * filter.PageSize
	var products []models.Product
	query.Offset(offset).Limit(filter.PageSize).Find(&products)

	// Build the response payload
	var productJSON []gin.H
	for _, product := range products {
		productJSON = append(productJSON, gin.H{
			"id":             product.ID,
			"name":           product.Name,
			"price":          product.Price,
			"image":          product.Image1 + product.Image2 + product.Image3,
			"brand":          product.Brand,
			"discount_price": product.Discount,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      true,
		"message":     "Products found",
		"data":        productJSON,
		"currentPage": filter.PageIndex,
		"pageSize":    filter.PageSize,
		"totalPages":  totalPages,
		"totalItems":  totalProducts,
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var products models.Product
	if err := database.Db.First(&products, "id=?", id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Cant find product",
			"error":   "check the id",
		})
		return
	}
	if err := database.Db.Delete(&products).Where("id=?", id); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": " cant delete product",
			"error":   "database error",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"status":  true,
		"message": "  deleted product",
		"data":    products.Name,
	})
}
