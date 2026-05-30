<template>
  <AdminLayout :pending-count="practices.length + verifications.length">
    <div class="moderation-page">
      <div class="page-header">
        <div>
          <h1 class="page-title">{{ mode === 'practices' ? 'Модерация практик' : 'Верификация ролей' }}</h1>
          <p class="page-sub">{{ mode === 'practices' ? `${practices.length} заявок ожидают проверки` : `${verifications.length} запросов на верификацию` }}</p>
        </div>
        <div class="mode-switcher">
          <button class="mode-btn" :class="{ active: mode === 'practices' }" @click="mode = 'practices'">Практики</button>
          <button class="mode-btn" :class="{ active: mode === 'verifications' }" @click="mode = 'verifications'">
            Верификация <span v-if="verifications.length" class="badge-num">{{ verifications.length }}</span>
          </button>
        </div>
      </div>

      <!-- Verifications mode -->
      <div v-if="mode === 'verifications'">
        <div v-if="vLoading" class="loading-state">Загрузка...</div>
        <div v-else-if="vError" class="error-state">{{ vError }}</div>
        <div v-else-if="verifications.length === 0" class="empty-table card" style="padding:40px;text-align:center;color:var(--c-muted)">
          Нет запросов на верификацию
        </div>
        <div v-else class="table-wrap card">
          <table class="mod-table">
            <thead>
              <tr>
                <th>Пользователь</th>
                <th>Роль</th>
                <th>Документ</th>
                <th>Дата</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="u in verifications" :key="u.id" class="table-row">
                <td>{{ u.email }}</td>
                <td><span class="role-badge">{{ roleLabel(u.social_role) }}</span></td>
                <td>
                  <a v-if="u.verification_doc_url" :href="u.verification_doc_url" target="_blank" class="doc-link">Открыть документ →</a>
                </td>
                <td class="col-date">{{ formatDate(u.created_at) }}</td>
                <td class="col-actions" style="min-width:200px">
                  <button class="action-approve" @click="approveVerif(u.id)" :disabled="vActionLoading === u.id">✓ Подтвердить</button>
                  <button class="action-reject" @click="rejectVerif(u.id)" :disabled="vActionLoading === u.id">✕ Отклонить</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Toolbar -->
      <div v-if="mode === 'practices'" class="toolbar">
        <div class="search-wrap">
          <span class="search-icon">
            <svg width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
            </svg>
          </span>
          <input v-model="search" class="search-input" placeholder="Поиск по названию или автору..." />
        </div>
        <button class="filter-btn">Категория ▾</button>
        <button class="filter-btn">Дата ▾</button>
      </div>

      <!-- Status tabs -->
      <div v-if="mode === 'practices'" class="status-tabs">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          class="status-tab"
          :class="{ active: activeTab === tab.value }"
          @click="activeTab = tab.value"
        >
          {{ tab.label }} ({{ tab.count }})
        </button>
      </div>

      <!-- Loading -->
      <div v-if="mode === 'practices' && loading" class="loading-state">Загрузка...</div>
      <div v-else-if="mode === 'practices' && error" class="error-state">{{ error }}</div>

      <!-- Table -->
      <div v-if="mode === 'practices' && !loading && !error" class="table-wrap card">
        <div v-if="filteredPractices.length === 0" class="empty-table">
          <p>Нет практик в этой категории</p>
        </div>
        <table v-else class="mod-table">
          <thead>
            <tr>
              <th>Практика</th>
              <th>Автор</th>
              <th>Категория</th>
              <th>Дата</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in filteredPractices" :key="p.id" class="table-row" @click="openDetail(p.id)">
              <td class="col-practice">
                <div class="practice-thumb"></div>
                <div>
                  <div class="practice-name">{{ p.title }}</div>
                  <span class="badge badge-pending">На модерации</span>
                </div>
              </td>
              <td class="col-author">Автор #{{ p.author_id }}</td>
              <td class="col-cat">{{ p.category }}</td>
              <td class="col-date">{{ formatDate(p.created_at) }}</td>
              <td class="col-actions" @click.stop>
                <button class="action-approve" @click="approve(p.id)" :disabled="actionLoading === p.id">
                  ✓ Одобрить
                </button>
                <button class="action-reject" @click="reject(p.id)" :disabled="actionLoading === p.id">
                  ✕ Отклонить
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Preview pane -->
      <div v-if="mode === 'practices' && selectedPractice" class="preview-card card">
        <h3 class="preview-title">{{ selectedPractice.title }}</h3>
        <p class="preview-meta">{{ selectedPractice.city }} · {{ selectedPractice.category }}</p>
        <p class="preview-desc">{{ selectedPractice.description }}</p>
        <div class="preview-actions">
          <button class="action-approve" @click="approve(selectedPractice.id)">✓ Одобрить</button>
          <button class="action-reject" @click="reject(selectedPractice.id)">✕ Отклонить</button>
          <router-link :to="`/admin/practice/${selectedPractice.id}`" class="btn btn-outline btn-sm">Открыть →</router-link>
        </div>
      </div>
      <div v-else-if="mode === 'practices' && !loading && filteredPractices.length > 0" class="preview-placeholder card">
        <p>Выберите практику в списке для просмотра деталей и принятия решения</p>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../../api'
