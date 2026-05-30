<template>
  <div class="profile-bg">
    <div class="profile-layout container">

      <!-- Sidebar -->
      <aside class="sidebar">
        <nav class="sidebar-nav">
          <button
            v-for="s in sections"
            :key="s.id"
            class="sidebar-item"
            :class="{ active: activeSection === s.id }"
            @click="scrollTo(s.id)"
          >
            <span class="sidebar-icon">{{ s.icon }}</span>
            <span>{{ s.label }}</span>
          </button>
        </nav>
      </aside>

      <!-- Main content -->
      <main class="profile-main">
        <div v-if="loading" class="state-msg">Загрузка...</div>
        <div v-else-if="loadError" class="state-msg err">{{ loadError }}</div>

        <template v-else>
          <div class="profile-title-row">
            <div>
              <h1 class="profile-title">Личный профиль</h1>
              <p class="profile-sub">Последнее обновление: {{ formatDate(user.created_at) }}</p>
            </div>
          </div>

          <!-- 1. Основная информация -->
          <section :id="'sec-info'" class="card profile-section">
            <div class="section-head">
              <span class="section-title">Основная информация</span>
              <button class="edit-btn" @click="editing.info = !editing.info">✎</button>
            </div>
            <div class="fields-grid">
              <div class="field-row">
                <label class="field-label">Полное имя (ФИО)</label>
                <input v-if="editing.info" v-model="form.name" class="field-input" placeholder="Иванов Иван Иванович" />
                <span v-else class="field-value">{{ user.name || '—' }}</span>
              </div>
              <div class="field-row">
                <label class="field-label">Email</label>
                <span class="field-value">{{ user.email }}</span>
              </div>
              <div class="field-row">
                <label class="field-label">Телефон</label>
                <input v-if="editing.info" v-model="form.phone" class="field-input" placeholder="+7 (___) ___-__-__" />
                <span v-else class="field-value">{{ user.phone || '—' }}</span>
              </div>
              <div class="field-row">
                <label class="field-label">Тип аккаунта</label>
                <span class="field-value">{{ roleLabelMap[user.social_role] ?? user.social_role }}</span>
              </div>
            </div>
          </section>

          <!-- 2. Город и район -->
          <section :id="'sec-geo'" class="card profile-section">
            <div class="section-head">
              <span class="section-title">Город и район</span>
              <button class="edit-btn" @click="editing.geo = !editing.geo">✎</button>
            </div>
            <p class="section-desc">Используется для расчёта географической близости к оцениваемым объектам</p>
            <div class="fields-grid">
              <div class="field-row">
                <label class="field-label">Город</label>
                <div v-if="editing.geo" class="ac-wrap">
                  <input v-model="form.city" class="field-input" placeholder="Начните вводить город..." autocomplete="off"
                    @input="onCityInput" @blur="cityBlur" />
                  <ul v-if="citySuggestions.length" class="ac-dropdown">
                    <li v-for="s in citySuggestions" :key="s.value" class="ac-item" @mousedown.prevent="selectCity(s)">
                      {{ s.data.city || s.value }}
                    </li>
                  </ul>
                </div>
                <span v-else class="field-value">{{ user.city || '—' }}</span>
              </div>
              <div class="field-row">
                <label class="field-label">Район</label>
                <div v-if="editing.geo" class="ac-wrap">
                  <input v-model="form.district" class="field-input" :placeholder="form.city ? 'Начните вводить район...' : 'Сначала выберите город'"
                    :disabled="!form.city" autocomplete="off" @input="onDistrictInput" @blur="districtBlur" />
                  <ul v-if="districtSuggestions.length" class="ac-dropdown">
                    <li v-for="s in districtSuggestions" :key="s.value" class="ac-item" @mousedown.prevent="selectDistrict(s)">
                      {{ s.data.city_district || s.value }}
                    </li>
                  </ul>
                </div>
                <span v-else class="field-value">{{ user.district || '—' }}</span>
              </div>
            </div>
          </section>

          <!-- 3. Социальный профиль -->
          <section :id="'sec-social'" class="card profile-section">
            <div class="section-head">
              <span class="section-title">Социальный профиль</span>
              <span class="section-info">ℹ Влияет на вес ваших оценок в системе</span>
            </div>

            <div class="subsection">
              <p class="subsection-label">Социальная роль</p>
              <p class="subsection-desc">Ваша роль в городском сообществе</p>
              <div class="choice-row">
                <label v-for="r in socialRoles" :key="r.value" class="choice-item" :class="{ active: form.social_role === r.value }">
                  <input type="radio" v-model="form.social_role" :value="r.value" class="choice-hidden" />
                  {{ r.label }}
                </label>
              </div>
            </div>

          </section>

          <!-- 4. Репутация -->
          <section :id="'sec-reputation'" class="card profile-section">
            <div class="section-head">
              <span class="section-title">Моя репутация в системе</span>
              <span class="section-readonly">Рассчитывается автоматически · только для чтения</span>
            </div>
            <div class="metrics-grid">
              <div class="metric-card metric-blue">
                <span class="metric-val">{{ user.reputation_score?.toFixed(2) ?? '0.00' }}</span>
                <span class="metric-name">Индекс активности</span>
              </div>
              <div class="metric-card metric-green">
                <span class="metric-val">{{ user.profile_weight?.toFixed(2) ?? '1.00' }}</span>
                <span class="metric-name">Профильный вес</span>
              </div>
              <div class="metric-card metric-indigo">
                <span class="metric-val">{{ user.is_role_verified ? '✓' : '—' }}</span>
                <span class="metric-name">Верификация роли</span>
              </div>
              <div class="metric-card metric-navy">
                <span class="metric-val">{{ user.is_role_verified ? user.social_role !== 'resident' ? roleLabelMap[user.social_role]?.split(' ')[0] : '—' : '—' }}</span>
                <span class="metric-name">Верифицированная роль</span>
              </div>
            </div>
          </section>

          <!-- 5. Верификация (если нужна) -->
          <section v-if="user.social_role !== 'resident'" :id="'sec-verif'" class="card profile-section">
            <div class="section-head">
              <span class="section-title">Верификация роли</span>
              <span v-if="user.is_role_verified" class="badge-verified">✓ Подтверждена</span>
              <span v-else class="badge-pending">Ожидает проверки</span>
            </div>
            <div v-if="!user.is_role_verified">
              <p class="section-desc">Загрузите документ, подтверждающий вашу роль (диплом, удостоверение, выписка ЕГРЮЛ).</p>
              <div v-if="user.verification_doc_url" class="verif-status">
                <span>📄 Документ загружен — ожидает модерации</span>
                <a :href="user.verification_doc_url" target="_blank" class="doc-link">Открыть</a>
              </div>
              <template v-else>
                <label class="upload-area" :class="{ 'has-file': docFile }">
                  <input type="file" accept=".pdf,.jpg,.jpeg,.png" class="upload-hidden" @change="e => docFile = e.target.files[0]" />
                  <span v-if="!docFile" class="upload-placeholder">
                    <span>📄</span><span>PDF, JPG, PNG — до 10 МБ</span>
                  </span>
                  <span v-else class="upload-filename">✓ {{ docFile.name }}</span>
                </label>
                <button v-if="docFile" class="btn-primary" :disabled="uploading" @click="uploadDoc">
                  {{ uploading ? 'Отправка...' : 'Отправить на проверку' }}
                </button>
              </template>
              <p v-if="uploadMsg" class="success-msg">{{ uploadMsg }}</p>
            </div>
          </section>

          <!-- Actions -->
          <div class="actions-row">
            <button class="btn-save" :disabled="saving" @click="save">
              {{ saving ? 'Сохранение...' : 'Сохранить изменения' }}
            </button>
            <button class="btn-cancel" @click="reset">Отмена</button>
            <p v-if="saveError" class="err-msg">{{ saveError }}</p>
            <p v-if="saveSuccess" class="success-msg">Изменения сохранены</p>
          </div>

        </template>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { api } from '../api'
