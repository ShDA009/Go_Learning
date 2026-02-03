# HTTP тестирование

---

## 💡 Ключевые идеи

1. **httptest** — стандартный пакет для тестирования HTTP
2. **httptest.NewRecorder** — записывает HTTP-ответ для проверки
3. **httptest.NewServer** — создаёт тестовый HTTP-сервер
4. **Изоляция** — тестируем handlers без реального сервера
5. **End-to-end** — тестирование с реальными HTTP-запросами

### Два подхода к тестированию HTTP

| Подход | Инструмент | Когда использовать |
|--------|------------|-------------------|
| Unit test handler | `httptest.NewRecorder` | Изолированное тестирование логики |
| Integration test | `httptest.NewServer` | Тестирование полного HTTP-цикла |

---

## 📖 Теория

### Зачем тестировать HTTP?

Веб-приложения — это HTTP-серверы. Когда вы пишете API, вам нужно проверить:
- Правильные ли статус-коды возвращаются?
- Корректный ли формат ответа (JSON)?
- Обрабатываются ли ошибки?
- Работает ли авторизация?

### Проблема: как тестировать без реального сервера?

Запускать настоящий сервер для каждого теста — медленно и сложно. К счастью, Go предоставляет пакет `net/http/httptest` с инструментами для тестирования HTTP без реального сервера.

### Два подхода к HTTP-тестированию

**1. ResponseRecorder (рекомендуется для unit-тестов)**

Вы создаёте "записывающий" ResponseWriter, который сохраняет ответ в память:

```go
// Создаём "записывающий" ответ
recorder := httptest.NewRecorder()

// Создаём запрос
request := httptest.NewRequest("GET", "/api/users", nil)

// Вызываем handler напрямую (без HTTP)
handler.ServeHTTP(recorder, request)

// Проверяем результат
fmt.Println(recorder.Code)        // 200
fmt.Println(recorder.Body.String()) // {"users": [...]}
```

**Преимущества:**
- Очень быстро (нет сетевого взаимодействия)
- Изолированное тестирование handler-логики
- Простая настройка

**2. Test Server (для интеграционных тестов)**

Создаёт настоящий HTTP-сервер на случайном порту:

```go
// Создаём реальный тестовый сервер
server := httptest.NewServer(handler)
defer server.Close()

// Делаем реальный HTTP-запрос
resp, err := http.Get(server.URL + "/api/users")
```

**Преимущества:**
- Тестирует полный HTTP-стек
- Проверяет маршрутизацию, middleware
- Ближе к реальному использованию

### Что проверять в HTTP-тестах?

1. **Status Code** — правильный ли код ответа?
```go
assert.Equal(t, http.StatusOK, recorder.Code)
assert.Equal(t, http.StatusNotFound, recorder.Code)
assert.Equal(t, http.StatusUnauthorized, recorder.Code)
```

2. **Content-Type** — правильный ли тип контента?
```go
assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
```

3. **Body** — правильное ли тело ответа?
```go
var response map[string]interface{}
json.NewDecoder(recorder.Body).Decode(&response)
assert.Equal(t, "John", response["name"])
```

4. **Headers** — установлены ли нужные заголовки?
```go
assert.NotEmpty(t, recorder.Header().Get("X-Request-ID"))
```

### Тестирование разных HTTP-методов

```go
// GET
req := httptest.NewRequest("GET", "/users", nil)

// POST с телом
body := strings.NewReader(`{"name": "John"}`)
req := httptest.NewRequest("POST", "/users", body)
req.Header.Set("Content-Type", "application/json")

// PUT
req := httptest.NewRequest("PUT", "/users/1", body)

// DELETE
req := httptest.NewRequest("DELETE", "/users/1", nil)
```

### Тестирование с авторизацией

```go
req := httptest.NewRequest("GET", "/api/protected", nil)
req.Header.Set("Authorization", "Bearer test-token")
```

---

## 📋 Синтаксис

### httptest.ResponseRecorder

```go
// Создание записывающего ответа
rr := httptest.NewRecorder()

// Создание запроса
req := httptest.NewRequest("GET", "/path", nil)

// Вызов handler
handler.ServeHTTP(rr, req)

// Проверка результата
rr.Code           // статус код
rr.Body.String()  // тело ответа
rr.Header()       // заголовки
```

### httptest.Server

