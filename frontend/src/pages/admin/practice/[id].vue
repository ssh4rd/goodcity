<template>
  <div v-if="loading">Загрузка...</div>
  <div v-else-if="error" class="error">{{ error }}</div>
  <div v-else-if="practice">
    <router-link to="/admin" class="back-link">← Назад к очереди</router-link>

    <div class="card" style="margin-top:1rem">
      <div style="display:flex;align-items:center;gap:1rem;margin-bottom:0.5rem">
        <h1>{{ practice.title }}</h1>
        <span :class="`badge badge-${practice.status}`">{{ statusLabel }}</span>
      </div>
      <p class="meta">{{ practice.city }} · {{ practice.category }} · #{{ practice.author_id }}</p>
      <div class="description">{{ practice.description }}</div>

      <div v-if="practice.status === 'pending'" class="actions">
        <button class="btn btn-success" @click="approve">Одобрить</button>
        <button class="btn btn-danger" @click="reject">Отклонить</button>
      </div>
      <p v-else style="margin-top:1rem;color:#6b7280">Решение уже принято</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../../../api'

const route = useRoute()
const router = useRouter()
const practice = ref(null)
const loading = ref(true)
const error = ref(null)

const statusLabels = { pending: 'На модерации', approved: 'Одобрено', rejected: 'Отклонено' }
const statusLabel = computed(() => statusLabels[practice.value?.status] || '')

onMounted(async () => {
  try {
    practice.value = await api.getPractice(route.params.id)
  } catch {
    error.value = 'Практика не найдена'
  } finally {
    loading.value = false
  }
})

async function approve() {
  await api.approve(practice.value.id)
  practice.value.status = 'approved'
}

async function reject() {
  await api.reject(practice.value.id)
  practice.value.status = 'rejected'
}
</script>

<style scoped>
.back-link { color: #2563eb; text-decoration: none; }
.meta { color: #6b7280; font-size: 0.9rem; margin-bottom: 1rem; }
.description { line-height: 1.7; white-space: pre-wrap; margin-bottom: 1.5rem; }
.actions { display: flex; gap: 0.75rem; }
</style>
