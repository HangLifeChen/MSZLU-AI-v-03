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
			UserId:    userID,
			Plan:      model.FreePlan,
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 1, 0),
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
	if existingSubscription != nil {
		subscription = existingSubscription
	}
	var endDate time.Time
	switch req.Duration {
	case model.Monthly:
		endDate = now.AddDate(0, 1, 0)
	case model.Quarterly:
		endDate = now.AddDate(0, 3, 0)
	case model.Yearly:
		endDate = now.AddDate(1, 0, 0)
	default:
		endDate = now.AddDate(0, 1, 0)
	}
	if existingSubscription == nil {
		// 创建新订阅
		subscription = &UserSubscription{
			Id:        uuid.New(),
			UserId:    userID,
			Plan:      req.PlanType,
			StartDate: now,
			EndDate:   endDate,
		}
		err = s.repo.createUserSubscription(ctx, subscription)
	} else {
		// 更新现有订阅
		existingSubscription.Plan = req.PlanType
		existingSubscription.StartDate = now
		existingSubscription.EndDate = endDate
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
		Id:       uuid.New(),
		UserId:   userID,
		Plan:     req.PlanType,
		Duration: req.Duration,
		Amount:   amount,
		Paid:     false,
	}
	//模拟支付成功
	order.Paid = true
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

	if order.Paid {
		// 创建或更新订阅
		_, err := s.updateSubscription(ctx, order.UserId, UpdateSubscriptionReq{
			PlanType: order.Plan,
			Duration: order.Duration,
		})
		if err != nil {
			logs.Errorf("创建或更新订阅失败: %v", err)
			return nil, err
		}
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
		configs = &model.PlanConfig{
			MaxAgents:            planConfig.MaxAgents,
			MaxWorkflows:         planConfig.MaxWorkflows,
			MaxKnowledgeBaseSize: planConfig.MaxKnowledgeBaseSize,
		}
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
	configs := &model.PlanConfig{
		MaxAgents:            config.MaxAgents,
		MaxWorkflows:         config.MaxWorkflows,
		MaxKnowledgeBaseSize: config.MaxKnowledgeBaseSize,
	}
	return &SubscriptionPlanConfigResponse{
		ID:          config.Id,
		Name:        config.Name,
		Plan:        string(config.Plan),
		Price:       config.Price,
		Description: config.Description,
		QuarterRate: config.QuarterRate,
		YearRate:    config.YearRate,
		Configs:     configs,
	}
}

// createPlanConfig 创建订阅计划配置
func (s *service) createPlanConfig(ctx context.Context, req CreatePlanConfigReq) (*SubscriptionPlanConfigResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查计划是否已存在
	existingConfig, err := s.repo.getPlanConfigByPlan(ctx, req.Plan)
	if err != nil {
		logs.Errorf("查询计划配置失败: %v", err)
		return nil, errs.DBError
	}
	if existingConfig != nil {
		return nil, biz.ErrPlanAlreadyExists
	}

	// 设置默认值
	quarterRate := req.QuarterRate
	if quarterRate <= 0 {
		quarterRate = 0.9
	}
	yearRate := req.YearRate
	if yearRate <= 0 {
		yearRate = 0.8
	}

	config := &SubscriptionPlanConfig{
		Name:                 req.Name,
		Plan:                 req.Plan,
		Price:                req.Price,
		Description:          req.Description,
		QuarterRate:          quarterRate,
		YearRate:             yearRate,
		MaxAgents:            req.MaxAgents,
		MaxWorkflows:         req.MaxWorkflows,
		MaxKnowledgeBaseSize: req.MaxKnowledgeBaseSize,
	}

	err = s.repo.createPlanConfig(ctx, config)
	if err != nil {
		logs.Errorf("创建计划配置失败: %v", err)
		return nil, errs.DBError
	}

	return toPlanConfigResponse(config), nil
}

// updatePlanConfig 更新订阅计划配置
func (s *service) updatePlanConfig(ctx context.Context, req UpdatePlanConfigReq) (*SubscriptionPlanConfigResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 获取现有配置
	var existingConfig *SubscriptionPlanConfig
	// 通过遍历所有配置来查找
	configs, err := s.repo.getAllPlanConfigs(ctx)
	if err != nil {
		logs.Errorf("查询计划配置失败: %v", err)
		return nil, errs.DBError
	}
	for _, c := range configs {
		if c.Id == req.ID {
			existingConfig = c
			break
		}
	}
	if existingConfig == nil {
		return nil, biz.ErrPlanNotFound
	}

	// 更新字段
	if req.Name != "" {
		existingConfig.Name = req.Name
	}
	if req.Plan != "" {
		existingConfig.Plan = req.Plan
	}
	if req.Price > 0 {
		existingConfig.Price = req.Price
	}
	if req.Description != "" {
		existingConfig.Description = req.Description
	}
	if req.QuarterRate > 0 {
		existingConfig.QuarterRate = req.QuarterRate
	}
	if req.YearRate > 0 {
		existingConfig.YearRate = req.YearRate
	}
	if req.MaxAgents > 0 {
		existingConfig.MaxAgents = req.MaxAgents
	}
	if req.MaxWorkflows > 0 {
		existingConfig.MaxWorkflows = req.MaxWorkflows
	}
	if req.MaxKnowledgeBaseSize > 0 {
		existingConfig.MaxKnowledgeBaseSize = req.MaxKnowledgeBaseSize
	}

	err = s.repo.updatePlanConfig(ctx, existingConfig)
	if err != nil {
		logs.Errorf("更新计划配置失败: %v", err)
		return nil, errs.DBError
	}

	return toPlanConfigResponse(existingConfig), nil
}

// deletePlanConfig 删除订阅计划配置
func (s *service) deletePlanConfig(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.repo.deletePlanConfig(ctx, id)
	if err != nil {
		logs.Errorf("删除计划配置失败: %v", err)
		return errs.DBError
	}

	return nil
}