```go
// Создание тестового сервера
ts := httptest.NewServer(handler)
defer ts.Close()

// URL сервера
ts.URL  // например: "http://127.0.0.1:12345"

// Запросы к серверу
resp, err := http.Get(ts.URL + "/path")
```

---

## 💻 Примеры кода

### Пример 1: Простой тест handler

```go
package api

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Hello, World!",
    })
}

func TestHelloHandler(t *testing.T) {
    // Создаём запрос
    req := httptest.NewRequest("GET", "/hello", nil)
    
    // Создаём ResponseRecorder для записи ответа
    rr := httptest.NewRecorder()
    
    // Вызываем handler
    HelloHandler(rr, req)
    
    // Проверяем статус код
    if rr.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", rr.Code)
    }
    
    // Проверяем Content-Type
    contentType := rr.Header().Get("Content-Type")
    if contentType != "application/json" {
        t.Errorf("expected Content-Type application/json, got %s", contentType)
    }
    
    // Проверяем тело ответа
    expected := `{"message":"Hello, World!"}`
    if strings.TrimSpace(rr.Body.String()) != expected {
        t.Errorf("expected body %s, got %s", expected, rr.Body.String())
    }
}
```

### Пример 2: Тест с параметрами URL

```go
package api

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gorilla/mux"
)

type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    
    user := User{ID: id, Name: "John Doe"}
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func TestGetUserHandler(t *testing.T) {
    // Создаём роутер
    router := mux.NewRouter()
    router.HandleFunc("/users/{id}", GetUserHandler).Methods("GET")
    
    // Создаём запрос
    req := httptest.NewRequest("GET", "/users/123", nil)
    rr := httptest.NewRecorder()
    
    // Выполняем через роутер
    router.ServeHTTP(rr, req)
    
    // Проверяем статус
    if rr.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", rr.Code)
    }
    
    // Парсим ответ
    var user User
    err := json.NewDecoder(rr.Body).Decode(&user)
    if err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }
    
    // Проверяем данные
    if user.ID != "123" {
        t.Errorf("expected ID 123, got %s", user.ID)
    }
}
```

### Пример 3: Тест POST запроса

```go
package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    
    if req.Name == "" || req.Email == "" {
        http.Error(w, "name and email required", http.StatusBadRequest)
        return
    }
    
    resp := CreateUserResponse{
        ID:    1,
        Name:  req.Name,
        Email: req.Email,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resp)
}

func TestCreateUserHandler(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    CreateUserRequest
        expectedStatus int
        checkResponse  bool
    }{
        {
            name:           "valid request",
            requestBody:    CreateUserRequest{Name: "John", Email: "john@example.com"},
            expectedStatus: http.StatusCreated,
            checkResponse:  true,
        },
        {
            name:           "missing name",
            requestBody:    CreateUserRequest{Email: "john@example.com"},
            expectedStatus: http.StatusBadRequest,
            checkResponse:  false,
        },
        {
            name:           "missing email",
            requestBody:    CreateUserRequest{Name: "John"},
            expectedStatus: http.StatusBadRequest,
            checkResponse:  false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Сериализуем тело запроса
            body, _ := json.Marshal(tt.requestBody)
            
            // Создаём запрос
            req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            
            rr := httptest.NewRecorder()
            
            CreateUserHandler(rr, req)
            
            // Проверяем статус
            if rr.Code != tt.expectedStatus {
                t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
            }
            
            // Проверяем ответ
            if tt.checkResponse {
                var resp CreateUserResponse
                json.NewDecoder(rr.Body).Decode(&resp)
                
                if resp.Name != tt.requestBody.Name {
                    t.Errorf("expected name %s, got %s", tt.requestBody.Name, resp.Name)
                }
            }
        })
    }
}
```

### Пример 4: Тест с httptest.Server

```go
package api

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestWithServer(t *testing.T) {
    // Создаём handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "status": "ok",
        })
    })
    
    // Создаём тестовый сервер
    server := httptest.NewServer(handler)
    defer server.Close()
    
    // Делаем реальный HTTP-запрос
    resp, err := http.Get(server.URL + "/health")
    if err != nil {
        t.Fatalf("failed to make request: %v", err)
    }
    defer resp.Body.Close()
    
    // Проверяем статус
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %d", resp.StatusCode)
    }
    
    // Парсим ответ
    var result map[string]string
    json.NewDecoder(resp.Body).Decode(&result)
    
    if result["status"] != "ok" {
        t.Errorf("expected status ok, got %s", result["status"])
    }
}
```

