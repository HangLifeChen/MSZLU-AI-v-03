package settings

import (
	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/req"
	"github.com/mszlu521/thunder/res"
)

type Handler struct {
	service *service
}

// GetSettings 获取系统设置
func (h *Handler) GetSettings(c *gin.Context) {
	resp, err := h.service.getSettings(c.Request.Context())
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

// SaveSettings 保存系统设置
func (h *Handler) SaveSettings(c *gin.Context) {
	var saveReq SaveSettingsReq
	if err := req.JsonParam(c, &saveReq); err != nil {
		return
	}
	resp, err := h.service.saveSettings(c.Request.Context(), saveReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

// UpdateSettings 更新系统设置
func (h *Handler) UpdateSettings(c *gin.Context) {
	var updateReq UpdateSettingsReq
	if err := req.JsonParam(c, &updateReq); err != nil {
		return
	}
	resp, err := h.service.updateSettings(c.Request.Context(), updateReq)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

// GetSettingsModule 获取特定模块设置
func (h *Handler) GetSettingsModule(c *gin.Context) {
	var module string
	if err := req.Path(c, "module", &module); err != nil {
		return
	}
	resp, err := h.service.getSettingsModule(c.Request.Context(), module)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, resp)
}

// UpdateSettingsModule 更新特定模块设置
func (h *Handler) UpdateSettingsModule(c *gin.Context) {
	var module string
	if err := req.Path(c, "module", &module); err != nil {
		return
	}
	// 直接解析请求体为 map，支持直接发送数据而不需要包裹在 data 字段中
	var data map[string]interface{}
	if err := req.JsonParam(c, &data); err != nil {
		return
	}
	err := h.service.updateSettingsModule(c.Request.Context(), module, data)
	if err != nil {
		res.Error(c, err)
		return
	}
	res.Success(c, map[string]bool{"success": true})
}

func NewHandler() *Handler {
	return &Handler{
		service: newService(),
	}
}
