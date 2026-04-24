import { Api } from './generated/apiSchema'

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const api = new Api({ baseURL })
