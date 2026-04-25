import { Api } from './generated/apiSchema'

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const AUTH_TOKEN_STORAGE_KEY = 'auth_token'

export const api = new Api({ baseURL })

api.instance.interceptors.request.use((config) => {
  const token = localStorage.getItem(AUTH_TOKEN_STORAGE_KEY)
  if (token) {
    config.headers.set('Authorization', `Bearer ${token}`)
  }
  return config
})
