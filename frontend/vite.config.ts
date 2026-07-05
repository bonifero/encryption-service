import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const proxyTarget = env.VITE_API_PROXY_TARGET || 'http://localhost:8080'

  return {
    plugins: [react()],
    server: {
      proxy: {
        '/api': proxyTarget,
        '/v3/api-docs': proxyTarget,
        '/swagger-ui.html': proxyTarget,
        '/swagger-ui': proxyTarget,
        '/webjars': proxyTarget,
      },
    },
  }
})
