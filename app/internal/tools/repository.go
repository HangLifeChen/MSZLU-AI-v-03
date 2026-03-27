package tools

import (
	"context"
	"model"

	"github.com/google/uuid"
)

type repository interface {
	getToolByName(ctx context.Context, name string) (*model.Tool, error)
	createTool(ctx context.Context, m *model.Tool) error
	listTools(ctx context.Context, userID uuid.UUID, filter toolFilter) ([]*model.Tool, int64, error)
	getTool(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.Tool, error)
	updateTool(ctx context.Context, info *model.Tool) error
	deleteTool(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	getToolsByIds(ctx context.Context, ids []uuid.UUID) ([]*model.Tool, error)
	get(ctx context.Context, id uuid.UUID) (*model.Tool, error)
	// 管理员接口
	listToolsAdmin(ctx context.Context, filter adminToolFilter) ([]*model.Tool, int64, error)
	getToolByID(ctx context.Context, id uuid.UUID) (*model.Tool, error)
	deleteToolAdmin(ctx context.Context, id uuid.UUID) error
	getUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	getUsersByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*model.User, error)
}
