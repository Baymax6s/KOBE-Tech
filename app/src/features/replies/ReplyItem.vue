<script setup lang="ts">
import { computed } from 'vue'
import { useTimeAgo, type UseTimeAgoMessages } from '@vueuse/core'
import type { ServerReplyJSONResponse } from '@/api/generated/apiSchema'

defineOptions({
  name: 'ReplyItem',
})

const props = defineProps<{
  reply: ServerReplyJSONResponse
  canReply: boolean
  replying: boolean
  isQuestionAuthor: boolean
}>()

defineEmits<{
  (e: 'toggle-reply'): void
  (e: 'toggle-best'): void
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
    :class="[
      'pa-4',
      reply.is_best ? 'bg-amber-lighten-5' : 'bg-grey-lighten-5',
    ]"
    :style="reply.is_best ? { borderLeft: '4px solid #FFC107' } : {}"
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
        <v-icon start icon="mdi-star" size="x-small" />
        ベストアンサー
      </v-chip>
      <span class="text-body-2 font-weight-medium">
        {{ reply.user_name }}
      </span>
      <span class="text-caption text-medium-emphasis">
        {{ formattedDate }}
      </span>
    </div>

    <div class="text-body-2 mb-2" style="white-space: pre-wrap">
      {{ reply.body }}
    </div>

    <div class="d-flex align-center ga-2">
      <v-btn
        v-if="canReply"
        variant="text"
        size="small"
        color="primary"
        :prepend-icon="replying ? 'mdi-close' : 'mdi-reply'"
        @click="$emit('toggle-reply')"
      >
        {{ replying ? 'キャンセル' : replyActionLabel }}
      </v-btn>
      <v-btn
        v-if="isQuestionAuthor && reply.kind === 'answer'"
        variant="text"
        size="small"
        :color="reply.is_best ? 'amber' : 'grey'"
        prepend-icon="mdi-star"
        @click="$emit('toggle-best')"
      >
        {{ reply.is_best ? 'ベストアンサーを解除' : 'ベストアンサーに選択' }}
      </v-btn>
    </div>
  </v-card>
</template>
