<script setup lang="ts">
import { computed, ref } from "vue"
import { useRouter } from "vue-router"
import { useAuthManager } from "@entities/auth/model"
import styles from "./AuthPage.module.css"

type AuthMode = "login" | "register"

const router = useRouter()
const { login, register, isAuthenticated } = useAuthManager()

const mode = ref<AuthMode>("login")
const username = ref("")
const password = ref("")
const errorMessage = ref<string | null>(null)
const isSubmitting = ref(false)

const title = computed(() => (mode.value === "login" ? "Login" : "Register"))
const submitLabel = computed(() => (mode.value === "login" ? "Sign in" : "Create account"))
const switchLabel = computed(() =>
    mode.value === "login" ? "No account yet?" : "Already have an account?"
)
const switchActionLabel = computed(() =>
    mode.value === "login" ? "Register" : "Login"
)

const toggleMode = () => {
    mode.value = mode.value === "login" ? "register" : "login"
    errorMessage.value = null
}

const parseAuthError = (error: unknown): string => {
    if (
        typeof error === "object" &&
        error !== null &&
        "error" in error &&
        typeof error.error === "string"
    ) {
        return error.error
    }

    return "Authentication failed"
}

const submit = async () => {
    if (isSubmitting.value) {
        return
    }

    if (!username.value.trim() || !password.value) {
        errorMessage.value = "Username and password are required"
        return
    }

    isSubmitting.value = true
    errorMessage.value = null
    try {
        if (mode.value === "login") {
            await login(username.value.trim(), password.value)
        } else {
            await register(username.value.trim(), password.value)
        }
        await router.push("/")
    } catch (error) {
        errorMessage.value = parseAuthError(error)
    } finally {
        isSubmitting.value = false
    }
}

if (isAuthenticated.value) {
    void router.replace("/")
}
</script>

<template>
    <section :class="styles.page">
        <div :class="styles.card">
            <h2 :class="styles.title">
                {{ title }}
            </h2>
            <p :class="styles.subtitle">
                Use your username and password
            </p>

            <form :class="styles.form" @submit.prevent="submit">
                <input v-model="username" :class="styles.input" type="text" name="username" autocomplete="username"
                    placeholder="Username" />
                <input v-model="password" :class="styles.input" type="password" name="password"
                    autocomplete="current-password" placeholder="Password" />

                <p v-if="errorMessage" :class="styles.error">
                    {{ errorMessage }}
                </p>

                <button :class="styles.button" type="submit" :disabled="isSubmitting">
                    {{ isSubmitting ? "Please wait..." : submitLabel }}
                </button>
            </form>

            <div :class="styles.switch">
                <span>{{ switchLabel }}</span>
                <button :class="styles.switchBtn" type="button" @click="toggleMode">
                    {{ switchActionLabel }}
                </button>
            </div>
        </div>
    </section>
</template>
