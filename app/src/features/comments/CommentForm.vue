<script setup lang="ts">
import { computed, ref } from 'vue'
import { api } from '@/api/client'
import type { ServerReplyJSONResponse } from '@/api/generated/apiSchema'

defineOptions({
  name: 'CommentForm',
})

const props = defineProps<{
  articleId: number
  parentId: number | null
  autofocus?: boolean
}>()

const emit = defineEmits<{
  (e: 'submitted', comment: ServerReplyJSONResponse): void
}>()

const body = ref('')
const submitting = ref(false)
const submitError = ref<string | null>(null)

const canSubmit = computed(() => !submitting.value && !!body.value.trim())

const placeholder = computed(() =>
  props.parentId === null ? 'コメントを書く' : 'このコメントに返信する',
)

const submit = async () => {
  if (submitting.value) return
  if (!canSubmit.value) return

  submitting.value = true
  submitError.value = null

  try {
    const { data } = await api.api.articlesRepliesCreate(
      props.articleId,
      {
        parent_id: props.parentId ?? undefined,
        kind: 'comment',
        body: body.value,
      },
      { skipGlobalErrorHandler: true },
    )
    emit('submitted', data)
    body.value = ''
  } catch {
    submitError.value = '投稿に失敗しました。もう一度お試しください。'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <v-form class="d-flex flex-column ga-3" @submit.prevent="submit">
    <v-alert
      v-if="submitError"
      type="error"
      density="compact"
      closable
      @click:close="submitError = null"
    >
      {{ submitError }}
    </v-alert>

    <v-textarea
      v-model="body"
      :placeholder="placeholder"
      :autofocus="autofocus"
      variant="outlined"
      density="comfortable"
      rows="3"
      auto-grow
      counter
      maxlength="2000"
      hide-details="auto"
    />

    <div class="d-flex justify-end">
      <v-btn
        type="submit"
        color="primary"
        :loading="submitting"
        :disabled="!canSubmit"
      >
        投稿する
      </v-btn>
    </div>
  </v-form>
</template>
