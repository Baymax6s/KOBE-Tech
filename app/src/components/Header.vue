<template>
  <v-app-bar app color="primary" dark>
    <v-img src="/kdtech-icon.png" alt="Logo" max-width="200" class="mr-4" />
    <v-spacer></v-spacer>
    <v-btn v-if="auth.isAuthenticated" @click="logout" color="white">
      ログアウト
    </v-btn>
    <v-btn v-else @click="goToLogin" color="white" variant="outlined">
      ログイン
    </v-btn>
  </v-app-bar>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const logout = () => {
  auth.clearToken()
  router.push('/')
}

const goToLogin = () => {
  router.push('/login')
}
</script>
