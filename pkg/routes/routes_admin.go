package routes

import (
	"ga/middleware"
	"ga/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(Admin *gin.Engine) {

	Admin.POST("/admin", controllers.AdminLogin)
	Admin.POST("/admin/signup", controllers.AdminSignup)
	Admin.POST("/admin/add", middleware.AdminAuth(), controllers.AdminSignup)

	Admin.GET("/admin/users", middleware.AdminAuth(), controllers.ViewUsers)
	Admin.PATCH("/admin/users/block", middleware.AdminAuth(), controllers.BlockUser)
	Admin.PATCH("/admin/users/unblock", middleware.AdminAuth(), controllers.UnBlockUser)
	Admin.DELETE("/admin/users/delete", middleware.AdminAuth(), controllers.DeleteUser)

	Admin.POST("/admin/category/add", middleware.AdminAuth(), controllers.AddCategory)
	Admin.GET("/admin/category/view", middleware.AdminAuth(), controllers.ViewCategory)
	Admin.GET("/admin/category/view/byid", middleware.AdminAuth(), controllers.ViewProductByCategory)
	Admin.PATCH("/admin/category/edit", middleware.AdminAuth(), controllers.EditCategory)
	Admin.DELETE("/admin/category/delete", middleware.AdminAuth(), controllers.DeletECategory)

	Admin.POST("/admin/product/add", middleware.AdminAuth(), controllers.AdminAddProduct)
	Admin.GET("/admin/product/view", middleware.AdminAuth(), controllers.ViewProducts)
	Admin.PATCH("/admin/product/edit/:id", middleware.AdminAuth(), controllers.EditProduct)
	Admin.DELETE("/admin/product/delete/:id", middleware.AdminAuth(), controllers.DeleteProduct)

	Admin.GET("/admin/order/view", middleware.AdminAuth(), controllers.ViewOrders)
	Admin.PATCH("/admin/order/view", middleware.AdminAuth(), controllers.EditOrder)

	Admin.GET("/admin/list/coupons", middleware.AdminAuth(), controllers.ListCoupons)
	Admin.GET("/admin/discount", middleware.AdminAuth(), controllers.GetDiscountsWithFilters)
	Admin.POST("/admin/add/discount", middleware.AdminAuth(), controllers.AddDiscount)
	Admin.DELETE("/admin/delete/discount", middleware.AdminAuth(), controllers.DeleteDiscount)
	Admin.POST("/admin/add/coupon", middleware.AdminAuth(), controllers.AddCoupon)
}
