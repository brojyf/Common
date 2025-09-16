package handlers

import (
	"Backend/internal/config"
	"Backend/internal/services"
	"Backend/internal/x"
	"log"
	"net"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *authHandler) HandleLogin(c *gin.Context) {
}

func (h *authHandler) HandleLogout(c *gin.Context) {}

func (h *authHandler) HandleLogoutAll(c *gin.Context) {}

func (h *authHandler) HandleRefresh(c *gin.Context) {}

func (h *authHandler) HandleSetUsername(c *gin.Context) {
	ctx := c.Request.Context()

	// 1) 鉴权：从中间件拿 uid，必须存在且非 0
	uid := c.GetUint64("uid")
	if uid == 0 {
		if x.ShouldSkipWrite(c, nil) {
			return
		}
		x.Unauthorized(c)
		return
	}

	// 2) 解析
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.BadReq(c) // 400: invalid body
		return
	}

	// 3) 规范化 + 校验
	req.Username = strings.TrimSpace(req.Username)
	if !isValidUsername(req.Username) {
		if x.ShouldSkipWrite(c, nil) {
			return
		}
		x.BadReq(c) // 400: invalid username
		return
	}

	// 4) 业务
	err := h.authSvc.UpdateUsername(ctx, uid, req.Username)
	if err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Internal(c) // 500
		return
	}

	// 5) OK
	c.JSON(200, gin.H{
		"username": req.Username,
	})
}

func (h *authHandler) HandleResetPassword(c *gin.Context) {}

func (h *authHandler) HandleForgetPassword(c *gin.Context) {}

func (h *authHandler) HandleCreateAccount(c *gin.Context) {
	ctx := c.Request.Context()

	// 401: Invalid Scene
	scene := c.GetString("scene")
	if scene != "signup" {
		if x.ShouldSkipWrite(c, nil) {
			return
		}
		x.Unauthorized(c)
		return
	}

	// 400: Invalid req body and password
	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.BadReq(c)
		return
	}
	if !isValidPassword(req.Password) {
		if x.ShouldSkipWrite(c, nil) {
			return
		}
		x.BadReq(c)
		return
	}

	// 409 & 500: Store user in db
	email := c.GetString("email")
	uid, err := h.authSvc.CreateAccount(ctx, email, req.Password)
	if err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Internal(c)
		return
	} // 500
	if uid == 0 {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Conflict(c)
		return
	} // 409

	// 500: Store device id
	deviceID := uuid.New().String()
	if err = h.authSvc.StoreDeviceID(ctx, deviceID, uid); err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Internal(c)
		return
	}

	// 500: Sign ARTK
	atk, rtk, err := h.authSvc.SignARTK(ctx, uid, deviceID)
	if err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Internal(c)
		return
	}

	// 201
	if x.ShouldSkipWrite(c, nil) {
		return
	}
	c.JSON(201, gin.H{
		"access_token":  atk,
		"token_type":    "Bearer",
		"expires_in":    config.C.JWT.ATK,
		"refresh_token": rtk,
		"user_id":       uid,
		"device_id":     deviceID,
	})
}

func (h *authHandler) HandleVerifyCode(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		ID    string `json:"otp_id"`
		Code  string `json:"code"`
		Email string `json:"email"`
		Scene string `json:"scene"`
	}
	// 400
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
	// 401
	pass, err := h.authSvc.VerifyCode(ctx, req.Email, req.Scene, req.Code, req.ID)
	if err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Internal(c)
		return
	} // 500
	if !pass {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Unauthorized(c)
		return
	}
	// 500
	jti := uuid.NewString()
	token, err := h.authSvc.SignOTP(ctx, req.Email, req.Scene, jti)
	if err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Internal(c)
		return
	}
	//200
	if x.ShouldSkipWrite(c, nil) {
		return
	}
	c.JSON(200, gin.H{
		"token": token,
	})
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
		x.TooManyReq(c)
		return
	}
	// 500: Internal Server
	codeID := uuid.New().String()
	if err := h.authSvc.RequestCode(ctx, req.Email, req.Scene, codeID); err != nil {
		if x.ShouldSkipWrite(c, err) {
			return
		}
		x.Internal(c)
		return
	}
	// 200: OK
	if x.ShouldSkipWrite(c, nil) {
		return
	}
	log.Printf("[DEV] codeID: %s", codeID)
	c.JSON(200, gin.H{
		"otp_id": codeID,
	})
}

func isValidPassword(pwd string) bool {
	if len(pwd) < 8 || len(pwd) > 20 {
		return false
	}
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, c := range pwd {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasLower && hasUpper && hasDigit && hasSpecial
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

func isValidUsername(s string) bool {
	var usernameRe = regexp.MustCompile(`^[A-Za-z0-9 ]+$`)
	if len(s) == 0 {
		return false
	}
	if utf8.RuneCountInString(s) > 20 {
		return false
	}
	return usernameRe.MatchString(s)
}
