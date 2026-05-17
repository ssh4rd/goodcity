import { createRouter, createMemoryHistory, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('./pages/index.vue') },
  { path: '/practices/:id', component: () => import('./pages/practices/[id].vue') },
  { path: '/submit', component: () => import('./pages/submit.vue'), meta: { requiresAuth: true } },
  { path: '/auth/login', component: () => import('./pages/auth/login.vue') },
  { path: '/auth/register', component: () => import('./pages/auth/register.vue') },
  { path: '/admin', component: () => import('./pages/admin/index.vue'), meta: { requiresModerator: true } },
  { path: '/admin/practice/:id', component: () => import('./pages/admin/practice/[id].vue'), meta: { requiresModerator: true } },
]

export function createAppRouter() {
  const isServer = typeof window === 'undefined'
  const history = isServer ? createMemoryHistory() : createWebHistory()
  return createRouter({ history, routes })
}
