<template>
  <div class="card auth-card">
    <h1>Регистрация</h1>
    <form @submit.prevent="submit">
      <div class="form-group">
        <label>Email</label>
        <input v-model="email" type="email" required autocomplete="email" />
      </div>
      <div class="form-group">
        <label>Пароль</label>
        <input v-model="password" type="password" required minlength="6" autocomplete="new-password" />
      </div>
      <p v-if="error" class="error">{{ error }}</p>
      <button type="submit" class="btn btn-primary" :disabled="loading">
        {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
      </button>
    </form>
    <p style="margin-top:1rem">Уже есть аккаунт? <router-link to="/auth/login">Войти</router-link></p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref(null)

async function submit() {
  loading.value = true
  error.value = null
  try {
    await auth.register(email.value, password.value)
    router.push('/')
  } catch (e) {
    error.value = e.data?.error || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-card { max-width: 420px; margin: 0 auto; }
.auth-card h1 { margin-bottom: 1.5rem; }
</style>
