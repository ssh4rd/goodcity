import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(null)

  const isModerator = computed(() => user.value?.role === 'moderator')

  function init() {
    if (typeof localStorage === 'undefined') return
    const stored = localStorage.getItem('auth')
    if (stored) {
      try {
        const data = JSON.parse(stored)
        user.value = data.user
        token.value = data.token
      } catch {}
    }
  }

  async function register(email, password, socialRole, name, city, district) {
    const data = await api.register(email, password, socialRole, name, city, district)
    setAuth(data)
    return data
  }

  async function login(email, password) {
    const data = await api.login(email, password)
    setAuth(data)
    return data
  }

  function logout() {
    user.value = null
    token.value = null
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem('auth')
      localStorage.removeItem('token')
    }
  }

  function setAuth(data) {
    user.value = data.user
    token.value = data.token
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem('auth', JSON.stringify(data))
      localStorage.setItem('token', data.token)
    }
  }

  return { user, token, isModerator, init, register, login, logout }
})
