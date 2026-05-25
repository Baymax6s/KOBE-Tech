<script setup lang="ts">
import { computed } from 'vue'
import { useDateFormat } from '@vueuse/core'
import type { ServerArticleJSONResponse } from '@/api/generated/apiSchema'

const props = defineProps<{
  article: ServerArticleJSONResponse
  selectedTags: string[]
}>()

const emit = defineEmits<{
  (e: 'select-tag', tagName: string): void
}>()

const isSelected = (tagName: string) => props.selectedTags.includes(tagName)

// 質問ステータスごとのチップ表示。'none'（質問が無い通常記事）はバッジを出さない。
type QuestionStatus = ServerArticleJSONResponse['question_status']
const questionStatusBadges: Record<
  Exclude<QuestionStatus, 'none'>,
  { label: string; color: string; icon: string }
> = {
  unanswered: { label: '未回答', color: 'orange', icon: 'mdi-help-circle' },
  answered: { label: '回答あり', color: 'blue', icon: 'mdi-comment-text' },
  solved: { label: '解決済', color: 'green', icon: 'mdi-check-circle' },
}

const statusBadge = computed(() =>
  props.article.question_status === 'none'
    ? null
    : questionStatusBadges[props.article.question_status],
)

const formattedDate = useDateFormat(
  () => props.article.created_at ?? '',
  'YYYY/MM/DD',
)
</script>

<template>
  <v-card :to="`/articles/${article.id}`" class="pa-4">
    <v-card-title
      class="d-flex align-center ga-2 text-body-1 font-weight-medium"
    >
      <span>{{ article.title }}</span>
      <v-chip
        v-if="statusBadge"
        size="small"
        variant="flat"
        :color="statusBadge.color"
        :prepend-icon="statusBadge.icon"
      >
        {{ statusBadge.label }}
      </v-chip>
    </v-card-title>

    <v-card-text v-if="article.tags?.length" class="py-0 px-4">
      <div class="d-flex ga-1 flex-wrap">
        <v-chip
          v-for="tag in article.tags"
          :key="tag.id"
          size="small"
          :variant="isSelected(tag.name) ? 'flat' : 'outlined'"
          :color="isSelected(tag.name) ? 'primary' : undefined"
          @click.stop.prevent="emit('select-tag', tag.name)"
        >
          {{ tag.name }}
        </v-chip>
      </div>
    </v-card-text>

    <v-card-subtitle class="d-flex align-center justify-space-between mt-2">
      <span class="text-caption text-medium-emphasis">{{ formattedDate }}</span>

      <div class="d-flex align-center">
        <v-icon
          :icon="article.liked_by_me ? 'mdi-heart' : 'mdi-heart-outline'"
          size="small"
          color="red"
          class="me-1"
        />
        <span class="text-caption">{{ article.likes_count ?? 0 }}</span>
      </div>
    </v-card-subtitle>
  </v-card>
</template>
