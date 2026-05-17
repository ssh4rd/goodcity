<template>
  <div>
    <h1 style="margin-bottom:1.5rem">Модерация практик</h1>

    <div v-if="loading">Загрузка...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else>
      <div v-if="practices.length === 0" class="card" style="text-align:center;color:#6b7280">
        Нет практик на модерации
      </div>
      <div v-for="p in practices" :key="p.id" class="card moderation-item">
        <div class="item-header">
          <h3>
            <router-link :to="`/admin/practice/${p.id}`">{{ p.title }}</router-link>
          </h3>
          <span class="badge badge-pending">На модерации</span>
        </div>
        <p class="meta">{{ p.city }} · {{ p.category }}</p>
        <p class="desc">{{ truncate(p.description) }}</p>
        <div class="item-actions">
          <button class="btn btn-success" @click="approve(p.id)">Одобрить</button>
          <button class="btn btn-danger" @click="reject(p.id)">Отклонить</button>
          <router-link :to="`/admin/practice/${p.id}`" class="btn btn-secondary">Детали</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { api } from '../../api'

const practices = ref([])
const loading = ref(true)
const error = ref(null)

async function load() {
  try {
    practices.value = await api.getPending()
  } catch (e) {
    error.value = 'Ошибка загрузки'
  } finally {
    loading.value = false
  }
}

async function approve(id) {
  await api.approve(id)
  practices.value = practices.value.filter(p => p.id !== id)
}

async function reject(id) {
  await api.reject(id)
  practices.value = practices.value.filter(p => p.id !== id)
}

function truncate(text) {
  return text?.length > 200 ? text.slice(0, 200) + '...' : text
}

onMounted(load)
</script>

<style scoped>
.moderation-item { margin-bottom: 1rem; }
.item-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.25rem; }
.item-header h3 a { color: #2563eb; text-decoration: none; }
.meta { font-size: 0.85rem; color: #6b7280; margin-bottom: 0.5rem; }
.desc { color: #555; margin-bottom: 0.75rem; }
.item-actions { display: flex; gap: 0.5rem; }
</style>
