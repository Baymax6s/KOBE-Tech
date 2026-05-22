import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { api, AUTH_TOKEN_STORAGE_KEY } from '@/api/client'
import type { ServerMeResponse } from '@/api/generated/apiSchema'
import { useAuthNotificationStore } from './authNotification'

const USER_STORAGE_KEY = 'auth_user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(AUTH_TOKEN_STORAGE_KEY))
  const user = ref<{ id: number; name: string } | null>(
    JSON.parse(localStorage.getItem(USER_STORAGE_KEY) ?? 'null'),
  )

  const isAuthenticated = computed(() => token.value !== null)
  const userId = computed(() => user.value?.id ?? null)

  const persistUser = (u: { id: number; name: string } | null) => {
    user.value = u
    if (u) {
      localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(u))
    } else {
      localStorage.removeItem(USER_STORAGE_KEY)
    }
  }

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, newToken)
  }

  const clearToken = () => {
    token.value = null
    persistUser(null)
    localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY)

    // 通知フラグもリセットして、再ログイン時や未ログイン時の誤表示を防ぐ
    const authNotification = useAuthNotificationStore()
    authNotification.consumeLoggedIn()
  }

  const fetchUser = async () => {
    if (!token.value) return null
    try {
      const res = await api.api.authMeList()
      const u = { id: res.data.id!, name: res.data.name! }
      persistUser(u)
      return u
    } catch {
      clearToken()
      return null
    }
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
    await fetchUser()
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
    userId,
    login,
    fetchMe,
    clearToken,
    fetchUser,
  }
})
