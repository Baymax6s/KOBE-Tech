import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { api, AUTH_TOKEN_STORAGE_KEY } from '@/api/client'
import { useAuthNotificationStore } from './authNotification'

const USER_ID_STORAGE_KEY = 'auth_user_id'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(AUTH_TOKEN_STORAGE_KEY))
  const userId = ref<number | null>(
    Number(localStorage.getItem(USER_ID_STORAGE_KEY)) || null,
  )

  const isAuthenticated = computed(() => token.value !== null)

  const persistUserId = (id: number | null) => {
    userId.value = id
    if (id !== null) {
      localStorage.setItem(USER_ID_STORAGE_KEY, String(id))
    } else {
      localStorage.removeItem(USER_ID_STORAGE_KEY)
    }
  }

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, newToken)
  }

  const clearToken = () => {
    token.value = null
    persistUserId(null)
    localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY)

    // 通知フラグもリセットして、再ログイン時や未ログイン時の誤表示を防ぐ
    const authNotification = useAuthNotificationStore()
    authNotification.consumeLoggedIn()
  }

  let fetchingUser = false

  const fetchUser = async () => {
    if (!token.value) return null
    if (fetchingUser) return userId.value
    fetchingUser = true
    try {
      const res = await api.api.authMeList()
      const id = res.data.id!
      persistUserId(id)
      return id
    } catch {
      clearToken()
      return null
    } finally {
      fetchingUser = false
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
    const { data } = await api.api.authMeList(options)
    persistUserId(data.id ?? null)
    return data
  }

  return {
    userId,
    isAuthenticated,
    login,
    fetchMe,
    clearToken,
    fetchUser,
  }
})
