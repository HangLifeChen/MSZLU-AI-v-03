# 管理后台 - 智能体管理 API

Base URL: `/api/v1/agent/admin`

## 统一响应格式

```json
{
  "code": 200,
  "msg": "success",
  "data": {}
}
```

---

## 1. 创建智能体

**POST** `/api/v1/agent/admin/`

### 请求体 (JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 智能体名称 |
| `description` | string | 否 | 描述信息 |
| `status` | string | 否 | 状态，可选值：`"draft"`、`"published"`、`"archived"`，默认 `"draft"` |
| `systemPrompt` | string | 否 | 系统提示词 |
| `modelProvider` | string | 否 | 模型提供商，如 `"openai"`、`"ollama"`、`"qwen"`、`"deepseek"` |
| `modelName` | string | 否 | 模型名称，如 `"gpt-4o"`、`"qwen-plus"` |
| `modelParameters` | object | 否 | 模型参数配置 |
| `openingDialogue` | string | 否 | 开场白 |
| `creatorId` | string(uuid) | 是 | 所属用户ID |

### modelParameters 结构

| 字段 | 类型 | 说明 |
|------|------|------|
| `maxTokens` | int | 最大生成长度 |
| `temperature` | float | 随机性控制 0.0~2.0 |
| `topP` | float | 核采样阈值 |
| `n` | int | 生成数量 |
| `stop` | string[] | 停止词列表 |
| `presencePenalty` | float | 话题新鲜度惩罚 |
| `frequencyPenalty` | float | 重复度惩罚 |

### 请求示例

```json
{
  "name": "客服助手",
  "description": "智能客服问答智能体",
  "status": "draft",
  "systemPrompt": "你是一个专业的客服助手，请礼貌地回答用户的问题。",
  "modelProvider": "openai",
  "modelName": "gpt-4o",
  "modelParameters": {
    "maxTokens": 4096,
    "temperature": 0.7,
    "topP": 1.0
  },
  "openingDialogue": "你好！我是客服助手，有什么可以帮你的吗？",
  "creatorId": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 响应 data (AgentDetailAdminResponse)

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string(uuid) | 智能体ID |
| `name` | string | 智能体名称 |
| `description` | string | 描述信息 |
| `icon` | string | 图标 |
| `systemPrompt` | string | 系统提示词 |
| `modelProvider` | string | 模型提供商 |
| `modelName` | string | 模型名称 |
| `modelParameters` | object | 模型参数 |
| `openingDialogue` | string | 开场白 |
| `suggestedQuestions` | object | 建议问题列表 |
| `version` | uint | 版本号 |
| `status` | string | 状态 |
| `visibility` | string | 可见性 |
| `invocationCount` | uint64 | 调用次数 |
| `creatorId` | string(uuid) | 创建者ID |
| `creatorName` | string | 创建者用户名 |
| `creatorEmail` | string | 创建者邮箱 |
| `createdAt` | string | 创建时间 (RFC3339) |
| `updatedAt` | string | 更新时间 (RFC3339) |

### 响应示例

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "客服助手",
    "description": "智能客服问答智能体",
    "icon": "",
    "systemPrompt": "你是一个专业的客服助手，请礼貌地回答用户的问题。",
    "modelProvider": "openai",
    "modelName": "gpt-4o",
    "modelParameters": {
      "maxTokens": 4096,
      "temperature": 0.7,
      "topP": 1.0
    },
    "openingDialogue": "你好！我是客服助手，有什么可以帮你的吗？",
    "suggestedQuestions": {},
    "version": 1,
    "status": "draft",
    "visibility": "private",
    "invocationCount": 0,
    "creatorId": "550e8400-e29b-41d4-a716-446655440000",
    "creatorName": "张三",
    "creatorEmail": "zhangsan@example.com",
    "createdAt": "2026-04-17T09:00:00Z",
    "updatedAt": "2026-04-17T09:00:00Z"
  }
}
```

---

## 2. 查询智能体列表

**GET** `/api/v1/agent/admin/`

### Query 参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 否 | 按智能体名称模糊搜索 |
| `status` | string | 否 | 按状态筛选：`"draft"`、`"published"`、`"archived"` |
| `creatorId` | string(uuid) | 否 | 按创建者ID筛选 |
| `page` | int | 否 | 页码，默认 1 |
| `pageSize` | int | 否 | 每页数量，默认 10 |

