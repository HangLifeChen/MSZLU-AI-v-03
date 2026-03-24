package settings

import (
	"model"
)

// GetSettingsReq 获取系统设置请求
type GetSettingsReq struct{}

// SaveSettingsReq 保存系统设置请求
type SaveSettingsReq struct {
	Basic        model.BasicSettings        `json:"basic"`
	Model        model.ModelSettings        `json:"model"`
	Security     model.SecuritySettings     `json:"security"`
	Notification model.NotificationSettings `json:"notification"`
	Storage      model.StorageSettings      `json:"storage"`
}

// UpdateSettingsReq 更新系统设置请求
type UpdateSettingsReq struct {
	ID           string                     `json:"id"`
	Basic        model.BasicSettings        `json:"basic"`
	Model        model.ModelSettings        `json:"model"`
	Security     model.SecuritySettings     `json:"security"`
	Notification model.NotificationSettings `json:"notification"`
	Storage      model.StorageSettings      `json:"storage"`
}

// GetSettingsModuleReq 获取特定模块设置请求
type GetSettingsModuleReq struct {
	Module string `uri:"module" binding:"required"`
}

// UpdateSettingsModuleReq 更新特定模块设置请求
type UpdateSettingsModuleReq struct {
	Data map[string]interface{} `json:"data"`
}
