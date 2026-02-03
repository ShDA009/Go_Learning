# Context Best Practices

---

## 💡 Ключевые идеи

1. **Первый параметр** — Context всегда первый аргумент функции
2. **Не храните в структурах** — передавайте через параметры
3. **Propagation** — пробрасывайте контекст через все слои
4. **Graceful Shutdown** — корректное завершение с context
5. **Testing** — тестирование с context

---

## 📖 Теория

### Золотые правила работы с Context

Пакет context имеет чёткие правила использования, нарушение которых ведёт к багам и утечкам ресурсов.

### Правило 1: Context — первый параметр

```go
// ✅ Правильно
func GetUser(ctx context.Context, id int) (*User, error)
func Process(ctx context.Context, data []byte) error

// ❌ Неправильно
func GetUser(id int, ctx context.Context) (*User, error)
func Process(data []byte, ctx context.Context) error
```

Это не просто конвенция — это стандарт Go. Весь код в стандартной библиотеке следует этому правилу.

### Правило 2: Никогда не передавайте nil

```go
// ❌ ПЛОХО — может вызвать panic
doSomething(nil)

// ✅ ХОРОШО
doSomething(context.Background())
doSomething(context.TODO())
```

Если не знаете, какой контекст использовать — используйте `context.TODO()` как временное решение.

### Правило 3: Не храните Context в структурах

```go
// ❌ ПЛОХО
type Server struct {
    ctx context.Context
}

// ✅ ХОРОШО — передавайте через методы
type Server struct {
    // другие поля
}

func (s *Server) Handle(ctx context.Context, req Request) error {
    // используем ctx
}
```

Почему? Контекст привязан к конкретной операции (запросу). Хранение в структуре привязывает его к времени жизни объекта — это неправильно.

**Исключение:** если структура представляет одну операцию:
```go
type Request struct {
    ctx context.Context  // OK — Request живёт один запрос
    // ...
}
```

### Правило 4: Всегда вызывайте cancel()

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()  // ОБЯЗАТЕЛЬНО

// даже если операция успешна — нужно освободить ресурсы
result, err := doWork(ctx)
// cancel() будет вызван при выходе из функции
```

### Правило 5: Пробрасывайте Context через все слои

```go
func Handler(ctx context.Context, r *http.Request) {
    // Пробрасываем через все слои
    user, err := userService.GetUser(ctx, userID)
    if err != nil { ... }
    
    orders, err := orderService.GetOrders(ctx, user.ID)
    if err != nil { ... }
    
    // ...
}

// В сервисе
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // Пробрасываем в репозиторий
    return s.repo.FindByID(ctx, id)
}

// В репозитории
func (r *UserRepo) FindByID(ctx context.Context, id int) (*User, error) {
    // Пробрасываем в SQL-запрос
    row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)
    // ...
}
```

### Graceful Shutdown с Context

```go
func main() {
    // Создаём контекст, отменяемый при сигнале
    ctx, stop := signal.NotifyContext(context.Background(), 
        syscall.SIGINT, syscall.SIGTERM)
    defer stop()
    
    // Запускаем сервер
    server := &http.Server{Addr: ":8080", Handler: handler}
    
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // Ждём сигнала
    <-ctx.Done()
    
    // Graceful shutdown с таймаутом
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    server.Shutdown(shutdownCtx)
}
```

### Тестирование с Context

```go
func TestSlowOperation(t *testing.T) {
    // Тест с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()
    
    err := slowOperation(ctx)
    if !errors.Is(err, context.DeadlineExceeded) {
        t.Error("expected timeout")
    }
}

func TestCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    
    // Запускаем операцию
    done := make(chan error)
    go func() {
        done <- longOperation(ctx)
    }()
    
    // Отменяем
    cancel()
    
    // Проверяем, что операция отменилась
    err := <-done
    if !errors.Is(err, context.Canceled) {
        t.Error("expected cancellation")
    }
}
```

---

## 📋 Правила использования Context

### ✅ DO (Делайте)

```go
// 1. Context — первый параметр
func GetUser(ctx context.Context, id int) (*User, error)

// 2. Всегда вызывайте cancel
ctx, cancel := context.WithTimeout(ctx, time.Second)
defer cancel()

