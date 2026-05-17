<template>
  <div class="comments">
    <h3>Комментарии ({{ comments.length }})</h3>

    <div v-if="auth.user" class="comment-form">
      <textarea v-model="text" placeholder="Написать комментарий..." rows="3"></textarea>
      <button class="btn btn-primary" @click="submit" :disabled="!text.trim() || loading">
        Отправить
      </button>
    </div>
    <p v-else class="auth-hint"><router-link to="/auth/login">Войдите</router-link>, чтобы оставить комментарий</p>

    <div v-if="comments.length === 0" class="no-comments">Комментариев пока нет</div>
    <div v-for="c in comments" :key="c.id" class="comment">
      <div class="comment-meta">Пользователь #{{ c.user_id }} · {{ formatDate(c.created_at) }}</div>
      <p>{{ c.text }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const props = defineProps({ practiceId: { type: Number, required: true } })

const auth = useAuthStore()
const comments = ref([])
const text = ref('')
const loading = ref(false)

onMounted(async () => {
  try {
    comments.value = await api.getComments(props.practiceId)
  } catch {}
})

async function submit() {
  if (!text.value.trim()) return
  loading.value = true
  try {
    const c = await api.addComment(props.practiceId, text.value)
    comments.value.push(c)
    text.value = ''
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function formatDate(d) {
  return new Date(d).toLocaleString('ru-RU')
}
</script>

<style scoped>
.comments { margin-top: 2rem; }
.comments h3 { margin-bottom: 1rem; }
.comment-form { margin-bottom: 1.5rem; }
.comment-form textarea { width: 100%; padding: 0.5rem; border: 1px solid #d1d5db; border-radius: 4px; margin-bottom: 0.5rem; }
.comment { padding: 0.75rem 0; border-bottom: 1px solid #f3f4f6; }
.comment-meta { font-size: 0.8rem; color: #9ca3af; margin-bottom: 0.25rem; }
.no-comments { color: #9ca3af; font-style: italic; padding: 1rem 0; }
.auth-hint { color: #6b7280; margin-bottom: 1rem; }
.auth-hint a { color: #2563eb; }
</style>
