package subscriptions

import (
	"model"

	"github.com/google/uuid"
)

// UpdateSubscriptionReq 创建或更新订阅请求
type UpdateSubscriptionReq struct {
	PlanType model.SubscriptionPlan `json:"planType"`
	Duration model.PaymentDuration  `json:"duration"`
}

// CreateWeChatPaymentOrderReq 创建微信支付订单请求
type CreateWeChatPaymentOrderReq struct {
	PlanType model.SubscriptionPlan `json:"planType"`
	Duration model.PaymentDuration  `json:"duration"`
	Amount   int64                  `json:"amount"` // 以分为单位
}

// CheckPaymentStatusReq 检查支付状态请求（路径参数）
type CheckPaymentStatusReq struct {
	OrderID uuid.UUID `uri:"orderId" binding:"required"`
}

// CreatePlanConfigReq 创建订阅计划配置请求
type CreatePlanConfigReq struct {
	Name                 string                 `json:"name" binding:"required"`
	Plan                 model.SubscriptionPlan `json:"plan" binding:"required"`
	Price                int64                  `json:"price" binding:"required"`
	Description          string                 `json:"description"`
	QuarterRate          float64                `json:"quarterRate"`
	YearRate             float64                `json:"yearRate"`
	MaxAgents            int64                  `json:"maxAgents"`
	MaxWorkflows         int64                  `json:"maxWorkflows"`
	MaxKnowledgeBaseSize int64                  `json:"maxKnowledgeBaseSize"`
}

// UpdatePlanConfigReq 更新订阅计划配置请求
type UpdatePlanConfigReq struct {
	ID                   int64                  `json:"id" binding:"required"`
	Name                 string                 `json:"name"`
	Plan                 model.SubscriptionPlan `json:"plan"`
	Price                int64                  `json:"price"`
	Description          string                 `json:"description"`
	QuarterRate          float64                `json:"quarterRate"`
	YearRate             float64                `json:"yearRate"`
	MaxAgents            int64                  `json:"maxAgents"`
	MaxWorkflows         int64                  `json:"maxWorkflows"`
	MaxKnowledgeBaseSize int64                  `json:"maxKnowledgeBaseSize"`
}

// DeletePlanConfigReq 删除订阅计划配置请求
type DeletePlanConfigReq struct {
	ID int64 `json:"id" binding:"required"`
}
