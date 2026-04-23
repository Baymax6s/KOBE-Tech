<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ArticleCard from './ArticleCard.vue'

// TODO: スキーマ生成後は api/generated/apiSchema.ts の型に差し替える
interface Article {
  id: number
  title: string
  created_at: string
}

const articles = ref<Article[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

// TODO: スキーマ生成後は生成クライアントの呼び出しに差し替える
const mockArticles: Article[] = [
  { id: 1, title: '神戸大学でのハッカソン体験記', created_at: '2026-04-01T09:00:00Z' },
  { id: 2, title: 'Goで作るREST API入門', created_at: '2026-04-10T12:30:00Z' },
  { id: 3, title: 'Vue 3 + Vuetifyで学ぶフロントエンド開発', created_at: '2026-04-20T15:00:00Z' },
]

onMounted(async () => {
  loading.value = true
  try {
    await new Promise((resolve) => setTimeout(resolve, 500))
    articles.value = mockArticles
  } catch {
    error.value = '記事の取得に失敗しました'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="max-w-2xl mx-auto px-4 py-8">

    <div v-if="loading" class="flex justify-center py-12">
      <v-progress-circular indeterminate color="primary" />
    </div>

    <v-alert v-else-if="error" type="error" class="mb-4">
      {{ error }}
    </v-alert>

    <div v-else class="flex flex-col gap-4">
      <ArticleCard
        v-for="article in articles"
        :key="article.id"
        :article="article"
      />
    </div>
  </div>
</template>
