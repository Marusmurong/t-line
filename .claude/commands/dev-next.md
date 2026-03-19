# /dev-next — 推进到下一阶段

完成当前阶段的验收检查，推进工作流到下一阶段。

## 执行步骤

1. **获取当前阶段**
   ```bash
   bash .workflow/workflow.sh stage
   ```

2. **执行阶段验收检查**
   - research: 需求文档完整性检查
   - design: 架构文档 & 接口定义检查
   - plan: 任务清单 & 依赖关系检查
   - implement: 代码完成度 & 编译检查
   - review: 代码审查通过 & 安全扫描
   - test: 覆盖率 ≥ 80% & 全部测试通过
   - deploy: 部署成功 & 健康检查通过
   - verify: 回归测试通过 & 验收确认

3. **推进阶段**
   ```bash
   bash .workflow/workflow.sh next
   ```

4. **调度下一阶段代理**
   - 根据阶段定义自动调度对应代理

## 用法

```
/dev-next
/dev-next --force   # 跳过验收检查（谨慎使用）
```
