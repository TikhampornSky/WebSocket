package middleware

import (
	"fmt"
	"net/http"
	"server/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := service.JWTAuthService().ValidateToken(tokenString)
		if token.Valid {
			fmt.Println(token.Claims)
		} else {
			fmt.Println("unauthorized: ", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
