<script setup lang="ts">
import type { StudyRecord } from '@/api/record-list/fetchRecordList'

defineProps<{
  records: StudyRecord[]
}>()
</script>

<template>
  <section class="table-card">
    <div class="heading">
      <div>
        <p class="eyebrow">table</p>
        <h2>記録一覧</h2>
      </div>
    </div>

    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>date</th>
            <th>study</th>
            <th>focus</th>
            <th>note</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in records" :key="record.id">
            <td>{{ record.date }}</td>
            <td>{{ record.studyMinutes }} min</td>
            <td>{{ record.focusScore }}</td>
            <td>{{ record.note }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<style scoped>
.table-card {
  padding: 24px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-surface);
  box-shadow: var(--shadow-soft);
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

.table-wrap {
  margin-top: 20px;
  overflow-x: auto;
}

table {
  width: 100%;
  min-width: 640px;
}

th,
td {
  padding: 14px 12px;
  border-bottom: 1px solid rgba(15, 23, 32, 0.08);
  text-align: left;
}

th {
  font-size: 0.76rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--color-muted);
}

td {
  color: var(--color-ink);
}

tbody tr:last-child td {
  border-bottom: none;
}
</style>
