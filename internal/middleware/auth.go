package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	"github.com/s-usmonalizoda25/api-gateway/pkg/jwt"
	"go.uber.org/zap"
)

func AuthMiddleware(parser *jwt.Parser, log *zap.Logger) gin.HandlerFunc {
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

		claims, err := parser.Parse(tokenString)
		if err != nil {
			errs.HandleAuthError(c, log, errs.MsgUnauthorized)
			c.Abort()
			return
		}

		if tokenType, _ := claims["type"].(string); tokenType != "access" {
			errs.HandleAuthError(c, log, "this endpoint requires an access token, not a refresh token")
			c.Abort()
			return
		}

		c.Set("role", claims["role"])
		c.Set("user_id", claims["user_id"])
		c.Set("token", tokenString)

		c.Next()
	}
}
