package settings

import (
	"model"
)

// GetSettingsResponse 获取系统设置响应
type GetSettingsResponse struct {
	Settings *model.SystemSettings `json:"settings"`
}

// SaveSettingsResponse 保存系统设置响应
type SaveSettingsResponse struct {
	Settings *model.SystemSettings `json:"settings"`
}

// UpdateSettingsResponse 更新系统设置响应
type UpdateSettingsResponse struct {
	Settings *model.SystemSettings `json:"settings"`
}

// GetSettingsModuleResponse 获取特定模块设置响应
type GetSettingsModuleResponse struct {
	Data interface{} `json:"data"`
}

// UpdateSettingsModuleResponse 更新特定模块设置响应
type UpdateSettingsModuleResponse struct {
	Success bool `json:"success"`
}
