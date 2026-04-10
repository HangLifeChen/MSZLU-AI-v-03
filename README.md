# XY-AI

XY-AI 是一个基于 Go 语言构建的 AI 智能体平台，集成了多模型管理、知识库（RAG）、工作流编排、MCP 工具协议和 A2A（Agent-to-Agent）通信等核心能力，支持快速构建和部署企业级 AI 应用。

## 架构概览

```
XY-AI/
├── app/            # 主应用服务（API 层）
├── core/           # 核心引擎（AI 工作流、知识库、工具）
├── model/          # 数据模型层
├── common/         # 公共工具库
├── a2a-server/     # A2A 智能体通信服务
├── mcp-server/     # MCP 工具协议服务
├── deployment/     # 部署配置
├── k8s/            # Kubernetes 编排文件
├── docs/           # 项目文档
└── skills/         # 技能插件
```

项目使用 [Go Workspace](https://go.dev/doc/modules/workspaces)（`go.work`）管理多模块，基于 Go 1.25.0 开发。

## 核心功能

### 智能体（Agent）管理
- 创建、配置和发布 AI 智能体
- 自定义系统提示词、模型参数（Temperature、TopP、MaxTokens 等）
- 支持开场白、推荐问题、版本管理和可见性控制（私有/公开/仅链接）
- 智能体可绑定工具、知识库和工作流
- 聊天会话与消息持久化

### 多模型支持（LLM）
- **对话模型**：OpenAI、DeepSeek、Qwen、Ollama、Gemini、Claude、千帆、混元、方舟（ARK）
- **向量模型**：OpenAI、Ollama、DashScope（通义）、方舟（ARK）
- **视觉模型**：Qwen-VL 等
- 自定义厂商配置（API Key、Base URL）
- 模型参数精细调控

### 知识库（RAG）
- 支持文档上传与解析（PDF、DOCX、Markdown、HTML、EPUB）
- 文档自动切片与向量化
- 向量存储支持 **Elasticsearch** 和 **Milvus** 双引擎
- 基于语义的文档检索
- Token 消耗统计与文档状态管理

### 工作流编排
- 可视化节点式工作流（基于 [CloudWeGo Eino](https://github.com/cloudwego/eino) 的 `compose.Workflow`）
- 内置节点类型：
  - `textDisplay` - 文本展示
  - `textCombine` - 文本合并
  - `htmlDisplay` - HTML 展示
  - `qwenVL` - 视觉语言节点
- 支持自定义节点扩展（`NodeFactory` 注册机制）
- 工作流版本管理与状态控制

### 工具系统
- **系统工具**：天气查询、Git 操作、Kubernetes 资源管理（查询/日志/操作/健康检查）
- **MCP 工具**：支持通过 MCP 协议动态接入外部工具
- 工具可绑定到智能体，由智能体在对话中自动调用

### A2A 通信
- 基于 [CloudWeGo Eino A2A](https://github.com/cloudwego/eino-ext) 的 Agent-to-Agent 协议
- Agent 市场：注册和发现远程智能体服务
- 支持 JSON-RPC 传输

### 系统管理
- 用户认证与权限管理（JWT）
- 系统设置（基本设置、模型设置、安全设置、通知设置、存储设置）
- 订阅计划（免费版/基础版/高级版/企业版）
- 云存储集成（阿里云 OSS、七牛云）
- 文件上传管理

## 技术栈

| 类别 | 技术 |
|------|------|
| 语言 | Go 1.25.0 |
| Web 框架 | Gin（app）、Hertz（a2a-server） |
| AI 框架 | CloudWeGo Eino |
| 数据库 | PostgreSQL（GORM） |
| 缓存 | Redis |
| 向量数据库 | Elasticsearch 8、Milvus |
| 容器编排 | Kubernetes |
| 协议 | MCP、A2A、SSE |
| 日志 | Thunder Logs |
| 配置 | Thunder Config + YAML |

## 快速开始

### 环境要求

- Go 1.25.0+
- PostgreSQL
- Redis
- Elasticsearch 8 或 Milvus（知识库功能）

### 配置

在 `app/etc/` 和 `mcp-server/etc/` 目录下创建 `config.yml` 配置文件，配置数据库、Redis、AI 模型 API Key 等信息。

### 构建 & 运行

```bash
# 主应用服务
go run app/main.go

# MCP 工具服务
go run mcp-server/main.go

# A2A 智能体通信服务
go run a2a-server/main.go
```

### Kubernetes 部署

项目提供了完整的 Kubernetes 部署文件：

```bash
kubectl apply -f k8s/
```

包含：Deployment、Service、Ingress、HPA、ConfigMap、PVC、RBAC 等资源。

## 项目结构详情

### app 模块
主应用服务，包含各业务模块的 Handler/Service/Repository 分层实现：
- `internal/agents` - 智能体管理
- `internal/llms` - 大模型管理
- `internal/knowledges` - 知识库管理
- `internal/tools` - 工具管理
- `internal/workflows` - 工作流管理
- `internal/a2a` - A2A 智能体市场
- `internal/subscriptions` - 订阅管理
- `internal/settings` - 系统设置
- `internal/nodes` - 节点管理
- `internal/employees` - 员工管理
- `internal/router` - 路由注册

### core 模块
核心引擎，提供 AI 能力的底层实现：
- `ai/` - 工作流执行器、模板引擎、消息处理
- `ai/nodes/` - 工作流节点实现
- `ai/kbs/` - 知识库向量存储（ES、Milvus）
- `ai/mcps/` - MCP 协议集成
- `ai/tools/` - 系统工具实现
- `ai/interview/` - 面试相关 AI 功能
- `upload/` - 文件上传

### model 模块
共享数据模型，定义所有数据库实体及其关联关系。

### common 模块
公共工具库，包含时间处理、Token 计算（tiktoken）、Markdown 解析、EPUB 解析等。

## 许可证

