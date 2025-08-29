// composables/useApi.js
import { $fetch } from 'ofetch'

// Estructura de error estándar de tu backend:
// { code: number, message: string, details?: any }
class ApiError extends Error {
  constructor(payload, raw) {
    super(payload?.message || 'Error')
    this.name = 'ApiError'
    this.code = payload?.code || 500
    this.details = payload?.details
    this.raw = raw
  }
}

export const useApi = () => {
  const { public: { apiBase } } = useRuntimeConfig()

  const fullUrl = (path) => (path.startsWith('/') ? apiBase + path : `${apiBase}/${path}`)

  const normalizeError = (e) => {
    const payload = (e && e.data) || e || {}
    const code = typeof payload.code === 'number' ? payload.code : e?.status || 500
    const message = payload.message || e?.message || 'Error de red'
    throw new ApiError({ code, message, details: payload.details }, e)
  }

  const get = async (path, opts = {}) => {
    try {
      return await $fetch(fullUrl(path), { method: 'GET', ...opts })
    } catch (e) { normalizeError(e) }
  }

  const post = async (path, body, opts = {}) => {
    try {
      return await $fetch(fullUrl(path), { method: 'POST', body, ...opts })
    } catch (e) { normalizeError(e) }
  }

  const put = async (path, body, opts = {}) => {
    try {
      return await $fetch(fullUrl(path), { method: 'PUT', body, ...opts })
    } catch (e) { normalizeError(e) }
  }

  const del = async (path, opts = {}) => {
    try {
      return await $fetch(fullUrl(path), { method: 'DELETE', ...opts })
    } catch (e) { normalizeError(e) }
  }

  return { get, post, put, del }
}
