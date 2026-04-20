package router

import (
	"app/internal/admin"

	"github.com/gin-gonic/gin"
)

type AdminRouter struct{}

func (a *AdminRouter) Register(engine *gin.Engine) {
	adminGroup := engine.Group("/api/v1/admin/statistics")
	{
		handler := admin.NewHandler()
		adminGroup.GET("/users/growth", handler.GetUserGrowthTrend)
		adminGroup.GET("/knowledge-bases/growth", handler.GetKnowledgeBaseGrowthTrend)
		adminGroup.GET("/agents/growth", handler.GetAgentGrowthTrend)
	}
}

func (a *AdminRouter) Close() error {
	return nil
}
