import { storeToRefs } from "pinia"
import { useTodoStore } from "./store"

export const useTodoManager = () => {
    const store = useTodoStore()
    const {todos, todosCount} = storeToRefs(store)

    const actions = {
        createTodo: async (value: string) => {
            return store.createTodo({ value })
        },
        deleteTodo: async (id: number) => {
            await store.deleteTodo(id)
        },
        clearTodos: () => {
            store.clearTodos()
        }
    }

    return { ...actions, todos, todosCount, refresh: store.fetchTodos }
}