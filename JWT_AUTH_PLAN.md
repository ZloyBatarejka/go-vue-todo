# JWT авторизация — план внедрения (goTodo)

Документ описывает план внедрения JWT-авторизации **от базы данных до фронтенда** для текущего монорепо.

## Цели

- Добавить пользователей  (минимум: user).
- Добавить регистрацию/логин.
- Защитить приватные эндпоинты (минимум: создание/удаление todo).
- Подготовить фронт к работе с авторизацией (UI/стор, хранение токенов, авто-рефреш).
- Поддержать Swagger (кнопка Authorize, security схемы).

## Ключевые решения (зафиксированы на 1-й этап “для обучения”)

### 1) Где хранить access token на фронте

- ✅ **Выбрано: `localStorage` (C)** — осознанно как простой старт для обучения.
  - Важно: это менее безопасно (XSS-риск), позже мигрируем на httpOnly cookie.

### 2) Нужен ли refresh token
- ✅ **Нет**: на 1-м этапе только access token.
  - Когда access истёк → пользователь перелогинивается.

### 3) Что считаем “Todo пользователя”
- ✅ **User-owned todos (A)**: у todo есть `user_id`, эндпоинты возвращают только свои.

### 4) Как передаём access token
- ✅ **`Authorization: Bearer <access>`**
  - Токен берём из `localStorage` и добавляем к запросам централизованно (interceptor / `securityWorker`).

### 5) Пользователи
- ✅ Логин по **`username`**

### 6) Регистрация/аккаунты (UX)
- ✅ Регистрация **открытая**
- ✅ Email-верификация **не нужна** на первом шаге
- ✅ Мультисессии **не нужны**

## План миграции (2-й этап): перейти на httpOnly cookie + refresh

Когда базовый JWT-логин и защита эндпоинтов заработают, делаем “правильнее”:

1) Перенос access из storage
   - вариант: access in-memory, refresh в httpOnly cookie
   - вариант: access тоже в httpOnly cookie (проще фронту, но добавляет CSRF-тему)
2) Добавить refresh token
   - `POST /api/auth/refresh`
   - rotation refresh токена
   - хранить **хеш** refresh в БД
3) Настроить cookie/CORS (credentials + SameSite/Secure)
4) На фронте: авто refresh при 401 + повтор запроса

## Операционные шаги по стекам (для команд "делаем шаг X")

Ниже шаги в формате `STACK.N`, чтобы можно было ссылаться коротко:
- `DB.N` — база данных
- `BE.N` — backend (Go API)
- `FE.N` — frontend (Vue)

### Этап 1 (текущий, учебный: localStorage + Bearer, без refresh)

#### DB стек
- `DB.1` Создать таблицу `users` (`id`, `username`, `password_hash`, `created_at`).
- `DB.2` Добавить `todos.user_id` + FK на `users(id)` + `ON DELETE CASCADE`.
- `DB.3` Добавить индекс для пользовательских todo:
  - минимум: `CREATE INDEX ON todos(user_id);`
  - рекомендовано: `CREATE INDEX ON todos(user_id, id DESC);`
- `DB.4` Убедиться, что у `todos.id` настроен автоген (`DEFAULT nextval(...)` или `IDENTITY`).
- `DB.5` Завести тестового пользователя для временной заглушки (до полной auth-схемы).

#### BE стек
- `BE.1` Добавить модели `User`, `RegisterRequest`, `LoginRequest`, `AuthResponse`.
- `BE.2` Реализовать `UserRepository` (`CreateUser`, `FindByUsername`, `FindByID`).
- `BE.3` Реализовать auth-сервис: `bcrypt` + генерация/валидация JWT access.
- `BE.4` Добавить ручки `POST /api/auth/register`, `POST /api/auth/login`.
- `BE.5` Добавить `AuthMiddleware` (чтение Bearer токена, проверка JWT, `userId` в context).
- `BE.6` Защитить todo-ручки и перейти на user-owned запросы:
  - `GET /api/todos` → только `WHERE user_id = currentUserId`
  - `DELETE /api/todos/{id}` → только `WHERE id = $1 AND user_id = currentUserId`
  - `POST /api/todos` → вставка с `user_id = currentUserId` (убрать временную заглушку).
- `BE.7` Обновить Swagger аннотации и перегенерить спеки.
- `BE.8` Проверить CORS для Bearer (`Authorization` в allowed headers).

#### FE стек
- `FE.1` Добавить auth-store (token в `localStorage`, user, auth-status).
- `FE.2` Подключить Bearer к generated API через `securityWorker`.
- `FE.3` Сделать UI логина/регистрации.
- `FE.4` Добавить router-guard для защищённых страниц.
- `FE.5` Обновить Todo API flow под авторизованные запросы.
- `FE.6` Добавить logout (очистка `localStorage` + сброс auth-store).

### Этап 2 (миграция: httpOnly cookie + refresh)

#### DB стек
- `DB.6` Добавить поля/таблицу под refresh token hash.
- `DB.7` Подготовить схему под rotation refresh токена.

#### BE стек
- `BE.9` Добавить `POST /api/auth/refresh` и `POST /api/auth/logout`.
- `BE.10` Реализовать refresh rotation + хранение hash refresh токена.
- `BE.11` Настроить cookie-политику (`HttpOnly`, `SameSite`, `Secure`, `Path/Domain`).
- `BE.12` Включить `AllowCredentials: true` + корректный CORS для cookie.

#### FE стек
- `FE.7` Перенести access из storage (в память или cookie-схему).
- `FE.8` Включить `credentials: include` для refresh/cookie flow.
- `FE.9` Добавить авто-refresh при 401 и retry исходного запроса.
- `FE.10` Удалить legacy-логику хранения access в storage (после стабилизации).

## TODO (отдельно): миграция Bearer/storage -> httpOnly cookie

Этот блок — отдельный чеклист финального перехода на cookie-схему после стабилизации этапа 1.

- [ ] `MIG.1` Backend: добавить `POST /api/auth/refresh` и `POST /api/auth/logout`.
- [ ] `MIG.2` Backend: внедрить refresh rotation + хранение hash refresh токена в БД.
- [ ] `MIG.3` Backend: выставлять refresh в `Set-Cookie` (`HttpOnly`, `SameSite`, `Secure`, `Path/Domain`).
- [ ] `MIG.4` Backend: включить `AllowCredentials: true` и проверить CORS для cookie-flow.
- [ ] `MIG.5` Frontend: включить `credentials: include` для refresh/cookie запросов.
- [ ] `MIG.6` Frontend: добавить авто-refresh при 401 и повтор исходного запроса.
- [ ] `MIG.7` Frontend: убрать хранение access в storage после валидации нового flow.
- [ ] `MIG.8` QA: проверить logout, истечение токена, повторный логин, сценарии "несколько вкладок".


