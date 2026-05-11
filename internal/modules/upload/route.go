package upload

import (
	"github.com/gin-gonic/gin"

	"doctor-go/internal/infrastructure/mysql"
	"doctor-go/internal/infrastructure/oss"
)

func RegisterRoutes(rg *gin.RouterGroup, ossClient *oss.Client, db *mysql.DB) {
	repo := NewRepository(db)
	service := NewService(ossClient, repo)
	handler := NewHandler(service)

	group := rg.Group("/upload")
	group.POST("/signed-url", handler.SignedURL)
}
