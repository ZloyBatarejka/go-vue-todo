import { api } from "@shared/api"
import type { TodoType } from "../model"
import { toCreateTodoRequest, toTodoType, toTodoTypeList } from "./converters"
import type {
	CreateTodoDto,
	DeleteTodoDto,
	ITodoApiService,
} from "./todo-api.contract"

export class TodoApiService implements ITodoApiService {
	fetchAll(): Promise<TodoType[]> {
		return api.todos.todosList().then(toTodoTypeList)
	}

	create(dto: CreateTodoDto): Promise<TodoType> {
		return api.todos.todosCreate(toCreateTodoRequest(dto)).then(toTodoType)
	}

	delete(dto: DeleteTodoDto): Promise<void> { 
		return api.todos.todosDelete(dto.id)
	}
}  

export const todoApiService = new TodoApiService()
 