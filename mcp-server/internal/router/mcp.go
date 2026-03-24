package router

import (
	"mcp-server/internal/tool"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type McpRouter struct {
}

func (u *McpRouter) Register(engine *gin.Engine) {
	//需要两个 /sse， /message
	//这里需要使用mcp-go 创建mcp服务
	mcpServer := server.NewMCPServer(
		// 1. 创建 MCP 服务实例（支持工具、资源、提示词能力）
		"mszlu mcp server",
		mcp.LATEST_PROTOCOL_VERSION,
		server.WithToolCapabilities(true),           // 启用工具调用
		server.WithResourceCapabilities(true, true), // 启用资源
		server.WithPromptCapabilities(true),         // 启用提示词
	)
	// 2. 注册自定义工具（天气查询）
	weather := tool.NewWeatherTool("weather")
	mcpServer.AddTool(weather.Build(), weather.Invoke)
	// 3. 创建 SSE 服务器（AI 长连接通信）
	sseServer := server.NewSSEServer(
		mcpServer,
		server.WithBaseURL("http://localhost:7777"),
		server.WithSSEEndpoint("/sse"),         // SSE 连接端点
		server.WithMessageEndpoint("/message"), // 消息收发端点
		server.WithKeepAlive(true),             // 保持长连接
	)
	// 4.挂载到 Gin 路由
	engine.GET("/sse", gin.WrapH(sseServer.SSEHandler()))
	engine.POST("/message", gin.WrapH(sseServer.MessageHandler()))
}
