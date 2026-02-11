# Swagger artifacts (generated)

Эта папка содержит **сгенерированные** OpenAPI/Swagger артефакты для использования фронтендом (автогенерация API, типов и т.п.).

- `swagger/swagger.json`
- `swagger/swagger.yaml`

## Как обновить

Из папки `swagger/`:

```powershell
.\gen-swagger.ps1
```

Примечание: Swagger UI на бэкенде использует пакет `backend/docs` (его тоже обновляет скрипт).


