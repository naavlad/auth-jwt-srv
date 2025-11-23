# 🧪 Тестирование

## Структура тестов

Проект использует стандартную библиотеку `testing` Go и `testify` для assertions.

```
internal/
├── config/
│   ├── config.go
│   └── config_test.go          # Unit тесты конфигурации
├── handlers/
│   ├── *.go
│   └── handlers_test.go        # Unit тесты handlers
├── tokens/
│   ├── jwt.go
│   └── jwt_test.go             # Unit тесты JWT
└── service/
    ├── auth.go
    └── auth_test.go            # (будущие тесты)
```

---

## Запуск тестов

### Все тесты

```bash
make test
```

### С покрытием кода

```bash
make test-coverage
```

Создаст файл `coverage.html` - откройте в браузере для визуализации покрытия.

### Только конкретный пакет

```bash
go test ./internal/tokens
go test ./internal/handlers
```

### Verbose mode (подробный вывод)

```bash
go test -v ./...
```

### Конкретный тест

```bash
go test -v -run TestGenerateAccessToken ./internal/tokens
```

---

## Unit тесты

### Пример: JWT токены (`internal/tokens/jwt_test.go`)

```go
func TestGenerateAccessToken(t *testing.T) {
    manager := NewManager("test-secret-key", 15*time.Minute, 168*time.Hour)

    token, err := manager.GenerateAccessToken(1, "testuser")

    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}
```

**Что тестируем:**

- ✅ Генерация access токена
- ✅ Генерация refresh токена
- ✅ Валидация токена
- ✅ Обработка невалидного токена
- ✅ Обработка истекшего токена

### Пример: HTTP Handlers (`internal/handlers/handlers_test.go`)

```go
func TestHealth(t *testing.T) {
    handler := &Handler{}

    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    w := httptest.NewRecorder()

    handler.Health(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response HealthResponse
    json.NewDecoder(w.Body).Decode(&response)
    assert.Equal(t, "ok", response.Status)
}
```

**Что тестируем:**

- ✅ Health check эндпойнт
- ✅ JSON serialization (writeJSON)
- ✅ Error handling (writeError)
- ✅ Валидация входных данных

### Пример: Конфигурация (`internal/config/config_test.go`)

```go
func TestLoad(t *testing.T) {
    os.Setenv("DATABASE_URL", "postgres://test@localhost/test")
    os.Setenv("JWT_SECRET", "secret")
    defer os.Unsetenv("DATABASE_URL")
    defer os.Unsetenv("JWT_SECRET")

    cfg, err := Load()

    assert.NoError(t, err)
    assert.NotNil(t, cfg)
    assert.Equal(t, "postgres://test@localhost/test", cfg.Database.URL)
}
```

**Что тестируем:**

- ✅ Загрузка конфигурации из environment
- ✅ Валидация обязательных полей
- ✅ Дефолтные значения

---

## Паттерны тестирования

### Table-Driven Tests (рекомендуется)

```go
func TestValidateToken(t *testing.T) {
    tests := []struct {
        name    string
        token   string
        wantErr bool
    }{
        {
            name:    "valid token",
            token:   "eyJhbGc...",
            wantErr: false,
        },
        {
            name:    "invalid token",
            token:   "invalid",
            wantErr: true,
        },
        {
            name:    "empty token",
            token:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := manager.ValidateToken(tt.token)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### HTTP тестирование

```go
func TestLogin_Success(t *testing.T) {
    // Создаем mock service
    mockService := &MockAuthService{}
    handler := NewHandler(mockService)

    // Готовим request
    body := `{"username":"user","password":"pass"}`
    req := httptest.NewRequest(http.MethodPost, "/auth/login",
        strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    // Записываем response
    w := httptest.NewRecorder()

    // Выполняем handler
    handler.Login(w, req)

    // Проверяем результат
    assert.Equal(t, http.StatusOK, w.Code)

    var response LoginResponse
    json.NewDecoder(w.Body).Decode(&response)
    assert.NotEmpty(t, response.AccessToken)
}
```

### Тестирование с моками

```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
    args := m.Called(ctx, username)
    return args.Get(0).(*User), args.Error(1)
}

