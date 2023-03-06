package controllers

import (
	"fmt"
	"ga/pkg/database"
	"ga/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

func RazorPay(c *gin.Context) {
	var user models.Users
	useremail := c.GetString("user")
	database.Db.Raw("select id,phone from users where email=?", useremail).Scan(&user)
	var order models.Orders
	database.Db.Where("users_id= ? and payment_method=?", user.ID, "razorpay").First(&order)
	client := razorpay.NewClient("rzp_test_Nfnipdccvgb8fW", "UfwKXCGjiUrcfTEXpWlupcrN")
	razpayvalue := order.TotalAmount * 100
	data := map[string]interface{}{
		"amount":   razpayvalue,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	value := body["id"]
	if err != nil {
		c.JSON(404, gin.H{
			"msg": "error creating order",
		})
	}
	c.HTML(200, "app.html", gin.H{

		"UserID":       user.ID,
		"total_price":  order.TotalAmount,
		"total":        razpayvalue,
		"orderid":      value,
		"amount":       order.TotalAmount,
		"Email":        useremail,
		"Phone_Number": user.Phone,
	})
	if err != nil {
		c.JSON(200, gin.H{
			"msg": value,
		})
	}
}

func RazorpaySuccess(c *gin.Context) {

	var user models.Users
	useremail := c.GetString("user")
	database.Db.Raw("select id,phone from users where email=?", useremail).Scan(&user)

	userid := c.Query("user_id")
	userID, _ := strconv.Atoi(userid)
	orderid := c.Query("order_id")
	signature := c.Query("signature")
	var order models.Orders
	database.Db.Where("users_id= ?", user.ID).First(&order)
	Rpay := models.RazorPay{
		UserID:          uint(userID),
		RazorPaymentId:  orderid,
		Signature:       signature,
		RazorPayOrderID: order.Orderid,
		AmountPaid:      order.TotalAmount,
	}
	err := database.Db.Create(&Rpay)
	if err.Error != nil {
		fmt.Println("error")
	}
	c.JSON(200, gin.H{

		"status": true,
	})

}
func Success(c *gin.Context) {
	c.HTML(200, "success.html", nil)

}
