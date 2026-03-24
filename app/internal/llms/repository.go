package llms

import (
	"context"
	"model"

	"github.com/google/uuid"
)

type repository interface {
	createProviderConfig(ctx context.Context, m *model.ProviderConfig) error
	listProviderConfigs(ctx context.Context, userId uuid.UUID) ([]*model.ProviderConfig, int64, error)
	createLLM(ctx context.Context, llm *model.LLM) error
	listLLMs(ctx context.Context, userID uuid.UUID, filter LLMFilter) ([]*model.LLM, int64, error)
	getProviderConfig(ctx context.Context, provider string) (*model.ProviderConfig, error)
	getProviderConfigByID(ctx context.Context, id uuid.UUID) (*model.ProviderConfig, error)
	deleteProviderConfig(ctx context.Context, id uuid.UUID) error
	updateProviderConfig(ctx context.Context, config *model.ProviderConfig) error
	listLLMAll(ctx context.Context, userID uuid.UUID, filter LLMFilter) ([]*model.LLM, error)
	getByModelID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*model.LLM, error)
	update(ctx context.Context, llm *model.LLM) error
	delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	getByModelName(ctx context.Context, userId uuid.UUID, modelName string, modelType model.LLMType) (*model.LLM, error)
}
