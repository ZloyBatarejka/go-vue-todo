<script setup lang="ts">

import styles from './HomePage.module.css'
import { TodoList } from '@widgets/todoList'

</script>

<template>
  <div :class="styles.welcome">
    <TodoList />
  </div>
</template>
