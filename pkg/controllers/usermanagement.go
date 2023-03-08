package controllers

import (
	"errors"
	"fmt"
	"ga/pkg/database"
	"ga/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ViewUsers(c *gin.Context) {
	var (
		users []models.Users
		count int64
		query = database.Db.Model(&models.Users{})
	)

	if name := c.Query("name"); name != "" {
		query = query.Where("first_name iLIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	if id := c.Query("id"); id != "" {
		query = query.Where("id = ?", id)
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	query = query.Limit(pageSize).Offset(offset)

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Can't find user",
			"error":   "Database error",
		})
		return
	}

	query.Count(&count)

	var data []gin.H
	for _, u := range users {
		data = append(data, gin.H{
			"id":           u.ID,
			"Name":         u.FirstName + " " + u.LastName,
			"Email":        u.Email,
			"Block Status": u.Block_status,
			"Phone number": u.Phone,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Users",
		"data":    data,
		"page":    page,
		"total":   count,
	})
}
func BlockUser(c *gin.Context) {
	var user models.Users
	var updateStatus bool = true
	var body struct {
		Id uint `json:"id" binding:"required"`
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

	database.Db.First(&user, body.Id)
	if err := database.Db.Model(&user).Where("id=?", body.Id).Update("block_status", updateStatus); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Cant find user",
			"error":   "Database error",
		})
		return
	}
	if body.Id != user.ID {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": " User Does Not Exist",
			"data":    "check user id",
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  true,
		"data":    user.UserName + " is Blocked",
		"message": "Blocked successfully",
	})

}
func UnBlockUser(c *gin.Context) {
	var user models.Users
	var updateStatus bool = false
	var body struct {
		Id uint `json:"id" binding:"required"`
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

	database.Db.First(&user, body.Id)
	if err := database.Db.Model(&user).Where("id=?", body.Id).Update("block_status", updateStatus); err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Cant find user",
			"error":   "Database error",
		})
		return
	}
	if body.Id != user.ID {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": " User Does Not Exist",
			"data":    "check user id",
		})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"status":  true,
		"data":    user.UserName + " is UnBlocked",
		"message": "Unblocked successfully",
	})
}
func DeleteUser(c *gin.Context) {
	var body struct {
		Id uint `json:"id" binding:"required"`
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
	var user models.Users
	result := database.Db.First(&user, body.Id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": " User Does Not Exist",
			"data":    "check user id",
		})
	} else {
		database.Db.Delete(&models.Users{}, body.Id)
		c.JSON(http.StatusFound, gin.H{
			"status":  true,
			"message": "Deleted succesfully",
			"data":    user.UserName + " is Deleted",
		})

	}
}
