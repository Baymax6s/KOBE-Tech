<script setup lang="ts">
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import type { ServerReplyJSONResponse } from '@/api/generated/apiSchema'
import CommentItem from './CommentItem.vue'
import CommentForm from './CommentForm.vue'

defineOptions({
  name: 'CommentThread',
})

const props = defineProps<{
  comment: ServerReplyJSONResponse
  childrenByParent: Map<number, ServerReplyJSONResponse[]>
  depth: number
  articleId: number
}>()

const emit = defineEmits<{
  (e: 'submitted', comment: ServerReplyJSONResponse): void
}>()

const { isAuthenticated } = storeToRefs(useAuthStore())

const children = computed<ServerReplyJSONResponse[]>(
  () => props.childrenByParent.get(props.comment.id) ?? [],
)

const expanded = ref(props.depth === 0)
const showReplyForm = ref(false)

const toggleReplyForm = () => {
  showReplyForm.value = !showReplyForm.value
}

const handleSubmitted = (newComment: ServerReplyJSONResponse) => {
  emit('submitted', newComment)
  showReplyForm.value = false
  expanded.value = true
}
</script>

<template>
  <div class="d-flex flex-column ga-3">
    <CommentItem
      :comment="comment"
      :can-reply="isAuthenticated"
      :replying="showReplyForm"
      @toggle-reply="toggleReplyForm"
    />

    <div v-if="showReplyForm" class="ml-8">
      <CommentForm
        :article-id="articleId"
        :parent-id="comment.id"
        autofocus
        @submitted="handleSubmitted"
      />
    </div>

    <div v-if="children.length > 0" class="ml-8">
      <div v-if="expanded" class="d-flex flex-column ga-3">
        <CommentThread
          v-for="child in children"
          :key="child.id"
          :comment="child"
          :children-by-parent="childrenByParent"
          :depth="depth + 1"
          :article-id="articleId"
          @submitted="(c) => emit('submitted', c)"
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
        返信 {{ children.length }} 件を表示
      </v-btn>
    </div>
  </div>
</template>
