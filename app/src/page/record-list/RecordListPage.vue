<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { fetchRecordList, type RecordListResponse } from '@/api/record-list/fetchRecordList'
import RecordListGraph from '@/components/record-list/RecordListGraph.vue'
import RecordListSummary from '@/components/record-list/RecordListSummary.vue'
import RecordListTable from '@/components/record-list/RecordListTable.vue'

const recordList = ref<RecordListResponse | null>(null)

onMounted(async () => {
  recordList.value = await fetchRecordList()
})
</script>

<template>
  <main class="page">
    <section class="hero">
      <p class="eyebrow">Record List</p>
      <h1>page / components / api に寄せた簡素な frontend</h1>
      <p class="lead">
        画面は `page/record-list` にまとめ、表示部品は `components/record-list`、データ取得は
        `api/record-list` に寄せています。
      </p>
    </section>

    <template v-if="recordList">
      <RecordListSummary
        :daily-average-minutes="recordList.dailyAverageMinutes"
        :streak-days="recordList.streakDays"
        :completion-rate="recordList.completionRate"
      />
      <RecordListGraph :records="recordList.records" />
      <RecordListTable :records="recordList.records" />
    </template>
  </main>
</template>

<style scoped>
.page {
  width: min(1120px, calc(100% - 32px));
  margin: 0 auto;
  padding: 48px 0 96px;
  display: grid;
  gap: 24px;
}

.hero {
  padding: 28px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
  box-shadow: var(--shadow-soft);
}

.eyebrow {
  font-size: 0.76rem;
  font-weight: 700;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--color-accent);
}

h1 {
  margin-top: 12px;
  max-width: 12ch;
  font-size: clamp(2.8rem, 6.5vw, 5rem);
  line-height: 0.95;
  color: var(--color-ink);
}

.lead {
  margin-top: 18px;
  max-width: 42rem;
  color: var(--color-muted);
}

@media (min-width: 768px) {
  .page {
    width: min(1120px, calc(100% - 56px));
    padding: 56px 0 120px;
    gap: 28px;
  }

  .hero {
    padding: 40px;
  }
}
</style>
