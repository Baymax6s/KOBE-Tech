<script setup lang="ts">
import { computed } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import Header from './components/Header.vue'

const route = useRoute()
const showHeader = computed(() => route.meta?.hideHeader !== true)
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
