package controllers

// var letters = []rune("123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

// func OrderIdGeneration(value uint) string {
// 	b := make([]rune, value)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(b)
// }

// func OrderInfo(c *gin.Context) {

// 	var Orders []models.Orders
// 	database.Db.Find(&Orders)
// 	for _, i := range Orders {
// 		c.JSON(200, gin.H{
// 			"status":         true,
// 			"UserId":         i.UsersID,
// 			"OrderID":        i.OrderID,
// 			"Discount":       i.Discount,
// 			"CouponDiscount": i.CouponDiscount,
// 			"CouponCode":     i.CouponCode,
// 			"PaymentMethod":  i.Payment_Method,
// 			"TotalAmount":    i.Total_Amount,
// 		})
// 	}

// }

// func OrderedItems(c *gin.Context) {
// 	// Fetching user id from jwt
// 	useremail := c.GetString("user")
// 	var UsersID uint
// 	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
// 	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "user coudnt find",
// 		})

// 	}
// 	var items []models.Orders
// 	database.Db.Where("users_id = ?", UsersID).Find(&items)

// 	for _, i := range items {
// 		c.JSON(200, gin.H{
// 			"status":      true,
// 			"id":          i.OrderID,
// 			"Amount_Paid": i.Total_Amount,

// 			"Discount":        i.Discount,
// 			"Coupon_Discount": i.CouponDiscount,
// 			"Order Status":    i.OrderStatus,
// 		})
// 	}
// }

// func CancelOrder(c *gin.Context) {
// 	useremail := c.GetString("user")
// 	var UsersID uint
// 	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
// 	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "user coudnt find",
// 		})
// 		var items models.Orders
// 		var updateStatus string = "CANCELLED"
// 		id := c.Query("orderid")

// 		database.Db.First(&items, id)
// 		if items.OrderStatus == updateStatus {
// 			c.JSON(400, gin.H{
// 				"status":  false,
// 				"message": "Order already Cancelled",
// 			})
// 			return
// 		}
// 		database.Db.Model(&items).Where("id=?", id).Update("order_status", updateStatus)

// 		var price uint
// 		database.Db.Raw("SELECT price FROM ordereditems WHERE id = ?", id).Scan(&price)

// 		var balance uint
// 		database.Db.Raw("SELECT balance FROM wallets WHERE users_id = ?", UsersID).Scan(&balance)
// 		newBalance := balance + price

// 		if items.Payment_Method == "COD" {
// 			c.JSON(200, gin.H{
// 				"status":  true,
// 				"message": "Order Cancelled",
// 			})
// 			return
// 		}

// 		WalletHistory := models.Wallethistory{UsersID: uint(UsersID), Debit: 0, Credit: price}
// 		database.Db.Create(&WalletHistory)

// 		var totalAmount uint
// 		database.Db.Raw("SELECT total_amaount FROM orders WHERE users_id = ?", UsersID).Scan(&totalAmount)
// 		Ntotal := totalAmount - balance
// 		// Updating wallet on order cancellation
// 		database.Db.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
// 		database.Db.Model(&models.Orders{}).Where("users_id = ?", UsersID).Update("total_amount", Ntotal)
// 		c.JSON(200, gin.H{
// 			"status":  true,
// 			"message": "Order Cancelled",
// 		})
// 	}
// }
// func CreatesOrderId() string {
// 	rand.Seed(time.Now().UnixNano())
// 	value := rand.Intn(9999999999-1000000000) + 1000000000
// 	id := strconv.Itoa(value)
// 	orderID := "OID" + id
// 	return orderID
// }

// func OrderNow(c *gin.Context) {
// 	useremail := c.GetString("user")
// 	fmt.Println(useremail)
// 	var UsersID uint
// 	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
// 	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "user coudnt find",
// 		})

// 	}
// 	notcompRazorpay := "Needs to complete razorpay payment"
// 	var body struct {
// 		ProductsID     uint   `json:"productsid"`
// 		AddressID      uint   `json:"addressid"`
// 		Payment_Method string `json:"paymentmethod"`
// 	}

// 	c.Bind(&body)
// 	var price uint
// 	database.Db.Raw("select price from products where id=?", body.ProductsID).Scan(&price)
// 	if body.Payment_Method == "COD" {
// 		order := models.Orders{
// 			OrderID:        CreateOrderId(),
// 			UsersID:        UsersID,
// 			AddressID:      body.AddressID,
// 			Payment_Method: body.Payment_Method,
// 			Total_Amount:   price,
// 			PaymentStatus:  "Pending",
// 			OrderStatus:    "order placed",
// 			Ordertype:      "BuyNow",
// 		}
// 		orders := database.Db.Create(&order)
// 		if orders.Error != nil {
// 			c.JSON(400, gin.H{
// 				"message": "Error",
// 			})
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"status":  true,
// 			"message": " orderplaced",
// 		})

