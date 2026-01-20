<script lang="ts" setup>
import { useTodoManager } from '@entities/todo/model'
import { ref } from 'vue'
import styles from './AddTodo.module.css'

const { createTodo } = useTodoManager()

const value = ref('')

const submit = async () => {
    await createTodo(value.value)
    value.value = ''
}
</script>


<template>
    <form :class="styles.form" @submit.prevent="submit">
        <input v-model="value" :class="styles.input" type="text" placeholder="New todo..." autocomplete="off" />
        <button type="submit" :class="styles.button">
            Add
        </button>
    </form>
</template>