<template>
  <div class="practice-card card">
    <div class="practice-header">
      <h3><router-link :to="`/practices/${practice.id}`">{{ practice.title }}</router-link></h3>
      <span :class="`badge badge-${practice.status}`">{{ statusLabel }}</span>
    </div>
    <p class="practice-meta">{{ practice.city }} · {{ practice.category }}</p>
    <p class="practice-desc">{{ truncated }}</p>
    <div class="practice-footer">
      <span class="vote-count">👍 {{ practice.vote_count }}</span>
      <span class="date">{{ formattedDate }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  practice: { type: Object, required: true }
})

const statusLabels = { pending: 'На модерации', approved: 'Одобрено', rejected: 'Отклонено' }
const statusLabel = computed(() => statusLabels[props.practice.status] || props.practice.status)

const truncated = computed(() => {
  const d = props.practice.description || ''
  return d.length > 150 ? d.slice(0, 150) + '...' : d
})

const formattedDate = computed(() =>
  new Date(props.practice.created_at).toLocaleDateString('ru-RU')
)
</script>

<style scoped>
.practice-card { margin-bottom: 1rem; }
.practice-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 0.25rem; }
.practice-header h3 { font-size: 1.05rem; }
.practice-header a { color: #2563eb; text-decoration: none; }
.practice-header a:hover { text-decoration: underline; }
.practice-meta { font-size: 0.85rem; color: #6b7280; margin-bottom: 0.5rem; }
.practice-desc { color: #555; font-size: 0.95rem; line-height: 1.5; }
.practice-footer { display: flex; justify-content: space-between; margin-top: 0.75rem; font-size: 0.85rem; color: #6b7280; }
</style>
