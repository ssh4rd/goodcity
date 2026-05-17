<template>
  <div id="app">
    <nav class="navbar">
      <router-link to="/" class="navbar-brand">GoodCity</router-link>
      <div class="navbar-links">
        <router-link to="/">Каталог</router-link>
        <template v-if="auth.user">
          <router-link to="/submit">Добавить практику</router-link>
          <router-link v-if="auth.isModerator" to="/admin">Модерация</router-link>
          <button @click="auth.logout">Выйти ({{ auth.user.email }})</button>
        </template>
        <template v-else>
          <router-link to="/auth/login">Войти</router-link>
          <router-link to="/auth/register">Регистрация</router-link>
        </template>
      </div>
    </nav>
    <main class="container">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from './stores/auth'
const auth = useAuthStore()
</script>

<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: system-ui, sans-serif; background: #f5f5f5; color: #333; }

.navbar {
  background: #2563eb;
  color: white;
  padding: 0.75rem 1.5rem;
  display: flex;
  align-items: center;
  gap: 2rem;
}
.navbar-brand { color: white; text-decoration: none; font-size: 1.25rem; font-weight: bold; }
.navbar-links { display: flex; gap: 1rem; align-items: center; margin-left: auto; }
.navbar-links a { color: white; text-decoration: none; opacity: 0.9; }
.navbar-links a:hover { opacity: 1; }
.navbar-links button { background: rgba(255,255,255,0.2); border: none; color: white; cursor: pointer; padding: 0.25rem 0.75rem; border-radius: 4px; }

.container { max-width: 1100px; margin: 0 auto; padding: 2rem 1rem; }

.card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.btn {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-size: 0.9rem;
  text-decoration: none;
}
.btn-primary { background: #2563eb; color: white; }
.btn-primary:hover { background: #1d4ed8; }
.btn-danger { background: #dc2626; color: white; }
.btn-success { background: #16a34a; color: white; }
.btn-secondary { background: #e5e7eb; color: #333; }

.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.25rem; font-weight: 500; }
.form-group input, .form-group textarea, .form-group select {
  width: 100%; padding: 0.5rem; border: 1px solid #d1d5db; border-radius: 4px;
  font-size: 1rem;
}
.form-group textarea { height: 120px; resize: vertical; }

.badge {
  display: inline-block;
  padding: 0.15rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}
.badge-pending { background: #fef3c7; color: #92400e; }
.badge-approved { background: #d1fae5; color: #065f46; }
.badge-rejected { background: #fee2e2; color: #991b1b; }

.error { color: #dc2626; font-size: 0.9rem; margin-top: 0.5rem; }
</style>