// 3. Пробрасывайте context
func Handler(ctx context.Context) {
    data := fetchData(ctx)       // пробрасываем
    result := process(ctx, data) // пробрасываем
    save(ctx, result)            // пробрасываем
}

// 4. Проверяйте ctx.Done() в долгих операциях
for i := range items {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        process(items[i])
    }
}
```

### ❌ DON'T (Не делайте)

```go
// 1. Не храните в структурах
type Server struct {
    ctx context.Context  // ❌ ПЛОХО
}

// 2. Не используйте nil context
doSomething(nil)  // ❌ ПЛОХО
doSomething(context.Background())  // ✅ OK

// 3. Не передавайте бизнес-данные через Values
ctx = context.WithValue(ctx, "order", order)  // ❌ ПЛОХО

// 4. Не игнорируйте context в публичных API
func PublicAPI() error  // ❌ Нет context
func PublicAPI(ctx context.Context) error  // ✅ Есть context
```

---

## 💻 Примеры кода

### Пример 1: Правильная структура сервиса

```go
package service

import (
    "context"
    "time"
)

// ❌ ПЛОХО — context в структуре
type BadService struct {
    ctx context.Context
    db  *Database
}

// ✅ ХОРОШО — context в методах
type GoodService struct {
    db      *Database
    timeout time.Duration
}

func NewService(db *Database) *GoodService {
    return &GoodService{
        db:      db,
        timeout: 5 * time.Second,
    }
}

func (s *GoodService) GetUser(ctx context.Context, id int) (*User, error) {
    // Добавляем таймаут к операции
    ctx, cancel := context.WithTimeout(ctx, s.timeout)
    defer cancel()
    
    return s.db.GetUser(ctx, id)
}

func (s *GoodService) CreateUser(ctx context.Context, user *User) error {
    ctx, cancel := context.WithTimeout(ctx, s.timeout)
    defer cancel()
    
    return s.db.CreateUser(ctx, user)
}
```

### Пример 2: Graceful Shutdown

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

type Application struct {
    server  *http.Server
    workers []*Worker
    wg      sync.WaitGroup
}

type Worker struct {
    id   int
    stop chan struct{}
}

func (w *Worker) Run(ctx context.Context, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("Worker %d started\n", w.id)
    
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d stopping...\n", w.id)
            // Cleanup
            time.Sleep(100 * time.Millisecond)
            fmt.Printf("Worker %d stopped\n", w.id)
            return
        default:
            // Работа
            time.Sleep(500 * time.Millisecond)
            fmt.Printf("Worker %d: tick\n", w.id)
        }
    }
}

func (app *Application) Start(ctx context.Context) {
    // Запускаем workers
    for _, worker := range app.workers {
        app.wg.Add(1)
        go worker.Run(ctx, &app.wg)
    }
    
    // Запускаем HTTP сервер
    go func() {
        fmt.Println("Server starting on :8080")
        if err := app.server.ListenAndServe(); err != http.ErrServerClosed {
            fmt.Printf("Server error: %v\n", err)
        }
    }()
}

func (app *Application) Shutdown(ctx context.Context) error {
    fmt.Println("Shutting down...")
    
    // Останавливаем HTTP сервер
    if err := app.server.Shutdown(ctx); err != nil {
        return fmt.Errorf("server shutdown: %w", err)
    }
    
    // Ждём завершения workers
    done := make(chan struct{})
    go func() {
        app.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        fmt.Println("All workers stopped")
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func main() {
    app := &Application{
        server: &http.Server{Addr: ":8080"},
        workers: []*Worker{
            {id: 1},
            {id: 2},
            {id: 3},
        },
    }
    
    // Контекст для работы приложения
    ctx, cancel := context.WithCancel(context.Background())
    
    // Запускаем приложение
    app.Start(ctx)
    
    // Ждём сигнал завершения
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    fmt.Println("\nReceived shutdown signal")
    
    // Отменяем контекст
    cancel()
    
    // Graceful shutdown с таймаутом
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()
    
    if err := app.Shutdown(shutdownCtx); err != nil {
        fmt.Printf("Shutdown error: %v\n", err)
    }
    
    fmt.Println("Application stopped")
}
```

### Пример 3: Пробрасывание через слои

