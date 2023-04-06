package middleware

import (
	"fmt"
	"net/http"
	"server/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
			c.Set("userID", token.Claims.(jwt.MapClaims)["id"])
			c.Next()
		} else {
			fmt.Println("unauthorized: ", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
