<script setup lang="ts">
import { useDateFormat } from '@vueuse/core'
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

const formattedDate = useDateFormat(
  () => props.comment.created_at,
  'YYYY/MM/DD HH:mm',
)
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
