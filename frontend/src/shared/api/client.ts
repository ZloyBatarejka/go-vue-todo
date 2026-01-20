const API_BASE_URL = 'http://localhost:8080/api'

type RequestOptions = {
    method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
    headers?: Record<string, string>
    body?: unknown
}


export class ApiClient {
    private baseUrl: string

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl
    }

    async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
        const { method = 'GET', headers = {}, body } = options

        const config: RequestInit = {
            method,
            headers: {
                'Content-Type': 'application/json',
                ...headers,
            },
            body: body ? JSON.stringify(body) : undefined,
        }

        const response = await fetch(`${this.baseUrl}${endpoint}`, config)
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`)
        }

        return response.json()

    }

    async get<T>(endpoint: string): Promise<T> {
        return this.request<T>(endpoint, { method: 'GET' })
    }

    async post<T>(endpoint: string, body: unknown): Promise<T> {
        return this.request<T>(endpoint, { method: 'POST', body })
    }

    async delete<T>(endpoint: string): Promise<T> {
        return this.request<T>(endpoint, { method: 'DELETE' })
    }

}


export const httpClient = new ApiClient(API_BASE_URL)