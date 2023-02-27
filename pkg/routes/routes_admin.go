package routes

import (
	"ga/middleware"
	"ga/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(Admin *gin.Engine) {

	// admin login create
	Admin.POST("/admin", controllers.AdminLogin)                                        //done
	Admin.POST("/admin/signup", controllers.AdminSignup)                                //done
	Admin.POST("/admin/add", middleware.TokenAuthMiddleware(), controllers.AdminSignup) //done

	// usermanagement
	Admin.GET("/admin/users", middleware.TokenAuthMiddleware(), controllers.ViewUsers)                 //done
	Admin.PATCH("/admin/users/block/:id", middleware.TokenAuthMiddleware(), controllers.BlockUser)     //done
	Admin.PATCH("/admin/users/unblock/:id", middleware.TokenAuthMiddleware(), controllers.UnBlockUser) //done
	Admin.DELETE("/admin/users/delete/:id", middleware.TokenAuthMiddleware(), controllers.DeleteUser)  //done

	//category management
	Admin.POST("/admin/category/add", middleware.TokenAuthMiddleware(), controllers.AddCategory)                //done
	Admin.GET("/admin/category/view", middleware.TokenAuthMiddleware(), controllers.ViewCategory)               //done
	Admin.GET("/admin/category/view/byid", middleware.TokenAuthMiddleware(), controllers.ViewProductByCategory) //done
	Admin.PATCH("/admin/category/edit", middleware.TokenAuthMiddleware(), controllers.EditCategory)             //done
	Admin.DELETE("/admin/category/delete", middleware.TokenAuthMiddleware(), controllers.DeletECategory)        //done

	// product maanagement
	Admin.POST("/admin/product/add", middleware.TokenAuthMiddleware(), controllers.AdminAddProduct)        //done
	Admin.GET("/admin/product/view", middleware.TokenAuthMiddleware(), controllers.ViewProducts)           //done
	Admin.PATCH("/admin/product/edit/:id", middleware.TokenAuthMiddleware(), controllers.EditProduct)      //done
	Admin.DELETE("/admin/product/delete/:id", middleware.TokenAuthMiddleware(), controllers.DeleteProduct) //done

	//cart
	//paymentmethod
	Admin.POST("/admin/paymentmethod", middleware.TokenAuthMiddleware(), controllers.AddPaymentMethod) //done
	// Admin.GET("/admin/order/view", middleware.AdminAuth(), controllers.ViewOrders)
	// Admin.PATCH("/admin/order/update/:id", middleware.AdminAuth(), controllers.EditOrder)
	Admin.GET("/admin/list/coupons", middleware.TokenAuthMiddleware(), controllers.ListCoupons)
	// Admin.GET("/admin/list/discounts", middleware.AdminAuth(), controllers.ListDiscount)
	Admin.POST("/admin/add/coupon", middleware.TokenAuthMiddleware(), controllers.AddCoupon)
	// Admin.POST("/admin/add/discount", middleware.AdminAuth(), controllers.D)

}
