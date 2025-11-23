# GitHub Copilot Instructions

## Контекст проекта

Это микросервис аутентификации на Go с использованием JWT токенов. Проект следует принципам Clean Architecture и использует современный стек технологий.

## Архитектурные принципы

### Структура проекта

- `cmd/api/` - точка входа приложения
- `internal/config/` - конфигурация с использованием cleanenv
- `internal/handlers/` - HTTP обработчики (тонкий слой)
- `internal/service/` - бизнес-логика
- `internal/repository/` - доступ к данным (генерируется sqlc)
- `internal/tokens/` - JWT утилиты
- `queries/` - SQL запросы для sqlc

### Соблюдай Clean Architecture

1. **Handlers** - только парсинг HTTP, валидация входных данных, вызов service
2. **Service** - вся бизнес-логика, работа с repository и tokens
3. **Repository** - только SQL запросы (генерируется sqlc)
4. Зависимости идут от handlers → service → repository

## Стандарты кодирования

### Именование

- Пакеты: короткие, lowercase, без подчеркиваний (`handlers`, `service`)
- Функции/методы: CamelCase (`GetUserByID`, `GenerateAccessToken`)
- Константы: CamelCase или UPPER_CASE для экспортируемых
- Переменные: camelCase для локальных, CamelCase для экспортируемых

### Структуры и типы

```go
// Всегда добавляй JSON теги
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Используй указатели для больших структур в response
type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}
```

### Обработка ошибок

```go
// Всегда оборачивай ошибки с контекстом
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// В handlers возвращай понятные HTTP ошибки
if err != nil {
    writeError(w, http.StatusUnauthorized, "invalid credentials")
    return
}
```

### HTTP Handlers

```go
// Шаблон для handlers:
func (h *Handler) MethodName(w http.ResponseWriter, r *http.Request) {
    // 1. Парсинг и валидация входных данных
    var req RequestType
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    // 2. Простая валидация
    if req.Field == "" {
        writeError(w, http.StatusBadRequest, "field is required")
        return
    }

    // 3. Вызов service
    result, err := h.service.Method(r.Context(), req.Field)
    if err != nil {
        writeError(w, http.StatusSomeError, "error message")
        return
    }

    // 4. Возврат ответа
    writeJSON(w, http.StatusOK, result)
}
```

### Context usage

- Всегда передавай `context.Context` первым параметром
- Используй `r.Context()` в handlers
- Передавай контекст в repository запросы

## Работа с базой данных

### sqlc

- SQL запросы пишутся в `queries/*.sql`
- Используй именованные параметры: `$1, $2`
- Всегда добавляй комментарии с именем запроса:

```sql
-- name: GetUserByUsername :one
SELECT id, username, password
FROM users
WHERE username = $1
LIMIT 1;
```

### Типы запросов sqlc

- `:one` - возвращает одну запись
- `:many` - возвращает массив записей
- `:exec` - выполнение без возврата данных

## JWT и безопасность

### JWT токены

- Access token: короткий срок жизни (15 минут по умолчанию)
- Refresh token: длинный срок жизни (7 дней по умолчанию)
- Всегда проверяй токены через `tokenManager.ValidateToken()`

### Извлечение токена из header

```go
authHeader := r.Header.Get("Authorization")
parts := strings.Split(authHeader, " ")
if len(parts) != 2 || parts[0] != "Bearer" {
    // error
}
token := parts[1]
```

## Конфигурация

### Используй cleanenv

```go
type Config struct {
    Field string `env:"FIELD_NAME" env-required:"true"`
    OptionalField string `env:"OPTIONAL" env-default:"default-value"`
}
```

### Переменные окружения

- Всегда документируй в `.env.example`
- Обязательные поля: `env-required:"true"`
- Дефолтные значения: `env-default:"value"`

## HTTP Response helpers

### Используй существующие функции

```go
writeJSON(w, status, data)      // для успешных ответов
writeError(w, status, message)  // для ошибок
```

