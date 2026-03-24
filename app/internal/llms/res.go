package llms

import "model"

type ListProviderConfigsResponse struct {
	Total           int64                   `json:"total"`
	ProviderConfigs []*model.ProviderConfig `json:"providerConfigs"`
}

type ListLLMsResponse struct {
	Total int64        `json:"total"`
	LLMs  []*model.LLM `json:"llms"`
}

type UpdateProviderConfigResponse struct {
	Success bool `json:"success"`
}

type DeleteProviderConfigResponse struct {
	Success bool `json:"success"`
}

type GetLLMResponse struct {
	LLM *model.LLM `json:"llm"`
}

type UpdateLLMResponse struct {
	Success bool `json:"success"`
}

type DeleteLLMResponse struct {
	Success bool `json:"success"`
}
