package llms

import (
	"model"

	"github.com/google/uuid"
)

type CreateProviderConfigReq struct {
	Name        string          `json:"name"`
	Provider    string          `json:"provider"`
	Description string          `json:"description"`
	APIKey      string          `json:"apiKey"`
	APIBase     string          `json:"apiBase"`
	Status      model.LLMStatus `json:"status"`
}

type CreateLLMReq struct {
	Name             string          `json:"name"`
	Description      string          `json:"description"`
	ProviderConfigID uuid.UUID       `json:"providerConfigId"`
	ModelName        string          `json:"modelName"`
	ModelType        model.LLMType   `json:"modelType"`
	Config           model.LLMConfig `json:"config"`
	Status           model.LLMStatus `json:"status"`
}

type ListLLMsReq struct {
	ModelID   string        `json:"modelId" form:"modelId"`
	Name      string        `json:"name" form:"name"`
	Provider  string        `json:"provider" form:"provider"`
	ModelType model.LLMType `json:"modelType"`
	Status    string        `json:"status" form:"status"`
	Limit     int           `json:"limit" form:"limit"`
	Offset    int           `json:"offset" form:"offset"`
}

type UpdateProviderConfigReq struct {
	Provider    string `json:"provider"`    // 提供商名称
	Description string `json:"description"` // 描述
	APIKey      string `json:"apiKey"`      // API密钥
	APIBase     string `json:"apiBase"`     // API地址
	Status      string `json:"status"`      // 状态
}
type UpdateLLMReq struct {
	Name             string           `json:"name"`             // 模型名称
	Description      string           `json:"description"`      // 描述
	ProviderConfigID uuid.UUID        `json:"providerConfigId"` // 关联的厂商配置ID
	ModelName        string           `json:"modelName"`        // 模型标识
	ModelType        string           `json:"modelType"`        // 模型类型
	Config           *model.LLMConfig `json:"config"`           // 其他关键配置
	Status           string           `json:"status"`           // 状态
}
