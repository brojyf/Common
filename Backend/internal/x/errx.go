package x

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Write(c *gin.Context, status int, msg string) {
	if ShouldSkipWrite(c, nil) {
		return
	}
	c.JSON(status, gin.H{"error": msg})
}

func BadReq(c *gin.Context)                 { Write(c, http.StatusBadRequest, "invalid request body") }
func BadRequest(c *gin.Context, msg string) { Write(c, http.StatusBadRequest, msg) }
func TooManyReq(c *gin.Context)             { Write(c, http.StatusTooManyRequests, "too many requests") }
func Internal(c *gin.Context)               { Write(c, http.StatusInternalServerError, "internal server error") }
func Unauthorized(c *gin.Context)           { Write(c, http.StatusUnauthorized, "unauthorized") }
