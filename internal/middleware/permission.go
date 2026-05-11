package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appErrors "doctor-go/internal/pkg/errors"
	"doctor-go/internal/pkg/response"
)

type PermissionChecker interface {
	HasPermission(adminID uint64, permissionCode string) (bool, error)
}

func RequirePermission(checker PermissionChecker, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, ok := c.Get(AdminIDKey)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, appErrors.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}
		adminID, ok := value.(uint64)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, appErrors.CodeUnauthorized, "登录已失效")
			c.Abort()
			return
		}

		allowed, err := checker.HasPermission(adminID, permissionCode)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "权限检查失败")
			c.Abort()
			return
		}
		if !allowed {
			response.Fail(c, http.StatusForbidden, appErrors.CodeForbidden, "无权限")
			c.Abort()
			return
		}
		c.Next()
	}
}
