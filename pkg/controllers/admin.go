package controllers

import (
	"ga/middleware"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func AdminSignup(c *gin.Context) {
	var body struct {
		Name     string
		Email    string
		Password string
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "binding json faild",
			"data":    "error ",
		})
		return
	}
	var check []models.Admin
	database.Db.Find(&check)
	for _, i := range check {
		if i.Email == body.Email {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "email Already Exist",
				"data":    "error please enter valid information",
			})
			return
		}
	}
	for _, i := range check {
		if i.Name == body.Name {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Username Already Exist",
				"data":    "error please enter valid information",
			})
			return
		}
	}

	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,

			"error": "Name is required",
			"data":  "error please enter valid information",
		})
		return
	}
	if body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Email is required",
			"data":   "error please enter valid information",
		})
		return
	}
	if body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "Password is required",
			"data":   "error please enter valid information",
		})
		return
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(body.Email) == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"error":  "please enter a valid email",
			"data":   "error please enter valid information",
		})
		return
	}
	admin := models.Admin{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	}
	result := database.Db.Create(&admin)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed to create admin",
			"data":    "error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Account Created",
		"data":    admin,
	})

}

func AdminLogin(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "binding json faild",
			"error":   "error ",
		})
		return
	}
	if body.Email == "" {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  false,
			"message": "Email is required",
			"error":   "error",
		})
		return
	}
	if body.Password == "" {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  false,
			"message": "Password is required",
			"error":   "error",
		})

		return
	}

	var admin models.Admin
	if err := database.Db.First(&admin, "email = ?", body.Email); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Cant find user",
			"error":   "error please enter valid information",
		})
		return
	}
	if admin.Password != body.Password {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  false,
			"message": "Incorrect Password",
			"error":   "error please enter valid information",
		})
		return
	}
	if admin.Email != body.Email {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  false,
			"message": "Incorrect email",
			"error":   "error please enter valid information",
		})
		return
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
