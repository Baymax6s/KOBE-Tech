import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { api, AUTH_TOKEN_STORAGE_KEY } from '@/api/client'
import type { ServerMeResponse } from '@/api/generated/apiSchema'

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
  }

  const login = async (name: string, password: string) => {
    const { data } = await api.api.authLoginCreate({ name, password })
    if (!data.token) {
      throw new Error('トークンが返却されませんでした')
    }
    setToken(data.token)
  }

  const fetchMe = async () => {
    const { data } = await api.api.authMeList()
    user.value = data
    return data
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
