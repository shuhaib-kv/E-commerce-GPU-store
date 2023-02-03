package controllers

import (
	"errors"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ViewUsers(c *gin.Context) {

	var user []models.Users
	database.Db.Find(&user)
	for _, i := range user {
		c.JSON(200, gin.H{
			"id":           i.ID,
			"Name":         i.FirstName + " " + i.LastName,
			"Email":        i.Email,
			"Block Status": i.Block_status,
			"Phone number": i.Phone,
		})
	}
}
func BlockUser(c *gin.Context) {

	var user models.Users
	var updateStatus bool = true
	id := c.Param("id")
	idu, _ := strconv.Atoi(id)
	database.Db.First(&user, id)
	database.Db.Model(&user).Where("id=?", id).Update("block_status", updateStatus)
	if idu != user.ID {
		c.JSON(200, gin.H{
			"status":  false,
			"message": " User Does Not Exist",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": user.UserName + " is Blocked",
	})

}
func UnBlockUser(c *gin.Context) {
	var user models.Users
	var updateStatus bool = false
	id := c.Param("id")
	idu, _ := strconv.Atoi(id)

	database.Db.First(&user, id)
	database.Db.Model(&user).Where("id=?", id).Update("block_status", updateStatus)
	if idu != user.ID {
		c.JSON(200, gin.H{
			"status":  false,
			"message": " User Does Not Exist",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"message": user.UserName + " is UnBlocked",
	})
}
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.Users
	result := database.Db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{
			"message": "can't find user",
		})
	} else {
		database.Db.Delete(&models.Users{}, id)
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Deleted succesfully",
		})

	}
}
func Validate(c *gin.Context) {
	//user := c.GetInt("id")
	check, _ := c.Get("user")

	id := c.GetUint("id")

	c.JSON(http.StatusOK, gin.H{
		"message": id,
		"user":    check,
	})

}
