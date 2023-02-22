package controllers

import (
	"errors"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ViewUsers(c *gin.Context) {
	var user []models.Users

	if err := database.Db.Find(&user); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Cant find user",
			"error":   "Database error",
		})
		return
	}
	for _, i := range user {
		c.JSON(http.StatusFound, gin.H{
			"status":  true,
			"message": "Users",
			"data": gin.H{
				"id":           i.ID,
				"Name":         i.FirstName + " " + i.LastName,
				"Email":        i.Email,
				"Block Status": i.Block_status,
				"Phone number": i.Phone,
			},
		})
	}
}
func BlockUser(c *gin.Context) {
	var user models.Users
	var updateStatus bool = true
	var body struct {
		userid uint
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

	database.Db.First(&user, body.userid)
	if err := database.Db.Model(&user).Where("id=?", body.userid).Update("block_status", updateStatus); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Cant find user",
			"error":   "Database error",
		})
		return
	}
	if body.userid != user.ID {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": " User Does Not Exist",
			"data":    "check user id",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": user.UserName + " is Blocked",
		"data":    user,
	})

}
func UnBlockUser(c *gin.Context) {
	var user models.Users
	var updateStatus bool = false
	var body struct {
		userid uint
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

	database.Db.First(&user, body.userid)
	if err := database.Db.Model(&user).Where("id=?", body.userid).Update("block_status", updateStatus); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Cant find user",
			"error":   "Database error",
		})
		return
	}
	if body.userid != user.ID {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": " User Does Not Exist",
			"data":    "check user id",
		})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"status":  true,
		"message": user.UserName + " is UNBlocked",
		"data":    user,
	})
}
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.Users
	result := database.Db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": " User Does Not Exist",
			"data":    "check user id",
		})
	} else {
		database.Db.Delete(&models.Users{}, id)
		c.JSON(http.StatusFound, gin.H{
			"status":  true,
			"message": "Deleted succesfully",
			"data":    user,
		})

	}
}
