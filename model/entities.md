# 数据库实体文档

## BaseModel（基础模型，嵌入多个实体）

| **序号** | **字段名**   | **数据类型** | **备注** | **是否允许为空** |
| -------- | ------------ | ------------ | -------- | ---------------- |
| 1        | id           | uuid         | 唯一ID   | 否               |
| 2        | created_at   | timestamp    | 创建时间 | 否               |
| 3        | updated_at   | timestamp    | 更新时间 | 否               |
| 4        | deleted_at   | timestamp    | 删除时间 | 是               |

---

## users（用户表）

| **序号** | **字段名**       | **数据类型** | **备注**     | **是否允许为空** |
| -------- | ---------------- | ------------ | ------------ | ---------------- |
| 1        | id               | uuid         | 唯一ID       | 否               |
| 2        | username         | varchar      | 用户名       | 否               |
| 3        | password         | varchar      | 密码         | 是               |
| 4        | avatar           | varchar      | 头像链接     | 是               |
| 5        | status           | smallint     | 状态         | 是               |
| 6        | last_login_time  | timestamptz  | 登录时间     | 是               |
| 7        | current_plan     | varchar(20)  | 当前订阅计划 | 否               |
| 8        | email            | varchar(100) | 邮箱         | 否               |
| 9        | email_verified   | boolean      | 邮箱验证     | 是               |
| 10       | introduction     | varchar(255) | 个人简介     | 是               |
| 11       | telephone_number | varchar(20)  | 手机号       | 是               |
| 12       | created_at       | timestamp    | 创建时间     | 否               |
| 13       | updated_at       | timestamp    | 更新时间     | 否               |
| 14       | deleted_at       | timestamp    | 删除时间     | 是               |

---

## employees（员工表）

| **序号** | **字段名**   | **数据类型** | **备注** | **是否允许为空** |
| -------- | ------------- | ------------ | -------- | ---------------- |
| 1        | id            | bigint       | 唯一ID   | 否               |
| 2        | employee_no   | varchar(50)  | 员工编号 | 否               |
| 3        | name          | varchar(100) | 姓名     | 否               |
| 4        | gender        | smallint     | 性别     | 是               |
| 5        | phone         | varchar(20)  | 电话     | 是               |
| 6        | email         | varchar(100) | 邮箱     | 是               |
| 7        | department_id | bigint       | 部门ID   | 是               |
| 8        | position      | varchar(100) | 职位     | 是               |
| 9        | hire_date     | date         | 入职日期 | 是               |
| 10       | status        | smallint     | 状态     | 是               |
| 11       | username      | varchar(100) | 用户名   | 是               |
| 12       | password      | varchar(255) | 密码     | 是               |
| 13       | remark        | text         | 备注     | 是               |
| 14       | created_at    | timestamp    | 创建时间 | 否               |
| 15       | updated_at    | timestamp    | 更新时间 | 否               |

---

## system_settings（系统设置表）

| **序号** | **字段名**   | **数据类型** | **备注**   | **是否允许为空** |
| -------- | ------------ | ------------ | ---------- | ---------------- |
| 1        | id           | uuid         | 唯一ID     | 否               |
| 2        | basic        | jsonb        | 基本设置   | 是               |
| 3        | model        | jsonb        | 模型设置   | 是               |
| 4        | security     | jsonb        | 安全设置   | 是               |
| 5        | notification | jsonb        | 通知设置   | 是               |
| 6        | storage      | jsonb        | 云存储设置 | 是               |
| 7        | created_at   | timestamp    | 创建时间   | 否               |
| 8        | updated_at   | timestamp    | 更新时间   | 否               |
| 9        | deleted_at   | timestamp    | 删除时间   | 是               |

---

## provider_configs（厂商配置表）

| **序号** | **字段名**  | **数据类型** | **备注**   | **是否允许为空** |
| -------- | ----------- | ------------ | ---------- | ---------------- |
| 1        | id          | uuid         | 唯一ID     | 否               |
| 2        | user_id     | uuid         | 用户ID     | 否               |
| 3        | name        | varchar(255) | 提供商名称 | 否               |
| 4        | provider    | varchar(50)  | 提供商标识 | 否               |
| 5        | description | text         | 描述       | 是               |
| 6        | api_key     | varchar(255) | API密钥    | 是               |
| 7        | api_base    | varchar(255) | API地址    | 是               |
| 8        | status      | varchar(20)  | 状态       | 是               |
| 9        | created_at  | timestamp    | 创建时间   | 否               |
| 10       | updated_at  | timestamp    | 更新时间   | 否               |
| 11       | deleted_at  | timestamp    | 删除时间   | 是               |

