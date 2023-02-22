package controllers

// func AddDiscount(c *gin.Context) {
// 	va body struct{
// 		Name
// 		Persentagec
// 		Product
// 	}
// 	discountName := c.PostForm("discountName")
// 	DPercentage := c.PostForm("discountPercentage")
// 	discountPercentage, _ := strconv.Atoi(DPercentage)
// 	PId := c.PostForm("productId")
// 	productId, _ := strconv.Atoi(PId)
// 	discount := models.Discount{DiscountName: discountName, DiscountPercentage: discountPercentage, ProductId: productId}
// 	var checkDisc []models.Discount
// 	database.Db.Find(&checkDisc)
// 	for _, i := range checkDisc {
// 		if i.DiscountName == discountName {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"status":  false,
// 				"message": "Discount Name Already Exist",
// 			})
// 			return
// 		}
// 	}
// 	result := database.Db.Create(&discount)
// 	if result.Error != nil {
// 		c.JSON(400, gin.H{
// 			"message": "Error Creating Coupon",
// 		})
// 		return
// 	}

// 	c.JSON(200, gin.H{
// 		"status":  true,
// 		"message": "Discount Created",
// 	})
// }
// func DeleteDiscount(c *gin.Context) {
// 	var discount models.Discount
// 	discountName := c.Query("discountName")
// 	database.Db.Where("coupon_name = ?", discountName).Delete(&discount)
// 	//database.DB.Raw("DELETE FROM coupons WHERE coupon_name=?", couponName).Scan(&coupon)
// 	c.JSON(200, gin.H{
// 		"status":  true,
// 		"message": "Deleted succesfully",
// 	})
// }
// func ListDiscount(c *gin.Context) {
// 	var discount []models.Discount
// 	result := database.Db.Find(&discount)
// 	if result.Error != nil {
// 		c.JSON(400, gin.H{
// 			"status":  false,
// 			"message": "No discount found",
// 		})
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"status": true,
// 		"data":   discount,
// 	})

// }
