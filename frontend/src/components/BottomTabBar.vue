<template>
  <nav class="bottom-tab-bar">
    <router-link to="/" class="tab-item" exact-active-class="active">
      <span class="tab-icon">
        <svg width="22" height="22" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
          <path d="M3 12L12 3l9 9M4 10v10a1 1 0 001 1h5v-6h4v6h5a1 1 0 001-1V10"/>
        </svg>
      </span>
      <span class="tab-label">Главная</span>
    </router-link>

    <router-link to="/catalog" class="tab-item" active-class="active">
      <span class="tab-icon">
        <svg width="22" height="22" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
          <rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/>
          <rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/>
        </svg>
      </span>
      <span class="tab-label">Каталог</span>
    </router-link>

    <router-link to="/submit" class="tab-add">
      <svg width="24" height="24" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
        <path d="M12 5v14M5 12h14"/>
      </svg>
    </router-link>

    <a href="#" class="tab-item">
      <span class="tab-icon">
        <svg width="22" height="22" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
          <path d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4"/>
        </svg>
      </span>
      <span class="tab-label">Карта</span>
    </a>

    <router-link v-if="auth.user" to="/admin" class="tab-item" active-class="active">
      <span class="tab-icon">
        <svg width="22" height="22" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
          <path d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
        </svg>
      </span>
      <span class="tab-label">Профиль</span>
    </router-link>
    <router-link v-else to="/auth/login" class="tab-item" active-class="active">
      <span class="tab-icon">
        <svg width="22" height="22" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
          <path d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
        </svg>
      </span>
      <span class="tab-label">Профиль</span>
    </router-link>
  </nav>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'
const auth = useAuthStore()
</script>

<style scoped>
.bottom-tab-bar {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: white;
  border-top: 1px solid var(--c-border);
  z-index: 200;
  align-items: center;
  justify-content: space-around;
  padding: 0 4px;
  box-shadow: 0 -2px 12px rgba(0,0,0,0.06);
}

@media (max-width: 768px) {
  .bottom-tab-bar { display: flex; }
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  padding: 6px 12px;
  color: #9CA3AF;
  text-decoration: none;
  flex: 1;
  transition: color 0.15s;
}
.tab-item.active { color: var(--c-green); }
.tab-item:hover { color: var(--c-text); }

.tab-icon { display: flex; }
.tab-label { font-size: 10px; font-weight: 500; }

.tab-add {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  background: var(--c-green);
  border-radius: 50%;
  color: white;
  flex-shrink: 0;
  margin-bottom: 8px;
  box-shadow: 0 4px 12px rgba(46, 204, 135, 0.4);
  transition: background 0.15s;
}
.tab-add:hover { background: var(--c-green-dark); }
</style>