import AdminLayout from '../../components/AdminLayout.vue'

const router = useRouter()
const mode = ref('practices')

// Practices
const practices = ref([])
const loading = ref(true)
const error = ref(null)
const activeTab = ref('pending')
const search = ref('')
const actionLoading = ref(null)
const selectedPractice = ref(null)

// Verifications
const verifications = ref([])
const vLoading = ref(false)
const vError = ref(null)
const vActionLoading = ref(null)

const socialRoleLabels = {
  specialist: 'Эксперт / специалист',
  official: 'Представитель администрации',
  activist: 'Бизнес / организация',
  resident: 'Житель города',
}
function roleLabel(role) { return socialRoleLabels[role] || role }

const tabs = computed(() => [
  { value: 'all', label: 'Все', count: practices.value.length },
  { value: 'pending', label: 'Ожидают', count: practices.value.length },
  { value: 'approved', label: 'Одобрено', count: 0 },
  { value: 'rejected', label: 'Отклонено', count: 0 },
])

const filteredPractices = computed(() => {
  let list = practices.value
  if (search.value) {
    const q = search.value.toLowerCase()
    list = list.filter(p => p.title.toLowerCase().includes(q) || String(p.author_id).includes(q))
  }
  return list
})

function openDetail(id) {
  selectedPractice.value = practices.value.find(p => p.id === id) || null
}

async function approve(id) {
  actionLoading.value = id
  try {
    await api.approve(id)
    practices.value = practices.value.filter(p => p.id !== id)
    if (selectedPractice.value?.id === id) selectedPractice.value = null
  } catch {} finally {
    actionLoading.value = null
  }
}

async function reject(id) {
  actionLoading.value = id
  try {
    await api.reject(id)
    practices.value = practices.value.filter(p => p.id !== id)
    if (selectedPractice.value?.id === id) selectedPractice.value = null
  } catch {} finally {
    actionLoading.value = null
  }
}

function formatDate(d) {
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

async function approveVerif(id) {
  vActionLoading.value = id
  try {
    await api.approveVerification(id)
    verifications.value = verifications.value.filter(u => u.id !== id)
  } catch {} finally {
    vActionLoading.value = null
  }
}

async function rejectVerif(id) {
  vActionLoading.value = id
  try {
    await api.rejectVerification(id)
    verifications.value = verifications.value.filter(u => u.id !== id)
  } catch {} finally {
    vActionLoading.value = null
  }
}

onMounted(async () => {
  try {
    practices.value = await api.getPending()
  } catch {
    error.value = 'Ошибка загрузки'
  } finally {
    loading.value = false
  }

  vLoading.value = true
  try {
    verifications.value = await api.getPendingVerifications()
  } catch {
    vError.value = 'Ошибка загрузки верификаций'
  } finally {
    vLoading.value = false
  }
})
</script>

<style scoped>
.moderation-page {}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 24px;
}
.page-title { font-size: 26px; font-weight: 700; margin-bottom: 4px; }
.page-sub { font-size: 14px; color: var(--c-muted); }

