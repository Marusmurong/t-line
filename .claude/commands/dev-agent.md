# /dev-agent — 调度指定代理

手动调度指定代理执行任务。

## 执行步骤

1. **解析参数**
   - 代理名称（必需）
   - 任务描述（可选）

2. **验证代理配置**
   ```bash
   bash .workflow/workflow.sh dispatch <agent-name>
   ```

3. **启动代理**
   - 加载代理 YAML 配置
   - 注入当前工作流上下文
   - 执行代理任务

4. **等待结果**
   - 收集代理输出
   - 更新工作流状态

## 可用代理

| 代理 | 说明 |
|------|------|
| pm | 项目经理 — 需求分析、任务拆解 |
| architect | 架构师 — 系统设计、技术选型 |
| tester | 测试 — 测试编写、覆盖率验证 |
| devops | 运维 — 构建、部署、监控 |
| your-agent | 自定义 — 按业务定制 |

## 用法

```
/dev-agent pm "分析新的用户反馈需求"
/dev-agent architect "评审数据库 Schema 设计"
/dev-agent tester "为登录模块编写集成测试"
/dev-agent devops "部署 v1.2.0 到预发环境"
```