### Пример 5: Тестирование middleware

```go
package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

// Middleware для проверки авторизации
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        if token == "" {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        
        if token != "Bearer valid-token" {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

func TestAuthMiddleware(t *testing.T) {
    // Создаём защищённый handler
    handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("success"))
    }))
    
    tests := []struct {
        name           string
        authHeader     string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "no token",
            authHeader:     "",
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   "unauthorized",
        },
        {
            name:           "invalid token",
            authHeader:     "Bearer invalid-token",
            expectedStatus: http.StatusForbidden,
            expectedBody:   "forbidden",
        },
        {
            name:           "valid token",
            authHeader:     "Bearer valid-token",
            expectedStatus: http.StatusOK,
            expectedBody:   "success",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/protected", nil)
            if tt.authHeader != "" {
                req.Header.Set("Authorization", tt.authHeader)
            }
            
            rr := httptest.NewRecorder()
            handler.ServeHTTP(rr, req)
            
            if rr.Code != tt.expectedStatus {
                t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
            }
            
            if !strings.Contains(rr.Body.String(), tt.expectedBody) {
                t.Errorf("expected body to contain %s, got %s", 
                    tt.expectedBody, rr.Body.String())
            }
        })
    }
}
```

### Пример 6: Тестирование API с моком

```go
package api

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

type UserService interface {
    GetByID(id int) (*User, error)
}

type UserHandler struct {
    service UserService
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    // Получаем ID из URL
    id := 1  // упрощённо
    
    user, err := h.service.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// Мок сервиса
type MockUserService struct {
    user *User
    err  error
}

func (m *MockUserService) GetByID(id int) (*User, error) {
    return m.user, m.err
}

func TestUserHandlerGetUser(t *testing.T) {
    t.Run("user found", func(t *testing.T) {
        mockService := &MockUserService{
            user: &User{ID: "1", Name: "John"},
        }
        
        handler := &UserHandler{service: mockService}
        
        req := httptest.NewRequest("GET", "/users/1", nil)
        rr := httptest.NewRecorder()
        
        handler.GetUser(rr, req)
        
        if rr.Code != http.StatusOK {
            t.Errorf("expected 200, got %d", rr.Code)
        }
    })
    
    t.Run("user not found", func(t *testing.T) {
        mockService := &MockUserService{
            err: errors.New("not found"),
        }
        
        handler := &UserHandler{service: mockService}
        
        req := httptest.NewRequest("GET", "/users/999", nil)
        rr := httptest.NewRecorder()
        
        handler.GetUser(rr, req)
        
        if rr.Code != http.StatusNotFound {
            t.Errorf("expected 404, got %d", rr.Code)
        }
    })
}
```

### Пример 7: Тест загрузки файла

```go
package upload

import (
    "bytes"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "testing"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    // Парсим multipart form
    err := r.ParseMultipartForm(10 << 20)  // 10 MB
    if err != nil {
        http.Error(w, "failed to parse form", http.StatusBadRequest)
        return
    }
    
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "failed to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    w.Write([]byte(fmt.Sprintf("Uploaded: %s", header.Filename)))
}

func TestUploadHandler(t *testing.T) {
    // Создаём multipart form
    var buf bytes.Buffer
    writer := multipart.NewWriter(&buf)
    
    // Добавляем файл
    part, err := writer.CreateFormFile("file", "test.txt")
    if err != nil {
        t.Fatal(err)
    }
    part.Write([]byte("file content"))
    
    writer.Close()
    
    // Создаём запрос
    req := httptest.NewRequest("POST", "/upload", &buf)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    rr := httptest.NewRecorder()
    
    UploadHandler(rr, req)
    
    if rr.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rr.Code)
    }
    
    if !strings.Contains(rr.Body.String(), "test.txt") {
        t.Errorf("expected filename in response, got %s", rr.Body.String())
    }
}
```

### Пример 8: Мок внешнего API

