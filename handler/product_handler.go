package handler

import (
	"encoding/json"
	"ga/domain"
	"ga/usecase"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	productUC *usecase.ProductUsecase
}

func NewProductHandler(puc *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUC: puc}
}

func (h *ProductHandler) AdminAddProduct(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		respondError(c, http.StatusBadRequest, "Product name is required")
		return
	}

	price, err := strconv.Atoi(c.PostForm("price"))
	if err != nil || price <= 0 {
		respondError(c, http.StatusBadRequest, "Valid price is required")
		return
	}
	modelNo, err := strconv.Atoi(c.PostForm("modelno"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Valid model number is required")
		return
	}
	stock, err := strconv.Atoi(c.PostForm("stock"))
	if err != nil || stock < 0 {
		respondError(c, http.StatusBadRequest, "Valid stock is required")
		return
	}

	categoryIDStr := c.PostForm("category_id")
	categoryID, err := primitive.ObjectIDFromHex(categoryIDStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Valid category ID is required")
		return
	}

	description := c.PostForm("description")
	brand := c.PostForm("brand")
	discountIDStr := c.PostForm("discount_id")
	var discountID primitive.ObjectID
	if discountIDStr != "" {
		discountID, err = primitive.ObjectIDFromHex(discountIDStr)
		if err != nil {
			respondError(c, http.StatusBadRequest, "Invalid discount ID")
			return
		}
	}

	// Handle images
	img1 := saveUploadedFile(c, "image1")
	img2 := saveUploadedFile(c, "image2")
	img3 := saveUploadedFile(c, "image3")

	// Parse specifications
	var specs bson.M
	specsStr := c.PostForm("specifications")
	if specsStr != "" {
		if err := json.Unmarshal([]byte(specsStr), &specs); err != nil {
			respondError(c, http.StatusBadRequest, "Invalid specifications JSON")
			return
		}
	}

	product := &domain.Product{
		Name:           name,
		Price:          uint(price),
		ModelNo:        uint(modelNo),
		Image1:         img1,
		Image2:         img2,
		Image3:         img3,
		Stock:          uint(stock),
		CategoryID:     categoryID,
		Description:    description,
		Brand:          brand,
		DiscountID:     discountID,
		Specifications: specs,
	}

	if err := h.productUC.AddProduct(c.Request.Context(), product); err != nil {
		respondInternalError(c, err, "AdminAddProduct")
		return
	}
	respondSuccess(c, http.StatusOK, "Product added", nil)
}

func (h *ProductHandler) EditProduct(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	contentType := c.ContentType()
	update := bson.M{}

	if contentType == "application/json" {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			respondError(c, http.StatusBadRequest, err.Error())
			return
		}
		for k, v := range body {
			update[k] = v
		}
	} else {
		// multipart/form-data — parse text fields
		if v := c.PostForm("name"); v != "" {
			update["name"] = v
		}
		if v := c.PostForm("brand"); v != "" {
			update["brand"] = v
		}
		if v := c.PostForm("description"); v != "" {
			update["description"] = v
		}
		if v, err := strconv.Atoi(c.PostForm("price")); err == nil && v > 0 {
			update["price"] = uint(v)
		}
		if v, err := strconv.Atoi(c.PostForm("modelno")); err == nil {
			update["model_no"] = uint(v)
		}
		if v, err := strconv.Atoi(c.PostForm("stock")); err == nil {
			update["stock"] = uint(v)
		}
		if v := c.PostForm("category_id"); v != "" {
			if oid, err := primitive.ObjectIDFromHex(v); err == nil {
				update["category_id"] = oid
			}
		}
		if v := c.PostForm("discount_id"); v != "" {
			if oid, err := primitive.ObjectIDFromHex(v); err == nil {
				update["discount_id"] = oid
			}
		}
		// Handle image uploads
		if img := saveUploadedFile(c, "image1"); img != "" {
			update["image1"] = img
		}
		if img := saveUploadedFile(c, "image2"); img != "" {
			update["image2"] = img
		}
		if img := saveUploadedFile(c, "image3"); img != "" {
			update["image3"] = img
		}
	}

	// Convert category_id and discount_id strings to ObjectID (for JSON path)
	if catID, ok := update["category_id"].(string); ok {
		if oid, err := primitive.ObjectIDFromHex(catID); err == nil {
			update["category_id"] = oid
		}
	}
	if disID, ok := update["discount_id"].(string); ok {
		if oid, err := primitive.ObjectIDFromHex(disID); err == nil {
			update["discount_id"] = oid
		}
	}

	if err := h.productUC.EditProduct(c.Request.Context(), id, update); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, "Product updated", nil)
}

