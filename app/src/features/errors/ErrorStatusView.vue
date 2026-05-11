<script setup lang="ts">
import { useRouter } from 'vue-router'

defineProps<{
  statusCode: number
  title: string
  message: string
}>()

const router = useRouter()
</script>

<template>
  <v-container class="error-status py-12" fluid>
    <div class="error-status__content">
      <p class="error-status__code">{{ statusCode }}</p>
      <h1 class="error-status__title">{{ title }}</h1>
      <p class="error-status__message">{{ message }}</p>

      <div class="error-status__actions">
        <v-btn
          color="primary"
          size="large"
          prepend-icon="mdi-home"
          @click="router.push('/')"
        >
          ホームへ戻る
        </v-btn>
        <v-btn
          variant="outlined"
          size="large"
          prepend-icon="mdi-refresh"
          @click="router.go(0)"
        >
          再読み込み
        </v-btn>
      </div>
    </div>
  </v-container>
</template>

<style scoped>
.error-status {
  min-height: calc(100vh - 64px);
  display: grid;
  place-items: center;
}

.error-status__content {
  width: min(100%, 560px);
  text-align: center;
  padding: 32px 20px;
}

.error-status__code {
  margin: 0 0 12px;
  color: rgb(var(--v-theme-primary));
  font-size: 0.875rem;
  font-weight: 700;
}

.error-status__title {
  margin: 0;
  font-size: 2rem;
  line-height: 1.2;
  font-weight: 700;
}

.error-status__message {
  margin: 16px auto 0;
  max-width: 36rem;
  color: rgba(var(--v-theme-on-surface), 0.72);
  font-size: 1rem;
  line-height: 1.8;
}

.error-status__actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 12px;
  margin-top: 32px;
}

@media (min-width: 600px) {
  .error-status__title {
    font-size: 2.5rem;
  }
}

@media (min-width: 960px) {
  .error-status__title {
    font-size: 3rem;
  }
}
</style>
