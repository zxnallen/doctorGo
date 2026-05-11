package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appErrors "doctor-go/internal/pkg/errors"
	"doctor-go/internal/pkg/response"
	"doctor-go/internal/pkg/validator"
)

type OperationLogHandler struct {
	service *OperationLogService
}

func NewOperationLogHandler(service *OperationLogService) *OperationLogHandler {
	return &OperationLogHandler{service: service}
}

type OperationLogQuery struct {
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	AdminID  uint64 `form:"admin_id"`
	Action   string `form:"action"`
	Resource string `form:"resource"`
}

// List godoc
// @Summary 后台操作日志列表
// @Tags Admin Operation Logs
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Param admin_id query int false "管理员 ID"
// @Param action query string false "动作"
// @Param resource query string false "资源"
// @Success 200 {object} response.Body
// @Router /admin/operation-logs [get]
func (h *OperationLogHandler) List(c *gin.Context) {
	var req OperationLogQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	var adminID *uint64
	if req.AdminID > 0 {
		adminID = &req.AdminID
	}

	items, total, err := h.service.List(OperationLogListRequest{
		Page:     req.Page,
		Size:     req.Size,
		AdminID:  adminID,
		Action:   req.Action,
		Resource: req.Resource,
	})
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "获取操作日志失败")
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	response.OK(c, gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
