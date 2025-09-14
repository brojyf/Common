package handlers

import (
	"Backend/internal/services"
	"Backend/internal/x"
	"net"
	"net/mail"
	"regexp"
	"strings"

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

func (h *authHandler) HandleLogin(c *gin.Context) {
}

func (h *authHandler) HandleRequestCode(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		Email string `json:"email"`
		Scene string `json:"scene"`
	}

	// 400: Invalid Req Body
	if err := c.ShouldBindJSON(&req); err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.BadReq(c)
		return
	}
	if !isValidEmail(req.Email) || !isValidScene(req.Scene) {
		if x.ShouldSkipWrite(c, nil) {
			return
		}
		x.BadRequest(c, "invalid email or scene")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	// 429: Too Many Req
	if err := h.authSvc.CheckRequestCodeThrottle(ctx, req.Email, req.Scene); err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Error(c, "authSvc.CheckRequestCodeThrottle", err)
		x.TooManyReq(c)
		return
	}
	// 500: Internal Server Error
	if err := h.authSvc.RequestCode(ctx, req.Email, req.Scene); err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Error(c, "authSvc.RequestCode", err)
		x.Internal(c)
		return
	}
	// 200: OK
	if x.ShouldSkipWrite(c, nil) {
		return
	}
	c.Status(200)
}

func isValidScene(scene string) bool {
	return scene == "signup" || scene == "reset_password"
}

func isValidEmail(email string) bool {
	var (
		localPartRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+$`)
		domainRegex    = regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	)

	// net/email
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	// 0. Get addr
	address := addr.Address
	parts := strings.Split(address, "@")
	if len(parts) != 2 {
		return false
	}
	localPart, domain := parts[0], parts[1]
	// 1. Simple characters
	if !localPartRegex.MatchString(localPart) {
		return false
	}
	// 2. No IP
	if !domainRegex.MatchString(domain) {
		return false
	}
	if net.ParseIP(domain) != nil {
		return false
	}
	// 3. Length
	if len(domain) > 255 {
		return false
	}

	return true
}
