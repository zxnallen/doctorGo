package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	appErrors "doctor-go/internal/pkg/errors"
	"doctor-go/internal/pkg/response"
	"doctor-go/internal/pkg/validator"
)

type RBACHandler struct {
	service *RBACService
}

func NewRBACHandler(service *RBACService) *RBACHandler {
	return &RBACHandler{service: service}
}

func (h *RBACHandler) ListRoles(c *gin.Context) {
	items, err := h.service.ListRoles()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "获取角色失败")
		return
	}
	response.OK(c, items)
}

func (h *RBACHandler) ListPermissions(c *gin.Context) {
	items, err := h.service.ListPermissions()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "获取权限失败")
		return
	}
	response.OK(c, items)
}

func (h *RBACHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}
	result, err := h.service.CreateRole(req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "创建角色失败")
		return
	}
	response.Created(c, result)
}

func (h *RBACHandler) UpdateRole(c *gin.Context) {
	id, ok := parseUintID(c, "id")
	if !ok {
		return
	}
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}
	result, err := h.service.UpdateRole(id, req)
	if errors.Is(err, ErrRoleNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "角色不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "更新角色失败")
		return
	}
	response.OK(c, result)
}

func (h *RBACHandler) SetRolePermissions(c *gin.Context) {
	id, ok := parseUintID(c, "id")
	if !ok {
		return
	}
	var req SetRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}
	err := h.service.SetRolePermissions(id, req.PermissionIDs)
	if errors.Is(err, ErrRoleNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "角色不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "分配权限失败")
		return
	}
	response.OK(c, gin.H{"updated": true})
}

func (h *RBACHandler) SetAdminRoles(c *gin.Context) {
	adminID, ok := parseUintID(c, "admin_id")
	if !ok {
		return
	}
	var req SetAdminRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}
	if err := h.service.SetAdminRoles(adminID, req.RoleIDs); err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "分配角色失败")
		return
	}
	response.OK(c, gin.H{"updated": true})
}

func parseUintID(c *gin.Context, name string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, "ID 不正确")
		return 0, false
	}
	return id, true
}
