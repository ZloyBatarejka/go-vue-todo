<script setup lang="ts">
import { ref } from 'vue'
import styles from './HomePage.module.css'
import { httpClient } from '@shared/api'
import {TodoEndpoints} from '@shared/config'

const count = ref(0)

const increment = () => {
  httpClient.get(TodoEndpoints.getAll).then(console.log)
  count.value++
}
</script>

<template>
  <div :class="styles.welcome">
    <button @click="increment">Click me</button>
    <p>{{ count }}</p>
  </div>
</template>

