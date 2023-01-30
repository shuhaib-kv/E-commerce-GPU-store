package controllers

import (
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
	var sumtotal int
	database.Db.Raw("select sum(total) from carts where user_id=?", user.ID).Scan(&sumtotal)

	client := razorpay.NewClient("rzp_test_Nfnipdccvgb8fW", "UfwKXCGjiUrcfTEXpWlupcrN")
	razpayvalue := sumtotal * 100
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
		"total_price":  sumtotal,
		"total":        razpayvalue,
		"orderid":      value,
		"amount":       sumtotal,
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
	userid := c.Query("user_id")
	userID, _ := strconv.Atoi(userid)
	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	id := c.Query("orderid")
	totalamount := c.Query("total")
	Rpay := models.RazorPay{
		UserID:          userID,
		OrderId:         id,
		RazorPaymentId:  paymentid,
		Signature:       signature,
		RazorPayOrderID: orderid,
		AmountPaid:      totalamount,
	}
	database.Db.Create(&Rpay)
	var cart models.Cart
	database.Db.Raw("delete from carts where user_id=?", userID).Scan(&cart)
	OrderPlaced(userID, orderid)

	c.JSON(200, gin.H{

		"status": true,
	})

}
func Success(c *gin.Context) {
	c.HTML(200, "success.html", nil)

}
func OrderPlaced(Uid int, orderId string) {
	userid := Uid
	orderid := orderId

	var orders models.Orders
	database.Db.Raw("update orders set order_status=?,payment_status=?,order_type=?,orderid=? where users_id=?", "order completed", "payment done", "cart", orderid, userid).Scan(&orders)
	database.Db.Raw("update cart_infos set orderid=? where users_id=?", orderid, userid).Scan(&orders)

}
