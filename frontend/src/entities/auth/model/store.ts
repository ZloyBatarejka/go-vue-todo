import { defineStore } from "pinia"
import { computed, ref } from "vue"
import { setApiAuthToken } from "@/shared/api"
import { authApiService } from "../api"
import type { AuthSession, AuthStatus, AuthUser } from "./types"

export const useAuthStore = defineStore("auth", () => {
    const token = ref<string | null>(null)
    const user = ref<AuthUser | null>(null)
    const authStatus = ref<AuthStatus>("unauthenticated")

    const isAuthenticated = computed(() => authStatus.value === "authenticated")

    const initAuth = async () => {
        try {
            const session = await authApiService.refresh()
            setSession(session)
        } catch {
            clearSession()
        }
    }

    const setSession = (session: AuthSession) => {
        token.value = session.token
        user.value = session.user
        authStatus.value = "authenticated"
        setApiAuthToken(session.token)
    }

    const clearSession = () => {
        token.value = null
        user.value = null
        authStatus.value = "unauthenticated"
        setApiAuthToken(null)
    }

    const setUser = (nextUser: AuthUser | null) => {
        user.value = nextUser
    }

    const login = async (username: string, password: string) => {
        const session = await authApiService.login({ username, password })
        setSession(session)
    }

    const register = async (username: string, password: string) => {
        const session = await authApiService.register({ username, password })
        setSession(session)
    }

    const logout = async () => {
        try {
            await authApiService.logout()
        } catch (error) {
            console.error("Failed to logout", error)
        } finally {
            clearSession()
        }
    }

    return {
        token,
        user,
        authStatus,
        isAuthenticated,
        initAuth,
        setSession,
        clearSession,
        setUser,
        login,
        register,
        logout
    }
})



