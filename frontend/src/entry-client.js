import { createApp } from './app'
import { useAuthStore } from './stores/auth'

const { app, pinia, router } = createApp()

// Restore auth from localStorage before first navigation
const auth = useAuthStore(pinia)
auth.init()

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !auth.user) {
    return '/auth/login'
  }
  if (to.meta.requiresModerator && !auth.isModerator) {
    return '/'
  }
})

router.isReady().then(() => {
  app.mount('#app')
})
