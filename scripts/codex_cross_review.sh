#!/usr/bin/env bash
# ZC 项目 — Codex 交叉审计脚本
# 用法: codex_cross_review.sh [workflow_id]
# 用 OpenAI Codex 对 Claude 写的 Go 代码做独立安全/逻辑审查
# 支持: codex exec review / codex -q / OpenAI API 三级降级

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
WORKFLOW_DIR="$PROJECT_DIR/.workflow"
WF_ID="${1:-}"
REVIEW_TARGET="${2:-code}"  # code（默认）或 architecture

# ============================================================
# 0. 加载 API Key（多个来源按优先级）
# ============================================================
if [[ -z "${OPENAI_API_KEY:-}" ]]; then
  # 来源1: 项目 .env
  [[ -f "$PROJECT_DIR/.env" ]] && source "$PROJECT_DIR/.env" 2>/dev/null || true
  # 来源2: autogenstudio 配置
  [[ -z "${OPENAI_API_KEY:-}" && -f "$HOME/.autogenstudio/temp_env_vars.env" ]] && source "$HOME/.autogenstudio/temp_env_vars.env" 2>/dev/null || true
  # 来源3: codex 自身配置
  [[ -z "${OPENAI_API_KEY:-}" && -f "$HOME/.codex/config.json" ]] && {
    OPENAI_API_KEY=$(jq -r '.apiKey // empty' "$HOME/.codex/config.json" 2>/dev/null) || true
  }
  # 来源4: shell profile
  [[ -z "${OPENAI_API_KEY:-}" && -f "$HOME/.zshrc" ]] && {
    OPENAI_API_KEY=$(grep "^export OPENAI_API_KEY=" "$HOME/.zshrc" 2>/dev/null | tail -1 | sed 's/^export OPENAI_API_KEY=//' | tr -d '"' | tr -d "'") || true
  }
fi

if [[ -z "${OPENAI_API_KEY:-}" ]]; then
  echo "ERROR: OPENAI_API_KEY 未设置。请在 .env 或环境变量中配置" >&2
  exit 1
fi
export OPENAI_API_KEY
export OPENAI_BASE_URL="${OPENAI_BASE_URL:-https://api.openai.com/v1}"