import { useCitySearch, useDistrictSearch } from '../composables/useDadata'

const user = ref(null)
const loading = ref(true)
const loadError = ref(null)
const saving = ref(false)
const saveError = ref(null)
const saveSuccess = ref(false)
const uploading = ref(false)
const uploadMsg = ref(null)
const docFile = ref(null)
const activeSection = ref('sec-info')

const editing = reactive({ info: false, geo: false })

const form = reactive({
  name: '', phone: '', city: '', district: '',
  social_role: 'resident', income_bracket: 'middle',
})

const roleLabelMap = {
  resident: 'Житель города',
  specialist: 'Эксперт / специалист',
  official: 'Представитель администрации',
  activist: 'Бизнес / организация',
}
const incomeLabelMap = { low: 'Невысокий', middle: 'Средний', high: 'Высокий' }

const socialRoles = [
  { value: 'resident', label: 'Житель' },
  { value: 'specialist', label: 'Эксперт' },
  { value: 'official', label: 'Чиновник' },
  { value: 'activist', label: 'Бизнес' },
]
const incomes = [
  { value: 'low', label: 'Невысокий' },
  { value: 'middle', label: 'Средний' },
  { value: 'high', label: 'Высокий' },
]

const sections = [
  { id: 'sec-info', icon: '👤', label: 'Основная информация' },
  { id: 'sec-geo', icon: '🏙️', label: 'Город и район' },
  { id: 'sec-social', icon: '📊', label: 'Социальный профиль' },
  { id: 'sec-reputation', icon: '📈', label: 'Моя репутация' },
  { id: 'sec-verif', icon: '🔒', label: 'Верификация' },
]

