package subscriptions

import (
	"common/biz"
	"context"
	"model"
	"time"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
)

type service struct {
	repo repository
}

func newService() *service {
	return &service{
		repo: newModels(database.GetPostgresDB().GormDB),
	}
}

// getCurrentSubscription 获取当前用户的订阅信息
func (s *service) getCurrentSubscription(ctx context.Context, userID uuid.UUID) (*SubscriptionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subscription, err := s.repo.getUserSubscription(ctx, userID)
	if err != nil {
		logs.Errorf("查询用户订阅失败: %v", err)
		return nil, errs.DBError
	}

	if subscription == nil {
		// 用户没有订阅记录，返回nil表示无订阅
		subscription = &UserSubscription{
			UserId: userID,
			Plan:   model.FreePlan,
		}
		subscription.Plan = model.FreePlan
	}

	// 获取计划配置
	planConfig, err := s.repo.getPlanConfigByPlan(ctx, subscription.Plan)
	if err != nil {
		logs.Errorf("查询计划配置失败: %v", err)
		return nil, errs.DBError
	}

	return toSubscriptionResponse(subscription, planConfig), nil
}

// getSubscriptionPlans 获取所有订阅计划
func (s *service) getSubscriptionPlans(ctx context.Context) ([]*SubscriptionPlanConfigResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	configs, err := s.repo.getAllPlanConfigs(ctx)
	if err != nil {
		logs.Errorf("查询订阅计划失败: %v", err)
		return nil, errs.DBError
	}

	var responses []*SubscriptionPlanConfigResponse
	for _, config := range configs {
		responses = append(responses, toPlanConfigResponse(config))
	}

	return responses, nil
}

// updateSubscription 创建或更新订阅
func (s *service) updateSubscription(ctx context.Context, userID uuid.UUID, req UpdateSubscriptionReq) (*SubscriptionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 获取计划配置
	planConfig, err := s.repo.getPlanConfigByPlan(ctx, req.PlanType)
	if err != nil {
		logs.Errorf("查询计划配置失败: %v", err)
		return nil, errs.DBError
	}
	if planConfig == nil {
		return nil, biz.ErrPlanNotFound
	}

	// 查询现有订阅
	existingSubscription, err := s.repo.getUserSubscription(ctx, userID)
	if err != nil {
		logs.Errorf("查询用户订阅失败: %v", err)
		return nil, errs.DBError
	}

	now := time.Now()
	var subscription *UserSubscription

	if existingSubscription == nil {
		// 创建新订阅
		subscription = &UserSubscription{
			UserId:    userID,
			Plan:      req.PlanType,
			StartDate: now,
			EndDate:   now.AddDate(0, 1, 0), // 默认一个月
		}
		err = s.repo.createUserSubscription(ctx, subscription)
	} else {
		// 更新现有订阅
		existingSubscription.Plan = req.PlanType
		existingSubscription.StartDate = now
		existingSubscription.EndDate = now.AddDate(0, 1, 0)
		subscription = existingSubscription
		err = s.repo.updateUserSubscription(ctx, subscription)
	}

	if err != nil {
		logs.Errorf("保存用户订阅失败: %v", err)
		return nil, errs.DBError
	}

	return toSubscriptionResponse(subscription, planConfig), nil
}

// cancelSubscription 取消订阅
func (s *service) cancelSubscription(ctx context.Context, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 查询现有订阅
	existingSubscription, err := s.repo.getUserSubscription(ctx, userID)
	if err != nil {
		logs.Errorf("查询用户订阅失败: %v", err)
		return errs.DBError
	}

	if existingSubscription == nil {
		return biz.ErrSubscriptionNotFound
	}

	// 删除订阅
	err = s.repo.deleteUserSubscription(ctx, userID)
	if err != nil {
		logs.Errorf("取消订阅失败: %v", err)
		return errs.DBError
	}

	return nil
}

