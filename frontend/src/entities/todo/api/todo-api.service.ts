import { TodoEndpoints } from "@/shared/config"
import { ITodoApiService, CreateTodoDto, DeleteTodoDto } from "./todo-api.contract"
import { httpClient, ApiClient } from "@shared/api"
import { TodoType } from "../model"

export class TodoApiService implements ITodoApiService {
    constructor(private readonly httpClient: ApiClient) { }

    fetchAll(): Promise<TodoType[]> {
        return this.httpClient.get<TodoType[]>(TodoEndpoints.getAll)
    }

    create(dto: CreateTodoDto): Promise<TodoType> {
        return this.httpClient.post<TodoType>(TodoEndpoints.create, dto)
    }

    delete(dto: DeleteTodoDto): Promise<void> {
        return this.httpClient.delete<void>(TodoEndpoints.delete(dto.id))
    }
}

export const todoApiService = new TodoApiService(httpClient)