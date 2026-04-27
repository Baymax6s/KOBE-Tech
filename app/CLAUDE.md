# CLAUDE.md

vuetifyのUIコンポーネントを使用してください。
使用するときはできるだけ素（そのまま）のvuetifyを使用するように心がけてください。
理由としては初学者が多いため、cssなどで悩むのを避けたいからです。

## レイアウト・スタイリングのルール

原則として、Vuetifyのコンポーネントとユーティリティクラスだけで完結させてください。

- レイアウト: `v-container` / `v-row` / `v-col` の12グリッドシステムを使う
- 余白・配置: Vuetify組み込みのユーティリティクラス（`d-flex` / `pa-4` / `ga-4` / `justify-center` / `align-center` など）を使う
- 配色: Vuetifyのカラークラス（`bg-grey-lighten-4` / `text-primary` など）を使う
- 縦の領域確保: `v-main` の中では `min-height: 100vh`（Tailwind の `min-h-screen` 含む）を使わない。`v-container` の `fill-height` を使う。`min-h-screen` はヘッダーぶんを二重計上してスクロールバーが出る原因になるため

Tailwind CSSは「Vuetifyのユーティリティでは表現できないどうしても必要なケース」の **逃げ道** という位置づけです。原則使わず、使う場合はレビュー時に「なぜVuetifyで書けないか」を説明できる状態にしてください。

理由: 2系統のCSSフレームワークが混在すると「どちらで書くのが正解か」で迷う温床になります。基本はVuetifyだけ、と決めておくほうが学習コストが低く、コードの一貫性も保てます。

### Tailwind → Vuetify 対応表（よく使うもの）

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

なお `mb-4` / `pa-6` / `py-8` などは**Vuetifyにも同名のクラスがある**ので、見た目はTailwindと同じですがVuetifyのユーティリティとして利用できます（書き換え不要）。

## フォーム送信処理のルール

非同期な送信処理（ログイン・登録・更新など）を実装する際は、必ず以下を両方実装してください。

1. `v-btn` の `:disabled` / `:loading` で UI 上の多重クリックを抑止する
2. 送信ハンドラ関数の冒頭で進行中フラグを確認し、早期 return でガードする

```ts
const submitting = ref(false)

const onSubmit = async () => {
  if (submitting.value) return // 関数レベルのガード（必須）
  submitting.value = true
  try {
    // 送信処理
  } finally {
    submitting.value = false
  }
}
```

理由: `:disabled` だけでは、Enter キー連打やプログラムからの直接呼び出しで二重送信される余地が残るため。UI と関数の両層で防ぐ。
