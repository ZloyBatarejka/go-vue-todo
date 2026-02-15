import { api } from "@shared/api"
import type { ModelsAuthResponse } from "@shared/api/generated"
import { isAuthUser } from "../lib/guards"
import type { AuthSession } from "../model"
import type { IAuthApiService, LoginDto, RegisterDto } from "./auth-api.contract"

const toAuthSession = (response: ModelsAuthResponse): AuthSession => {
    if (!response.accessToken || !response.user) {
        throw new Error("Invalid auth response shape")
    }

    if (!isAuthUser(response.user)) {
        throw new Error("Invalid auth user response shape")
    }

    return {
        token: response.accessToken,
        user: response.user,
    }
}

export class AuthApiService implements IAuthApiService {
    login(dto: LoginDto): Promise<AuthSession> {
        return api.auth.loginCreate(dto, {
            secure: false,
        })
            .then(toAuthSession)
    }

    register(dto: RegisterDto): Promise<AuthSession> {
        return api.auth.registerCreate(dto, {
            secure: false,
        })
            .then(toAuthSession)
    }

    refresh(): Promise<AuthSession> {
        return api.auth.refreshCreate({
            secure: false,
        })
            .then(toAuthSession)
    }

    logout(): Promise<void> {
        return api.auth.logoutCreate({
            secure: false,
        })
    }
}

export const authApiService = new AuthApiService()


