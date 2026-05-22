/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL?: string
  readonly VITE_USE_MSW?: string
  // Markdown 記法チートシート記事の ID。未設定なら NewArticleView は
  // 外部遷移ではなくダイアログでチートシートを表示するフォールバックに切替える。
  readonly VITE_HELP_ARTICLE_MARKDOWN_ID?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
