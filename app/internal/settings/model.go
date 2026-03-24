package settings

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

func (m *models) getSettings(ctx context.Context) (*model.SystemSettings, error) {
	var settings model.SystemSettings
	err := m.db.WithContext(ctx).First(&settings).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &settings, err
}

func (m *models) createSettings(ctx context.Context, settings *model.SystemSettings) error {
	return m.db.WithContext(ctx).Create(settings).Error
}

func (m *models) updateSettings(ctx context.Context, settings *model.SystemSettings) error {
	return m.db.WithContext(ctx).Save(settings).Error
}

func (m *models) getSettingsByID(ctx context.Context, id uuid.UUID) (*model.SystemSettings, error) {
	var settings model.SystemSettings
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&settings).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &settings, err
}

func (m *models) transaction(ctx context.Context, f func(tx *gorm.DB) error) error {
	return m.db.WithContext(ctx).Transaction(f)
}

func newModels(db *gorm.DB) *models {
	return &models{
		db: db,
	}
}
