<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { ServerTagRankingJSONResponse } from '@/api/generated/apiSchema'
import { api } from '@/api/client'

const rankings = ref<ServerTagRankingJSONResponse[]>([])
const loading = ref(false)

// 順位ごとのメダル色。RANK() は同点で番号が飛ぶので 1/2/3 以外（同点4位など）は既定色にフォールバックする。
const rankColors: Record<number, string> = {
  1: 'amber',
  2: 'blue-grey',
  3: 'brown',
}

onMounted(async () => {
  loading.value = true
  try {
    const res = await api.api.tagsPopularArticlesList()
    rankings.value = res.data.rankings ?? []
  } catch {
    rankings.value = []
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <section v-if="rankings.length" class="mb-6">
    <h2 class="text-h6 font-weight-bold mb-3">カテゴリ別の注目記事</h2>

    <v-slide-group show-arrows>
      <v-slide-group-item v-for="ranking in rankings" :key="ranking.tag_id">
        <v-card width="280" class="me-4 pa-2" variant="outlined">
          <v-card-subtitle class="d-flex align-center ga-1 px-2 pt-2">
            <v-icon icon="mdi-tag" size="small" />
            <span class="font-weight-medium">{{ ranking.tag_name }}</span>
          </v-card-subtitle>

          <v-list density="compact" lines="one">
            <v-list-item
              v-for="popular in ranking.articles"
              :key="popular.article_id"
              :to="`/articles/${popular.article_id}`"
            >
              <template #prepend>
                <v-chip
                  size="small"
                  variant="flat"
                  :color="rankColors[popular.rank]"
                  class="me-2"
                >
                  {{ popular.rank }}
                </v-chip>
              </template>

              <v-list-item-title>{{ popular.title }}</v-list-item-title>

              <template #append>
                <span
                  class="d-flex align-center text-caption text-medium-emphasis"
                >
                  <v-icon
                    icon="mdi-heart"
                    size="x-small"
                    color="red"
                    class="me-1"
                  />
                  {{ popular.likes_count }}
                </span>
              </template>
            </v-list-item>
          </v-list>
        </v-card>
      </v-slide-group-item>
    </v-slide-group>
  </section>
</template>
