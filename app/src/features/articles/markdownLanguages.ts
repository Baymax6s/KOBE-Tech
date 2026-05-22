import type { LanguageFn } from 'highlight.js'

// サポート対象言語の動的 import を Record で静的に列挙する。
// 動的 import の引数が完全リテラルだと Vite が per-language で chunk を切り、
// 実際に登場した言語の chunk だけが実行時に fetch される。
// 一方で `import(`.../${lang}.js`)` のような可変パスは node_modules 配下では
// Vite の dynamic-import-vars プラグインが解決できず、production で 404 になる。
// 追加したい言語があればここにエントリを増やすだけで済む。
type LanguageImporter = () => Promise<{ default: LanguageFn }>

export const LANGUAGE_IMPORTERS: Record<string, LanguageImporter> = {
  bash: () => import('highlight.js/lib/languages/bash'),
  c: () => import('highlight.js/lib/languages/c'),
  clojure: () => import('highlight.js/lib/languages/clojure'),
  cmake: () => import('highlight.js/lib/languages/cmake'),
  cpp: () => import('highlight.js/lib/languages/cpp'),
  csharp: () => import('highlight.js/lib/languages/csharp'),
  css: () => import('highlight.js/lib/languages/css'),
  dart: () => import('highlight.js/lib/languages/dart'),
  diff: () => import('highlight.js/lib/languages/diff'),
  dockerfile: () => import('highlight.js/lib/languages/dockerfile'),
  elixir: () => import('highlight.js/lib/languages/elixir'),
  erlang: () => import('highlight.js/lib/languages/erlang'),
  fsharp: () => import('highlight.js/lib/languages/fsharp'),
  go: () => import('highlight.js/lib/languages/go'),
  graphql: () => import('highlight.js/lib/languages/graphql'),
  groovy: () => import('highlight.js/lib/languages/groovy'),
  haskell: () => import('highlight.js/lib/languages/haskell'),
  ini: () => import('highlight.js/lib/languages/ini'),
  java: () => import('highlight.js/lib/languages/java'),
  javascript: () => import('highlight.js/lib/languages/javascript'),
  json: () => import('highlight.js/lib/languages/json'),
  kotlin: () => import('highlight.js/lib/languages/kotlin'),
  latex: () => import('highlight.js/lib/languages/latex'),
  less: () => import('highlight.js/lib/languages/less'),
  lua: () => import('highlight.js/lib/languages/lua'),
  makefile: () => import('highlight.js/lib/languages/makefile'),
  markdown: () => import('highlight.js/lib/languages/markdown'),
  nginx: () => import('highlight.js/lib/languages/nginx'),
  objectivec: () => import('highlight.js/lib/languages/objectivec'),
  ocaml: () => import('highlight.js/lib/languages/ocaml'),
  perl: () => import('highlight.js/lib/languages/perl'),
  php: () => import('highlight.js/lib/languages/php'),
  plaintext: () => import('highlight.js/lib/languages/plaintext'),
  powershell: () => import('highlight.js/lib/languages/powershell'),
  properties: () => import('highlight.js/lib/languages/properties'),
  protobuf: () => import('highlight.js/lib/languages/protobuf'),
  python: () => import('highlight.js/lib/languages/python'),
  r: () => import('highlight.js/lib/languages/r'),
  ruby: () => import('highlight.js/lib/languages/ruby'),
  rust: () => import('highlight.js/lib/languages/rust'),
  scala: () => import('highlight.js/lib/languages/scala'),
  scss: () => import('highlight.js/lib/languages/scss'),
  shell: () => import('highlight.js/lib/languages/shell'),
  sql: () => import('highlight.js/lib/languages/sql'),
  swift: () => import('highlight.js/lib/languages/swift'),
  typescript: () => import('highlight.js/lib/languages/typescript'),
  vbnet: () => import('highlight.js/lib/languages/vbnet'),
  xml: () => import('highlight.js/lib/languages/xml'),
  yaml: () => import('highlight.js/lib/languages/yaml'),
}

// フェンスのよく使われる短縮表記を canonical な言語名へ寄せる。
// highlight.js は登録済み言語の alias を自動解決するので、canonical 名で
// register すれば `ts` `js` 等もそのまま getLanguage で引ける。
export const LANGUAGE_ALIASES: Record<string, string> = {
  js: 'javascript',
  jsx: 'javascript',
  ts: 'typescript',
  tsx: 'typescript',
  py: 'python',
  sh: 'bash',
  zsh: 'bash',
  yml: 'yaml',
  md: 'markdown',
  rb: 'ruby',
  kt: 'kotlin',
  cs: 'csharp',
  golang: 'go',
  html: 'xml',
  vue: 'xml',
  text: 'plaintext',
  txt: 'plaintext',
  toml: 'ini',
}

// アプリで利用可能な言語を「canonical 名」と「その alias 群」の組で公開する。
// LANGUAGE_IMPORTERS / LANGUAGE_ALIASES を真のソースとして派生させるので、
// 言語を追加・削除すれば SupportedLanguagesDialog の表示も自動的に追随する。
export type SupportedLanguage = {
  name: string
  aliases: string[]
}

const aliasesByCanonical = Object.entries(LANGUAGE_ALIASES).reduce<
  Record<string, string[]>
>((acc, [alias, canonical]) => {
  if (!acc[canonical]) acc[canonical] = []
  acc[canonical].push(alias)
  return acc
}, {})

export const SUPPORTED_LANGUAGES: SupportedLanguage[] = Object.keys(
  LANGUAGE_IMPORTERS,
)
  .sort()
  .map((name) => ({
    name,
    aliases: (aliasesByCanonical[name] ?? []).sort(),
  }))
