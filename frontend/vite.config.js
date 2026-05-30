import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  define: {
    __DADATA_TOKEN__: JSON.stringify(process.env.VITE_DADATA_TOKEN || ''),
  },
  server: {
    host: true,
    proxy: {
      '/api': {
        target: process.env.API_URL || 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
