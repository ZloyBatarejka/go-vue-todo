import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@entities/auth/model'
import HomePage from '@/pages/home/ui/HomePage.vue'
import AuthPage from '@/pages/auth/ui/AuthPage.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomePage,
            meta: { requiresAuth: true },
        },
        {
            path: '/auth',
            name: 'auth',
            component: AuthPage,
            meta: { guestOnly: true },
        },
    ]
})

router.beforeEach((to) => {
    const authStore = useAuthStore()
    const isAuthenticated = authStore.isAuthenticated

    if (to.meta.requiresAuth && !isAuthenticated) {
        return { name: 'auth' }
    }

    if (to.meta.guestOnly && isAuthenticated) {
        return { name: 'home' }
    }

    return true
})

export default router

