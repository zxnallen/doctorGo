package news

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"doctor-go/internal/middleware"
	"doctor-go/internal/modules/admin"
	appErrors "doctor-go/internal/pkg/errors"
	"doctor-go/internal/pkg/response"
	"doctor-go/internal/pkg/validator"
)

type Handler struct {
	service   *Service
	logWriter *admin.OperationLogService
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func NewHandlerWithLog(service *Service, logWriter *admin.OperationLogService) *Handler {
	return &Handler{service: service, logWriter: logWriter}
}

// List godoc
// @Summary 前台资讯列表
// @Tags News
// @Produce json
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Param category_id query int false "分类 ID"
// @Param keyword query string false "关键词"
// @Success 200 {object} response.Body
// @Router /news [get]
func (h *Handler) List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.List(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "获取资讯列表失败")
		return
	}
	response.OK(c, result)
}

// Detail godoc
// @Summary 前台资讯详情
// @Tags News
// @Produce json
// @Param id path int true "资讯 ID"
// @Success 200 {object} response.Body
// @Router /news/{id} [get]
func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, "资讯 ID 不正确")
		return
	}

	result, err := h.service.Detail(c.Request.Context(), id)
	if errors.Is(err, ErrNewsNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "资讯不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "获取资讯详情失败")
		return
	}
	response.OK(c, result)
}

// AdminList godoc
// @Summary 后台资讯列表
// @Tags Admin News
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Param category_id query int false "分类 ID"
// @Param status query int false "状态"
// @Param keyword query string false "关键词"
// @Success 200 {object} response.Body
// @Router /admin/news [get]
func (h *Handler) AdminList(c *gin.Context) {
	var req AdminListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.AdminList(req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "获取后台资讯列表失败")
		return
	}
	response.OK(c, result)
}

func (h *Handler) AdminDetail(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	result, err := h.service.AdminDetail(id)
	if errors.Is(err, ErrNewsNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "资讯不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "获取后台资讯详情失败")
		return
	}
	response.OK(c, result)
}

// Create godoc
// @Summary 创建资讯
// @Tags Admin News
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateRequest true "资讯参数"
// @Success 201 {object} response.Body
// @Router /admin/news [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "创建资讯失败")
		return
	}
	h.writeLog(c, "create", "news", result.ID, true, "创建资讯")
	response.Created(c, result)
}

// Update godoc
// @Summary 更新资讯
// @Tags Admin News
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "资讯 ID"
// @Param request body UpdateRequest true "资讯参数"
// @Success 200 {object} response.Body
// @Router /admin/news/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.Update(id, req)
	if errors.Is(err, ErrNewsNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "资讯不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "更新资讯失败")
		return
	}
	h.writeLog(c, "update", "news", id, true, "更新资讯")
	response.OK(c, result)
}

// UpdateStatus godoc
// @Summary 更新资讯状态
// @Tags Admin News
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "资讯 ID"
// @Param request body UpdateStatusRequest true "状态参数"
// @Success 200 {object} response.Body
// @Router /admin/news/{id}/status [patch]
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.UpdateStatus(id, req.Status)
	if errors.Is(err, ErrNewsNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "资讯不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "更新资讯状态失败")
		return
	}
	h.writeLog(c, "change_status", "news", id, true, "更新资讯状态")
	response.OK(c, result)
}

// Delete godoc
// @Summary 删除资讯
// @Tags Admin News
// @Security BearerAuth
// @Produce json
// @Param id path int true "资讯 ID"
// @Success 200 {object} response.Body
// @Router /admin/news/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	err := h.service.Delete(id)
	if errors.Is(err, ErrNewsNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "资讯不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "删除资讯失败")
		return
	}
	h.writeLog(c, "delete", "news", id, true, "删除资讯")
	response.OK(c, gin.H{"deleted": true})
}

// ListCategories godoc
// @Summary 资讯分类列表
// @Tags News Categories
// @Produce json
// @Param only_enabled query bool false "只看启用分类"
// @Success 200 {object} response.Body
// @Router /news/categories [get]
func (h *Handler) ListCategories(c *gin.Context) {
	var req CategoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.ListCategories(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "获取资讯分类失败")
		return
	}
	response.OK(c, result)
}

// CreateCategory godoc
// @Summary 创建资讯分类
// @Tags Admin News Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "分类参数"
// @Success 201 {object} response.Body
// @Router /admin/news-categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.CreateCategory(req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "创建资讯分类失败")
		return
	}
	h.writeLog(c, "create", "news_category", result.ID, true, "创建资讯分类")
	response.Created(c, result)
}

// UpdateCategory godoc
// @Summary 更新资讯分类
// @Tags Admin News Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "分类 ID"
// @Param request body UpdateCategoryRequest true "分类参数"
// @Success 200 {object} response.Body
// @Router /admin/news-categories/{id} [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.UpdateCategory(id, req)
	if errors.Is(err, ErrCategoryNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "资讯分类不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "更新资讯分类失败")
		return
	}
	h.writeLog(c, "update", "news_category", id, true, "更新资讯分类")
	response.OK(c, result)
}

// UpdateCategoryStatus godoc
// @Summary 更新资讯分类状态
// @Tags Admin News Categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "分类 ID"
// @Param request body UpdateCategoryStatusRequest true "状态参数"
// @Success 200 {object} response.Body
// @Router /admin/news-categories/{id}/status [patch]
func (h *Handler) UpdateCategoryStatus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var req UpdateCategoryStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, validator.Message(err))
		return
	}

	result, err := h.service.UpdateCategoryStatus(id, req.Status)
	if errors.Is(err, ErrCategoryNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "资讯分类不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "更新资讯分类状态失败")
		return
	}
	h.writeLog(c, "change_status", "news_category", id, true, "更新资讯分类状态")
	response.OK(c, result)
}

// DeleteCategory godoc
// @Summary 删除资讯分类
// @Tags Admin News Categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "分类 ID"
// @Success 200 {object} response.Body
// @Router /admin/news-categories/{id} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	err := h.service.DeleteCategory(id)
	if errors.Is(err, ErrCategoryNotFound) {
		response.Fail(c, http.StatusNotFound, appErrors.CodeNotFound, "资讯分类不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, appErrors.CodeInternal, "删除资讯分类失败")
		return
	}
	h.writeLog(c, "delete", "news_category", id, true, "删除资讯分类")
	response.OK(c, gin.H{"deleted": true})
}

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, appErrors.CodeInvalidParams, "资讯 ID 不正确")
		return 0, false
	}
	return id, true
}

func (h *Handler) writeLog(c *gin.Context, action string, resource string, resourceID uint64, success bool, remark string) {
	if h.logWriter == nil {
		return
	}
	adminID, _ := c.Get(middleware.AdminIDKey)
	username := ""
	if value, ok := c.Get("admin_username"); ok {
		username, _ = value.(string)
	}
	_ = h.logWriter.Create(admin.NewOperationLog(
		adminID.(uint64),
		username,
		action,
		resource,
		resourceID,
		c.Request.Method,
		c.Request.URL.Path,
		c.Writer.Status(),
		success,
		remark,
	))
}
