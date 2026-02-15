import type { AuthSession } from "../model"

export type LoginDto = {
    username: string
    password: string
}

export type RegisterDto = {
    username: string
    password: string
}

export interface IAuthApiService {
    login(dto: LoginDto): Promise<AuthSession>
    register(dto: RegisterDto): Promise<AuthSession>
    refresh(): Promise<AuthSession>
    logout(): Promise<void>
}


