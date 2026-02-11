import { ITodoApiService, CreateTodoDto, DeleteTodoDto } from "./todo-api.contract"
import { TodoType } from "../model"
import { api } from "@shared/api"
import { toCreateTodoRequest, toTodoType, toTodoTypeList } from "./converters"

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