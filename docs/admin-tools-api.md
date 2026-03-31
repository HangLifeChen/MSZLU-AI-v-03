# 管理后台 - 工具管理 API 文档

## 概述

本文档描述了管理后台工具管理相关的 API 接口，用于管理员对系统中的工具进行增删改查和统计操作。

**基础路径**: `/api/v1/admin/tools`

**认证方式**: 需要管理员权限认证

---

## 目录

1. [创建工具](#1-创建工具)
2. [查询工具列表](#2-查询工具列表)
3. [获取工具详情](#3-获取工具详情)
4. [更新工具](#4-更新工具)
5. [删除工具](#5-删除工具)
6. [获取工具统计信息](#6-获取工具统计信息)
7. [数据模型](#数据模型)

---

## 1. 创建工具

### 接口信息

- **接口路径**: `POST /api/v1/admin/tools`
- **功能描述**: 管理员创建新工具，可以指定创建者
- **权限要求**: 管理员权限

### 请求参数

**Headers**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**Body (JSON)**:

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 工具名称，最大长度255字符 |
| description | string | 否 | 工具描述 |
| toolType | string | 是 | 工具类型，可选值：`mcp`、`system` |
| isEnable | boolean | 是 | 是否启用，默认为 `true` |
| mcpConfig | object | 否 | MCP配置信息（当toolType为mcp时需要） |
| mcpConfig.type | string | 否 | MCP类型，如 `sse` |
| mcpConfig.url | string | 否 | MCP服务URL |
| mcpConfig.authenticationRequired | boolean | 否 | 是否需要认证 |
| mcpConfig.credentialType | string | 否 | 凭证类型 |
| creatorId | string | 是 | 创建者用户ID（UUID格式） |

### 请求示例

```json
{
  "name": "天气查询工具",
  "description": "用于查询天气信息的工具",
  "toolType": "mcp",
  "isEnable": true,
  "mcpConfig": {
    "type": "sse",
    "url": "https://example.com/mcp/weather",
    "authenticationRequired": true,
    "credentialType": "api_key"
  },
  "creatorId": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 响应示例

**成功响应 (200 OK)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "天气查询工具",
    "description": "用于查询天气信息的工具",
    "toolType": "mcp",
    "isEnable": true,
    "parametersSchema": {},
    "mcpConfig": {
      "type": "sse",
      "url": "https://example.com/mcp/weather",
      "authenticationRequired": true,
      "credentialType": "api_key"
    },
    "creatorId": "550e8400-e29b-41d4-a716-446655440000",
    "creatorName": "张三",
    "creatorEmail": "zhangsan@example.com",
    "createdAt": "2024-01-15T08:30:00Z",
    "updatedAt": "2024-01-15T08:30:00Z"
  }
}
```

**错误响应 (400 Bad Request)**:

```json
{
  "code": 400,
  "message": "参数验证失败",
  "data": null
}
```

---

## 2. 查询工具列表

### 接口信息

- **接口路径**: `GET /api/v1/admin/tools`
- **功能描述**: 管理员查询工具列表，支持多条件筛选和分页
- **权限要求**: 管理员权限

### 请求参数

**Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 工具名称（支持模糊搜索） |
| type | string | 否 | 工具类型，可选值：`mcp`、`system` |
| creatorId | string | 否 | 创建者用户ID（UUID格式） |
| isEnable | boolean | 否 | 是否启用状态 |
| page | integer | 否 | 页码，默认为 `1` |
| pageSize | integer | 否 | 每页数量，默认为 `10` |

### 请求示例

```
GET /api/v1/admin/tools?name=天气&type=mcp&page=1&pageSize=10
```

### 响应示例

**成功响应 (200 OK)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "天气查询工具",
        "description": "用于查询天气信息的工具",
        "toolType": "mcp",
        "isEnable": true,
        "creatorId": "550e8400-e29b-41d4-a716-446655440000",
        "creatorName": "张三",
        "creatorEmail": "zhangsan@example.com",
        "createdAt": "2024-01-15T08:30:00Z",
        "updatedAt": "2024-01-15T08:30:00Z"
      },
      {
        "id": "550e8400-e29b-41d4-a716-446655440002",
        "name": "代码生成工具",
        "description": "用于生成代码的工具",
        "toolType": "system",
        "isEnable": false,
        "creatorId": "550e8400-e29b-41d4-a716-446655440000",
        "creatorName": "张三",
        "creatorEmail": "zhangsan@example.com",
        "createdAt": "2024-01-16T10:20:00Z",
        "updatedAt": "2024-01-16T10:20:00Z"
      }
    ],
    "total": 2,
    "currentPage": 1,
    "pageSize": 10
  }
}
```

**空列表响应 (200 OK)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 0,
    "currentPage": 1,
    "pageSize": 10
  }
}
```

---

## 3. 获取工具详情

### 接口信息

- **接口路径**: `GET /api/v1/admin/tools/{id}`
- **功能描述**: 管理员获取指定工具的详细信息，包括创建者信息和参数配置
- **权限要求**: 管理员权限

### 请求参数

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 工具ID（UUID格式） |

### 请求示例

```
GET /api/v1/admin/tools/550e8400-e29b-41d4-a716-446655440001
```

### 响应示例

**成功响应 (200 OK)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "天气查询工具",
    "description": "用于查询天气信息的工具",
    "toolType": "mcp",
    "isEnable": true,
    "parametersSchema": {
      "location": {
        "type": "string",
        "description": "查询位置"
      },
      "date": {
        "type": "string",
        "description": "查询日期"
      }
    },
    "mcpConfig": {
      "type": "sse",
      "url": "https://example.com/mcp/weather",
      "authenticationRequired": true,
      "credentialType": "api_key"
    },
    "creatorId": "550e8400-e29b-41d4-a716-446655440000",
    "creatorName": "张三",
    "creatorEmail": "zhangsan@example.com",
    "createdAt": "2024-01-15T08:30:00Z",
    "updatedAt": "2024-01-15T08:30:00Z"
  }
}
```

**错误响应 (404 Not Found)**:

```json
{
  "code": 404,
  "message": "工具不存在",
  "data": null
}
```

---

## 4. 更新工具

### 接口信息

- **接口路径**: `PUT /api/v1/admin/tools/{id}`
- **功能描述**: 管理员更新指定工具的信息
- **权限要求**: 管理员权限

### 请求参数

**Headers**:
```
Content-Type: application/json
Authorization: Bearer {token}
```

**Path Parameters**:

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 工具ID（UUID格式） |

**Body (JSON)**:

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 工具ID（UUID格式），需与路径参数一致 |
| name | string | 是 | 工具名称，最大长度255字符 |
| description | string | 否 | 工具描述 |
| toolType | string | 是 | 工具类型，可选值：`mcp`、`system` |
| isEnable | boolean | 是 | 是否启用 |
| mcpConfig | object | 否 | MCP配置信息（当toolType为mcp时需要） |
| mcpConfig.type | string | 否 | MCP类型，如 `sse` |
| mcpConfig.url | string | 否 | MCP服务URL |
| mcpConfig.authenticationRequired | boolean | 否 | 是否需要认证 |
| mcpConfig.credentialType | string | 否 | 凭证类型 |

### 请求示例

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "天气查询工具（更新版）",
  "description": "用于查询天气信息的工具，支持多城市",
  "toolType": "mcp",
  "isEnable": false,
  "mcpConfig": {
    "type": "sse",
    "url": "https://example.com/mcp/weather/v2",
    "authenticationRequired": false,
    "credentialType": ""
  }
}
```

### 响应示例

**成功响应 (200 OK)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "天气查询工具（更新版）",
    "description": "用于查询天气信息的工具，支持多城市",
    "toolType": "mcp",
    "isEnable": false,
    "parametersSchema": {},
    "mcpConfig": {
      "type": "sse",
      "url": "https://example.com/mcp/weather/v2",
      "authenticationRequired": false,
      "credentialType": ""
    },
    "creatorId": "550e8400-e29b-41d4-a716-446655440000",
    "creatorName": "张三",
    "creatorEmail": "zhangsan@example.com",
    "createdAt": "2024-01-15T08:30:00Z",
    "updatedAt": "2024-01-17T14:20:00Z"
  }
}
```

**错误响应 (404 Not Found)**:

```json
{
  "code": 404,
  "message": "工具不存在",
  "data": null
}
```

---

## 5. 删除工具

### 接口信息

- **接口路径**: `DELETE /api/v1/admin/tools/{id}`
- **功能描述**: 管理员删除指定的工具
- **权限要求**: 管理员权限

### 请求参数

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 工具ID（UUID格式） |

### 请求示例

```
DELETE /api/v1/admin/tools/550e8400-e29b-41d4-a716-446655440001
```

### 响应示例

**成功响应 (200 OK)**:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

**错误响应 (404 Not Found)**:

```json
{
  "code": 404,
  "message": "工具不存在",
  "data": null
}
```

---

## 6. 获取工具统计信息

### 接口信息

- **接口路径**: `GET /api/v1/admin/tools/stats`
- **功能描述**: 获取工具的统计信息，包括总数、启用/禁用数量、MCP/系统工具数量等
- **权限要求**: 管理员权限

### 请求参数

**Headers**:
```
Authorization: Bearer {token}
```

### 请求示例

```
GET /api/v1/admin/tools/stats
```

### 响应示例

**成功响应 (200 OK)**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "totalTools": 150,
    "enabledTools": 120,
    "disabledTools": 30,
    "mcpTools": 80,
    "systemTools": 70
  }
}
```

