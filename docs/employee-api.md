# 员工管理接口文档

## 概述

本文档描述了员工管理模块的 API 接口，包括员工的增删改查功能。

**基础路径**: `/api/v1/employees`

**认证方式**: 需要在请求头中携带有效的 Authorization Token

---

## 1. 创建员工

### 接口信息

| 项目 | 内容 |
|------|------|
| URL | `/api/v1/employees/create` |
| Method | `POST` |
| Content-Type | `application/json` |

### 请求参数

```json
{
  "employeeNo": "EMP001",           // 必填，员工编号，最大50字符
  "name": "张三",                    // 必填，员工姓名，最大100字符
  "gender": 1,                      // 可选，性别：1-男，2-女
  "phone": "13800138000",           // 可选，电话号码，最大20字符
  "email": "zhangsan@example.com",  // 可选，邮箱地址，最大100字符
  "departmentId": 1,                // 可选，部门ID
  "position": "软件工程师",          // 可选，职位，最大100字符
  "hireDate": "2024-01-15",         // 可选，入职日期，格式：YYYY-MM-DD
  "status": 1,                      // 可选，状态：1-在职(默认)，2-离职，3-休假
  "username": "zhangsan",           // 可选，登录用户名，最大100字符，唯一
  "password": "password123",        // 可选，登录密码，最大255字符
  "remark": "备注信息"               // 可选，备注信息
}
```

### 响应示例

**成功响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "employeeNo": "EMP001",
    "name": "张三",
    "gender": 1,
    "phone": "13800138000",
    "email": "zhangsan@example.com",
    "departmentId": 1,
    "position": "软件工程师",
    "hireDate": "2024-01-15T00:00:00Z",
    "status": 1,
    "username": "zhangsan",
    "remark": "备注信息",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

