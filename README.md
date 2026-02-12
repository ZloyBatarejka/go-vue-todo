# goTodo (monorepo)

Монорепозиторий с **Go backend** + **Vue 3 (Vite) frontend** для Todo-приложения.

## Структура

- `backend/` — Go API + PostgreSQL + Swagger UI
- `frontend/` — Vue 3 + TypeScript + Vite (FSD-структура)

## Быстрый старт

### Backend

1) Перейди в `backend/`.
2) Создай файл `backend/.env` (пример: `backend/env.example`).
3) Запусти:

```bash
go run main.go
```

Backend стартует на `http://localhost:8080`.

- **API**: `http://localhost:8080/api`
- **Swagger UI**: `http://localhost:8080/swagger/index.html`

### Frontend

1) Перейди в `frontend/`.
2) Установи зависимости и запусти dev-сервер:

```bash
npm i
npm run dev
```

Frontend по умолчанию стартует на `http://localhost:5173`.

## Переменные окружения (backend)

Backend читает настройки БД из переменных окружения (и автоматически подхватывает `backend/.env`, если он есть):

- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5432`)
- `DB_USER` (default: `postgres`)
- `DB_PASSWORD` (default: пусто)
- `DB_NAME` (default: `postgres`)
- `DB_SSLMODE` (default: `disable`)

## CORS

Backend разрешает запросы с origin `http://localhost:5173` (Vite dev server).

## Линт (frontend)

Из `frontend/`:

```bash
npm run lint
npm run lint:fix
```

## Планы / заметки

- JWT авторизация (план внедрения): `JWT_AUTH_PLAN.md`