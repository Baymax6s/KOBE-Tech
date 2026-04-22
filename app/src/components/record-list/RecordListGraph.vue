<script setup lang="ts">
import { computed } from 'vue'

import type { StudyRecord } from '@/api/record-list/fetchRecordList'

const props = defineProps<{
  records: StudyRecord[]
}>()

const maxStudyMinutes = computed(() =>
  props.records.reduce((max, record) => Math.max(max, record.studyMinutes), 0),
)

const barWidth = (studyMinutes: number) => {
  if (maxStudyMinutes.value === 0) {
    return '0%'
  }

  return `${Math.max(0, Math.min((studyMinutes / maxStudyMinutes.value) * 100, 100))}%`
}
</script>

<template>
  <section class="graph-card">
    <div class="heading">
      <div>
        <p class="eyebrow">graph</p>
        <h2>学習時間の推移</h2>
      </div>
      <p class="caption">直近 {{ props.records.length }} 日</p>
    </div>

    <div class="bars">
      <div v-for="record in props.records" :key="record.id" class="bar-row">
        <div class="bar-meta">
          <span>{{ record.date.slice(5) }}</span>
          <strong>{{ record.studyMinutes }}m</strong>
        </div>
        <div class="track">
          <div class="bar" :style="{ width: barWidth(record.studyMinutes) }" />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.graph-card {
  padding: 24px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
  box-shadow: var(--shadow-soft);
}

.heading {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 12px;
}

.eyebrow {
  font-size: 0.76rem;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--color-accent);
}

h2 {
  margin-top: 8px;
  font-size: 1.6rem;
  color: var(--color-ink);
}

.caption {
  color: var(--color-muted);
}

.bars {
  margin-top: 22px;
  display: grid;
  gap: 14px;
}

.bar-row {
  display: grid;
  gap: 8px;
}

.bar-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  color: var(--color-muted);
}

.bar-meta strong {
  color: var(--color-ink);
}

.track {
  overflow: hidden;
  height: 16px;
  border-radius: 999px;
  background: rgba(15, 23, 32, 0.08);
}

.bar {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #d86d46, #264a60);
}
</style>
