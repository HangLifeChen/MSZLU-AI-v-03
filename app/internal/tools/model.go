package tools

import (
	"context"
	"model"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/gorms"
	"gorm.io/gorm"
)

type models struct {
	db *gorm.DB
}

func (m *models) getToolsByIds(ctx context.Context, ids []uuid.UUID) ([]*model.Tool, error) {
	var tools []*model.Tool
	return tools, m.db.WithContext(ctx).Where("id in ?", ids).Find(&tools).Error
}

func (m *models) deleteTool(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return m.db.WithContext(ctx).Where("id = ? and creator_id=?", id, userID).Delete(&model.Tool{}).Error
}

func (m *models) getTool(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Tool, error) {
	var tool model.Tool
	err := m.db.WithContext(ctx).Where("id = ? and creator_id=?", id, userID).First(&tool).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &tool, err
}

func (m *models) updateTool(ctx context.Context, info *model.Tool) error {
	return m.db.WithContext(ctx).Updates(info).Error
}

func (m *models) listTools(ctx context.Context, userID uuid.UUID, filter toolFilter) ([]*model.Tool, int64, error) {
	var tools []*model.Tool
	var count int64
	query := m.db.WithContext(ctx).Model(&model.Tool{}).Where("creator_id = ?", userID)
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.ToolType != "" {
		query = query.Where("tool_type = ?", filter.ToolType)
	}
	query = query.Count(&count)
	if filter.Limit != 0 {
		query = query.Limit(filter.Limit).Offset(filter.Offset)
	}
	err := query.Find(&tools).Error
	return tools, count, err
}

type toolFilter struct {
	Name     string
	ToolType model.ToolType
	Limit    int
	Offset   int
}

func (m *models) getToolByName(ctx context.Context, name string) (*model.Tool, error) {
	var tool model.Tool
	err := m.db.WithContext(ctx).Where("name = ?", name).First(&tool).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &tool, err
}

func (m *models) createTool(ctx context.Context, tool *model.Tool) error {
	return m.db.WithContext(ctx).Create(tool).Error
}

func (m *models) get(ctx context.Context, id uuid.UUID) (*model.Tool, error) {
	var tool model.Tool
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&tool).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &tool, err
}

// ============== 管理员数据库方法 ==============

// listToolsAdmin 管理员查询工具列表（可查看所有用户的工具）
func (m *models) listToolsAdmin(ctx context.Context, filter adminToolFilter) ([]*model.Tool, int64, error) {
	var tools []*model.Tool
	var count int64
	query := m.db.WithContext(ctx).Model(&model.Tool{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.ToolType != "" {
		query = query.Where("tool_type = ?", filter.ToolType)
	}
	if filter.CreatorID != uuid.Nil {
		query = query.Where("creator_id = ?", filter.CreatorID)
	}
	if filter.IsEnable != nil {
		query = query.Where("is_enable = ?", *filter.IsEnable)
	}

	// 先获取总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit).Offset(filter.Offset)
	}

	// 按创建时间倒序
	query = query.Order("created_at DESC")

	err := query.Find(&tools).Error
	return tools, count, err
}

// getToolByID 根据ID获取工具详情
func (m *models) getToolByID(ctx context.Context, id uuid.UUID) (*model.Tool, error) {
	var tool model.Tool
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&tool).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &tool, err
}

// getUserByID 根据ID获取用户信息
func (m *models) getUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &user, err
}

// getUsersByIDs 批量获取用户信息
func (m *models) getUsersByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*model.User, error) {
	var users []*model.User
	if len(ids) == 0 {
		return make(map[uuid.UUID]*model.User), nil
	}
	err := m.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	userMap := make(map[uuid.UUID]*model.User)
	for _, user := range users {
		userMap[user.Id] = user
	}
	return userMap, nil
}

type adminToolFilter struct {
	Name      string
	ToolType  model.ToolType
	CreatorID uuid.UUID
	IsEnable  *bool
	Limit     int
	Offset    int
}

// deleteToolAdmin 管理员删除工具
func (m *models) deleteToolAdmin(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Tool{}).Error
}

func newModels(db *gorm.DB) *models {
	return &models{
		db: db,
	}
}
