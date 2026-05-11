package news

import (
	"github.com/gin-gonic/gin"

	"doctor-go/internal/config"
	"doctor-go/internal/infrastructure/mysql"
	"doctor-go/internal/infrastructure/redis"
	"doctor-go/internal/middleware"
	"doctor-go/internal/modules/admin"
	"doctor-go/internal/pkg/token"
)

func RegisterRoutes(rg *gin.RouterGroup, db *mysql.DB, cache *redis.Client) {
	repo := NewRepository(db)
	service := NewService(repo, cache)
	handler := NewHandler(service)

	group := rg.Group("/news")
	group.GET("/categories", handler.ListCategories)
	group.GET("", handler.List)
	group.GET("/:id", handler.Detail)
}

func RegisterAdminRoutes(rg *gin.RouterGroup, cfg *config.Config, db *mysql.DB, cache *redis.Client, tokenService *token.Service) {
	repo := NewRepository(db)
	service := NewService(repo, cache)
	adminRepo := admin.NewRepository(db)
	logRepo := admin.NewOperationLogRepository(db)
	logService := admin.NewOperationLogService(logRepo)
	logHandler := admin.NewOperationLogHandler(logService)
	handler := NewHandlerWithLog(service, logService)

	group := rg.Group("/admin/news")
	group.Use(middleware.AdminAuthWithTokenService(cfg.JWT.Secret, tokenService))
	group.GET("", middleware.RequirePermission(adminRepo, "news:list"), handler.AdminList)
	group.POST("", middleware.RequirePermission(adminRepo, "news:create"), handler.Create)
	group.GET("/:id", middleware.RequirePermission(adminRepo, "news:list"), handler.AdminDetail)
	group.PUT("/:id", middleware.RequirePermission(adminRepo, "news:update"), handler.Update)
	group.PATCH("/:id/status", middleware.RequirePermission(adminRepo, "news:status"), handler.UpdateStatus)
	group.DELETE("/:id", middleware.RequirePermission(adminRepo, "news:delete"), handler.Delete)

	categoryGroup := rg.Group("/admin/news-categories")
	categoryGroup.Use(middleware.AdminAuthWithTokenService(cfg.JWT.Secret, tokenService))
	categoryGroup.GET("", middleware.RequirePermission(adminRepo, "news_category:list"), handler.ListCategories)
	categoryGroup.POST("", middleware.RequirePermission(adminRepo, "news_category:create"), handler.CreateCategory)
	categoryGroup.PUT("/:id", middleware.RequirePermission(adminRepo, "news_category:update"), handler.UpdateCategory)
	categoryGroup.PATCH("/:id/status", middleware.RequirePermission(adminRepo, "news_category:status"), handler.UpdateCategoryStatus)
	categoryGroup.DELETE("/:id", middleware.RequirePermission(adminRepo, "news_category:delete"), handler.DeleteCategory)

	logGroup := rg.Group("/admin/operation-logs")
	logGroup.Use(middleware.AdminAuthWithTokenService(cfg.JWT.Secret, tokenService))
	logGroup.GET("", middleware.RequirePermission(adminRepo, "operation_log:list"), logHandler.List)
}
