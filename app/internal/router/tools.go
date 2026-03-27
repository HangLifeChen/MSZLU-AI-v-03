package router

import (
	"app/internal/tools"

	"github.com/gin-gonic/gin"
)

type ToolRouter struct {
}

func (t *ToolRouter) Register(r *gin.Engine) {
	toolsGroup := r.Group("/api/v1/tools")
	{
		toolsHandler := tools.NewHandler()
		toolsGroup.POST("/", toolsHandler.CreateTool)
		toolsGroup.GET("/", toolsHandler.ListTools)
		toolsGroup.GET("/:id", toolsHandler.GetTool)
		toolsGroup.PUT("/:id", toolsHandler.UpdateTool)
		toolsGroup.DELETE("/:id", toolsHandler.DeleteTool)
		toolsGroup.POST("/:id/test", toolsHandler.TestTool)
		toolsGroup.GET("/mcp/:mcpId/tools", toolsHandler.GetMcpTools)
	}

	// 管理员路由
	adminToolsGroup := r.Group("/api/v1/admin/tools")
	{
		adminToolsHandler := tools.NewHandler()
		adminToolsGroup.POST("/", adminToolsHandler.CreateToolAdmin)
		adminToolsGroup.GET("/", adminToolsHandler.ListToolsAdmin)
		adminToolsGroup.GET("/stats", adminToolsHandler.GetToolStats)
		adminToolsGroup.GET("/:id", adminToolsHandler.GetToolAdmin)
		adminToolsGroup.PUT("/:id", adminToolsHandler.UpdateToolAdmin)
		adminToolsGroup.DELETE("/:id", adminToolsHandler.DeleteToolAdmin)
	}
}
