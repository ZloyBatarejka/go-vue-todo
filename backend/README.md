# Go TODO Backend

Профессиональный TODO проект на Go с архитектурой, разделенной на слои.

## Архитектура

Проект следует принципам чистой архитектуры и разделен на слои:

- **models** - модели данных (структуры)
- **database** - подключение к базе данных
- **repository** - слой работы с БД (data access layer)
- **handlers** - HTTP обработчики (presentation layer)
- **main.go** - точка входа, инициализация и роутинг

## Требования

- Go 1.21 или выше
- PostgreSQL (запущенный на localhost:5432)

## Конфигурация базы данных (env / .env)

Настройки подключения берутся из переменных окружения:

- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5432`)
- `DB_USER` (default: `postgres`)
- `DB_PASSWORD` (default: пусто)
- `DB_NAME` (default: `postgres`)
- `DB_SSLMODE` (default: `disable`)

Для локальной разработки можно создать файл `backend/.env` (он подхватится автоматически при старте).
В репозитории есть пример: `backend/env.example` — просто переименуй его в `.env` и заполни пароль.

## Установка зависимостей

```bash
go mod download
```

## Запуск

### Обычный запуск:
```bash
go run main.go
```

### С hot reload (автоматическая перезагрузка при изменениях):
```bash
air
```

Если `air` не установлен:
```bash
go install github.com/air-verse/air@latest
```

Сервер запустится на `http://localhost:8080`

## Swagger (API документация)

1) Запусти сервер (см. раздел **Запуск**) — он стартует на `http://localhost:8080`.
2) Открой Swagger UI в браузере:

`http://localhost:8080/swagger/index.html`

### Генерация swagger.json/yaml (для фронта / автогенерации API)

В репозитории есть корневая папка `swagger/` (рядом с `frontend/` и `backend/`) — туда складываются артефакты:
- `swagger/swagger.json`
- `swagger/swagger.yaml`

Чтобы обновить Swagger, из **корня репозитория** запусти:

```powershell
.\swagger\gen-swagger.ps1
```

## API Endpoints

### Создать Todo

**POST** `/api/todos`

**Request Body:**
```json
{
  "value": "Название задачи"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "value": "Название задачи",
  "date": "2024-01-15T12:34:56Z"
}
```

### Получить все Todo

**GET** `/api/todos`

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "value": "Название задачи",
    "date": "2024-01-15"
  }
]
```

### Получить Todo по ID

**GET** `/api/todos/{id}`

**Response (200 OK):**
```json
{
  "id": 1,
  "value": "Название задачи",
  "date": "2024-01-15"
}
```

### Удалить Todo

**DELETE** `/api/todos/{id}`

**Response (200 OK):**
```json
{
  "message": "Todo deleted successfully"
}
```

## Примеры использования

### Создать задачу
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"value": "Купить молоко"}'
```

### Получить все задачи
```bash
curl http://localhost:8080/api/todos
```

### Получить задачу
```bash
curl http://localhost:8080/api/todos/1
```

### Удалить задачу
```bash
curl -X DELETE http://localhost:8080/api/todos/1
```

## Особенности реализации

- ✅ ID генерируется автоматически на бэкенде (PostgreSQL SERIAL)
- ✅ Разделение на слои (handlers → repository → database)
- ✅ Использование интерфейсов для легкого тестирования
- ✅ Правильная обработка ошибок
- ✅ Валидация входных данных
- ✅ Настройка пула соединений с БД
- ✅ Dependency Injection
- ✅ CORS поддержка для фронтенда
- ✅ Hot reload с помощью `air`
