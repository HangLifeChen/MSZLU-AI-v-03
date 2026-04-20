package admin

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type repository interface {
	countByMonth(ctx context.Context, tableName string, year int) ([]int64, error)
}

type models struct {
	db *gorm.DB
}

func (m *models) countByMonth(ctx context.Context, tableName string, year int) ([]int64, error) {
	results := make([]int64, 12)
	now := time.Now()
	for month := 1; month <= 12; month++ {
		endOfMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
		if endOfMonth.After(now) {
			if month == 1 {
				results[month-1] = 0
			} else {
				results[month-1] = results[month-2]
			}
			continue
		}
		var count int64
		err := m.db.WithContext(ctx).
			Table(tableName).
			Where("created_at <= ?", endOfMonth).
			Count(&count).Error
		if err != nil {
			return nil, err
		}
		results[month-1] = count
	}
	return results, nil
}

func newModels(db *gorm.DB) *models {
	return &models{db: db}
}
