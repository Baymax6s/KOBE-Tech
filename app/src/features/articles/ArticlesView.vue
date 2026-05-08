<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import ArticleCard from './ArticleCard.vue'
import type { ServerArticleJSONResponse } from '@/api/generated/apiSchema'
import { api } from '@/api/client'
import { useArticleNotificationStore } from '@/stores/articleNotification'

const articles = ref<ServerArticleJSONResponse[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const showCreatedAlert = ref(false)
const notificationStore = useArticleNotificationStore()
const selectedTag = ref<string | null>(null)

const filteredArticles = computed(() => {
  if (!selectedTag.value) return articles.value
  return articles.value.filter((article) =>
    article.tags?.some((tag: { name: string }) => tag.name === selectedTag.value)
  )
})

const clearTag = () => {
  selectedTag.value = null
}

onMounted(async () => {
  if (notificationStore.consumeCreated()) {
    showCreatedAlert.value = true
  }

  loading.value = true
  try {
    const response = await api.api.articlesList()
    articles.value = response.data.articles ?? []
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
        <v-alert
          v-if="showCreatedAlert"
          type="success"
          class="mb-4"
          closable
          @click:close="showCreatedAlert = false"
        >
          記事を投稿しました
        </v-alert>

        <v-fade-transition>
          <div v-if="selectedTag" class="mb-4 d-flex align-center">
            <span class="text-subtitle-2 me-2">絞り込み中:</span>
            <v-chip closable color="primary" label @click:close="clearTag">
              {{ selectedTag }}
            </v-chip>
          </div>
        </v-fade-transition>

        <div v-if="loading" class="d-flex justify-center py-12">
          <v-progress-circular indeterminate color="primary" />
        </div>

        <v-alert v-else-if="error" type="error" class="mb-4">
          {{ error }}
        </v-alert>

        <div v-else class="d-flex flex-column ga-4">
          <ArticleCard
            v-for="article in filteredArticles"
            :key="article.id"
            :article="article"
            @select-tag="selectedTag = $event"
          />

          <v-alert v-if="filteredArticles.length === 0" type="info" variant="tonal">
            該当する記事はありません
          </v-alert>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>