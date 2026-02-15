import { defineStore } from "pinia"
import { computed, ref } from "vue"
import { setApiAuthToken } from "@/shared/api"
import { isAuthUser } from "../lib/guards"
import { authApiService } from "../api"
import type { AuthSession, AuthStatus, AuthUser } from "./types"

const TOKEN_STORAGE_KEY = "goTodo.auth.token"
const USER_STORAGE_KEY = "goTodo.auth.user"

const readTokenFromStorage = (): string | null => {
    const token = localStorage.getItem(TOKEN_STORAGE_KEY)
    return token && token.trim() !== "" ? token : null
}

const readUserFromStorage = (): AuthUser | null => {
    const raw = localStorage.getItem(USER_STORAGE_KEY)
    if (!raw) {
        return null
    }

    try {
        const parsed: unknown = JSON.parse(raw)
        if (isAuthUser(parsed)) {
            return parsed
        }

        console.error("Stored auth user has invalid shape")
        localStorage.removeItem(USER_STORAGE_KEY)
        return null
    } catch (error) {
        console.error("Failed to parse stored auth user", error)
        localStorage.removeItem(USER_STORAGE_KEY)
        return null
    }
}

const writeSessionToStorage = (session: AuthSession) => {
    localStorage.setItem(TOKEN_STORAGE_KEY, session.token)
    localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(session.user))
}

const clearSessionFromStorage = () => {
    localStorage.removeItem(TOKEN_STORAGE_KEY)
    localStorage.removeItem(USER_STORAGE_KEY)
}

export const useAuthStore = defineStore("auth", () => {
    const token = ref<string | null>(null)
    const user = ref<AuthUser | null>(null)
    const authStatus = ref<AuthStatus>("unauthenticated")

    const isAuthenticated = computed(() => authStatus.value === "authenticated")

    const initAuth = () => {
        token.value = readTokenFromStorage()
        user.value = readUserFromStorage()
        authStatus.value = token.value ? "authenticated" : "unauthenticated"
        setApiAuthToken(token.value)
    }

    const setSession = (session: AuthSession) => {
        token.value = session.token
        user.value = session.user
        authStatus.value = "authenticated"
        setApiAuthToken(session.token)
        writeSessionToStorage(session)
    }

    const setUser = (nextUser: AuthUser | null) => {
        user.value = nextUser
        if (nextUser) {
            localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(nextUser))
            return
        }

        localStorage.removeItem(USER_STORAGE_KEY)
    }

    const login = async (username: string, password: string) => {
        const session = await authApiService.login({ username, password })
        setSession(session)
    }

    const register = async (username: string, password: string) => {
        const session = await authApiService.register({ username, password })
        setSession(session)
    }

    const logout = () => {
        token.value = null
        user.value = null
        authStatus.value = "unauthenticated"
        setApiAuthToken(null)
        clearSessionFromStorage()
    }

    return {
        token,
        user,
        authStatus,
        isAuthenticated,
        initAuth,
        setSession,
        setUser,
        login,
        register,
        logout
    }
})



