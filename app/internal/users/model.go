package users

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

func (m *models) createUser(ctx context.Context, user *model.User) error {
	return m.db.WithContext(ctx).Create(user).Error
}

func (m *models) getUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &user, nil
}

func (m *models) updateUser(ctx context.Context, user *model.User) error {
	return m.db.WithContext(ctx).Updates(user).Error
}

func (m *models) deleteUser(ctx context.Context, id uuid.UUID) error {
	return m.db.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&model.User{}).Error
}

func (m *models) listUsers(ctx context.Context, filter UserFilter) ([]*model.User, int64, error) {
	var users []*model.User
	var count int64
	query := m.db.WithContext(ctx).Model(&model.User{})
	if filter.Username != "" {
		query = query.Where("username like ?", "%"+filter.Username+"%")
	}
	if filter.Email != "" {
		query = query.Where("email like ?", "%"+filter.Email+"%")
	}
	if filter.Status != 0 {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Page > 0 && filter.PageSize > 0 {
		query = query.Limit(filter.PageSize).Offset((filter.Page - 1) * filter.PageSize)
	}
	return users, count, query.Find(&users).Error
}

func (m *models) getUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := m.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &user, nil
}

func (m *models) getUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := m.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	return &user, nil
}

func newModels(db *gorm.DB) *models {
	return &models{
		db: db,
	}
}
