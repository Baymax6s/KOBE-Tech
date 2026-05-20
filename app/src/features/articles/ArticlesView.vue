<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArticleCard from './ArticleCard.vue'
import type { ServerArticleJSONResponse } from '@/api/generated/apiSchema'
import { api } from '@/api/client'
import { useArticleNotificationStore } from '@/stores/articleNotification'

const route = useRoute()
const router = useRouter()

const articles = ref<ServerArticleJSONResponse[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const showCreatedAlert = ref(false)
const notificationStore = useArticleNotificationStore()

// URL の ?tag=... を唯一の真実として selectedTags を扱う。
// リロード / ブックマーク / 共有リンクで絞り込み状態を再現できるようにするため、
// ローカル ref ではなく route.query と双方向にバインドする。
const selectedTags = computed<string[]>({
  get() {
    const q = route.query.tag
    if (Array.isArray(q)) {
      return q.filter((t): t is string => typeof t === 'string' && t.length > 0)
    }
    if (typeof q === 'string' && q.length > 0) return [q]
    return []
  },
  set(next) {
    // フィルタ操作 1 回ごとに履歴を積むと「戻る」が壊れるので push ではなく replace。
    void router.replace({
      query: { ...route.query, tag: next.length ? next : undefined },
    })
  },
})

const toggleTag = (tagName: string) => {
  const current = selectedTags.value
  const index = current.indexOf(tagName)
  selectedTags.value =
    index === -1 ? [...current, tagName] : current.filter((_, i) => i !== index)
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

const fetchArticles = async () => {
  loading.value = true
  error.value = null
  try {
    const response = await api.api.articlesList({
      tag: selectedTags.value.length ? selectedTags.value : undefined,
    })
    articles.value = response.data.articles ?? []
  } catch {
    error.value = '記事の取得に失敗しました。時間をおいて再度お試しください。'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (notificationStore.consumeCreated()) {
    showCreatedAlert.value = true
  }
})

// URL の tag クエリが変わるたびにサーバ側で再フィルタした結果を取得する。
watch(selectedTags, () => void fetchArticles(), { immediate: true })
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
          :items="articleTagNames"
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
            v-for="article in articles"
            :key="article.id"
            :article="article"
            :selected-tags="selectedTags"
            @select-tag="toggleTag"
          />

          <v-alert v-if="articles.length === 0" type="info" variant="tonal">
            該当する記事はありません
          </v-alert>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>
