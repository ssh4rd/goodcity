<template>
  <AdminLayout>
    <div v-if="loading" class="loading-state">Загрузка...</div>
    <div v-else-if="error" class="error-state">{{ error }}</div>
    <template v-else-if="practice">
      <div class="detail-header">
        <router-link to="/admin" class="back-link">← Назад к очереди</router-link>
        <div class="detail-title-row">
          <h1 class="page-title">{{ practice.title }}</h1>
          <span :class="`badge badge-${practice.status}`">{{ statusLabel }}</span>
        </div>
        <p class="detail-meta">{{ practice.city }} · {{ practice.category }} · Автор #{{ practice.author_id }}</p>
      </div>

      <div class="detail-layout">
        <div class="detail-main card">
          <div class="detail-image">
            <svg width="48" height="38" viewBox="0 0 48 38" fill="none">
              <rect x="2" y="2" width="44" height="34" rx="3" fill="#D1D5DB"/>
              <text x="24" y="22" text-anchor="middle" fill="#9CA3AF" font-size="10" font-family="Inter,sans-serif">Фото практики</text>
            </svg>
          </div>
          <div class="detail-content">
            <h3 class="section-label">Описание</h3>
            <p class="detail-desc">{{ practice.description }}</p>
          </div>
        </div>

        <div class="detail-side">
          <div class="card info-panel">
            <div class="info-row">
              <span class="info-key">Статус</span>
              <span :class="`badge badge-${practice.status}`">{{ statusLabel }}</span>
            </div>
            <div class="info-row">
              <span class="info-key">Категория</span>
              <span class="info-val">{{ practice.category }}</span>
            </div>
            <div class="info-row">
              <span class="info-key">Город</span>
              <span class="info-val">{{ practice.city }}</span>
            </div>
            <div class="info-row">
              <span class="info-key">Дата добавления</span>
              <span class="info-val">{{ formattedDate }}</span>
            </div>
            <div class="info-row">
              <span class="info-key">Голоса</span>
              <span class="info-val">{{ practice.vote_count }}</span>
            </div>
          </div>

          <div v-if="practice.status === 'pending'" class="action-panel card">
            <h3 class="action-title">Принять решение</h3>
            <p class="action-sub">Проверьте практику и одобрите или отклоните её публикацию</p>
            <div class="action-buttons">
              <button class="action-approve-lg" @click="approve" :disabled="actioning">
                ✓ Одобрить практику
              </button>
              <button class="action-reject-lg" @click="reject" :disabled="actioning">
                ✕ Отклонить
              </button>
            </div>
          </div>
          <div v-else class="decided-panel card">
            <p>Решение уже принято: <strong>{{ statusLabel }}</strong></p>
          </div>
        </div>
      </div>
    </template>
  </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../../../api'
import AdminLayout from '../../../components/AdminLayout.vue'

const route = useRoute()
const practice = ref(null)
const loading = ref(true)
const error = ref(null)
const actioning = ref(false)

const statusLabels = { pending: 'На модерации', approved: 'Одобрено', rejected: 'Отклонено' }
const statusLabel = computed(() => statusLabels[practice.value?.status] || '')
const formattedDate = computed(() =>
  practice.value ? new Date(practice.value.created_at).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' }) : ''
)

async function approve() {
  actioning.value = true
  try {
    await api.approve(practice.value.id)
    practice.value.status = 'approved'
  } catch {} finally { actioning.value = false }
}
async function reject() {
  actioning.value = true
  try {
    await api.reject(practice.value.id)
    practice.value.status = 'rejected'
  } catch {} finally { actioning.value = false }
}

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
.loading-state, .error-state { padding: 40px 0; text-align: center; color: var(--c-muted); }

.detail-header { margin-bottom: 28px; }
.back-link { font-size: 14px; color: var(--c-muted); display: inline-block; margin-bottom: 16px; transition: color 0.15s; }
.back-link:hover { color: var(--c-text); }
.detail-title-row { display: flex; align-items: center; gap: 14px; margin-bottom: 8px; }
.page-title { font-size: 26px; font-weight: 700; }
.detail-meta { font-size: 14px; color: var(--c-muted); }

.detail-layout { display: grid; grid-template-columns: 1fr 280px; gap: 20px; align-items: start; }

.detail-main { padding: 24px; display: flex; flex-direction: column; gap: 20px; }
.detail-image {
  background: #F3F4F6;
  border-radius: 10px;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.section-label { font-size: 14px; font-weight: 600; color: var(--c-text); margin-bottom: 8px; }
.detail-desc { font-size: 14px; line-height: 1.7; color: var(--c-text); }

.detail-side { display: flex; flex-direction: column; gap: 16px; }

.info-panel { padding: 20px; display: flex; flex-direction: column; gap: 12px; }
.info-row { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.info-key { font-size: 13px; color: var(--c-muted); }
.info-val { font-size: 13px; font-weight: 500; }

.action-panel { padding: 20px; }
.action-title { font-size: 15px; font-weight: 600; margin-bottom: 6px; }
.action-sub { font-size: 13px; color: var(--c-muted); margin-bottom: 16px; line-height: 1.5; }
.action-buttons { display: flex; flex-direction: column; gap: 8px; }
.action-approve-lg {
  width: 100%;
  padding: 11px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  font-family: inherit;
  background: #D1FAE5;
  color: #065F46;
  border: 1px solid #A7F3D0;
  cursor: pointer;
  transition: background 0.15s;
}
.action-approve-lg:hover:not(:disabled) { background: #A7F3D0; }
.action-approve-lg:disabled { opacity: 0.5; }
.action-reject-lg {
  width: 100%;
  padding: 11px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  font-family: inherit;
  background: #FEE2E2;
  color: #991B1B;
  border: 1px solid #FECACA;
  cursor: pointer;
  transition: background 0.15s;
}
.action-reject-lg:hover:not(:disabled) { background: #FECACA; }
.action-reject-lg:disabled { opacity: 0.5; }

.decided-panel { padding: 16px; font-size: 14px; color: var(--c-muted); }

@media (max-width: 768px) {
  .page-title { font-size: 20px; }
  .detail-layout { grid-template-columns: 1fr; }
  .detail-side { order: -1; }
}
</style>
