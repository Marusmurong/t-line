# ============================================================
# Makefile — 项目快捷命令
# ============================================================

.PHONY: help init status next review clean

WORKFLOW := bash .workflow/workflow.sh
CODEX_REVIEW := bash scripts/codex_cross_review.sh

help: ## 显示帮助
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

init: ## 初始化工作流
	$(WORKFLOW) init

status: ## 查看工作流状态
	$(WORKFLOW) status

next: ## 推进到下一阶段
	$(WORKFLOW) next

review: ## 运行 Codex 交叉审计
	$(CODEX_REVIEW) review

clean: ## 清理运行时状态
	rm -rf .workflow/active/*
	@echo "运行时状态已清理"

# ── 代理调度 ─────────────────────────────────────────────

dispatch-pm: ## 调度 PM 代理
	$(WORKFLOW) dispatch pm

dispatch-architect: ## 调度架构师代理
	$(WORKFLOW) dispatch architect

dispatch-tester: ## 调度测试代理
	$(WORKFLOW) dispatch tester

dispatch-devops: ## 调度运维代理
	$(WORKFLOW) dispatch devops
