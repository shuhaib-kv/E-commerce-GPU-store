package main

import (
	"ga/initializers"
	"ga/pkg/database"
	"ga/pkg/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)
var vp = viper.New()
var app = gin.Default()

func init() {
	database.ConnectDB()
	initializers.LoadEnvVariables()
	app.LoadHTMLGlob("templates/*.html")
	
	vp.SetConfigName(".env")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		print(err)
	}
	
	// dsn := os.Getenv("REDIS_DSN")
	//
	//	if len(dsn) == 0 {
	//		dsn = "localhost:6379"
	//	}
	//
	//	client := redis.NewClient(&redis.Options{
	//		Addr: dsn, //redis port
	//	})
	//
	// _, err := client.Ping().Result()
	//
	//	if err != nil {
	//		panic(err)
	//	}
}

func main() {

	// app.Use(gin.Logger())
	routes.UserRoutes(app)
	routes.AdminRoutes(app)
	app.Run() // Port Declaration to serve the routes

}
