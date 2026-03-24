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
