package router

import (
	"app/internal/subscriptions"

	"github.com/gin-gonic/gin"
)

type SubscriptionRouter struct {
}

func (u *SubscriptionRouter) Register(engine *gin.Engine) {
	handler := subscriptions.NewHandler()
	group := engine.Group("/api/v1/subscription")
	{
		// 获取当前用户的订阅信息
		group.GET("/current", handler.GetUserSubscription)
		// 获取所有订阅计划
		group.GET("/plans", handler.GetSubscriptionPlans)
		// 创建或更新订阅
		group.POST("", handler.UpdateSubscription)
		// 取消订阅
		group.DELETE("", handler.CancelSubscription)
		// 创建微信支付订单
		group.POST("/wechat/order", handler.CreateWeChatPaymentOrder)
		// 检查支付状态
		group.GET("/wechat/status/:orderId", handler.CheckPaymentStatus)
	}
}
