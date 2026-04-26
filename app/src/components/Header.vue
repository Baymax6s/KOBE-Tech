<template>
  <v-app-bar app color="primary" dark>
    <v-img src="/kdtech-icon.png" alt="Logo" max-width="200" class="mr-4" />
    <v-spacer></v-spacer>
    <v-btn v-if="showLogout" @click="logout" color="white"> ログアウト </v-btn>
  </v-app-bar>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const emit = defineEmits<{
  (e: 'onLogout'): void
}>()

const route = useRoute()
const showLogout = computed(() => route.path !== '/login')

const logout = () => {
  emit('onLogout')
}
</script>
