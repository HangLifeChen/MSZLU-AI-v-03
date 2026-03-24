package subscriptions

import (
	"context"
	"model"
	"time"

	"github.com/google/uuid"

	"github.com/mszlu521/thunder/gorms"
	"gorm.io/gorm"
)

// UserSubscription 用户订阅模型
type UserSubscription struct {
	Id                    uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserId                uuid.UUID              `json:"userId" gorm:"type:uuid;not null;uniqueIndex"`
	Plan                  model.SubscriptionPlan `json:"plan" gorm:"type:varchar(20);not null"`
	StartDate             time.Time              `json:"startDate" gorm:"type:timestamp;not null"`
	EndDate               time.Time              `json:"endDate" gorm:"type:timestamp;not null"`
	UsedAgents            int64                  `json:"usedAgents" gorm:"default:0"`
	UsedWorkflows         int64                  `json:"usedWorkflows" gorm:"default:0"`
	UsedKnowledgeBaseSize int64                  `json:"usedKnowledgeBaseSize" gorm:"default:0"` // 以MB为单位
	CreatedAt             time.Time              `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt             time.Time              `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (UserSubscription) TableName() string {
	return "user_subscriptions"
}

// SubscriptionPlanConfig 订阅计划配置模型
type SubscriptionPlanConfig struct {
	Id                   int64                  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                 string                 `json:"name" gorm:"type:varchar(50);not null"`
	Plan                 model.SubscriptionPlan `json:"plan" gorm:"type:varchar(20);uniqueIndex;not null"`
	Price                int64                  `json:"price" gorm:"not null"` // 月费价格，以分为单位
	Description          string                 `json:"description" gorm:"type:varchar(255)"`
	QuarterRate          float64                `json:"quarterRate" gorm:"default:0.9"` // 季度折扣率
	YearRate             float64                `json:"yearRate" gorm:"default:0.8"`    // 年度折扣率
	MaxAgents            int64                  `json:"maxAgents" gorm:"default:0"`
	MaxWorkflows         int64                  `json:"maxWorkflows" gorm:"default:0"`
	MaxKnowledgeBaseSize int64                  `json:"maxKnowledgeBaseSize" gorm:"default:0"`
	CreatedAt            time.Time              `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt            time.Time              `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (SubscriptionPlanConfig) TableName() string {
	return "subscription_plan_configs"
}

// WeChatPaymentOrder 微信支付订单模型
type WeChatPaymentOrder struct {
	Id        uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserId    uuid.UUID              `json:"userId" gorm:"type:uuid;not null;index"`
	Plan      model.SubscriptionPlan `json:"plan" gorm:"type:varchar(20);not null"`
	Duration  model.PaymentDuration  `json:"duration" gorm:"type:varchar(20);not null"`
	Amount    int64                  `json:"amount" gorm:"not null"` // 以分为单位
	QrCodeUrl string                 `json:"qrCodeUrl" gorm:"type:varchar(500)"`
	Paid      bool                   `json:"paid" gorm:"default:false"`
	CreatedAt time.Time              `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time              `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (WeChatPaymentOrder) TableName() string {
	return "wechat_payment_orders"
}

// repository 数据访问层接口
type repository interface {
	// 用户订阅相关
	getUserSubscription(ctx context.Context, userID uuid.UUID) (*UserSubscription, error)
	createUserSubscription(ctx context.Context, subscription *UserSubscription) error
	updateUserSubscription(ctx context.Context, subscription *UserSubscription) error
	deleteUserSubscription(ctx context.Context, userID uuid.UUID) error

	// 订阅计划配置相关
	getAllPlanConfigs(ctx context.Context) ([]*SubscriptionPlanConfig, error)
	getPlanConfigByPlan(ctx context.Context, plan model.SubscriptionPlan) (*SubscriptionPlanConfig, error)

	// 微信支付订单相关
	createWeChatPaymentOrder(ctx context.Context, order *WeChatPaymentOrder) error
	getWeChatPaymentOrder(ctx context.Context, orderID uuid.UUID) (*WeChatPaymentOrder, error)
	updateWeChatPaymentOrder(ctx context.Context, order *WeChatPaymentOrder) error
}

// models 实现repository接口
type models struct {
	db *gorm.DB
}

func (m *models) getUserSubscription(ctx context.Context, userID uuid.UUID) (*UserSubscription, error) {
	var subscription UserSubscription
	err := m.db.WithContext(ctx).Where("user_id = ?", userID).First(&subscription).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (m *models) createUserSubscription(ctx context.Context, subscription *UserSubscription) error {
	return m.db.WithContext(ctx).Create(subscription).Error
}

func (m *models) updateUserSubscription(ctx context.Context, subscription *UserSubscription) error {
	return m.db.WithContext(ctx).Save(subscription).Error
}

func (m *models) deleteUserSubscription(ctx context.Context, userID uuid.UUID) error {
	return m.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&UserSubscription{}).Error
}

func (m *models) getAllPlanConfigs(ctx context.Context) ([]*SubscriptionPlanConfig, error) {
	var configs []*SubscriptionPlanConfig
	err := m.db.WithContext(ctx).Find(&configs).Error
	return configs, err
}

func (m *models) getPlanConfigByPlan(ctx context.Context, plan model.SubscriptionPlan) (*SubscriptionPlanConfig, error) {
	var config SubscriptionPlanConfig
	err := m.db.WithContext(ctx).Where("plan = ?", plan).First(&config).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (m *models) createWeChatPaymentOrder(ctx context.Context, order *WeChatPaymentOrder) error {
	return m.db.WithContext(ctx).Create(order).Error
}

func (m *models) getWeChatPaymentOrder(ctx context.Context, orderID uuid.UUID) (*WeChatPaymentOrder, error) {
	var order WeChatPaymentOrder
	err := m.db.WithContext(ctx).Where("id = ?", orderID).First(&order).Error
	if gorms.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (m *models) updateWeChatPaymentOrder(ctx context.Context, order *WeChatPaymentOrder) error {
	return m.db.WithContext(ctx).Save(order).Error
}

func newModels(db *gorm.DB) *models {
	return &models{
		db: db,
	}
}
