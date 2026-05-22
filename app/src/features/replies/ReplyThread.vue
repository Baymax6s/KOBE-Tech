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
  maxDepth: number
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

const localReveal = ref(false)

// 3階層目（depth === 2）以上のときに制限をかけます
const isMaxDepth = computed(() => props.depth >= 2)

const effectiveReveal = computed(() => props.revealAll || localReveal.value)

const visibleChildren = computed(() => {
  if (isMaxDepth.value && !effectiveReveal.value) {
    return []
  }
  return effectiveReveal.value
    ? children.value
    : children.value.filter((c) => props.bestAnswerPathIds.has(c.id))
})

const hiddenCount = computed(() => {
  if (effectiveReveal.value) return 0
  if (isMaxDepth.value) {
    return props.descendantCountByParent.get(props.reply.id) ?? 0
  }
  return props.hiddenDescendantCountByReplyId.get(props.reply.id) ?? 0
})

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
      :can-reply="!isMaxDepth && isAuthenticated"
      :replying="showReplyForm"
      :current-user-id="currentUserId"
      :question-author-by-reply-id="questionAuthorByReplyId"
      @toggle-reply="toggleReplyForm"
      @best-updated="handleBestUpdated"
    />

    <div v-if="showReplyForm" :class="{ 'ml-8': !isMaxDepth }">
      <ReplyForm
        :article-id="articleId"
        :parent-id="reply.id"
        :parent-kind="reply.kind"
        autofocus
        @submitted="handleSubmitted"
      />
    </div>

    <div v-if="children.length > 0" :class="{ 'ml-8': !isMaxDepth }">
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
          :max-depth="maxDepth"
          :article-id="articleId"
          :current-user-id="currentUserId"
          :question-author-by-reply-id="questionAuthorByReplyId"
          @submitted="emit('submitted', $event)"
          @best-updated="handleBestUpdated"
        />
      </div>

      <v-btn
        v-if="(depth === 0 || !isMaxDepth) && hiddenCount > 0"
        variant="text"
        size="small"
        color="primary"
        prepend-icon="mdi-chevron-down"
        @click="localReveal = true"
      >
        返信 {{ hiddenCount }} 件を表示
      </v-btn>
      <v-btn
        v-else-if="(depth === 0 || !isMaxDepth) && localReveal"
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