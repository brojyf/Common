package handlers

import (
	"Backend/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	HandleRefresh(*gin.Context)
	HandleLogin(*gin.Context)
	HandleRequestCode(*gin.Context)
	HandleVerifyCode(*gin.Context)
	HandleCreateAccount(*gin.Context)
	HandleForgetPassword(*gin.Context)
	HandleResetPassword(*gin.Context)
	HandleSetUsername(*gin.Context)
	HandleLogout(*gin.Context)
	HandleLogoutAll(*gin.Context)
}

type authHandler struct {
	authSvc services.AuthService
}

func NewAuthHandler(svc services.AuthService) AuthHandler {
	return &authHandler{authSvc: svc}
}

func (h *authHandler) HandleLogout(c *gin.Context) {}

func (h *authHandler) HandleLogoutAll(c *gin.Context) {}

func (h *authHandler) HandleRefresh(c *gin.Context) {}

func (h *authHandler) HandleSetUsername(c *gin.Context) {}

func (h *authHandler) HandleResetPassword(c *gin.Context) {}

func (h *authHandler) HandleForgetPassword(c *gin.Context) {}

func (h *authHandler) HandleCreateAccount(c *gin.Context) {}

func (h *authHandler) HandleVerifyCode(c *gin.Context) {}

func (h *authHandler) HandleRequestCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Scene string `json:"scene"`
	}
	ctx := c.Request.Context()
	fmt.Print(ctx)

	// 400: Invalid Req Body
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 429: Too Many Req
	if err := h.authSvc.CheckRequestCodeThrottle(req.Email, req.Scene); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}
	// 500: Internal Server Error
	if err := h.authSvc.RequestCode(req.Email, req.Scene); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 200: OK
	c.Status(200)
}

func (h *authHandler) HandleLogin(c *gin.Context) {
}
