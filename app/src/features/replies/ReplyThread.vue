<script setup lang="ts">
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import type { ServerReplyJSONResponse } from '@/api/generated/apiSchema'
import ReplyItem from './ReplyItem.vue'
import ReplyForm from './ReplyForm.vue'

defineOptions({
  name: 'ReplyThread',
})

const props = defineProps<{
  reply: ServerReplyJSONResponse
  childrenByParent: Map<number, ServerReplyJSONResponse[]>
  descendantCountByParent: Map<number, number>
  bestAnswerPathIds: Set<number>
  hiddenDescendantCountByReplyId: Map<number, number>
  revealAll: boolean
  depth: number
  articleId: number
  currentUserId: number | null
  questionAuthorByReplyId: Map<number, number>
}>()

const emit = defineEmits<{
  (e: 'submitted', reply: ServerReplyJSONResponse): void
  (e: 'best-updated', replyId: number, isBest: boolean): void
}>()

const { isAuthenticated } = storeToRefs(useAuthStore())

const children = computed<ServerReplyJSONResponse[]>(
  () => props.childrenByParent.get(props.reply.id) ?? [],
)

// localReveal はルートで「N 件を表示」が押されているかどうか。
// 「返信を隠す」で false に戻り、初期状態に戻せる。
const localReveal = ref(false)

// 親が全表示モードなら自分も全表示。ルート側でのみ localReveal を持つ運用にする。
const effectiveReveal = computed(() => props.revealAll || localReveal.value)

// 初期表示で見せる子の集合。
// 全表示モードなら children すべて、そうでなければベストアンサー経路に乗っている子だけ。
const visibleChildren = computed(() =>
  effectiveReveal.value
    ? children.value
    : children.value.filter((c) => props.bestAnswerPathIds.has(c.id)),
)

// 隠れている件数はサブツリー全体で集計済みのものを参照する。
// ネストの奥（例: ベストアンサーの下の返信）も合算したうえで、ボタン 1 つで全部開けるようにするため。
const hiddenCount = computed(() =>
  effectiveReveal.value
    ? 0
    : (props.hiddenDescendantCountByReplyId.get(props.reply.id) ?? 0),
)

const showReplyForm = ref(false)

const toggleReplyForm = () => {
  showReplyForm.value = !showReplyForm.value
}

const handleSubmitted = (newReply: ServerReplyJSONResponse) => {
  emit('submitted', newReply)
  showReplyForm.value = false
  localReveal.value = true
}

const handleBestUpdated = (replyId: number, isBest: boolean) => {
  emit('best-updated', replyId, isBest)
}
</script>

<template>
  <div class="d-flex flex-column ga-3">
    <ReplyItem
      :reply="reply"
      :can-reply="isAuthenticated"
      :replying="showReplyForm"
      :current-user-id="currentUserId"
      :question-author-by-reply-id="questionAuthorByReplyId"
      @toggle-reply="toggleReplyForm"
      @best-updated="handleBestUpdated"
    />

    <div v-if="showReplyForm" class="ml-8">
      <ReplyForm
        :article-id="articleId"
        :parent-id="reply.id"
        :parent-kind="reply.kind"
        autofocus
        @submitted="handleSubmitted"
      />
    </div>

    <div v-if="children.length > 0" class="ml-8">
      <div
        v-if="visibleChildren.length > 0"
        class="d-flex flex-column ga-3 mb-2"
      >
        <ReplyThread
          v-for="child in visibleChildren"
          :key="child.id"
          :reply="child"
          :children-by-parent="childrenByParent"
          :descendant-count-by-parent="descendantCountByParent"
          :best-answer-path-ids="bestAnswerPathIds"
          :hidden-descendant-count-by-reply-id="hiddenDescendantCountByReplyId"
          :reveal-all="effectiveReveal"
          :depth="depth + 1"
          :article-id="articleId"
          :current-user-id="currentUserId"
          :question-author-by-reply-id="questionAuthorByReplyId"
          @submitted="emit('submitted', $event)"
          @best-updated="handleBestUpdated"
        />
      </div>
      <v-btn
        v-if="depth === 0 && hiddenCount > 0"
        variant="text"
        size="small"
        color="primary"
        prepend-icon="mdi-chevron-down"
        @click="localReveal = true"
      >
        返信 {{ hiddenCount }} 件を表示
      </v-btn>
      <v-btn
        v-else-if="depth === 0 && localReveal"
        variant="text"
        size="small"
        color="primary"
        prepend-icon="mdi-chevron-up"
        @click="localReveal = false"
      >
        返信を隠す
      </v-btn>
    </div>
  </div>
</template>
