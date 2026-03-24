package subscriptions

import (
	"model"

	"github.com/google/uuid"
)

// SubscriptionResponse 用户订阅信息响应
type SubscriptionResponse struct {
	ID                    uuid.UUID         `json:"id"`
	UserID                uuid.UUID         `json:"userId"`
	Plan                  string            `json:"plan"`
	StartDate             string            `json:"startDate"`
	EndDate               string            `json:"endDate"`
	Configs               *model.PlanConfig `json:"configs"`
	UsedAgents            int64             `json:"usedAgents"`
	UsedWorkflows         int64             `json:"usedWorkflows"`
	UsedKnowledgeBaseSize int64             `json:"usedKnowledgeBaseSize"`
	CreatedAt             string            `json:"createdAt"`
	UpdatedAt             string            `json:"updatedAt"`
}

// SubscriptionPlanConfigResponse 订阅计划配置响应
type SubscriptionPlanConfigResponse struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Plan        string            `json:"plan"`
	Price       int64             `json:"price"` // 月费价格，以分为单位
	Description string            `json:"description"`
	QuarterRate float64           `json:"quarterRate"`
	YearRate    float64           `json:"yearRate"`
	Configs     *model.PlanConfig `json:"configs"`
}

// WeChatPaymentOrderResponse 微信支付订单响应
type WeChatPaymentOrderResponse struct {
	OrderID   uuid.UUID `json:"orderId"`
	QrCodeURL string    `json:"qrCodeUrl"`
}

// PaymentStatusResponse 支付状态响应
type PaymentStatusResponse struct {
	Paid bool `json:"paid"`
}
