# Context Values

---

## 💡 Ключевые идеи

1. **Context Values** — передача request-scoped данных через контекст
2. **Типизированные ключи** — избегание коллизий имён
3. **Иммутабельность** — каждый WithValue создаёт новый контекст
4. **Request-scoped** — только для данных, привязанных к запросу
5. **Не для опциональных параметров** — Values не заменяют аргументы функций

### Когда использовать Values

| ✅ Правильно | ❌ Неправильно |
|-------------|---------------|
| Request ID | Бизнес-логика |
| User ID из токена | Конфигурация |
| Trace ID | Опциональные параметры |
| Deadline info | База данных |

---

## 📖 Теория

### Что такое Context Values?

`context.WithValue` позволяет прикрепить к контексту пару ключ-значение:

```go
ctx := context.WithValue(parentCtx, "key", "value")
value := ctx.Value("key")  // "value"
```

### Когда использовать Values?

Context Values предназначены **ТОЛЬКО** для данных, привязанных к запросу (request-scoped):

✅ **Правильное использование:**
- Request ID для логирования
- User ID из JWT-токена
- Trace ID для distributed tracing
- Локаль пользователя

❌ **Неправильное использование:**
- Подключение к базе данных
- Конфигурация приложения
- Бизнес-логика и параметры функций
- Логгер (хотя это спорно)

### Почему типизированные ключи важны?

Проблема со строковыми ключами:

```go
// Пакет A
ctx = context.WithValue(ctx, "userID", 123)

// Пакет B (не знает о пакете A)
ctx = context.WithValue(ctx, "userID", "admin")  // перезаписал!

// Пакет A
id := ctx.Value("userID").(int)  // PANIC! Это string
```

Решение — **приватные типы ключей**:

```go
// Пакет A
type contextKey int
const userIDKey contextKey = 1
ctx = context.WithValue(ctx, userIDKey, 123)

// Пакет B (не может использовать приватный тип)
// Даже если создаст свой type contextKey int — это ДРУГОЙ тип
```

### Паттерн: функции-хелперы

Вместо прямого доступа к Values создавайте функции:

```go
// Определение в пакете auth
type contextKey int
const userIDKey contextKey = 1

func WithUserID(ctx context.Context, userID int) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(userIDKey).(int)
    return userID, ok
}

// Использование
ctx = auth.WithUserID(ctx, 123)
if userID, ok := auth.UserIDFromContext(ctx); ok {
    fmt.Println("User:", userID)
}
```

Преимущества:
- Инкапсуляция — детали скрыты
- Типобезопасность — функции возвращают конкретные типы
- Проверка наличия — возвращается bool

### Как работает поиск значений?

Контексты образуют цепочку. Поиск значения идёт от текущего к корню:

```
Background
    └── WithValue(userID=123)
            └── WithValue(requestID="abc")
                    └── WithTimeout(5s)  ← текущий ctx
```

При `ctx.Value(userID)`:
1. Проверяет WithTimeout — нет
2. Проверяет WithValue(requestID) — нет
3. Проверяет WithValue(userID) — НАШЁЛ!

### Иммутабельность

Каждый `WithValue` создаёт **новый** контекст:

```go
ctx1 := context.Background()
ctx2 := context.WithValue(ctx1, key1, "value1")
ctx3 := context.WithValue(ctx2, key2, "value2")

// ctx1 не содержит key1 и key2
// ctx2 содержит key1, но не key2
// ctx3 содержит и key1, и key2
```

### Производительность Values

Поиск значения — O(n) от глубины цепочки контекстов. Для большинства случаев это не проблема (глубина обычно < 10), но не храните много значений.

### Альтернативы Context Values

Если данных много или они сложные — рассмотрите:

1. **Явные параметры функций:**
```go
func ProcessOrder(ctx context.Context, userID int, order Order) error
```

2. **Структура запроса:**
```go
type Request struct {
    UserID    int
    RequestID string
    Locale    string
}
func Process(ctx context.Context, req Request) error
```

---

## 💻 Примеры кода

### Пример 1: Типизированные ключи

