export const articles = [
  {
    id: 1,
    title: 'MSWを導入してみた',
    content: `# MSWを導入してみた

フロントだけで API のモックを切り替えたかったので、Mock Service Worker (MSW) を入れてみました。

## セットアップ

\`\`\`bash
npm install --save-dev msw
npx msw init public/ --save
\`\`\`

## ハンドラの定義

\`\`\`ts
import { http, HttpResponse } from 'msw'

export const handlers = [
  http.get('/api/articles', () =>
    HttpResponse.json({ articles: [{ id: 1, title: 'hello' }] }),
  ),
]
\`\`\`

## 起動フラグで切り替え

\`\`\`ts
if (import.meta.env.DEV && import.meta.env.VITE_USE_MSW === 'true') {
  const { worker } = await import('./mocks/browser')
  await worker.start({ onUnhandledRequest: 'bypass' })
}
\`\`\`

> 環境変数で切り替えると、本物の API を叩きたい日も同じビルドで対応できます。

## 良かった点

- バックエンドが落ちていてもフロント開発が止まらない
- E2E テストでネットワーク依存を排除できる
- DevTools の Network タブで普通に確認できる
`,
    user_id: 1,
    likes_count: 3,
    tags: [{ id: 1, name: 'dev' }],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: 2,
    title: 'Vue3 + TS構成まとめ',
    content: `# Vue3 + TS構成まとめ

Composition API を使った Vue 3 + TypeScript の構成をまとめます。

## defineProps / defineEmits の型

\`\`\`ts
const props = defineProps<{
  articleId: number
  initialLiked?: boolean
}>()

const emit = defineEmits<{
  (e: 'like', id: number): void
}>()
\`\`\`

## ref と computed

\`\`\`ts
import { ref, computed } from 'vue'

const count = ref(0)
const isEven = computed(() => count.value % 2 === 0)
\`\`\`

## Pinia ストア

\`\`\`ts
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useCounter = defineStore('counter', () => {
  const count = ref(0)
  const increment = () => count.value++
  return { count, increment }
})
\`\`\`

## ポイント

- **型は \`defineProps\` の型引数で渡す**（runtime declaration より型が効く）
- **\`<script setup>\`** で書く（テンプレートに自動公開される）
- Pinia は **setup store** が Composition API と一貫していておすすめ
`,
    user_id: 2,
    likes_count: 7,
    tags: [{ id: 5, name: 'typescript' }],
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
]