---

## 数据模型

### ToolType（工具类型）

| 值 | 说明 |
|----|------|
| `mcp` | MCP工具 |
| `system` | 系统工具 |

### McpConfig（MCP配置）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| type | string | MCP类型，如 `sse` |
| url | string | MCP服务URL |
| authenticationRequired | boolean | 是否需要认证 |
| credentialType | string | 凭证类型 |

### ToolDetailResponse（工具详情响应）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | string | 工具ID（UUID格式） |
| name | string | 工具名称 |
| description | string | 工具描述 |
| toolType | string | 工具类型 |
| isEnable | boolean | 是否启用 |
| parametersSchema | object | 参数配置Schema |
| mcpConfig | object | MCP配置信息 |
| creatorId | string | 创建者用户ID |
| creatorName | string | 创建者用户名 |
| creatorEmail | string | 创建者邮箱 |
| createdAt | string | 创建时间（ISO 8601格式） |
| updatedAt | string | 更新时间（ISO 8601格式） |

### ToolListResponse（工具列表项响应）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | string | 工具ID（UUID格式） |
| name | string | 工具名称 |
| description | string | 工具描述 |
| toolType | string | 工具类型 |
| isEnable | boolean | 是否启用 |
| creatorId | string | 创建者用户ID |
| creatorName | string | 创建者用户名 |
| creatorEmail | string | 创建者邮箱 |
| createdAt | string | 创建时间（ISO 8601格式） |
| updatedAt | string | 更新时间（ISO 8601格式） |

