<template>
  <div class="auth-page">
    <!-- Left panel (same as login) -->
    <div class="auth-left">
      <div class="auth-brand">
        <div class="auth-logo">
          <span class="logo-icon"></span>
          <span class="logo-name">GoodCity</span>
        </div>
        <h2 class="auth-tagline">Платформа лучших практик городского благоустройства</h2>
      </div>
      <div class="auth-features">
        <div class="auth-feature">1 240 практик из 342 городов России</div>
        <div class="auth-feature">Система экспертной оценки и рейтингов</div>
        <div class="auth-feature">8 500+ участников сообщества</div>
      </div>
      <div class="auth-float-card">
        <div class="float-card-score">4.9 ★</div>
        <div class="float-card-sub">Средний рейтинг</div>
        <div class="float-card-note">+342 практики за месяц</div>
      </div>
      <div class="auth-circles">
        <div class="ac ac-1"></div>
        <div class="ac ac-2"></div>
      </div>
    </div>

    <!-- Right form -->
    <div class="auth-right">
      <div class="auth-form-wrap">
        <div class="form-header">
          <h1 class="form-title">Создать аккаунт</h1>
          <p class="already">Уже есть аккаунт? <router-link to="/auth/login" class="switch-link">Войти</router-link></p>
        </div>

        <form class="auth-form" @submit.prevent="submit">
          <div class="fields-section">
            <p class="section-label">Основная информация</p>
            <div class="fields-grid-2">
              <div class="form-field">
                <label class="form-label">Полное имя (ФИО) *</label>
                <input v-model="form.name" type="text" class="form-input" placeholder="Иванов Иван Иванович" />
              </div>
              <div class="form-field">
                <label class="form-label">Email *</label>
                <input v-model="form.email" type="email" class="form-input" placeholder="your@email.ru" required autocomplete="email" />
              </div>
              <div class="form-field">
                <label class="form-label">Пароль *</label>
                <input v-model="form.password" type="password" class="form-input" placeholder="••••••••" required minlength="6" autocomplete="new-password" />
              </div>
              <div class="form-field">
                <label class="form-label">Подтверждение пароля *</label>
                <input v-model="form.passwordConfirm" type="password" class="form-input" placeholder="••••••••" required />
              </div>
            </div>
          </div>

          <div class="fields-section">
            <div class="section-header-row">
              <p class="section-label">Тип аккаунта *</p>
              <span class="section-hint">Определяет права доступа и уровень ответственности</span>
            </div>
            <div class="account-types">
              <label v-for="type in accountTypes" :key="type.value" class="account-type" :class="{ selected: form.accountType === type.value }">
                <input type="radio" v-model="form.accountType" :value="type.value" class="radio-hidden" />
                <span class="type-icon">{{ type.icon }}</span>
                <span class="type-name">{{ type.name }}</span>
                <span class="type-desc">{{ type.desc }}</span>
                <span class="type-check" v-if="form.accountType === type.value">✓</span>
              </label>
            </div>
          </div>

          <div class="fields-section">
            <p class="section-label">Город и социальный профиль</p>
            <div class="fields-grid-2">
              <!-- City autocomplete -->
              <div class="form-field ac-wrap">
                <label class="form-label">Город *</label>
                <input
                  v-model="form.city"
                  type="text"
                  class="form-input"
                  placeholder="Начните вводить город..."
                  required
                  autocomplete="off"
                  @input="onCityInput"
                  @blur="cityBlur"
                  @focus="form.city.length >= 2 && onCityInput()"
                />
                <ul v-if="citySuggestions.length" class="ac-dropdown">
                  <li
                    v-for="s in citySuggestions"
                    :key="s.value"
                    class="ac-item"
                    @mousedown.prevent="selectCity(s)"
                  >{{ s.data.city || s.value }}</li>
                </ul>
              </div>
              <!-- District autocomplete -->
              <div class="form-field ac-wrap">
                <label class="form-label">Район</label>
                <input
                  v-model="form.district"
                  type="text"
                  class="form-input"
                  :placeholder="form.city ? 'Начните вводить район...' : 'Сначала выберите город'"
                  :disabled="!form.city"
                  autocomplete="off"
                  @input="onDistrictInput"
                  @blur="districtBlur"
                />
                <ul v-if="districtSuggestions.length" class="ac-dropdown">
                  <li
                    v-for="s in districtSuggestions"
                    :key="s.value"
                    class="ac-item"
                    @mousedown.prevent="selectDistrict(s)"
                  >{{ s.data.city_district || s.value }}</li>
                </ul>
              </div>
            </div>
          </div>

          <!-- Document upload for non-resident roles -->
          <div v-if="form.accountType !== 'resident'" class="fields-section">
            <div class="section-header-row">
              <p class="section-label">Документ для верификации роли</p>
              <span class="section-hint">Необязательно — можно загрузить позже</span>
            </div>
            <label class="upload-area" :class="{ 'has-file': docFile }">
              <input type="file" accept=".pdf,.jpg,.jpeg,.png" class="upload-hidden" @change="onFileChange" />
              <span v-if="!docFile" class="upload-placeholder">
                <span class="upload-icon">📄</span>
                <span>Загрузите диплом, удостоверение или выписку ЕГРЮЛ</span>
                <span class="upload-hint">PDF, JPG, PNG — до 10 МБ</span>
              </span>
              <span v-else class="upload-filename">✓ {{ docFile.name }}</span>
            </label>
          </div>

          <label class="checkbox-agree">
            <input type="checkbox" v-model="form.agreed" class="checkbox" required />
            <span>Я согласен с условиями использования платформы и политикой конфиденциальности</span>
          </label>

          <p v-if="error" class="error-msg">{{ error }}</p>

          <button type="submit" class="btn-submit" :disabled="loading || !form.agreed">
            {{ loading ? 'Создание...' : 'Создать аккаунт' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'
import { api } from '../../api'
import { useCitySearch, useDistrictSearch } from '../../composables/useDadata'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const error = ref(null)
const docFile = ref(null)

const form = ref({
  name: '', email: '', password: '', passwordConfirm: '',
  accountType: 'resident', city: '', district: '', agreed: false
})

// DaData autocomplete
const { suggestions: citySuggestions, search: searchCity, clear: clearCitySuggestions } = useCitySearch()
const { suggestions: districtSuggestions, search: searchDistrict, clear: clearDistrictSuggestions } = useDistrictSearch()

let cityTimer = null
function onCityInput() {
  clearTimeout(cityTimer)
  cityTimer = setTimeout(() => searchCity(form.value.city), 300)
}
function cityBlur() { setTimeout(clearCitySuggestions, 150) }
function selectCity(s) {
  form.value.city = s.data.city || s.value
  form.value.district = ''
  clearCitySuggestions()
}

let districtTimer = null
function onDistrictInput() {
  clearTimeout(districtTimer)
  districtTimer = setTimeout(() => searchDistrict(form.value.district, form.value.city), 300)
}
function districtBlur() { setTimeout(clearDistrictSuggestions, 150) }
function selectDistrict(s) {
  form.value.district = s.data.city_district || s.value
  clearDistrictSuggestions()
}

function onFileChange(e) {
  docFile.value = e.target.files[0] || null
}

const accountTypes = [
  { value: 'resident', icon: '🏘', name: 'Житель города', desc: 'Оцениваю практики, делюсь опытом', socialRole: 'resident' },
  { value: 'expert', icon: '🎓', name: 'Эксперт / специалист', desc: 'Профессиональная оценка по области знаний', socialRole: 'specialist' },
  { value: 'admin', icon: '🏛', name: 'Представитель администрации', desc: 'Управляю городскими проектами и инициативами', socialRole: 'official' },
  { value: 'business', icon: '🏢', name: 'Бизнес / организация', desc: 'Реализую практики в рамках проектов', socialRole: 'activist' },
]

async function submit() {
  if (form.value.password !== form.value.passwordConfirm) {
    error.value = 'Пароли не совпадают'
    return
  }
  loading.value = true
  error.value = null
  try {
    const selected = accountTypes.find(t => t.value === form.value.accountType)
    const socialRole = selected?.socialRole || 'resident'
    await auth.register(form.value.email, form.value.password, socialRole, form.value.name, form.value.city, form.value.district)
    if (docFile.value) {
      const fd = new FormData()
      fd.append('document', docFile.value)
      await api.submitVerification(fd).catch(() => {})
    }
    router.push('/')
  } catch (e) {
    error.value = e.data?.error || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page { display: flex; min-height: 100vh; }

.auth-left {
  width: 380px;
  flex-shrink: 0;
  background: var(--c-navy);
  padding: 40px 36px;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
}
.auth-brand { margin-bottom: 32px; }
.auth-logo { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.logo-icon { width: 32px; height: 32px; background: var(--c-green); border-radius: 8px; }
.logo-name { font-size: 17px; font-weight: 700; color: white; }
.auth-tagline { font-size: 22px; font-weight: 700; color: white; line-height: 1.35; }
.auth-features { display: flex; flex-direction: column; gap: 8px; margin-bottom: 28px; position: relative; z-index: 1; }
.auth-feature { background: rgba(255,255,255,0.1); border-radius: 8px; padding: 10px 14px; font-size: 13px; color: rgba(255,255,255,0.85); border-left: 3px solid var(--c-green); }
.auth-float-card { background: rgba(255,255,255,0.12); border: 1px solid rgba(255,255,255,0.2); border-radius: 14px; padding: 16px 18px; position: relative; z-index: 1; margin-top: auto; }
.float-card-score { font-size: 22px; font-weight: 700; color: white; }
.float-card-sub { font-size: 12px; color: rgba(255,255,255,0.65); margin: 3px 0; }
.float-card-note { font-size: 11px; color: rgba(255,255,255,0.5); }
.auth-circles { position: absolute; inset: 0; pointer-events: none; }
.ac { position: absolute; border-radius: 50%; opacity: 0.15; }
.ac-1 { width: 220px; height: 220px; background: var(--c-green); bottom: 80px; right: -60px; }
.ac-2 { width: 180px; height: 160px; background: #4B9EFF; bottom: 0; right: 40px; }

.auth-right {
  flex: 1;
  display: flex;
  justify-content: center;
  background: white;
  padding: 40px 48px;
  overflow-y: auto;
}
.auth-form-wrap { max-width: 620px; }

.form-header { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 28px; }
.form-title { font-size: 26px; font-weight: 700; }
.already { font-size: 14px; color: var(--c-muted); }
.switch-link { color: var(--c-green); font-weight: 500; }

.auth-form { display: flex; flex-direction: column; gap: 24px; }
.fields-section {}
.section-label { font-size: 14px; font-weight: 600; color: var(--c-text); margin-bottom: 14px; }
.section-header-row { display: flex; align-items: center; gap: 12px; margin-bottom: 14px; }
.section-hint { font-size: 12px; color: var(--c-muted); }

.fields-grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 12px 16px; }
.fields-grid-3 { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 12px 16px; margin-bottom: 12px; }
.form-field { margin-bottom: 0; }

/* Account types */
.account-types { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.account-type {
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding: 14px;
  border: 1.5px solid var(--c-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  position: relative;
  transition: border-color 0.15s;
}
.account-type.selected { border-color: var(--c-green); background: #F0FDF4; }
.radio-hidden { position: absolute; opacity: 0; width: 0; height: 0; }
.type-icon { font-size: 20px; margin-bottom: 4px; }
.type-name { font-size: 14px; font-weight: 600; }
.type-desc { font-size: 12px; color: var(--c-muted); line-height: 1.4; }
.type-check {
  position: absolute;
  top: 10px; right: 10px;
  width: 20px; height: 20px;
  background: var(--c-green);
  color: white;
  border-radius: 50%;
  font-size: 11px;
  display: flex;
  align-items: center;
  justify-content: center;
}


.ac-wrap { position: relative; }
.ac-dropdown {
  position: absolute;
  top: calc(100% + 2px);
  left: 0; right: 0;
  background: white;
  border: 1.5px solid var(--c-border);
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.10);
  z-index: 100;
  list-style: none;
  margin: 0; padding: 4px 0;
  max-height: 220px;
  overflow-y: auto;
}
.ac-item {
  padding: 9px 14px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.1s;
}
.ac-item:hover { background: #F0FDF4; }

.upload-area {
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px dashed var(--c-border);
  border-radius: var(--radius-sm);
  padding: 20px;
  cursor: pointer;
  transition: border-color 0.15s;
  min-height: 80px;
  position: relative;
}
.upload-area:hover { border-color: var(--c-green); }
.upload-area.has-file { border-color: var(--c-green); background: #F0FDF4; }
.upload-hidden { position: absolute; inset: 0; opacity: 0; cursor: pointer; width: 100%; }
.upload-placeholder { display: flex; flex-direction: column; align-items: center; gap: 4px; text-align: center; }
.upload-icon { font-size: 24px; }
.upload-placeholder span:nth-child(2) { font-size: 13px; color: var(--c-text); }
.upload-hint { font-size: 12px; color: var(--c-muted); }
.upload-filename { font-size: 13px; color: #065F46; font-weight: 500; }

.checkbox-agree {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font-size: 13px;
  color: var(--c-muted);
  cursor: pointer;
}
.checkbox { accent-color: var(--c-green); width: 16px; height: 16px; margin-top: 2px; flex-shrink: 0; }

.btn-submit {
  width: 100%;
  padding: 14px;
  background: var(--c-green);
  color: white;
  font-size: 15px;
  font-weight: 600;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font-family: inherit;
  transition: background 0.15s;
}
.btn-submit:hover:not(:disabled) { background: var(--c-green-dark); }
.btn-submit:disabled { opacity: 0.55; cursor: not-allowed; }

.error-msg { color: var(--c-red-dark); font-size: 13px; }

@media (max-width: 768px) {
  .auth-page { flex-direction: column; }

  .auth-left {
    width: 100%;
    padding: 24px 20px;
    min-height: auto;
    align-items: center;
    text-align: center;
  }
  .auth-brand { margin-bottom: 0; }
  .auth-tagline { font-size: 16px; }
  .auth-features { display: none; }
  .auth-float-card { display: none; }

  .auth-right {
    padding: 24px 20px 40px;
    width: 100%;
  }
  .auth-form-wrap { max-width: 100%; }
  .form-header { flex-direction: column; gap: 4px; align-items: flex-start; }
  .form-title { font-size: 22px; }

  .fields-grid-2 { grid-template-columns: 1fr; }
  .account-types { grid-template-columns: 1fr; }
}
</style>
