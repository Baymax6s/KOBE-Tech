#!/usr/bin/env bash
set -uo pipefail

FILE=$(jq -r '.tool_input.file_path // empty')

[[ -z "$FILE" ]] && exit 0
[[ "$FILE" != *"/app/"* ]] && exit 0
[[ "$FILE" == *"/app/src/api/generated/"* ]] && exit 0

case "$FILE" in
  *.vue|*.ts|*.tsx|*.js|*.mjs|*.cjs|*.json|*.css|*.scss|*.md)
    cd "${CLAUDE_PROJECT_DIR}/app"
    npx --no-install prettier --write --log-level=warn "$FILE" >&2 || true
    case "$FILE" in
      *.vue|*.ts|*.tsx|*.js|*.mjs|*.cjs)
        npx --no-install eslint --fix "$FILE" >&2 || true
        ;;
    esac
    ;;
esac
exit 0
