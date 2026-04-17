package agents

import (
	"model"
	"time"

	"github.com/google/uuid"
)

type ListAgentResponse struct {
	Agents []*model.Agent `json:"agents"`
	Total  int64          `json:"total"`
}

type chatSessionResponse struct {
	ID        uuid.UUID `json:"id"`
	AgentId   uuid.UUID `json:"agentId"`
	Title     string    `json:"title"`
	UserId    uuid.UUID `json:"userId"`
	CreatedAt int64     `json:"createdAt"`
	UpdatedAt int64     `json:"updatedAt"`
}

func toChatSessionResponse(session *model.ChatSession) *chatSessionResponse {
	return &chatSessionResponse{
		ID:        session.ID,
		AgentId:   session.AgentID,
		Title:     session.Title,
		UserId:    session.UserID,
		CreatedAt: session.CreatedAt.UnixMilli(),
		UpdatedAt: session.UpdatedAt.UnixMilli(),
	}
}

type chatMessageResponse struct {
	Id        uuid.UUID `json:"id"`
	SessionId uuid.UUID `json:"sessionId"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt int64     `json:"createdAt"`
}

func toChatMessageResponse(message *model.ChatMessage) *chatMessageResponse {
	return &chatMessageResponse{
		Id:        message.ID,
		SessionId: message.SessionID,
		Role:      message.Role,
		Content:   message.Content,
		CreatedAt: message.CreatedAt.UnixMilli(),
	}
}
func toChatMessageResponses(messages []*model.ChatMessage) []*chatMessageResponse {
	var res []*chatMessageResponse
	for _, message := range messages {
		res = append(res, toChatMessageResponse(message))
	}
	return res
}

type AgentDetailAdminResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	Icon               string     `json:"icon"`
	SystemPrompt       string     `json:"systemPrompt"`
	ModelProvider      string     `json:"modelProvider"`
	ModelName          string     `json:"modelName"`
	ModelParameters    model.JSON `json:"modelParameters"`
	OpeningDialogue    string     `json:"openingDialogue"`
	SuggestedQuestions model.JSON `json:"suggestedQuestions"`
	Version            uint       `json:"version"`
	Status             string     `json:"status"`
	Visibility         string     `json:"visibility"`
	InvocationCount    uint64     `json:"invocationCount"`
	CreatorID          uuid.UUID  `json:"creatorId"`
	CreatorName        string     `json:"creatorName"`
	CreatorEmail       string     `json:"creatorEmail"`
	CreatedAt          string     `json:"createdAt"`
	UpdatedAt          string     `json:"updatedAt"`
}

type AgentListAdminResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Icon            string    `json:"icon"`
	Status          string    `json:"status"`
	Visibility      string    `json:"visibility"`
	ModelProvider   string    `json:"modelProvider"`
	ModelName       string    `json:"modelName"`
	InvocationCount uint64    `json:"invocationCount"`
	CreatorID       uuid.UUID `json:"creatorId"`
	CreatorName     string    `json:"creatorName"`
	CreatorEmail    string    `json:"creatorEmail"`
	CreatedAt       string    `json:"createdAt"`
	UpdatedAt       string    `json:"updatedAt"`
}

type ListAgentsAdminResponse struct {
	List        []*AgentListAdminResponse `json:"list"`
	Total       int64                     `json:"total"`
	CurrentPage int                       `json:"currentPage"`
	PageSize    int                       `json:"pageSize"`
}

func toAgentDetailAdminResponse(agent *model.Agent, user *model.User) *AgentDetailAdminResponse {
	resp := &AgentDetailAdminResponse{
		ID:                 agent.ID,
		Name:               agent.Name,
		Description:        agent.Description,
		Icon:               agent.Icon,
		SystemPrompt:       agent.SystemPrompt,
		ModelProvider:      agent.ModelProvider,
		ModelName:          agent.ModelName,
		ModelParameters:    agent.ModelParameters,
		OpeningDialogue:    agent.OpeningDialogue,
		SuggestedQuestions: agent.SuggestedQuestions,
		Version:            agent.Version,
		Status:             string(agent.Status),
		Visibility:         string(agent.Visibility),
		InvocationCount:    agent.InvocationCount,
		CreatorID:          agent.CreatorID,
		CreatedAt:          agent.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          agent.UpdatedAt.Format(time.RFC3339),
	}
	if user != nil {
		resp.CreatorName = user.Username
		resp.CreatorEmail = user.Email
	}
	return resp
}

func toAgentListAdminResponse(agent *model.Agent, user *model.User) *AgentListAdminResponse {
	resp := &AgentListAdminResponse{
		ID:              agent.ID,
		Name:            agent.Name,
		Description:     agent.Description,
		Icon:            agent.Icon,
		Status:          string(agent.Status),
		Visibility:      string(agent.Visibility),
		ModelProvider:   agent.ModelProvider,
		ModelName:       agent.ModelName,
		InvocationCount: agent.InvocationCount,
		CreatorID:       agent.CreatorID,
		CreatedAt:       agent.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       agent.UpdatedAt.Format(time.RFC3339),
	}
	if user != nil {
		resp.CreatorName = user.Username
		resp.CreatorEmail = user.Email
	}
	return resp
}
