import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vitejs.dev/config/
export default defineConfig(( { mode } ) => {
  const env = loadEnv(mode, `${process.cwd()}/..`, '')

  return {
    plugins: [vue()],
    define: {
      serverHost: JSON.stringify(env.SERVER_HOST),
      serverPort: JSON.stringify(env.SERVER_PORT)
    },
    build: {
      outDir: 'dist',
      assetsDir: 'static',
    }
  }
})
