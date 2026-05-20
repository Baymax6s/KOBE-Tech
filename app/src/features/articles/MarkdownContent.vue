<script setup lang="ts">
import { shallowRef, watch } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js/lib/core'
import type { LanguageFn } from 'highlight.js'

defineOptions({ name: 'MarkdownContent' })

const props = defineProps<{ source: string }>()

// サポート対象言語の動的 import を Record で静的に列挙する。
// 動的 import の引数が完全リテラルだと Vite が per-language で chunk を切り、
// 実際に登場した言語の chunk だけが実行時に fetch される。
// 一方で `import(\`.../${lang}.js\`)` のような可変パスは node_modules 配下では
// Vite の dynamic-import-vars プラグインが解決できず、production で 404 になる。
// 追加したい言語があればここにエントリを増やすだけで済む。
type LanguageImporter = () => Promise<{ default: LanguageFn }>

const LANGUAGE_IMPORTERS: Record<string, LanguageImporter> = {
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
const LANGUAGE_ALIASES: Record<string, string> = {
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

const resolveLanguage = (lang: string): string =>
  LANGUAGE_ALIASES[lang.toLowerCase()] ?? lang.toLowerCase()

const loadedLanguages = new Set<string>()
const loadingPromises = new Map<string, Promise<void>>()

const loadLanguage = (lang: string): Promise<void> => {
  const canonical = resolveLanguage(lang)
  if (loadedLanguages.has(canonical) || hljs.getLanguage(canonical)) {
    loadedLanguages.add(canonical)
    return Promise.resolve()
  }
  const importer = LANGUAGE_IMPORTERS[canonical]
  if (!importer) {
    // 未サポート言語はハイライト無しで表示する。
    loadedLanguages.add(canonical)
    return Promise.resolve()
  }
  const existing = loadingPromises.get(canonical)
  if (existing) return existing

  const promise = importer()
    .then((mod) => {
      hljs.registerLanguage(canonical, mod.default)
      loadedLanguages.add(canonical)
    })
    .catch(() => {
      loadedLanguages.add(canonical)
    })
    .finally(() => {
      loadingPromises.delete(canonical)
    })

  loadingPromises.set(canonical, promise)
  return promise
}

const md: MarkdownIt = new MarkdownIt({
  html: false,
  linkify: true,
  breaks: false,
  highlight: (str, lang) => {
    const canonical = resolveLanguage(lang)
    if (lang && hljs.getLanguage(canonical)) {
      try {
        const { value } = hljs.highlight(str, {
          language: canonical,
          ignoreIllegals: true,
        })
        return `<pre class="hljs"><code class="hljs language-${md.utils.escapeHtml(canonical)}">${value}</code></pre>`
      } catch {
        // パース失敗時はエスケープしたプレーン表示にフォールバックする。
      }
    }
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
  },
})

const FENCE_LANG_RE = /^[ \t]*```([\w+#-]+)/gm

const extractFenceLanguages = (src: string): string[] => {
  const langs = new Set<string>()
  for (const match of src.matchAll(FENCE_LANG_RE)) {
    if (match[1]) langs.add(match[1])
  }
  return [...langs]
}

const html = shallowRef('')

const render = (src: string) => {
  html.value = md.render(src)
}

watch(
  () => props.source,
  async (src) => {
    const source = src ?? ''
    render(source)

    const langs = extractFenceLanguages(source)
    if (!langs.length) return

    const pending = langs.filter(
      (l) => !hljs.getLanguage(l) && !loadedLanguages.has(resolveLanguage(l)),
    )
    if (!pending.length) return

    await Promise.all(pending.map(loadLanguage))

    // 後から register された言語を反映するため再レンダリングする。
    // props.source が再度切り替わっていた場合は古い結果で上書きしない。
    if (props.source === src) render(source)
  },
  { immediate: true },
)
</script>

<template>
  <div class="markdown-content" v-html="html" />
</template>

<style scoped>
.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3),
.markdown-content :deep(h4),
.markdown-content :deep(h5),
.markdown-content :deep(h6) {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
  font-weight: 600;
  line-height: 1.3;
}

.markdown-content :deep(h1) {
  font-size: 1.75em;
  border-bottom: 1px solid rgba(0, 0, 0, 0.12);
  padding-bottom: 0.3em;
}

.markdown-content :deep(h2) {
  font-size: 1.5em;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  padding-bottom: 0.2em;
}

.markdown-content :deep(h3) {
  font-size: 1.25em;
}

.markdown-content :deep(p) {
  margin: 0 0 1em;
  line-height: 1.7;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  padding-left: 1.5em;
  margin: 0 0 1em;
}

.markdown-content :deep(li) {
  margin-bottom: 0.25em;
}

.markdown-content :deep(blockquote) {
  margin: 0 0 1em;
  padding: 0 1em;
  border-left: 4px solid rgba(0, 0, 0, 0.12);
  color: rgba(0, 0, 0, 0.6);
}

.markdown-content :deep(a) {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
}

.markdown-content :deep(a:hover) {
  text-decoration: underline;
}

.markdown-content :deep(hr) {
  border: 0;
  border-top: 1px solid rgba(0, 0, 0, 0.12);
  margin: 1.5em 0;
}

.markdown-content :deep(table) {
  border-collapse: collapse;
  margin: 0 0 1em;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid rgba(0, 0, 0, 0.12);
  padding: 0.4em 0.8em;
}

.markdown-content :deep(code) {
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono',
    'Courier New', monospace;
  font-size: 0.9em;
}

.markdown-content :deep(:not(pre) > code) {
  background-color: rgba(175, 184, 193, 0.2);
  padding: 0.15em 0.4em;
  border-radius: 4px;
}

.markdown-content :deep(pre) {
  margin: 0 0 1em;
  padding: 12px 16px;
  border-radius: 6px;
  overflow-x: auto;
  background-color: #f6f8fa;
}

.markdown-content :deep(pre code) {
  padding: 0;
  background: transparent;
}
</style>
