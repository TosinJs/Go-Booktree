package middleware

import (
	"net/http"
	"strings"
	"tosinjs/gin-booktree/pkg/auth"

	"github.com/gin-gonic/gin"
)

func RequiresAuthToken() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "authentication required/invalid auth token",
			})
			return
		}
		authArray := strings.Split(authorizationHeader, " ")
		if len(authArray) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "authentication required/invalid auth token",
			})
			return
		}
		isValid, err := auth.VerifyJWTToken(authArray[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
			return
		}
		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "authentication required/invalid auth token",
			})
			return
		}
		c.Next()
	}
	return fn
}
