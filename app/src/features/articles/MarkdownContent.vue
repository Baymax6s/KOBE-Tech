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
  <MdPreview
    :id="previewId"
    :model-value="source"
    class="markdown-content"
    language="en-US"
    theme="light"
    preview-theme="github"
    code-theme="github"
    :show-code-row-number="false"
    no-katex
    no-mermaid
  />
</template>

<style scoped>
/* リンク色を Vuetify テーマの primary に合わせる */
.markdown-content :deep(a) {
  color: rgb(var(--v-theme-primary));
}
</style>
