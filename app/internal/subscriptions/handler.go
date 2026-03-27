package subscriptions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mszlu521/thunder/req"
	"github.com/mszlu521/thunder/res"
)

type Handler struct {
	service *service
}

func NewHandler() *Handler {
	return &Handler{
		service: newService(),
	}
}

// GetUserSubscription 获取当前用户的订阅信息
// GET /api/v1/subscription/current
func (h *Handler) GetUserSubscription(c *gin.Context) {
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}

	resp, err := h.service.getCurrentSubscription(c.Request.Context(), userID)
	if err != nil {
		res.Error(c, err)
		return
	}

	// 如果用户没有订阅，返回null
	if resp == nil {
		res.Success(c, nil)
		return
	}

	res.Success(c, resp)
}

// GetSubscriptionPlans 获取所有订阅计划
// GET /api/v1/subscription/plans
func (h *Handler) GetSubscriptionPlans(c *gin.Context) {
	plans, err := h.service.getSubscriptionPlans(c.Request.Context())
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, plans)
}

// UpdateSubscription 创建或更新订阅
// POST /api/v1/subscription
func (h *Handler) UpdateSubscription(c *gin.Context) {
	var updateReq UpdateSubscriptionReq
	if err := req.JsonParam(c, &updateReq); err != nil {
		return
	}

	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}

	resp, err := h.service.updateSubscription(c.Request.Context(), userID, updateReq)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, resp)
}

// CancelSubscription 取消订阅
// DELETE /api/v1/subscription
func (h *Handler) CancelSubscription(c *gin.Context) {
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}

	err := h.service.cancelSubscription(c.Request.Context(), userID)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, true)
}

// CreateWeChatPaymentOrder 创建微信支付订单
// POST /api/v1/subscription/wechat/order
func (h *Handler) CreateWeChatPaymentOrder(c *gin.Context) {
	var orderReq CreateWeChatPaymentOrderReq
	if err := req.JsonParam(c, &orderReq); err != nil {
		return
	}

	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}

	resp, err := h.service.createWeChatPaymentOrder(c.Request.Context(), userID, orderReq)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, resp)
}

// CheckPaymentStatus 检查支付状态
// GET /api/v1/subscription/wechat/status/:orderId
func (h *Handler) CheckPaymentStatus(c *gin.Context) {
	var orderID uuid.UUID
	if err := req.Path(c, "orderId", &orderID); err != nil {
		return
	}

	resp, err := h.service.checkPaymentStatus(c.Request.Context(), orderID)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, resp)
}

// CreatePlanConfig 创建订阅计划配置
// POST /api/v1/subscription/admin/plan
func (h *Handler) CreatePlanConfig(c *gin.Context) {
	var createReq CreatePlanConfigReq
	if err := req.JsonParam(c, &createReq); err != nil {
		return
	}

	resp, err := h.service.createPlanConfig(c.Request.Context(), createReq)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, resp)
}

// UpdatePlanConfig 更新订阅计划配置
// PUT /api/v1/subscription/admin/plan
func (h *Handler) UpdatePlanConfig(c *gin.Context) {
	var updateReq UpdatePlanConfigReq
	if err := req.JsonParam(c, &updateReq); err != nil {
		return
	}

	resp, err := h.service.updatePlanConfig(c.Request.Context(), updateReq)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, resp)
}

// DeletePlanConfig 删除订阅计划配置
// DELETE /api/v1/subscription/admin/plan/:id
func (h *Handler) DeletePlanConfig(c *gin.Context) {
	var deleteReq DeletePlanConfigReq
	if err := req.Path(c, "id", &deleteReq.ID); err != nil {
		return
	}

	err := h.service.deletePlanConfig(c.Request.Context(), deleteReq.ID)
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, true)
}
