import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    // recharts alone is ~400 kB; acceptable for an internal dashboard
    chunkSizeWarningLimit: 1000,
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8002',
    },
  },
})
