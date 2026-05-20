import axios, { AxiosHeaders } from 'axios'
import { Api } from './generated/apiSchema'

declare module 'axios' {
  export interface AxiosRequestConfig {
    // ページ全体のエラー画面遷移を抑止したい呼び出しに付ける
    // （例：記事詳細内のコメントなど、サブ機能の取得失敗で記事ごと吹き飛ばしたくないケース）
    skipGlobalErrorHandler?: boolean
  }
}

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const AUTH_TOKEN_STORAGE_KEY = 'auth_token'

export const api = new Api({ baseURL })

// 配列クエリは `?tag=a&tag=b`（OpenAPI の collectionFormat: multi）形式で送る。
// axios のデフォルトは `tag[]=a&tag[]=b` で、Gin の c.QueryArray("tag") が拾えないため。
api.instance.defaults.paramsSerializer = { indexes: null }

type ApiErrorHandler = (status: number | null) => void

let apiErrorHandler: ApiErrorHandler | null = null

export const setApiErrorHandler = (handler: ApiErrorHandler) => {
  apiErrorHandler = handler
}

api.instance.interceptors.request.use((config) => {
  const token = localStorage.getItem(AUTH_TOKEN_STORAGE_KEY)
  if (token) {
    config.headers = AxiosHeaders.from(config.headers)
    config.headers.set('Authorization', `Bearer ${token}`)
  }
  return config
})

api.instance.interceptors.response.use(
  (response) => response,
  (error: unknown) => {
    if (axios.isAxiosError(error) && error.code !== 'ERR_CANCELED') {
      const skip = error.config?.skipGlobalErrorHandler === true

      if (!skip) {
        const status = error.response?.status ?? null

        if (status === 404 || (status !== null && status >= 500)) {
          apiErrorHandler?.(status)
        }

        if (status === null && error.request) {
          apiErrorHandler?.(status)
        }
      }
    }

    return Promise.reject(error)
  },
)