```go
package client

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

type APIClient struct {
    baseURL string
    client  *http.Client
}

func (c *APIClient) GetUser(id int) (*User, error) {
    resp, err := c.client.Get(fmt.Sprintf("%s/users/%d", c.baseURL, id))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var user User
    json.NewDecoder(resp.Body).Decode(&user)
    return &user, nil
}

func TestAPIClient(t *testing.T) {
    // Создаём мок-сервер
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Проверяем путь
        if r.URL.Path != "/users/1" {
            t.Errorf("unexpected path: %s", r.URL.Path)
        }
        
        // Возвращаем мок-данные
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(User{
            ID:   "1",
            Name: "Mock User",
        })
    }))
    defer server.Close()
    
    // Создаём клиент с URL мок-сервера
    client := &APIClient{
        baseURL: server.URL,
        client:  http.DefaultClient,
    }
    
    // Тестируем
    user, err := client.GetUser(1)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if user.Name != "Mock User" {
        t.Errorf("expected Mock User, got %s", user.Name)
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Не закрыли тело ответа

```go
// ❌ ПЛОХО — утечка ресурсов
resp, _ := http.Get(server.URL)
// забыли resp.Body.Close()

// ✅ ХОРОШО
resp, _ := http.Get(server.URL)
defer resp.Body.Close()
```

### 2. Не закрыли тестовый сервер

```go
// ❌ ПЛОХО
server := httptest.NewServer(handler)
// забыли server.Close()

// ✅ ХОРОШО
server := httptest.NewServer(handler)
defer server.Close()
```

### 3. Не установили Content-Type

```go
// ❌ ПЛОХО — JSON без Content-Type
req := httptest.NewRequest("POST", "/api", bytes.NewBuffer(jsonBody))

// ✅ ХОРОШО
req := httptest.NewRequest("POST", "/api", bytes.NewBuffer(jsonBody))
req.Header.Set("Content-Type", "application/json")
```

---

## 🏋️ Практические задания

### Задание 1: Тест GET endpoint

Напишите тест для простого GET handler.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

```go
// handler_test.go
package api

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHealthHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/health", nil)
    rr := httptest.NewRecorder()
    
    HealthHandler(rr, req)
    
    // Проверьте статус код
    if rr.Code != http.StatusOK {
        t.Errorf("expected 200, got %d", rr.Code)
    }
    
    // Проверьте Content-Type
    contentType := rr.Header().Get("Content-Type")
    if contentType != "application/json" {
        t.Errorf("expected application/json, got %s", contentType)
    }
    
    // Проверьте тело ответа
    var resp Response
    json.NewDecoder(rr.Body).Decode(&resp)
    
    if resp.Message != "OK" {
        t.Errorf("expected OK, got %s", resp.Message)
    }
}
```

**Ожидаемый вывод:**
```
=== RUN   TestHealthHandler
--- PASS: TestHealthHandler (0.00s)
PASS
```

**Баллы:** 10

### Задание 2: Тест POST с JSON

Протестируйте handler, принимающий JSON.

**Начальный код:**
```go
package main

import (
    "encoding/json"
    "fmt"
)

// TODO: Работа с JSON

func main() {
    // Ваш код здесь
    
}
```

```go
// user_handler_test.go
package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateUserHandler(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        body := `{"name": "John", "email": "john@example.com"}`
        req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(body))
        req.Header.Set("Content-Type", "application/json")
        
        rr := httptest.NewRecorder()
        CreateUserHandler(rr, req)
        
        if rr.Code != http.StatusCreated {
            t.Errorf("expected 201, got %d", rr.Code)
        }
        
        var user User
        json.NewDecoder(rr.Body).Decode(&user)
        
        if user.Name != "John" {
            t.Errorf("expected John, got %s", user.Name)
        }
    })
    
    t.Run("invalid json", func(t *testing.T) {
        req := httptest.NewRequest("POST", "/users", bytes.NewBufferString("invalid"))
        rr := httptest.NewRecorder()
        
        CreateUserHandler(rr, req)
        
        if rr.Code != http.StatusBadRequest {
            t.Errorf("expected 400, got %d", rr.Code)
        }
    })
    
    // Добавьте тест для пустых полей
    // Добавьте тест для неправильного метода
}
```

**Ожидаемый вывод:**
```
=== RUN   TestCreateUserHandler
=== RUN   TestCreateUserHandler/success
=== RUN   TestCreateUserHandler/invalid_json
=== RUN   TestCreateUserHandler/empty_fields
=== RUN   TestCreateUserHandler/wrong_method
--- PASS: TestCreateUserHandler (0.00s)
PASS
```

**Баллы:** 15

### Задание 3: Тест с httptest.Server

Протестируйте HTTP клиент с использованием тестового сервера.

**Начальный код:**
```go
package main

