# /dev-start — 启动开发工作流

初始化并启动项目开发工作流。

## 执行步骤

1. **初始化工作流**
   ```bash
   bash .workflow/workflow.sh init
   ```

2. **加载代理配置**
   - 读取 `.workflow/agents/` 下所有 YAML 配置
   - 验证必需代理（pm, architect, tester, devops）存在

3. **启动 PM 代理**
   - 分析需求输入
   - 生成任务拆解
   - 输出初始计划

4. **显示状态**
   ```bash
   bash .workflow/workflow.sh status
   ```

## 用法

```
/dev-start [需求描述或 PRD 文件路径]
```

## 示例

```
/dev-start 实现用户登录功能，支持手机号验证码登录
/dev-start docs/prd/feature-x.md
```
