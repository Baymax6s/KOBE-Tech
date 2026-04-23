<template>
  <v-app-bar app color="primary" dark>
    <v-img
      src="/ChatGPT_Image_2026422_23_48_11.png"
      alt="Logo"
      max-width="50"
      class="mr-4"
    />
    <v-spacer></v-spacer>
    <v-btn @click="logout" color="secondary">
      ログアウト
    </v-btn>
  </v-app-bar>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const logout = () => {
  // ログアウト処理: トークンを削除など
  localStorage.removeItem('authToken')
  // ログインページにリダイレクト
  router.push('/login')
}
</script>