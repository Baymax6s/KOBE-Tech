---
name: vuetify-rules
description: MUST use when editing files under `app/` (Vue/Vuetify frontend). Loads this project's Vuetify-first UI rules from `app/AGENTS.md` — Vuetify components and utility classes only, Tailwind is a last-resort escape hatch, layout via `v-container`/`v-row`/`v-col`, form-submit double-click guards, and the `app/src/api/generated/**` no-touch rule. Load when touching .vue, .ts, .tsx files under `app/src/`, or when implementing forms, layouts, or styling on the frontend.
---

# Vuetify ルール（KOBE-Tech フロントエンド）

このスキルは `app/` 配下を編集する前に発火します。**ルール本文は [`app/AGENTS.md`](../../../app/AGENTS.md) に集約されています**。Read ツールでそのファイルを読み込んでから作業してください。

`app/AGENTS.md` は Claude Code / Codex / OpenCode の共通指示ファイルです。ここを真のソースとし、各エージェント設定からポインタで参照する構成にしています（二重管理を避けるため）。

## ルールの種類（要約）

- Vuetify を素で使う（CSS で悩む選択肢を減らす）
- レイアウトは `v-container` / `v-row` / `v-col`、ユーティリティクラスは Vuetify のものを使う
- Tailwind は逃げ道のみ
- フォーム送信は `:disabled`/`:loading` と関数レベルの早期 return ガードの両方
- `app/src/api/generated/**` は触らない（自動生成）

詳細・対応表・コード例は `app/AGENTS.md` を参照してください。
