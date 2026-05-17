<template>
  <button class="vote-btn" :class="{ voted, disabled: !auth.user }" @click="handleVote" :disabled="!auth.user || loading">
    <span class="vote-icon">★</span>
    <span>{{ auth.user ? (voted ? 'Оценено' : 'Оценить практику') : 'Войдите, чтобы оценить' }}</span>
    <span v-if="count > 0" class="vote-count">{{ count }}</span>
  </button>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const props = defineProps({
  practiceId: { type: Number, required: true },
  initialCount: { type: Number, default: 0 },
  initialVoted: { type: Boolean, default: false },
})

const auth = useAuthStore()
const count = ref(props.initialCount)
const voted = ref(props.initialVoted)
const loading = ref(false)

async function handleVote() {
  if (!auth.user) return
  loading.value = true
  try {
    const data = await api.vote(props.practiceId)
    if (data.voted) { count.value++; voted.value = true }
    else { count.value--; voted.value = false }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.vote-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 22px;
  background: var(--c-green);
  color: white;
  font-size: 14px;
  font-weight: 600;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font-family: inherit;
  transition: background 0.15s;
}
.vote-btn:hover:not(:disabled) { background: var(--c-green-dark); }
.vote-btn.voted { background: #F3F4F6; color: var(--c-text); }
.vote-btn.voted:hover:not(:disabled) { background: #E5E7EB; }
.vote-btn.disabled { background: #F3F4F6; color: var(--c-muted); }
.vote-btn:disabled { cursor: not-allowed; }
.vote-icon { font-size: 16px; }
.vote-count {
  background: rgba(255,255,255,0.3);
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}
.vote-btn.voted .vote-count,
.vote-btn.disabled .vote-count { background: var(--c-border); }
</style>