```go
package main

import (
    "context"
    "database/sql"
    "net/http"
)

// Handler layer
func UserHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    userID := 1 // из URL
    
    user, err := userService.GetUser(ctx, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // respond
}

// Service layer
type UserService struct {
    repo *UserRepository
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // Добавляем бизнес-логику
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Дополнительные операции
    if err := s.enrichUser(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *UserService) enrichUser(ctx context.Context, user *User) error {
    // Context пробрасывается дальше
    return nil
}

// Repository layer
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
    // Используем ctx в запросе к БД
    row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)
    
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    return &user, err
}
```

### Пример 4: Тестирование с Context

```go
package service

import (
    "context"
    "testing"
    "time"
)

func TestGetUser(t *testing.T) {
    service := NewUserService(mockRepo)
    
    t.Run("success", func(t *testing.T) {
        ctx := context.Background()
        
        user, err := service.GetUser(ctx, 1)
        
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if user.ID != 1 {
            t.Errorf("expected ID 1, got %d", user.ID)
        }
    })
    
    t.Run("timeout", func(t *testing.T) {
        // Контекст с очень коротким таймаутом
        ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
        defer cancel()
        
        time.Sleep(time.Millisecond)
        
        _, err := service.GetUser(ctx, 1)
        
        if err != context.DeadlineExceeded {
            t.Errorf("expected DeadlineExceeded, got %v", err)
        }
    })
    
    t.Run("cancellation", func(t *testing.T) {
        ctx, cancel := context.WithCancel(context.Background())
        
        // Отменяем сразу
        cancel()
        
        _, err := service.GetUser(ctx, 1)
        
        if err != context.Canceled {
            t.Errorf("expected Canceled, got %v", err)
        }
    })
}

func TestSlowOperation(t *testing.T) {
    t.Run("respects context", func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
        defer cancel()
        
        start := time.Now()
        err := slowOperation(ctx)  // должна занять 1 секунду
        duration := time.Since(start)
        
        // Должна прерваться раньше
        if duration > 200*time.Millisecond {
            t.Errorf("operation took too long: %v", duration)
        }
        
        if err != context.DeadlineExceeded {
            t.Errorf("expected DeadlineExceeded, got %v", err)
        }
    })
}
```

### Пример 5: Context в middleware цепочке

```go
package middleware

import (
    "context"
    "net/http"
    "time"
)

// Timeout middleware
func Timeout(duration time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), duration)
            defer cancel()
            
            // Заменяем context в request
            r = r.WithContext(ctx)
            
            done := make(chan struct{})
            go func() {
                next.ServeHTTP(w, r)
                close(done)
            }()
            
            select {
            case <-done:
                // Нормальное завершение
            case <-ctx.Done():
                http.Error(w, "Request Timeout", http.StatusGatewayTimeout)
            }
        })
    }
}

// Recovery middleware с context
func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                // Логируем с context данными
                requestID := GetRequestID(r.Context())
                log.Printf("[%s] Panic: %v", requestID, err)
                
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

// Использование
func main() {
    handler := http.HandlerFunc(myHandler)
    
    // Цепочка middleware
    withMiddleware := Recovery(
        Timeout(30 * time.Second)(
            RequestID(
                Auth(handler),
            ),
        ),
    )
    
    http.ListenAndServe(":8080", withMiddleware)
}
```

### Пример 6: Error handling с Context

```go
package main

import (
    "context"
    "errors"
    "fmt"
)

// Кастомные ошибки
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
)

func processRequest(ctx context.Context, id int) error {
    // Проверяем контекст в начале
    if err := ctx.Err(); err != nil {
        return fmt.Errorf("context error before processing: %w", err)
    }
    
    // Получаем данные
    data, err := fetchData(ctx, id)
    if err != nil {
        // Оборачиваем ошибку с контекстом
        if errors.Is(err, context.Canceled) {
            return fmt.Errorf("request was cancelled: %w", err)
        }
        if errors.Is(err, context.DeadlineExceeded) {
            return fmt.Errorf("request timeout: %w", err)
        }
        return fmt.Errorf("fetch failed: %w", err)
    }
    
    // Обрабатываем данные
    result, err := processData(ctx, data)
    if err != nil {
        return fmt.Errorf("process failed: %w", err)
    }
    
    // Сохраняем результат
    if err := saveResult(ctx, result); err != nil {
        return fmt.Errorf("save failed: %w", err)
    }
    
    return nil
}

// Проверка причины ошибки
func handleError(err error) {
    switch {
    case errors.Is(err, context.Canceled):
        fmt.Println("Request was cancelled by client")
    case errors.Is(err, context.DeadlineExceeded):
        fmt.Println("Request timed out")
    case errors.Is(err, ErrNotFound):
        fmt.Println("Resource not found")
    case errors.Is(err, ErrUnauthorized):
        fmt.Println("Unauthorized access")
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
}
```

