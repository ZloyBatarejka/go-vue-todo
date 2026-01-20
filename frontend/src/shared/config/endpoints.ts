export const TodoEndpoints = {
    getAll: '/todos',
    create: '/todos',
    getById: (id: number) => `/todos/${id}`,
    delete: (id: number) => `/todos/${id}`,
} as const