/* Toolbar */
.toolbar {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.search-wrap { position: relative; flex: 1; min-width: 240px; }
.search-icon {
  position: absolute;
  left: 12px; top: 50%;
  transform: translateY(-50%);
  color: var(--c-muted);
  display: flex;
}
.search-input {
  width: 100%;
  padding: 10px 14px 10px 36px;
  border: 1.5px solid var(--c-border);
  border-radius: 8px;
  font-size: 14px;
  font-family: inherit;
  outline: none;
  background: white;
}
.search-input:focus { border-color: var(--c-green); }
.filter-btn {
  padding: 9px 16px;
  border: 1.5px solid var(--c-border);
  border-radius: 8px;
  background: white;
  font-size: 14px;
  font-family: inherit;
  cursor: pointer;
  color: var(--c-muted);
  transition: border-color 0.15s;
}
.filter-btn:hover { border-color: #9CA3AF; }

/* Tabs */
.status-tabs { display: flex; gap: 4px; margin-bottom: 20px; }
.status-tab {
  padding: 9px 18px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  border: 1.5px solid var(--c-border);
  background: white;
  color: var(--c-muted);
  cursor: pointer;
  transition: all 0.15s;
}
.status-tab.active {
  background: var(--c-navy);
  color: white;
  border-color: var(--c-navy);
}

/* Table */
.table-wrap { overflow: hidden; margin-bottom: 20px; }
.mod-table { width: 100%; border-collapse: collapse; }
.mod-table th {
  padding: 12px 16px;
  font-size: 12px;
  font-weight: 600;
  color: var(--c-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  background: #F8FAFB;
  border-bottom: 1px solid var(--c-border);
  text-align: left;
}
.table-row {
  border-bottom: 1px solid var(--c-border);
  cursor: pointer;
  transition: background 0.1s;
}
.table-row:hover { background: #F8FAFB; }
.table-row:last-child { border-bottom: none; }
.mod-table td { padding: 14px 16px; font-size: 14px; vertical-align: middle; }

.col-practice { display: flex; align-items: center; gap: 12px; }
.practice-thumb {
  width: 40px;
  height: 40px;
  background: #E5E7EB;
  border-radius: 6px;
  flex-shrink: 0;
}
.practice-name { font-weight: 500; margin-bottom: 4px; }
.col-author { color: var(--c-muted); }
.col-cat { color: var(--c-text); }
.col-date { color: var(--c-muted); white-space: nowrap; }
.col-actions { display: flex; gap: 8px; align-items: center; }

.action-approve {
  padding: 7px 14px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  background: #D1FAE5;
  color: #065F46;
  border: 1px solid #A7F3D0;
  cursor: pointer;
  transition: background 0.15s;
  white-space: nowrap;
}
.action-approve:hover:not(:disabled) { background: #A7F3D0; }
.action-approve:disabled { opacity: 0.5; }

.action-reject {
  padding: 7px 14px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  background: #FEE2E2;
  color: #991B1B;
  border: 1px solid #FECACA;
  cursor: pointer;
  transition: background 0.15s;
  white-space: nowrap;
}
.action-reject:hover:not(:disabled) { background: #FECACA; }
.action-reject:disabled { opacity: 0.5; }

/* Preview */
.preview-card { padding: 24px; }
.preview-title { font-size: 18px; font-weight: 600; margin-bottom: 8px; }
.preview-meta { font-size: 13px; color: var(--c-muted); margin-bottom: 12px; }
.preview-desc { font-size: 14px; color: var(--c-text); line-height: 1.6; margin-bottom: 16px; }
.preview-actions { display: flex; gap: 10px; align-items: center; }

.preview-placeholder {
  padding: 28px;
  text-align: center;
  color: var(--c-muted);
  font-size: 14px;
}

.loading-state, .error-state { padding: 32px 0; color: var(--c-muted); text-align: center; }
.empty-table { padding: 40px; text-align: center; color: var(--c-muted); font-size: 14px; }

.mode-switcher { display: flex; gap: 6px; }
.mode-btn {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  border: 1.5px solid var(--c-border);
  background: white;
  color: var(--c-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.15s;
  font-family: inherit;
}
.mode-btn.active { background: var(--c-navy); color: white; border-color: var(--c-navy); }
.badge-num {
  background: #EF4444;
  color: white;
  font-size: 11px;
  font-weight: 700;
  border-radius: 10px;
  padding: 1px 6px;
}
.role-badge {
  display: inline-block;
  padding: 3px 8px;
  background: #EFF6FF;
  color: #1D4ED8;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}
.doc-link { color: var(--c-green); font-size: 13px; text-decoration: none; }
.doc-link:hover { text-decoration: underline; }

@media (max-width: 768px) {
  .page-title { font-size: 20px; }
  .toolbar { flex-direction: column; }
  .search-wrap { min-width: 100%; }
  .filter-btn { width: 100%; }
  .status-tabs { overflow-x: auto; flex-wrap: nowrap; scrollbar-width: none; }
  .status-tabs::-webkit-scrollbar { display: none; }
  .status-tab { flex-shrink: 0; }

  /* Stack table as cards on mobile */
  .mod-table thead { display: none; }
  .mod-table, .mod-table tbody, .table-row, .mod-table td { display: block; width: 100%; }
  .table-row { padding: 12px 16px; border-bottom: 1px solid var(--c-border); }
  .mod-table td { padding: 4px 0; border: none; }
  .col-practice { flex-direction: row; align-items: center; }
  .col-author::before { content: 'Автор: '; color: var(--c-muted); font-size: 12px; }
  .col-cat::before { content: 'Категория: '; color: var(--c-muted); font-size: 12px; }
  .col-date::before { content: 'Дата: '; color: var(--c-muted); font-size: 12px; }
  .col-actions { margin-top: 8px; }
}
</style>
