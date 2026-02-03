# Введение в Context

---

## 💡 Ключевые идеи

1. **Context** — механизм передачи сигналов отмены, дедлайнов и данных между горутинами
2. **Отмена (Cancellation)** — прерывание операций при закрытии контекста
3. **Timeout/Deadline** — ограничение времени выполнения операций
4. **Values** — передача request-scoped данных (осторожно!)
5. **Иерархия** — дочерние контексты наследуют свойства родительских

### Зачем нужен Context?

| Проблема | Решение с Context |
|----------|-------------------|
| Утечка горутин | Сигнал отмены |
| Зависшие запросы | Timeout |
| Передача request ID | Values |
| Graceful shutdown | Cancellation |

---

## 📖 Теория

### Что такое Context?

**Context** — это механизм для передачи "контекста выполнения" между горутинами. Контекст содержит:
- Сигнал отмены — "прекрати работу"
- Дедлайн — "заверши работу до этого времени"
- Значения — данные, связанные с запросом

### Зачем нужен Context?

Представьте HTTP-сервер, который при запросе:
1. Обращается к базе данных
2. Вызывает внешний API
3. Обрабатывает данные

Что если клиент закрыл соединение? Без Context сервер продолжит выполнять все операции впустую — это трата ресурсов. С Context все операции получают сигнал "прекрати работу" и освобождают ресурсы.

### Аналогия из жизни

Context — как радиосвязь в команде:
- **Главный** создаёт контекст и может сказать "операция отменена"
- **Все участники** слушают канал и реагируют на команды
- **Дедлайн** — "операция должна закончиться к 15:00"
- **Иерархия** — команды передаются вниз по цепочке

### Иерархия контекстов

Контексты образуют **дерево**. Когда родительский контекст отменяется, все дочерние тоже отменяются:

```
context.Background() (корень)
    └── WithTimeout (API запрос, 30s)
            ├── WithValue (request ID)
            │       └── WithTimeout (DB query, 5s)
            └── WithCancel (фоновая задача)
```

Если отменить "API запрос" — отменятся и "DB query", и "фоновая задача".

### Четыре способа создания Context

**1. context.Background()** — корневой контекст, используется в main() или как основа для других контекстов.

**2. context.TODO()** — временный контекст, когда вы ещё не знаете, какой контекст использовать. Это маркер "доделать позже".

**3. context.WithCancel(parent)** — создаёт контекст с функцией отмены. Вызов cancel() отменяет контекст.

**4. context.WithTimeout/WithDeadline(parent, time)** — автоматически отменяется через заданное время.

### Как работает отмена?

Когда контекст отменяется:
1. Закрывается канал `ctx.Done()`
2. `ctx.Err()` возвращает причину (Canceled или DeadlineExceeded)
3. Все дочерние контексты тоже отменяются

```go
select {
case <-ctx.Done():
    // контекст отменён
    return ctx.Err()
case result := <-resultChan:
    // получили результат
    return result
}
```

### Почему важно вызывать cancel()?

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()  // ОБЯЗАТЕЛЬНО!
```

Если не вызвать `cancel()`, ресурсы (таймеры, горутины) не освободятся до истечения таймаута. Это утечка ресурсов!

**Правило:** Всегда вызывайте `defer cancel()` сразу после создания контекста.

### Context в стандартной библиотеке

Context интегрирован во многие пакеты Go:
- `net/http` — все запросы имеют контекст
- `database/sql` — QueryContext, ExecContext
- `os/exec` — CommandContext
- `net` — DialContext

```go
// HTTP клиент с таймаутом
req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := client.Do(req)

// SQL запрос с контекстом
rows, err := db.QueryContext(ctx, "SELECT * FROM users")
```

---

## 📋 Синтаксис

### Создание контекста

```go
// Пустой контекст (корневой)
ctx := context.Background()

// Контекст для тестов
ctx := context.TODO()

// С отменой
ctx, cancel := context.WithCancel(parent)

// С таймаутом
ctx, cancel := context.WithTimeout(parent, 5*time.Second)

// С дедлайном
ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))

