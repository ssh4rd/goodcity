import { createRouter, createMemoryHistory, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('./pages/index.vue') },
  { path: '/catalog', component: () => import('./pages/catalog.vue') },
  { path: '/practices/:id', component: () => import('./pages/practices/[id].vue') },
  { path: '/submit', component: () => import('./pages/submit.vue'), meta: { requiresAuth: true } },
  { path: '/auth/login', component: () => import('./pages/auth/login.vue'), meta: { layout: 'auth' } },
  { path: '/auth/register', component: () => import('./pages/auth/register.vue'), meta: { layout: 'auth' } },
  { path: '/profile', component: () => import('./pages/profile.vue'), meta: { requiresAuth: true } },
  { path: '/admin', component: () => import('./pages/admin/index.vue'), meta: { layout: 'admin', requiresModerator: true } },
  { path: '/admin/practice/:id', component: () => import('./pages/admin/practice/[id].vue'), meta: { layout: 'admin', requiresModerator: true } },
]

export function createAppRouter() {
  const isServer = typeof window === 'undefined'
  const history = isServer ? createMemoryHistory() : createWebHistory()
  return createRouter({ history, routes })
}
