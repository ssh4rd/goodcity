<template>
  <button class="btn vote-btn" :class="{ voted }" @click="handleVote" :disabled="!auth.user || loading">
    👍 {{ count }}
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
    if (data.voted) {
      count.value++
      voted.value = true
    } else {
      count.value--
      voted.value = false
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.vote-btn { background: #e5e7eb; color: #374151; font-size: 1rem; }
.vote-btn.voted { background: #dbeafe; color: #1d4ed8; }
.vote-btn:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
