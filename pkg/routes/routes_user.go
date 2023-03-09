package routes

import (
	"ga/middleware"
	"ga/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(User *gin.Engine) {
	User.POST("/user/login/otp", controllers.OtpLog)
	User.POST("/user/signup", controllers.UserSignUp)
	User.POST("/user/login/otp/validate", controllers.CheckOtp)
	User.POST("/user/login", controllers.UserLogin)

	User.GET("/user/home", middleware.UserAuth(), controllers.UserHome)

	User.GET("/user/viewproducts", middleware.UserAuth(), controllers.ViewProductsUser)

	User.POST("/user/add/address", middleware.UserAuth(), controllers.AddAddress)
	User.PATCH("/user/edit/address", middleware.UserAuth(), controllers.EditAddress)
	User.GET("/user/view/address", middleware.UserAuth(), controllers.EditAddress)

	User.GET("/payment-success", middleware.UserAuth(), controllers.RazorpaySuccess)
	User.GET("/success", middleware.UserAuth(), controllers.Success)
	User.GET("/user/address", middleware.UserAuth(), controllers.ShowAddress)
	User.POST("/cart/order", middleware.UserAuth(), controllers.OrderCart)
	User.GET("/user/orderview", middleware.UserAuth(), controllers.ListOrders)
	User.POST("/user/cancel/order", middleware.UserAuth(), controllers.CancelOrder)
	User.POST("/cart/add", middleware.UserAuth(), controllers.AddToCart)
	User.GET("/razorpay", middleware.UserAuth(), controllers.RazorPay)
	User.GET("/user/cart/view", middleware.UserAuth(), controllers.ViewCart)
	User.GET("/user/wallet/history", middleware.UserAuth(), controllers.WalletInfo)
}
