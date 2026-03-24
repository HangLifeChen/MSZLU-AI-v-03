package llms

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mszlu521/thunder/logs"
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
func (h *Handler) CreateProviderConfig(c *gin.Context) {
	var createReq CreateProviderConfigReq
	if err := req.JsonParam(c, &createReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	createResp, err := h.service.createProviderConfig(c.Request.Context(), userID, createReq)
	if err != nil {
		return
	}
	res.Success(c, createResp)
}

func (h *Handler) GetProviderConfig(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}
	_, exist := req.GetUserIdUUID(c)
	if !exist {
		return
	}
	resp, err := h.service.getProviderConfig(c.Request.Context(), id)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) ListProviderConfigs(c *gin.Context) {
	userId, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.listProviderConfigs(c.Request.Context(), userId)
	if err != nil {
		return
	}
	res.Success(c, resp)
}

func (h *Handler) DeleteProviderConfig(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}
	userId, exist := req.GetUserIdUUID(c)
	if !exist {
		return
	}

	response, err := h.service.DeleteProviderConfig(c.Request.Context(), userId, id)
	if err != nil {
		logs.Errorf("删除厂商配置失败: %v", err)
		res.Error(c, err)
		return
	}

	res.Success(c, response)
}

func (h *Handler) UpdateProviderConfig(c *gin.Context) {
	var id uuid.UUID
	if err := req.Path(c, "id", &id); err != nil {
		return
	}

	userId, exist := req.GetUserIdUUID(c)
	if !exist {
		return
	}

	var updateReq UpdateProviderConfigReq
	if err := req.JsonParam(c, &updateReq); err != nil {
		return
	}

	response, err := h.service.UpdateProviderConfig(c.Request.Context(), userId, id, updateReq)
	if err != nil {
		logs.Errorf("更新厂商配置失败: %v", err)
		res.Error(c, err)
		return
	}

	res.Success(c, response)
}

func (h *Handler) CreateLLM(c *gin.Context) {
	var createReq CreateLLMReq
	if err := req.JsonParam(c, &createReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	createResp, err := h.service.createLLM(c.Request.Context(), userID, createReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, createResp)
}

func (h *Handler) ListLLMs(c *gin.Context) {
	var listReq ListLLMsReq
	if err := req.QueryParam(c, &listReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.listLLMs(c.Request.Context(), userID, listReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) ListLLMAll(c *gin.Context) {
	var listReq ListLLMsReq
	if err := req.QueryParam(c, &listReq); err != nil {
		return
	}
	userID, ok := req.GetUserIdUUID(c)
	if !ok {
		return
	}
	resp, err := h.service.listLLMAll(c.Request.Context(), userID, listReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

func (h *Handler) GetLLM(c *gin.Context) {
	var id uuid.UUID
	err := req.Path(c, "id", &id)
	if err != nil {
		return
	}
	userId, exist := req.GetUserIdUUID(c)
	if !exist {
		return
	}
	response, err := h.service.GetLLM(c.Request.Context(), userId, id)
	if err != nil {
		logs.Errorf("获取模型失败: %v", err)
		res.Error(c, err)
		return
	}

	res.Success(c, response)
}

func (h *Handler) UpdateLLM(c *gin.Context) {
	var id uuid.UUID
	err := req.Path(c, "id", &id)
	if err != nil {
		return
	}
	userId, exist := req.GetUserIdUUID(c)
	if !exist {
		return
	}

	var updateReq UpdateLLMReq
	if err := req.JsonParam(c, &updateReq); err != nil {
		return
	}
	response, err := h.service.UpdateLLM(c.Request.Context(), userId, id, updateReq)
	if err != nil {
		logs.Errorf("更新模型失败: %v", err)
		res.Error(c, err)
		return
	}

	res.Success(c, response)
}

func (h *Handler) DeleteLLM(c *gin.Context) {
	var id uuid.UUID
	err := req.Path(c, "id", &id)
	if err != nil {
		logs.Errorf("参数错误: %v", err)
		res.Error(c, err)
		return
	}

	userID, exists := req.GetUserIdUUID(c)
	if !exists {
		return
	}

	response, err := h.service.DeleteLLM(c.Request.Context(), userID, id)
	if err != nil {
		logs.Errorf("删除模型失败: %v", err)
		res.Error(c, err)
		return
	}

	res.Success(c, response)
}
