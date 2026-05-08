<script setup lang="ts">
import { ref, watch } from 'vue'
import { useDateFormat } from '@vueuse/core'
import { useRouter } from 'vue-router'
import { api } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import type { ServerGetArticleJSONResponse } from '@/api/generated/apiSchema'

defineOptions({
  name: 'ArticleDetailView',
})

const props = defineProps<{ articleId: number }>()

const router = useRouter()
const authStore = useAuthStore()

const article = ref<ServerGetArticleJSONResponse | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const isLiked = ref(false)

/*
  いいね処理
*/
const toggleLike = () => {
  if (!article.value) return

  // 未ログインならログイン画面へ
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }

  // いいね切り替え
  isLiked.value = !isLiked.value

  // いいね数変更
  if (isLiked.value) {
    article.value.likes_count = (article.value.likes_count ?? 0) + 1
  } else {
    article.value.likes_count = Math.max(
      0,
      (article.value.likes_count ?? 1) - 1,
    )
  }
}

/*
  投稿日フォーマット
*/
const formattedDate = useDateFormat(
  () => article.value?.created_at,
  'YYYY/MM/DD',
)

/*
  記事取得
*/
watch(
  () => props.articleId,
  async (id) => {
    article.value = null
    error.value = null
    loading.value = true
    isLiked.value = false

    try {
      const response = await api.api.articlesDetail(id)
      article.value = response.data
    } catch {
      error.value = '記事の取得に失敗しました'
    } finally {
      loading.value = false
    }
  },
  { immediate: true },
)
</script>

<template>
  <v-sheet color="grey-lighten-4" min-height="100%">
    <v-container class="py-12">
      <v-row justify="center">
        <v-col cols="12" md="8" lg="7">
          <!-- ローディング -->
          <div v-if="loading" class="d-flex justify-center py-12">
            <v-progress-circular indeterminate color="primary" />
          </div>

          <!-- エラー -->
          <v-alert v-else-if="error" type="error">
            {{ error }}
          </v-alert>

          <!-- 記事 -->
          <template v-else-if="article">
            <!-- タイトル -->
            <h1 class="text-h4 font-weight-bold mb-4">
              {{ article.title }}
            </h1>

            <!-- 著者・投稿日・いいね -->
            <div
              class="text-body-2 text-medium-emphasis mb-6 d-flex align-center"
            >
              <div>
                <div>著者 {{ article.author?.name }}</div>

                <div>投稿日 {{ formattedDate }}</div>
              </div>

              <v-spacer />

              <!-- いいね -->
              <div class="d-flex align-center">
                <v-btn
                  variant="text"
                  icon
                  color="red-lighten-2"
                  @click="toggleLike"
                >
                  <v-icon :icon="isLiked ? 'mdi-heart' : 'mdi-heart-outline'" />
                </v-btn>

                <span class="text-subtitle-1 ml-1">
                  {{ article.likes_count ?? 0 }}
                </span>
              </div>
            </div>

            <!-- 本文 -->
            <v-card flat rounded="lg" class="pa-8">
              <div class="text-body-1" style="white-space: pre-wrap">
                {{ article.content }}
              </div>
            </v-card>
          </template>
        </v-col>
      </v-row>
    </v-container>
  </v-sheet>
</template>
