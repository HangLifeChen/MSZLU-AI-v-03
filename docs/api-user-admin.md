# 用户管理员接口文档

> 基础路径：`/api/user/admin`
> 所有接口均需管理员权限认证，通过请求头携带 Token 鉴权。

---

## 1. 获取用户列表

**GET** `/api/user/admin/list`

### 请求参数（Query String）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 否 | 按用户名模糊搜索 |
| email | string | 否 | 按邮箱模糊搜索 |
| status | int | 否 | 按状态筛选，可选值：`1`(正常)、`2`(禁用)、`3`(待审核) |
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页条数，默认 10 |

### 请求示例

```
GET /api/user/admin/list?page=1&pageSize=10&status=1&username=张
```

### 响应示例

```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "张三",
        "avatar": "https://oss.example.com/avatar/xxx.png",
        "status": 1,
        "lastLoginTime": 1745184000000,
        "currentPlan": "free",
        "email": "zhangsan@example.com",
        "emailVerified": true,
        "introduction": "管理员",
        "telephoneNumber": "13800138000"
      }
    ],
    "total": 50,
    "currentPage": 1,
    "pageSize": 10
  }
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string(UUID) | 用户ID |
| username | string | 用户名 |
| avatar | string | 头像URL |
| status | int | 状态：`1`=正常，`2`=禁用，`3`=待审核 |
| lastLoginTime | int64 | 最后登录时间（毫秒时间戳） |
| currentPlan | string | 当前套餐：`"free"`、`"basic"`、`"pro"` |
| email | string | 邮箱 |
| emailVerified | bool | 邮箱是否已验证 |
| introduction | string | 个人简介 |
| telephoneNumber | string | 手机号 |
| total | int64 | 总记录数 |
| currentPage | int | 当前页码 |
| pageSize | int | 每页条数 |

---

## 枚举值参考

### status 用户状态

| 值 | 说明 |
|----|------|
| `1` | 正常 |
| `2` | 禁用 |
| `3` | 待审核 |

### currentPlan 订阅套餐

| 值 | 说明 |
|----|------|
| `free` | 免费版 |
| `basic` | 基础版 |
| `pro` | 专业版 |

---

## 错误码

| HTTP 状态码 | code | 说明 |
|-------------|------|------|
| 400 | - | 参数错误 |
| 401 | - | 未认证（Token 缺失或无效） |
| 403 | - | 无管理员权限 |
| 500 | - | 服务器内部错误 |
