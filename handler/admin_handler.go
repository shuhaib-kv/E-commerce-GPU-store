package handler

import (
	"ga/config"
	"ga/domain"
	"ga/middleware"
	"ga/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminUC   *usecase.AdminUsecase
	orderUC   *usecase.OrderUsecase
	productUC *usecase.ProductUsecase
	cfg       *config.Config
}

func NewAdminHandler(auc *usecase.AdminUsecase, ouc *usecase.OrderUsecase, puc *usecase.ProductUsecase, cfg *config.Config) *AdminHandler {
	return &AdminHandler{adminUC: auc, orderUC: ouc, productUC: puc, cfg: cfg}
}

func (h *AdminHandler) Signup(c *gin.Context) {
	var body struct {
		Name           string                `json:"name" binding:"required"`
		Email          string                `json:"email" binding:"required,email"`
		Password       string                `json:"password" binding:"required"`
		StoreName      string                `json:"store_name"`
		PaymentGateway domain.PaymentGateway `json:"payment_gateway"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	admin := &domain.Admin{
		Name:           body.Name,
		Email:          body.Email,
		Password:       body.Password,
		Role:           "admin",
		StoreName:      body.StoreName,
		PaymentGateway: body.PaymentGateway,
	}
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

	role := admin.Role
	if role == "" {
		role = "admin"
	}

	token, err := middleware.GenerateJWT(admin.Email, admin.ID.Hex(), role, h.cfg.JWTSecret)
	if err != nil {
		respondInternalError(c, err, "AdminLogin.GenerateJWT")
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Adminjwt", token, 86400, "/", "", false, true)
	respondSuccess(c, http.StatusOK, "Login successful", gin.H{"role": role})
}

func (h *AdminHandler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Adminjwt", "", -1, "/", "", false, true)
	respondSuccess(c, http.StatusOK, "Logged out", nil)
}

// --- Super Admin handlers ---

func (h *AdminHandler) ListAdmins(c *gin.Context) {
	admins, err := h.adminUC.ListAdmins(c.Request.Context())
	if err != nil {
		respondInternalError(c, err, "ListAdmins")
		return
	}
	respondSuccess(c, http.StatusOK, "Admins fetched", admins)
}

func (h *AdminHandler) GetAdmin(c *gin.Context) {
	id := c.Param("id")
	admin, err := h.adminUC.GetAdmin(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, "Admin fetched", admin)
}

func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	id := c.Param("id")
	admin, err := h.adminUC.GetAdmin(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error())
		return
	}

	var body struct {
		StoreName      *string                `json:"store_name"`
		PaymentGateway *domain.PaymentGateway `json:"payment_gateway"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if body.StoreName != nil {
		admin.StoreName = *body.StoreName
	}
	if body.PaymentGateway != nil {
		admin.PaymentGateway = *body.PaymentGateway
	}

	if err := h.adminUC.UpdateAdmin(c.Request.Context(), admin); err != nil {
		respondInternalError(c, err, "UpdateAdmin")
		return
	}
	respondSuccess(c, http.StatusOK, "Admin updated", admin)
}

func (h *AdminHandler) BlockAdmin(c *gin.Context) {
	id := c.Param("id")
	if err := h.adminUC.BlockAdmin(c.Request.Context(), id); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, "Admin blocked", nil)
}

func (h *AdminHandler) UnblockAdmin(c *gin.Context) {
	id := c.Param("id")
	if err := h.adminUC.UnblockAdmin(c.Request.Context(), id); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}
	respondSuccess(c, http.StatusOK, "Admin unblocked", nil)
}

func (h *AdminHandler) AdminStats(c *gin.Context) {
	ctx := c.Request.Context()

	products, totalProducts, _ := h.productUC.ViewProducts(ctx, domain.ProductFilter{PageIndex: 1, PageSize: 1})
	_ = products

	orders, totalOrders, _ := h.orderUC.ListAllOrders(ctx, map[string]interface{}{}, 1, 10000)

	var totalRevenue uint
	var paidOrders int
	var pendingOrders int
	for _, o := range orders {
		if o.PaymentStatus {
			totalRevenue += o.TotalAmount
			paidOrders++
		}
		if !o.Status {
			pendingOrders++
		}
	}

	respondSuccess(c, http.StatusOK, "Stats fetched", gin.H{
		"total_products": totalProducts,
		"total_orders":   totalOrders,
		"total_revenue":  totalRevenue,
		"paid_orders":    paidOrders,
		"pending_orders": pendingOrders,
	})
}
