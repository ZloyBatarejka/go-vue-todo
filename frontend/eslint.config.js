import js from '@eslint/js'
import pluginVue from 'eslint-plugin-vue'
import tseslint from 'typescript-eslint'
import vueParser from 'vue-eslint-parser'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

export default [
  {
    ignores: ['dist/**', 'node_modules/**'],
  },

  js.configs.recommended,

  ...tseslint.configs.recommended,

  ...pluginVue.configs['flat/recommended'],

  {
    files: ['**/*.vue'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: tseslint.parser,
        extraFileExtensions: ['.vue'],
        projectService: true,
        tsconfigRootDir: __dirname,
      },
    },
    rules: {
      // "дефолт" под Vue 3: большинство SFC сейчас используют <script setup>
      'vue/multi-word-component-names': 'off',
      // на ранней стадии проекта удобнее не душить форматированием
      'vue/html-self-closing': 'off',
    },
  },
]


