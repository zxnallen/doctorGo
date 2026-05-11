package auth

import (
	"time"

	"github.com/gin-gonic/gin"

	"doctor-go/internal/config"
	"doctor-go/internal/infrastructure/mysql"
	"doctor-go/internal/infrastructure/redis"
	"doctor-go/internal/middleware"
	"doctor-go/internal/pkg/token"
)

func RegisterRoutes(rg *gin.RouterGroup, cfg *config.Config, db *mysql.DB, cache *redis.Client) {
	repo := NewRepository(db)
	tokenService := token.NewService(cfg.JWT, cache)
	service := NewService(cfg.JWT, repo, tokenService)
	handler := NewHandler(service)

	group := rg.Group("/auth")
	group.POST("/register", middleware.RateLimit(cache, "auth:register", 5, time.Minute), handler.Register)
	group.POST("/login", middleware.RateLimit(cache, "auth:login", 10, time.Minute), handler.Login)
	group.POST("/refresh", handler.Refresh)
	group.POST("/logout", middleware.AuthWithTokenService(cfg.JWT.Secret, tokenService), handler.Logout)
}
