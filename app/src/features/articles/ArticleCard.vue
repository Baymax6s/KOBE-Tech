<script setup lang="ts">
import { useDateFormat } from '@vueuse/core'
import type { ServerArticleJSONResponse } from '@/api/generated/apiSchema'

const props = defineProps<{ article: ServerArticleJSONResponse }>()

const formattedDate = useDateFormat(
  () => props.article.created_at,
  'YYYY/MM/DD',
)
</script>

<template>
  <v-card :to="`/articles/${article.id}`" class="p-4">
    <v-card-title class="text-base font-medium">
      {{ article.title }}
    </v-card-title>
    
    <v-card-subtitle class="text-sm text-gray-500 d-flex align-center justify-space-between">
      <span>{{ formattedDate }}</span>

      <div class="d-flex align-center">
        <v-icon 
          icon="mdi-heart-outline" 
          size="small" 
          color="red-lighten-2" 
          class="me-1" 
        />
        <span>{{ article.likes_count ?? 0 }}</span>
      </div>
    </v-card-subtitle>
  </v-card>
</template>