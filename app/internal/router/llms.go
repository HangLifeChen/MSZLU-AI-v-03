package router

import (
	"app/internal/llms"

	"github.com/gin-gonic/gin"
)

type LLMRouter struct {
}

func (u *LLMRouter) Register(engine *gin.Engine) {
	llmGroup := engine.Group("/api/v1/provider-configs")
	{
		llmHandler := llms.NewHandler()
		llmGroup.POST("/", llmHandler.CreateProviderConfig)
		llmGroup.GET("/", llmHandler.ListProviderConfigs)
		llmGroup.GET("/:id", llmHandler.GetProviderConfig)
		llmGroup.PUT("/:id", llmHandler.UpdateProviderConfig)
		llmGroup.DELETE("/:id", llmHandler.DeleteProviderConfig)
	}
	llmsGroup := engine.Group("/api/v1/llms")
	{
		llmsHandler := llms.NewHandler()
		llmsGroup.POST("/", llmsHandler.CreateLLM)
		llmsGroup.GET("/", llmsHandler.ListLLMs)
		llmsGroup.GET("/all", llmsHandler.ListLLMAll)
		llmsGroup.GET("/:id", llmsHandler.GetLLM)
		llmsGroup.PUT("/:id", llmsHandler.UpdateLLM)
		llmsGroup.DELETE("/:id", llmsHandler.DeleteLLM)
	}
}
