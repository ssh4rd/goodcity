<template>
  <div class="auth-page">
    <!-- Left panel -->
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
        <h1 class="form-title">Добро пожаловать</h1>
        <p class="form-sub">Войдите в аккаунт, чтобы добавлять и оценивать практики</p>

        <form class="auth-form" @submit.prevent="submit">
          <div class="form-field">
            <label class="form-label">Email</label>
            <input v-model="email" type="email" class="form-input" placeholder="your@email.ru" required autocomplete="email" />
          </div>
          <div class="form-field">
            <label class="form-label">Пароль</label>
            <input v-model="password" type="password" class="form-input" placeholder="••••••••" required autocomplete="current-password" />
          </div>

          <div class="form-row">
            <label class="checkbox-label">
              <input type="checkbox" class="checkbox" />
              <span>Запомнить меня</span>
            </label>
            <a href="#" class="forgot-link">Забыли пароль?</a>
          </div>

          <p v-if="error" class="error-msg">{{ error }}</p>

          <button type="submit" class="btn-submit" :disabled="loading">
            {{ loading ? 'Вход...' : 'Войти в аккаунт' }}
          </button>

          <p class="switch-auth">
            Нет аккаунта?
            <router-link to="/auth/register" class="switch-link">Зарегистрироваться</router-link>
          </p>
        </form>
      </div>
    </div>
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
    await auth.login(email.value, password.value)
    router.push('/')
  } catch (e) {
    error.value = e.data?.error || 'Неверный email или пароль'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  display: flex;
  min-height: 100vh;
}

/* Left */
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

.auth-brand { margin-bottom: 40px; }
.auth-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
}
.logo-icon {
  width: 32px;
  height: 32px;
  background: var(--c-green);
  border-radius: 8px;
}
.logo-name { font-size: 17px; font-weight: 700; color: white; }
.auth-tagline {
  font-size: 24px;
  font-weight: 700;
  color: white;
  line-height: 1.3;
}

.auth-features {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-bottom: 32px;
  position: relative;
  z-index: 1;
}
.auth-feature {
  background: rgba(255,255,255,0.1);
  border-radius: 8px;
  padding: 12px 16px;
  font-size: 14px;
  color: rgba(255,255,255,0.85);
  border-left: 3px solid var(--c-green);
}

.auth-float-card {
  background: rgba(255,255,255,0.12);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 14px;
  padding: 18px 20px;
  position: relative;
  z-index: 1;
  margin-top: auto;
}
.float-card-score { font-size: 24px; font-weight: 700; color: white; }
.float-card-sub { font-size: 13px; color: rgba(255,255,255,0.65); margin: 4px 0; }
.float-card-note { font-size: 12px; color: rgba(255,255,255,0.5); }

.auth-circles { position: absolute; inset: 0; pointer-events: none; }
.ac {
  position: absolute;
  border-radius: 50%;
  opacity: 0.15;
}
.ac-1 {
  width: 220px; height: 220px;
  background: var(--c-green);
  bottom: 80px; right: -60px;
}
.ac-2 {
  width: 180px; height: 160px;
  background: #4B9EFF;
  bottom: 0; right: 40px;
}

/* Right */
.auth-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  padding: 40px;
}
.auth-form-wrap { width: 100%; max-width: 420px; }

.form-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--c-text);
  margin-bottom: 8px;
}
.form-sub {
  font-size: 14px;
  color: var(--c-muted);
  margin-bottom: 32px;
}
.auth-form { display: flex; flex-direction: column; gap: 0; }
.auth-form .form-field { margin-bottom: 16px; }

.form-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--c-muted);
  cursor: pointer;
}
.checkbox { accent-color: var(--c-green); width: 16px; height: 16px; }
.forgot-link { font-size: 14px; color: var(--c-green); text-decoration: none; }
.forgot-link:hover { text-decoration: underline; }

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
  margin-bottom: 20px;
  transition: background 0.15s;
}
.btn-submit:hover:not(:disabled) { background: var(--c-green-dark); }
.btn-submit:disabled { opacity: 0.55; cursor: not-allowed; }

.divider {
  text-align: center;
  position: relative;
  margin-bottom: 16px;
  color: var(--c-muted);
  font-size: 13px;
}
.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--c-border);
}
.divider span {
  background: white;
  padding: 0 12px;
  position: relative;
}

.switch-auth {
  text-align: center;
  font-size: 14px;
  color: var(--c-muted);
}
.switch-link { color: var(--c-green); font-weight: 500; text-decoration: none; }
.switch-link:hover { text-decoration: underline; }

.error-msg { color: var(--c-red-dark); font-size: 13px; margin-bottom: 12px; }

@media (max-width: 768px) {
  .auth-page { flex-direction: column; }

  .auth-left {
    width: 100%;
    padding: 32px 24px 28px;
    flex-direction: column;
    min-height: auto;
    align-items: center;
    text-align: center;
  }
  .auth-brand { margin-bottom: 0; }
  .auth-tagline { font-size: 18px; }
  .auth-features { display: none; }
  .auth-float-card { display: none; }

  .auth-right {
    flex: none;
    padding: 28px 20px 40px;
    align-items: flex-start;
    width: 100%;
  }
  .auth-form-wrap { max-width: 100%; }
  .form-title { font-size: 22px; }
  .form-sub { margin-bottom: 20px; }
}
</style>
