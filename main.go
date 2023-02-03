package main

import (
	"ga/initializers"
	"ga/pkg/database"
	"ga/pkg/routes"

	"github.com/gin-gonic/gin"
)

var app = gin.Default()

func init() {
	database.ConnectDB()
	initializers.LoadEnvVariables()
	app.LoadHTMLGlob("templates/*.html")
}

func main() {

	// app.Use(gin.Logger())
	routes.UserRoutes(app)
	routes.AdminRoutes(app)
	app.Run(":8085") // Port Declaration to serve the routes

}
