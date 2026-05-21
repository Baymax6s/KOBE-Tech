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
  <v-container class="fill-height">
    <v-row justify="center" align="center">
      <v-col cols="12" sm="10" class="text-center">
        <p class="text-caption font-weight-bold text-primary mb-3">
          {{ statusCode }}
        </p>
        <h1 class="text-h3 font-weight-bold mb-4">{{ title }}</h1>
        <p class="text-body-1 text-medium-emphasis mb-8 mx-auto" style="max-width: 36rem">
          {{ message }}
        </p>

        <div class="d-flex justify-center ga-3">
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
      </v-col>
    </v-row>
  </v-container>
</template>
