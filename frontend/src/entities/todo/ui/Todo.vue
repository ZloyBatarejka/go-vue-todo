<script setup lang="ts">
import { useTodoManager } from '../model'
import type { TodoType } from '../model/types'
import styles from './Todo.module.css'
import { computed } from 'vue'
import { formatTodoDate } from '../lib'

const props = defineProps<{
  todo: TodoType
}>()

const { deleteTodo } = useTodoManager()

const handleDeleteTodo = async () => {
  await deleteTodo(props.todo.id)
}

const formattedDate = computed(() => {
  return formatTodoDate(props.todo.date)
})

</script>


<template>
  <article :class="styles.card">
    <div :class="styles.main">
      <div :class="styles.value">
        {{ todo.value }}
      </div>

      <div :class="styles.meta">
        <span :class="styles.badge">#{{ todo.id }}</span>
        <time :class="styles.date" :datetime="todo.date">
          {{ formattedDate }}
        </time>
      </div>
    </div>

    <button type="button" :class="styles.delete" @click="handleDeleteTodo">
      Delete
    </button>
  </article>
</template>