import (
    "fmt"
    "net/http"
)

// TODO: Реализуйте HTTP обработчик

func main() {
    // Ваш код здесь
    
}
```

```go
// client_test.go
package client

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAPIClient(t *testing.T) {
    // Создаём тестовый сервер
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/users/1" {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(User{ID: 1, Name: "Test User"})
            return
        }
        w.WriteHeader(http.StatusNotFound)
    }))
    defer server.Close()
    
    client := &APIClient{
        BaseURL: server.URL,
        Client:  http.DefaultClient,
    }
    
    t.Run("user exists", func(t *testing.T) {
        user, err := client.GetUser(1)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if user.Name != "Test User" {
            t.Errorf("expected Test User, got %s", user.Name)
        }
    })
    
    t.Run("user not found", func(t *testing.T) {
        _, err := client.GetUser(999)
        if err == nil {
            t.Error("expected error, got nil")
        }
    })
}
```

**Ожидаемый вывод:**
```
=== RUN   TestAPIClient
=== RUN   TestAPIClient/user_exists
=== RUN   TestAPIClient/user_not_found
--- PASS: TestAPIClient (0.00s)
PASS
```

**Баллы:** 15

### Задание 4: Тест middleware

Протестируйте middleware для авторизации.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

```go
// middleware_test.go
package api

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAuthMiddleware(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("success"))
    })
    
    protected := AuthMiddleware(handler)
    
    tests := []struct {
        name       string
        token      string
        wantStatus int
    }{
        {"valid token", "Bearer valid-token", http.StatusOK},
        {"no token", "", http.StatusUnauthorized},
        {"invalid token", "Bearer wrong", http.StatusForbidden},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/", nil)
            if tt.token != "" {
                req.Header.Set("Authorization", tt.token)
            }
            
            rr := httptest.NewRecorder()
            protected.ServeHTTP(rr, req)
            
            if rr.Code != tt.wantStatus {
                t.Errorf("expected %d, got %d", tt.wantStatus, rr.Code)
            }
        })
    }
}
```

**Ожидаемый вывод:**
```
=== RUN   TestAuthMiddleware
=== RUN   TestAuthMiddleware/valid_token
=== RUN   TestAuthMiddleware/no_token
=== RUN   TestAuthMiddleware/invalid_token
--- PASS: TestAuthMiddleware (0.00s)
PASS
```

**Баллы:** 10

### Задание 5: Тест query параметров

Протестируйте handler с query параметрами и пагинацией.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

```go
// list_handler_test.go
func TestListUsersHandler(t *testing.T) {
    tests := []struct {
        name      string
        query     string
        wantPage  int
        wantLimit int
    }{
        {"default", "", 1, 10},
        {"with page", "?page=5", 5, 10},
        {"with limit", "?limit=20", 1, 20},
        {"invalid page", "?page=-1", 1, 10},
        {"limit too high", "?limit=500", 1, 10},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/users"+tt.query, nil)
            rr := httptest.NewRecorder()
            
            ListUsersHandler(rr, req)
            
            var resp ListResponse
            json.NewDecoder(rr.Body).Decode(&resp)
            
            if resp.Page != tt.wantPage {
                t.Errorf("page: expected %d, got %d", tt.wantPage, resp.Page)
            }
            if resp.Limit != tt.wantLimit {
                t.Errorf("limit: expected %d, got %d", tt.wantLimit, resp.Limit)
            }
        })
    }
}
```

**Ожидаемый вывод:**
```
=== RUN   TestListUsersHandler
=== RUN   TestListUsersHandler/default
=== RUN   TestListUsersHandler/with_page
=== RUN   TestListUsersHandler/with_limit
=== RUN   TestListUsersHandler/invalid_page
=== RUN   TestListUsersHandler/limit_too_high
--- PASS: TestListUsersHandler (0.00s)
PASS
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [httptest Package](https://pkg.go.dev/net/http/httptest)
- [Testing HTTP Handlers in Go](https://blog.questionable.services/article/testing-http-handlers-go/)
