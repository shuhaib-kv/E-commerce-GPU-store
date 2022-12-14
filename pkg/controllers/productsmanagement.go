package controllers

import (
	"ga/pkg/database"
	"ga/pkg/models"
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
	p_Discount_price := c.PostForm("discountprice")
	discountprice, _ := strconv.Atoi(p_Discount_price)
	product := models.Product{
		Name:                  P_Name,
		Price:                 P_Price,
		ModelNo:               Modelno,
		Image1:                img1,
		Image2:                img2,
		Image3:                img3,
		Stock:                 P_Stock,
		CategoryID:            P_CategoryID,
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
		Memmory_size:          P_Memmory_size,
		Recomented_resolution: P_Recomented_resolution,
		DirectX:               P_DirectX,
		Memmory_bus_width:     P_Memmory_bus_width,
		Production_status:     P_Production_status,
		Text_mapping_unit:     P_Text_mapping_unit,
		Slots:                 P_Slots,
		Rops:                  P_Rops,
		Power_Connecters:      P_Power_Connecters,
		Discount:              discount,
		Discount_Price:        discountprice,
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
		Price                 int
		ModelNo               int
		Image1                string
		Image2                string
		Image3                string
		Stock                 int
		CategoryID            int
		SubCategoryID         int
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
		Memmory_size          int
		Recomented_resolution string
		DirectX               string
		Memmory_bus_width     string
		Production_status     string
		Text_mapping_unit     string
		Slots                 string
		Rops                  string
		Power_Connecters      string
		Discount              int
		Discount_Price        int
	}
	c.Bind(&body)
	var products []models.Product
	var count int
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
			Discount_Price:        body.Discount_Price,
		})
		c.JSON(200, gin.H{
			"message": "Product updated",
		})

	}

}
func ViewProducts(c *gin.Context) {
	var products []models.Product
	database.Db.Find(&products)
	for _, i := range products {
		c.JSON(200, gin.H{
			"id":             i.ID,
			"Name":           i.Name,
			"price":          i.Price,
			"Discount_price": i.Discount,
			"image":          i.Image1 + i.Image2 + i.Image3,
			"brand":          i.Brand,
		})
	}

}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	DE := database.Db.Raw("select*from products where id =?", id)
	if DE != nil {
		c.JSON(200, gin.H{
			"status":  false,
			"message": "can't find product",
		})
		return
	}
	var count int
	// record := database.Db.Delete(&models.Product{}, id)
	database.Db.Raw("select count(product_id) from products where product_id=?", id).Scan(&count)
	if count <= 0 {
		c.JSON(404, gin.H{
			"msg": "product doesnot exist",
		})
		c.Abort()
		return
	} else {
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Deleted succesfully",
		})
	}

}