---

## 📋 Чек-лист

### При создании API

- [ ] Context — первый параметр
- [ ] Функция поддерживает отмену
- [ ] Документация описывает поведение при отмене
- [ ] Таймаут по умолчанию разумный

### При использовании Context

- [ ] Вызвал `defer cancel()`
- [ ] Пробросил context во все вызовы
- [ ] Проверяю `ctx.Done()` в долгих операциях
- [ ] Не храню context в структурах

### При тестировании

- [ ] Тест с успешным выполнением
- [ ] Тест с таймаутом
- [ ] Тест с отменой
- [ ] Тест graceful shutdown

---

## 🏋️ Практические задания

### Задание 1: Рефакторинг под Context

Добавьте поддержку context в существующий код.

**Начальный код (до рефакторинга):**
```go
package main

import (
    "fmt"
    "time"
)

// ДО: без context
func fetchUser(id int) (string, error) {
    time.Sleep(100 * time.Millisecond)
    return fmt.Sprintf("User-%d", id), nil
}

func fetchOrders(userID int) ([]string, error) {
    time.Sleep(200 * time.Millisecond)
    return []string{"Order-1", "Order-2"}, nil
}

func processRequest(userID int) error {
    user, err := fetchUser(userID)
    if err != nil {
        return err
    }
    
    orders, err := fetchOrders(userID)
    if err != nil {
        return err
    }
    
    fmt.Printf("User: %s, Orders: %v\n", user, orders)
    return nil
}
```

**После рефакторинга:**
```go
package main

import (
    "context"
    "fmt"
    "time"
)

// ПОСЛЕ: с context
func fetchUser(ctx context.Context, id int) (string, error) {
    select {
    case <-time.After(100 * time.Millisecond):
        return fmt.Sprintf("User-%d", id), nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func fetchOrders(ctx context.Context, userID int) ([]string, error) {
    select {
    case <-time.After(200 * time.Millisecond):
        return []string{"Order-1", "Order-2"}, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func processRequest(ctx context.Context, userID int) error {
    user, err := fetchUser(ctx, userID)
    if err != nil {
        return fmt.Errorf("fetch user: %w", err)
    }
    
    orders, err := fetchOrders(ctx, userID)
    if err != nil {
        return fmt.Errorf("fetch orders: %w", err)
    }
    
    fmt.Printf("User: %s, Orders: %v\n", user, orders)
    return nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()
    
    err := processRequest(ctx, 1)
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

**Баллы:** 15

### Задание 2: Timeout Policy

Создайте сервис с разными таймаутами для разных операций.

**Начальный код:**
```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"
)

// TODO: Создайте функцию согласно заданию
// TODO: Используйте select для работы с каналами
// TODO: Используйте defer

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Read: Data-1, err=<nil>
Data written: test
Write: err=<nil>
Fast Read: err=context deadline exceeded
```

**Баллы:** 15

### Задание 3: Middleware Chain

Создайте цепочку middleware с правильной передачей context.

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

**Баллы:** 10

### Задание 4: Context в тестах

Напишите тесты с использованием context.

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

**Баллы:** 10

### Задание 5: Structured Logging

Реализуйте структурированное логирование с context.

**Начальный код:**
```go
package main

import (
    "context"
    "fmt"
)

// TODO: Определите структуру
// TODO: Используйте context

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод (JSON):**
```json
{"level":"info","message":"Processing started","request_id":"abc-123","timestamp":"...","user_id":42}
{"level":"info","message":"Order found","order_id":999,"request_id":"abc-123","timestamp":"...","user_id":42}
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Go Blog — Context](https://go.dev/blog/context)
- [Context and Structs](https://go.dev/blog/context-and-structs)
- [Context Package](https://pkg.go.dev/context)
