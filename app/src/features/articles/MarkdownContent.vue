<script setup lang="ts">
import { useId } from 'vue'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'

defineOptions({ name: 'MarkdownContent' })

defineProps<{ source: string }>()

// 同一ページに複数マウントしても DOM (アンカー id 等) が衝突しないよう一意な id を振る
const previewId = useId()
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
