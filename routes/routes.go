package routes

import (
	"ga/handler"
	"ga/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	r *gin.Engine,
	jwtSecret string,
	userH *handler.UserHandler,
	adminH *handler.AdminHandler,
	productH *handler.ProductHandler,
	categoryH *handler.CategoryHandler,
	cartH *handler.CartHandler,
	orderH *handler.OrderHandler,
	couponH *handler.CouponHandler,
	discountH *handler.DiscountHandler,
	walletH *handler.WalletHandler,
	paymentH *handler.PaymentHandler,
) {
	adminAuth := middleware.AdminAuth(jwtSecret)
	userAuth := middleware.UserAuth(jwtSecret)

	// User auth (no middleware)
	r.POST("/user/signup", userH.Signup)
	r.POST("/user/login", userH.Login)

	// User protected routes
	user := r.Group("/user", userAuth)
	{
		user.GET("/home", userH.Home)
		user.GET("/viewproducts", productH.ViewProductsUser)
		user.POST("/add/address", userH.AddAddress)
		user.PATCH("/edit/address", userH.EditAddress)
		user.GET("/address", userH.ShowAddress)
		user.GET("/orderview", orderH.ListOrders)
		user.POST("/cancel/order", orderH.CancelOrder)
		user.GET("/wallet/history", walletH.WalletInfo)
	}

	// Cart routes
	r.POST("/cart/add", userAuth, cartH.AddToCart)
	r.GET("/user/cart/view", userAuth, cartH.ViewCart)
	r.POST("/cart/order", userAuth, orderH.OrderCart)

	// Payment routes
	r.GET("/razorpay", userAuth, paymentH.RazorPay)
	r.GET("/payment-success", userAuth, paymentH.RazorpaySuccess)
	r.GET("/success", userAuth, paymentH.Success)

	// Admin auth (no middleware)
	r.POST("/admin", adminH.Login)
	r.POST("/admin/signup", adminH.Signup)

	// Admin protected routes
	admin := r.Group("/admin", adminAuth)
	{
		// User management
		admin.GET("/users", userH.ViewUsers)
		admin.PATCH("/users/block", userH.BlockUser)
		admin.PATCH("/users/unblock", userH.UnblockUser)
		admin.DELETE("/users/delete", userH.DeleteUser)

		// Category management
		admin.POST("/category/add", categoryH.AddCategory)
		admin.GET("/category/view", categoryH.ViewCategory)
		admin.GET("/category/view/byid", categoryH.ViewProductByCategory)
		admin.PATCH("/category/edit", categoryH.EditCategory)
		admin.DELETE("/category/delete", categoryH.DeleteCategory)

		// Product management
		admin.POST("/product/add", productH.AdminAddProduct)
		admin.GET("/product/view", productH.ViewProducts)
		admin.PATCH("/product/edit/:id", productH.EditProduct)
		admin.DELETE("/product/delete/:id", productH.DeleteProduct)

		// Order management
		admin.GET("/order/view", orderH.ViewOrders)
		admin.PATCH("/order/view", orderH.EditOrder)

		// Coupon management
		admin.GET("/list/coupons", couponH.ListCoupons)
		admin.POST("/add/coupon", couponH.AddCoupon)
		admin.DELETE("/delete/coupon", couponH.DeleteCoupon)

		// Discount management
		admin.POST("/add/discount", discountH.AddDiscount)
		admin.GET("/discount", discountH.GetDiscounts)
		admin.DELETE("/delete/discount", discountH.DeleteDiscount)

		// Attribute definitions
		admin.POST("/category/:id/attributes", productH.AddAttributeDefinition)
		admin.GET("/category/:id/attributes", productH.ListAttributeDefinitions)
		admin.DELETE("/category/attributes/:id", productH.DeleteAttributeDefinition)
	}
}
