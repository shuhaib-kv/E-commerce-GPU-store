import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 3000,
    proxy: {
      '/user': 'http://localhost:8080',
      '/admin': 'http://localhost:8080',
      '/cart': 'http://localhost:8080',
      '/razorpay': 'http://localhost:8080',
      '/payment-success': 'http://localhost:8080',
      '/success': 'http://localhost:8080',
      '/public': 'http://localhost:8080',
    },
  },
})