```go
package main

import (
    "context"
    "fmt"
)

// ПРАВИЛЬНО: приватный тип для ключей
type contextKey int

const (
    userIDKey contextKey = iota
    requestIDKey
    traceIDKey
)

// Функции-хелперы для работы с контекстом
func WithUserID(ctx context.Context, userID int) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(userIDKey).(int)
    return userID, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, requestIDKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
    if requestID, ok := ctx.Value(requestIDKey).(string); ok {
        return requestID
    }
    return ""
}

func processRequest(ctx context.Context) {
    userID, ok := UserIDFromContext(ctx)
    if !ok {
        fmt.Println("User ID не найден")
        return
    }
    
    requestID := RequestIDFromContext(ctx)
    
    fmt.Printf("[%s] Обработка запроса для пользователя %d\n", requestID, userID)
}

func main() {
    ctx := context.Background()
    ctx = WithUserID(ctx, 42)
    ctx = WithRequestID(ctx, "req-abc-123")
    
    processRequest(ctx)
}
```

### Пример 2: Middleware для HTTP

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    
    "github.com/google/uuid"
)

type ctxKey string

const (
    requestIDCtxKey ctxKey = "requestID"
    userCtxKey      ctxKey = "user"
)

type User struct {
    ID    int
    Email string
    Role  string
}

// Middleware: добавляет Request ID
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        ctx := context.WithValue(r.Context(), requestIDCtxKey, requestID)
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Middleware: добавляет User из токена
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        // Упрощённая "авторизация"
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // В реальности — валидация JWT
        user := User{ID: 1, Email: "user@example.com", Role: "admin"}
        
        ctx := context.WithValue(r.Context(), userCtxKey, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Хелперы
func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(requestIDCtxKey).(string); ok {
        return id
    }
    return "unknown"
}

func GetUser(ctx context.Context) (User, bool) {
    user, ok := ctx.Value(userCtxKey).(User)
    return user, ok
}

// Handler
func profileHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    requestID := GetRequestID(ctx)
    user, ok := GetUser(ctx)
    if !ok {
        http.Error(w, "User not found", http.StatusInternalServerError)
        return
    }
    
    fmt.Printf("[%s] Profile request for user %d\n", requestID, user.ID)
    fmt.Fprintf(w, "Hello, %s! (Request: %s)", user.Email, requestID)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/profile", profileHandler)
    
    // Цепочка middleware
    handler := RequestIDMiddleware(AuthMiddleware(mux))
    
    fmt.Println("Server on :8080")
    http.ListenAndServe(":8080", handler)
}
```

### Пример 3: Логирование с контекстом

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
)

type loggerKey struct{}

type Logger struct {
    requestID string
    userID    int
}

func (l *Logger) Info(msg string) {
    log.Printf("[INFO] [req:%s] [user:%d] %s", l.requestID, l.userID, msg)
}

func (l *Logger) Error(msg string, err error) {
    log.Printf("[ERROR] [req:%s] [user:%d] %s: %v", l.requestID, l.userID, msg, err)
}

func WithLogger(ctx context.Context, requestID string, userID int) context.Context {
    logger := &Logger{requestID: requestID, userID: userID}
    return context.WithValue(ctx, loggerKey{}, logger)
}

func LoggerFromContext(ctx context.Context) *Logger {
    if logger, ok := ctx.Value(loggerKey{}).(*Logger); ok {
        return logger
    }
    return &Logger{requestID: "unknown", userID: 0}
}

// Сервис
func processOrder(ctx context.Context, orderID int) error {
    logger := LoggerFromContext(ctx)
    
    logger.Info(fmt.Sprintf("Начало обработки заказа %d", orderID))
    
    time.Sleep(100 * time.Millisecond)
    
    logger.Info(fmt.Sprintf("Заказ %d обработан", orderID))
    
    return nil
}

func main() {
    ctx := context.Background()
    ctx = WithLogger(ctx, "req-123", 42)
    
    processOrder(ctx, 1001)
}
```

### Пример 4: Tracing

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type traceKey struct{}