// 	}
// 	if body.Payment_Method == "RAZORPAY" {
// 		orders := models.Orders{
// 			UsersID:        UsersID,
// 			AddressID:      body.AddressID,
// 			OrderID:        CreateOrderId(),
// 			OrderStatus:    "pending",
// 			Payment_Method: body.Payment_Method,
// 			PaymentStatus:  notcompRazorpay,
// 			Total_Amount:   price,
// 		}

// 		result := database.Db.Create(&orders)
// 		if result.Error == nil {
// 			c.JSON(300, gin.H{
// 				"msg": "Go to the Razorpay Page for Order completion",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		var ordereditems models.Orderd_Items

// 		database.Db.Raw("update orderd_items set  order_status=?,payment_status=?,payment_method=? where user_id=?", "orderplaced", notcompRazorpay, razorpay, user.ID).Scan(&ordereditems)
// 		if result.Error == nil {
// 			c.JSON(300, gin.H{
// 				"msg": "Go to the Razorpay Page for Order completion",
// 			})
// 			c.Abort()
// 			return
// 		}

// 	}
// }

// func ViewOrdersUser(c *gin.Context) {
// 	useremail := c.GetString("user")
// 	fmt.Println(useremail)
// 	var UsersID uint
// 	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
// 	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "user coudnt find",
// 		})

// 	}
// 	var order []models.Orders
// 	database.Db.Find(&order)
// 	for _, i := range order {
// 		c.JSON(200, gin.H{
// 			"id":              i.ID,
// 			"user id":         i.UsersID,
// 			"price":           i.Total_Amount,
// 			"Adressid":        i.AddressID,
// 			"order status":    i.OrderStatus,
// 			"payment methode": i.Payment_Method,
// 			"payment status":  i.PaymentStatus,
// 			"order type":      i.OrderType,
// 		})
// 	}

// }
// func VieworderDetails(c *gin.Context) {
// 	id := c.Param("id")
// 	orderid, _ := strconv.Atoi(id)
// 	var address []struct {
// 		Address_id uint
// 		UserId     uint
// 		Pincode    uint
// 		House      string
// 		Area       string
// 		Landmark   string
// 		City       string
// 	}
// 	var order []struct {
// 		ID             uint
// 		UsersID        uint
// 		AddressID      uint
// 		Payment_Method string
// 		OrderDate      string
// 		Total_Amount   uint
// 		PaymentStatus  string
// 		OrderStatus    string
// 		CreatedAt      time.Time
// 	}
// 	var ad string

// 	database.Db.Select("address_id").Table("orders").Where("id=?", orderid).Scan(&ad)

// 	database.Db.Select("users_id", "address_id", "payment_method", "total_amount", "order_status", "created_at").Table("orders").Where("address_id=?", ad).Find(&order)
// 	c.JSON(200, gin.H{
// 		"Products": order,
// 	})
// 	database.Db.Select("name", "phone_number", "pincode", "house", "area", "landmark", "city").Table("addresses").Where("address_id=?", ad).Find(&address)
// 	c.JSON(200, gin.H{
// 		"Address": address,
// 	})
// }
// func Cancelorders(c *gin.Context) {
// 	useremail := c.GetString("user")
// 	fmt.Println(useremail)
// 	var UsersID uint
// 	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
// 	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "user coudnt find",
// 		})

// 	}
// 	var user models.Users
// 	var ordered_items models.Orders
// 	userEmail := c.GetString("user")
// 	orderid := c.Param("id")
// 	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
// 	record := database.Db.Raw("select users_id,productid,total_amount,orderid,order_status,payment_status,payment_method,total_amount from orders where users_id =?", UsersID).Scan(&ordered_items)
// 	if record.Error != nil {
// 		c.JSON(404, gin.H{"err": record.Error.Error()})
// 		c.Abort()
// 		return
// 	}
// 	ni := "order canceld"
// 	fmt.Println(orderid)
// 	database.Db.Raw("update orders set order_status=? where id=?", ni, orderid).Scan(&ordered_items)
// 	c.JSON(200, gin.H{"orders": ordered_items})
// }
