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
  depth: number
  articleId: number
  isQuestionAuthor: boolean
}>()

const emit = defineEmits<{
  (e: 'submitted', reply: ServerReplyJSONResponse): void
  (e: 'toggle-best', replyId: number): void
}>()

const { isAuthenticated } = storeToRefs(useAuthStore())

const children = computed<ServerReplyJSONResponse[]>(
  () => props.childrenByParent.get(props.reply.id) ?? [],
)

const descendantCount = computed(
  () => props.descendantCountByParent.get(props.reply.id) ?? 0,
)

// depth 0 のみ初期折りたたみ。クリックで自分以下のサブツリーをまとめて開く。
const expanded = ref(props.depth >= 1)
const showReplyForm = ref(false)

const toggleReplyForm = () => {
  showReplyForm.value = !showReplyForm.value
}

const handleSubmitted = (newReply: ServerReplyJSONResponse) => {
  emit('submitted', newReply)
  showReplyForm.value = false
  expanded.value = true
}

const handleToggleBest = (replyId: number) => {
  emit('toggle-best', replyId)
}
</script>

<template>
  <div class="d-flex flex-column ga-3">
    <ReplyItem
      :reply="reply"
      :can-reply="isAuthenticated"
      :replying="showReplyForm"
      :is-question-author="isQuestionAuthor"
      @toggle-reply="toggleReplyForm"
      @toggle-best="handleToggleBest(reply.id)"
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
      <div v-if="expanded" class="d-flex flex-column ga-3">
        <ReplyThread
          v-for="child in children"
          :key="child.id"
          :reply="child"
          :children-by-parent="childrenByParent"
          :descendant-count-by-parent="descendantCountByParent"
          :depth="depth + 1"
          :article-id="articleId"
          :is-question-author="isQuestionAuthor"
          @submitted="emit('submitted', $event)"
          @toggle-best="emit('toggle-best', $event)"
        />
      </div>
      <v-btn
        v-else
        variant="text"
        size="small"
        color="primary"
        prepend-icon="mdi-chevron-down"
        @click="expanded = true"
      >
        返信 {{ descendantCount }} 件を表示
      </v-btn>
    </div>
  </div>
</template>
