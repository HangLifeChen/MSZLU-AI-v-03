package subscriptions

import (
	"app/shared"
	"common/biz"
	"context"

	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/event"
	"github.com/mszlu521/thunder/logs"
)

type PublicService struct {
	repo repository
}

func NewPublicService() *PublicService {
	return &PublicService{
		repo: newModels(database.GetPostgresDB().GormDB),
	}
}

func (s *PublicService) UpdateCurrentSubscription(e event.Event) (any, error) {
	req := e.Data.(*shared.UpdateCurrentBaseSubscriptionReq)
	ctx := context.Background()
	userID := req.UserId
	// 查询现有订阅
	subscription, err := s.repo.getUserSubscription(ctx, userID)
	if err != nil {
		logs.Errorf("查询用户订阅失败: %v", err)
		return nil, errs.DBError
	}
	if subscription == nil {
		return nil, biz.ErrSubscriptionNotFound
	}
	if req.UsedAgents != 0 {
		subscription.UsedAgents = req.UsedAgents
	}
	if req.UsedWorkflows != 0 {
		subscription.UsedWorkflows = req.UsedWorkflows
	}
	if req.UsedKnowledgeBaseSize != 0 {
		subscription.UsedKnowledgeBaseSize = req.UsedKnowledgeBaseSize / 1024 / 1024
	}
	// // 更新订阅
	// subscription.UsedAgents = req.UsedAgents
	// subscription.UsedWorkflows = req.UsedWorkflows
	// subscription.UsedKnowledgeBaseSize = req.UsedKnowledgeBaseSize

	err = s.repo.updateUserSubscription(ctx, subscription)
	if err != nil {
		logs.Errorf("更新用户订阅失败: %v", err)
		return nil, errs.DBError
	}

	return subscription, nil
}