type TraceInfo struct {
    TraceID    string
    SpanID     string
    ParentSpan string
    StartTime  time.Time
}

func StartTrace(ctx context.Context, traceID string) context.Context {
    trace := &TraceInfo{
        TraceID:   traceID,
        SpanID:    generateSpanID(),
        StartTime: time.Now(),
    }
    return context.WithValue(ctx, traceKey{}, trace)
}

func StartSpan(ctx context.Context, name string) (context.Context, func()) {
    parent := TraceFromContext(ctx)
    
    span := &TraceInfo{
        TraceID:    parent.TraceID,
        SpanID:     generateSpanID(),
        ParentSpan: parent.SpanID,
        StartTime:  time.Now(),
    }
    
    fmt.Printf("[TRACE] Start span '%s' (trace=%s, span=%s, parent=%s)\n",
        name, span.TraceID, span.SpanID, span.ParentSpan)
    
    newCtx := context.WithValue(ctx, traceKey{}, span)
    
    end := func() {
        duration := time.Since(span.StartTime)
        fmt.Printf("[TRACE] End span '%s' (duration=%v)\n", name, duration)
    }
    
    return newCtx, end
}

func TraceFromContext(ctx context.Context) *TraceInfo {
    if trace, ok := ctx.Value(traceKey{}).(*TraceInfo); ok {
        return trace
    }
    return &TraceInfo{TraceID: "unknown", SpanID: "unknown"}
}

func generateSpanID() string {
    return fmt.Sprintf("span-%d", time.Now().UnixNano()%10000)
}

// Сервисы
func handleRequest(ctx context.Context) {
    ctx, end := StartSpan(ctx, "handleRequest")
    defer end()
    
    fetchData(ctx)
    processData(ctx)
}

func fetchData(ctx context.Context) {
    ctx, end := StartSpan(ctx, "fetchData")
    defer end()
    
    time.Sleep(50 * time.Millisecond)
}

func processData(ctx context.Context) {
    ctx, end := StartSpan(ctx, "processData")
    defer end()
    
    time.Sleep(30 * time.Millisecond)
}

func main() {
    ctx := StartTrace(context.Background(), "trace-abc-123")
    handleRequest(ctx)
}
```

### Пример 5: Tenant/Organization Context

```go
package main

import (
    "context"
    "fmt"
)

type tenantKey struct{}

type Tenant struct {
    ID       string
    Name     string
    Settings map[string]interface{}
}

func WithTenant(ctx context.Context, tenant *Tenant) context.Context {
    return context.WithValue(ctx, tenantKey{}, tenant)
}

func TenantFromContext(ctx context.Context) *Tenant {
    if tenant, ok := ctx.Value(tenantKey{}).(*Tenant); ok {
        return tenant
    }
    return nil
}

// Repository с multi-tenancy
type UserRepository struct{}

func (r *UserRepository) GetUsers(ctx context.Context) ([]string, error) {
    tenant := TenantFromContext(ctx)
    if tenant == nil {
        return nil, fmt.Errorf("tenant not found in context")
    }
    
    // В реальности — фильтрация по tenant_id в БД
    fmt.Printf("Fetching users for tenant: %s\n", tenant.Name)
    
    return []string{"user1@" + tenant.ID, "user2@" + tenant.ID}, nil
}

func main() {
    tenant := &Tenant{
        ID:   "acme",
        Name: "Acme Corp",
        Settings: map[string]interface{}{
            "maxUsers": 100,
        },
    }
    
    ctx := WithTenant(context.Background(), tenant)
    
    repo := &UserRepository{}
    users, _ := repo.GetUsers(ctx)
    fmt.Println("Users:", users)
}
```

### Пример 6: Feature Flags

```go
package main

import (
    "context"
    "fmt"
)

type featuresKey struct{}

type FeatureFlags struct {
    flags map[string]bool
}

func NewFeatureFlags() *FeatureFlags {
    return &FeatureFlags{flags: make(map[string]bool)}
}

func (f *FeatureFlags) Enable(flag string) {
    f.flags[flag] = true
}

func (f *FeatureFlags) IsEnabled(flag string) bool {
    return f.flags[flag]
}

