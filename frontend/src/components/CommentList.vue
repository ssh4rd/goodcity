<template>
  <div class="comments-section">
    <h2 class="comments-title">Комментарии ({{ comments.length }})</h2>

    <!-- Form -->
    <div class="comment-form-wrap card">
      <template v-if="auth.user">
        <div class="comment-form">
          <div class="comment-avatar me">{{ userInitial }}</div>
          <input
            v-model="text"
            class="comment-input"
            placeholder="Написать комментарий..."
            @keydown.enter.exact.prevent="submit"
          />
          <button class="btn btn-primary btn-sm" @click="submit" :disabled="!text.trim() || loading">
            Отправить
          </button>
        </div>
      </template>
      <p v-else class="auth-hint">
        <router-link to="/auth/login" class="auth-link">Войдите</router-link>, чтобы оставить комментарий
      </p>
    </div>

    <!-- List -->
    <div v-if="comments.length === 0 && !loading" class="no-comments">Комментариев пока нет. Будьте первым!</div>

    <div class="comment-list">
      <div v-for="c in comments" :key="c.id" class="comment-item">
        <div class="comment-avatar">{{ avatarLetter(c.user_id) }}</div>
        <div class="comment-body">
          <div class="comment-header">
            <span class="comment-author">Пользователь #{{ c.user_id }}</span>
            <span class="comment-time">{{ formatDate(c.created_at) }}</span>
          </div>
          <p class="comment-text">{{ c.text }}</p>
          <div class="comment-actions">
            <button class="comment-action">★★★★★</button>
            <button class="comment-action">Ответить</button>
            <button class="comment-action">Пожаловаться</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const props = defineProps({ practiceId: { type: Number, required: true } })

const auth = useAuthStore()
const comments = ref([])
const text = ref('')
const loading = ref(false)

const userInitial = computed(() => auth.user?.email?.[0]?.toUpperCase() || '?')
function avatarLetter(userId) {
  const letters = 'АБВГДЕЖЗИЙКЛМНОПРСТУФ'
  return letters[(userId || 0) % letters.length]
}

onMounted(async () => {
  try { comments.value = await api.getComments(props.practiceId) } catch {}
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
  return new Date(d).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' })
}
</script>

<style scoped>
.comments-section { margin-top: 40px; }
.comments-title { font-size: 22px; font-weight: 700; margin-bottom: 20px; }

.comment-form-wrap { padding: 16px 20px; margin-bottom: 24px; }
.comment-form { display: flex; align-items: center; gap: 12px; }
.comment-avatar {
  width: 36px; height: 36px;
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  font-size: 13px; font-weight: 700;
  flex-shrink: 0;
  color: white;
}
.comment-avatar.me { background: var(--c-green); }
.comment-input {
  flex: 1;
  padding: 10px 16px;
  border: 1.5px solid var(--c-border);
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  outline: none;
  background: #F8FAFB;
  transition: border-color 0.15s;
}
.comment-input:focus { border-color: var(--c-green); background: white; }
.comment-input::placeholder { color: #9CA3AF; }

.auth-hint { font-size: 14px; color: var(--c-muted); }
.auth-link { color: var(--c-green); font-weight: 500; }

.no-comments { text-align: center; padding: 32px 0; color: var(--c-muted); font-size: 15px; }

.comment-list { display: flex; flex-direction: column; gap: 0; }
.comment-item {
  display: flex;
  gap: 14px;
  padding: 20px 0;
  border-bottom: 1px solid var(--c-border);
}
.comment-item:last-child { border-bottom: none; }

.comment-avatar {
  width: 38px; height: 38px;
  border-radius: 50%;
  background: var(--c-blue);
  color: white;
  font-size: 14px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.comment-body { flex: 1; }
.comment-header { display: flex; align-items: center; gap: 12px; margin-bottom: 6px; }
.comment-author { font-size: 14px; font-weight: 600; }
.comment-time { font-size: 12px; color: var(--c-muted); }
.comment-text { font-size: 14px; color: var(--c-text); line-height: 1.6; margin-bottom: 8px; }
.comment-actions { display: flex; gap: 12px; }
.comment-action {
  background: none;
  border: none;
  font-size: 12px;
  color: var(--c-muted);
  cursor: pointer;
  padding: 0;
  font-family: inherit;
}
.comment-action:first-child { color: var(--c-amber); }
.comment-action:hover { color: var(--c-text); }

@media (max-width: 480px) {
  .comment-form { flex-wrap: wrap; }
  .comment-input { min-width: 0; }
  .comment-form .btn { width: 100%; }
}
</style>