---

## llms（大语言模型表）

| **序号** | **字段名**         | **数据类型** | **备注**       | **是否允许为空** |
| -------- | ------------------ | ------------ | -------------- | ---------------- |
| 1        | id                 | uuid         | 唯一ID         | 否               |
| 2        | user_id            | uuid         | 用户ID         | 否               |
| 3        | name               | varchar(255) | 模型名称       | 否               |
| 4        | description        | text         | 描述           | 是               |
| 5        | provider_config_id | uuid         | 关联厂商配置ID | 是               |
| 6        | model_name         | varchar(255) | 模型标识       | 否               |
| 7        | model_type         | varchar(20)  | 模型类型       | 是               |
| 8        | config             | jsonb        | 关键配置       | 是               |
| 9        | status             | varchar(20)  | 状态           | 是               |
| 10       | created_at         | timestamp    | 创建时间       | 否               |
| 11       | updated_at         | timestamp    | 更新时间       | 否               |
| 12       | deleted_at         | timestamp    | 删除时间       | 是               |

---

## workflows（工作流表）

| **序号** | **字段名** | **数据类型** | **备注** | **是否允许为空** |
| -------- | ---------- | ------------ | -------- | ---------------- |
| 1        | id         | uuid         | 唯一ID   | 否               |
| 2        | user_id    | uuid         | 用户ID   | 否               |
| 3        | name       | varchar(255) | 名称     | 否               |
| 4        | description | varchar(511)| 描述     | 是               |
| 5        | type       | varchar(31)  | 类型     | 否               |
| 6        | status     | varchar(31)  | 状态     | 否               |
| 7        | version    | int          | 版本号   | 否               |
| 8        | config     | jsonb        | 配置     | 否               |
| 9        | data       | jsonb        | 图数据   | 是               |
| 10       | created_at | timestamp    | 创建时间 | 否               |
| 11       | updated_at | timestamp    | 更新时间 | 否               |
| 12       | deleted_at | timestamp    | 删除时间 | 是               |

---

## agents（智能体表）

| **序号** | **字段名**         | **数据类型** | **备注**               | **是否允许为空** |
| -------- | ------------------ | ------------ | ---------------------- | ---------------- |
| 1        | id                 | uuid         | 唯一ID                 | 否               |
| 2        | creator_id         | uuid         | 创建者ID               | 否               |
| 3        | name               | varchar(255) | 名称                   | 否               |
| 4        | description        | text         | 描述信息               | 是               |
| 5        | icon               | varchar(512) | 图标URL                | 是               |
| 6        | system_prompt      | text         | 系统提示词             | 是               |
| 7        | model_provider     | varchar(50)  | 模型提供商             | 否               |
| 8        | model_name         | varchar(100) | 模型名称               | 否               |
| 9        | model_parameters   | jsonb        | 模型参数配置           | 是               |
| 10       | opening_dialogue   | text         | 开场白对话内容         | 是               |
| 11       | suggested_questions | jsonb        | 建议问题列表           | 是               |
| 12       | version            | int          | 版本号                 | 否               |
| 13       | status             | varchar(20)  | 状态（草稿/发布/归档） | 否               |
| 14       | visibility         | varchar(20)  | 可见性                 | 否               |
| 15       | invocation_count   | bigint       | 调用次数统计           | 否               |
| 16       | published_at       | timestamptz  | 发布时间               | 是               |
| 17       | created_at         | timestamp    | 创建时间               | 否               |
| 18       | updated_at         | timestamp    | 更新时间               | 否               |
| 19       | deleted_at         | timestamp    | 删除时间               | 是               |

---

## agent_market（智能体市场表）

| **序号** | **字段名**    | **数据类型** | **备注**      | **是否允许为空** |
| -------- | ------------- | ------------ | ------------- | ---------------- |
| 1        | id            | uuid         | 唯一ID        | 否               |
| 2        | url           | varchar(255) | Agent URL地址 | 否               |
| 3        | name          | varchar(255) | Agent名称     | 否               |
| 4        | description   | text         | Agent描述信息 | 是               |
| 5        | handler_path  | varchar(255) | Handler路径   | 否               |

