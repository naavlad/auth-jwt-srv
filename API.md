# 📡 API Документация

## Swagger UI

Интерактивная документация API доступна через Swagger UI:

```
http://localhost:8080/swagger/index.html
```

### Возможности Swagger UI:

- ✅ Интерактивное тестирование всех эндпойнтов
- ✅ Автоматическая генерация примеров запросов
- ✅ Просмотр схем данных (models)
- ✅ Авторизация через Bearer токен
- ✅ Экспорт OpenAPI спецификации

### Как использовать Swagger UI:

1. **Откройте** `http://localhost:8080/swagger/index.html` в браузере
2. **Раскройте** нужный эндпойнт
3. **Нажмите** "Try it out"
4. **Заполните** параметры запроса
5. **Нажмите** "Execute"
6. **Просмотрите** ответ

### Авторизация в Swagger UI:

1. Выполните `/auth/login` и скопируйте `access_token`
2. Нажмите кнопку **"Authorize"** (🔒) вверху страницы
3. Введите: `Bearer <ваш_токен>`
4. Нажмите **"Authorize"**
5. Теперь можете тестировать защищенные эндпойнты

---

## Эндпойнты

### 1. Health Check

Проверка доступности сервиса.

**Endpoint:** `GET /health`

**Response:**

```json
{
  "status": "ok"
}
```

**Пример curl:**

```bash
curl http://localhost:8080/health
```

**Пример JavaScript (fetch):**

```javascript
fetch("http://localhost:8080/health")
  .then((res) => res.json())
  .then((data) => console.log(data));
```

**Пример Python (requests):**

```python
import requests
response = requests.get('http://localhost:8080/health')
print(response.json())
```

---

### 2. Login (Аутентификация)

Получение JWT токенов для аутентификации.

**Endpoint:** `POST /auth/login`

**Request Body:**

```json
{
  "username": "john_doe",
  "password": "secure_password"
}
```

**Response (200 OK):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImpvaG5fZG9lIiwiZXhwIjoxNzAwMDAwMDAwfQ...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImpvaG5fZG9lIiwiZXhwIjoxNzAwNjA0ODAwfQ..."
}
```

**Коды ошибок:**

- `400 Bad Request` - Невалидный request body или отсутствуют обязательные поля
- `401 Unauthorized` - Неверные учетные данные

**Пример curl:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "secure_password"
  }'
```

**Пример JavaScript:**

```javascript
const response = await fetch("http://localhost:8080/auth/login", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    username: "john_doe",
    password: "secure_password",
  }),
});
const { access_token, refresh_token } = await response.json();
```

**Пример Python:**

```python
import requests

response = requests.post('http://localhost:8080/auth/login', json={
    'username': 'john_doe',
    'password': 'secure_password'
})
tokens = response.json()
access_token = tokens['access_token']
```

**Пример Go:**

```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

body, _ := json.Marshal(LoginRequest{
    Username: "john_doe",
    Password: "secure_password",
})

resp, err := http.Post("http://localhost:8080/auth/login",
    "application/json", bytes.NewBuffer(body))
```

---

### 3. Refresh Token (Обновление токена)

Получение нового access токена без повторной аутентификации.

**Endpoint:** `POST /auth/refresh`

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Коды ошибок:**

- `400 Bad Request` - Refresh token отсутствует
- `401 Unauthorized` - Невалидный или истекший refresh token

**Пример curl:**

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your_refresh_token_here"
  }'
```

**Пример JavaScript:**

```javascript
const response = await fetch("http://localhost:8080/auth/refresh", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    refresh_token: localStorage.getItem("refresh_token"),
  }),
});
const { access_token } = await response.json();
localStorage.setItem("access_token", access_token);
```

**Когда использовать:**

- Access token истек (получили 401 ошибку)
- Перед выполнением важного запроса (проактивное обновление)
- При загрузке приложения (если access token скоро истечет)

---

### 4. Get User Info (Получение данных пользователя)

Получение информации о текущем аутентифицированном пользователе.

**Endpoint:** `GET /auth/me`

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response (200 OK):**

```json
{
  "id": 1,
  "username": "john_doe"
}
```

**Коды ошибок:**

- `401 Unauthorized` - Authorization header отсутствует, невалиден или токен истек

**Пример curl:**

```bash
curl http://localhost:8080/auth/me \
  -H "Authorization: Bearer your_access_token_here"
