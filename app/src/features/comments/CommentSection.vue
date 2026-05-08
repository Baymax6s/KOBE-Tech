<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'
import { api } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import type { ServerReplyJSONResponse } from '@/api/generated/apiSchema'
import CommentForm from './CommentForm.vue'
import CommentThread from './CommentThread.vue'

defineOptions({
  name: 'CommentSection',
})

const props = defineProps<{ articleId: number }>()

const { isAuthenticated } = storeToRefs(useAuthStore())
const route = useRoute()
const loginRedirect = computed(
  () => `/login?redirect=${encodeURIComponent(route.fullPath)}`,
)

const comments = ref<ServerReplyJSONResponse[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const childrenByParent = computed(() => {
  const map = new Map<number, ServerReplyJSONResponse[]>()
  for (const c of comments.value) {
    if (c.parent_id === null) continue
    const list = map.get(c.parent_id) ?? []
    list.push(c)
    map.set(c.parent_id, list)
  }
  for (const list of map.values()) {
    list.sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    )
  }
  return map
})

const rootComments = computed(() =>
  comments.value
    .filter((c) => c.parent_id === null)
    .slice()
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    ),
)

const isCommentKind = (c: ServerReplyJSONResponse) => c.kind === 'comment'

const fetchComments = async (id: number) => {
  loading.value = true
  error.value = null
  try {
    const { data } = await api.api.articleRepliesList(id)
    comments.value = (data.replies ?? []).filter(isCommentKind)
  } catch {
    error.value = 'コメントの取得に失敗しました'
  } finally {
    loading.value = false
  }
}

const handleSubmitted = (newComment: ServerReplyJSONResponse) => {
  if (!isCommentKind(newComment)) return
  comments.value = [...comments.value, newComment]
}

watch(
  () => props.articleId,
  (id) => {
    comments.value = []
    void fetchComments(id)
  },
  { immediate: true },
)
</script>

<template>
  <section class="d-flex flex-column ga-6">
    <h2 class="text-h6 font-weight-bold">
      コメント
      <span v-if="!loading" class="text-medium-emphasis text-body-2 ml-1">
        ({{ comments.length }})
      </span>
    </h2>

    <v-card flat rounded="lg" class="pa-4">
      <CommentForm
        v-if="isAuthenticated"
        :article-id="articleId"
        :parent-id="null"
        @submitted="handleSubmitted"
      />
      <div v-else class="text-body-2 d-flex align-center ga-2">
        <v-icon icon="mdi-lock-outline" size="small" />
        <span>コメントを投稿するには</span>
        <RouterLink :to="loginRedirect" class="text-primary">
          ログイン
        </RouterLink>
        <span>してください</span>
      </div>
    </v-card>

    <v-skeleton-loader
      v-if="loading"
      type="paragraph, paragraph"
      class="bg-transparent"
    />

    <v-alert v-else-if="error" type="error" closable>
      {{ error }}
      <template #append>
        <v-btn
          variant="text"
          size="small"
          @click="fetchComments(articleId)"
        >
          再試行
        </v-btn>
      </template>
    </v-alert>

    <div
      v-else-if="rootComments.length === 0"
      class="text-body-2 text-medium-emphasis text-center py-6"
    >
      まだコメントはありません
    </div>

    <div v-else class="d-flex flex-column ga-4">
      <CommentThread
        v-for="comment in rootComments"
        :key="comment.id"
        :comment="comment"
        :children-by-parent="childrenByParent"
        :depth="0"
        :article-id="articleId"
        @submitted="handleSubmitted"
      />
    </div>
  </section>
</template>
