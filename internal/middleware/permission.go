package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/s-usmonalizoda25/api-gateway/models/permission"
	"github.com/s-usmonalizoda25/api-gateway/pkg/errs"
	"go.uber.org/zap"
)

func CheckPermission(log *zap.Logger, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			errs.HandleForbiddenError(c, log, "Role not found in token")
			c.Abort()
			return
		}

		var roleStr string
		switch v := role.(type) {
		case string:
			roleStr = v
		case float64:
			if int(v) == 2 {
				roleStr = "ADMIN"
			} else if int(v) == 1 {
				roleStr = "USER"
			} else {
				errs.HandleForbiddenError(c, log, "invalid role value in token")
				c.Abort()
				return
			}
		default:
			errs.HandleForbiddenError(c, log, "invalid role format in token")
			c.Abort()
			return
		}

		if permissions, ok := permission.RolePermission[roleStr]; ok {
			if _, has := permissions[requiredPermission]; has {
				c.Next()
				return
			}
		}

		errs.HandleForbiddenError(c, log, errs.MsgForbidden)
		c.Abort()
	}
}