func (h *ProductHandler) ViewProducts(c *gin.Context) {
	var filter domain.ProductFilter
	c.BindQuery(&filter)

	products, total, err := h.productUC.ViewProducts(c.Request.Context(), filter)
	if err != nil {
		respondInternalError(c, err, "ViewProducts")
		return
	}

	pageSize := filter.PageSize
	if pageSize == 0 {
		pageSize = 10
	}
	pageIndex := filter.PageIndex
	if pageIndex == 0 {
		pageIndex = 1
	}
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	var productJSON []gin.H
	for _, p := range products {
		productJSON = append(productJSON, gin.H{
			"id": p.ID, "name": p.Name, "price": p.Price,
			"image1": p.Image1, "image2": p.Image2, "image3": p.Image3,
			"brand": p.Brand, "discount_id": p.DiscountID,
			"specifications": p.Specifications,
			"stock": p.Stock, "model_no": p.ModelNo,
			"category_id": p.CategoryID, "description": p.Description,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true, "message": "Products found", "data": productJSON,
		"currentPage": pageIndex, "pageSize": pageSize,
		"totalPages": totalPages, "totalItems": total,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid product ID")
		return
	}
	if err := h.productUC.DeleteProduct(c.Request.Context(), id); err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, "Product deleted", nil)
}

func (h *ProductHandler) ViewProductsUser(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)

	filter := bson.M{}
	if name := c.Query("name"); name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if brand := c.Query("brand"); brand != "" {
		filter["brand"] = bson.M{"$regex": brand, "$options": "i"}
	}
	if minPrice, err := strconv.Atoi(c.Query("minPrice")); err == nil {
		filter["price"] = bson.M{"$gte": minPrice}
	}
	if maxPrice, err := strconv.Atoi(c.Query("maxPrice")); err == nil {
		if existing, ok := filter["price"].(bson.M); ok {
			existing["$lte"] = maxPrice
		} else {
			filter["price"] = bson.M{"$lte": maxPrice}
		}
	}
	if idStr := c.Query("id"); idStr != "" {
		if id, err := primitive.ObjectIDFromHex(idStr); err == nil {
			filter["_id"] = id
		}
	}

	response, total, err := h.productUC.ViewProductsUser(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		respondInternalError(c, err, "ViewProductsUser")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true, "page": page, "pageSize": pageSize, "total": total, "data": response,
	})
}

// Attribute definitions
func (h *ProductHandler) AddAttributeDefinition(c *gin.Context) {
	categoryID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var body struct {
		AttributeKey  string   `json:"attribute_key" binding:"required"`
		DisplayName   string   `json:"display_name" binding:"required"`
		AttributeType string   `json:"attribute_type" binding:"required"`
		Required      bool     `json:"required"`
		Options       []string `json:"options"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	validTypes := map[string]bool{"string": true, "number": true, "enum": true}
	if !validTypes[body.AttributeType] {
		respondError(c, http.StatusBadRequest, "attribute_type must be 'string', 'number', or 'enum'")
		return
	}

	attr := &domain.ProductAttributeDefinition{
		CategoryID:    categoryID,
		AttributeKey:  body.AttributeKey,
		DisplayName:   body.DisplayName,
		AttributeType: body.AttributeType,
		Required:      body.Required,
		Options:       body.Options,
	}

	if err := h.productUC.AddAttributeDefinition(c.Request.Context(), attr); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusCreated, "Attribute definition created", attr)
}

func (h *ProductHandler) ListAttributeDefinitions(c *gin.Context) {
	categoryID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid category ID")
		return
	}
	attrs, err := h.productUC.ListAttributeDefinitions(c.Request.Context(), categoryID)
	if err != nil {
		respondInternalError(c, err, "ListAttributeDefinitions")
		return
	}
	respondSuccess(c, http.StatusOK, "", attrs)
}

func (h *ProductHandler) DeleteAttributeDefinition(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.productUC.DeleteAttributeDefinition(c.Request.Context(), id); err != nil {
		respondInternalError(c, err, "DeleteAttributeDefinition")
		return
	}
	respondSuccess(c, http.StatusOK, "Attribute definition deleted", nil)
}

func saveUploadedFile(c *gin.Context, key string) string {
	file, err := c.FormFile(key)
	if err != nil {
		return ""
	}
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	if err := c.SaveUploadedFile(file, "./public/images/"+filename); err != nil {
		log.Println("saveUploadedFile error:", err)
		return ""
	}
	return filename
}
