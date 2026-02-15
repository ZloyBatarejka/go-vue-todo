import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import '../shared/assets/main.css'
import { createPinia } from 'pinia'
import { useAuthStore } from '@entities/auth/model'

const bootstrap = async () => {
    const app = createApp(App)
    const pinia = createPinia()

    app.use(pinia)
    await useAuthStore(pinia).initAuth()

    app.use(router).mount('#app')
}

void bootstrap()