```

**Пример JavaScript:**

```javascript
const response = await fetch("http://localhost:8080/auth/me", {
  headers: {
    Authorization: `Bearer ${localStorage.getItem("access_token")}`,
  },
});
const user = await response.json();
```

**Пример Python:**

```python
headers = {'Authorization': f'Bearer {access_token}'}
response = requests.get('http://localhost:8080/auth/me', headers=headers)
user = response.json()
```

---

## Модели данных

### LoginRequest

```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

### LoginResponse

```json
{
  "access_token": "string (JWT)",
  "refresh_token": "string (JWT)"
}
```

### RefreshRequest

```json
{
  "refresh_token": "string (JWT, required)"
}
```

### RefreshResponse

```json
{
  "access_token": "string (JWT)"
}
```

### UserInfo

```json
{
  "id": "integer",
  "username": "string"
}
```

### ErrorResponse

```json
{
  "error": "string (описание ошибки)"
}
```

---

## JWT Токены

### Access Token

- **Назначение**: Доступ к защищенным эндпойнтам
- **Срок жизни**: 15 минут (по умолчанию)
- **Использование**: Передается в заголовке `Authorization: Bearer <token>`
- **Payload**:
  ```json
  {
    "user_id": 1,
    "username": "john_doe",
    "exp": 1700000000,
    "iat": 1699999100
  }
  ```

### Refresh Token

- **Назначение**: Обновление access токена
- **Срок жизни**: 7 дней (по умолчанию)
- **Использование**: Отправляется в body запроса `/auth/refresh`
- **Payload**: Аналогичен access токену, но с большим сроком действия

---

## Паттерны использования

### Базовая аутентификация

```javascript
// 1. Login
const loginResponse = await fetch("/auth/login", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ username, password }),
});
const { access_token, refresh_token } = await loginResponse.json();

// Сохраняем токены
localStorage.setItem("access_token", access_token);
localStorage.setItem("refresh_token", refresh_token);

// 2. Используем access token
const userResponse = await fetch("/auth/me", {
  headers: { Authorization: `Bearer ${access_token}` },
});
const user = await userResponse.json();
```

### Автоматическое обновление токена

```javascript
async function fetchWithAuth(url, options = {}) {
  let token = localStorage.getItem("access_token");

  // Добавляем токен в заголовки
  options.headers = {
    ...options.headers,
    Authorization: `Bearer ${token}`,
  };

  let response = await fetch(url, options);

  // Если получили 401 - обновляем токен
  if (response.status === 401) {
    const refreshToken = localStorage.getItem("refresh_token");
    const refreshResponse = await fetch("/auth/refresh", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (refreshResponse.ok) {
      const { access_token } = await refreshResponse.json();
      localStorage.setItem("access_token", access_token);

      // Повторяем запрос с новым токеном
      options.headers["Authorization"] = `Bearer ${access_token}`;
      response = await fetch(url, options);
    }
  }

  return response;
}
```

---

## Коды состояния HTTP

| Код | Значение              | Когда возвращается                     |
| --- | --------------------- | -------------------------------------- |
| 200 | OK                    | Успешный запрос                        |
| 400 | Bad Request           | Невалидные данные в запросе            |
| 401 | Unauthorized          | Отсутствует или невалиден токен/пароль |
| 500 | Internal Server Error | Ошибка на стороне сервера              |

---

## Тестирование API

### Через Swagger UI

Самый простой способ - используйте Swagger UI на `http://localhost:8080/swagger/index.html`

### Через curl

```bash
# 1. Login
TOKEN=$(curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"pass"}' \
  | jq -r '.access_token')

# 2. Используем токен
curl http://localhost:8080/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

### Через Postman

1. Импортируйте OpenAPI спецификацию: `http://localhost:8080/swagger/doc.json`
2. Создайте environment с переменными `base_url`, `access_token`
3. Используйте Pre-request Scripts для автоматической авторизации

### Через HTTPie

```bash
# Login
http POST :8080/auth/login username=user password=pass

# Get user info
http :8080/auth/me Authorization:"Bearer <token>"
```

---

## Регенерация документации

После изменения кода handlers обновите Swagger документацию:

```bash
make swagger
```

Это пересоздаст файлы в директории `docs/`.
