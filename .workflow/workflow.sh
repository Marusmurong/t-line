#!/usr/bin/env bash
# ============================================================
# Workflow Engine — ZC 版本（可自定义阶段）
# ============================================================
set -euo pipefail

WORKFLOW_DIR="$(cd "$(dirname "$0")" && pwd)"
ACTIVE_DIR="${WORKFLOW_DIR}/active"
AGENTS_DIR="${WORKFLOW_DIR}/agents"

# ── 阶段定义（按顺序执行） ──────────────────────────────
STAGES=(
  "research"      # 需求研究 & 可行性分析
  "design"        # 架构设计 & 技术方案
  "plan"          # 任务拆解 & 排期
  "implement"     # 编码实现（TDD）
  "review"        # 代码审查 & 交叉审计
  "test"          # 集成测试 & E2E
  "deploy"        # 构建 & 部署
  "verify"        # 上线验证 & 回归
)

# ── 状态管理 ─────────────────────────────────────────────
init_workflow() {
  mkdir -p "${ACTIVE_DIR}"
  echo "research" > "${ACTIVE_DIR}/current_stage"
  echo "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "${ACTIVE_DIR}/started_at"
  echo "initialized" > "${ACTIVE_DIR}/status"
  echo "[workflow] 初始化完成，当前阶段: research"
}

get_current_stage() {
  cat "${ACTIVE_DIR}/current_stage" 2>/dev/null || echo "none"
}

get_status() {
  cat "${ACTIVE_DIR}/status" 2>/dev/null || echo "unknown"
}

advance_stage() {
  local current
  current="$(get_current_stage)"
  local next=""

  for i in "${!STAGES[@]}"; do
    if [[ "${STAGES[$i]}" == "$current" ]]; then
      local next_idx=$((i + 1))
      if [[ $next_idx -lt ${#STAGES[@]} ]]; then
        next="${STAGES[$next_idx]}"
      fi
      break
    fi
  done

  if [[ -z "$next" ]]; then
    echo "completed" > "${ACTIVE_DIR}/status"
    echo "[workflow] 所有阶段已完成"
    return 0
  fi

  echo "$next" > "${ACTIVE_DIR}/current_stage"
  echo "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "${ACTIVE_DIR}/stage_started_at"
  echo "in_progress" > "${ACTIVE_DIR}/status"
  echo "[workflow] 推进到阶段: $next"
}

show_status() {
  echo "================================"
  echo "  Workflow Status"
  echo "================================"
  echo "  当前阶段: $(get_current_stage)"
  echo "  状态:     $(get_status)"
  echo "  启动时间: $(cat "${ACTIVE_DIR}/started_at" 2>/dev/null || echo 'N/A')"
  echo "================================"
  echo "  阶段列表:"
  local current
  current="$(get_current_stage)"
  for stage in "${STAGES[@]}"; do
    if [[ "$stage" == "$current" ]]; then
      echo "    ➤ $stage  ← 当前"
    else
      echo "      $stage"
    fi
  done
  echo "================================"
}

# ── 代理调度 ─────────────────────────────────────────────
dispatch_agent() {
  local agent_name="$1"
  local agent_file="${AGENTS_DIR}/${agent_name}.yaml"

  if [[ ! -f "$agent_file" ]]; then
    echo "[error] 代理配置不存在: ${agent_file}"
    return 1
  fi

  echo "[dispatch] 调度代理: ${agent_name}"
  # 实际调度逻辑由 Claude Code 或外部编排器处理
}

# ── 入口 ─────────────────────────────────────────────────
case "${1:-help}" in
  init)     init_workflow ;;
  status)   show_status ;;
  next)     advance_stage ;;
  stage)    get_current_stage ;;
  dispatch) dispatch_agent "${2:?需要指定代理名称}" ;;
  help)
    echo "用法: workflow.sh <command>"
    echo "  init      初始化工作流"
    echo "  status    查看当前状态"
    echo "  next      推进到下一阶段"
    echo "  stage     获取当前阶段名"
    echo "  dispatch  调度代理 (需要代理名称)"
    ;;
  *)
    echo "[error] 未知命令: $1"
    exit 1
    ;;
esac
