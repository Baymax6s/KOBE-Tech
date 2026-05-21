<script setup lang="ts">
import { computed, ref } from 'vue'
import { useTimeAgo, type UseTimeAgoMessages } from '@vueuse/core'
import { api } from '@/api/client'
import type { ServerReplyJSONResponse } from '@/api/generated/apiSchema'

defineOptions({
  name: 'ReplyItem',
})

const props = defineProps<{
  reply: ServerReplyJSONResponse
  canReply: boolean
  replying: boolean
  currentUserId: number | null
  questionAuthorByReplyId: Map<number, number>
}>()

const emit = defineEmits<{
  (e: 'toggle-reply'): void
  (e: 'best-updated', replyId: number, isBest: boolean): void
}>()

const kindBadge = computed(() => {
  switch (props.reply.kind) {
    case 'question':
      return { color: 'warning', label: '質問' }
    case 'answer':
      return { color: 'success', label: '回答' }
    default:
      return { color: 'primary', label: 'コメント' }
  }
})

// 返信ボタン文言は親（=この reply）の kind から決める。
// コメント配下は「コメントする」、質問/回答配下は「回答する」。
const replyActionLabel = computed(() =>
  props.reply.kind === 'comment' ? 'コメントする' : '回答する',
)

const canMarkBest = computed(() => {
  if (props.reply.kind !== 'answer') return false
  if (props.currentUserId == null) return false
  const questionUserId = props.questionAuthorByReplyId.get(props.reply.id)
  return questionUserId != null && questionUserId === props.currentUserId
})

const submittingBest = ref(false)
const bestError = ref<string | null>(null)

const markAsBest = async () => {
  if (submittingBest.value) return
  bestError.value = null
  submittingBest.value = true
  try {
    await api.api.repliesBestCreate(props.reply.id, {
      skipGlobalErrorHandler: true,
    })
    emit('best-updated', props.reply.id, true)
  } catch {
    bestError.value = 'ベストアンサーの設定に失敗しました'
  } finally {
    submittingBest.value = false
  }
}

const unmarkBest = async () => {
  if (submittingBest.value) return
  bestError.value = null
  submittingBest.value = true
  try {
    await api.api.repliesBestUpdate(props.reply.id, {
      skipGlobalErrorHandler: true,
    })
    emit('best-updated', props.reply.id, false)
  } catch {
    bestError.value = 'ベストアンサーの解除に失敗しました'
  } finally {
    submittingBest.value = false
  }
}

// 投稿から 6 日以内は「N分前 / N時間前 / N日前」と表示し、7 日以上経過したら絶対日付に切り替える。
// useTimeAgo は内部で setInterval により値をリアクティブに更新するので、画面を開きっぱなしでも表示が古びない。
const formattedDate = useTimeAgo(() => props.reply.created_at, {
  max: 6 * 24 * 60 * 60 * 1000,
  fullDateFormatter: (date) => {
    const y = date.getFullYear()
    const m = String(date.getMonth() + 1).padStart(2, '0')
    const d = String(date.getDate()).padStart(2, '0')
    return `${y}/${m}/${d}`
  },
  // VueUse の messages 型は builtin と Record の交差型で TS が past/future を unit と同列に扱おうとして衝突する。
  // 値の中身は仕様どおりなので、構築済みオブジェクトを UseTimeAgoMessages にアサートして渡す。
  messages: {
    justNow: 'たった今',
    past: (n) => (/\d/.test(n) ? `${n}前` : n),
    future: (n) => (/\d/.test(n) ? `${n}後` : n),
    year: (n) => `${n}年`,
    month: (n) => `${n}か月`,
    week: (n) => `${n}週間`,
    day: (n) => `${n}日`,
    hour: (n) => `${n}時間`,
    minute: (n) => `${n}分`,
    second: (n) => `${n}秒`,
    invalid: '',
  } as UseTimeAgoMessages,
})
</script>

<template>
  <v-card
    flat
    rounded="lg"
    class="pa-4"
    :class="reply.is_best ? 'bg-yellow-lighten-4' : 'bg-grey-lighten-5'"
  >
    <div class="d-flex align-center ga-2 mb-2">
      <v-chip :color="kindBadge.color" size="small" variant="tonal" label>
        {{ kindBadge.label }}
      </v-chip>
      <v-chip
        v-if="reply.is_best"
        color="amber"
        size="small"
        variant="flat"
        label
      >
        <v-icon start icon="mdi-crown" size="x-small" />
        ベストアンサー
      </v-chip>
      <span class="text-body-2 font-weight-medium">
        {{ reply.user_name }}
      </span>
      <span class="text-body-2 text-medium-emphasis">
        {{ formattedDate }}
      </span>
    </div>

    <div class="text-body-2 mb-2" style="white-space: pre-wrap">
      {{ reply.body }}
    </div>

    <v-alert
      v-if="bestError"
      type="error"
      density="compact"
      variant="tonal"
      class="mb-2"
      closable
      @click:close="bestError = null"
    >
      {{ bestError }}
    </v-alert>

    <div class="d-flex ga-2">
      <div v-if="canReply">
        <v-btn
          variant="text"
          size="small"
          color="primary"
          :prepend-icon="replying ? 'mdi-close' : 'mdi-reply'"
          @click="$emit('toggle-reply')"
        >
          {{ replying ? 'キャンセル' : replyActionLabel }}
        </v-btn>
      </div>

      <div v-if="canMarkBest">
        <v-btn
          v-if="reply.is_best"
          variant="text"
          size="small"
          color="amber"
          :loading="submittingBest"
          prepend-icon="mdi-crown-off"
          @click="unmarkBest"
        >
          取り消す
        </v-btn>
        <v-btn
          v-else
          variant="text"
          size="small"
          color="amber"
          :loading="submittingBest"
          prepend-icon="mdi-crown-outline"
          @click="markAsBest"
        >
          ベストアンサーに選ぶ
        </v-btn>
      </div>
    </div>
  </v-card>
</template>
