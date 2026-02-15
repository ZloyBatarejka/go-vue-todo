import { storeToRefs } from "pinia"
import type { AuthSession, AuthUser } from "./types"
import { useAuthStore } from "./store"

export const useAuthManager = () => {
    const store = useAuthStore()
    const { token, user, authStatus, isAuthenticated } = storeToRefs(store)

    const actions = {
        initAuth: () => {
            store.initAuth()
        },
        setSession: (session: AuthSession) => {
            store.setSession(session)
        },
        setUser: (nextUser: AuthUser | null) => {
            store.setUser(nextUser)
        },
        login: async (username: string, password: string) => {
            await store.login(username, password)
        },
        register: async (username: string, password: string) => {
            await store.register(username, password)
        },
        logout: () => {
            store.logout()
        }
    }

    return { token, user, authStatus, isAuthenticated, ...actions }
}



