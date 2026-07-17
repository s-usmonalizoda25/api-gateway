package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	"go.uber.org/zap"
)

func AuthMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errs.HandleAuthError(c, log, errs.MsgUnauthorized)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			errs.HandleAuthError(c, log, "Invalid auth header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			errs.HandleAuthError(c, log, errs.MsgUnauthorized)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("role", claims["role"])
			c.Set("user_id", claims["user_id"])
		} else {
			errs.HandleAuthError(c, log, errs.MsgUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}
