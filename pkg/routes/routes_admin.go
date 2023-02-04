package routes

import (
	"ga/middleware"
	"ga/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(Admin *gin.Engine) {

	// admin login create
	Admin.POST("/admin", controllers.AdminLogin)                              //done
	Admin.POST("/admin/signup", controllers.AdminSignup)                      //done
	Admin.POST("/admin/add", middleware.AdminAuth(), controllers.AdminSignup) //done

	// usermanagement
	Admin.GET("/admin/users", middleware.AdminAuth(), controllers.ViewUsers)                 //done
	Admin.PATCH("/admin/users/block/:id", middleware.AdminAuth(), controllers.BlockUser)     //done
	Admin.PATCH("/admin/users/unblock/:id", middleware.AdminAuth(), controllers.UnBlockUser) //done
	Admin.DELETE("/admin/users/delete/:id", middleware.AdminAuth(), controllers.DeleteUser)  //done

	//category management
	Admin.POST("/admin/category/add", middleware.AdminAuth(), controllers.AddCategory)
	Admin.GET("/admin/category/view", middleware.AdminAuth(), controllers.ViewCategory)
	Admin.PATCH("/admin/category/edit/:id", middleware.AdminAuth(), controllers.EditCategory)
	Admin.DELETE("/admin/category/delete/:id", middleware.AdminAuth(), controllers.DeletECategory)

	// product maanagement
	Admin.POST("/admin/product/add", middleware.AdminAuth(), controllers.AdminAddProduct)
	Admin.GET("/admin/product/view", middleware.AdminAuth(), controllers.ViewProducts) //done
	Admin.PATCH("/admin/product/edit/:id", middleware.AdminAuth(), controllers.EditProduct)
	Admin.DELETE("/admin/product/delete/:id", middleware.AdminAuth(), controllers.DeleteProduct) //done

	//cart
	//paymentmethod
	Admin.POST("/admin/paymentmethod", middleware.AdminAuth(), controllers.AddPaymentMethod)
	// Admin.GET("/admin/order/view", middleware.AdminAuth(), controllers.ViewOrders)
	// Admin.PATCH("/admin/order/update/:id", middleware.AdminAuth(), controllers.EditOrder)
	Admin.GET("/admin/list/coupons", middleware.AdminAuth(), controllers.ListCoupons)
	Admin.GET("/admin/list/discounts", middleware.AdminAuth(), controllers.ListDiscount)
	Admin.POST("/admin/add/coupon", middleware.AdminAuth(), controllers.AddCoupon)
	Admin.POST("/admin/add/discount", middleware.AdminAuth(), controllers.AddDiscount)

}
