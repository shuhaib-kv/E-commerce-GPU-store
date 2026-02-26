package handler

import (
	"ga/domain"
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
		respondError(c, http.StatusBadRequest, "Invalid user")
		return
	}

	wallet, history, err := h.walletUC.GetWalletInfo(c.Request.Context(), userID)
	if err != nil {
		// No wallet yet — return zero balance
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"balance": 0,
			"history": []interface{}{},
		})
		return
	}

	balance := uint(0)
	if wallet != nil {
		balance = wallet.Balance
	}
	if history == nil {
		history = []domain.WalletHistory{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"balance": balance,
		"history": history,
	})
}
