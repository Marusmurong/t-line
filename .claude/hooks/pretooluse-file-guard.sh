#!/usr/bin/env bash
# ============================================================
# PreToolUse Hook: File Guard
# 防止对受保护文件/目录的意外修改
# ============================================================
set -euo pipefail

# ── 受保护路径（正则匹配） ───────────────────────────────
PROTECTED_PATTERNS=(
  "\.env$"
  "\.env\."
  "credentials"
  "secrets?"
  "\.pem$"
  "\.key$"
  "node_modules/"
  "\.git/"
  "dist/"
  "build/"
  "\.workflow/active/"
)

# ── 从 stdin 读取工具调用信息 ────────────────────────────
INPUT=$(cat)

# 提取文件路径（适配 Write / Edit 工具的 file_path 参数）
FILE_PATH=$(echo "$INPUT" | grep -oP '"file_path"\s*:\s*"([^"]+)"' | head -1 | sed 's/.*: *"//;s/"$//' || true)

if [[ -z "$FILE_PATH" ]]; then
  # 非文件操作，放行
  exit 0
fi

# ── 检查是否匹配受保护路径 ───────────────────────────────
for pattern in "${PROTECTED_PATTERNS[@]}"; do
  if echo "$FILE_PATH" | grep -qE "$pattern"; then
    echo "[file-guard] BLOCKED: 尝试修改受保护文件: $FILE_PATH"
    echo "[file-guard] 匹配规则: $pattern"
    echo "[file-guard] 如需修改，请先确认并手动操作。"
    exit 1
  fi
done

# 放行
exit 0
