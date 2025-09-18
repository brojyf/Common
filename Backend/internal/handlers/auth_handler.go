package handlers

import (
	"Backend/internal/pkg/httpx"
	"Backend/internal/services"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	HandleRequestCode(*gin.Context)
}

func (h *authHandler) HandleRequestCode(c *gin.Context) {

	// 0. Get context
	ctx := c.Request.Context()

	// 1. Bind JSON
	var req struct {
		Email string `json:"email" binding:"required,email,max=255"`
		Scene string `json:"scene" binding:"required,oneof= signup reset_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("invalid request body: %v", err)
		httpx.WriteBadReq(c)
		return
	}

	// 2. Call service
	otpID, err := h.svc.RequestCode(ctx, req.Email, req.Scene)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrBadRequest):
			httpx.WriteBadReq(c)
		case errors.Is(err, services.ErrTooManyRequest):
			httpx.WriteTooManyReq(c)
		case errors.Is(err, services.ErrCtxError):
			httpx.WriteCtxError(c, err)
		default:
			httpx.WriteInternal(c)
		}
		return
	}

	// 3. Write JSON
	httpx.TryWriteJSON(c, ctx, 200, gin.H{
		"otp_id": otpID,
	})
}

type authHandler struct {
	svc services.AuthService
}

func NewAuthHandler(authSvc services.AuthService) AuthHandler {
	return &authHandler{svc: authSvc}
}