// Со значением
ctx := context.WithValue(parent, key, value)
```

### Проверка отмены

```go
select {
case <-ctx.Done():
    return ctx.Err()
default:
    // продолжаем работу
}
```

### Получение значения

```go
value := ctx.Value(key)
if v, ok := value.(MyType); ok {
    // используем v
}
```

---

## 💻 Примеры кода

### Пример 1: Базовый Context

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func doWork(ctx context.Context, name string) {
    for i := 1; i <= 5; i++ {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: отменено на итерации %d\n", name, i)
            return
        default:
            fmt.Printf("%s: итерация %d\n", name, i)
            time.Sleep(500 * time.Millisecond)
        }
    }
    fmt.Printf("%s: завершено\n", name)
}

func main() {
    // Создаём контекст с отменой
    ctx, cancel := context.WithCancel(context.Background())
    
    go doWork(ctx, "Worker")
    
    // Ждём 2 секунды и отменяем
    time.Sleep(2 * time.Second)
    cancel()
    
    // Даём время на завершение
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Main: завершено")
}
```

### Пример 2: Timeout

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func slowOperation(ctx context.Context) error {
    select {
    case <-time.After(5 * time.Second):
        fmt.Println("Операция завершена")
        return nil
    case <-ctx.Done():
        return ctx.Err()  // context deadline exceeded
    }
}

func main() {
    // Контекст с таймаутом 2 секунды
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()  // важно! освобождаем ресурсы
    
    err := slowOperation(ctx)
    if err != nil {
        fmt.Println("Ошибка:", err)  // context deadline exceeded
    }
}
```

### Пример 3: Deadline

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func processTask(ctx context.Context, taskID int) error {
    // Проверяем, сколько времени осталось
    deadline, ok := ctx.Deadline()
    if ok {
        fmt.Printf("Task %d: дедлайн через %v\n", taskID, time.Until(deadline))
    }
    
    // Имитация работы
    workDuration := time.Duration(taskID) * 500 * time.Millisecond
    
    select {
    case <-time.After(workDuration):
        fmt.Printf("Task %d: выполнена за %v\n", taskID, workDuration)
        return nil
    case <-ctx.Done():
        return fmt.Errorf("task %d: %w", taskID, ctx.Err())
    }
}

func main() {
    // Дедлайн через 2 секунды
    deadline := time.Now().Add(2 * time.Second)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()
    
    // Запускаем несколько задач
    for i := 1; i <= 5; i++ {
        err := processTask(ctx, i)
        if err != nil {
            fmt.Println("Ошибка:", err)
            break
        }
    }
}
```

### Пример 4: Передача значений

```go
package main

import (
    "context"
    "fmt"
)

// Тип для ключей (избегаем коллизий)
type contextKey string

const (
    userIDKey    contextKey = "userID"
    requestIDKey contextKey = "requestID"
)

func processRequest(ctx context.Context) {
    // Получаем значения из контекста
    userID, ok := ctx.Value(userIDKey).(int)
    if !ok {
        fmt.Println("userID не найден")
        return
    }
    
    requestID, ok := ctx.Value(requestIDKey).(string)
    if !ok {
        fmt.Println("requestID не найден")
        return
    }
    
    fmt.Printf("Обработка запроса %s для пользователя %d\n", requestID, userID)
}

func main() {
    // Создаём контекст со значениями
    ctx := context.Background()
    ctx = context.WithValue(ctx, userIDKey, 42)
    ctx = context.WithValue(ctx, requestIDKey, "req-123-456")
    
    processRequest(ctx)
}
```

### Пример 5: HTTP сервер с Context

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

func slowHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    fmt.Println("Handler: начало обработки")
    
    select {
    case <-time.After(10 * time.Second):
        fmt.Fprintf(w, "Операция завершена")
    case <-ctx.Done():
        // Клиент отменил запрос или таймаут
        fmt.Println("Handler: запрос отменён")
        http.Error(w, ctx.Err().Error(), http.StatusRequestTimeout)
    }
    
    fmt.Println("Handler: завершение")
}

func main() {
    http.HandleFunc("/slow", slowHandler)
    
    server := &http.Server{
        Addr:         ":8080",
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    
    fmt.Println("Сервер запущен на :8080")
    server.ListenAndServe()
}
```

### Пример 6: Graceful Shutdown

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    server := &http.Server{Addr: ":8080"}
    
    // Канал для сигналов ОС
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    // Запускаем сервер в горутине
    go func() {
        fmt.Println("Сервер запущен на :8080")
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            fmt.Printf("Ошибка сервера: %v\n", err)
        }
    }()
    
    // Ждём сигнал завершения
    <-quit
    fmt.Println("\nПолучен сигнал завершения...")
    
    // Контекст с таймаутом для graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Graceful shutdown
    if err := server.Shutdown(ctx); err != nil {
        fmt.Printf("Ошибка при завершении: %v\n", err)
    }
    
    fmt.Println("Сервер остановлен")
}
```

### Пример 7: Context в базе данных

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/lib/pq"
)

