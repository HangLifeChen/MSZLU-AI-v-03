package admin

import (
	"context"
	"time"

	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
)

type service struct {
	repo repository
}

func (s *service) getUserGrowthTrend(ctx context.Context, year int) (*GrowthTrendResponse, error) {
	data, err := s.repo.countByMonth(ctx, "users", year)
	if err != nil {
		logs.Errorf("getUserGrowthTrend error: %v", err)
		return nil, errs.DBError
	}
	return &GrowthTrendResponse{
		Months: buildMonths(),
		Data:   data,
	}, nil
}

func (s *service) getKnowledgeBaseGrowthTrend(ctx context.Context, year int) (*GrowthTrendResponse, error) {
	data, err := s.repo.countByMonth(ctx, "knowledge_bases", year)
	if err != nil {
		logs.Errorf("getKnowledgeBaseGrowthTrend error: %v", err)
		return nil, errs.DBError
	}
	return &GrowthTrendResponse{
		Months: buildMonths(),
		Data:   data,
	}, nil
}

func (s *service) getAgentGrowthTrend(ctx context.Context, year int) (*GrowthTrendResponse, error) {
	data, err := s.repo.countByMonth(ctx, "agents", year)
	if err != nil {
		logs.Errorf("getAgentGrowthTrend error: %v", err)
		return nil, errs.DBError
	}
	return &GrowthTrendResponse{
		Months: buildMonths(),
		Data:   data,
	}, nil
}

func buildMonths() []string {
	return []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}
}

func newService() *service {
	return &service{
		repo: newModels(database.GetPostgresDB().GormDB),
	}
}

func getYear(year int) int {
	if year <= 0 {
		return time.Now().Year()
	}
	return year
}