**错误响应**:
```json
{
  "code": 90002,
  "msg": "员工编号已存在",
  "data": null
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 90002 | 员工编号已存在 |
| 90003 | 员工用户名已存在 |

---

## 2. 获取员工详情

### 接口信息

| 项目 | 内容 |
|------|------|
| URL | `/api/v1/employees/:id` |
| Method | `GET` |

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 员工ID |

### 请求示例

```
GET /api/v1/employees/1
```

### 响应示例

**成功响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "employeeNo": "EMP001",
    "name": "张三",
    "gender": 1,
    "phone": "13800138000",
    "email": "zhangsan@example.com",
    "departmentId": 1,
    "position": "软件工程师",
    "hireDate": "2024-01-15T00:00:00Z",
    "status": 1,
    "username": "zhangsan",
    "remark": "备注信息",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

**错误响应**:
```json
{
  "code": 90001,
  "msg": "员工不存在",
  "data": null
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 90001 | 员工不存在 |

---

## 3. 更新员工信息

### 接口信息

| 项目 | 内容 |
|------|------|
| URL | `/api/v1/employees/update` |
| Method | `PUT` |
| Content-Type | `application/json` |

### 请求参数

```json
{
  "id": 1,                          // 必填，员工ID
  "employeeNo": "EMP001",           // 可选，员工编号
  "name": "张三",                    // 可选，员工姓名
  "gender": 1,                      // 可选，性别
  "phone": "13800138001",           // 可选，电话号码
  "email": "zhangsan@example.com",  // 可选，邮箱地址
  "departmentId": 2,                // 可选，部门ID
  "position": "高级软件工程师",      // 可选，职位
  "hireDate": "2024-01-15",         // 可选，入职日期
  "status": 1,                      // 可选，状态
  "username": "zhangsan",           // 可选，登录用户名
  "password": "newpassword123",     // 可选，登录密码
  "remark": "更新备注信息"           // 可选，备注信息
}
```

### 响应示例

**成功响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "employeeNo": "EMP001",
    "name": "张三",
    "gender": 1,
    "phone": "13800138001",
    "email": "zhangsan@example.com",
    "departmentId": 2,
    "position": "高级软件工程师",
    "hireDate": "2024-01-15T00:00:00Z",
    "status": 1,
    "username": "zhangsan",
    "remark": "更新备注信息",
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T11:00:00Z"
  }
}
```

**错误响应**:
```json
{
  "code": 90001,
  "msg": "员工不存在",
  "data": null
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 90001 | 员工不存在 |

---

## 4. 删除员工

### 接口信息

| 项目 | 内容 |
|------|------|
| URL | `/api/v1/employees/:id` |
| Method | `DELETE` |

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 员工ID |

### 请求示例

```
DELETE /api/v1/employees/1
```

### 响应示例

**成功响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

**错误响应**:
```json
{
  "code": 90001,
  "msg": "员工不存在",
  "data": null
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 90001 | 员工不存在 |

---

## 5. 查询员工列表

### 接口信息

| 项目 | 内容 |
|------|------|
| URL | `/api/v1/employees/list` |
| Method | `GET` |

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| employeeNo | string | 否 | 员工编号，支持模糊查询 |
| name | string | 否 | 员工姓名，支持模糊查询 |
| departmentId | int64 | 否 | 部门ID |
| status | int | 否 | 状态：1-在职，2-离职，3-休假 |
| page | int | 否 | 页码，默认1 |
| pageSize | int | 否 | 每页数量，默认10 |

### 请求示例

```
GET /api/v1/employees/list?name=张&page=1&pageSize=10&status=1
```

### 响应示例

**成功响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "employees": [
      {
        "id": 1,
        "employeeNo": "EMP001",
        "name": "张三",
        "gender": 1,
        "phone": "13800138000",
        "email": "zhangsan@example.com",
        "departmentId": 1,
        "position": "软件工程师",
        "hireDate": "2024-01-15T00:00:00Z",
        "status": 1,
        "username": "zhangsan",
        "remark": "备注信息",
        "createdAt": "2024-01-15T10:30:00Z",
        "updatedAt": "2024-01-15T10:30:00Z"
      },
      {
        "id": 2,
        "employeeNo": "EMP002",
        "name": "张四",
        "gender": 2,
        "phone": "13800138001",
        "email": "zhangsi@example.com",
        "departmentId": 1,
        "position": "产品经理",
        "hireDate": "2024-02-01T00:00:00Z",
        "status": 1,
        "username": "zhangsi",
        "remark": "",
        "createdAt": "2024-02-01T09:00:00Z",
        "updatedAt": "2024-02-01T09:00:00Z"
      }
    ],
    "total": 2
  }
}
```

---

## 数据字典

### 员工状态 (status)

| 值 | 说明 |
|----|------|
| 1 | 在职 |
| 2 | 离职 |
| 3 | 休假 |

### 性别 (gender)

| 值 | 说明 |
|----|------|
| 1 | 男 |
| 2 | 女 |

---

## 员工对象结构

### EmployeeResponse

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 员工ID |
| employeeNo | string | 员工编号 |
| name | string | 员工姓名 |
| gender | int16 | 性别 |
| phone | string | 电话号码 |
| email | string | 邮箱地址 |
| departmentId | int64 | 部门ID |
| position | string | 职位 |
| hireDate | string | 入职日期 (ISO 8601格式) |
| status | int16 | 员工状态 |
| username | string | 登录用户名 |
| remark | string | 备注信息 |
| createdAt | string | 创建时间 (ISO 8601格式) |
| updatedAt | string | 更新时间 (ISO 8601格式) |

---

## 通用错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10004 | Token无效 |
| 90001 | 员工不存在 |
| 90002 | 员工编号已存在 |
| 90003 | 员工用户名已存在 |

---

## 前端调用示例

### JavaScript/TypeScript (Axios)

```typescript
import axios from 'axios';

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

// 创建员工
async function createEmployee(data: CreateEmployeeReq) {
  const response = await api.post('/employees/create', data);
  return response.data;
}

// 获取员工详情
async function getEmployee(id: number) {
  const response = await api.get(`/employees/${id}`);
  return response.data;
}

// 更新员工信息
async function updateEmployee(data: UpdateEmployeeReq) {
  const response = await api.put('/employees/update', data);
  return response.data;
}

// 删除员工
async function deleteEmployee(id: number) {
  const response = await api.delete(`/employees/${id}`);
  return response.data;
}

// 查询员工列表
async function listEmployees(params: ListEmployeesReq) {
  const response = await api.get('/employees/list', { params });
  return response.data;
}
```

### TypeScript 类型定义

```typescript
// 员工状态枚举
enum EmployeeStatus {
  Active = 1,    // 在职
  Resigned = 2,  // 离职
  OnLeave = 3    // 休假
}

// 性别枚举
enum Gender {
  Male = 1,    // 男
  Female = 2   // 女
}

// 创建员工请求
interface CreateEmployeeReq {
  employeeNo: string;
  name: string;
  gender?: number;
  phone?: string;
  email?: string;
  departmentId?: number;
  position?: string;
  hireDate?: string;
  status?: EmployeeStatus;
  username?: string;
  password?: string;
  remark?: string;
}

// 更新员工请求
interface UpdateEmployeeReq {
  id: number;
  employeeNo?: string;
  name?: string;
  gender?: number;
  phone?: string;
  email?: string;
  departmentId?: number;
  position?: string;
  hireDate?: string;
  status?: EmployeeStatus;
  username?: string;
  password?: string;
  remark?: string;
}

// 查询员工列表请求
interface ListEmployeesReq {
  employeeNo?: string;
  name?: string;
  departmentId?: number;
  status?: EmployeeStatus;
  page?: number;
  pageSize?: number;
}

// 员工响应
interface EmployeeResponse {
  id: number;
  employeeNo: string;
  name: string;
  gender?: number;
  phone: string;
  email: string;
  departmentId?: number;
  position: string;
  hireDate?: string;
  status: EmployeeStatus;
  username: string;
  remark: string;
  createdAt: string;
  updatedAt: string;
}

// 员工列表响应
interface ListEmployeeResponse {
  employees: EmployeeResponse[];
  total: number;
}

// API响应包装
interface ApiResponse<T> {
  code: number;
  msg: string;
  data: T;
}
```
