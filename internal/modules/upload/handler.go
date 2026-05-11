package upload

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"doctor-go/internal/middleware"
	appErrors "doctor-go/internal/pkg/errors"
	"doctor-go/internal/pkg/response"
	"doctor-go/internal/pkg/validator"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SignedURL(c *gin.Context) {
	var req SignedURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	userID, _ := c.Get(middleware.UserIDKey)
	result, err := h.service.CreateSignedURL(userID.(uint64), req)
	if errors.Is(err, ErrOSSNotConfigured) {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeOSSNotConfigured, "OSS 未配置")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "创建上传地址失败")
		return
	}
	response.OK(c, result)
}
