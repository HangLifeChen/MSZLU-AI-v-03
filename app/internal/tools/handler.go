package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mszlu521/thunder/req"
	"github.com/mszlu521/thunder/res"
)

type Handler struct {
	service *service
}

func (h *Handler) CreateTool(c *gin.Context) {
	var createReq CreateToolReq
	if err := req.JsonParam(c, &createReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	tool, err := h.service.createTool(c.Request.Context(), userID, createReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, tool)
}

func (h *Handler) ListTools(c *gin.Context) {
	var listReq ListToolsReq
	if err := req.QueryParam(c, &listReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	//判断一下分页 如果没有就查询全部
	tools, err := h.service.listTools(c.Request.Context(), userID, listReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, tools)
}

func (h *Handler) UpdateTool(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}
	var updateReq UpdateToolReq
	if err := req.JsonParam(c, &updateReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	tool, err := h.service.updateTool(c.Request.Context(), userID, id, updateReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, tool)
}

func (h *Handler) DeleteTool(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	err := h.service.deleteTool(c.Request.Context(), userID, id)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, nil)
}

func (h *Handler) GetTool(c *gin.Context) {
	// 获取路径参数
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}
	// 调用服务层获取工具详情
	response, err := h.service.getTool(c.Request.Context(), id)
	if err != nil {
		res.Error(c, err)
		return
	}
	// 返回成功响应
	res.Success(c, response)
}

func (h *Handler) GetMcpTools(c *gin.Context) {
	var mcpId uuid.UUID
	if err := req.Path(c, "mcpId", &mcpId); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	tools, err := h.service.getMcpTools(c.Request.Context(), userID, mcpId)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, tools)
}

func (h *Handler) TestTool(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}
	var testReq TestToolReq
	if err := req.JsonParam(c, &testReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.testTool(c.Request.Context(), userID, id, testReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

// ============== 管理员接口 ==============

// CreateToolAdmin 创建工具（管理员）
// @Summary 创建工具
// @Description 管理员创建工具
// @Tags 管理后台-工具管理
// @Accept json
// @Produce json
// @Param body body CreateToolAdminReq true "创建工具请求"
// @Success 200 {object} ToolDetailResponse
// @Router /api/v1/admin/tools [post]
func (h *Handler) CreateToolAdmin(c *gin.Context) {
	var createReq CreateToolAdminReq
	if err := req.JsonParam(c, &createReq); err != nil {
		return
	}

	resp, err := h.service.createToolAdmin(c.Request.Context(), createReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

// ListToolsAdmin 查询工具列表（管理员）
// @Summary 查询工具列表
// @Description 管理员查询工具列表，可按用户筛选
// @Tags 管理后台-工具管理
// @Accept json
// @Produce json
// @Param name query string false "工具名称"
// @Param type query string false "工具类型"
// @Param creatorId query string false "创建者ID"
// @Param isEnable query bool false "是否启用"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} ListToolsAdminResponse
// @Router /api/v1/admin/tools [get]
func (h *Handler) ListToolsAdmin(c *gin.Context) {
	var listReq ListToolsAdminReq
	if err := req.QueryParam(c, &listReq); err != nil {
		return
	}

	// 设置默认分页
	if listReq.Page <= 0 {
		listReq.Page = 1
	}
	if listReq.PageSize <= 0 {
		listReq.PageSize = 10
	}

	resp, err := h.service.listToolsAdmin(c.Request.Context(), listReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

// GetToolAdmin 获取工具详情（管理员）
// @Summary 获取工具详情
// @Description 管理员获取工具详情，包含创建者信息
// @Tags 管理后台-工具管理
// @Accept json
// @Produce json
// @Param id path string true "工具ID"
// @Success 200 {object} ToolDetailResponse
// @Router /api/v1/admin/tools/{id} [get]
func (h *Handler) GetToolAdmin(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}

	resp, err := h.service.getToolAdmin(c.Request.Context(), id)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

// UpdateToolAdmin 更新工具（管理员）
// @Summary 更新工具
// @Description 管理员更新工具
// @Tags 管理后台-工具管理
// @Accept json
// @Produce json
// @Param id path string true "工具ID"
// @Param body body UpdateToolAdminReq true "更新工具请求"
// @Success 200 {object} ToolDetailResponse
// @Router /api/v1/admin/tools/{id} [put]
func (h *Handler) UpdateToolAdmin(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}

	var updateReq UpdateToolAdminReq
	if err := req.JsonParam(c, &updateReq); err != nil {
		return
	}
	updateReq.ID = id

	resp, err := h.service.updateToolAdmin(c.Request.Context(), updateReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

// DeleteToolAdmin 删除工具（管理员）
// @Summary 删除工具
// @Description 管理员删除工具
// @Tags 管理后台-工具管理
// @Accept json
// @Produce json
// @Param id path string true "工具ID"
// @Success 200 {object} nil
// @Router /api/v1/admin/tools/{id} [delete]
func (h *Handler) DeleteToolAdmin(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}

	err := h.service.deleteToolAdmin(c.Request.Context(), id)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, nil)
}

// GetToolStats 获取工具统计信息
// @Summary 获取工具统计信息
// @Description 获取工具总数、启用/禁用数量、MCP/系统工具数量等统计信息
// @Tags 管理后台-工具管理
// @Accept json
// @Produce json
// @Success 200 {object} ToolStats
// @Router /api/v1/admin/tools/stats [get]
func (h *Handler) GetToolStats(c *gin.Context) {
	stats, err := h.service.getToolStats(c.Request.Context())
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, stats)
}

func NewHandler() *Handler {
	return &Handler{
		service: newService(),
	}
}
