import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useAuthNotificationStore = defineStore('authNotification', () => {
  const loggedIn = ref(false)

  const markLoggedIn = () => {
    loggedIn.value = true
  }

  const consumeLoggedIn = () => {
    if (!loggedIn.value) return false
    loggedIn.value = false
    return true
  }

  return { loggedIn, markLoggedIn, consumeLoggedIn }
})
