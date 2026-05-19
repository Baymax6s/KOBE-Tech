# Codex CLI 設定手順

Codex CLI は OpenAI 製で、設定はユーザー単位（`~/.codex/config.toml`）です。リポジトリ単位の設定機構が無いため、各自一度だけ以下を実施してください。

## 1. MCP（context7）を登録

`~/.codex/config.toml` に以下を追記します（既存設定があれば末尾に追加）。

```toml
[mcp_servers.context7]
command = "npx"
args = ["-y", "@upstash/context7-mcp"]
```

stdio 経由で接続します。`npx` が必要なので Node.js（v20 以上）が入っていること。

> HTTP 接続（`url = "https://mcp.context7.com/mcp"`）も Codex CLI 0.29 以降の `experimental_use_rmcp_client = true` で利用できますが、安定運用は stdio 版を推奨します。

## 2. AGENTS.md は自動で読まれる

Codex CLI は CWD から親方向のすべての `AGENTS.md` を自動で読み込みます。このリポジトリのルート [`AGENTS.md`](../AGENTS.md) と [`app/AGENTS.md`](../app/AGENTS.md) は何もしなくても適用されます。

## 3. 動作確認

```sh
codex
```

起動後、`/mcp` 系のコマンドで `context7` がリストアップされるか確認してください（Codex のバージョンによりコマンド名は異なります）。

## このリポジトリで対応できないこと

- **Hooks（編集後の自動 format / generated 編集のブロック）**: Codex CLI には Claude Code の hooks や OpenCode の plugins に相当する機構がありません。`AGENTS.md` でルールとして書くにとどめ、フォーマットは `cd app && npm run fix` を手動で実行してください。`app/src/api/generated/**` を触らないルールは AGENTS.md に明記済みです。
