<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const goToProfile = () => {
  router.push('/settings/profile')
}

const goToPasswordSettings = () => {
  router.push('/settings/password')
}

const logout = () => {
  auth.clearToken()
  router.push('/')
}

const goToLogin = () => {
  router.push('/login')
}
</script>
<template>
  <v-app-bar app color="primary" dark>
    <v-img src="/kdtech-icon.png" alt="Logo" max-width="200" class="mr-4" />

    <v-spacer></v-spacer>

    <!-- 設定メニュー -->
    <v-menu offset-y>
      <template #activator="{ props }">
        <v-btn
          v-if="auth.isAuthenticated"
          icon="mdi-cog"
          aria-label="設定"
          v-bind="props"
        />
      </template>

      <v-list>
        <v-list-item @click="goToProfile">
          <template #prepend>
            <v-icon>mdi-account</v-icon>
          </template>

          <v-list-item-title> プロフィール </v-list-item-title>
        </v-list-item>

        <v-list-item @click="goToPasswordSettings">
          <template #prepend>
            <v-icon>mdi-lock</v-icon>
          </template>

          <v-list-item-title> パスワード変更 </v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>

    <!-- ログアウト -->
    <v-btn v-if="auth.isAuthenticated" @click="logout" color="white">
      ログアウト
    </v-btn>

    <!-- ログイン -->
    <v-btn v-else @click="goToLogin" color="white" variant="outlined">
      ログイン
    </v-btn>
  </v-app-bar>
</template>