### 请求示例

```
GET /api/v1/agent/admin/?name=助手&status=published&page=1&pageSize=10
```

### 响应 data (ListAgentsAdminResponse)

| 字段 | 类型 | 说明 |
|------|------|------|
| `list` | AgentListAdminResponse[] | 智能体列表 |
| `total` | int64 | 总数 |
| `currentPage` | int | 当前页码 |
| `pageSize` | int | 每页数量 |

#### AgentListAdminResponse

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string(uuid) | 智能体ID |
| `name` | string | 智能体名称 |
| `description` | string | 描述信息 |
| `icon` | string | 图标 |
| `status` | string | 状态 |
| `visibility` | string | 可见性 |
| `modelProvider` | string | 模型提供商 |
| `modelName` | string | 模型名称 |
| `invocationCount` | uint64 | 调用次数 |
| `creatorId` | string(uuid) | 创建者ID |
| `creatorName` | string | 创建者用户名 |
| `creatorEmail` | string | 创建者邮箱 |
| `createdAt` | string | 创建时间 (RFC3339) |
| `updatedAt` | string | 更新时间 (RFC3339) |

### 响应示例

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "客服助手",
        "description": "智能客服问答智能体",
        "icon": "",
        "status": "published",
        "visibility": "public",
        "modelProvider": "openai",
        "modelName": "gpt-4o",
        "invocationCount": 128,
        "creatorId": "550e8400-e29b-41d4-a716-446655440000",
        "creatorName": "张三",
        "creatorEmail": "zhangsan@example.com",
        "createdAt": "2026-04-17T09:00:00Z",
        "updatedAt": "2026-04-17T10:30:00Z"
      }
    ],
    "total": 1,
    "currentPage": 1,
    "pageSize": 10
  }
}
```

---

## 3. 获取智能体详情

**GET** `/api/v1/agent/admin/:id`

### 路径参数

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string(uuid) | 智能体ID |

### 请求示例

```
GET /api/v1/agent/admin/660e8400-e29b-41d4-a716-446655440001
```

### 响应 data

同创建接口的 `AgentDetailAdminResponse`，返回完整详情信息。

---

## 4. 更新智能体

**PUT** `/api/v1/agent/admin/:id`

### 路径参数

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string(uuid) | 智能体ID |

### 请求体 (JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 否 | 智能体名称（非空时更新） |
| `description` | string | 否 | 描述信息（非空时更新） |
| `status` | string | 否 | 状态（非空时更新） |
| `systemPrompt` | string | 否 | 系统提示词（非空时更新） |
| `modelProvider` | string | 否 | 模型提供商（非空时更新） |
| `modelName` | string | 否 | 模型名称（非空时更新） |
| `modelParameters` | object | 否 | 模型参数（非nil时更新） |
| `openingDialogue` | string | 否 | 开场白（非空时更新） |

> `id` 由路径参数自动填充，请求体中无需传。

### 请求示例

```json
{
  "name": "高级客服助手",
  "status": "published",
  "modelParameters": {
    "maxTokens": 8192,
    "temperature": 0.5,
    "topP": 0.9
  }
}
```

### 响应 data

同创建接口的 `AgentDetailAdminResponse`。

---

## 5. 删除智能体

**DELETE** `/api/v1/agent/admin/:id`

### 路径参数

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string(uuid) | 智能体ID |

### 请求示例

```
DELETE /api/v1/agent/admin/660e8400-e29b-41d4-a716-446655440001
```

### 响应 data

`null`（无返回数据）

### 响应示例

```json
{
  "code": 200,
  "msg": "success",
  "data": null
}
```

> **注意**：删除智能体会级联删除关联的工具绑定、知识库绑定、子智能体绑定、工作流绑定。

---

## 枚举值说明

### status（智能体状态）

| 值 | 说明 |
|------|------|
| `"draft"` | 草稿 |
| `"published"` | 已发布 |
| `"archived"` | 已归档 |

### visibility（可见性）

| 值 | 说明 |
|------|------|
| `"private"` | 私有 |
| `"public"` | 公开 |
| `"link_only"` | 仅链接可访问 |

### modelProvider（模型提供商）

| 值 | 说明 |
|------|------|
| `"openai"` | OpenAI / 兼容接口 |
| `"ollama"` | Ollama 本地部署 |
| `"qwen"` | 通义千问 |
| `"deepseek"` | DeepSeek |
