# Conventional Commits 规范参考

## 规范概述

Conventional Commits 是一种用于给提交信息增加人机可读含义的规范。
规范地址：https://www.conventionalcommits.org/zh-hans/v1.0.0/

## 提交信息结构
<type><description>
[optional body]
[optional footer(s)]

## 类型（Type）定义

| 类型 | 含义 | 使用场景 |
|------|------|----------|
| **feat** | 新功能 | 新增功能特性，对应`MINOR`版本号升级 |
| **fix** | Bug修复 | 修复代码缺陷，对应`PATCH`版本号升级 |
| **docs** | 文档 | 仅文档变更（README、注释等） |
| **style** | 格式 | 不影响代码含义的变更（空格、分号、缩进等） |
| **refactor** | 重构 | 既不修复bug也不添加功能的代码变更 |
| **perf** | 性能优化 | 提升性能的代码变更 |
| **test** | 测试 | 添加或修正测试代码 |
| **chore** | 构建/工具 | 构建过程、辅助工具、依赖更新的变更 |
| **ci** | 持续集成 | CI配置、脚本变更（GitHub Actions等） |
| **build** | 构建 | 影响构建系统或外部依赖的变更（webpack、npm等） |
| **revert** | 回滚 | 撤销之前的提交 |

## 范围（Scope）

- **可选**：但建议在中大型项目中使用
- **格式**：括号包围，如 `feat(auth):`、 `fix(api):`
- **命名**：使用模块/组件名，如 `auth`, `api`, `ui`, `docs`
- **多范围**：使用 `*` 表示多个模块，或用逗号分隔（不推荐）

## 描述（Description）

- **语气**：使用祈使句，现在时（"change"而非"changed"或"changes"）
- **首字母**：小写（如 "add" 而非 "Add"）
- **结尾**：不加句号（.）
- **长度**：不超过50个字符（建议）

✅ **正确示例**：
feat(auth): add OAuth2 login support
fix(api): resolve user data leak in public endpoint

❌ **错误示例**：
feat(auth): Added OAuth2 login support.  // 使用了过去时，加了句号
feat(auth): Add OAuth2 login support     // 首字母大写

## 正文（Body）

- **格式**：每行不超过72字符（便于git log阅读）
- **内容**：详细说明变更的动机、对比修改前后的行为
- **空行**：描述与正文之间必须空一行

示例：
fix(auth): correct password validation logic
Previously, the password validation only checked length >= 6.
Now it requires at least one uppercase, one lowercase, and one number.
This aligns with the security requirements in SEC-202.

## 页脚（Footer）

### Breaking Changes（破坏性变更）

任何破坏性变更都必须在页脚标注，以 `BREAKING CHANGE:` 开头：
feat(api): change user authentication endpoint
BREAKING CHANGE: /api/v1/login is removed, use /api/v2/auth instead.
Migration guide: https://docs.example.com/migration/v2

也可在类型后加 `!` 标记：
feat(api)!: remove deprecated user endpoints

### 关联 Issue

- **Closes**: 修复了某个Issue（自动关闭）
- **Fixes**: 同上
- **Refs**: 引用某个Issue（不关闭）
- **Relates to**: 相关Issue
fix(auth): resolve memory leak in session handling
Closes #234
Refs #156

## 完整示例

### 新功能
feat(auth): implement JWT token refresh mechanism
Add automatic token refresh 5 minutes before expiry.
Implements sliding session for better UX.
Closes #456

### Bug修复
fix(api): handle null pointer in user search
Check for nil user object before accessing Email field.
Add unit test to cover edge case.
Fixes #789

### 破坏性变更
feat!: drop support for Node.js 14
BREAKING CHANGE: Minimum Node.js version is now 16.0.0
This allows us to use native fetch API and ES2022 features.

### 多范围变更
refactor(*): restructure project folders
Move all business logic from handlers to services.
Update import paths across the codebase.

## 验证检查清单

提交前确认：
- [ ] 类型正确且在允许列表中
- [ ] 描述使用祈使句、小写开头、无句号
- [ ] 如有正文，每行不超过72字符
- [ ] 破坏性变更已标记 BREAKING CHANGE 或使用 !
- [ ] 关联的Issue已正确引用