// createWeChatPaymentOrder 创建微信支付订单
func (s *service) createWeChatPaymentOrder(ctx context.Context, userID uuid.UUID, req CreateWeChatPaymentOrderReq) (*WeChatPaymentOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 获取计划配置
	planConfig, err := s.repo.getPlanConfigByPlan(ctx, req.PlanType)
	if err != nil {
		logs.Errorf("查询计划配置失败: %v", err)
		return nil, errs.DBError
	}
	if planConfig == nil {
		return nil, biz.ErrPlanNotFound
	}

	// 计算实际金额
	amount := req.Amount
	if amount <= 0 {
		// 如果没有传入金额，根据计划价格和时长计算
		amount = planConfig.Price
		switch req.Duration {
		case model.Quarterly:
			amount = int64(float64(amount) * 3 * planConfig.QuarterRate)
		case model.Yearly:
			amount = int64(float64(amount) * 12 * planConfig.YearRate)
		}
	}

	// 创建订单
	order := &WeChatPaymentOrder{
		UserId:   userID,
		Plan:     req.PlanType,
		Duration: req.Duration,
		Amount:   amount,
		Paid:     false,
	}

	err = s.repo.createWeChatPaymentOrder(ctx, order)
	if err != nil {
		logs.Errorf("创建微信支付订单失败: %v", err)
		return nil, errs.DBError
	}

	// TODO: 调用微信支付API获取二维码URL
	// 这里需要集成微信支付SDK，暂时返回模拟数据
	qrCodeURL := "weixin://wxpay/bizpayurl?pr=" + order.Id.String()

	// 更新订单的二维码URL
	order.QrCodeUrl = qrCodeURL
	err = s.repo.updateWeChatPaymentOrder(ctx, order)
	if err != nil {
		logs.Errorf("更新微信支付订单失败: %v", err)
		return nil, errs.DBError
	}

	return &WeChatPaymentOrderResponse{
		OrderID:   order.Id,
		QrCodeURL: order.QrCodeUrl,
	}, nil
}

// checkPaymentStatus 检查支付状态
func (s *service) checkPaymentStatus(ctx context.Context, orderID uuid.UUID) (*PaymentStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	order, err := s.repo.getWeChatPaymentOrder(ctx, orderID)
	if err != nil {
		logs.Errorf("查询支付订单失败: %v", err)
		return nil, errs.DBError
	}

	if order == nil {
		return nil, biz.ErrOrderNotFound
	}

	// TODO: 如果订单未支付，可以调用微信支付API查询实际状态
	// 这里暂时直接返回数据库中的状态

	return &PaymentStatusResponse{
		Paid: order.Paid,
	}, nil
}

// 转换函数
func toSubscriptionResponse(subscription *UserSubscription, planConfig *SubscriptionPlanConfig) *SubscriptionResponse {
	if subscription == nil {
		return nil
	}

	var configs *model.PlanConfig
	if planConfig != nil {
		configs = planConfig.Configs
	}

	return &SubscriptionResponse{
		ID:                    subscription.Id,
		UserID:                subscription.UserId,
		Plan:                  string(subscription.Plan),
		StartDate:             subscription.StartDate.Format(time.RFC3339),
		EndDate:               subscription.EndDate.Format(time.RFC3339),
		Configs:               configs,
		UsedAgents:            subscription.UsedAgents,
		UsedWorkflows:         subscription.UsedWorkflows,
		UsedKnowledgeBaseSize: subscription.UsedKnowledgeBaseSize,
		CreatedAt:             subscription.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             subscription.UpdatedAt.Format(time.RFC3339),
	}
}

func toPlanConfigResponse(config *SubscriptionPlanConfig) *SubscriptionPlanConfigResponse {
	if config == nil {
		return nil
	}

	return &SubscriptionPlanConfigResponse{
		ID:          config.Id,
		Name:        config.Name,
		Plan:        string(config.Plan),
		Price:       config.Price,
		Description: config.Description,
		QuarterRate: config.QuarterRate,
		YearRate:    config.YearRate,
		Configs:     config.Configs,
	}
}
