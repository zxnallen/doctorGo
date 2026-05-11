package admin

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

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

// Login godoc
// @Summary 后台管理员登录
// @Tags Admin Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录参数"
// @Success 200 {object} response.Body
// @Router /admin/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if errors.Is(err, ErrBadCredentials) {
		response.Fail(c, http.StatusUnauthorized, appErrors.CodeBadCredentials, "用户名或密码错误")
		return
	}
	if errors.Is(err, ErrAdminDisabled) {
		response.Fail(c, http.StatusForbidden, appErrors.CodeForbidden, "管理员已禁用")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "登录失败")
		return
	}
	response.OK(c, result)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}
	result, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, appErrors.CodeUnauthorized, "刷新登录失败")
		return
	}
	response.OK(c, result)
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	_ = c.ShouldBindJSON(&req)
	h.service.Logout(c.Request.Context(), bearerToken(c), req.RefreshToken)
	response.OK(c, gin.H{"logged_out": true})
}

func bearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}
