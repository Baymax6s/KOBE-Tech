import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useArticleNotificationStore = defineStore('articleNotification', () => {
  const created = ref(false)

  const markCreated = () => {
    created.value = true
  }

  const consumeCreated = () => {
    if (!created.value) return false
    created.value = false
    return true
  }

  return { created, markCreated, consumeCreated }
})
