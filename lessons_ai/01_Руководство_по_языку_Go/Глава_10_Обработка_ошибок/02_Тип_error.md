# Тип error

## 💡 Ключевые идеи

1. **error** — встроенный интерфейс с одним методом `Error() string`
2. **errors.New** — создаёт простую ошибку из строки
3. **fmt.Errorf** — создаёт форматированную ошибку
4. **Кастомные ошибки** — любой тип с методом `Error()` реализует интерфейс error
5. **Обёртка ошибок** — `%w` в fmt.Errorf для сохранения цепочки
6. **errors.Is/As** — проверка и извлечение типов ошибок

---

## 📋 Синтаксис

### Интерфейс error

```go
type error interface {
    Error() string
}
```

### Создание ошибок

```go
// Простая ошибка
err := errors.New("something went wrong")

// Форматированная ошибка
err := fmt.Errorf("failed to process %d items", count)

// С обёрткой другой ошибки
err := fmt.Errorf("operation failed: %w", originalErr)
```

### Кастомный тип ошибки

```go
type MyError struct {
    Code    int
    Message string
}

func (e MyError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}
```

---

## 💻 Примеры кода

### Использование errors.New

```go
package main

import (
    "errors"
    "fmt"
)

var (
    ErrDivisionByZero = errors.New("division by zero")
    ErrNegativeNumber = errors.New("negative number not allowed")
)

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, ErrDivisionByZero
    }
    return a / b, nil
}

func sqrt(x float64) (float64, error) {
    if x < 0 {
        return 0, ErrNegativeNumber
    }
    // Упрощённый расчёт
    return x * 0.5, nil
}

func main() {
    _, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
    }
    
    // Сравнение ошибок
    if err == ErrDivisionByZero {
        fmt.Println("Cannot divide by zero!")
    }
}
```

### Использование fmt.Errorf

```go
package main

import (
    "fmt"
)

func processOrder(orderID int) error {
    if orderID <= 0 {
        return fmt.Errorf("invalid order ID: %d", orderID)
    }
    
    // Обработка заказа...
    return nil
}

func validateUser(userID int, name string) error {
    if userID <= 0 {
        return fmt.Errorf("invalid user ID: %d", userID)
    }
    if name == "" {
        return fmt.Errorf("user %d: name is required", userID)
    }
    return nil
}

func main() {
    err := processOrder(-1)
    if err != nil {
        fmt.Println(err)  // invalid order ID: -1
    }
    
    err = validateUser(0, "Alice")
    if err != nil {
        fmt.Println(err)  // invalid user ID: 0
    }
}
```

### Кастомный тип ошибки

```go
package main

import (
    "fmt"
)

// Кастомный тип ошибки
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on '%s': %s", e.Field, e.Message)
}

// Ещё один тип ошибки
type NotFoundError struct {
    Resource string
    ID       int
}

func (e NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

func validateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field:   "age",
            Message: "cannot be negative",
        }
    }
    if age > 150 {
        return ValidationError{
            Field:   "age",
            Message: "unrealistic value",
        }
    }
    return nil
}

func getUser(id int) (*struct{ Name string }, error) {
    if id == 0 {
        return nil, NotFoundError{Resource: "User", ID: id}
    }
    return &struct{ Name string }{"Alice"}, nil
}

func main() {
    err := validateAge(-5)
    if err != nil {
        fmt.Println(err)  // validation error on 'age': cannot be negative
    }
    
    _, err = getUser(0)
    if err != nil {
        fmt.Println(err)  // User with ID 0 not found
    }
}
```

### Обёртка ошибок (Error Wrapping)

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

func readConfig(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        // %w сохраняет оригинальную ошибку
        return nil, fmt.Errorf("readConfig: %w", err)
    }
    return data, nil
}

func loadApp() error {
    _, err := readConfig("config.json")
    if err != nil {
        return fmt.Errorf("loadApp: %w", err)
    }
    return nil
}

func main() {
    err := loadApp()
    if err != nil {
        fmt.Println("Error:", err)
        // Error: loadApp: readConfig: open config.json: no such file or directory
        
        // Проверяем оригинальную ошибку
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("File does not exist!")
        }
    }
}
```

### errors.Is — проверка типа ошибки

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

var ErrPermissionDenied = errors.New("permission denied")

func accessResource() error {
    return fmt.Errorf("cannot access resource: %w", ErrPermissionDenied)
}

func main() {
    err := accessResource()
    
    // errors.Is проверяет всю цепочку обёрток
    if errors.Is(err, ErrPermissionDenied) {
        fmt.Println("Access denied!")
    }
    
    // Также работает с системными ошибками
    _, err = os.Open("nonexistent.txt")
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("File not found")
    }
}
```

### errors.As — извлечение типа ошибки

```go
package main

import (
    "errors"
    "fmt"
)

type HTTPError struct {
    Code    int
    Message string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

func fetchData() error {
    return fmt.Errorf("fetch failed: %w", &HTTPError{
        Code:    404,
        Message: "Not Found",
    })
}

func main() {
    err := fetchData()
    
    var httpErr *HTTPError
    if errors.As(err, &httpErr) {
        fmt.Printf("HTTP Error Code: %d\n", httpErr.Code)
        fmt.Printf("HTTP Error Message: %s\n", httpErr.Message)
    }
}
```

