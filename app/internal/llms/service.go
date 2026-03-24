package llms

import (
	"context"
	"model"
	"time"

	"github.com/google/uuid"
	"github.com/mszlu521/thunder/database"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
)

type service struct {
	repo repository
}

func (s *service) createProviderConfig(ctx context.Context, userID uuid.UUID, req CreateProviderConfigReq) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	config := model.ProviderConfig{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Provider:    req.Provider,
		Status:      req.Status,
		APIKey:      req.APIKey,
		APIBase:     req.APIBase,
	}
	err := s.repo.createProviderConfig(ctx, &config)
	if err != nil {
		logs.Errorf("create provider config error: %v", err)
		return nil, errs.DBError
	}
	return config, nil
}

func (s *service) getProviderConfig(ctx context.Context, id uuid.UUID) (*model.ProviderConfig, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	return s.repo.getProviderConfigByID(ctx, id)
}

func (s *service) listProviderConfigs(ctx context.Context, userId uuid.UUID) (*ListProviderConfigsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	list, total, err := s.repo.listProviderConfigs(ctx, userId)
	if err != nil {
		logs.Errorf("list provider configs error: %v", err)
		return nil, errs.DBError
	}
	return &ListProviderConfigsResponse{
		ProviderConfigs: list,
		Total:           total,
	}, nil
}

func (s *service) DeleteProviderConfig(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*DeleteProviderConfigResponse, error) {
	// 先检查配置是否存在且属于该用户
	config, err := s.repo.getProviderConfigByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if config.UserID != userID {
		return nil, nil // 或者返回错误
	}

	err = s.repo.deleteProviderConfig(ctx, id)
	if err != nil {
		return nil, err
	}

	return &DeleteProviderConfigResponse{Success: true}, nil
}

func (s *service) UpdateProviderConfig(ctx context.Context, userId uuid.UUID, id uuid.UUID, req UpdateProviderConfigReq) (*UpdateProviderConfigResponse, error) {
	config, err := s.repo.getProviderConfigByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 检查是否属于该用户
	if config.UserID != userId {
		return nil, nil // 或者返回错误
	}

	// 更新配置信息
	if req.Provider != "" {
		config.Provider = req.Provider
	}

	if req.Description != "" {
		config.Description = req.Description
	}

	if req.APIKey != "" {
		config.APIKey = req.APIKey
	}

	if req.APIBase != "" {
		config.APIBase = req.APIBase
	}

	if req.Status != "" {
		switch req.Status {
		case "active":
			config.Status = model.LLMStatusActive
		case "inactive":
			config.Status = model.LLMStatusInactive
		}
	}

	err = s.repo.updateProviderConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &UpdateProviderConfigResponse{Success: true}, nil
}

func (s *service) createLLM(ctx context.Context, userID uuid.UUID, req CreateLLMReq) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	llm := model.LLM{
		BaseModel: model.BaseModel{
			ID: uuid.New(),
		},
		UserID:           userID,
		Name:             req.Name,
		Description:      req.Description,
		ProviderConfigID: req.ProviderConfigID,
		ModelName:        req.ModelName,
		ModelType:        req.ModelType,
		Config:           req.Config,
		Status:           req.Status,
	}
	err := s.repo.createLLM(ctx, &llm)
	if err != nil {
		logs.Errorf("create llm error: %v", err)
		return nil, errs.DBError
	}
	return llm, nil
}

func (s *service) listLLMs(ctx context.Context, userID uuid.UUID, req ListLLMsReq) (*ListLLMsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	filter := LLMFilter{
		ModelType: req.ModelType,
	}
	list, total, err := s.repo.listLLMs(ctx, userID, filter)
	if err != nil {
		logs.Errorf("list llms error: %v", err)
		return nil, errs.DBError
	}
	return &ListLLMsResponse{
		LLMs:  list,
		Total: total,
	}, nil
}

func (s *service) listLLMAll(ctx context.Context, userID uuid.UUID, req ListLLMsReq) ([]*model.LLM, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	filter := LLMFilter{
		ModelType: req.ModelType,
	}
	list, err := s.repo.listLLMAll(ctx, userID, filter)
	if err != nil {
		logs.Errorf("list llms error: %v", err)
		return nil, errs.DBError
	}
	return list, nil
}
func (s *service) GetLLM(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*GetLLMResponse, error) {
	llm, err := s.repo.getByModelID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	return &GetLLMResponse{LLM: llm}, nil
}

// UpdateLLM 更新用户自定义模型
func (s *service) UpdateLLM(ctx context.Context, userID uuid.UUID, id uuid.UUID, req UpdateLLMReq) (*UpdateLLMResponse, error) {
	llm, err := s.repo.getByModelID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	// 更新模型信息
	if req.Name != "" {
		llm.Name = req.Name
	}

	if req.Description != "" {
		llm.Description = req.Description
	}

	if req.ProviderConfigID != uuid.Nil {
		llm.ProviderConfigID = req.ProviderConfigID
	}

	if req.ModelName != "" {
		llm.ModelName = req.ModelName
	}

	if req.ModelType != "" {
		switch req.ModelType {
		case "chat":
			llm.ModelType = model.LLMTypeChat
		case "embedding":
			llm.ModelType = model.LLMTypeEmbedding
		}
	}

	if req.Config != nil {
		if req.Config.MaxTokens != 0 {
			llm.Config.MaxTokens = req.Config.MaxTokens
		}
		if req.Config.Temperature != 0 {
			llm.Config.Temperature = req.Config.Temperature
		}
		if req.Config.TopP != 0 {
			llm.Config.TopP = req.Config.TopP
		}
	}

	if req.Status != "" {
		switch req.Status {
		case "active":
			llm.Status = model.LLMStatusActive
		case "inactive":
			llm.Status = model.LLMStatusInactive
		}
	}

	err = s.repo.update(ctx, llm)
	if err != nil {
		return nil, err
	}

	return &UpdateLLMResponse{Success: true}, nil
}

// DeleteLLM 删除用户自定义模型
func (s *service) DeleteLLM(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*DeleteLLMResponse, error) {
	err := s.repo.delete(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	return &DeleteLLMResponse{Success: true}, nil
}
func newService() *service {
	return &service{
		repo: newModels(database.GetPostgresDB().GormDB),
	}
}
