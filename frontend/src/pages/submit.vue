<template>
  <div class="card" style="max-width:600px">
    <h1 style="margin-bottom:1.5rem">Добавить практику</h1>
    <form @submit.prevent="submit">
      <div class="form-group">
        <label>Название *</label>
        <input v-model="form.title" required placeholder="Название практики" />
      </div>
      <div class="form-group">
        <label>Описание *</label>
        <textarea v-model="form.description" required placeholder="Подробное описание практики"></textarea>
      </div>
      <div class="form-group">
        <label>Город *</label>
        <input v-model="form.city" required placeholder="Например: Москва" />
      </div>
      <div class="form-group">
        <label>Категория *</label>
        <select v-model="form.category" required>
          <option value="">Выберите категорию</option>
          <option v-for="c in categories" :key="c" :value="c">{{ c }}</option>
        </select>
      </div>
      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="success" style="color:#16a34a">Практика отправлена на модерацию!</p>
      <button type="submit" class="btn btn-primary" :disabled="loading">
        {{ loading ? 'Отправка...' : 'Отправить на модерацию' }}
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../api'

const router = useRouter()
const categories = ['Транспорт', 'Благоустройство', 'ЖКХ', 'Экология', 'Безопасность', 'Культура', 'Образование', 'Здоровье']

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
