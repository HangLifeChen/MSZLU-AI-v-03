package tools

import (
	"common/biz"
	"context"
	"core/ai/mcps"
	"core/ai/tools"
	"encoding/json"
	"fmt"
	"model"
	"time"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/ai/einos"
	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
	"github.com/mszlu521/thunder/res"
)

type service struct {
	repo repository
}

func (s *service) createTool(ctx context.Context, userId uuid.UUID, req CreateToolReq) (*model.Tool, error) {
	//先查询tool名字是否存在 防止重复
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	toolInfo, err := s.repo.getToolByName(ctx, req.Name)
	if err != nil {
		logs.Errorf("get tool by name error: %v", err)
		return nil, errs.DBError
	}
	if toolInfo != nil {
		return nil, biz.ErrToolNameExisted
	}
	//创建tool
	tool := model.Tool{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		ToolType:  req.ToolType,
		IsEnable:  true,
		CreatorID: userId,
	}
	//这个地方我们需要先检查tool是否存在，启动时，我们将tool注册了
	//注意 这个地方 我们只能注册 我们系统中已经开发好的tool
	//这个地方因为有mcp工具的存在，所以这里我们先判断一下
	if req.ToolType == model.McpToolType {
		if req.McpConfig != nil {
			tool.McpConfig = req.McpConfig
		}
		tool.Name = req.Name
		tool.Description = req.Description
	} else {
		//这是系统工具
		invokeParamTool := tools.FindTool(req.Name)
		if invokeParamTool == nil {
			return nil, biz.ErrToolNotExisted
		}
		info, err := invokeParamTool.Info(ctx)
		if err != nil {
			logs.Errorf("get tool info error: %v", err)
			return nil, errs.DBError
		}
		tool.Name = info.Name
		tool.Description = info.Desc
		tool.ParametersSchema = invokeParamTool.Params()
	}
	err = s.repo.createTool(ctx, &tool)
	if err != nil {
		logs.Errorf("create tool error: %v", err)
		return nil, errs.DBError
	}
	return &tool, nil
}

func (s *service) listTools(ctx context.Context, userID uuid.UUID, req ListToolsReq) (*res.Page, error) {
	//构建过滤条件
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	filter := toolFilter{
		Name:     req.Name,
		Limit:    req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		ToolType: req.Type,
	}
	toolList, total, err := s.repo.listTools(ctx, userID, filter)
	if err != nil {
		logs.Errorf("list tools error: %v", err)
		return nil, errs.DBError
	}
	return &res.Page{
		List:        toolList,
		Total:       total,
		CurrentPage: int64(req.Page),
		PageSize:    int64(req.PageSize),
	}, nil
}

func (s *service) updateTool(ctx context.Context, userID uuid.UUID, id uuid.UUID, req UpdateToolReq) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	//先根据id进行查询
	toolInfo, err := s.repo.getTool(ctx, userID, id)
	if err != nil {
		logs.Errorf("get tool error: %v", err)
		return nil, errs.DBError
	}
	if toolInfo == nil {
		return nil, biz.ErrToolNotExisted
	}
	//然后判断名字是否重复
	if req.Name != toolInfo.Name {
		toolInfo1, err := s.repo.getToolByName(ctx, req.Name)
		if err != nil {
			logs.Errorf("get tool by name error: %v", err)
			return nil, errs.DBError
		}
		if toolInfo1 != nil {
			return nil, biz.ErrToolNameExisted
		}
	}
	if req.ToolType == model.McpToolType {
		if req.McpConfig != nil {
			toolInfo.McpConfig = req.McpConfig
		}
	}
	toolInfo.Name = req.Name
	toolInfo.Description = req.Description
	toolInfo.IsEnable = req.IsEnable
	err = s.repo.updateTool(ctx, toolInfo)
	if err != nil {
		logs.Errorf("update tool error: %v", err)
		return nil, errs.DBError
	}
	return toolInfo, nil
}

func (s *service) deleteTool(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := s.repo.deleteTool(ctx, userID, id)
	if err != nil {
		logs.Errorf("delete tool error: %v", err)
		return errs.DBError
	}
	return nil
}

func (s *service) testTool(ctx context.Context, userId uuid.UUID, id uuid.UUID, req TestToolReq) (*TestToolResponse, error) {
	//获取 tool
	toolInfo, err := s.repo.getTool(ctx, userId, id)
	if err != nil {
		logs.Errorf("get tool error: %v", err)
		return nil, errs.DBError
	}
	//查找系统中注册的tool
	invokeParamTool := tools.FindTool(toolInfo.Name)
	if invokeParamTool == nil {
		return nil, biz.ErrToolNotExisted
	}
	//参数转换成json
	params, _ := json.Marshal(req.Params)
	result, err := invokeParamTool.InvokableRun(ctx, string(params))
	if err != nil {
		logs.Errorf("invoke tool error: %v", err)
		return &TestToolResponse{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		}, nil
	}
	return &TestToolResponse{
		Message: "success",
		Success: true,
		Data:    result,
	}, nil
}

