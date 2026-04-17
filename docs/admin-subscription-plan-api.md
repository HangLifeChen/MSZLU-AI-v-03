# 管理后台 - 订阅计划配置管理 API

Base URL: `/api/v1/subscription/admin`

## 统一响应格式

```json
{
  "code": 200,
  "msg": "success",
  "data": {}
}
```

---

## 1. 创建订阅计划配置

**POST** `/api/v1/subscription/admin/plan`

### 请求体 (JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `name` | string | 是 | 计划名称（最大50字符） |
| `plan` | string | 是 | 计划标识，可选值见下方枚举说明 |
| `price` | int64 | 是 | 月费价格，单位：分 |
| `description` | string | 否 | 计划描述（最大255字符） |
| `quarterRate` | float64 | 否 | 季度折扣率，默认 0.9 |
| `yearRate` | float64 | 否 | 年度折扣率，默认 0.8 |
| `maxAgents` | int64 | 否 | 最大 Agent 数量，默认 0（不限制） |
| `maxWorkflows` | int64 | 否 | 最大工作流数量，默认 0（不限制） |
| `maxKnowledgeBaseSize` | int64 | 否 | 最大知识库大小（MB），默认 0（不限制） |

### 请求示例

```json
{
  "name": "高级版",
  "plan": "pro",
  "price": 9900,
  "description": "适合专业用户，包含更多功能",
  "quarterRate": 0.85,
  "yearRate": 0.75,
  "maxAgents": 50,
  "maxWorkflows": 100,
  "maxKnowledgeBaseSize": 5120
}
```

### 响应 data

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 计划配置ID |
| `name` | string | 计划名称 |
| `plan` | string | 计划标识 |
| `price` | int64 | 月费价格（分） |
| `description` | string | 计划描述 |
| `quarterRate` | float64 | 季度折扣率 |
| `yearRate` | float64 | 年度折扣率 |
| `configs` | object | 资源配额配置 |
| `configs.maxAgents` | int64 | 最大 Agent 数量 |
| `configs.maxWorkflows` | int64 | 最大工作流数量 |
| `configs.maxKnowledgeBaseSize` | int64 | 最大知识库大小（MB） |

### 响应示例

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "高级版",
    "plan": "pro",
    "price": 9900,
    "description": "适合专业用户，包含更多功能",
    "quarterRate": 0.85,
    "yearRate": 0.75,
    "configs": {
      "maxAgents": 50,
      "maxWorkflows": 100,
      "maxKnowledgeBaseSize": 5120
    }
  }
}
```

> **注意**：每个 `plan` 标识只能创建一个配置，重复创建会返回错误。

---

## 2. 更新订阅计划配置

**PUT** `/api/v1/subscription/admin/plan`

### 请求体 (JSON)

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `id` | int64 | 是 | 计划配置ID |
| `name` | string | 否 | 计划名称 |
| `plan` | string | 否 | 计划标识 |
| `price` | int64 | 否 | 月费价格，单位：分（传值 > 0 时更新） |
| `description` | string | 否 | 计划描述（非空时更新） |
| `quarterRate` | float64 | 否 | 季度折扣率（传值 > 0 时更新） |
| `yearRate` | float64 | 否 | 年度折扣率（传值 > 0 时更新） |
| `maxAgents` | int64 | 否 | 最大 Agent 数量（传值 > 0 时更新） |
| `maxWorkflows` | int64 | 否 | 最大工作流数量（传值 > 0 时更新） |
| `maxKnowledgeBaseSize` | int64 | 否 | 最大知识库大小 MB（传值 > 0 时更新） |

> **更新策略**：仅更新非零值字段。`name`、`plan`、`description` 非空字符串时更新；`price`、`quarterRate`、`yearRate`、`maxAgents`、`maxWorkflows`、`maxKnowledgeBaseSize` 大于 0 时更新。

### 请求示例

```json
{
  "id": 1,
  "price": 12900,
  "quarterRate": 0.8,
  "yearRate": 0.7,
  "maxAgents": 100
}
```

### 响应 data

同创建接口响应结构。

### 响应示例

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "高级版",
    "plan": "pro",
    "price": 12900,
    "description": "适合专业用户，包含更多功能",
    "quarterRate": 0.8,
    "yearRate": 0.7,
    "configs": {
      "maxAgents": 100,
      "maxWorkflows": 100,
      "maxKnowledgeBaseSize": 5120
    }
  }
}
```

---

## 3. 删除订阅计划配置

**DELETE** `/api/v1/subscription/admin/plan/:id`

### 路径参数

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 计划配置ID |

### 请求示例

```
DELETE /api/v1/subscription/admin/plan/1
```

### 响应示例

```json
{
  "code": 200,
  "msg": "success",
  "data": true
}
```

---

## 4. 获取所有订阅计划

**GET** `/api/v1/subscription/admin/plans`

### 请求示例

```
GET /api/v1/subscription/admin/plans
```

### 响应 data

`SubscriptionPlanConfigResponse[]` 数组，元素结构同创建接口响应。

### 响应示例

```json
{
  "code": 200,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "name": "免费版",
      "plan": "free",
      "price": 0,
      "description": "免费体验基础功能",
      "quarterRate": 0.9,
      "yearRate": 0.8,
      "configs": {
        "maxAgents": 3,
        "maxWorkflows": 5,
        "maxKnowledgeBaseSize": 100
      }
    },
    {
      "id": 2,
      "name": "基础版",
      "plan": "basic",
      "price": 2900,
      "description": "适合个人用户",
      "quarterRate": 0.9,
      "yearRate": 0.8,
      "configs": {
        "maxAgents": 10,
        "maxWorkflows": 30,
        "maxKnowledgeBaseSize": 1024
      }
    },
    {
      "id": 3,
      "name": "高级版",
      "plan": "pro",
      "price": 9900,
      "description": "适合专业用户",
      "quarterRate": 0.85,
      "yearRate": 0.75,
      "configs": {
        "maxAgents": 50,
        "maxWorkflows": 100,
        "maxKnowledgeBaseSize": 5120
      }
    },
    {
      "id": 4,
      "name": "企业版",
      "plan": "enterprise",
      "price": 29900,
      "description": "适合企业团队",
      "quarterRate": 0.8,
      "yearRate": 0.7,
      "configs": {
        "maxAgents": 0,
        "maxWorkflows": 0,
        "maxKnowledgeBaseSize": 0
      }
    }
  ]
}
```

---

## 枚举值说明

### plan（订阅计划标识）

| 值 | 说明 |
|------|------|
| `"free"` | 免费版 |
| `"basic"` | 基础版 |
| `"pro"` | 高级版 |
| `"enterprise"` | 企业版 |

### 折扣率说明

| 字段 | 说明 | 示例 |
|------|------|------|
| `quarterRate` | 季度折扣率，实际价格 = 月费 × 3 × quarterRate | 月费 9900 分，季度价 = 9900 × 3 × 0.85 = 25245 分 |
| `yearRate` | 年度折扣率，实际价格 = 月费 × 12 × yearRate | 月费 9900 分，年度价 = 9900 × 12 × 0.75 = 89100 分 |
