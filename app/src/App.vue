<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import Header from './components/Header.vue'
import { useAuthStore } from './stores/auth'

const route = useRoute()
const showHeader = computed(() => route.meta?.hideHeader !== true)
const auth = useAuthStore()

// リロード時にトークンはあるがuserIdがnullのため、全ページ共通でユーザー情報を復元する
onMounted(() => {
  if (auth.isAuthenticated) {
    auth.fetchCurrentUser()
  }
})
</script>

<template>
  <v-app>
    <v-layout>
      <Header v-if="showHeader" />
      <v-main>
        <RouterView />
      </v-main>
    </v-layout>
  </v-app>
</template>
