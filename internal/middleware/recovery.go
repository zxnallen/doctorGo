package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"doctor-go/internal/infrastructure/logger"
	appErrors "doctor-go/internal/pkg/errors"
	"doctor-go/internal/pkg/response"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Error("panic recovered",
			zap.Any("error", recovered),
			zap.String("request_id", requestID(c)),
		)
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "服务器错误")
	})
}
