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
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	admin := &domain.Admin{Name: body.Name, Email: body.Email, Password: body.Password}
	if err := h.adminUC.Signup(c.Request.Context(), admin); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, "Admin created", nil)
}

func (h *AdminHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	admin, err := h.adminUC.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		respondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := middleware.GenerateJWT(admin.Email, admin.ID.Hex(), "admin", h.cfg.JWTSecret)
	if err != nil {
		respondInternalError(c, err, "AdminLogin.GenerateJWT")
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Adminjwt", token, 86400, "/", "", false, true)
	respondSuccess(c, http.StatusOK, "Login successful", nil)
}

func (h *AdminHandler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Adminjwt", "", -1, "/", "", false, true)
	respondSuccess(c, http.StatusOK, "Logged out", nil)
}
