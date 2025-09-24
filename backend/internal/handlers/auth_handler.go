package handlers

import (
	"backend/internal/pkg/httpx"
	"backend/internal/services"
	"errors"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	HandleLogin(c *gin.Context)
	HandleCreateAccount(c *gin.Context)
	HandleVerifyCode(c *gin.Context)
	HandleRequestCode(c *gin.Context)
}

func (h *authHandler) HandleLogin(c *gin.Context) {

	// 0. Get context
	ctx := c.Request.Context()

	// 1. Bind JSON
	var req struct {
		Email    string `json:"email" binding:"required,email,max=255"`
		Password string `json:"password" binding:"required,max=20,min=8"`
		DeviceID string `json:"device_id" binding:"required,uuid4"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteBadReq(c, "Invalid email or password")
		return
	}

	// 2. Call service
	resp, err := h.svc.Login(ctx, req.Email, req.Password, req.DeviceID)
	if err != nil {
	}

	// 3. Write JSON
	c.JSON(200, resp)
}

func (h *authHandler) HandleCreateAccount(c *gin.Context) {

	// 0. Get context
	ctx := c.Request.Context()
	email := c.GetString("email")
	scene := c.GetString("scene")
	jti := c.GetString("jti")

	// 1. Bind JSON
	var req struct {
		Password string `json:"password" binding:"required,max=20,min=8"`
		DeviceID string `json:"device_id" binding:"required,uuid4"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteBadReq(c, "The length of password should be between 8 and 20.")
		return
	}

	// 2  Call service: Create Account
	resp, err := h.svc.CreateAccount(ctx, email, scene, jti, req.Password, req.DeviceID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrBadRequest):
			httpx.WriteBadReq(c, "Invalid email or password.")
		case errors.Is(err, services.ErrUnauthorized):
			httpx.WriteUnauthorized(c, "Try to signup one more time.")
		case errors.Is(err, services.ErrConflict):
			httpx.WriteConflict(c, "Email already exists. Please login.")
		case errors.Is(err, services.ErrCtxError):
			httpx.WriteCtxError(c, err)
		default:
			httpx.WriteInternal(c)
		}
		return
	}

	// 3. Write JSON
	c.JSON(200, resp)
}

func (h *authHandler) HandleVerifyCode(c *gin.Context) {

	// 0. Get context
	ctx := c.Request.Context()

	// 1. Bind JSON
	var req struct {
		Email  string `json:"email" binding:"required,email,max=255"`
		Scene  string `json:"scene" binding:"required,oneof=signup reset_password"`
		Code   string `json:"code" binding:"required,len=6"`
		CodeID string `json:"code_id" binding:"required,uuid4"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteBadReq(c, "Please enter a valid code")
		return
	}

	// 2. Call service
	token, err := h.svc.VerifyCodeAndGenToken(ctx, req.Email, req.Scene, req.CodeID, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrBadRequest):
			httpx.WriteBadReq(c, "Please enter a valid code.")
		case errors.Is(err, services.ErrUnauthorized):
			httpx.WriteUnauthorized(c, "Verification code is invalid or expired.")
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
		httpx.WriteBadReq(c, "Please check your email.")
		return
	}

	// 2. Call service
	codeID, err := h.svc.RequestCode(ctx, req.Email, req.Scene)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrBadRequest):
			httpx.WriteBadReq(c, "Please check your email.")
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
