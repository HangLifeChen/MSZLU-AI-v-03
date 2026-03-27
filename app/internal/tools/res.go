package tools

import (
	"model"
	"time"

	"github.com/google/uuid"
)

type TestToolResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ToolResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Config      model.JSON `json:"config"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
}

// ============== 管理员响应类型 ==============

// ToolDetailResponse 工具详情响应（管理员）
type ToolDetailResponse struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	ToolType         model.ToolType         `json:"toolType"`
	IsEnable         bool                   `json:"isEnable"`
	ParametersSchema model.ParametersSchema `json:"parametersSchema"`
	McpConfig        *model.McpConfig       `json:"mcpConfig"`
	CreatorID        string                 `json:"creatorId"`
	CreatorName      string                 `json:"creatorName"`
	CreatorEmail     string                 `json:"creatorEmail"`
	CreatedAt        string                 `json:"createdAt"`
	UpdatedAt        string                 `json:"updatedAt"`
}

// ToolListResponse 工具列表响应项（管理员）
type ToolListResponse struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	ToolType     model.ToolType `json:"toolType"`
	IsEnable     bool           `json:"isEnable"`
	CreatorID    string         `json:"creatorId"`
	CreatorName  string         `json:"creatorName"`
	CreatorEmail string         `json:"creatorEmail"`
	CreatedAt    string         `json:"createdAt"`
	UpdatedAt    string         `json:"updatedAt"`
}

// ListToolsAdminResponse 工具列表响应（管理员）
type ListToolsAdminResponse struct {
	List        []*ToolListResponse `json:"list"`
	Total       int64               `json:"total"`
	CurrentPage int64               `json:"currentPage"`
	PageSize    int64               `json:"pageSize"`
}

// CreatorInfo 创建者信息
type CreatorInfo struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// ToolStats 工具统计信息
type ToolStats struct {
	TotalTools    int64 `json:"totalTools"`
	EnabledTools  int64 `json:"enabledTools"`
	DisabledTools int64 `json:"disabledTools"`
	McpTools      int64 `json:"mcpTools"`
	SystemTools   int64 `json:"systemTools"`
}

// ToToolDetailResponse 将model.Tool转换为ToolDetailResponse
func ToToolDetailResponse(tool *model.Tool, user *model.User) *ToolDetailResponse {
	resp := &ToolDetailResponse{
		ID:               tool.ID.String(),
		Name:             tool.Name,
		Description:      tool.Description,
		ToolType:         tool.ToolType,
		IsEnable:         tool.IsEnable,
		ParametersSchema: tool.ParametersSchema,
		McpConfig:        tool.McpConfig,
		CreatorID:        tool.CreatorID.String(),
		CreatedAt:        tool.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        tool.UpdatedAt.Format(time.RFC3339),
	}
	if user != nil {
		resp.CreatorName = user.Username
		resp.CreatorEmail = user.Email
	}
	return resp
}

// ToToolListResponse 将model.Tool转换为ToolListResponse
func ToToolListResponse(tool *model.Tool, user *model.User) *ToolListResponse {
	resp := &ToolListResponse{
		ID:          tool.ID.String(),
		Name:        tool.Name,
		Description: tool.Description,
		ToolType:    tool.ToolType,
		IsEnable:    tool.IsEnable,
		CreatorID:   tool.CreatorID.String(),
		CreatedAt:   tool.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   tool.UpdatedAt.Format(time.RFC3339),
	}
	if user != nil {
		resp.CreatorName = user.Username
		resp.CreatorEmail = user.Email
	}
	return resp
}
