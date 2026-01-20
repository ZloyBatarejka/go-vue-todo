import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'url'

export default defineConfig({
    plugins: [vue()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
            '@src': fileURLToPath(new URL('./src', import.meta.url)),
            '@app': fileURLToPath(new URL('./src/app', import.meta.url)),
            '@entities': fileURLToPath(new URL('./src/entities', import.meta.url)),
            '@features': fileURLToPath(new URL('./src/features', import.meta.url)),
            '@pages': fileURLToPath(new URL('./src/pages', import.meta.url)),
            '@shared': fileURLToPath(new URL('./src/shared', import.meta.url)),
            '@widgets': fileURLToPath(new URL('./src/widgets', import.meta.url)),
        }
    }
})

