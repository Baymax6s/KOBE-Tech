import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { api, AUTH_TOKEN_STORAGE_KEY } from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(AUTH_TOKEN_STORAGE_KEY))

  const isAuthenticated = computed(() => token.value !== null)

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, newToken)
  }

  const clearToken = () => {
    token.value = null
    localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY)
  }

  const login = async (name: string, password: string) => {
    const { data } = await api.api.authLoginCreate({ name, password })
    if (!data.token) {
      throw new Error('トークンが返却されませんでした')
    }
    setToken(data.token)
  }

  return {
    token,
    isAuthenticated,
    login,
    clearToken,
  }
})