// DaData
const { suggestions: citySuggestions, search: searchCity, clear: clearCity } = useCitySearch()
const { suggestions: districtSuggestions, search: searchDistrict, clear: clearDistrict } = useDistrictSearch()
let cityTimer, districtTimer
function onCityInput() { clearTimeout(cityTimer); cityTimer = setTimeout(() => searchCity(form.city), 300) }
function cityBlur() { setTimeout(clearCity, 150) }
function selectCity(s) { form.city = s.data.city || s.value; form.district = ''; clearCity() }
function onDistrictInput() { clearTimeout(districtTimer); districtTimer = setTimeout(() => searchDistrict(form.district, form.city), 300) }
function districtBlur() { setTimeout(clearDistrict, 150) }
function selectDistrict(s) { form.district = s.data.city_district || s.value; clearDistrict() }

function fillForm(u) {
  form.name = u.name ?? ''
  form.phone = u.phone ?? ''
  form.city = u.city ?? ''
  form.district = u.district ?? ''
  form.social_role = u.social_role ?? 'resident'
  form.income_bracket = u.income_bracket ?? 'middle'
}

function reset() {
  if (user.value) fillForm(user.value)
  editing.info = false
  editing.geo = false
  saveSuccess.value = false
  saveError.value = null
}

function scrollTo(id) {
  activeSection.value = id
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
}

async function save() {
  saving.value = true
  saveError.value = null
  saveSuccess.value = false
  try {
    user.value = await api.updateMe({
      name: form.name, phone: form.phone,
      city: form.city, district: form.district,
      social_role: form.social_role,
      income_bracket: form.income_bracket,
    })
    fillForm(user.value)
    editing.info = false
    editing.geo = false
    saveSuccess.value = true
  } catch {
    saveError.value = 'Не удалось сохранить изменения'
  } finally {
    saving.value = false
  }
}

