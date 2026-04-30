<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useDateFormat } from '@vueuse/core'
import { api } from '@/api/client'
import type { ServerGetArticleJSONResponse } from '@/api/generated/apiSchema'

defineOptions({
  name: 'ArticleDetailView',
})

const props = defineProps<{ articleId: number }>()

const article = ref<ServerGetArticleJSONResponse | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const formattedDate = useDateFormat(
  () => article.value?.created_at,
  'YYYY/MM/DD',
)

onMounted(async () => {
  loading.value = true
  try {
    const response = await api.api.articlesDetail(props.articleId)
    article.value = response.data
  } catch {
    error.value = '記事の取得に失敗しました'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <v-container class="py-8">
    <v-row justify="center">
      <v-col cols="12" md="8" lg="6">
        <div v-if="loading" class="d-flex justify-center py-12">
          <v-progress-circular indeterminate color="primary" />
        </div>

        <v-alert v-else-if="error" type="error" class="mb-4">
          {{ error }}
        </v-alert>

        <v-card v-else-if="article">
          <v-card-title>{{ article.title }}</v-card-title>
          <v-card-subtitle>Author: {{ article.author?.name }}</v-card-subtitle>
          <v-card-subtitle>{{ formattedDate }}に公開</v-card-subtitle>
          <v-card-text style="white-space: pre-wrap">{{
            article.content
          }}</v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
