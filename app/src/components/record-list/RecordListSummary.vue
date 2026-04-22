<script setup lang="ts">
defineProps<{
  dailyAverageMinutes: number
  streakDays: number
  completionRate: number
}>()
</script>

<template>
  <section class="summary">
    <article class="card">
      <p class="label">daily average</p>
      <strong>{{ dailyAverageMinutes }} min</strong>
      <span>1 日あたりの学習時間</span>
    </article>
    <article class="card">
      <p class="label">streak</p>
      <strong>{{ streakDays }} days</strong>
      <span>連続で記録できた日数</span>
    </article>
    <article class="card">
      <p class="label">completion</p>
      <strong>{{ completionRate }}%</strong>
      <span>集中度から見た達成感</span>
    </article>
  </section>
</template>

<style scoped>
.summary {
  display: grid;
  gap: 16px;
}

.card {
  padding: 22px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: rgba(255, 255, 255, 0.58);
  box-shadow: var(--shadow-soft);
}

.label {
  font-size: 0.76rem;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--color-accent);
}

strong {
  display: block;
  margin-top: 10px;
  font-size: 2rem;
  line-height: 1;
  color: var(--color-ink);
}

span {
  display: block;
  margin-top: 10px;
  color: var(--color-muted);
}

@media (min-width: 768px) {
  .summary {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}
</style>