### Sentinel Errors (ошибки-маркеры)

```go
package main

import (
    "errors"
    "fmt"
)

// Sentinel errors — предопределённые ошибки
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrBadRequest   = errors.New("bad request")
)

func getResource(id int, hasAccess bool) (string, error) {
    if id <= 0 {
        return "", ErrBadRequest
    }
    if !hasAccess {
        return "", ErrUnauthorized
    }
    if id > 100 {
        return "", ErrNotFound
    }
    return fmt.Sprintf("Resource %d", id), nil
}

func main() {
    _, err := getResource(-1, true)
    
    switch {
    case errors.Is(err, ErrBadRequest):
        fmt.Println("Invalid request parameters")
    case errors.Is(err, ErrUnauthorized):
        fmt.Println("Please login first")
    case errors.Is(err, ErrNotFound):
        fmt.Println("Resource doesn't exist")
    case err != nil:
        fmt.Println("Unknown error:", err)
    }
}
```

### Практический пример: API клиент

```go
package main

import (
    "errors"
    "fmt"
)

// Типы ошибок
type APIError struct {
    StatusCode int
    Message    string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

type NetworkError struct {
    Reason string
}

func (e *NetworkError) Error() string {
    return fmt.Sprintf("network error: %s", e.Reason)
}

// Sentinel errors
var (
    ErrTimeout = errors.New("request timeout")
    ErrRateLimit = errors.New("rate limit exceeded")
)

// API функция
func callAPI(endpoint string) (string, error) {
    switch endpoint {
    case "/timeout":
        return "", fmt.Errorf("calling %s: %w", endpoint, ErrTimeout)
    case "/error":
        return "", &APIError{StatusCode: 500, Message: "Internal Server Error"}
    case "/network":
        return "", &NetworkError{Reason: "connection refused"}
    default:
        return "Success!", nil
    }
}

func handleError(err error) {
    var apiErr *APIError
    var netErr *NetworkError
    
    switch {
    case errors.Is(err, ErrTimeout):
        fmt.Println("→ Retry later (timeout)")
    case errors.Is(err, ErrRateLimit):
        fmt.Println("→ Slow down requests")
    case errors.As(err, &apiErr):
        fmt.Printf("→ API returned %d: %s\n", apiErr.StatusCode, apiErr.Message)
    case errors.As(err, &netErr):
        fmt.Printf("→ Check network: %s\n", netErr.Reason)
    default:
        fmt.Println("→ Unknown error:", err)
    }
}

func main() {
    endpoints := []string{"/timeout", "/error", "/network", "/success"}
    
    for _, ep := range endpoints {
        result, err := callAPI(ep)
        if err != nil {
            fmt.Printf("Error on %s: %v\n", ep, err)
            handleError(err)
        } else {
            fmt.Printf("Success on %s: %s\n", ep, result)
        }
        fmt.Println()
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Сравнение ошибок через ==

```go
// ❌ Не работает для обёрнутых ошибок
if err == ErrNotFound {
    // может не сработать!
}

// ✅ Используйте errors.Is
if errors.Is(err, ErrNotFound) {
    // работает с обёртками
}
```

### 2. Type assertion без errors.As

```go
// ❌ Не работает для обёрнутых ошибок
if e, ok := err.(*MyError); ok {
    // может не сработать!
}

// ✅ Используйте errors.As
var e *MyError
if errors.As(err, &e) {
    // работает с обёртками
}
```

### 3. Забыли %w при обёртке

```go
// ❌ Теряет оригинальную ошибку
return fmt.Errorf("failed: %v", err)

// ✅ Сохраняет ошибку для errors.Is/As
return fmt.Errorf("failed: %w", err)
```

### 4. Изменение sentinel error

```go
// ❌ Sentinel errors должны быть неизменяемыми
var ErrNotFound = errors.New("not found")
// где-то в коде: ErrNotFound = errors.New("другое")

// ✅ Используйте константы или приватные переменные
var errNotFound = errors.New("not found")  // приватная
```

---

## 📝 Практика

### Задача 1: Кастомная ошибка
Создайте `FileError` с полями Path и Operation.

### Задача 2: HTTP ошибки
Создайте типы ошибок для разных HTTP статусов (4xx, 5xx).

### Задача 3: Цепочка ошибок
Реализуйте функции с обёрткой ошибок и используйте errors.Is/As.

### Задача 4: Sentinel errors
Создайте набор sentinel errors для базы данных (ErrDuplicate, ErrConstraint).

### Задача 5: Unwrap
Реализуйте метод Unwrap() для кастомной ошибки.

### Задача 6: Множественные ошибки
Реализуйте тип MultiError, собирающий несколько ошибок.

### Задача 7: Retry logic
Используйте типы ошибок для определения, стоит ли повторять операцию.

### Задача 8: Error handling middleware
Создайте обработчик ошибок, форматирующий разные типы по-разному.
