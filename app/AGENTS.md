# app/ ルール（フロントエンド）

`app/` 配下を編集するすべての AI コーディングエージェント（Claude Code / Codex / OpenCode 等）はここを参照してください。初学者メンバーが多いため、CSS で悩む選択肢を減らすことを優先します。

## 大原則: Vuetify を素で使う

vuetify の UI コンポーネントを使ってください。できるだけ素（そのまま）の vuetify を使うように心がけてください。

## レイアウト・スタイリング

Vuetify のコンポーネントとユーティリティクラスだけで完結させてください。

- レイアウト: `v-container` / `v-row` / `v-col` の 12 グリッド
- 余白・配置: `d-flex` / `pa-4` / `ga-4` / `justify-center` / `align-center` などの Vuetify ユーティリティ
- 配色: `bg-grey-lighten-4` / `text-primary` などの Vuetify カラークラス
- 縦の領域確保: `v-main` の中では `min-height: 100vh`（Tailwind の `min-h-screen` 含む）を使わない。`v-container` の `fill-height` を使う

理由（`min-h-screen` 禁止）: ヘッダーぶんを二重計上してスクロールバーが出る原因になるため。

## Tailwind は逃げ道

Tailwind CSS は「Vuetify のユーティリティでは表現できないどうしても必要なケース」のみ。使う場合はレビュー時に「なぜ Vuetify で書けないか」を説明できる状態にしてください。

理由: 2 系統の CSS フレームワークが混在すると「どちらで書くのが正解か」で迷う温床になるため。

### Tailwind → Vuetify 対応表

| Tailwind                   | Vuetify                             |
| -------------------------- | ----------------------------------- |
| `flex`                     | `d-flex`                            |
| `flex-col`                 | `flex-column`                       |
| `gap-4`                    | `ga-4`                              |
| `justify-center`           | `justify-center`（同名）            |
| `items-center`             | `align-center`                      |
| `min-h-screen`（v-main内） | `v-container` + `fill-height`       |
| `max-w-md` / `max-w-2xl`   | `v-col` の `cols/sm/md/lg` で幅指定 |
| `mx-auto`                  | `v-row` + `v-col` で自然に中央寄せ  |
| `bg-gray-100`              | `bg-grey-lighten-4`                 |

`mb-4` / `pa-6` / `py-8` などは Vuetify にも同名クラスがあるので、Vuetify のユーティリティとしてそのまま使えます。

## フォーム送信処理

非同期な送信処理（ログイン・登録・更新など）は、必ず以下を両方実装してください。

1. `v-btn` の `:disabled` / `:loading` で UI 上の多重クリックを抑止する
2. 送信ハンドラ関数の冒頭で進行中フラグを確認し、早期 return でガードする

```ts
const submitting = ref(false)

const onSubmit = async () => {
  if (submitting.value) return
  submitting.value = true
  try {
    // 送信処理
  } finally {
    submitting.value = false
  }
}
```

理由: `:disabled` だけでは Enter キー連打やプログラムからの直接呼び出しで二重送信される余地が残るため。UI と関数の両層で防ぐ。

## 触ってはいけないファイル

- `app/src/api/generated/**` — `swagger-typescript-api` が `api/swagger/openapi.yml` から生成しているコントラクト層。手で編集すると次の再生成で消えます。変更が必要なときは Go 側のハンドラ／DTO を直し、`make swagger`（api/）→ `npm run generate:api`（app/）で再生成してください。
