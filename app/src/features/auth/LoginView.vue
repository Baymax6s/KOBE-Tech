<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <v-card class="w-full max-w-md p-6">
      <v-card-title class="text-h6 text-center">
        ログイン
      </v-card-title>

      <v-card-text>
        <v-form @submit.prevent="login">
          <v-text-field
            v-model="name"
            label="ユーザー名"
            variant="outlined"
            required
            class="mb-4"
          />

          <v-text-field
            v-model="password"
            label="パスワード"
            type="password"
            variant="outlined"
            required
            class="mb-4"
          />

          <v-btn
            color="primary"
            block
            :loading="loading"
            @click="login"
          >
          ログイン
          </v-btn>

          <div v-if="error" class="mt-4 text-red-500">
            {{ error }}
          </div>
        </v-form>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { api } from '@/api/client'

const name = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const login = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await api.api.authLoginCreate({
      name: name.value,
      password: password.value,
    })
    if (response.data.token) {
      localStorage.setItem('token', response.data.token)
      // ログイン成功の処理（例: リダイレクト）
      console.log('ログイン成功')
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'ログインに失敗しました'
  } finally {
    loading.value = false
  }
}
</script>