// GetTool 获取工具详情
func (s *service) getTool(ctx context.Context, id uuid.UUID) (*ToolResponse, error) {
	// 获取工具
	tool, err := s.repo.get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取工具失败: %w", err)
	}

	// 返回工具信息
	return convertToolToResponse(tool), nil
}
func (s *service) getMcpTools(ctx context.Context, userId uuid.UUID, toolId uuid.UUID) ([]*model.Tool, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	//获取mcp
	tool, err := s.repo.getTool(ctx, userId, toolId)
	if err != nil {
		logs.Errorf("get tool error: %v", err)
		return nil, errs.DBError
	}
	config := tool.McpConfig
	if config == nil {
		return nil, biz.ErrMcpConfigNotExisted
	}
	//获取mcp的tool列表，这里我们需要用到mcp-go这个库
	mcpConfig := einos.McpConfig{
		BaseUrl: config.Url,
		Token:   config.CredentialType,
		Name:    "mszlu-AI",
		Version: "1.0.0",
	}
	mcpTools, err := mcps.GetMCPTool(ctx, &mcpConfig)
	if err != nil {
		logs.Errorf("get mcp tool error: %v", err)
		return nil, biz.ErrGetMcpTools
	}
	//转换为model.Tool
	var toolList []*model.Tool
	for _, mcpTool := range mcpTools {
		toolList = append(toolList, &model.Tool{
			BaseModel: model.BaseModel{
				ID: uuid.New(),
			},
			Name:             mcpTool.Name,
			Description:      mcpTool.Description,
			ToolType:         model.McpToolType,
			CreatorID:        userId,
			ParametersSchema: einos.ConvertSchema(mcpTool.InputSchema),
		})
	}
	return toolList, nil
}
func newService() *service {

	return &service{
		repo: newModels(database.GetPostgresDB().GormDB),
	}
}

