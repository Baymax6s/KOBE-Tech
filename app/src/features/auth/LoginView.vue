<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <v-card class="w-full max-w-md p-6">
      <v-card-title class="text-h6 text-center">
        ログイン
      </v-card-title>

      <v-card-text>
        <v-alert v-if="errorMessage" type="error" class="mb-4" closable @click:close="errorMessage = null">
          {{ errorMessage }}
        </v-alert>

        <v-form @submit.prevent="onSubmit">
          <v-text-field
            v-model="name"
            label="ユーザー名"
            variant="outlined"
            required
            class="mb-4"
            :disabled="submitting"
          />

          <v-text-field
            v-model="password"
            label="パスワード"
            type="password"
            variant="outlined"
            required
            class="mb-4"
            :disabled="submitting"
          />

          <v-btn
            type="submit"
            color="primary"
            block
            :loading="submitting"
            :disabled="!name || !password || submitting"
          >
            ログイン
          </v-btn>
        </v-form>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const name = ref('')
const password = ref('')
const submitting = ref(false)
const errorMessage = ref<string | null>(null)

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const onSubmit = async () => {
  submitting.value = true
  errorMessage.value = null
  try {
    await auth.login(name.value, password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/articles'
    router.push(redirect)
  } catch (e) {
    // 401 のような認証エラーはサーバのメッセージを優先、それ以外は汎用文言
    if (axios.isAxiosError(e) && e.response?.status === 401) {
      errorMessage.value = e.response.data?.message ?? 'ユーザー名またはパスワードが正しくありません'
    } else {
      errorMessage.value = 'ログインに失敗しました。時間をおいて再度お試しください'
    }
  } finally {
    submitting.value = false
  }
}
</script>
