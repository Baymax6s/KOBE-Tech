<script setup lang="ts">
import { useId } from 'vue'
import { MdPreview, type CustomIcon } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'

defineOptions({ name: 'MarkdownContent' })

defineProps<{ source: string }>()

// 同一ページに複数マウントしても DOM (アンカー id 等) が衝突しないよう一意な id を振る
const previewId = useId()

// コピーボタンをテキスト「Copy」からアイコンに置き換える。customIcon.copy を渡すと
// md-editor-v3 が data-is-icon モードに切り替わり、クリック時の「Copied!」は
// ツールチップ (data-tips / language の successTips) として表示される。
// サイズは md-editor-icon クラスに対する既定 CSS (15x15) に委ねる。
const customIcon: CustomIcon = {
  copy: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="md-editor-icon"><rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/></svg>',
}
</script>

<template>
  <div class="markdown-content">
    <MdPreview
      :id="previewId"
      :model-value="source"
      language="en-US"
      theme="light"
      preview-theme="github"
      code-theme="github"
      :show-code-row-number="false"
      :code-foldable="false"
      :custom-icon="customIcon"
      no-katex
      no-mermaid
    />
  </div>
</template>

<style scoped>
/* リンク色を Vuetify テーマの primary に合わせる */
.markdown-content :deep(a) {
  color: rgb(var(--v-theme-primary));
}

/* コードブロック左上の Mac 風 3 色ドットを非表示にする。
   md-editor-v3 側が詳細度の高い (クラス 5 個) セレクタで表示しているため !important で上書きする */
.markdown-content :deep(.md-editor-code-flag span) {
  display: none !important;
}
</style>
