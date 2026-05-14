<script setup lang="ts">
import { useDateFormat } from '@vueuse/core'
import type { ServerArticleJSONResponse } from '@/api/generated/apiSchema'

const props = defineProps<{
  article: ServerArticleJSONResponse
  selectedTags: string[]
}>()

const emit = defineEmits<{
  (e: 'select-tag', tagName: string): void
}>()



const isSelected = (tagName: string) =>
  props.selectedTags.includes(tagName)



const formattedDate = useDateFormat(
  () => props.article.created_at ?? '',
  'YYYY/MM/DD',
)
</script>

<template>
  <v-card :to="`/articles/${article.id}`" class="p-4">
    <v-card-title class="text-base font-medium">
      {{ article.title }}
    </v-card-title>

  <v-card-text v-if="article.tags?.length" class="py-0 px-4">
    <div class="d-flex ga-1 flex-wrap">
      <v-chip
        v-for="tag in article.tags"
          :key="tag.id"
          size="x-small"
          :variant="isSelected(tag.name) ? 'flat' : 'outlined'"
          :color="isSelected(tag.name) ? 'primary' : undefined"
          @click.stop.prevent="emit('select-tag', tag.name)"
        >
          {{ tag.name }}
        </v-chip>
      </div>
    </v-card-text>


    <v-card-subtitle
      class="text-sm text-gray-500 d-flex align-center justify-space-between mt-2"
    >
      <span>{{ formattedDate }}</span>

      <div class="d-flex align-center">
        <v-icon
          icon="mdi-heart-outline"
          size="small"
          color="red"
          class="me-1"
        />
        <span>{{ article.likes_count ?? 0 }}</span>
      </div>
    </v-card-subtitle>
  </v-card>
</template>
