package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerdocs "doctor-go/docs/swagger"
	"doctor-go/internal/config"
	"doctor-go/internal/middleware"
	"doctor-go/internal/modules/admin"
	"doctor-go/internal/modules/auth"
	"doctor-go/internal/modules/news"
	"doctor-go/internal/modules/upload"
	"doctor-go/internal/pkg/response"
	"doctor-go/internal/pkg/token"
)

func NewRouter(cfg *config.Config, dep *Dependencies) *gin.Engine {
	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	if err := r.SetTrustedProxies(cfg.App.TrustedProxies); err != nil {
		panic(err)
	}

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(dep.Logger))
	r.Use(middleware.Recovery(dep.Logger))
	r.Use(middleware.CORS())

	if cfg.App.Env != "prod" {
		swaggerdocs.SwaggerInfo.BasePath = "/api/v1"
		r.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.PersistAuthorization(true),
		))
	}

	r.GET("/", func(c *gin.Context) {
		response.OK(c, gin.H{
			"name":   cfg.App.Name,
			"status": "running",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{
			"status": "ok",
		})
	})
	r.NoRoute(func(c *gin.Context) {
		response.Fail(c, http.StatusNotFound, 10004, "接口不存在")
	})

	api := r.Group("/api/v1")
	admin.RegisterRoutes(api, cfg, dep.MySQL, dep.Redis)

	auth.RegisterRoutes(api, cfg, dep.MySQL, dep.Redis)
	news.RegisterRoutes(api, dep.MySQL, dep.Redis)
	tokenService := token.NewService(cfg.JWT, dep.Redis)
	news.RegisterAdminRoutes(api, cfg, dep.MySQL, dep.Redis, tokenService)

	protected := api.Group("")
	protected.Use(middleware.AuthWithTokenService(cfg.JWT.Secret, tokenService))
	upload.RegisterRoutes(protected, dep.OSS, dep.MySQL)

	return r
}
