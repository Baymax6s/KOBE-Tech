#!/usr/bin/env bash
set -euo pipefail

FILE=$(jq -r '.tool_input.file_path // empty')

case "$FILE" in
  */app/src/api/generated/*)
    cat >&2 <<'EOF'
Refusing to edit generated API client.

`app/src/api/generated/` は swagger-typescript-api が `api/swagger/openapi.yml` から
生成しているコントラクト層です。手で編集すると次の再生成で消えます。
変更したい場合は Go 側のハンドラ／DTO を直し、`make swagger` → `npm --prefix app run generate:api` で再生成してください。
EOF
    exit 2
    ;;
esac
exit 0
