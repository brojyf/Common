package middleware

import "github.com/gin-gonic/gin"

func OPT() gin.HandlerFunc { 
	return func(c *gin.Context){
		c.Next()
	}
 }

func ATK() gin.HandlerFunc { 
	return func(c *gin.Context){
		c.Next()
	}
 }
