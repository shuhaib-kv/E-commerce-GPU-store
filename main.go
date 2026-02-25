package main

import (
	"ga/config"
	"ga/database"
	"ga/handler"
	"ga/repository"
	"ga/routes"
	"ga/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Database
	db := database.ConnectMongoDB(cfg.MongoURI, cfg.MongoDatabase)
	defer db.Disconnect()
	mongoDB := db.Database

	// Repositories
	userRepo := repository.NewUserRepo(mongoDB)
	adminRepo := repository.NewAdminRepo(mongoDB)
	productRepo := repository.NewProductRepo(mongoDB)
	categoryRepo := repository.NewCategoryRepo(mongoDB)
	cartRepo := repository.NewCartRepo(mongoDB)
	orderRepo := repository.NewOrderRepo(mongoDB)
	couponRepo := repository.NewCouponRepo(mongoDB)
	discountRepo := repository.NewDiscountRepo(mongoDB)
	walletRepo := repository.NewWalletRepo(mongoDB)
	paymentRepo := repository.NewPaymentRepo(mongoDB)

	// Usecases
	userUC := usecase.NewUserUsecase(userRepo, walletRepo, cartRepo)
	adminUC := usecase.NewAdminUsecase(adminRepo)
	productUC := usecase.NewProductUsecase(productRepo, categoryRepo, discountRepo)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo, productRepo)
	cartUC := usecase.NewCartUsecase(cartRepo, productRepo)
	orderUC := usecase.NewOrderUsecase(orderRepo, cartRepo, couponRepo, walletRepo, discountRepo, productRepo)
	couponUC := usecase.NewCouponUsecase(couponRepo)
	discountUC := usecase.NewDiscountUsecase(discountRepo)
	walletUC := usecase.NewWalletUsecase(walletRepo)
	paymentUC := usecase.NewPaymentUsecase(paymentRepo, orderRepo)

	// Handlers
	userH := handler.NewUserHandler(userUC, cfg)
	adminH := handler.NewAdminHandler(adminUC, cfg)
	productH := handler.NewProductHandler(productUC)
	categoryH := handler.NewCategoryHandler(categoryUC, productUC)
	cartH := handler.NewCartHandler(cartUC)
	orderH := handler.NewOrderHandler(orderUC)
	couponH := handler.NewCouponHandler(couponUC)
	discountH := handler.NewDiscountHandler(discountUC)
	walletH := handler.NewWalletHandler(walletUC)
	paymentH := handler.NewPaymentHandler(paymentUC, cfg)

	// Router
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	routes.Setup(r, cfg.JWTSecret, userH, adminH, productH, categoryH, cartH, orderH, couponH, discountH, walletH, paymentH)

	r.Run(cfg.Port)
}