---

## knowledge_bases（知识库表）

| **序号** | **字段名**                | **数据类型** | **备注**       | **是否允许为空** |
| -------- | ------------------------- | ------------ | -------------- | ---------------- |
| 1        | id                        | uuid         | 唯一ID         | 否               |
| 2        | creator_id                | uuid         | 创建者ID       | 否               |
| 3        | name                      | varchar(255) | 名称           | 否               |
| 4        | description               | text         | 描述           | 是               |
| 5        | chat_model_name           | varchar(255) | 对话模型名称   | 是               |
| 6        | chat_model_provider       | varchar(50)  | 对话模型提供商 | 是               |
| 7        | embedding_model_name      | varchar(255) | 向量模型名称   | 是               |
| 8        | embedding_model_provider  | varchar(50)  | 向量模型提供商 | 是               |
| 9        | embedding_dimension       | integer      | 向量维度       | 否               |
| 10       | storage_type              | varchar(50)  | 存储类型       | 否               |
| 11       | storage_config            | jsonb        | 存储配置       | 是               |
| 12       | document_count            | integer      | 文档数量       | 否               |
| 13       | tags                      | jsonb        | 标签           | 是               |
| 14       | status                    | varchar(20)  | 状态           | 否               |
| 15       | created_at                | timestamp    | 创建时间       | 否               |
| 16       | updated_at                | timestamp    | 更新时间       | 否               |
| 17       | deleted_at                | timestamp    | 删除时间       | 是               |

---

## documents（文档表）

| **序号** | **字段名**    | **数据类型** | **备注**                   | **是否允许为空** |
| -------- | ------------- | ------------ | -------------------------- | ---------------- |
| 1        | id            | uuid         | 唯一ID                     | 否               |
| 2        | kb_id         | uuid         | 归属知识库ID               | 否               |
| 3        | creator_id    | uuid         | 创建者ID                   | 否               |
| 4        | name          | varchar(255) | 文件名                     | 否               |
| 5        | file_type     | varchar(50)  | 文件类型                   | 否               |
| 6        | size          | bigint       | 文件大小(字节)             | 否               |
| 7        | token_count   | integer      | 总Token数消耗统计          | 是               |
| 8        | storage_key   | varchar(512) | S3/OSS存储路径             | 否               |
| 9        | file_hash     | varchar(64)  | SHA256哈希(防重复上传)     | 是               |
| 10       | status        | varchar(20)  | 处理状态                   | 否               |
| 11       | error_message | text         | 错误信息                   | 是               |
| 12       | meta_info     | jsonb        | 解析结果元数据             | 是               |
| 13       | enabled       | boolean      | 是否启用                   | 否               |
| 14       | created_at    | timestamp    | 创建时间                   | 否               |
| 15       | updated_at    | timestamp    | 更新时间                   | 否               |
| 16       | deleted_at    | timestamp    | 删除时间                   | 是               |

---

## document_chunks（文档切片表）

| **序号** | **字段名**  | **数据类型** | **备注**               | **是否允许为空** |
| -------- | ----------- | ------------ | ---------------------- | ---------------- |
| 1        | id          | uuid         | 唯一ID                 | 否               |
| 2        | document_id | uuid         | 关联文档ID             | 否               |
| 3        | kb_id       | uuid         | 归属知识库ID           | 否               |
| 4        | es_id       | varchar(100) | ES中的ID               | 是               |
| 5        | chunk_index | integer      | 切片在原文中的顺序     | 否               |
| 6        | content     | text         | 切片文本内容           | 否               |
| 7        | token_count | integer      | 该切片的Token数        | 是               |
| 8        | meta_info   | jsonb        | 元数据                 | 是               |
| 9        | status      | varchar(20)  | 状态                   | 否               |
| 10       | created_at  | timestamp    | 创建时间               | 否               |
| 11       | updated_at  | timestamp    | 更新时间               | 否               |
| 12       | deleted_at  | timestamp    | 删除时间               | 是               |

---

## tools（工具表）

