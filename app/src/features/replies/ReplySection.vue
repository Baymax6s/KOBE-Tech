<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'
import { api } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import type { ServerReplyJSONResponse } from '@/api/generated/apiSchema'
import ReplyForm from './ReplyForm.vue'
import ReplyThread from './ReplyThread.vue'

defineOptions({
  name: 'ReplySection',
})

const props = defineProps<{ articleId: number }>()

const { isAuthenticated } = storeToRefs(useAuthStore())
const route = useRoute()
const loginRedirect = computed(
  () => `/login?redirect=${encodeURIComponent(route.fullPath)}`,
)

const replies = ref<ServerReplyJSONResponse[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const currentUserId = ref<number | null>(null)

onMounted(async () => {
  if (!isAuthenticated.value) return
  try {
    const res = await api.api.authMeList({ skipGlobalErrorHandler: true })
    currentUserId.value = res.data.id ?? null
  } catch {
    currentUserId.value = null
  }
})

const childrenByParent = computed(() => {
  const map = new Map<number, ServerReplyJSONResponse[]>()
  for (const r of replies.value) {
    if (r.parent_id == null) continue
    const list = map.get(r.parent_id) ?? []
    list.push(r)
    map.set(r.parent_id, list)
  }
  for (const list of map.values()) {
    list.sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    )
  }
  return map
})

// 各リプライの「自分以下にぶら下がる返信の総数」をメモ化付き DFS で一括計算する。
// 「返信 N 件を表示」ボタンは展開時にサブツリー全体を開く挙動なので、N もサブツリー全体の件数に揃える。
const descendantCountByParent = computed(() => {
  const counts = new Map<number, number>()
  const compute = (id: number): number => {
    const cached = counts.get(id)
    if (cached !== undefined) return cached
    const kids = childrenByParent.value.get(id) ?? []
    let total = 0
    for (const k of kids) total += 1 + compute(k.id)
    counts.set(id, total)
    return total
  }
  for (const r of replies.value) compute(r.id)
  return counts
})

const rootReplies = computed(() =>
  replies.value
    .filter((r) => r.parent_id == null)
    .slice()
    .sort(
      (a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
    ),
)

// 各リプライについて、その親が質問（kind=question）なら親の user_id を記録する Map。
// key = 子の reply.id, value = 質問の user_id。
const questionAuthorByReplyId = computed(() => {
  const map = new Map<number, number>()
  const replyMap = new Map(replies.value.map((r) => [r.id, r]))
  for (const r of replies.value) {
    if (r.parent_id == null) continue
    const parent = replyMap.get(r.parent_id)
    if (parent && parent.kind === 'question') {
      map.set(r.id, parent.user_id)
    }
  }
  return map
})

const fetchReplies = async (id: number) => {
  loading.value = true
  error.value = null
  try {
    const { data } = await api.api.articlesRepliesList(id, {
      skipGlobalErrorHandler: true,
    })
    replies.value = data.replies ?? []
  } catch {
    error.value = 'リプライの取得に失敗しました'
  } finally {
    loading.value = false
  }
}

const handleSubmitted = (newReply: ServerReplyJSONResponse) => {
  replies.value = [...replies.value, newReply]
}

const handleBestUpdated = (replyId: number, isBest: boolean) => {
  replies.value = replies.value.map((r) =>
    r.id === replyId ? { ...r, is_best: isBest } : r,
  )
}

watch(
  () => props.articleId,
  (id) => {
    replies.value = []
    void fetchReplies(id)
  },
  { immediate: true },
)
</script>

<template>
  <section class="d-flex flex-column ga-6">
    <h2 class="text-h6 font-weight-bold">
      リプライ
      <span v-if="!loading" class="text-medium-emphasis text-body-2 ml-1">
        ({{ replies.length }})
      </span>
    </h2>

    <v-card flat rounded="lg" class="pa-4">
      <ReplyForm
        v-if="isAuthenticated"
        :article-id="articleId"
        :parent-id="null"
        :parent-kind="null"
        @submitted="handleSubmitted"
      />
      <div v-else class="text-body-2 d-flex align-center ga-2">
        <v-icon icon="mdi-lock-outline" size="small" />
        <span>リプライを投稿するには</span>
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
        <v-btn variant="text" size="small" @click="fetchReplies(articleId)">
          再試行
        </v-btn>
      </template>
    </v-alert>

    <div
      v-else-if="rootReplies.length === 0"
      class="text-body-2 text-medium-emphasis text-center py-6"
    >
      まだリプライはありません
    </div>

    <div v-else class="d-flex flex-column ga-4">
      <ReplyThread
        v-for="reply in rootReplies"
        :key="reply.id"
        :reply="reply"
        :children-by-parent="childrenByParent"
        :descendant-count-by-parent="descendantCountByParent"
        :depth="0"
        :article-id="articleId"
        :current-user-id="currentUserId"
        :question-author-by-reply-id="questionAuthorByReplyId"
        @submitted="handleSubmitted"
        @best-updated="handleBestUpdated"
      />
    </div>
  </section>
</template>
