<script setup lang="ts">
import { useAuthManager } from "@entities/auth/model"
import { useTodoManager } from "@entities/todo/model"
import { watch } from "vue"
import { useRouter } from "vue-router"
import styles from './HomePage.module.css'
import { TodoList } from '@widgets/todoList'
import { RouterLink } from "vue-router"

const { isAuthenticated, user, logout } = useAuthManager()
const { clearTodos } = useTodoManager()
const router = useRouter()

watch(isAuthenticated, (nextValue) => {
  if (!nextValue) {
    void router.push("/auth")
  }
})

const logoutAndResetState = () => {
  clearTodos()
  logout()
}

</script>

<template>
  <div :class="styles.welcome">
    <div v-if="isAuthenticated">
      <p>Signed in as {{ user?.username }}</p>
      <button type="button" @click="logoutAndResetState">
        Logout
      </button>
    </div>
    <div v-else>
      <RouterLink to="/auth">
        <button type="button">
          Login / Register
        </button>
      </RouterLink>
    </div>
    <TodoList />
  </div>
</template>
