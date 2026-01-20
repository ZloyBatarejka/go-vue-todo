import { useTodoStore } from "./store"

export const useTodoManager = () => {
    const { todos, todosCount, fetchTodos, createTodo, deleteTodo } = useTodoStore()

    const actions = {
        createTodo: async (value: string) => {
            return createTodo({ value })
        },
        deleteTodo: async (id: number) => {
            await deleteTodo(id)
        }
    }

    return { ...actions, todos, todosCount, refresh: fetchTodos }
}