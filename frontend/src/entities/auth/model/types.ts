export type AuthStatus = "authenticated" | "unauthenticated"

export type AuthUser = {
    id: number
    username: string
    createdAt: string
}

export type AuthSession = {
    token: string
    user: AuthUser
}