async function uploadDoc() {
  if (!docFile.value) return
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('document', docFile.value)
    await api.submitVerification(fd)
    uploadMsg.value = 'Документ отправлен на проверку!'
    docFile.value = null
    user.value = await api.getMe()
  } catch {
    uploadMsg.value = 'Ошибка загрузки'
  } finally {
    uploading.value = false
  }
}

onMounted(async () => {
  try {
    user.value = await api.getMe()
    fillForm(user.value)
  } catch {
    loadError.value = 'Не удалось загрузить профиль'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.profile-bg { background: #f5f7fa; min-height: 100vh; }
.profile-layout {
  display: flex;
  gap: 28px;
  padding-top: 36px;
  padding-bottom: 60px;
  align-items: flex-start;
}

/* Sidebar */
.sidebar { width: 220px; flex-shrink: 0; position: sticky; top: 88px; }
.sidebar-nav { display: flex; flex-direction: column; gap: 2px; }
.sidebar-item {
  display: flex; align-items: center; gap: 10px;
  padding: 10px 14px; border-radius: 8px;
  font-size: 13px; color: #6b7280;
  background: none; border: none; cursor: pointer;
  text-align: left; font-family: inherit;
  transition: background 0.15s, color 0.15s;
}
.sidebar-item:hover { background: #e8edf3; color: #1a1a2e; }
.sidebar-item.active { background: #e8edf3; color: #1a1a2e; font-weight: 600; }
.sidebar-icon { font-size: 16px; flex-shrink: 0; }

/* Main */
.profile-main { flex: 1; display: flex; flex-direction: column; gap: 20px; min-width: 0; }
.profile-title-row { margin-bottom: 4px; }
.profile-title { font-size: 24px; font-weight: 700; color: #1a1a2e; margin-bottom: 4px; }
.profile-sub { font-size: 13px; color: #9ca3af; }

/* Sections */
.profile-section { padding: 24px; }
.section-head { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.section-title { font-size: 15px; font-weight: 700; color: #1a1a2e; flex: 1; }
.section-desc { font-size: 12px; color: #9ca3af; margin-bottom: 16px; margin-top: -10px; }
.section-info { font-size: 11px; color: #9ca3af; }
.section-readonly { font-size: 12px; color: #9ca3af; }
.edit-btn {
  background: none; border: none; cursor: pointer;
  color: #9ca3af; font-size: 14px; padding: 2px 6px;
  border-radius: 4px; font-family: inherit;
  transition: background 0.1s, color 0.1s;
}
.edit-btn:hover { background: #f3f4f6; color: #1a1a2e; }

/* Fields */
.fields-grid { display: flex; flex-direction: column; }
.field-row {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 0; border-bottom: 1px solid #e0e4e9;
  gap: 16px;
}
.field-row:last-child { border-bottom: none; }
.field-label { font-size: 13px; font-weight: 500; color: #6b7280; flex-shrink: 0; min-width: 160px; }
.field-value { font-size: 14px; color: #1a1a2e; }
.field-input {
  flex: 1; padding: 7px 12px; border: 1.5px solid #e0e4e9;
  border-radius: 8px; font-size: 14px; font-family: inherit;
  outline: none; color: #1a1a2e;
}
.field-input:focus { border-color: #2d9b6e; }

/* Autocomplete */
.ac-wrap { position: relative; flex: 1; }
.ac-dropdown {
  position: absolute; top: calc(100% + 2px); left: 0; right: 0;
  background: white; border: 1.5px solid #e0e4e9; border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.10); z-index: 100;
  list-style: none; margin: 0; padding: 4px 0; max-height: 200px; overflow-y: auto;
}
.ac-item { padding: 9px 14px; font-size: 13px; cursor: pointer; }
.ac-item:hover { background: #f0fdf4; }

/* Social choices */
.subsection { margin-bottom: 20px; }
.subsection:last-child { margin-bottom: 0; }
.subsection-label { font-size: 13px; font-weight: 600; color: #1a1a2e; margin-bottom: 4px; }
.subsection-desc { font-size: 11px; color: #9ca3af; margin-bottom: 10px; }
.choice-row { display: flex; gap: 8px; flex-wrap: wrap; }
.choice-item {
  padding: 7px 16px; border-radius: 20px; border: 1.5px solid #e0e4e9;
  font-size: 13px; color: #6b7280; cursor: pointer;
  transition: all 0.15s; user-select: none;
}
.choice-item.active { border-color: #2d9b6e; background: #f0fdf4; color: #1a1a2e; font-weight: 500; }
.choice-hidden { display: none; }

/* Metrics */
.metrics-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px; }
.metric-card {
  border-radius: 12px; padding: 18px 16px;
  display: flex; flex-direction: column; gap: 6px;
  color: white;
}
.metric-blue { background: #4a90d9; }
.metric-green { background: #16a34a; }
.metric-indigo { background: #4f46e5; }
.metric-navy { background: #1e3a5f; }
.metric-val { font-size: 26px; font-weight: 700; }
.metric-name { font-size: 11px; opacity: 0.85; }

/* Verification */
.badge-verified { font-size: 12px; background: #d1fae5; color: #065f46; padding: 3px 10px; border-radius: 20px; }
.badge-pending { font-size: 12px; background: #fef9c3; color: #92400e; padding: 3px 10px; border-radius: 20px; }
.verif-status { display: flex; align-items: center; gap: 10px; font-size: 14px; padding: 12px; background: #fffbeb; border-radius: 8px; border: 1px solid #fde68a; }
.doc-link { color: #2d9b6e; margin-left: auto; font-size: 13px; }
.upload-area {
  display: flex; align-items: center; justify-content: center;
  border: 2px dashed #e0e4e9; border-radius: 8px; padding: 20px;
  cursor: pointer; min-height: 72px; position: relative;
}
.upload-area:hover { border-color: #2d9b6e; }
.upload-area.has-file { border-color: #2d9b6e; background: #f0fdf4; }
.upload-hidden { position: absolute; inset: 0; opacity: 0; cursor: pointer; width: 100%; }
.upload-placeholder { display: flex; gap: 8px; align-items: center; font-size: 13px; color: #6b7280; }
.upload-filename { font-size: 13px; color: #065f46; font-weight: 500; }

/* Actions */
.actions-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.btn-save {
  padding: 11px 24px; background: #2d9b6e; color: white;
  border: none; border-radius: 8px; font-size: 14px; font-weight: 700;
  cursor: pointer; font-family: inherit; transition: background 0.15s;
}
.btn-save:hover:not(:disabled) { background: #238c60; }
.btn-save:disabled { opacity: 0.55; cursor: not-allowed; }
.btn-cancel {
  padding: 11px 24px; background: white; color: #6b7280;
  border: 1.5px solid #e0e4e9; border-radius: 8px; font-size: 14px;
  cursor: pointer; font-family: inherit; transition: border-color 0.15s;
}
.btn-cancel:hover { border-color: #9ca3af; }
.btn-primary {
  margin-top: 12px; padding: 9px 20px; background: #2d9b6e; color: white;
  border: none; border-radius: 8px; font-size: 14px; font-weight: 600;
  cursor: pointer; font-family: inherit;
}
.btn-primary:disabled { opacity: 0.55; }
.success-msg { font-size: 13px; color: #065f46; }
.err-msg { font-size: 13px; color: #dc2626; }
.state-msg { padding: 40px; text-align: center; color: #9ca3af; }
.state-msg.err { color: #dc2626; }

@media (max-width: 900px) {
  .profile-layout { flex-direction: column; }
  .sidebar { width: 100%; position: static; }
  .sidebar-nav { flex-direction: row; flex-wrap: wrap; }
  .metrics-grid { grid-template-columns: 1fr 1fr; }
}
@media (max-width: 560px) {
  .metrics-grid { grid-template-columns: 1fr; }
  .field-row { flex-direction: column; align-items: flex-start; }
  .field-label { min-width: unset; }
}
</style>
