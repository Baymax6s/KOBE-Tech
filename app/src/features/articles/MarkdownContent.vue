<script setup lang="ts">
import { shallowRef, watch } from 'vue'
import {
  extractFenceLanguages,
  isLanguageReady,
  loadLanguage,
  md,
} from './markdownIt'

// markdown-it 自身は CSS を持たないため、レンダリング後の HTML
// (見出し / 表 / リスト / コードブロック等) のタイポグラフィは
// github-markdown-css に委ねる。prefers-color-scheme で
// light/dark を自動切替するバリアントを採用。
// コードハイライトの配色は main.ts で highlight.js/styles/github.css を
// 読み込み済み (.hljs クラスに対する着色)。
import 'github-markdown-css/github-markdown.css'

defineOptions({ name: 'MarkdownContent' })

const props = defineProps<{ source: string }>()

const html = shallowRef('')

const render = (src: string) => {
  html.value = md.render(src)
}

watch(
  () => props.source,
  async (src) => {
    const source = src ?? ''
    render(source)

    const pending = extractFenceLanguages(source).filter(
      (l) => !isLanguageReady(l),
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
  <div class="markdown-body markdown-content" v-html="html" />
</template>

<style scoped>
/*
 * markdown-it が生成した HTML を v-html で挿入しているため、個々の要素に
 * ビルド時 class を付与できず Vuetify / Tailwind のユーティリティが当たらない。
 * タイポグラフィ全般は github-markdown-css (.markdown-body) に委譲し、
 * ここでは Vuetify テーマと統一したい箇所のみを :deep() で上書きする。
 */

/* 本文フォントを Vuetify アプリ全体のフォントへ揃える */
.markdown-content {
  font-family: inherit;
}

/* リンク色を Vuetify テーマの primary に合わせる (github-markdown-css の青より優先) */
.markdown-content :deep(a) {
  color: rgb(var(--v-theme-primary));
}
</style>
