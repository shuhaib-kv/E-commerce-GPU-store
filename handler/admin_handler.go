package handler

import (
	"ga/config"
	"ga/middleware"
	"ga/usecase"
	"ga/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUC *usecase.AdminUsecase
	cfg     *config.Config
}

func NewAdminHandler(auc *usecase.AdminUsecase, cfg *config.Config) *AdminHandler {
	return &AdminHandler{adminUC: auc, cfg: cfg}
}

func (h *AdminHandler) Signup(c *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=4"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	admin := &domain.Admin{Name: body.Name, Email: body.Email, Password: body.Password}
	if err := h.adminUC.Signup(c.Request.Context(), admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": true, "message": "Admin created"})
}

func (h *AdminHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": err.Error()})
		return
	}

	admin, err := h.adminUC.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": err.Error()})
		return
	}

	token, err := middleware.GenerateJWT(admin.Email, admin.ID.Hex(), h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Failed to generate token"})
		return
	}

	c.SetCookie("Adminjwt", token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "Login successful"})
}
