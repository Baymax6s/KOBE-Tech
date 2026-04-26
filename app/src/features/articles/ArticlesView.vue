<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ArticleCard from './ArticleCard.vue'
import type { Article } from './types'
import { api } from '@/api/client'
import { useArticleNotificationStore } from '@/stores/articleNotification'

const articles = ref<Article[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const showCreatedAlert = ref(false)
const notificationStore = useArticleNotificationStore()

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

        <div v-if="loading" class="d-flex justify-center py-12">
          <v-progress-circular indeterminate color="primary" />
        </div>

        <v-alert v-else-if="error" type="error" class="mb-4">
          {{ error }}
        </v-alert>

        <div v-else class="d-flex flex-column ga-4">
          <ArticleCard
            v-for="article in articles"
            :key="article.id"
            :article="article"
          />
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>
