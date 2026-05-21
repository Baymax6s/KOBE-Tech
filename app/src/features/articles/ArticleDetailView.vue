<script setup lang="ts">
import axios from 'axios'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useDateFormat } from '@vueuse/core'
import { api } from '@/api/client'
import type { ServerGetArticleJSONResponse } from '@/api/generated/apiSchema'
import ReplySection from '@/features/replies/ReplySection.vue'
import { useAuthStore } from '@/stores/auth'

defineOptions({ name: 'ArticleDetailView' })

const props = defineProps<{ articleId: number }>()

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const article = ref<ServerGetArticleJSONResponse | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
const isLiked = computed(() => article.value?.liked_by_me ?? false)
const likeSubmitting = ref(false)
const likeError = ref<string | null>(null)

const likeArticle = async () => {
  if (!article.value || likeSubmitting.value) return

  likeError.value = null

  if (!auth.isAuthenticated) {
    await router.push({ path: '/login', query: { redirect: route.fullPath } })
    return
  }

  likeSubmitting.value = true

  try {
    if (isLiked.value) {
      const response = await api.api.articlesLikeDelete(props.articleId)
      article.value.liked_by_me =
        response.data.liked_by_me ?? article.value.liked_by_me
      article.value.likes_count =
        response.data.likes_count ?? article.value.likes_count
    } else {
      const response = await api.api.articlesLikeCreate(props.articleId)
      article.value.liked_by_me =
        response.data.liked_by_me ?? article.value.liked_by_me
      article.value.likes_count =
        response.data.likes_count ?? article.value.likes_count
    }
  } catch (err: unknown) {
    if (axios.isAxiosError(err) && err.response?.status === 409) {
      article.value.liked_by_me = true
      return
    }

    if (axios.isAxiosError(err) && err.response?.status === 404) {
      article.value.liked_by_me = false
      return
    }

    if (axios.isAxiosError(err) && err.response?.status === 401) {
      auth.clearToken()
      await router.push({ path: '/login', query: { redirect: route.fullPath } })
      return
    }

    likeError.value =
      'いいねの処理に失敗しました。時間をおいて再度お試しください。'
  } finally {
    likeSubmitting.value = false
  }
}

const formattedDate = useDateFormat(
  () => article.value?.created_at ?? '',
  'YYYY/MM/DD',
)

watch(
  () => props.articleId,
  async (id) => {
    article.value = null
    error.value = null
    likeError.value = null
    loading.value = true
    likeSubmitting.value = false

    try {
      const response = await api.api.articlesDetail(id)
      article.value = response.data
    } catch {
      error.value = '記事の取得に失敗しました。時間をおいて再度お試しください。'
    } finally {
      loading.value = false
    }
  },
  { immediate: true },
)
</script>

<template>
  <v-sheet color="grey-lighten-4" min-height="100%">
    <v-container class="py-8">
      <v-row justify="center">
        <v-col cols="12" sm="10">
          <div v-if="loading" class="d-flex justify-center py-12">
            <v-progress-circular indeterminate color="primary" />
          </div>

          <v-alert v-else-if="error" type="error">
            {{ error }}
          </v-alert>

          <template v-else-if="article">
            <h1 class="text-h4 font-weight-bold mb-4">
              {{ article.title }}
            </h1>

            <div v-if="article.tags?.length" class="mb-4 d-flex flex-wrap ga-2">
              <v-chip
                v-for="tag in article.tags"
                :key="tag.id"
                size="small"
                color="primary"
                variant="flat"
                label
              >
                <v-icon start icon="mdi-tag-outline" size="small" />
                {{ tag.name }}
              </v-chip>
            </div>

            <div class="text-body-2 text-medium-emphasis mb-6">
              <div>著者 {{ article.author?.name }}</div>
              <div>投稿日 {{ formattedDate }}</div>
            </div>

            <v-alert v-if="likeError" type="error" class="mb-4">
              {{ likeError }}
            </v-alert>

            <v-card flat rounded="lg" class="pa-8">
              <div class="text-body-1" style="white-space: pre-wrap">
                {{ article.content }}
              </div>
            </v-card>

            <div class="d-flex align-center ga-3 my-6">
              <v-btn
                variant="text"
                :color="isLiked ? 'red' : undefined"
                size="large"
                :loading="likeSubmitting"
                @click="likeArticle"
              >
                <v-icon
                  :icon="isLiked ? 'mdi-heart' : 'mdi-heart-outline'"
                  size="large"
                  start
                />
                {{ article.likes_count ?? 0 }}
              </v-btn>
            </div>

            <div class="mt-10">
              <ReplySection :article-id="article.id" />
            </div>
          </template>
        </v-col>
      </v-row>
    </v-container>
  </v-sheet>
</template>
