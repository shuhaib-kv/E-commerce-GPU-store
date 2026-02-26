package handler

import (
	"ga/config"
	"ga/domain"
	"ga/middleware"
	"ga/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userUC *usecase.UserUsecase
	cfg    *config.Config
}

func NewUserHandler(uuc *usecase.UserUsecase, cfg *config.Config) *UserHandler {
	return &UserHandler{userUC: uuc, cfg: cfg}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var body struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		UserName  string `json:"user_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=4"`
		Phone     string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	user := &domain.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		UserName:  body.UserName,
		Email:     body.Email,
		Password:  body.Password,
		Phone:     body.Phone,
	}

	if err := h.userUC.Signup(c.Request.Context(), user); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, "Account Created", "welcome to Store")
}

func (h *UserHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userUC.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		respondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := middleware.GenerateJWT(user.Email, user.ID.Hex(), h.cfg.JWTSecret)
	if err != nil {
		respondInternalError(c, err, "UserLogin.GenerateJWT")
		return
	}

	c.SetCookie("UserAuth", token, 3600, "/", "", false, true)
	respondSuccess(c, http.StatusOK, "Login successful", nil)
}

func (h *UserHandler) Home(c *gin.Context) {
	email := c.GetString("user_email")
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Welcome", "email": email})
}

func (h *UserHandler) AddAddress(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid user")
		return
	}

	var body struct {
		Name        string `json:"name" binding:"required"`
		PhoneNumber string `json:"phone_number" binding:"required"`
		Pincode     string `json:"pincode" binding:"required"`
		House       string `json:"house"`
		Area        string `json:"area"`
		Landmark    string `json:"landmark"`
		City        string `json:"city" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	addr := &domain.Address{
		UserID:      userID,
		Name:        body.Name,
		PhoneNumber: body.PhoneNumber,
		Pincode:     body.Pincode,
		House:       body.House,
		Area:        body.Area,
		Landmark:    body.Landmark,
		City:        body.City,
	}

	if err := h.userUC.AddAddress(c.Request.Context(), addr); err != nil {
		respondInternalError(c, err, "AddAddress")
		return
	}
	respondSuccess(c, http.StatusCreated, "Address added", nil)
}

func (h *UserHandler) ShowAddress(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid user")
		return
	}
	addr, err := h.userUC.GetAddress(c.Request.Context(), userID)
	if err != nil {
		respondError(c, http.StatusNotFound, "Address not found")
		return
	}
	respondSuccess(c, http.StatusOK, "", addr)
}

func (h *UserHandler) EditAddress(c *gin.Context) {
	userID, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid user")
		return
	}
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.userUC.EditAddress(c.Request.Context(), userID, body); err != nil {
		respondInternalError(c, err, "EditAddress")
		return
	}
	respondSuccess(c, http.StatusOK, "Address updated", nil)
}

// Admin user management handlers
func (h *UserHandler) ViewUsers(c *gin.Context) {
	page := getIntQuery(c, "page", 1)
	pageSize := getIntQuery(c, "pageSize", 10)

	filter := map[string]interface{}{}
	if name := c.Query("name"); name != "" {
		filter["name"] = name
	}
	if idStr := c.Query("id"); idStr != "" {
		if id, err := primitive.ObjectIDFromHex(idStr); err == nil {
			filter["id"] = id
		}
	}

	users, total, err := h.userUC.ListUsers(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		respondInternalError(c, err, "ViewUsers")
		return
	}

	var data []gin.H
	for _, u := range users {
		data = append(data, gin.H{
			"id": u.ID, "first_name": u.FirstName, "last_name": u.LastName,
			"user_name": u.UserName, "email": u.Email, "phone": u.Phone,
			"block_status": u.BlockStatus,
		})
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": data, "total": total})
}

func (h *UserHandler) BlockUser(c *gin.Context) {
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
	if err := h.userUC.BlockUser(c.Request.Context(), id); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, "User blocked", nil)
}

func (h *UserHandler) UnblockUser(c *gin.Context) {
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
	if err := h.userUC.UnblockUser(c.Request.Context(), id); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, "User unblocked", nil)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
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
	if err := h.userUC.DeleteUser(c.Request.Context(), id); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, "User deleted", nil)
}
