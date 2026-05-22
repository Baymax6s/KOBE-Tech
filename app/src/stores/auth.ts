import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { api, AUTH_TOKEN_STORAGE_KEY } from '@/api/client'
import type { ServerMeResponse } from '@/api/generated/apiSchema'
import { useAuthNotificationStore } from './authNotification'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(AUTH_TOKEN_STORAGE_KEY))
  const user = ref<ServerMeResponse | null>(null)

  const isAuthenticated = computed(() => token.value !== null)

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, newToken)
  }

  const clearToken = () => {
    token.value = null
    user.value = null
    localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY)

    // 通知フラグもリセットして、再ログイン時や未ログイン時の誤表示を防ぐ
    const authNotification = useAuthNotificationStore()
    authNotification.consumeLoggedIn()
  }

  const login = async (name: string, password: string) => {
    const { data } = await api.api.authLoginCreate(
      { name, password },
      { skipGlobalErrorHandler: true },
    )
    if (!data.token) {
      throw new Error('トークンが返却されませんでした')
    }
    setToken(data.token)
  }

  const fetchMe = async (options?: { skipGlobalErrorHandler?: boolean }) => {
    try {
      const { data } = await api.api.authMeList(options)
      user.value = data
    } catch (e) {
      user.value = null
      throw e
    }
  }

  return {
    token,
    user,
    isAuthenticated,
    login,
    fetchMe,
    clearToken,
  }
})
