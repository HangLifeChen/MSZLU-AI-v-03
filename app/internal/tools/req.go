package tools

import (
	"model"

	"github.com/google/uuid"
)

type CreateToolReq struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ToolType    model.ToolType   `json:"toolType"`
	IsEnable    bool             `json:"isEnable"`
	McpConfig   *model.McpConfig `json:"mcpConfig"`
}

type ListToolsReq struct {
	Name     string         `json:"name" form:"name"`
	Type     model.ToolType `json:"type" form:"type"`
	Page     int            `json:"page" form:"page"`
	PageSize int            `json:"pageSize" form:"pageSize"`
}

type UpdateToolReq struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ToolType    model.ToolType   `json:"toolType"`
	IsEnable    bool             `json:"isEnable"`
	McpConfig   *model.McpConfig `json:"mcpConfig"`
}

type TestToolReq struct {
	Params map[string]interface{} `json:"params"`
}

// ============== 管理员请求类型 ==============

// CreateToolAdminReq 创建工具请求（管理员）
type CreateToolAdminReq struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ToolType    model.ToolType   `json:"toolType"`
	IsEnable    bool             `json:"isEnable"`
	McpConfig   *model.McpConfig `json:"mcpConfig"`
	CreatorID   uuid.UUID        `json:"creatorId"`
}

// UpdateToolAdminReq 更新工具请求（管理员）
type UpdateToolAdminReq struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ToolType    model.ToolType   `json:"toolType"`
	IsEnable    bool             `json:"isEnable"`
	McpConfig   *model.McpConfig `json:"mcpConfig"`
}

// ListToolsAdminReq 查询工具列表请求（管理员）
type ListToolsAdminReq struct {
	Name      string         `json:"name" form:"name"`
	Type      model.ToolType `json:"type" form:"type"`
	CreatorID uuid.UUID      `json:"creatorId" form:"creatorId"`
	IsEnable  *bool          `json:"isEnable" form:"isEnable"`
	Page      int            `json:"page" form:"page"`
	PageSize  int            `json:"pageSize" form:"pageSize"`
}
