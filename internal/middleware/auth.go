package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appErrors "doctor-go/internal/pkg/errors"
	appJWT "doctor-go/internal/pkg/jwt"
	"doctor-go/internal/pkg/response"
	"doctor-go/internal/pkg/token"
)

const UserIDKey = "user_id"
const AdminIDKey = "admin_id"

func Auth(secret string) gin.HandlerFunc {
	return AuthWithTokenService(secret, nil)
}

func AuthWithTokenService(secret string, tokenService *token.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, appErrors.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}

		claims, err := appJWT.Parse(secret, strings.TrimPrefix(header, "Bearer "))
		if err != nil || claims.TokenType != "access" {
			response.Fail(c, http.StatusUnauthorized, appErrors.CodeUnauthorized, "登录已失效")
			c.Abort()
			return
		}
		if tokenService != nil && tokenService.IsBlacklisted(c.Request.Context(), claims.ID) {
			response.Fail(c, http.StatusUnauthorized, appErrors.CodeUnauthorized, "登录已失效")
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

func AdminAuth(secret string) gin.HandlerFunc {
	return AdminAuthWithTokenService(secret, nil)
}

func AdminAuthWithTokenService(secret string, tokenService *token.Service) gin.HandlerFunc {
	return authWithRole(secret, "admin", AdminIDKey, tokenService)
}

func authWithRole(secret string, role string, contextKey string, tokenService *token.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, appErrors.CodeUnauthorized, "未登录")
			c.Abort()
			return
		}

		claims, err := appJWT.Parse(secret, strings.TrimPrefix(header, "Bearer "))
		if err != nil || claims.Role != role || claims.TokenType != "access" {
			response.Fail(c, http.StatusUnauthorized, appErrors.CodeUnauthorized, "登录已失效")
			c.Abort()
			return
		}
		if tokenService != nil && tokenService.IsBlacklisted(c.Request.Context(), claims.ID) {
			response.Fail(c, http.StatusUnauthorized, appErrors.CodeUnauthorized, "登录已失效")
			c.Abort()
			return
		}

		c.Set(contextKey, claims.UserID)
		c.Next()
	}
}
