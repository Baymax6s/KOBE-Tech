<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import ArticleCard from './ArticleCard.vue'
import type { ServerArticleJSONResponse } from '@/api/generated/apiSchema'
import { api } from '@/api/client'
import { useArticleNotificationStore } from '@/stores/articleNotification'

type ArticleTag = NonNullable<ServerArticleJSONResponse['tags']>[number]

const articles = ref<ServerArticleJSONResponse[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const showCreatedAlert = ref(false)
const notificationStore = useArticleNotificationStore()

const selectedTags = ref<string[]>([])
const tagCandidates = ref<string[]>([])

const toggleTag = (tagName: string) => {
  const index = selectedTags.value.indexOf(tagName)

  if (index === -1) {
    selectedTags.value.push(tagName)
  } else {
    selectedTags.value.splice(index, 1)
  }
}

const clearTag = () => {
  selectedTags.value = []
}

const articleTagNames = computed(() => {
  const tagNames = new Set<string>()

  for (const article of articles.value) {
    for (const tag of article.tags ?? []) {
      tagNames.add(tag.name)
    }
  }

  return [...tagNames].sort((a, b) => a.localeCompare(b))
})

const filterTagItems = computed(() => {
  if (tagCandidates.value.length > 0) {
    return tagCandidates.value
  }

  return articleTagNames.value
})

const filteredArticles = computed(() => {
  if (selectedTags.value.length === 0) return articles.value

  return articles.value.filter((article) =>
    selectedTags.value.every((tagName) =>
      article.tags?.some((tag: ArticleTag) => tag.name === tagName),
    ),
  )
})

onMounted(async () => {
  if (notificationStore.consumeCreated()) {
    showCreatedAlert.value = true
  }

  loading.value = true
  try {
    const response = await api.api.articlesList()
    articles.value = response.data.articles ?? []
  } catch {
    error.value = '記事の取得に失敗しました。時間をおいて再度お試しください。'
  } finally {
    loading.value = false
  }

  tagCandidates.value = []
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

        <v-select
          v-model="selectedTags"
          :items="filterTagItems"
          label="タグで絞り込み"
          prepend-inner-icon="mdi-tag-outline"
          variant="outlined"
          density="comfortable"
          multiple
          chips
          closable-chips
          clearable
          hide-details
          class="mb-4"
        />

        <v-fade-transition>
          <div
            v-if="selectedTags.length"
            class="mb-4 d-flex align-center flex-wrap ga-2"
          >
            <span class="text-subtitle-2 me-2">絞り込み中:</span>

            <v-chip
              v-for="tag in selectedTags"
              :key="tag"
              closable
              color="primary"
              label
              @click:close="toggleTag(tag)"
            >
              {{ tag }}
            </v-chip>

            <v-btn
              aria-label="絞り込みをすべて解除"
              icon="mdi-close-circle"
              size="small"
              color="grey"
              variant="text"
              class="ms-1"
              @click="clearTag"
            />
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
            :selected-tags="selectedTags"
            @select-tag="toggleTag"
          />

          <v-alert
            v-if="filteredArticles.length === 0"
            type="info"
            variant="tonal"
          >
            該当する記事はありません
          </v-alert>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>
