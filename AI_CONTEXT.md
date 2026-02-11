# AI_CONTEXT (goTodo)

Этот файл — **внутренние заметки для быстрого погружения** в монорепозиторий `goTodo` (backend + frontend). Держу здесь “карту проекта”: что где лежит, как запустить, куда смотреть при доработках.

## TL;DR (как запустить локально)

### Backend (Go + PostgreSQL)
- **Папка**: `backend/`
- **Env**: создай `backend/.env` на базе `backend/env.example`
- **Запуск**:

```bash
cd backend
go run main.go
```

- **URL**:
  - API: `http://localhost:8080/api`
  - Health: `http://localhost:8080/health`
  - Swagger UI: `http://localhost:8080/swagger/index.html`

### Frontend (Vue 3 + Vite + Pinia)
- **Папка**: `frontend/`
- **Запуск**:

```bash
cd frontend
npm i
npm run dev
```

- **URL**: `http://localhost:5173`

## Архитектура и ключевые места

### Backend (`backend/`)
- **Точка входа / роутинг / CORS / Swagger**: `backend/main.go`
  - Роуты:
    - `GET    /api/todos`
    - `POST   /api/todos`
    - `GET    /api/todos/{id}`
    - `DELETE /api/todos/{id}`
  - CORS разрешён для origin `http://localhost:5173`
- **HTTP-слой**: `backend/handlers/todo.go`
  - Валидация: `value` обязателен
  - Ошибки возвращаются как `{"error": "..."}`
- **Repository слой (SQL)**: `backend/repository/todo.go`
  - Работает с таблицей `todos` и полями `id, value, date`
  - При `Create()` делает `CREATE SEQUENCE IF NOT EXISTS todos_id_seq` и сам выставляет дату на бэкенде
- **DB подключение**: `backend/database/db.go` (PostgreSQL через `github.com/lib/pq`)
- **Модели**: `backend/models/todo.go`
- **Swagger артефакты**: `backend/docs/*` (сгенерированные файлы + `swagger.yaml/json`)

### Frontend (`frontend/`)
- **Точка входа**: `frontend/src/app/main.ts` (Vue + Pinia + Router)
- **Маршрутизация**: `frontend/src/app/router/index.ts` (сейчас одна страница `/`)
- **Структура**: FSD-подобная (`app/`, `pages/`, `widgets/`, `features/`, `entities/`, `shared/`)

#### Поток данных “Todo”
- UI:
  - `HomePage` → `TodoList` (`frontend/src/pages/home/ui/HomePage.vue`, `frontend/src/widgets/todoList/ui/TodoList.vue`)
  - создание: `AddTodo` (`frontend/src/features/addTodo/ui/AddTodo.vue`)
  - отображение/удаление: `Todo` (`frontend/src/entities/todo/ui/Todo.vue`)
- “Менеджер” (удобная обёртка над стором):
  - `useTodoManager()` (`frontend/src/entities/todo/model/manager.ts`)
- Pinia store:
  - `useTodoStore()` (`frontend/src/entities/todo/model/store.ts`)
  - хранит `todos`, `isLoading`, `todosCount`; методы: `fetchTodos`, `createTodo`, `deleteTodo`
- API слой:
  - `TodoApiService` (`frontend/src/entities/todo/api/todo-api.service.ts`)
  - endpoints: `frontend/src/shared/config/endpoints.ts` (пути без `/api`, т.к. базовый URL уже содержит `/api`)
  - HTTP-клиент: `frontend/src/shared/api/client.ts`

## API контракт (как сейчас)

### Backend модели
- `Todo`: `{ id: number, value: string, date: string }`
- Create request: `{ value: string }`

### Важная заметка про `DELETE`
- Бэкенд возвращает JSON вида `{"message":"Todo deleted successfully"}`.
- Фронт сейчас вызывает `httpClient.delete<void>(...)` и **не использует тело ответа**.
  - Это обычно ок, но сейчас `ApiClient` всегда делает `response.json()`, т.е. на практике он попытается распарсить JSON даже для `void`.
  - Если когда-нибудь бэкенд поменяется на `204 No Content`, этот код начнёт падать на `response.json()`. (См. `frontend/src/shared/api/client.ts`.)

## База данных (важно для локального старта)

Репозиторий ожидает таблицу `todos`, но **миграций/DDL в репо нет** (таблица не создаётся автоматически). Нужно создать вручную.

Минимальная схема, совместимая с запросами:

```sql
CREATE TABLE IF NOT EXISTS todos (
  id   BIGINT PRIMARY KEY,
  value TEXT NOT NULL,
  date  TEXT NOT NULL
);
```

Примечание: репозиторий использует `todos_id_seq` и вставляет `id` через `nextval('todos_id_seq')`, так что `id` не обязан быть `SERIAL`, но последовательность должна существовать (она создаётся кодом).

## Конфигурация / env

### Backend
Переменные (см. `backend/env.example`): `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`.

### Frontend
Сейчас base URL API **захардкожен** в `frontend/src/shared/api/client.ts`:
- `API_BASE_URL = 'http://localhost:8080/api'`

Если понадобится деплой/разные окружения — выносить в `import.meta.env` (Vite env).

## Где править, если…
- **Добавить поля Todo**: `backend/models/todo.go`, SQL в `backend/repository/todo.go`, типы в `frontend/src/entities/todo/model/types.ts`, UI-компоненты.
- **Добавить новые endpoints**: роуты в `backend/main.go`, handler + repo; на фронте — `shared/config/endpoints.ts` + сервис `todo-api.service.ts` + store.
- **Починить/улучшить обработку ошибок на фронте**: `frontend/src/shared/api/client.ts` (сейчас кидает общий `Error` без тела ответа).