func TestLogin_WithMock(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("GetUserByUsername", mock.Anything, "user").
        Return(&User{ID: 1, Username: "user"}, nil)

    service := NewAuthService(mockRepo, tokenManager)

    result, err := service.Login(context.Background(), "user", "pass")

    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

---

## Покрытие кода

### Просмотр покрытия

```bash
# Генерация отчета
go test -coverprofile=coverage.out ./...

# HTML отчет
go tool cover -html=coverage.out -o coverage.html

# Или через make
make test-coverage
```

### Целевые показатели

- **Критичные компоненты** (tokens, auth): > 80%
- **Handlers**: > 70%
- **Utils/helpers**: > 60%

### Текущее покрытие

```bash
go test -cover ./internal/...
```

Результат:

```
ok      internal/config     coverage: 85.7%
ok      internal/handlers   coverage: 72.3%
ok      internal/tokens     coverage: 91.2%
```

---

## Интеграционные тесты

### Будущая структура

```
tests/
├── integration/
│   ├── api_test.go          # Тесты полного API flow
│   └── db_test.go           # Тесты с реальной БД
└── e2e/
    └── auth_flow_test.go    # End-to-end тесты
```

### Пример интеграционного теста

```go
// +build integration

func TestLoginFlow_Integration(t *testing.T) {
    // Подключаемся к тестовой БД
    db := setupTestDB(t)
    defer db.Close()

    // Создаем тестового пользователя
    insertTestUser(t, db, "testuser", "password")

    // Запускаем сервер
    srv := startTestServer(t, db)
    defer srv.Close()

    // Выполняем login
    resp := doLogin(t, srv.URL, "testuser", "password")
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    var tokens LoginResponse
    json.NewDecoder(resp.Body).Decode(&tokens)
    assert.NotEmpty(t, tokens.AccessToken)

    // Используем токен
    userResp := getMe(t, srv.URL, tokens.AccessToken)
    assert.Equal(t, http.StatusOK, userResp.StatusCode)
}
```

Запуск интеграционных тестов:

```bash
go test -tags=integration ./tests/integration/...
```

---

## Бенчмарки

### Пример бенчмарка

```go
func BenchmarkGenerateAccessToken(b *testing.B) {
    manager := NewManager("secret", 15*time.Minute, 168*time.Hour)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        manager.GenerateAccessToken(1, "user")
    }
}
```

Запуск:

```bash
go test -bench=. ./internal/tokens
```

Результат:

```
BenchmarkGenerateAccessToken-8    50000    28543 ns/op
```

---

## Best Practices

### 1. Именование тестов

```go
// ✅ Хорошо
func TestGenerateAccessToken(t *testing.T) {}
func TestValidateToken_ExpiredToken(t *testing.T) {}
func TestLogin_InvalidCredentials(t *testing.T) {}

// ❌ Плохо
func TestFunc1(t *testing.T) {}
func Test1(t *testing.T) {}
```

### 2. Arrange-Act-Assert

```go
func TestLogin(t *testing.T) {
    // Arrange (подготовка)
    handler := &Handler{service: mockService}
    req := httptest.NewRequest(http.MethodPost, "/auth/login", body)
    w := httptest.NewRecorder()

    // Act (действие)
    handler.Login(w, req)

    // Assert (проверка)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### 3. Изоляция тестов

```go
// ✅ Каждый тест независим
func TestA(t *testing.T) {
    manager := NewManager(...)  // Создаем новый объект
    // тест
}

func TestB(t *testing.T) {
    manager := NewManager(...)  // Создаем новый объект
    // тест
}
```

### 4. Используйте testify

```go
// ✅ С testify
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.NotEmpty(t, token)

// ❌ Без testify
if actual != expected {
    t.Errorf("expected %v, got %v", expected, actual)
}
```

### 5. Cleanup в тестах

```go
func TestWithCleanup(t *testing.T) {
    os.Setenv("TEST_VAR", "value")
    defer os.Unsetenv("TEST_VAR")  // Очищаем после теста

    // тест
}
```

---

## Debugging тестов

### Verbose mode

```bash
go test -v ./internal/tokens
```

### Логирование в тестах

```go
func TestSomething(t *testing.T) {
    t.Logf("Testing with value: %v", value)

    result, err := DoSomething(value)
    if err != nil {
        t.Fatalf("Failed: %v", err)
    }

    t.Logf("Result: %v", result)
}
```

### Пропуск тестов

```go
func TestLongRunning(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping long test")
    }
    // длинный тест
}
```

Запуск без длинных тестов:

```bash
go test -short ./...
```

---

## CI/CD Integration

### GitHub Actions пример

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: "1.25"
      - run: go test -v -coverprofile=coverage.out ./...
      - run: go tool cover -html=coverage.out -o coverage.html
      - uses: actions/upload-artifact@v3
        with:
          name: coverage
          path: coverage.html
```

---

## Дальнейшее изучение

### Рекомендуемые ресурсы:

1. **Go Testing**: https://go.dev/doc/tutorial/add-a-test
2. **Testify**: https://github.com/stretchr/testify
3. **Table-Driven Tests**: https://go.dev/wiki/TableDrivenTests
4. **Testing Best Practices**: https://golang.org/doc/effective_go#testing

### Следующие шаги:

- [ ] Добавить тесты для `service` слоя
- [ ] Написать интеграционные тесты с реальной БД
- [ ] Добавить e2e тесты
- [ ] Настроить CI/CD с автоматическим запуском тестов
- [ ] Добавить тесты производительности (benchmarks)
