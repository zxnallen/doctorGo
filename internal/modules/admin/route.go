package admin

import (
	"time"

	"github.com/gin-gonic/gin"

	"doctor-go/internal/config"
	"doctor-go/internal/infrastructure/mysql"
	"doctor-go/internal/infrastructure/redis"
	"doctor-go/internal/middleware"
	"doctor-go/internal/pkg/token"
)

func RegisterRoutes(rg *gin.RouterGroup, cfg *config.Config, db *mysql.DB, cache *redis.Client) *gin.RouterGroup {
	repo := NewRepository(db)
	tokenService := token.NewService(cfg.JWT, cache)
	service := NewService(cfg.JWT, repo, tokenService)
	handler := NewHandler(service)
	rbacService := NewRBACService(repo)
	rbacHandler := NewRBACHandler(rbacService)

	admin := rg.Group("/admin")
	admin.POST("/auth/login", middleware.RateLimit(cache, "admin:login", 5, time.Minute), handler.Login)
	admin.POST("/auth/refresh", handler.Refresh)
	admin.POST("/auth/logout", middleware.AdminAuthWithTokenService(cfg.JWT.Secret, tokenService), handler.Logout)

	rbac := admin.Group("/rbac")
	rbac.Use(middleware.AdminAuthWithTokenService(cfg.JWT.Secret, tokenService))
	rbac.GET("/roles", middleware.RequirePermission(repo, "rbac:list"), rbacHandler.ListRoles)
	rbac.POST("/roles", middleware.RequirePermission(repo, "rbac:manage"), rbacHandler.CreateRole)
	rbac.PUT("/roles/:id", middleware.RequirePermission(repo, "rbac:manage"), rbacHandler.UpdateRole)
	rbac.PUT("/roles/:id/permissions", middleware.RequirePermission(repo, "rbac:manage"), rbacHandler.SetRolePermissions)
	rbac.GET("/permissions", middleware.RequirePermission(repo, "rbac:list"), rbacHandler.ListPermissions)
	rbac.PUT("/admins/:admin_id/roles", middleware.RequirePermission(repo, "rbac:manage"), rbacHandler.SetAdminRoles)

	return admin
}
