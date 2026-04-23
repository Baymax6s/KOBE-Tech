<script setup lang="ts">
import { useDateFormat } from '@vueuse/core'

interface Article {
  id: number
  title: string
  created_at: string
}

const props = defineProps<{ article: Article }>()

const formattedDate = useDateFormat(props.article.created_at, 'YYYY/MM/DD')
</script>

<template>
  <v-card class="p-4">
    <v-card-title class="text-base font-medium">
      {{ article.title }}
    </v-card-title>
    <v-card-subtitle class="text-sm text-gray-500">
      {{ formattedDate }}
    </v-card-subtitle>
  </v-card>
</template>
