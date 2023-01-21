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
	User.GET("/user/home", middleware.UserAuth(), controllers.UserHome)
	//products
	User.GET("/user/viewproducts", middleware.UserAuth(), controllers.ViewProductsUser)
	User.GET("/user/viewproducts/:id", middleware.UserAuth(), controllers.ViewProductsUserbyid)
	User.GET("/user/product/viewbycategory/:id", middleware.UserAuth(), controllers.ViewProductByCategory)
	//user profile
	User.POST("/user/add/address", middleware.UserAuth(), controllers.AddAddress)
	User.PATCH("/user/edit/address/:id", middleware.UserAuth(), controllers.EditAddress)
	User.GET("/user/address", middleware.UserAuth(), controllers.ShowAddress)
	// order
	User.GET("/user/orderview", middleware.UserAuth(), controllers.OrderedItems)
	User.PATCH("/user/orde/cancel", middleware.UserAuth(), controllers.CancelOrder)
	//payment
	User.GET("/payment-success", middleware.UserAuth(), controllers.RazorpaySuccess)
	User.GET("/razorpay", middleware.UserAuth(), controllers.RazorPay)
	User.GET("/success", controllers.Success)
	//ordernow
	User.POST("/cart/add", middleware.UserAuth(), controllers.AddTOcart)
	User.GET("/user/cart/view", middleware.UserAuth(), controllers.CartList)

	User.POST("/chechout", middleware.UserAuth(), controllers.CheckoutCart)
	//wallet
	User.GET("/user/wallet/history", middleware.UserAuth(), controllers.WalletInfo)
	User.GET("/user/wallet/balance", middleware.UserAuth(), controllers.WalletBalance)

}
