<script setup lang="ts">
import { computed, ref } from 'vue'
import { api } from '@/api/client'
import type { ServerReplyJSONResponse } from '@/api/generated/apiSchema'

type ReplyKind = 'comment' | 'question' | 'answer'

defineOptions({
  name: 'ReplyForm',
})

const props = defineProps<{
  articleId: number
  parentId: number | null
  parentKind: ReplyKind | null
  autofocus?: boolean
}>()

const emit = defineEmits<{
  (e: 'submitted', reply: ServerReplyJSONResponse): void
}>()

const body = ref('')
const submitting = ref(false)
const submitError = ref<string | null>(null)

// 記事直下投稿のときだけユーザーが選べる。返信時は親 kind から自動導出する。
const selectedKind = ref<'comment' | 'question'>('comment')

const resolvedKind = computed<ReplyKind>(() => {
  if (props.parentKind === null) return selectedKind.value
  if (props.parentKind === 'comment') return 'comment'
  return 'answer'
})

const placeholder = computed(() => {
  if (props.parentKind === null) {
    return selectedKind.value === 'question' ? '質問を書く' : 'コメントを書く'
  }
  if (props.parentKind === 'comment') return 'このコメントに返信する'
  return '回答を書く'
})

const canSubmit = computed(() => !submitting.value && !!body.value.trim())

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
        kind: resolvedKind.value,
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

    <v-btn-toggle
      v-if="parentKind === null"
      v-model="selectedKind"
      mandatory
      density="comfortable"
      color="primary"
    >
      <v-btn value="comment">コメント</v-btn>
      <v-btn value="question">質問</v-btn>
    </v-btn-toggle>

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
