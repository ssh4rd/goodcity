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
          <router-link v-if="auth.isModerator" to="/admin" class="nav-link">Админ</router-link>
          <div class="user-menu" @click.stop="toggleDropdown" ref="menuRef">
            <div class="avatar">{{ initials }}</div>
            <div v-if="open" class="dropdown">
              <router-link to="/profile" class="dropdown-item" @click="open = false">Профиль</router-link>
              <button class="dropdown-item dropdown-item--danger" @click="handleLogout">Выйти</button>
            </div>
          </div>
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()
const open = ref(false)
const menuRef = ref(null)

const initials = computed(() => {
  const name = auth.user?.name || auth.user?.email || ''
  return name.split(' ').map(w => w[0]).slice(0, 2).join('').toUpperCase()
})

function toggleDropdown() { open.value = !open.value }

function handleClickOutside(e) {
  if (menuRef.value && !menuRef.value.contains(e.target)) open.value = false
}

onMounted(() => document.addEventListener('click', handleClickOutside))
onUnmounted(() => document.removeEventListener('click', handleClickOutside))

function handleLogout() {
  open.value = false
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

.user-menu { position: relative; cursor: pointer; }
.avatar {
  width: 36px; height: 36px;
  background: var(--c-green); color: white;
  border-radius: 50%; font-size: 13px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
  user-select: none;
}
.dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  background: white;
  border: 1px solid var(--c-border);
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
  min-width: 160px;
  overflow: hidden;
  z-index: 200;
}
.dropdown-item {
  display: block; width: 100%;
  padding: 11px 16px;
  font-size: 14px; color: var(--c-text);
  text-align: left; background: none; border: none;
  cursor: pointer; font-family: inherit;
  transition: background 0.1s;
  text-decoration: none;
}
.dropdown-item:hover { background: #F8FAFB; }
.dropdown-item--danger { color: #DC2626; }
.dropdown-item--danger:hover { background: #FEF2F2; }

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
