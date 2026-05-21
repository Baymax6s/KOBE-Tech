# KOBE-Tech AGENTS.md

このリポジトリで作業する AI コーディングエージェント（Claude Code / Codex / OpenCode 等）共通の指示です。

## リポジトリ構成

- `app/` — Vue 3 + Vuetify + Pinia + TypeScript のフロントエンド（Vite）
- `api/` — Go 1.25 + Gin + Postgres のバックエンド（swag で OpenAPI 生成）
- `app/src/api/generated/` — `api/swagger/openapi.yml` から生成される TypeScript の API クライアント。**手で編集しない**

## コントラクト層と再生成

API の型定義は `api/swagger/openapi.yml` を真のソースとし、`app/src/api/generated/` へ自動生成します。フローは:

1. Go のハンドラ／DTO を編集
2. `cd api && make swagger` で OpenAPI を更新
3. `cd app && npm run generate:api` で TypeScript クライアントを再生成

`app/src/api/generated/**` を手で編集することは禁止です。

## ディレクトリ別の追加ルール

- `app/` 配下の編集: [`app/AGENTS.md`](app/AGENTS.md) を参照（Vuetify ルール・フォーム送信ルール）
- `api/` 配下の編集: [`api/AGENTS.md`](api/AGENTS.md) を参照（Go handler 規約・コマンド）

## チーム方針

- 初学者メンバーが多い。**宣言的・読みやすい構造**を優先し、ref／フラグが増える素朴実装は避ける
- コメントは最低限。well-named な関数・型で意図を伝える
- コントラクト層（API/型）は厳密に、実装層は再生成可能に保つ
