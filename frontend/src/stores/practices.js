import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api'

export const usePracticesStore = defineStore('practices', () => {
  const practices = ref([])
  const total = ref(0)
  const loading = ref(false)
  const error = ref(null)

  async function fetchPractices(params = {}) {
    loading.value = true
    error.value = null
    try {
      const data = await api.getPractices(params)
      practices.value = data.data || []
      total.value = data.total || 0
    } catch (e) {
      error.value = e.message || 'Ошибка загрузки'
    } finally {
      loading.value = false
    }
  }

  return { practices, total, loading, error, fetchPractices }
})
