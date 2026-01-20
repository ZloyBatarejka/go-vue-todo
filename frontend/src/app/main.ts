import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import '../shared/assets/main.css'
import { createPinia } from 'pinia'

createApp(App).use(createPinia()).use(router).mount('#app')