### ListToolsAdminResponse（工具列表响应）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| list | array | 工具列表 |
| total | integer | 总记录数 |
| currentPage | integer | 当前页码 |
| pageSize | integer | 每页数量 |

### ToolStats（工具统计信息）

| 字段名 | 类型 | 说明 |
|--------|------|------|
| totalTools | integer | 工具总数 |
| enabledTools | integer | 已启用工具数 |
| disabledTools | integer | 已禁用工具数 |
| mcpTools | integer | MCP工具数 |
| systemTools | integer | 系统工具数 |

---

## 通用响应格式

所有接口都遵循统一的响应格式：

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": { /* 响应数据 */ }
}
```

### 错误响应

```json
{
  "code": {错误码},
  "message": "{错误信息}",
  "data": null
}
```

### 常见错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权或token无效 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 注意事项

1. **时间格式**: 所有时间字段均使用 ISO 8601 格式（如：`2024-01-15T08:30:00Z`）
2. **UUID格式**: 所有ID字段使用 UUID 格式
3. **分页**: 列表查询接口支持分页，默认每页10条记录
4. **模糊搜索**: 名称字段支持模糊搜索（LIKE查询）
5. **认证**: 所有接口都需要在请求头中携带有效的认证token
6. **权限**: 所有接口都需要管理员权限才能访问
7. **MCP配置**: 当工具类型为 `mcp` 时，建议提供完整的 `mcpConfig` 配置

---

## 前端集成示例

### 使用 Axios 请求示例

```javascript
import axios from 'axios';

// 配置基础URL和认证
const api = axios.create({
  baseURL: process.env.VUE_APP_API_BASE_URL,
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('token')}`
  }
});

// 1. 创建工具
async function createTool(toolData) {
  try {
    const response = await api.post('/api/v1/admin/tools', toolData);
    return response.data;
  } catch (error) {
    console.error('创建工具失败:', error);
    throw error;
  }
}

// 2. 查询工具列表
async function getToolList(params) {
  try {
    const response = await api.get('/api/v1/admin/tools', { params });
    return response.data;
  } catch (error) {
    console.error('查询工具列表失败:', error);
    throw error;
  }
}

// 3. 获取工具详情
async function getToolDetail(toolId) {
  try {
    const response = await api.get(`/api/v1/admin/tools/${toolId}`);
    return response.data;
  } catch (error) {
    console.error('获取工具详情失败:', error);
    throw error;
  }
}

// 4. 更新工具
async function updateTool(toolId, toolData) {
  try {
    const response = await api.put(`/api/v1/admin/tools/${toolId}`, toolData);
    return response.data;
  } catch (error) {
    console.error('更新工具失败:', error);
    throw error;
  }
}

// 5. 删除工具
async function deleteTool(toolId) {
  try {
    const response = await api.delete(`/api/v1/admin/tools/${toolId}`);
    return response.data;
  } catch (error) {
    console.error('删除工具失败:', error);
    throw error;
  }
}

// 6. 获取工具统计
async function getToolStats() {
  try {
    const response = await api.get('/api/v1/admin/tools/stats');
    return response.data;
  } catch (error) {
    console.error('获取工具统计失败:', error);
    throw error;
  }
}
```



---

## 更新日志

| 版本 | 日期 | 说明 |
|------|------|------|
| 1.0.0 | 2024-01-15 | 初始版本 |

---

## 联系方式

如有问题或建议，请联系开发团队。
