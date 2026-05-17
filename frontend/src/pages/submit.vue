<template>
  <div class="submit-page container">
    <div class="submit-header">
      <h1 class="submit-title">Добавить практику</h1>
      <p class="submit-sub">Поделитесь успешным опытом городского благоустройства</p>
    </div>

    <div class="submit-layout">
      <div class="submit-form-wrap card">
        <form @submit.prevent="submit">
          <div class="form-field">
            <label class="form-label required">Название практики</label>
            <input v-model="form.title" class="form-input" required placeholder="Например: Умные парковки с датчиками движения" />
          </div>

          <div class="fields-grid-2">
            <div class="form-field">
              <label class="form-label required">Город</label>
              <input v-model="form.city" class="form-input" required placeholder="Волгоград" />
            </div>
            <div class="form-field">
              <label class="form-label required">Категория</label>
              <select v-model="form.category" class="form-input form-select" required>
                <option value="">Выберите категорию</option>
                <option v-for="c in categories" :key="c" :value="c">{{ c }}</option>
              </select>
            </div>
          </div>

          <div class="form-field">
            <label class="form-label required">Описание</label>
            <textarea v-model="form.description" class="form-input form-textarea" required
              placeholder="Подробно опишите практику: что сделали, какой результат получили, какие ресурсы потребовались..."
              style="min-height:140px"></textarea>
          </div>

          <div class="upload-area">
            <div class="upload-icon">📷</div>
            <p class="upload-text">Перетащите фото или <span class="upload-link">выберите файл</span></p>
            <p class="upload-hint">PNG, JPG до 10 МБ</p>
          </div>

          <p v-if="error" class="error-msg">{{ error }}</p>
          <div v-if="success" class="success-msg">
            ✓ Практика отправлена на модерацию! Переходим на главную...
          </div>

          <div class="form-actions">
            <router-link to="/" class="btn btn-outline">Отмена</router-link>
            <button type="submit" class="btn btn-primary" :disabled="loading">
              {{ loading ? 'Отправка...' : 'Отправить на модерацию' }}
            </button>
          </div>
        </form>
      </div>

      <div class="submit-tips">
        <div class="tip-card card">
          <h3 class="tip-title">💡 Советы по описанию</h3>
          <ul class="tip-list">
            <li>Укажите конкретные результаты в цифрах</li>
            <li>Опишите масштаб реализации</li>
            <li>Упомяните затраченные ресурсы</li>
            <li>Добавьте временной период</li>
          </ul>
        </div>
        <div class="tip-card card">
          <h3 class="tip-title">⏱ Сроки модерации</h3>
          <p class="tip-text">Ваша практика будет проверена экспертами в течение 1–3 рабочих дней.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const categories = ['Транспорт', 'Благоустройство', 'ЖКХ', 'Экология', 'Безопасность', 'Культура', 'Образование', 'Здоровье', 'Инфраструктура']

const form = ref({ title: '', description: '', city: '', category: '' })
const loading = ref(false)
const error = ref(null)
const success = ref(false)

async function submit() {
  loading.value = true
  error.value = null
  try {
    await api.createPractice(form.value)
    success.value = true
    setTimeout(() => router.push('/'), 2000)
  } catch (e) {
    error.value = e.data?.error || 'Ошибка при отправке'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.submit-page { padding: 40px 24px 64px; }
.submit-header { margin-bottom: 32px; }
.submit-title { font-size: 32px; font-weight: 700; margin-bottom: 8px; }
.submit-sub { font-size: 16px; color: var(--c-muted); }

.submit-layout {
  display: grid;
  grid-template-columns: 1fr 280px;
  gap: 24px;
  align-items: start;
}

.submit-form-wrap { padding: 32px; }
.fields-grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }

.upload-area {
  border: 2px dashed var(--c-border);
  border-radius: var(--radius);
  padding: 32px;
  text-align: center;
  margin: 8px 0 24px;
  cursor: pointer;
  transition: border-color 0.15s;
}
.upload-area:hover { border-color: var(--c-green); }
.upload-icon { font-size: 28px; margin-bottom: 8px; }
.upload-text { font-size: 14px; color: var(--c-muted); }
.upload-link { color: var(--c-green); cursor: pointer; }
.upload-hint { font-size: 12px; color: #9CA3AF; margin-top: 4px; }

.form-actions { display: flex; justify-content: flex-end; gap: 12px; margin-top: 8px; }

.success-msg {
  background: #D1FAE5;
  color: #065F46;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 14px;
  margin-bottom: 16px;
}
.error-msg { margin-bottom: 16px; }

.tip-card { padding: 20px; margin-bottom: 16px; }
.tip-title { font-size: 15px; font-weight: 600; margin-bottom: 12px; }
.tip-list { padding-left: 16px; display: flex; flex-direction: column; gap: 8px; }
.tip-list li { font-size: 14px; color: var(--c-muted); }
.tip-text { font-size: 14px; color: var(--c-muted); line-height: 1.5; }

@media (max-width: 768px) {
  .submit-page { padding: 24px 16px 40px; }
  .submit-title { font-size: 24px; }
  .submit-layout { grid-template-columns: 1fr; }
  .submit-tips { display: none; }
  .fields-grid-2 { grid-template-columns: 1fr; }
  .form-actions { flex-direction: column-reverse; }
  .form-actions > * { width: 100%; justify-content: center; }
}
</style>
