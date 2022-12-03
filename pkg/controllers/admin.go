package controllers

import (
	"ga/middleware"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminSignup(c *gin.Context) {
	Email := c.PostForm("email")
	Name := c.PostForm("name")
	Password := c.PostForm("password")
	admin := models.Admin{
		Name:     Name,
		Email:    Email,
		Password: Password,
	}
	result := database.Db.Create(&admin)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create admin",
		})
		return
	}
	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "Account Created",
	})

}

func AdminLogin(c *gin.Context) {

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.Bind(&body)
	var admin models.Admin
	database.Db.First(&admin, "email = ?", body.Email)

	database.Db.Find(&admin)
	if admin.Password != body.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Incorrect Password",
		})
	}
	tokenstring, err := middleware.GenerateJWT(body.Email, int(admin.ID))
	c.SetCookie("Adminjwt", tokenstring, 3600*24*30, "", "", false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "ok",
		"data":    tokenstring,
	})

}
