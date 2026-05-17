<template>
  <div id="app">
    <!-- Auth pages: full-page, self-contained -->
    <router-view v-if="layout === 'auth'" />

    <!-- Admin pages: self-contained with sidebar -->
    <router-view v-else-if="layout === 'admin'" />

    <!-- Public pages: with AppHeader -->
    <template v-else>
      <AppHeader />
      <div class="page-content">
        <router-view />
      </div>
      <BottomTabBar />
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from './components/AppHeader.vue'
import BottomTabBar from './components/BottomTabBar.vue'

const route = useRoute()
const layout = computed(() => route.meta.layout || 'default')
</script>

<style>
:root {
  --c-navy: #1B2C4B;
  --c-navy-2: #243454;
  --c-green: #2ECC87;
  --c-green-dark: #22A06B;
  --c-blue: #2563EB;
  --c-amber: #F59E0B;
  --c-red: #EF4444;
  --c-red-dark: #DC2626;
  --c-text: #111827;
  --c-muted: #6B7280;
  --c-border: #E5E7EB;
  --c-bg: #F8FAFB;
  --c-white: #FFFFFF;
  --shadow-sm: 0 1px 3px rgba(0,0,0,0.08), 0 1px 2px rgba(0,0,0,0.04);
  --shadow-md: 0 4px 12px rgba(0,0,0,0.08);
  --radius: 12px;
  --radius-sm: 8px;
}

*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

html { font-size: 16px; }

body {
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
  background: var(--c-bg);
  color: var(--c-text);
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
}

a { color: inherit; text-decoration: none; }

button { font-family: inherit; cursor: pointer; border: none; }

.container {
  max-width: 1248px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Shared button styles */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 500;
  line-height: 1;
  cursor: pointer;
  transition: opacity 0.15s, background 0.15s;
  border: none;
  text-decoration: none;
}
.btn:disabled { opacity: 0.55; cursor: not-allowed; }
.btn-primary {
  background: var(--c-green);
  color: var(--c-white);
}
.btn-primary:hover:not(:disabled) { background: var(--c-green-dark); }
.btn-outline {
  background: transparent;
  color: var(--c-text);
  border: 1.5px solid var(--c-border);
}
.btn-outline:hover:not(:disabled) { border-color: #9CA3AF; }
.btn-danger {
  background: #FEE2E2;
  color: var(--c-red-dark);
  border: 1px solid #FECACA;
}
.btn-danger:hover:not(:disabled) { background: #FECACA; }
.btn-success {
  background: #D1FAE5;
  color: #065F46;
  border: 1px solid #A7F3D0;
}
.btn-success:hover:not(:disabled) { background: #A7F3D0; }
.btn-sm {
  padding: 6px 14px;
  font-size: 13px;
}

/* Shared card */
.card {
  background: var(--c-white);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--c-border);
}

/* Form elements */
.form-field { margin-bottom: 16px; }
.form-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--c-text);
  margin-bottom: 6px;
}
.form-label.required::after { content: ' *'; color: var(--c-red); }
.form-input {
  width: 100%;
  padding: 10px 14px;
  border: 1.5px solid var(--c-border);
  border-radius: var(--radius-sm);
  font-size: 15px;
  font-family: inherit;
  color: var(--c-text);
  background: var(--c-white);
  transition: border-color 0.15s;
  outline: none;
}
.form-input:focus { border-color: var(--c-green); }
.form-input::placeholder { color: #9CA3AF; }
.form-textarea {
  resize: vertical;
  min-height: 100px;
}
.form-select {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%236B7280' stroke-width='2'%3E%3Cpath d='M6 9l6 6 6-6'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 36px;
}

@media (max-width: 768px) {
  .page-content { padding-bottom: 64px; }
}

.error-msg {
  color: var(--c-red-dark);
  font-size: 13px;
  margin-top: 4px;
}

/* Stars */
.stars { color: var(--c-amber); letter-spacing: 1px; }
.star-score { font-weight: 600; font-size: 14px; margin-left: 4px; }

/* Badge */
.badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}
.badge-pending { background: #FEF3C7; color: #92400E; }
.badge-approved { background: #D1FAE5; color: #065F46; }
.badge-rejected { background: #FEE2E2; color: #991B1B; }

/* Rating bars */
.criteria-bar { display: flex; align-items: center; gap: 8px; font-size: 12px; color: var(--c-muted); margin-bottom: 4px; }
.criteria-bar-label { width: 52px; flex-shrink: 0; }
.criteria-bar-track { flex: 1; height: 4px; background: var(--c-border); border-radius: 2px; overflow: hidden; }
.criteria-bar-fill { height: 100%; border-radius: 2px; transition: width 0.3s; }
.criteria-bar-fill.blue { background: var(--c-blue); }
.criteria-bar-fill.green { background: #22C55E; }
.criteria-bar-fill.teal { background: #10B981; }
.criteria-bar-pct { width: 32px; text-align: right; flex-shrink: 0; }
</style>
