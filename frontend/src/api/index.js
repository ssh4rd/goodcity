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
  register: (email, password, socialRole = 'resident', name = '', city = '', district = '') =>
    api.fetch('/api/auth/register', { method: 'POST', body: { email, password, social_role: socialRole, name, city, district } }),
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

  // Ratings
  rate: (id, data) =>
    api.fetch(`/api/practices/${id}/rate`, { method: 'POST', body: data }),
  getMyRating: (id) => api.fetch(`/api/practices/${id}/rate`),
  getRatingStats: (id) => api.fetch(`/api/practices/${id}/rating`),

  // Comments
  getComments: (id) => api.fetch(`/api/practices/${id}/comments`),
  addComment: (id, text) =>
    api.fetch(`/api/practices/${id}/comments`, { method: 'POST', body: { text } }),

  // Moderation
  getPending: () => api.fetch('/api/moderation/pending'),
  approve: (id) => api.fetch(`/api/moderation/${id}/approve`, { method: 'POST' }),
  reject: (id) => api.fetch(`/api/moderation/${id}/reject`, { method: 'POST' }),

  // User
  getMe: () => api.fetch('/api/user/me'),
  updateMe: (data) => api.fetch('/api/user/me', { method: 'PUT', body: data }),

  // Role verification
  submitVerification: (formData) => api.fetch('/api/user/verification', { method: 'POST', body: formData }),
  getPendingVerifications: () => api.fetch('/api/moderation/verifications'),
  approveVerification: (id) => api.fetch(`/api/moderation/verifications/${id}/approve`, { method: 'POST' }),
  rejectVerification: (id) => api.fetch(`/api/moderation/verifications/${id}/reject`, { method: 'POST' }),
}
