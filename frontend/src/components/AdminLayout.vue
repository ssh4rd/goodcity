<template>
  <div class="admin-layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <div class="sidebar-logo">
          <span class="logo-icon"></span>
          <span class="logo-name">GoodCity</span>
        </div>
        <p class="sidebar-subtitle">Панель администратора</p>
      </div>

      <nav class="sidebar-nav">
        <a href="#" class="sidebar-link">
          <span class="sidebar-link-icon">●</span>
          Обзор
        </a>
        <router-link to="/admin" class="sidebar-link" active-class="active">
          <span class="sidebar-link-icon">+</span>
          Модерация
          <span v-if="pendingCount > 0" class="sidebar-badge">{{ pendingCount }}</span>
        </router-link>
        <a href="#" class="sidebar-link">
          <span class="sidebar-link-icon">≡</span>
          Пользователи
        </a>
        <a href="#" class="sidebar-link">
          <span class="sidebar-link-icon">◇</span>
          Категории
        </a>
        <a href="#" class="sidebar-link">
          <span class="sidebar-link-icon">⊞</span>
          Аналитика
        </a>
      </nav>

      <div class="sidebar-footer">
        <div class="sidebar-user">
          <div class="sidebar-avatar">АД</div>
          <div class="sidebar-user-info">
            <div class="sidebar-user-name">Администратор</div>
            <div class="sidebar-user-email">{{ auth.user?.email || 'admin@goodcity.ru' }}</div>
          </div>
        </div>
      </div>
    </aside>

    <main class="admin-content">
      <slot />
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

defineProps({
  pendingCount: { type: Number, default: 0 }
})
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: #F1F5F9;
}

.sidebar {
  width: 240px;
  flex-shrink: 0;
  background: var(--c-navy);
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
}

.sidebar-header {
  padding: 24px 20px 16px;
  border-bottom: 1px solid rgba(255,255,255,0.08);
}
.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 4px;
}
.logo-icon {
  width: 28px;
  height: 28px;
  background: var(--c-green);
  border-radius: 7px;
  flex-shrink: 0;
}
.logo-name {
  font-size: 16px;
  font-weight: 700;
  color: white;
}
.sidebar-subtitle {
  font-size: 12px;
  color: rgba(255,255,255,0.45);
  margin-left: 38px;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow-y: auto;
}
.sidebar-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  font-size: 14px;
  color: rgba(255,255,255,0.65);
  transition: background 0.15s, color 0.15s;
  text-decoration: none;
  position: relative;
}
.sidebar-link:hover { background: rgba(255,255,255,0.07); color: rgba(255,255,255,0.9); }
.sidebar-link.active { background: rgba(255,255,255,0.12); color: white; font-weight: 500; }
.sidebar-link-icon { font-size: 12px; width: 18px; text-align: center; flex-shrink: 0; }
.sidebar-badge {
  margin-left: auto;
  background: #EF4444;
  color: white;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 7px;
  border-radius: 10px;
  line-height: 1.4;
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid rgba(255,255,255,0.08);
}
.sidebar-user {
  display: flex;
  align-items: center;
  gap: 10px;
}
.sidebar-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--c-green);
  color: white;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.sidebar-user-name {
  font-size: 13px;
  font-weight: 500;
  color: white;
}
.sidebar-user-email {
  font-size: 11px;
  color: rgba(255,255,255,0.45);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 150px;
}

.admin-content {
  margin-left: 240px;
  flex: 1;
  min-height: 100vh;
  padding: 32px;
}

@media (max-width: 768px) {
  .admin-layout { flex-direction: column; }
  .sidebar {
    position: static;
    width: 100%;
    height: auto;
    flex-direction: row;
    flex-wrap: wrap;
    padding: 12px 16px;
  }
  .sidebar-header { padding: 0; border-bottom: none; margin-right: auto; }
  .sidebar-subtitle { display: none; }
  .sidebar-nav {
    flex-direction: row;
    flex-wrap: nowrap;
    overflow-x: auto;
    padding: 8px 0 0;
    gap: 4px;
    width: 100%;
    border-top: 1px solid rgba(255,255,255,0.08);
    margin-top: 8px;
    scrollbar-width: none;
  }
  .sidebar-nav::-webkit-scrollbar { display: none; }
  .sidebar-link { white-space: nowrap; padding: 7px 12px; font-size: 13px; }
  .sidebar-footer { display: none; }
  .admin-content { margin-left: 0; padding: 20px 16px; }
}
</style>
