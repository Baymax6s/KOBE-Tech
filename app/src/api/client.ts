import axios from 'axios'
import { Api } from './generated/apiSchema'

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const AUTH_TOKEN_STORAGE_KEY = 'auth_token'

export const api = new Api({ baseURL })

type ApiErrorHandler = (status: number | null) => void

let apiErrorHandler: ApiErrorHandler | null = null

export const setApiErrorHandler = (handler: ApiErrorHandler) => {
  apiErrorHandler = handler
}

api.instance.interceptors.request.use((config) => {
  const token = localStorage.getItem(AUTH_TOKEN_STORAGE_KEY)
  if (token) {
    config.headers.set('Authorization', `Bearer ${token}`)
  }
  return config
})

api.instance.interceptors.response.use(
  (response) => response,
  (error: unknown) => {
    if (axios.isAxiosError(error) && error.code !== 'ERR_CANCELED') {
      const status = error.response?.status ?? null

      if (status === 404 || (status !== null && status >= 500)) {
        apiErrorHandler?.(status)
      }

      if (status === null && error.request) {
        apiErrorHandler?.(status)
      }
    }

    return Promise.reject(error)
  },
)
