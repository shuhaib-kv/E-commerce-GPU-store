package controllers

import (
	"errors"
	"fmt"
	"ga/pkg/database"
	"ga/pkg/models"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrderId() string {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999999999-1000000000) + 1000000000
	id := strconv.Itoa(value)
	orderID := "OID" + id
	return orderID
}

func AddTOcart(c *gin.Context) {
	useremail := c.GetString("user")
	fmt.Println(useremail)
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var body struct {
		ProductsID int `json:"productsid"`
		Quantity   int `json:"quantity"`
	}
	c.Bind(&body)
	fmt.Println(body.ProductsID)
	fmt.Println(body.Quantity)
	var stock int
	database.Db.Raw("SELECT stock FROM products WHERE id = ?", body.ProductsID).Scan(&stock)

	if stock < body.Quantity {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Product is Out Of Stock",
		})
		return
	}
	var price int
	var name string
	var brand string
	var discription string
	var discountprice int
	var total int
	database.Db.Raw("SELECT price FROM products WHERE id = ?", body.ProductsID).Scan(&price)
	database.Db.Raw("SELECT name FROM products WHERE id = ?", body.ProductsID).Scan(&name)
	database.Db.Raw("SELECT brand FROM products WHERE id = ?", body.ProductsID).Scan(&brand)
	database.Db.Raw("SELECT description FROM products WHERE id = ?", body.ProductsID).Scan(&discription)
	database.Db.Raw("SELECT discount_price FROM products WHERE id = ?", body.ProductsID).Scan(&discountprice)

	cart := models.Cart{
		User_id:       UsersID,
		Product_ID:    body.ProductsID,
		Quantity:      body.Quantity,
		Product_Name:  name,
		Brand_Name:    brand,
		Description:   discription,
		Price:         body.Quantity * price,
		DiscountPrice: body.Quantity * discountprice,
		Total:         total,
	}
	if discountprice < price {
		fmt.Println("wtf")
		var check []models.Cart
		database.Db.Find(&check)
		for _, i := range check {
			if i.Product_ID != cart.Product_ID && i.User_id != cart.Product_ID && i.Product_ID != body.ProductsID {
				continue
			} else {
				v := cart.Quantity + i.Quantity
				s := i.Price + body.Quantity*price
				p := i.DiscountPrice + body.Quantity*discountprice
				fmt.Println(v)
				fmt.Println(s)
				database.Db.Model(&cart).Where("user_id = ?", UsersID).Update("quantity", v)
				database.Db.Model(&cart).Where("user_id = ?", UsersID).Update("price", s)
				database.Db.Model(&cart).Where("user_id = ?", UsersID).Update("discount_price", p)
				database.Db.Model(&cart).Where("user_id = ?", UsersID).Update("total", p)
				fmt.Println()

				c.JSON(200, gin.H{
					"status":  true,
					"message": "Product already exist in cart so updated the quantity",
				})
				return
			}
		}
		result := database.Db.Create(&cart)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"message": "Error",
			})
			return
		}
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Added To Cart",
		})

	} else {
		fmt.Println("boooooo")
		var check []models.Cart
		database.Db.Find(&check)
		for _, i := range check {
			if i.Product_ID != cart.Product_ID && i.User_id != cart.Product_ID && i.Product_ID != body.ProductsID {
				continue
			} else {
				v := cart.Quantity + i.Quantity
				s := i.Price + body.Quantity*price
				p := i.DiscountPrice + body.Quantity*discountprice
				fmt.Println(v)
				fmt.Println(s)
				database.Db.Model(&cart).Where("user_id = ?", UsersID).Update("quantity", v)
				database.Db.Model(&cart).Where("user_id = ?", UsersID).Update("price", s)
				database.Db.Model(&cart).Where("user_id = ?", UsersID).Update("discount_price", p)
				database.Db.Model(&cart).Where("user_id = ?", UsersID).Update("total", s)
				database.Db.Raw("update carts where user_id=? set total=?", UsersID, s)

				c.JSON(200, gin.H{
					"status":  true,
					"message": "Product already exist in cart so updated the quantity",
				})
				return
			}
		}
		result := database.Db.Create(&cart)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"message": "Error",
			})
			return
		}
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Added To Cart",
		})
	}

	// cartInfo := models.CartInfo{UsersId: UsersID, ProductsID: body.ProductsID, Price: cart.Price, Quantity: cart.Quantity}

	for i := 1; i < body.Quantity; i++ {
		var discount models.Discount
		database.Db.Where("product_id = ?", body.ProductsID).Find(&discount)

		deductdiscount := discount.DiscountPercentage
		fmt.Print(deductdiscount)

		cartInfo := models.CartInfo{UsersId: UsersID, ProductsID: body.ProductsID, Discount: deductdiscount, ProductName: name, BrandName: brand, ProductPrice: price}
		database.Db.Create(&cartInfo)
	}
}

func CartList(c *gin.Context) {
	var Subtotal int
	useremail := c.GetString("user")
	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var cart []struct {
		ID           int
		Product_id   int
		Product_Name string
		Brand_Name   string
		Description  string
		Total        int
		Quantity     int
	}
	database.Db.Select("id", "product_id", "product_name", "brand_name", "description", "price", "quantity", "total").Table("carts").Where("user_id=?", UsersID).Find(&cart)

	c.JSON(200, gin.H{
		"Products": cart,
	})

	for _, i := range cart {
		sum := i.Total
		Subtotal = Subtotal + sum
	}
	c.JSON(200, gin.H{
		"status":   true,
		"Subtotal": Subtotal,
	})

}

type Cartsinfo []struct {
	User_id      int
	Product_id   int
	Product_Name string
	Price        string
	Email        string
	Quantity     int
	Total_Amount int
	Total_Price  int
}

