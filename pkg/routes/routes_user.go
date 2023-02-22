package routes

import (
	"ga/middleware"
	"ga/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(User *gin.Engine) {
	//login and signup
	User.POST("/user/login/otp", controllers.OtpLog)
	User.POST("/user/signup", controllers.UserSignUp)
	User.POST("/user/login/otp/validate", controllers.CheckOtp)
	User.POST("/user/login", controllers.UserLogin)
	//userhome
	// User.GET("/user/home", middleware.UserAuth(), controllers.UserHome, controllers.PrintUserid)
	//products
	User.GET("/user/viewproducts", middleware.UserAuth(), controllers.ViewProductsUser)
	User.GET("/user/viewproducts/:id", middleware.UserAuth(), controllers.ViewProductsUserbyid)

	//category
	User.GET("/user/product/viewbycategory/:id", middleware.UserAuth(), controllers.ViewProductByCategory)
	//cart
	// User.GET("/cart", middleware.UserAuth(), controllers.ViewCart)
	// User.POST("/cart/add", middleware.UserAuth(), controllers.AddTOcart)

	// User.POST("/order/:id", middleware.UserAuth(), controllers.Order)
	User.POST("/user/add/address", middleware.UserAuth(), controllers.AddAddress)
	User.PATCH("/user/edit/address/:id", middleware.UserAuth(), controllers.EditAddress)

	// User.GET("/user/cart/view", middleware.UserAuth(), controllers.CartLists)
	User.GET("/payment-success", middleware.UserAuth(), controllers.RazorpaySuccess)
	User.GET("/success", controllers.Success)
	User.GET("/user/address", middleware.UserAuth(), controllers.ShowAddress)

	//ordernow
	// User.POST("/user/buynow", middleware.UserAuth(), controllers.OrderNow)
	User.GET("/user/orderview", middleware.UserAuth(), controllers.OrderedItems)
	// User.GET("/user/order/details/:id", middleware.UserAuth(), controllers.VieworderDetails)
	// // User.POST("/user/cart/order", middleware.UserAuth(), controllers.BuyFromCart)
	User.PATCH("/user/orde/cancel", middleware.UserAuth(), controllers.CancelOrder)
	User.POST("/cart/add", middleware.UserAuth(), controllers.AddToCart) //new
	User.GET("/razorpay", middleware.UserAuth(), controllers.RazorPay)
	// User.GET("/user/cart/view", middleware.UserAuth(), controllers.CartList)
	// User.POST("/chechout", middleware.UserAuth(), controllers.CheckoutCart) //new
	User.GET("/user/wallet/history", middleware.UserAuth(), controllers.WalletInfo)
	User.GET("/user/wallet/balance", middleware.UserAuth(), controllers.WalletBalance)

	// User.POST("/od1", middleware.UserAuth(), controllers.CartCheckoutDetails)

}
