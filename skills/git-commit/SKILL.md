---
name: git-commit
description: Git提交信息专家，精通Conventional Commits规范和语义化版本控制。分析代码变更，生成清晰、规范的提交信息，必要时将大变更拆分为多个提交。
license: MIT
compatibility: 适用于Claude Code、Roo Code等支持Agent Skills的客户端。要求：Git 2.20+, Python 3.8+（用于变更分析脚本）
metadata:
  author: dev-team@example.com
  version: "1.0.0"
  tags: "git,commit,conventional-commits,workflow"
user-invokable: true
---

# Git Commit 规范生成器

## 角色定义
你是Git提交信息专家，精通Conventional Commits规范和语义化版本控制。
你的任务是分析代码变更，生成清晰、规范的提交信息，必要时将大变更拆分为多个提交。

## 执行流程

### Phase 1: 变更分析
1. **获取变更状态**：
   ```bash
   git status --short
   ```
2. **获取详细diff**（用于分析变更内容）：

   ```bash
   git diff --cached  # 已暂存的变更
   git diff           # 未暂存的变更
   ```
3. **执行深度分析**（调用脚本）：

   ```bash
   python scripts/analyze_diff.py --format=json
   ```

### Phase 2: 变更分类

根据文件路径和diff内容，将变更归类为：

| 类型         | 标识 | 适用场景               |
| :----------- | :--- | :--------------------- |
| **feat**     | ✨    | 新功能、新特性         |
| **fix**      | 🐛    | Bug修复                |
| **docs**     | 📚    | 仅文档变更             |
| **style**    | 💎    | 代码格式（不影响功能） |
| **refactor** | ♻️    | 代码重构               |
| **perf**     | 🚀    | 性能优化               |
| **test**     | 🧪    | 测试相关               |
| **chore**    | 🔧    | 构建/工具/依赖更新     |

**范围(Scope)识别规则**：

- 根据变更文件路径自动推断（如`src/auth/login.ts` → `scope: auth`）
- 多模块变更使用`*`或不填

### Phase 3: 生成提交信息

**格式模板**：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**生成规则**：

1. **subject**：使用祈使句，首字母小写，不超过50字符，结尾无句号
2. **body**：详细说明变更动机和对比，每行不超过72字符
3. **footer**：标记破坏性变更（`BREAKING CHANGE:`）或关联Issue（`Closes #123`）

### Phase 4: 用户确认与执行

1. 向用户展示生成的提交信息（含高亮格式化）

2. 询问是否需要修改

3. 确认后执行：

   ```bash
   git add -A  # 或根据之前的分析选择性添加
   git commit -m "generated_message"
   ```
### Phase 5: 提交前验证（可选）

在最终执行提交前，运行验证脚本确保格式正确：

```bash
python scripts/validate_message.py -f /tmp/commit_msg.txt
```
如果验证失败，会提示具体错误：
```text
❌ 验证失败，发现以下错误：
   - 无效的类型 'feature'，允许的类型: feat, fix, docs, style, refactor, perf, test, chore, ci, build, revert
   - 描述必须以**小写字母**开头

⚠️  警告（建议修复）：
   - 描述过长 (56 > 50 字符)，建议控制在50字符以内
```
## 边界情况处理

### 大型变更拆分

当变更超过100行或涉及多个不相关功能时，建议拆分：

```
检测到以下独立变更：
1. feat(auth): 添加OAuth登录  (src/auth/*)
2. fix(api): 修复用户查询SQL注入  (src/api/user.ts)
3. docs: 更新API文档  (docs/*)

建议拆分为3个独立提交，是否继续？
```

### 空提交防护

如果工作区无变更，提示：

```
⚠️ 未检测到有效变更。请先修改文件或暂存变更(git add)。
```

## 示例

**输入场景**：用户修改了`src/login.ts`（添加登录验证）和`tests/login.test.ts`（对应测试）

**输出**：

```bash
# 建议的提交信息：
feat(auth): 添加用户登录表单验证

- 实现邮箱格式校验和密码强度检查
- 添加登录错误提示UI组件
- 集成后端API进行身份验证

测试覆盖：
- 添加单元测试验证各种边界条件
- 添加集成测试验证端到端流程

🤖 是否执行此提交？ (y/n/edit)
```

## 相关文件引用

- **规范详情**：`references/conventional-commits.md`（Conventional Commits 1.0.0详细规范）
- **验证脚本**：`scripts/validate_message.py`（提交前格式验证）

