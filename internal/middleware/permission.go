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

		roleStr := role.(string)
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
