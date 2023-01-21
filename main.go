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

}

func main() {

	routes.UserRoutes(app)
	routes.AdminRoutes(app)
	app.Run() // Port Declaration to serve the routes

}
