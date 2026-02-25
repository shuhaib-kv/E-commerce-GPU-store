package handler

import (
	"ga/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletHandler struct {
	walletUC *usecase.WalletUsecase
}

func NewWalletHandler(wuc *usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{walletUC: wuc}
}

func (h *WalletHandler) WalletInfo(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Invalid user"})
		return
	}

	wallet, history, err := h.walletUC.GetWalletInfo(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"balance": wallet.Balance,
		"history": history,
	})
}