# ============================================================
# 1. 确定工作流
# ============================================================
if [[ -z "$WF_ID" ]]; then
  ACTIVE_DIR="$WORKFLOW_DIR/active"
  ACTIVES=()
  for f in "$ACTIVE_DIR"/WF-*.json; do
    [[ -f "$f" ]] || continue
    STATUS=$(jq -r '.phase_status' "$f" 2>/dev/null)
    [[ "$STATUS" == "completed" || "$STATUS" == "cancelled" ]] && continue
    ACTIVES+=("$f")
  done
  if [[ ${#ACTIVES[@]} -eq 0 ]]; then
    echo "ERROR: 没有活跃工作流" >&2; exit 1
  elif [[ ${#ACTIVES[@]} -eq 1 ]]; then
    WF_ID=$(jq -r '.workflow_id' "${ACTIVES[0]}")
  else
    echo "ERROR: 多个活跃工作流，请指定 workflow_id:" >&2
    for f in "${ACTIVES[@]}"; do
      jq -r '"  \(.workflow_id) — \(.feature)"' "$f"
    done
    exit 1
  fi
fi

STATE_FILE="$WORKFLOW_DIR/active/${WF_ID}.json"
[[ -f "$STATE_FILE" ]] || { echo "ERROR: 工作流不存在: $WF_ID" >&2; exit 1; }

BRANCH=$(jq -r '.branch' "$STATE_FILE")
FEATURE=$(jq -r '.feature' "$STATE_FILE")
# 报告文件名区分：architecture 或 code
if [[ "$REVIEW_TARGET" == "architecture" ]]; then
  REPORT="$PROJECT_DIR/docs/CODEX_ARCH_REVIEW_${WF_ID}.md"
  REPORT_TYPE="架构审计"
else
  REPORT="$PROJECT_DIR/docs/CODEX_REVIEW_${WF_ID}.md"
  REPORT_TYPE="代码审计"
fi

echo "=== ZC Codex 交叉审计 (${REPORT_TYPE}) ==="
echo "工作流: $WF_ID | 功能: $FEATURE | 分支: $BRANCH"
echo ""

# ============================================================
# 2. 检查 Codex CLI
# ============================================================
if command -v codex &>/dev/null; then
  CODEX_CMD="codex"
elif command -v npx &>/dev/null; then
  CODEX_CMD="npx -y @openai/codex"
else
  CODEX_CMD=""
fi

# ============================================================
# 3. 审计 prompt（根据目标不同）
# ============================================================
if [[ "$REVIEW_TARGET" == "architecture" ]]; then
  # 架构阶段审计：审查设计文档和 TODO 清单
  REVIEW_PROMPT="你是一位资深系统架构审计专家，正在审查一个高并发 WhatsApp 群发平台（日请求5000万+，50-60万并发）的架构设计。

审计目标：docs/ 目录下的架构文档（ARCH_*.md、PRD_*.md）和 TODO 清单。

审计重点：
1. 架构合理性：服务拆分是否合理、模块耦合度、单点故障
2. 并发设计：50万+并发下的瓶颈、Kafka 分区策略、连接池规划
3. 数据一致性：分布式事务、幂等设计、消息重复处理
4. 安全设计：认证方案、密钥管理、API 鉴权、数据脱敏
5. 可运维性：监控方案、灰度发布、回滚策略
6. TODO 清单：任务粒度是否合理、优先级是否正确、遗漏项
7. 接口设计：API 规范、请求/响应格式、错误码设计

输出格式（Markdown）：
# Codex 架构审计报告 — ZC $WF_ID
## $FEATURE

### 🔴 CRITICAL (架构缺陷，必须修改)
- **文档:章节** — 描述 → 修改建议 → 指派: architect

### 🟡 HIGH (设计风险)
- **文档:章节** — 描述 → 修改建议 → 指派: architect

### 🔵 MEDIUM (优化建议)
...

### ⚪ LOW (注意事项)
...

### TODO 清单审查
- 遗漏的任务
- 粒度不合理的任务
- 优先级需调整的任务

### 总结
- 架构总体评分 (1-10)
- 是否建议继续到开发阶段"
else
  # 代码阶段审计
  REVIEW_PROMPT="你是一位资深 Go 安全审计专家，正在审查一个高并发 WhatsApp 群发平台（日请求5000万+，50-60万并发）。

审计重点（Go 项目特有）：
1. 并发安全：goroutine 泄漏、data race、channel 死锁、sync.Mutex 使用
2. 错误处理：error 是否被忽略、panic 恢复、defer 顺序
3. 内存管理：大切片未释放、string/[]byte 转换、sync.Pool 使用
4. 安全漏洞：SQL注入（sqlx参数化）、命令注入、密钥硬编码
5. Kafka 消费：offset 提交、重复消费、消息丢失
6. Redis 操作：连接泄漏、pipeline 使用、key 过期策略
7. HTTP 安全：CORS、Rate Limit、Token 校验
8. 配置问题：硬编码值、调试日志未关、环境差异

重要：审计结果中每个问题必须标注 → 指派: <agent名> (写这段代码的 Agent)

输出格式（Markdown）：
# Codex 交叉审计报告 — ZC $WF_ID
## $FEATURE

### 🔴 CRITICAL (立即修复)
- **文件:行号** — 描述 → 修复方案 → 指派: <agent>

### 🟡 HIGH (本次迭代修复)
- **文件:行号** — 描述 → 修复方案 → 指派: <agent>

### 🔵 MEDIUM (建议优化)
...

### ⚪ LOW (注意事项)
...

### 总结
- 总问题数 / 各级别数量
- 并发安全评估
- 是否建议合并"
fi

# ============================================================
# 4. 获取审计内容
# ============================================================
cd "$PROJECT_DIR"

REVIEW_CONTENT=""
if [[ "$REVIEW_TARGET" == "architecture" ]]; then
  # 架构审计：读取 docs/ 下的文档
  echo "收集架构文档..."
  for doc in "$PROJECT_DIR"/docs/{PRD_,ARCH_}*.md; do
    [[ -f "$doc" ]] || continue
    echo "  + $(basename "$doc")"
    REVIEW_CONTENT+="
=== $(basename "$doc") ===
$(head -c 8000 "$doc")
"
  done
  [[ -z "$REVIEW_CONTENT" ]] && { echo "ERROR: docs/ 下没有 PRD/ARCH 文档" >&2; exit 1; }
  USER_MSG="请审查以下架构文档：

${REVIEW_CONTENT}"
else
  # 代码审计：获取 git diff
  DIFF=""
  if git rev-parse --is-inside-work-tree &>/dev/null 2>&1; then
    BASE_BRANCH=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@' || echo "main")
    DIFF=$(git diff "${BASE_BRANCH}...HEAD" 2>/dev/null || git diff HEAD~10 2>/dev/null || git diff HEAD 2>/dev/null || echo "")
  fi
  if [[ -z "$DIFF" ]]; then
    # 没有 diff，直接读所有 .go 文件摘要
    echo "WARN: 无法获取 diff，读取 Go 源文件..." >&2
    DIFF=$(find "$PROJECT_DIR" -name "*.go" -not -path "*/vendor/*" -not -name "*_test.go" -exec head -50 {} + 2>/dev/null | head -c 15000)
  fi
  # 截取
  DIFF_TRUNCATED="${DIFF:0:15000}"
  [[ ${#DIFF} -gt 15000 ]] && DIFF_TRUNCATED+="
... (已截断，共 ${#DIFF} 字符)"
  USER_MSG="请审查以下 Go 代码变更:

\`\`\`diff
${DIFF_TRUNCATED}
\`\`\`"
fi

# ============================================================
# 5. 执行审查（优先 API 直接调用，最可靠）
# ============================================================
CODEX_OUTPUT=""

# 方式1: OpenAI API 直接调用（最可靠，不依赖 codex CLI 的交互行为）
echo "[1/2] 通过 OpenAI API 审计..."

# 模型选择（可通过环境变量覆盖）
MODEL="${CODEX_MODEL:-gpt-5.4}"
echo "使用模型: $MODEL"

CODEX_OUTPUT=$(curl -s --max-time 120 "${OPENAI_BASE_URL}/chat/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d "$(jq -n \
    --arg prompt "$REVIEW_PROMPT" \
    --arg user_msg "$USER_MSG" \
    --arg model "$MODEL" \
    '{
      "model": $model,
      "messages": [
        {"role": "system", "content": $prompt},
        {"role": "user", "content": $user_msg}
      ],
      "temperature": 0.2,
      "max_tokens": 4096
    }')" 2>/dev/null | jq -r '.choices[0].message.content // empty')

if [[ -n "$CODEX_OUTPUT" ]]; then
  echo "✅ OpenAI API 审计完成"
fi

# 方式2: codex exec（如果 API 失败且 codex 可用）
if [[ -z "$CODEX_OUTPUT" && -n "$CODEX_CMD" ]]; then
  echo "[2/2] 降级到 codex exec..."
  CODEX_OUTPUT=$($CODEX_CMD --full-auto exec "$REVIEW_PROMPT" 2>/dev/null) || CODEX_OUTPUT=""
  [[ -n "$CODEX_OUTPUT" ]] && echo "✅ codex exec 审计完成"
fi

# 方式3: 全部失败
if [[ -z "$CODEX_OUTPUT" ]]; then
  CODEX_OUTPUT="# Codex 交叉审计报告 — ZC $WF_ID

## ⚠️ 自动审计失败

无法完成 Codex 自动审计。原因可能是：
- codex CLI 未安装（npm install -g @openai/codex）
- OPENAI_API_KEY 未设置
- 网络连接问题

## 手动审计方法

### 方法A: VSCode Codex 插件
1. 在 VSCode 中打开项目
2. 打开 Codex 面板
3. 输入 \`/review\` 执行审查
4. 将结果复制到 docs/CODEX_REVIEW_${WF_ID}.md

### 方法B: Codex CLI
\`\`\`bash
cd $PROJECT_DIR
codex
# 进入后输入:
/review
\`\`\`

### 方法C: 终端命令
\`\`\`bash
codex exec review uncommitted
\`\`\`"
fi

# ============================================================
# 6. 写入报告
# ============================================================
mkdir -p "$PROJECT_DIR/docs"

cat > "$REPORT" <<EOF
---
workflow: $WF_ID
feature: $FEATURE
branch: $BRANCH
reviewer: OpenAI Codex (GPT-5.4)
cross_audit: true
auditor_model: Claude → Codex 交叉审计
date: $(date -u +%Y-%m-%dT%H:%M:%SZ)
---

$CODEX_OUTPUT
EOF

echo ""
echo "==============================="
echo "审计报告: $REPORT"
echo "==============================="
echo ""
echo "下一步操作:"
echo "  cat $REPORT                      # 查看报告"
echo "  # CRITICAL/HIGH 问题 → 通知 Agent 修复"
echo "  # 无严重问题 → PM 执行 /dev-next"
