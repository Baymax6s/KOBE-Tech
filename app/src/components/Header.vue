<template>
  <v-app-bar app color="primary" dark>
    <v-img
      src="/kdtech-icon.png"
      alt="Logo"
      max-width="200"
      class="mr-4"
    />
    <v-spacer></v-spacer>
    <v-btn @click="logout" color="white">
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