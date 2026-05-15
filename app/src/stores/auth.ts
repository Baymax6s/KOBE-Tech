import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import axios from 'axios'
import { api, AUTH_TOKEN_STORAGE_KEY } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(AUTH_TOKEN_STORAGE_KEY))
  const userId = ref<number | null>(null)

  const isAuthenticated = computed(() => token.value !== null)

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, newToken)
  }

  const clearToken = () => {
    token.value = null
    userId.value = null
    localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY)
  }

  const fetchCurrentUser = async () => {
    if (!token.value) return
    try {
      const { data } = await api.api.authMeList({
        skipGlobalErrorHandler: true,
      })
      userId.value = data.id ?? null
    } catch (err) {
      if (axios.isAxiosError(err) && err.response?.status === 401) {
        clearToken()
      }
    }
  }

  const login = async (name: string, password: string) => {
    const { data } = await api.api.authLoginCreate({ name, password })
    if (!data.token) {
      throw new Error('トークンが返却されませんでした')
    }
    setToken(data.token)
    await fetchCurrentUser()
  }

  return {
    token,
    userId,
    isAuthenticated,
    login,
    clearToken,
    fetchCurrentUser,
  }
})