| **序号** | **字段名**        | **数据类型** | **备注**   | **是否允许为空** |
| -------- | ----------------- | ------------ | ---------- | ---------------- |
| 1        | id                | uuid         | 唯一ID     | 否               |
| 2        | creator_id        | uuid         | 创建者ID   | 否               |
| 3        | name              | varchar(255) | 名称       | 否               |
| 4        | description       | text         | 描述       | 是               |
| 5        | tool_type         | varchar(50)  | 工具类型   | 否               |
| 6        | is_enable         | boolean      | 是否启用   | 是               |
| 7        | parameters_schema | jsonb        | 参数Schema | 是               |
| 8        | mcp_config        | jsonb        | MCP配置    | 是               |
| 9        | created_at        | timestamp    | 创建时间   | 否               |
| 10       | updated_at        | timestamp    | 更新时间   | 否               |
| 11       | deleted_at        | timestamp    | 删除时间   | 是               |

---

## agent_tools（智能体-工具关联表）

| **序号** | **字段名** | **数据类型** | **备注** | **是否允许为空** |
| -------- | ---------- | ------------ | -------- | ---------------- |
| 1        | agent_id   | uuid         | 智能体ID | 否               |
| 2        | tool_id    | uuid         | 工具ID   | 否               |
| 3        | status     | varchar(50)  | 状态     | 是               |
| 4        | created_at | timestamp    | 创建时间 | 是               |

---

## agent_knowledge_bases（智能体-知识库关联表）

| **序号** | **字段名**         | **数据类型**    | **备注** | **是否允许为空** |
| -------- | ------------------ | --------------- | -------- | ---------------- |
| 1        | agent_id           | bigint unsigned | 智能体ID | 否               |
| 2        | knowledge_base_id  | bigint unsigned | 知识库ID | 否               |
| 3        | status             | varchar(20)     | 状态     | 否               |

---

## agent_agents（智能体-市场关联表）

| **序号** | **字段名**      | **数据类型** | **备注**       | **是否允许为空** |
| -------- | --------------- | ------------ | -------------- | ---------------- |
| 1        | agent_id        | uuid         | 智能体ID       | 是               |
| 2        | agent_market_id | uuid         | 市场智能体ID   | 是               |

---

## agent_workflows（智能体-工作流关联表）

| **序号** | **字段名**        | **数据类型** | **备注** | **是否允许为空** |
| -------- | ----------------- | ------------ | -------- | ---------------- |
| 1        | agent_id          | uuid         | 智能体ID | 否               |
| 2        | workflow_id       | uuid         | 工作流ID | 否               |
| 3        | is_default        | boolean      | 是否默认 | 否               |
| 4        | trigger_condition | varchar(255) | 触发条件 | 是               |
| 5        | priority          | int          | 优先级   | 否               |
| 6        | status            | varchar(20)  | 状态     | 否               |
| 7        | created_at        | timestamptz  | 创建时间 | 否               |

---

## chat_sessions（聊天会话表）

| **序号** | **字段名** | **数据类型** | **备注** | **是否允许为空** |
| -------- | ---------- | ------------ | -------- | ---------------- |
| 1        | id         | uuid         | 唯一ID   | 否               |
| 2        | agent_id   | uuid         | 智能体ID | 否               |
| 3        | user_id    | uuid         | 用户ID   | 否               |
| 4        | title      | varchar(255) | 标题     | 是               |
| 5        | created_at | timestamp    | 创建时间 | 否               |
| 6        | updated_at | timestamp    | 更新时间 | 否               |
| 7        | deleted_at | timestamp    | 删除时间 | 是               |

---

## chat_messages（聊天消息表）

| **序号** | **字段名** | **数据类型** | **备注**                    | **是否允许为空** |
| -------- | ---------- | ------------ | --------------------------- | ---------------- |
| 1        | id         | uuid         | 唯一ID                      | 否               |
| 2        | session_id | uuid         | 会话ID                      | 否               |
| 3        | role       | varchar(50)  | 角色(user/assistant/system) | 否               |
| 4        | content    | text         | 消息内容                    | 否               |
| 5        | created_at | timestamp    | 创建时间                    | 否               |
| 6        | updated_at | timestamp    | 更新时间                    | 否               |
| 7        | deleted_at | timestamp    | 删除时间                    | 是               |
