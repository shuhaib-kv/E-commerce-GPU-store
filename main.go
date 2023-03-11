package main

import (
	"ga/initializers"
	"ga/pkg/database"
	"ga/pkg/routes"
	"os"

	"github.com/gin-gonic/gin"
)

var app = gin.Default()

func init() {
	database.ConnectDB()
	initializers.LoadEnvVariables()
	app.LoadHTMLGlob("templates/*.html")
}

func main() {
	port := os.Getenv("PORT")
	routes.UserRoutes(app)
	routes.AdminRoutes(app)
	app.Run(port)

}
