package settings

import (
	"context"
	"model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository interface {
	getSettings(ctx context.Context) (*model.SystemSettings, error)
	createSettings(ctx context.Context, settings *model.SystemSettings) error
	updateSettings(ctx context.Context, settings *model.SystemSettings) error
	getSettingsByID(ctx context.Context, id uuid.UUID) (*model.SystemSettings, error)
	transaction(ctx context.Context, f func(tx *gorm.DB) error) error
}
