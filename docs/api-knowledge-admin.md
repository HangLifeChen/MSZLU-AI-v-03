# 知识库管理员接口文档

> 基础路径：`/api/v1/knowledge/admin`
> 所有接口均需管理员权限认证，通过请求头携带 Token 鉴权。

---

## 1. 创建知识库

**POST** `/api/v1/knowledge/admin`

### 请求参数（JSON Body）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 知识库名称 |
| description | string | 否 | 描述 |
| chatModelName | string | 否 | 对话模型名称 |
| chatModelProvider | string | 否 | 对话模型提供商 |
| embeddingModelName | string | 否 | 向量模型名称 |
| embeddingModelProvider | string | 否 | 向量模型提供商 |
| storageType | string | 否 | 存储类型，可选值：`"es"`、`"milvus"` |
| tags | string[] | 否 | 标签列表 |
| creatorId | string(UUID) | 是 | 指定创建者用户ID |

### 请求示例

```json
{
  "name": "产品文档库",
  "description": "存放所有产品相关文档",
  "chatModelName": "qwen-plus",
  "chatModelProvider": "qwen",
  "embeddingModelName": "text-embedding-v3",
  "embeddingModelProvider": "qwen",
  "storageType": "milvus",
  "tags": ["产品", "文档"],
  "creatorId": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 响应示例

```json
{
  "code": 200,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "产品文档库",
    "description": "存放所有产品相关文档",
    "embeddingModelName": "text-embedding-v3",
    "embeddingModelProvider": "qwen",
    "chatModelName": "qwen-plus",
    "chatModelProvider": "qwen",
    "storageType": "milvus",
    "tags": ["产品", "文档"],
    "status": "active",
    "documentCount": 0,
    "totalSize": 0,
    "creatorId": "550e8400-e29b-41d4-a716-446655440000",
    "creatorName": "张三",
    "creatorEmail": "zhangsan@example.com",
    "createdAt": "2026-04-20T15:00:00Z",
    "updatedAt": "2026-04-20T15:00:00Z"
  }
}
```

---

## 2. 获取知识库列表

**GET** `/api/v1/knowledge/admin`

### 请求参数（Query String）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 按名称模糊搜索 |
| search | string | 否 | 按名称模糊搜索（与 name 等效） |
| creatorId | string | 否 | 按创建者ID筛选 |
| status | string | 否 | 按状态筛选，可选值：`"active"`、`"disabled"` |
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页条数，默认 10 |

### 请求示例

```
GET /api/v1/knowledge/admin?page=1&pageSize=10&status=active&name=产品
```

### 响应示例

```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "产品文档库",
        "description": "存放所有产品相关文档",
        "storageType": "milvus",
        "status": "active",
        "documentCount": 5,
        "totalSize": 1024000,
        "creatorId": "550e8400-e29b-41d4-a716-446655440000",
        "creatorName": "张三",
        "creatorEmail": "zhangsan@example.com",
        "createdAt": "2026-04-20T15:00:00Z",
        "updatedAt": "2026-04-20T16:00:00Z"
      }
    ],
    "total": 1,
    "currentPage": 1,
    "pageSize": 10
  }
}
```

---

## 3. 获取知识库详情

**GET** `/api/v1/knowledge/admin/:id`

### 路径参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string(UUID) | 是 | 知识库ID |

### 请求示例

```
GET /api/v1/knowledge/admin/660e8400-e29b-41d4-a716-446655440001
```

### 响应示例

```json
{
  "code": 200,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "产品文档库",
    "description": "存放所有产品相关文档",
    "embeddingModelName": "text-embedding-v3",
    "embeddingModelProvider": "qwen",
    "chatModelName": "qwen-plus",
    "chatModelProvider": "qwen",
    "storageType": "milvus",
    "tags": ["产品", "文档"],
    "status": "active",
    "documentCount": 5,
    "totalSize": 1024000,
    "creatorId": "550e8400-e29b-41d4-a716-446655440000",
    "creatorName": "张三",
    "creatorEmail": "zhangsan@example.com",
    "createdAt": "2026-04-20T15:00:00Z",
    "updatedAt": "2026-04-20T16:00:00Z"
  }
}
```

---

## 4. 更新知识库

**PUT** `/api/v1/knowledge/admin/:id`

### 路径参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string(UUID) | 是 | 知识库ID |

### 请求参数（JSON Body）

所有字段均为可选，仅传递需要更新的字段。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 知识库名称 |
| description | string | 否 | 描述 |
| chatModelName | string | 否 | 对话模型名称 |
| chatModelProvider | string | 否 | 对话模型提供商 |
| embeddingModelName | string | 否 | 向量模型名称 |
| embeddingModelProvider | string | 否 | 向量模型提供商 |
| tags | string[] | 否 | 标签列表（传空数组 `[]` 可清空标签） |
| status | string | 否 | 状态，可选值：`"active"`、`"disabled"` |

### 请求示例

```json
{
  "name": "产品文档库V2",
  "description": "更新后的描述",
  "status": "disabled",
  "tags": ["产品", "文档", "V2"]
}
```

### 响应示例

```json
{
  "code": 200,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "产品文档库V2",
    "description": "更新后的描述",
    "embeddingModelName": "text-embedding-v3",
    "embeddingModelProvider": "qwen",
    "chatModelName": "qwen-plus",
    "chatModelProvider": "qwen",
    "storageType": "milvus",
    "tags": ["产品", "文档", "V2"],
    "status": "disabled",
    "documentCount": 5,
    "totalSize": 1024000,
    "creatorId": "550e8400-e29b-41d4-a716-446655440000",
    "creatorName": "张三",
    "creatorEmail": "zhangsan@example.com",
    "createdAt": "2026-04-20T15:00:00Z",
    "updatedAt": "2026-04-20T17:00:00Z"
  }
}
```

---

## 5. 删除知识库

**DELETE** `/api/v1/knowledge/admin/:id`

### 路径参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string(UUID) | 是 | 知识库ID |

### 请求示例

```
DELETE /api/v1/knowledge/admin/660e8400-e29b-41d4-a716-446655440001
```

### 响应示例

```json
{
  "code": 200,
  "data": null
}
```

---

## 枚举值参考

### storageType 存储类型

| 值 | 说明 |
|----|------|
| `es` | Elasticsearch |
| `milvus` | Milvus 向量数据库 |

### status 知识库状态

| 值 | 说明 |
|----|------|
| `active` | 正常 |
| `disabled` | 已禁用 |

---

## 与普通用户接口的区别

| 特性 | 普通用户接口 `/api/v1/knowledge` | 管理员接口 `/api/v1/knowledge/admin` |
|------|------|------|
| 数据范围 | 仅当前用户创建的知识库 | 所有用户的知识库 |
| 创建 | 自动使用当前用户ID | 需指定 `creatorId` |
| 创建者信息 | 不返回 | 返回 `creatorName`、`creatorEmail` |
| 状态管理 | 不可修改 | 可通过 `status` 字段启用/禁用 |
| 存储类型 | 不可选择 | 可通过 `storageType` 指定 |
| 时间格式 | Unix 时间戳（秒） | RFC3339 字符串 |
| 列表筛选 | 仅支持 search | 支持 name、search、creatorId、status |
| 分页参数 | 在 JSON Body 中（POST） | 在 Query String 中（GET） |

---

## 错误码

| HTTP 状态码 | code | 说明 |
|-------------|------|------|
| 400 | - | 参数错误（缺少必填字段） |
| 401 | - | 未认证（Token 缺失或无效） |
| 403 | - | 无管理员权限 |
| 404 | - | 知识库不存在 |
| 500 | - | 服务器内部错误 |
