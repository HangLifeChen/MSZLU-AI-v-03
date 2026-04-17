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

	// 管理员接口 - 订阅计划配置管理
	adminGroup := engine.Group("/api/v1/subscription/admin")
	{
		// 创建订阅计划配置
		adminGroup.POST("/plan", handler.CreatePlanConfig)
		// 更新订阅计划配置
		adminGroup.PUT("/plan", handler.UpdatePlanConfig)
		// 删除订阅计划配置
		adminGroup.DELETE("/plan/:id", handler.DeletePlanConfig)
		// 获取所有订阅计划
		adminGroup.GET("/plans", handler.GetSubscriptionPlans)
	}
}
