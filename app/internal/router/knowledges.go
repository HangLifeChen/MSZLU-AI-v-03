package router

import (
	"app/internal/knowledges"

	"github.com/gin-gonic/gin"
)

type KnowledgeBaseRouter struct {
	handler *knowledges.Handler
}

func (u *KnowledgeBaseRouter) Register(engine *gin.Engine) {
	knowledgesGroup := engine.Group("/api/v1/knowledge")
	{
		knowledgesHandler := knowledges.NewHandler()
		u.handler = knowledgesHandler
		knowledgesGroup.POST("/", knowledgesHandler.CreateKnowledgeBase)
		knowledgesGroup.POST("/list", knowledgesHandler.ListKnowledgeBases)
		knowledgesGroup.GET("/:id", knowledgesHandler.GetKnowledgeBase)
		knowledgesGroup.PUT("/:id", knowledgesHandler.UpdateKnowledgeBase)
		knowledgesGroup.POST("/:id/search", knowledgesHandler.SearchKnowledgeBase)
		knowledgesGroup.DELETE("/:id", knowledgesHandler.DeleteKnowledgeBase)
		knowledgesGroup.GET("/:id/documents", knowledgesHandler.ListDocuments)
		knowledgesGroup.POST("/:id/documents", knowledgesHandler.UploadDocuments)
		knowledgesGroup.DELETE("/:id/documents/:documentId", knowledgesHandler.DeleteDocuments)
		//GET /api/v1/knowledge/{id}/documents/{documentId}/content
		knowledgesGroup.GET("/:id/documents/:documentId/content", knowledgesHandler.GetDocumentContent)
		//POST /api/v1/knowledge/{id}/documents/{documentId}/reindex
		knowledgesGroup.POST("/:id/documents/:documentId/reindex", knowledgesHandler.ReindexDocument)
	}
}

func (u *KnowledgeBaseRouter) Close() error {
	if u.handler != nil {
		return u.handler.Close()
	}
	return nil
}
