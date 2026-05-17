<template>
  <header class="header">
    <div class="header-inner container">
      <router-link to="/" class="logo">
        <span class="logo-icon"></span>
        <span class="logo-name">GoodCity</span>
      </router-link>

      <nav class="nav">
        <router-link to="/catalog" class="nav-link">Каталог</router-link>
        <a href="#" class="nav-link">Карта</a>
        <a href="#" class="nav-link">Рейтинг</a>
        <a href="#" class="nav-link">О проекте</a>
      </nav>

      <div class="actions">
        <template v-if="auth.user">
          <button class="btn-login" @click="handleLogout">Выйти</button>
          <router-link v-if="auth.isModerator" to="/admin" class="nav-link" style="margin-right:8px">Админ</router-link>
        </template>
        <template v-else>
          <router-link to="/auth/login" class="btn-login">Войти</router-link>
        </template>
        <router-link to="/submit" class="btn-add">Добавить +</router-link>
      </div>
    </div>
  </header>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

function handleLogout() {
  auth.logout()
  router.push('/')
}
</script>

<style scoped>
.header {
  background: var(--c-navy);
  position: sticky;
  top: 0;
  z-index: 100;
}
.header-inner {
  display: flex;
  align-items: center;
  height: 72px;
  gap: 40px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}
.logo-icon {
  width: 32px;
  height: 32px;
  background: var(--c-green);
  border-radius: 8px;
  display: block;
}
.logo-name {
  font-size: 17px;
  font-weight: 700;
  color: white;
}

.nav {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}
.nav-link {
  color: rgba(255,255,255,0.80);
  font-size: 15px;
  padding: 6px 14px;
  border-radius: 6px;
  transition: color 0.15s, background 0.15s;
}
.nav-link:hover,
.nav-link.router-link-active {
  color: white;
  background: rgba(255,255,255,0.08);
}

.actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}
.btn-login {
  color: rgba(255,255,255,0.85);
  font-size: 15px;
  padding: 8px 20px;
  border-radius: 8px;
  border: 1.5px solid rgba(255,255,255,0.25);
  background: transparent;
  cursor: pointer;
  font-family: inherit;
  transition: border-color 0.15s, color 0.15s;
}
.btn-login:hover { border-color: rgba(255,255,255,0.55); color: white; }
.btn-add {
  background: var(--c-green);
  color: white;
  font-size: 15px;
  font-weight: 600;
  padding: 9px 20px;
  border-radius: 8px;
  display: inline-block;
  transition: background 0.15s;
}
.btn-add:hover { background: var(--c-green-dark); }

@media (max-width: 768px) {
  .header-inner { gap: 0; padding: 0 16px; }
  .nav { display: none; }
  .actions { gap: 6px; }
  .btn-login { display: none; }
  .btn-add {
    font-size: 13px;
    padding: 7px 14px;
  }
}
</style>
