package shared

import "github.com/google/uuid"

type UpdateCurrentBaseSubscriptionReq struct {
	UserId                uuid.UUID `json:"userId"`
	UsedAgents            int64     `json:"usedAgents" gorm:"default:0"`
	UsedWorkflows         int64     `json:"usedWorkflows" gorm:"default:0"`
	UsedKnowledgeBaseSize int64     `json:"usedKnowledgeBaseSize" gorm:"default:0"` // 以MB为单位
}