func queryWithTimeout(db *sql.DB) error {
    // Контекст с таймаутом 3 секунды
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    // Запрос с контекстом
    rows, err := db.QueryContext(ctx, "SELECT pg_sleep(5)")  // долгий запрос
    if err != nil {
        return fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()
    
    // Обработка результатов
    for rows.Next() {
        // ...
    }
    
    return rows.Err()
}

func main() {
    db, err := sql.Open("postgres", "postgresql://localhost/testdb")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    
    err = queryWithTimeout(db)
    if err != nil {
        fmt.Println("Ошибка:", err)  // context deadline exceeded
    }
}
```

### Пример 8: Множественные горутины

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: остановлен (%v)\n", id, ctx.Err())
            return
        default:
            fmt.Printf("Worker %d: работаю\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    var wg sync.WaitGroup
    
    // Запускаем несколько workers
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(ctx, i, &wg)
    }
    
    // Ждём завершения всех workers
    wg.Wait()
    fmt.Println("Все workers остановлены")
}
```

---

## ⚠️ Частые ошибки

### 1. Забыли вызвать cancel

```go
// ❌ ПЛОХО — утечка ресурсов
ctx, cancel := context.WithTimeout(ctx, time.Second)
// забыли defer cancel()

// ✅ ХОРОШО
ctx, cancel := context.WithTimeout(ctx, time.Second)
defer cancel()  // всегда вызываем!
```

### 2. Хранение Context в структуре

```go
// ❌ ПЛОХО — контекст не должен храниться
type Service struct {
    ctx context.Context  // НЕ ДЕЛАЙТЕ ТАК
}

// ✅ ХОРОШО — передаём как первый параметр
func (s *Service) DoWork(ctx context.Context) error {
    // ...
}
```

### 3. Неправильное использование Values

```go
// ❌ ПЛОХО — передача бизнес-данных
ctx = context.WithValue(ctx, "user", user)  // строковый ключ!

// ✅ ХОРОШО — типизированный ключ, только request-scoped данные
type ctxKey int
const userIDKey ctxKey = 0
ctx = context.WithValue(ctx, userIDKey, userID)
```

### 4. Игнорирование ctx.Done()

```go
// ❌ ПЛОХО — не проверяем отмену
func process(ctx context.Context) {
    for i := 0; i < 1000000; i++ {
        doHeavyWork()  // не реагирует на отмену
    }
}

// ✅ ХОРОШО — регулярно проверяем
func process(ctx context.Context) error {
    for i := 0; i < 1000000; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            doHeavyWork()
        }
    }
    return nil
}
```

---

## 🏋️ Практические задания

### Задание 1: Первый Context

Напишите функцию с использованием context для отмены.

**Начальный код:**
```go
package main

import (
    "context"
    "fmt"
)

// TODO: Создайте функцию согласно заданию
// TODO: Используйте context

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Worker 1: iteration 1
Worker 1: iteration 2
Worker 1: cancelled
Error: context canceled
```

**Баллы:** 10

### Задание 2: HTTP запрос с таймаутом

Создайте HTTP клиент с таймаутом через context.

**Начальный код:**
```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

// TODO: Используйте context
// TODO: Реализуйте HTTP обработчик

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 15

### Задание 3: Worker Pool с отменой

Создайте пул воркеров с общей отменой.

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
Worker 0: processing job 1
Worker 1: processing job 2
Worker 2: processing job 3
Worker 0: processing job 4
Worker 1: processing job 5
Cancelling...
Worker 0: shutting down
Worker 1: shutting down
Worker 2: shutting down
All workers stopped
```

**Баллы:** 15

### Задание 4: Цепочка вызовов

Пробросьте context через несколько уровней функций.

**Начальный код:**
```go
package main

import (
    "context"
    "fmt"
)

// TODO: Создайте функцию согласно заданию
// TODO: Используйте context

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Service A: start
Service B: start
Service C: start
Service C: done
Service B: done
Service A: done
Success!
---
Service A: start
Service B: start
Service C: start
Service C: cancelled
Error: serviceA: serviceB: context deadline exceeded
```

**Баллы:** 10

### Задание 5: Graceful Shutdown

Реализуйте корректное завершение HTTP сервера.

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

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Go Context Package](https://pkg.go.dev/context)
- [Go Blog — Context](https://go.dev/blog/context)
- [Context Best Practices](https://go.dev/blog/context-and-structs)
