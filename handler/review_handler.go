package handler

import (
	"ga/usecase"
	"ga/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewHandler struct {
	reviewUC *usecase.ReviewUsecase
}

func NewReviewHandler(ruc *usecase.ReviewUsecase) *ReviewHandler {
	return &ReviewHandler{reviewUC: ruc}
}

func (h *ReviewHandler) AddReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)
	userObjID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid user")
		return
	}

	var body struct {
		ProductID string `json:"product_id" binding:"required"`
		Rating    int    `json:"rating" binding:"required"`
		Comment   string `json:"comment"`
		UserName  string `json:"user_name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	productObjID, err := primitive.ObjectIDFromHex(body.ProductID)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	review := &domain.Review{
		UserID:    userObjID,
		ProductID: productObjID,
		UserName:  body.UserName,
		Rating:    body.Rating,
		Comment:   body.Comment,
	}

	if err := h.reviewUC.AddReview(c.Request.Context(), review); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusCreated, "Review added", review)
}

func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID, err := primitive.ObjectIDFromHex(c.Param("productId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	reviews, avg, err := h.reviewUC.GetProductReviews(c.Request.Context(), productID)
	if err != nil {
		respondInternalError(c, err, "GetProductReviews")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         true,
		"reviews":        reviews,
		"average_rating": avg,
		"total":          len(reviews),
	})
}
