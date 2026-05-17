<template>
  <div>
    <button
      class="vote-btn"
      :class="{ rated: myRating, disabled: !auth.user }"
      @click="openModal"
      :disabled="!auth.user"
    >
      <span class="vote-icon">★</span>
      <span>
        {{ !auth.user
          ? 'Войдите, чтобы оценить'
          : myRating
            ? 'Ваша оценка: ' + myScore
            : 'Оценить практику' }}
      </span>
    </button>

    <!-- Rating Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-box">
          <div class="modal-header">
            <h3 class="modal-title">{{ myRating ? 'Обновить оценку' : 'Оценить практику' }}</h3>
            <button class="modal-close" @click="closeModal">✕</button>
          </div>

          <div class="criteria-list">
            <div v-for="c in criteria" :key="c.key" class="criteria-row">
              <span class="criteria-label">{{ c.label }}</span>
              <div class="stars">
                <button
                  v-for="n in 5"
                  :key="n"
                  class="star-btn"
                  :class="{ active: form[c.key] >= n, hover: hovered[c.key] >= n }"
                  @click="form[c.key] = n"
                  @mouseenter="hovered[c.key] = n"
                  @mouseleave="hovered[c.key] = 0"
                >★</button>
              </div>
              <span class="score-val">{{ form[c.key] }}/5</span>
            </div>
          </div>

          <div class="comment-wrap">
            <label class="comment-label">Комментарий <span class="optional">(необязательно)</span></label>
            <textarea v-model="form.comment" class="textarea" rows="3" placeholder="Ваш опыт с этой практикой..."></textarea>
          </div>

          <p v-if="modalError" class="modal-error">{{ modalError }}</p>

          <div class="modal-actions">
            <button class="btn-cancel" @click="closeModal">Отмена</button>
            <button class="btn-submit" @click="submit" :disabled="loading || !isValid">
              {{ loading ? 'Сохранение...' : 'Сохранить оценку' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { api } from '../api'

const props = defineProps({
  practiceId: { type: Number, required: true },
})

const emit = defineEmits(['rated'])

const auth = useAuthStore()
const showModal = ref(false)
const loading = ref(false)
const modalError = ref(null)
const myRating = ref(null)

const form = ref({ score_utility: 3, score_aesthetics: 3, score_ecology: 3, comment: '' })
const hovered = ref({ score_utility: 0, score_aesthetics: 0, score_ecology: 0 })

const criteria = [
  { key: 'score_utility', label: 'Польза' },
  { key: 'score_aesthetics', label: 'Эстетика' },
  { key: 'score_ecology', label: 'Экология' },
]

const isValid = computed(() =>
  form.value.score_utility >= 1 && form.value.score_aesthetics >= 1 && form.value.score_ecology >= 1
)

const myScore = computed(() => {
  if (!myRating.value) return '—'
  const raw = 0.4 * myRating.value.score_utility + 0.35 * myRating.value.score_aesthetics + 0.25 * myRating.value.score_ecology
  return raw.toFixed(1)
})

onMounted(async () => {
  if (!auth.user) return
  try {
    const data = await api.getMyRating(props.practiceId)
    if (data) {
      myRating.value = data
      form.value.score_utility = data.score_utility
      form.value.score_aesthetics = data.score_aesthetics
      form.value.score_ecology = data.score_ecology
      form.value.comment = data.comment || ''
    }
  } catch {
    // no existing rating
  }
})

function openModal() {
  if (!auth.user) return
  showModal.value = true
  modalError.value = null
}

function closeModal() {
  showModal.value = false
}

async function submit() {
  loading.value = true
  modalError.value = null
  try {
    const data = await api.rate(props.practiceId, {
      score_utility: form.value.score_utility,
      score_aesthetics: form.value.score_aesthetics,
      score_ecology: form.value.score_ecology,
      comment: form.value.comment,
    })
    myRating.value = data
    emit('rated', data)
    closeModal()
  } catch (e) {
    modalError.value = e.data?.error || 'Ошибка при сохранении оценки'
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
.vote-btn.rated { background: #F3F4F6; color: var(--c-text); }
.vote-btn.rated:hover:not(:disabled) { background: #E5E7EB; }
.vote-btn.disabled { background: #F3F4F6; color: var(--c-muted); }
.vote-btn:disabled { cursor: not-allowed; }
.vote-icon { font-size: 16px; }

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 16px;
}
.modal-box {
  background: white;
  border-radius: 16px;
  padding: 28px;
  width: 100%;
  max-width: 440px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.2);
}
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
.modal-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--c-text);
}
.modal-close {
  background: none;
  border: none;
  font-size: 18px;
  color: var(--c-muted);
  cursor: pointer;
  padding: 4px;
  line-height: 1;
}

.criteria-list { display: flex; flex-direction: column; gap: 16px; margin-bottom: 20px; }
.criteria-row { display: flex; align-items: center; gap: 12px; }
.criteria-label { width: 80px; font-size: 14px; font-weight: 500; color: var(--c-text); flex-shrink: 0; }
.stars { display: flex; gap: 4px; flex: 1; }
.star-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: #D1D5DB;
  cursor: pointer;
  padding: 2px;
  line-height: 1;
  transition: color 0.1s, transform 0.1s;
}
.star-btn.active { color: #F59E0B; }
.star-btn.hover { color: #FCD34D; transform: scale(1.1); }
.score-val { font-size: 13px; color: var(--c-muted); width: 28px; text-align: right; flex-shrink: 0; }

.comment-wrap { margin-bottom: 20px; }
.comment-label { font-size: 14px; font-weight: 500; color: var(--c-text); display: block; margin-bottom: 8px; }
.optional { font-weight: 400; color: var(--c-muted); }
.textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1.5px solid var(--c-border);
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  resize: vertical;
  color: var(--c-text);
  box-sizing: border-box;
}
.textarea:focus { outline: none; border-color: var(--c-green); }

.modal-error { color: var(--c-red-dark, #DC2626); font-size: 13px; margin-bottom: 12px; }

.modal-actions { display: flex; gap: 12px; justify-content: flex-end; }
.btn-cancel {
  padding: 10px 20px;
  border: 1.5px solid var(--c-border);
  border-radius: 8px;
  background: white;
  color: var(--c-text);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
}
.btn-submit {
  padding: 10px 24px;
  background: var(--c-green);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  font-family: inherit;
  transition: background 0.15s;
}
.btn-submit:hover:not(:disabled) { background: var(--c-green-dark); }
.btn-submit:disabled { opacity: 0.55; cursor: not-allowed; }
</style>
