import type { AuthUser } from "../model/types"

const isRecord = (value: unknown): value is Record<string, unknown> => {
    return typeof value === "object" && value !== null
}

export const isAuthUser = (value: unknown): value is AuthUser => {
    if (!isRecord(value)) {
        return false
    }

    return (
        typeof value.id === "number" &&
        typeof value.username === "string" &&
        typeof value.createdAt === "string"
    )
}


