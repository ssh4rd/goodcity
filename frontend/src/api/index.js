import { $fetch } from 'ofetch'

const BASE = typeof window === 'undefined'
  ? (process.env.API_URL || 'http://localhost:8080')
  : ''

function getToken() {
  if (typeof localStorage !== 'undefined') {
    return localStorage.getItem('token')
  }
  return null
}

export const api = {
  fetch(path, options = {}) {
    const token = getToken()
    return $fetch(`${BASE}${path}`, {
      ...options,
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options.headers,
      },
    })
  },

  // Auth
  register: (email, password) =>
    api.fetch('/api/auth/register', { method: 'POST', body: { email, password } }),
  login: (email, password) =>
    api.fetch('/api/auth/login', { method: 'POST', body: { email, password } }),

  // Practices
  getPractices: (params = {}) => {
    const q = new URLSearchParams(
      Object.fromEntries(Object.entries(params).filter(([, v]) => v))
    ).toString()
    return api.fetch(`/api/practices${q ? '?' + q : ''}`)
  },
  getPractice: (id) => api.fetch(`/api/practices/${id}`),
  createPractice: (data) =>
    api.fetch('/api/practices', { method: 'POST', body: data }),
  updatePractice: (id, data) =>
    api.fetch(`/api/practices/${id}`, { method: 'PUT', body: data }),
  deletePractice: (id) =>
    api.fetch(`/api/practices/${id}`, { method: 'DELETE' }),

  // Votes
  vote: (id) => api.fetch(`/api/practices/${id}/vote`, { method: 'POST' }),

  // Comments
  getComments: (id) => api.fetch(`/api/practices/${id}/comments`),
  addComment: (id, text) =>
    api.fetch(`/api/practices/${id}/comments`, { method: 'POST', body: { text } }),

  // Moderation
  getPending: () => api.fetch('/api/moderation/pending'),
  approve: (id) => api.fetch(`/api/moderation/${id}/approve`, { method: 'POST' }),
  reject: (id) => api.fetch(`/api/moderation/${id}/reject`, { method: 'POST' }),
}