func CheckoutCart(c *gin.Context) {
	// paymentmethod := c.Query("paymentmethod")
	Adressid := c.Query("adressid")
	addressID, _ := strconv.Atoi(Adressid)
	PaymentMethod := c.Query("paymentmethod")
	CoupenCode := c.Query("coupen")
	wallet := c.Query("Applaywallet")
	cod := "COD"
	razorpay := "RAZORPAY"
	notcompRazorpay := "Needs to complete razorpay payment"
	useremail := c.GetString("user")
	// record := database.Db.Delete(&models.Product{}, id)
	var count int
	database.Db.Raw("select count(coupon_code) from coupons where coupon_code=?", CoupenCode).Scan(&count)
	if count <= 0 {
		c.JSON(404, gin.H{
			"msg": "coupen doesnot exist",
		})
		c.Abort()
		return
	}

	var UsersID int
	err := database.Db.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var cart []struct {
		ID           int
		Product_id   int
		Product_Name string
		Brand_Name   string
		Description  string
		Price        int
		Quantity     int
	}
	var total int

	database.Db.Select("id", "product_id", "product_name", "price", "quantity").Table("carts").Where("user_id=?", UsersID).Find(&cart)

	// cod := "COD"
	// razorpay := "Razorpay"
	// notcompRazorpay := "Needs to complete razorpay payment"

	database.Db.Raw("select sum(total) as total from carts where user_id=?", UsersID).Scan(&total)
	fmt.Println(total)
	if CoupenCode != "" {
		var coupon models.Coupon
		database.Db.Select("coupon_code").Table("coupons").Where("coupon_code=?", CoupenCode).Find(&coupon)

	}

	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999999999-1000000000) + 1000000000
	id := strconv.Itoa(value)
	orderID := "OID" + id

	var balance int
	database.Db.Raw("SELECT balance FROM wallets WHERE users_id = ?", UsersID).Scan(&balance)
	c.JSON(200, gin.H{
		"Wallet Balance": balance,
	})
	if balance == 0 {
		c.JSON(200, gin.H{
			"Balance is zero": "",
		})
	} else {
		if wallet == "apply" {
			if balance > total {
				newBalance := balance - total
				database.Db.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
				WalletHistory := models.Wallethistory{UsersID: uint(UsersID), Debit: total, Credit: 0}
				database.Db.Create(&WalletHistory)
				total = 0
			} else if balance < total {
				total = total - balance
				newBalance := 0
				database.Db.Model(&models.Wallet{}).Where("users_id = ?", UsersID).Update("balance", newBalance)
				WalletHistory := models.Wallethistory{UsersID: uint(UsersID), Debit: balance, Credit: 0}
				database.Db.Create(&WalletHistory)
			}
		}
	}
	if CoupenCode == "" {
		var checkInfo models.Checkoutinfo
		database.Db.Raw("UPDATE checkoutinfos SET order_id = ?,discount = ?,coupon_discount = ?,coupon_code = ?,total_mrp = ?,total = ? WHERE users_id = ?", orderID, 0, 0, "Not Applied", total, total, UsersID).Scan(&checkInfo)
		c.JSON(200, gin.H{
			"messsage": "No COupon",
		})
		return
	} else {
		var coupon models.Coupon
		database.Db.Raw("SELECT coupon_percentage FROM coupons WHERE coupon_code = ?", CoupenCode).Find(&coupon)
		STotalAmpount := total * coupon.CouponPercentage / 100
		total = total - STotalAmpount
		c.JSON(200, gin.H{
			"Total MRP": total,
			// "Discount":        discountAmount,
			"Coupon Discount": STotalAmpount,
			// "Total":           TotalAmpount,
			"Wallet amount": balance,
		})
	}

	if PaymentMethod == razorpay {

		orders := models.Orders{
			UsersID:        UsersID,
			AddressID:      addressID,
			OrderID:        orderID,
			OrderStatus:    "pending",
			Ordertype:      "cart",
			Payment_Method: razorpay,
			PaymentStatus:  notcompRazorpay,
			Total_Amount:   total,
		}
		fmt.Println(total)
		result := database.Db.Create(&orders)
		if result.Error != nil {
			c.JSON(404, gin.H{"err": result.Error.Error()})
			c.Abort()
			return
		}
		var ordereditems models.Ordereditems
		database.Db.Raw("update ordereditems set  order_status=?,payment_status=?,payment_method=?,totalamount=? ,where users_id=?", "orderplaced", notcompRazorpay, razorpay, UsersID, total).Scan(&ordereditems)
		if result.Error == nil {
			c.JSON(300, gin.H{
				"msg": "Go to the Razorpay Page for Order completion",
			})
			c.Abort()
			return
		}
	} else if PaymentMethod == cod {
		fmt.Println("hai cod")
		orders := models.Orders{
			UsersID:        UsersID,
			AddressID:      addressID,
			OrderID:        orderID,
			OrderStatus:    "pending",
			Ordertype:      "cart",
			Payment_Method: cod,
			PaymentStatus:  "cash on delivery",
			Total_Amount:   total,
		}
		database.Db.Create(&orders)
		var ordereditems models.Ordereditems
		database.Db.Raw("update ordereditems set order_status=?,payment_method=? where users_id=?", "orderplaced", cod, UsersID).Scan(&ordereditems)
		c.JSON(300, gin.H{
			"msg": "orderplaced",
		})
		database.Db.Create(&ordereditems)
	} else {
		c.JSON(300, gin.H{
			"msg": "select payment method and address",
		})
		c.Abort()
		return
	}

	database.Db.Raw("delete from carts where user_id=?", UsersID).Scan(&cart)
}
