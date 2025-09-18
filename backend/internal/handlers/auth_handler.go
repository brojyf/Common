package handlers

import (
	"backend/internal/pkg/httpx"
	"backend/internal/services"
	"errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	HandleVerifyCode(c *gin.Context)
	HandleRequestCode(c *gin.Context)
}

func (h *authHandler) HandleVerifyCode(c *gin.Context) {

	// 0. Get context
	ctx := c.Request.Context()

	// 1. Bind JSON
	var req struct {
		Email  string `json:"email" binding:"required,email,max=255"`
		Scene  string `json:"scene" binding:"required,oneof=signup reset_password"`
		Code   string `json:"code" binding:"required,len=6"`
		CodeID string `json:"code_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteBadReq(c)
		return
	}

	// 2. Call service
	token, err := h.svc.VerifyCodeAndGenToken(ctx, req.Email, req.Scene, req.CodeID, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrBadRequest):
			httpx.WriteBadReq(c)
		case errors.Is(err, services.ErrUnauthorized):
			httpx.WriteUnauthorized(c)
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
		"token": token,
	})
}

func (h *authHandler) HandleRequestCode(c *gin.Context) {

	// 0. Get context
	ctx := c.Request.Context()

	// 1. Bind JSON
	var req struct {
		Email string `json:"email" binding:"required,email,max=255"`
		Scene string `json:"scene" binding:"required,oneof=signup reset_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteBadReq(c)
		return
	}

	// 2. Call service
	codeID, err := h.svc.RequestCode(ctx, req.Email, req.Scene)
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
		"code_id": codeID,
	})
}

type authHandler struct {
	svc services.AuthService
}

func NewAuthHandler(authSvc services.AuthService) AuthHandler {
	return &authHandler{svc: authSvc}
}
