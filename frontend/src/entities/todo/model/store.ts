import { defineStore } from "pinia"
import { computed, ref } from "vue"
import { useAuthStore } from "@entities/auth/model"
import { todoApiService } from "../api"
import type { CreateTodoDto } from "../api/todo-api.contract"
import type { TodoType } from "./types"

export const useTodoStore = defineStore("todo", () => {  
	const todos = ref<TodoType[]>([])
	const isLoading = ref(false)
	const todosCount = computed(() => todos.value.length)

	const clearTodos = () => {
		todos.value = []
	}

	const isUnauthorizedError = (error: unknown): error is { status: number } => {
		return (
			typeof error === "object" &&
			error !== null &&
			"status" in error &&
			typeof error.status === "number" &&
			error.status === 401
		)
	}

	const handleTodoRequestError = (operation: string, error: unknown) => {
		if (isUnauthorizedError(error)) {
			void useAuthStore().logout()
			clearTodos()
			console.error(`Failed to ${operation}: unauthorized`)
			return
		}

		console.error(`Failed to ${operation}`, error)
	}

	const fetchTodos = async () => {
		isLoading.value = true
		try {
			todos.value = await todoApiService.fetchAll()
		} catch (error) {
			handleTodoRequestError("fetch todos", error)
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
			handleTodoRequestError("create todo", error)
		} finally {
			isLoading.value = false
		}
	}

	const deleteTodo = async (id: number) => {
		isLoading.value = true
		try {
			await todoApiService.delete({ id })
			todos.value = todos.value.filter((todo) => todo.id !== id)
		} catch (error) {
			handleTodoRequestError("delete todo", error)
		} finally {
			isLoading.value = false
		}
	}

	return { todos, todosCount, fetchTodos, createTodo, deleteTodo, clearTodos }
})