func WithFeatures(ctx context.Context, features *FeatureFlags) context.Context {
    return context.WithValue(ctx, featuresKey{}, features)
}

func IsFeatureEnabled(ctx context.Context, flag string) bool {
    if features, ok := ctx.Value(featuresKey{}).(*FeatureFlags); ok {
        return features.IsEnabled(flag)
    }
    return false
}

// Использование
func processPayment(ctx context.Context, amount float64) {
    if IsFeatureEnabled(ctx, "new_payment_flow") {
        fmt.Println("Using NEW payment flow")
        // новая логика
    } else {
        fmt.Println("Using OLD payment flow")
        // старая логика
    }
}

func main() {
    features := NewFeatureFlags()
    features.Enable("new_payment_flow")
    features.Enable("dark_mode")
    
    ctx := WithFeatures(context.Background(), features)
    
    processPayment(ctx, 99.99)
    
    fmt.Println("Dark mode:", IsFeatureEnabled(ctx, "dark_mode"))
    fmt.Println("Beta feature:", IsFeatureEnabled(ctx, "beta"))
}
```

---

## ⚠️ Частые ошибки

### 1. Строковые ключи

```go
// ❌ ПЛОХО — возможны коллизии
ctx = context.WithValue(ctx, "userID", 42)

// ✅ ХОРОШО — типизированный ключ
type ctxKey int
const userIDKey ctxKey = 0
ctx = context.WithValue(ctx, userIDKey, 42)
```

### 2. Бизнес-логика в Values

```go
// ❌ ПЛОХО — не для бизнес-данных
ctx = context.WithValue(ctx, "order", order)
ctx = context.WithValue(ctx, "products", products)

// ✅ ХОРОШО — передавайте как параметры
func ProcessOrder(ctx context.Context, order Order, products []Product) {}
```

### 3. Мутация данных в контексте

```go
// ❌ ПЛОХО — мутируем данные
user := ctx.Value(userKey).(*User)
user.Name = "New Name"  // изменяем!

// ✅ ХОРОШО — данные иммутабельны
```

### 4. Проверка типа без ok

```go
// ❌ ПЛОХО — panic при неправильном типе
userID := ctx.Value(userIDKey).(int)

// ✅ ХОРОШО — безопасная проверка
userID, ok := ctx.Value(userIDKey).(int)
if !ok {
    return errors.New("userID not found")
}
```

---

## 🏋️ Практические задания

### Задание 1: Request ID Middleware

Создайте middleware, добавляющий уникальный request ID.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 10

### Задание 2: Typed Context Keys

Создайте типобезопасные ключи и хелперы для context.

**Начальный код:**
```go
package main

import (
    "context"
    "fmt"
)

// TODO: Используйте context

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Processing request for user 123 with role admin
```

**Баллы:** 10

### Задание 3: Audit Logger

Реализуйте логгер, извлекающий метаданные из context.

**Начальный код:**
```go
package main

import (
    "context"
    "fmt"
)

// TODO: Используйте context

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
[15:04:05] request=abc-123 user=42: Request started
[15:04:05] request=abc-123 user=42: Processing order
[15:04:05] request=abc-123 user=42: Request completed
```

**Баллы:** 15

### Задание 4: Multi-tenant Repository

Создайте репозиторий с автоматической фильтрацией по tenant.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Tenant A users: [{1 Alice tenant-a} {2 Bob tenant-a}]
Tenant B users: [{3 Charlie tenant-b} {4 Dave tenant-b}]
```

**Баллы:** 15

### Задание 5: Feature Flags

Реализуйте feature flags через context.

**Начальный код:**
```go
package main

import (
    "context"
    "fmt"
)

// TODO: Используйте context

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
=== Regular user ===
Processing 99.99 via old payment system
Showing UI...
  - Dark mode: ON

=== Beta tester ===
Processing 99.99 via NEW payment system
Showing UI...
  - Dark mode: ON
  - Beta features: ON
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Context Values](https://pkg.go.dev/context#WithValue)
- [Context Best Practices](https://go.dev/blog/context-and-structs)
