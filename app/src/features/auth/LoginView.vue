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
            outlined
            dense
            required
            class="mb-4"
          />

          <v-text-field
            v-model="password"
            label="パスワード"
            type="password"
            outlined
            dense
            required
            class="mb-4"
          />

          <v-btn
            type="submit"
            color="primary"
            block
          >
            ログイン
          </v-btn>
        </v-form>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const name = ref('')
const password = ref('')

const login = () => {
  console.log('名前:', name.value)
  console.log('パスワード:', password.value)
}
</script>