## Middleware

### Если нужно добавить middleware

```go
// В cmd/api/main.go используй chi middleware
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(middleware.RequestID)
```

## Docker

### Dockerfile

- Multi-stage build
- Alpine для production образа
- CGO_ENABLED=0 для статической компиляции

### Environment в docker-compose

- Всегда используй переменные из `.env`
- Формат: `${VARIABLE_NAME:-default}`

## Комментарии

### Когда комментировать

- Экспортируемые функции и типы (godoc format)
- Сложная бизнес-логика
- Не очевидные решения

### Формат godoc

```go
// GetUserInfo возвращает информацию о пользователе по токену.
// Валидирует access токен и делает запрос в БД.
func (s *Service) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
    // implementation
}
```

## Тестирование (если добавляешь)

### Структура тестов

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    WantType
        wantErr bool
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Логирование

### Используй стандартный log

```go
log.Printf("Starting server on port %s", port)
log.Fatalf("Failed to connect: %v", err)
```

### Не логируй sensitive данные

- Пароли
- JWT токены (целиком)
- Полные connection strings

## Git commit messages

### Формат

```
<type>: <subject>

<body>
```

### Типы

- `feat:` - новая функциональность
- `fix:` - исправление бага
- `refactor:` - рефакторинг без изменения функциональности
- `docs:` - изменения в документации
- `chore:` - обновление зависимостей, конфигурации

## Что НЕ делать

❌ Не смешивай бизнес-логику в handlers
❌ Не делай прямые SQL запросы (используй sqlc)
❌ Не игнорируй ошибки
❌ Не используй `panic` без крайней необходимости
❌ Не храни секреты в коде
❌ Не логируй чувствительные данные
❌ Не создавай глобальные переменные без необходимости

## Что ВСЕГДА делать

✅ Проверяй ошибки
✅ Используй context.Context
✅ Валидируй входные данные
✅ Возвращай понятные HTTP коды и сообщения
✅ Следуй принципам Clean Architecture
✅ Пиши godoc комментарии для экспортируемых функций
✅ Используй существующие helper функции
✅ Обрабатывай graceful shutdown

## Примеры кода

### Добавление нового эндпойнта

1. **SQL запрос** (`queries/users.sql`):

```sql
-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1;
```

2. **Сгенерировать код**: `make sqlc-generate`

3. **Service метод** (`internal/service/auth.go`):

```go
func (s *AuthService) UpdatePassword(ctx context.Context, userID int32, newPassword string) error {
    err := s.repo.UpdateUserPassword(ctx, repository.UpdateUserPasswordParams{
        ID:       userID,
        Password: newPassword,
    })
    if err != nil {
        return fmt.Errorf("failed to update password: %w", err)
    }
    return nil
}
```

4. **Handler** (`internal/handlers/auth.go`):

```go
type UpdatePasswordRequest struct {
    Password string `json:"password"`
}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
    var req UpdatePasswordRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    // Получаем userID из токена
    authHeader := r.Header.Get("Authorization")
    // ... извлечение токена ...

    if err := h.authService.UpdatePassword(r.Context(), userID, req.Password); err != nil {
        writeError(w, http.StatusInternalServerError, "failed to update password")
        return
    }

    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
```

5. **Роут** (`cmd/api/main.go`):

```go
r.Put("/auth/password", handler.UpdatePassword)
```

## Производительность

- Используй connection pool для БД (pgxpool)
- Устанавливай таймауты для HTTP сервера
- Закрывай ресурсы через `defer`
- Используй `context.WithTimeout` для долгих операций

## Версионирование API

При добавлении breaking changes:

- Используй версионирование в URL: `/v1/auth/login`, `/v2/auth/login`
- Или через заголовки: `Accept: application/vnd.api.v1+json`

---

**Помни**: Код должен быть простым, читаемым и поддерживаемым. Если что-то кажется слишком сложным - возможно, есть более простое решение.
