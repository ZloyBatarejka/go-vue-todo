<script setup lang="ts">
import { useTodoManager } from '@entities/todo/model'
import Todo from '@entities/todo/ui/Todo.vue'
import { onMounted } from 'vue'
import styles from './TodoList.module.css'
import { AddTodo } from '@features/addTodo'

const { todos, refresh } = useTodoManager()

onMounted(() => {
    refresh()
})

</script>

<template>
    <section :class="styles.container">
        <AddTodo />
        <div :class="styles.list">
            <Todo v-for="todo in todos" :key="todo.id" :todo="todo" />
        </div>
    </section>
</template>