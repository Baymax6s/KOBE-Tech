import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js/lib/core'
import { LANGUAGE_ALIASES, LANGUAGE_IMPORTERS } from './markdownLanguages'

const resolveLanguage = (lang: string): string =>
  LANGUAGE_ALIASES[lang.toLowerCase()] ?? lang.toLowerCase()

// アプリ寿命のキャッシュ。MarkdownContent が unmount/remount を繰り返しても
// 登録済み言語の二重ロードを防ぐ。
const loadedLanguages = new Set<string>()
const loadingPromises = new Map<string, Promise<void>>()

export const isLanguageReady = (lang: string): boolean => {
  const canonical = resolveLanguage(lang)
  return (
    loadedLanguages.has(canonical) || hljs.getLanguage(canonical) !== undefined
  )
}

export const loadLanguage = (lang: string): Promise<void> => {
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

const FENCE_LANG_RE = /^[ \t]*```([\w+#-]+)/gm

export const extractFenceLanguages = (src: string): string[] => {
  const langs = new Set<string>()
  for (const match of src.matchAll(FENCE_LANG_RE)) {
    if (match[1]) langs.add(match[1])
  }
  return [...langs]
}

export const md: MarkdownIt = new MarkdownIt({
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
