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

func WalletInfo(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.GetString("id"), 10, 32)
	var balance uint
	result := database.Db.Model(&models.Wallet{}).Where("users_id = ?", userID).Select("balance").Scan(&balance)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Failed to query wallet balance",
			})
		}
		return
	}

	var history []models.Wallethistory
	result = database.Db.Model(&models.Wallethistory{}).Where("users_id = ?", userID).Find(&history)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to query wallet history",
		})
		return
	}
	historyResponse := make([]gin.H, len(history))
	for i, item := range history {
		historyResponse[i] = gin.H{
			"date":   item.CreatedAt,
			"debit":  item.Debit,
			"credit": item.Credit,
			"userID": item.UsersID,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{
			"balance": balance,
			"userID":  userID,
			"history": historyResponse,
		},
		"message": "Your Wallet History",
	})
}
