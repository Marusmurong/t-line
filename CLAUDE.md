# Your Project

## 项目概述

本项目使用 ZC 工作流引擎进行开发管理，通过多代理协作完成从需求到部署的全流程。

## 工作流

### 启动开发
```bash
/dev-start <需求描述>
```

### 查看状态
```bash
/dev-status
```

### 推进阶段
```bash
/dev-next
```

### 调度代理
```bash
/dev-agent <agent-name> <task>
```

## 阶段定义

1. **research** — 需求研究 & 可行性分析
2. **design** — 架构设计 & 技术方案
3. **plan** — 任务拆解 & 排期
4. **implement** — 编码实现（TDD）
5. **review** — 代码审查 & 交叉审计
6. **test** — 集成测试 & E2E
7. **deploy** — 构建 & 部署
8. **verify** — 上线验证 & 回归

## 代理

| 代理 | 配置 | 职责 |
|------|------|------|
| PM | `.workflow/agents/pm.yaml` | 需求分析、任务拆解、进度跟踪 |
| Architect | `.workflow/agents/architect.yaml` | 系统设计、技术选型 |
| Tester | `.workflow/agents/tester.yaml` | 测试策略、覆盖率验证、Codex 交叉审计 |
| DevOps | `.workflow/agents/devops.yaml` | 构建、部署、监控 |

## 开发规范

- **TDD**: 测试先行，覆盖率 ≥ 80%
- **不可变**: 创建新对象，不修改已有对象
- **小文件**: 200-400 行，最多 800 行
- **安全优先**: 无硬编码密钥，输入验证，参数化查询

## 文件保护

通过 `.claude/hooks/pretooluse-file-guard.sh` 阻止对以下文件的意外修改：
- `.env` / 凭证文件
- `node_modules/` / `.git/`
- `dist/` / `build/`
- `.workflow/active/`（运行时状态）
