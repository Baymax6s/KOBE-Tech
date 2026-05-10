<script setup lang="ts">
import { useTimeAgo, type UseTimeAgoMessages } from '@vueuse/core'
import type { ServerReplyJSONResponse } from '@/api/generated/apiSchema'

defineOptions({
  name: 'CommentItem',
})

const props = defineProps<{
  comment: ServerReplyJSONResponse
  canReply: boolean
  replying: boolean
}>()

defineEmits<{
  (e: 'toggle-reply'): void
}>()

// 投稿から 6 日以内は「N分前 / N時間前 / N日前」と表示し、7 日以上経過したら絶対日付に切り替える。
// useTimeAgo は内部で setInterval により値をリアクティブに更新するので、画面を開きっぱなしでも表示が古びない。
const formattedDate = useTimeAgo(() => props.comment.created_at, {
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
  <v-card flat rounded="lg" class="pa-4 bg-grey-lighten-5">
    <div class="d-flex align-center ga-2 mb-2">
      <v-chip color="primary" size="small" variant="tonal" label>
        コメント
      </v-chip>
      <span class="text-body-2 font-weight-medium">
        {{ comment.user_name }}
      </span>
      <span class="text-caption text-medium-emphasis">
        {{ formattedDate }}
      </span>
    </div>

    <div class="text-body-2 mb-2" style="white-space: pre-wrap">{{
      comment.body
    }}</div>

    <div v-if="canReply" class="d-flex">
      <v-btn
        variant="text"
        size="small"
        color="primary"
        :prepend-icon="replying ? 'mdi-close' : 'mdi-reply'"
        @click="$emit('toggle-reply')"
      >
        {{ replying ? 'キャンセル' : 'コメントする' }}
      </v-btn>
    </div>
  </v-card>
</template>
