package knowledges

import (
	"model"

	"github.com/google/uuid"
)

type createKnowledgeBaseReq struct {
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	EmbeddingModelName     string   `json:"embeddingModelName"`
	EmbeddingModelProvider string   `json:"embeddingModelProvider"`
	ChatModelName          string   `json:"chatModelName"`
	ChatModelProvider      string   `json:"chatModelProvider"`
	Tags                   []string `json:"tags"`
}
type updateKnowledgeBaseReq struct {
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	EmbeddingModelName     string   `json:"embeddingModelName"`
	EmbeddingModelProvider string   `json:"embeddingModelProvider"`
	ChatModelName          string   `json:"chatModelName"`
	ChatModelProvider      string   `json:"chatModelProvider"`
	Tags                   []string `json:"tags"`
}
type listReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"size"`
	Search   string `json:"search"`
}
type searchReq struct {
	Params listReq `json:"params"`
}

type searchParams struct {
	Query string `json:"query"`
}
type listDocumentReq struct {
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
	Search    string `json:"search" form:"search"`
	SortBy    string `json:"sortBy" form:"sortBy"`
	Status    string `json:"status" form:"status"`
	SortOrder string `json:"sortOrder" form:"sortOrder"`
}

type CreateKnowledgeBaseAdminReq struct {
	Name                   string            `json:"name" binding:"required"`
	Description            string            `json:"description"`
	ChatModelName          string            `json:"chatModelName"`
	ChatModelProvider      string            `json:"chatModelProvider"`
	EmbeddingModelName     string            `json:"embeddingModelName"`
	EmbeddingModelProvider string            `json:"embeddingModelProvider"`
	StorageType            model.StorageType `json:"storageType"`
	Tags                   []string          `json:"tags"`
	CreatorID              uuid.UUID         `json:"creatorId" binding:"required"`
}

type UpdateKnowledgeBaseAdminReq struct {
	ID                     uuid.UUID                 `json:"id" binding:"required"`
	Name                   string                    `json:"name"`
	Description            string                    `json:"description"`
	ChatModelName          string                    `json:"chatModelName"`
	ChatModelProvider      string                    `json:"chatModelProvider"`
	EmbeddingModelName     string                    `json:"embeddingModelName"`
	EmbeddingModelProvider string                    `json:"embeddingModelProvider"`
	Tags                   []string                  `json:"tags"`
	Status                 model.KnowledgeBaseStatus `json:"status"`
}

type ListKnowledgeBasesAdminReq struct {
	Name      string `json:"name" form:"name"`
	Search    string `json:"search" form:"search"`
	CreatorID string `json:"creatorId" form:"creatorId"`
	Status    string `json:"status" form:"status"`
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"pageSize" form:"pageSize"`
}
