import { defineStore } from "pinia"
import { TodoType } from "./types"
import { computed, ref } from "vue"
import { todoApiService } from "../api"
import { CreateTodoDto } from "../api/todo-api.contract"

export const useTodoStore = defineStore('todo', () => {
    const todos = ref<TodoType[]>([])
    const isLoading = ref(false)
    const todosCount = computed(() => todos.value.length)

    const fetchTodos = async () => {
        isLoading.value = true
        try {
            todos.value = await todoApiService.fetchAll()
        } catch (error) {
            console.error('Failed to fetch todos', error)
        } finally {
            isLoading.value = false
        }
    }

    const createTodo = async (todo: CreateTodoDto) => {
        isLoading.value = true
        try {
            const createdTodo = await todoApiService.create(todo)
            todos.value.unshift(createdTodo)
        } catch (error) {
            console.error('Failed to create todo', error)
        } finally {
            isLoading.value = false
        }
    }

    const deleteTodo = async (id: number) => {
        isLoading.value = true
        try {
            await todoApiService.delete({ id })
            todos.value = todos.value.filter(todo => todo.id !== id)
        } catch (error) {
            console.error('Failed to delete todo', error)
        } finally {
            isLoading.value = false
        }
    }

    return { todos, todosCount, fetchTodos, createTodo, deleteTodo }
})