import { TodoType } from "../model"

export type CreateTodoDto = Pick<TodoType, 'value'>
export type DeleteTodoDto = Pick<TodoType, 'id'>

export interface ITodoApiService {
    fetchAll(): Promise<TodoType[]>
    create(dto: CreateTodoDto): Promise<TodoType>
    delete(dto: DeleteTodoDto): Promise<void>
}

