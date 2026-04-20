package knowledges

import (
	"model"
	"time"

	"github.com/google/uuid"
)

type ListResp struct {
	KnowledgeBases []*model.KnowledgeBase `json:"knowledgeBases"`
	Total          int64                  `json:"total"`
}

type KnowledgeBaseResponse struct {
	Id                     uuid.UUID         `json:"id"`
	Name                   string            `json:"name"`
	Tags                   []string          `json:"tags"`
	Description            string            `json:"description"`
	EmbeddingModelName     string            `json:"embeddingModelName"`
	EmbeddingModelProvider string            `json:"embeddingModelProvider"`
	ChatModelName          string            `json:"chatModelName"`
	ChatModelProvider      string            `json:"chatModelProvider"`
	StorageType            model.StorageType `json:"storageType"`
	StorageConfig          model.JSON        `json:"storageConfig"`
	DocumentCount          int               `json:"documentCount"`
	TotalSize              int64             `json:"totalSize"`
	CreatedAt              int64             `json:"createdAt"`
	UpdatedAt              int64             `json:"updatedAt"`
	CreatorId              uuid.UUID         `json:"creatorId"`
}

type ListDocumentsResp struct {
	Documents []*model.Document `json:"items"`
	Total     int64             `json:"total"`
}

type SearchResponse struct {
	Query   string          `json:"query"`
	Results []*SearchResult `json:"results"`
	Total   int64           `json:"total"`
	Took    int64           `json:"took"` //耗时
	KbId    uuid.UUID       `json:"kbId"`
}

type SearchResult struct {
	Id         uuid.UUID       `json:"id"`
	DocumentId uuid.UUID       `json:"documentId"`
	Content    string          `json:"content"`
	Score      float64         `json:"score"`
	Metadata   model.JSON      `json:"metadata"`
	Position   int             `json:"position"`
	Document   *model.Document `json:"document"`
}

type DocumentContent struct {
	Id             uuid.UUID            `json:"id"`
	DocId          uuid.UUID            `json:"docId"`
	KbId           uuid.UUID            `json:"kbId"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	FileType       string               `json:"fileType"`
	FileSize       int64                `json:"fileSize"`
	Status         model.DocumentStatus `json:"status"`
	EmbeddingModel string               `json:"embeddingModel"`
	ChunkCount     int64                `json:"chunkCount"`
	ProcessedCount int64                `json:"processedCount"`
	Metadata       model.JSON           `json:"metadata"`
	CreatedAt      int64                `json:"createdAt"`
	UpdatedAt      int64                `json:"updatedAt"`
}

type KnowledgeBaseDetailAdminResponse struct {
	Id                     uuid.UUID         `json:"id"`
	Name                   string            `json:"name"`
	Description            string            `json:"description"`
	EmbeddingModelName     string            `json:"embeddingModelName"`
	EmbeddingModelProvider string            `json:"embeddingModelProvider"`
	ChatModelName          string            `json:"chatModelName"`
	ChatModelProvider      string            `json:"chatModelProvider"`
	StorageType            model.StorageType `json:"storageType"`
	Tags                   []string          `json:"tags"`
	Status                 string            `json:"status"`
	DocumentCount          int               `json:"documentCount"`
	TotalSize              int64             `json:"totalSize"`
	CreatorID              uuid.UUID         `json:"creatorId"`
	CreatorName            string            `json:"creatorName"`
	CreatorEmail           string            `json:"creatorEmail"`
	CreatedAt              string            `json:"createdAt"`
	UpdatedAt              string            `json:"updatedAt"`
}

type KnowledgeBaseListAdminResponse struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StorageType   string    `json:"storageType"`
	Status        string    `json:"status"`
	DocumentCount int       `json:"documentCount"`
	TotalSize     int64     `json:"totalSize"`
	CreatorID     uuid.UUID `json:"creatorId"`
	CreatorName   string    `json:"creatorName"`
	CreatorEmail  string    `json:"creatorEmail"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
}

type ListKnowledgeBasesAdminResponse struct {
	List        []*KnowledgeBaseListAdminResponse `json:"list"`
	Total       int64                             `json:"total"`
	CurrentPage int                               `json:"currentPage"`
	PageSize    int                               `json:"pageSize"`
}

func toKnowledgeBaseDetailAdminResponse(kb *model.KnowledgeBase, user *model.User, totalSize int64, docCount int64) *KnowledgeBaseDetailAdminResponse {
	resp := &KnowledgeBaseDetailAdminResponse{
		Id:                     kb.ID,
		Name:                   kb.Name,
		Description:            kb.Description,
		EmbeddingModelName:     kb.EmbeddingModelName,
		EmbeddingModelProvider: kb.EmbeddingModelProvider,
		ChatModelName:          kb.ChatModelName,
		ChatModelProvider:      kb.ChatModelProvider,
		StorageType:            kb.StorageType,
		Tags:                   kb.Tags,
		Status:                 string(kb.Status),
		DocumentCount:          int(docCount),
		TotalSize:              totalSize,
		CreatorID:              kb.CreatorID,
		CreatedAt:              kb.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              kb.UpdatedAt.Format(time.RFC3339),
	}
	if user != nil {
		resp.CreatorName = user.Username
		resp.CreatorEmail = user.Email
	}
	return resp
}

func toKnowledgeBaseListAdminResponse(kb *model.KnowledgeBase, user *model.User) *KnowledgeBaseListAdminResponse {
	resp := &KnowledgeBaseListAdminResponse{
		Id:            kb.ID,
		Name:          kb.Name,
		Description:   kb.Description,
		StorageType:   string(kb.StorageType),
		Status:        string(kb.Status),
		DocumentCount: int(kb.DocumentCount),
		CreatorID:     kb.CreatorID,
		CreatedAt:     kb.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     kb.UpdatedAt.Format(time.RFC3339),
	}
	if user != nil {
		resp.CreatorName = user.Username
		resp.CreatorEmail = user.Email
	}
	return resp
}