func convertToolToResponse(tool *model.Tool) *ToolResponse {
	return &ToolResponse{
		Name:        tool.Name,
		Description: tool.Description,
		CreatedAt:   tool.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   tool.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ============== 管理员服务方法 ==============

// createToolAdmin 管理员创建工具
func (s *service) createToolAdmin(ctx context.Context, req CreateToolAdminReq) (*ToolDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查工具名称是否已存在
	toolInfo, err := s.repo.getToolByName(ctx, req.Name)
	if err != nil {
		logs.Errorf("get tool by name error: %v", err)
		return nil, errs.DBError
	}
	if toolInfo != nil {
		return nil, biz.ErrToolNameExisted
	}

	// 创建工具
	toolModel := &model.Tool{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		ToolType:    req.ToolType,
		IsEnable:    req.IsEnable,
		CreatorID:   req.CreatorID,
		Name:        req.Name,
		Description: req.Description,
	}

	// 根据工具类型处理
	if req.ToolType == model.McpToolType {
		if req.McpConfig != nil {
			toolModel.McpConfig = req.McpConfig
		}
	} else {
		// 系统工具
		invokeParamTool := tools.FindTool(req.Name)
		if invokeParamTool == nil {
			return nil, biz.ErrToolNotExisted
		}
		info, err := invokeParamTool.Info(ctx)
		if err != nil {
			logs.Errorf("get tool info error: %v", err)
			return nil, errs.DBError
		}
		toolModel.Name = info.Name
		toolModel.Description = info.Desc
		toolModel.ParametersSchema = invokeParamTool.Params()
	}

	err = s.repo.createTool(ctx, toolModel)
	if err != nil {
		logs.Errorf("create tool error: %v", err)
		return nil, errs.DBError
	}

	// 获取创建者信息
	user, err := s.repo.getUserByID(ctx, toolModel.CreatorID)
	if err != nil {
		logs.Errorf("get user by id error: %v", err)
	}

	return ToToolDetailResponse(toolModel, user), nil
}

// listToolsAdmin 管理员查询工具列表
func (s *service) listToolsAdmin(ctx context.Context, req ListToolsAdminReq) (*ListToolsAdminResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 构建过滤条件
	filter := adminToolFilter{
		Name:     req.Name,
		ToolType: req.Type,
		IsEnable: req.IsEnable,
	}
	if req.CreatorID != uuid.Nil {
		filter.CreatorID = req.CreatorID
	}
	if req.PageSize > 0 {
		filter.Limit = req.PageSize
		filter.Offset = (req.Page - 1) * req.PageSize
	}

	// 查询工具列表
	toolList, total, err := s.repo.listToolsAdmin(ctx, filter)
	if err != nil {
		logs.Errorf("list tools error: %v", err)
		return nil, errs.DBError
	}

	// 收集所有创建者ID
	creatorIDs := make([]uuid.UUID, 0)
	creatorIDSet := make(map[uuid.UUID]bool)
	for _, tool := range toolList {
		if !creatorIDSet[tool.CreatorID] {
			creatorIDs = append(creatorIDs, tool.CreatorID)
			creatorIDSet[tool.CreatorID] = true
		}
	}

	// 批量获取创建者信息
	userMap, err := s.repo.getUsersByIDs(ctx, creatorIDs)
	if err != nil {
		logs.Errorf("get users by ids error: %v", err)
	}

	// 转换响应
	list := make([]*ToolListResponse, 0, len(toolList))
	for _, tool := range toolList {
		user := userMap[tool.CreatorID]
		list = append(list, ToToolListResponse(tool, user))
	}

	return &ListToolsAdminResponse{
		List:        list,
		Total:       total,
		CurrentPage: int64(req.Page),
		PageSize:    int64(req.PageSize),
	}, nil
}

// getToolAdmin 管理员获取工具详情
func (s *service) getToolAdmin(ctx context.Context, id uuid.UUID) (*ToolDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tool, err := s.repo.getToolByID(ctx, id)
	if err != nil {
		logs.Errorf("get tool error: %v", err)
		return nil, errs.DBError
	}
	if tool == nil {
		return nil, biz.ErrToolNotExisted
	}

	// 获取创建者信息
	user, err := s.repo.getUserByID(ctx, tool.CreatorID)
	if err != nil {
		logs.Errorf("get user error: %v", err)
	}

	return ToToolDetailResponse(tool, user), nil
}

// updateToolAdmin 管理员更新工具
func (s *service) updateToolAdmin(ctx context.Context, req UpdateToolAdminReq) (*ToolDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 获取工具
	toolInfo, err := s.repo.getToolByID(ctx, req.ID)
	if err != nil {
		logs.Errorf("get tool error: %v", err)
		return nil, errs.DBError
	}
	if toolInfo == nil {
		return nil, biz.ErrToolNotExisted
	}

	// 检查名称是否重复
	if req.Name != toolInfo.Name {
		existTool, err := s.repo.getToolByName(ctx, req.Name)
		if err != nil {
			logs.Errorf("get tool by name error: %v", err)
			return nil, errs.DBError
		}
		if existTool != nil {
			return nil, biz.ErrToolNameExisted
		}
	}

	// 更新字段
	toolInfo.Name = req.Name
	toolInfo.Description = req.Description
	toolInfo.IsEnable = req.IsEnable
	if req.ToolType == model.McpToolType && req.McpConfig != nil {
		toolInfo.McpConfig = req.McpConfig
	}

	err = s.repo.updateTool(ctx, toolInfo)
	if err != nil {
		logs.Errorf("update tool error: %v", err)
		return nil, errs.DBError
	}

	// 获取创建者信息
	user, err := s.repo.getUserByID(ctx, toolInfo.CreatorID)
	if err != nil {
		logs.Errorf("get user error: %v", err)
	}

	return ToToolDetailResponse(toolInfo, user), nil
}

// deleteToolAdmin 管理员删除工具
func (s *service) deleteToolAdmin(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查工具是否存在
	tool, err := s.repo.getToolByID(ctx, id)
	if err != nil {
		logs.Errorf("get tool error: %v", err)
		return errs.DBError
	}
	if tool == nil {
		return biz.ErrToolNotExisted
	}

	err = s.repo.deleteToolAdmin(ctx, id)
	if err != nil {
		logs.Errorf("delete tool error: %v", err)
		return errs.DBError
	}
	return nil
}

// getToolStats 获取工具统计信息
func (s *service) getToolStats(ctx context.Context) (*ToolStats, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stats := &ToolStats{}

	// 获取所有工具
	tools, total, err := s.repo.listToolsAdmin(ctx, adminToolFilter{})
	if err != nil {
		return nil, fmt.Errorf("获取工具列表失败: %w", err)
	}

	stats.TotalTools = total
	for _, tool := range tools {
		if tool.IsEnable {
			stats.EnabledTools++
		} else {
			stats.DisabledTools++
		}
		if tool.ToolType == model.McpToolType {
			stats.McpTools++
		} else {
			stats.SystemTools++
		}
	}

	return stats, nil
}
