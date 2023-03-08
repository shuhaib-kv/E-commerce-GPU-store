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
	Admin.POST("/admin/category/add", middleware.AdminAuth(), controllers.AddCategory)                //done
	Admin.GET("/admin/category/view", middleware.AdminAuth(), controllers.ViewCategory)               //done
	Admin.GET("/admin/category/view/byid", middleware.AdminAuth(), controllers.ViewProductByCategory) //done
	Admin.PATCH("/admin/category/edit", middleware.AdminAuth(), controllers.EditCategory)             //done
	Admin.DELETE("/admin/category/delete", middleware.AdminAuth(), controllers.DeletECategory)        //done

	// product maanagement
	Admin.POST("/admin/product/add", middleware.AdminAuth(), controllers.AdminAddProduct)        //done
	Admin.GET("/admin/product/view", middleware.AdminAuth(), controllers.ViewProducts)           //done
	Admin.PATCH("/admin/product/edit/:id", middleware.AdminAuth(), controllers.EditProduct)      //done
	Admin.DELETE("/admin/product/delete/:id", middleware.AdminAuth(), controllers.DeleteProduct) //done

	//cart
	//paymentmethod
	Admin.POST("/admin/paymentmethod", middleware.AdminAuth(), controllers.AddPaymentMethod) //done

	Admin.GET("/admin/list/coupons", middleware.AdminAuth(), controllers.ListCoupons)
	Admin.GET("/admin/order/view", middleware.AdminAuth(), controllers.ViewOrders)
	Admin.GET("/admin/discount", middleware.AdminAuth(), controllers.GetDiscountsWithFilters)
	// Admin.PATCH("/admin/order/update/:id", middleware.AdminAuth(), controllers.EditOrder)
	Admin.POST("/admin/add/discount", middleware.AdminAuth(), controllers.AddDiscount)
	Admin.POST("/admin/add/coupon", middleware.AdminAuth(), controllers.AddCoupon)
	// Admin.POST("/admin/add/discount", middleware.AdminAuth(), controllers.D)

}
