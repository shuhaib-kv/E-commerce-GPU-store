package handler

import (
	"ga/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryHandler struct {
	categoryUC *usecase.CategoryUsecase
	productUC  *usecase.ProductUsecase
}

func NewCategoryHandler(cuc *usecase.CategoryUsecase, puc *usecase.ProductUsecase) *CategoryHandler {
	return &CategoryHandler{categoryUC: cuc, productUC: puc}
}

func (h *CategoryHandler) AddCategory(c *gin.Context) {
	var body struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	cat, err := h.categoryUC.AddCategory(c.Request.Context(), body.Name)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusCreated, "Category created", cat)
}

func (h *CategoryHandler) ViewCategory(c *gin.Context) {
	cats, err := h.categoryUC.ViewAll(c.Request.Context())
	if err != nil {
		respondInternalError(c, err, "ViewCategory")
		return
	}
	respondSuccess(c, http.StatusOK, "", cats)
}

func (h *CategoryHandler) EditCategory(c *gin.Context) {
	var body struct {
		ID   string `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := primitive.ObjectIDFromHex(body.ID)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.categoryUC.EditCategory(c.Request.Context(), id, body.Name); err != nil {
		respondInternalError(c, err, "EditCategory")
		return
	}
	respondSuccess(c, http.StatusOK, "Category updated", nil)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	var body struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := primitive.ObjectIDFromHex(body.ID)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.categoryUC.DeleteCategory(c.Request.Context(), id); err != nil {
		respondInternalError(c, err, "DeleteCategory")
		return
	}
	respondSuccess(c, http.StatusOK, "Category deleted", nil)
}

func (h *CategoryHandler) ViewProductByCategory(c *gin.Context) {
	idStr := c.Query("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}
	filter := bson.M{"category_id": id}
	products, total, err := h.productUC.ViewProductsUser(c.Request.Context(), filter, 1, 100)
	if err != nil {
		respondInternalError(c, err, "ViewProductByCategory")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": products, "total": total})
}
