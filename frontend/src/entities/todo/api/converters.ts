import type { TodoType } from '../model'
import type { CreateTodoDto } from './todo-api.contract'

import type { ModelsCreateTodoRequest, ModelsTodo } from '@/shared/api/generated'


export function toTodoType(input: ModelsTodo): TodoType {
    return {
        id: input.id!,
        value: input.value!,
        date: input.date!,
    }
}

export function toTodoTypeList(input: ModelsTodo[]): TodoType[] {
    return input.map(toTodoType)
}

export function toCreateTodoRequest(dto: CreateTodoDto): ModelsCreateTodoRequest {
    return {
        value: dto.value,
    }
}


