<template>
  <div v-if="loading">Загрузка...</div>
  <div v-else-if="error" class="error">{{ error }}</div>
  <div v-else-if="practice">
    <div class="card practice-detail">
      <div class="detail-header">
        <h1>{{ practice.title }}</h1>
        <span :class="`badge badge-${practice.status}`">{{ statusLabel }}</span>
      </div>
      <p class="meta">{{ practice.city }} · {{ practice.category }} · {{ formattedDate }}</p>
      <div class="description">{{ practice.description }}</div>
      <div class="actions">
        <VoteButton :practice-id="practice.id" :initial-count="practice.vote_count" />
      </div>
    </div>
    <CommentList :practice-id="practice.id" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../../api'
import VoteButton from '../../components/VoteButton.vue'
import CommentList from '../../components/CommentList.vue'

const route = useRoute()
const practice = ref(null)
const loading = ref(true)
const error = ref(null)

const statusLabels = { pending: 'На модерации', approved: 'Одобрено', rejected: 'Отклонено' }
const statusLabel = computed(() => statusLabels[practice.value?.status] || '')
const formattedDate = computed(() =>
  practice.value ? new Date(practice.value.created_at).toLocaleDateString('ru-RU') : ''
)

onMounted(async () => {
  try {
    practice.value = await api.getPractice(route.params.id)
  } catch {
    error.value = 'Практика не найдена'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.practice-detail { margin-bottom: 1.5rem; }
.detail-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.5rem; }
.detail-header h1 { font-size: 1.5rem; }
.meta { color: #6b7280; font-size: 0.9rem; margin-bottom: 1rem; }
.description { line-height: 1.7; white-space: pre-wrap; }
.actions { margin-top: 1.5rem; }
</style>
