# Отмена и Timeout

---

## 💡 Ключевые идеи

1. **WithCancel** — ручная отмена контекста
2. **WithTimeout** — автоматическая отмена через заданное время
3. **WithDeadline** — отмена в заданный момент времени
4. **ctx.Done()** — канал, закрывающийся при отмене
5. **ctx.Err()** — причина отмены (Canceled или DeadlineExceeded)

### Сравнение методов

| Метод | Когда использовать |
|-------|-------------------|
| `WithCancel` | Ручное управление отменой |
| `WithTimeout` | Ограничение длительности операции |
| `WithDeadline` | Абсолютное время завершения |

---

## 📖 Теория

### Три механизма ограничения выполнения

Go предоставляет три способа ограничить время выполнения операций:

**1. WithCancel — ручная отмена**
```go
ctx, cancel := context.WithCancel(parent)
// ... где-то позже
cancel()  // отменяем вручную
```

Используйте, когда отмена зависит от внешнего события (пользователь нажал "отмена", получен результат от другой горутины).

**2. WithTimeout — относительное время**
```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()
// контекст автоматически отменится через 5 секунд
```

Используйте, когда операция должна завершиться за определённое время.

**3. WithDeadline — абсолютное время**
```go
deadline := time.Now().Add(30 * time.Minute)
ctx, cancel := context.WithDeadline(parent, deadline)
defer cancel()
// контекст отменится в указанный момент
```

Используйте, когда есть конкретное время завершения (например, "до 18:00").

### Разница между Timeout и Deadline

```go
// Timeout: "через 5 секунд"
ctx, _ := context.WithTimeout(parent, 5*time.Second)

// Deadline: "в 15:30:00"
ctx, _ := context.WithDeadline(parent, time.Date(2024, 1, 1, 15, 30, 0, 0, time.UTC))

// Внутренне Timeout использует Deadline:
// WithTimeout(parent, 5s) = WithDeadline(parent, time.Now().Add(5s))
```

### Как проверять отмену

**Способ 1: select с ctx.Done()**
```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-workChan:
    return result
}
```

**Способ 2: в цикле**
```go
for i := 0; i < 1000000; i++ {
    select {
    case <-ctx.Done():
        return ctx.Err()  // выходим досрочно
    default:
        // продолжаем работу
    }
    processItem(i)
}
```

**Способ 3: ctx.Err() для проверки**
```go
if ctx.Err() != nil {
    return ctx.Err()
}
// контекст ещё активен
```

### Два типа ошибок отмены

```go
if errors.Is(ctx.Err(), context.Canceled) {
    // контекст был отменён вручную (cancel())
    fmt.Println("Операция отменена пользователем")
}

if errors.Is(ctx.Err(), context.DeadlineExceeded) {
    // истёк таймаут или дедлайн
    fmt.Println("Операция не успела выполниться")
}
```

### Наследование дедлайнов

Дочерний контекст не может иметь дедлайн позже родительского:

```go
parentCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)

// Попытка создать больший таймаут
childCtx, _ := context.WithTimeout(parentCtx, 10*time.Second)

// childCtx всё равно отменится через 5 секунд (вместе с parent)!
```

Это защита от "зависших" операций — дочерняя операция не может работать дольше родительской.

### Проверка оставшегося времени

```go
deadline, ok := ctx.Deadline()
if ok {
    remaining := time.Until(deadline)
    fmt.Printf("Осталось: %v\n", remaining)
}
```

Это полезно для адаптации поведения: если времени мало — пропустить некритичные операции.

### Паттерн: Timeout для внешних вызовов

```go
func callExternalAPI(ctx context.Context) (Result, error) {
    // Создаём собственный таймаут, не больше родительского
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := client.Do(req)
    // ...
}
```

Каждый внешний вызов получает свой таймаут, но не превышающий общий.

---

## 💻 Примеры кода

### Пример 1: Паттерн отмены

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func fetchData(ctx context.Context, source string) (string, error) {
    // Имитация медленного запроса
    resultChan := make(chan string, 1)
    
    go func() {
        time.Sleep(2 * time.Second)
        resultChan <- fmt.Sprintf("Data from %s", source)
    }()
    
    select {
    case result := <-resultChan:
        return result, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func main() {
    // Сценарий 1: Успешное выполнение
    ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel1()
    
    result, err := fetchData(ctx1, "API")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", result)  // Data from API
    }
    
    // Сценарий 2: Таймаут
    ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel2()
    
    result, err = fetchData(ctx2, "API")
    if err != nil {
        fmt.Println("Error:", err)  // context deadline exceeded
    }
}
```

### Пример 2: Каскадная отмена

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

func task(ctx context.Context, name string, duration time.Duration, wg *sync.WaitGroup) {
    defer wg.Done()
    
    select {
    case <-time.After(duration):
        fmt.Printf("%s: завершена успешно\n", name)
    case <-ctx.Done():
        fmt.Printf("%s: отменена (%v)\n", name, ctx.Err())
    }
}

func main() {
    // Родительский контекст
    parentCtx, parentCancel := context.WithCancel(context.Background())
    
    // Дочерний контекст с таймаутом
    childCtx, childCancel := context.WithTimeout(parentCtx, 5*time.Second)
    defer childCancel()
    
    var wg sync.WaitGroup
    
    // Запускаем задачи с дочерним контекстом
    wg.Add(3)
    go task(childCtx, "Task1", 2*time.Second, &wg)
    go task(childCtx, "Task2", 4*time.Second, &wg)
    go task(childCtx, "Task3", 6*time.Second, &wg)
    
    // Через 3 секунды отменяем родительский контекст
    time.Sleep(3 * time.Second)
    fmt.Println("Отменяем родительский контекст...")
    parentCancel()
    
    wg.Wait()
    fmt.Println("Все задачи завершены")
}
// Вывод:
// Task1: завершена успешно
// Отменяем родительский контекст...
// Task2: отменена (context canceled)
// Task3: отменена (context canceled)
```

### Пример 3: Timeout для HTTP запроса

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
)

func fetchURL(ctx context.Context, url string) ([]byte, error) {
    // Создаём запрос с контекстом
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    // Выполняем запрос
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return io.ReadAll(resp.Body)
}

func main() {
    // Быстрый запрос
    ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel1()
    
    data, err := fetchURL(ctx1, "https://httpbin.org/get")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Printf("Received %d bytes\n", len(data))
    }
    
    // Медленный запрос с коротким таймаутом
    ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel2()
    
    _, err = fetchURL(ctx2, "https://httpbin.org/delay/5")
    if err != nil {
        fmt.Println("Error:", err)  // context deadline exceeded
    }
}
```

### Пример 4: Retry с уменьшающимся таймаутом

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "time"
)

func unreliableOperation() error {
    if rand.Float32() < 0.7 {
        return errors.New("operation failed")
    }
    return nil
}

func retryWithTimeout(ctx context.Context, maxRetries int, baseTimeout time.Duration) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        // Проверяем, не отменён ли контекст
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        
        // Таймаут для текущей попытки
        attemptTimeout := baseTimeout * time.Duration(i+1)
        attemptCtx, cancel := context.WithTimeout(ctx, attemptTimeout)
        
        // Выполняем операцию
        done := make(chan error, 1)
        go func() {
            time.Sleep(100 * time.Millisecond)  // имитация работы
            done <- unreliableOperation()
        }()
        
        select {
        case err := <-done:
            cancel()
            if err == nil {
                fmt.Printf("Успех на попытке %d\n", i+1)
                return nil
            }
            lastErr = err
            fmt.Printf("Попытка %d не удалась: %v\n", i+1, err)
        case <-attemptCtx.Done():
            cancel()
            lastErr = attemptCtx.Err()
            fmt.Printf("Попытка %d: таймаут\n", i+1)
        }
        
        // Пауза перед следующей попыткой
        time.Sleep(100 * time.Millisecond)
    }
    
    return fmt.Errorf("все %d попыток не удались: %w", maxRetries, lastErr)
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    err := retryWithTimeout(ctx, 5, 500*time.Millisecond)
    if err != nil {
        fmt.Println("Финальная ошибка:", err)
    }
}
```

### Пример 5: Параллельные запросы с общим таймаутом

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Result struct {
    Source string
    Data   string
    Err    error
}

func fetchFromSource(ctx context.Context, source string) Result {
    // Случайная задержка
    delay := time.Duration(rand.Intn(3000)) * time.Millisecond
    
    select {
    case <-time.After(delay):
        return Result{
            Source: source,
            Data:   fmt.Sprintf("Data from %s (took %v)", source, delay),
        }
    case <-ctx.Done():
        return Result{
            Source: source,
            Err:    ctx.Err(),
        }
    }
}

func fetchAll(ctx context.Context, sources []string) []Result {
    results := make([]Result, len(sources))
    var wg sync.WaitGroup
    
    for i, source := range sources {
        wg.Add(1)
        go func(idx int, src string) {
            defer wg.Done()
            results[idx] = fetchFromSource(ctx, src)
        }(i, source)
    }
    
    wg.Wait()
    return results
}

func main() {
    sources := []string{"API1", "API2", "API3", "API4", "API5"}
    
    // Общий таймаут 2 секунды на все запросы
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    results := fetchAll(ctx, sources)
    
    for _, r := range results {
        if r.Err != nil {
            fmt.Printf("%s: ОШИБКА - %v\n", r.Source, r.Err)
        } else {
            fmt.Printf("%s: %s\n", r.Source, r.Data)
        }
    }
}
```

### Пример 6: First Response (первый ответ)

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "time"
)

func queryReplica(ctx context.Context, replicaID int, results chan<- string) {
    delay := time.Duration(rand.Intn(2000)) * time.Millisecond
    
    select {
    case <-time.After(delay):
        select {
        case results <- fmt.Sprintf("Replica %d responded in %v", replicaID, delay):
        case <-ctx.Done():
        }
    case <-ctx.Done():
    }
}

func queryFirst(ctx context.Context, numReplicas int) (string, error) {
    // Контекст для отмены остальных после первого ответа
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    results := make(chan string, numReplicas)
    
    // Запускаем запросы ко всем репликам
    for i := 1; i <= numReplicas; i++ {
        go queryReplica(ctx, i, results)
    }
    
    // Возвращаем первый результат
    select {
    case result := <-results:
        return result, nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    result, err := queryFirst(ctx, 5)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("First result:", result)
    }
}
```

### Пример 7: Проверка оставшегося времени

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func processWithDeadline(ctx context.Context) error {
    deadline, ok := ctx.Deadline()
    if !ok {
        fmt.Println("Нет дедлайна")
        return nil
    }
    
    fmt.Printf("Дедлайн: %v\n", deadline.Format(time.RFC3339))
    
    for {
        remaining := time.Until(deadline)
        
        if remaining <= 0 {
            return ctx.Err()
        }
        
        fmt.Printf("Осталось: %v\n", remaining.Round(time.Millisecond))
        
        // Проверяем, есть ли время на ещё одну итерацию
        if remaining < 500*time.Millisecond {
            fmt.Println("Недостаточно времени для следующей итерации")
            return nil
        }
        
        select {
        case <-time.After(500 * time.Millisecond):
            fmt.Println("Выполнена итерация")
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    err := processWithDeadline(ctx)
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Пример 8: Graceful degradation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type Response struct {
    Data       string
    FromCache  bool
    Degraded   bool
}

func getFromPrimary(ctx context.Context) (string, error) {
    select {
    case <-time.After(2 * time.Second):
        return "Fresh data from primary", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}

func getFromCache() string {
    return "Cached data"
}

func getData(ctx context.Context) Response {
    // Пробуем получить свежие данные с таймаутом
    primaryCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
    defer cancel()
    
    data, err := getFromPrimary(primaryCtx)
    if err == nil {
        return Response{Data: data, FromCache: false}
    }
    
    // Fallback на кэш
    fmt.Println("Primary недоступен, используем кэш")
    return Response{
        Data:      getFromCache(),
        FromCache: true,
        Degraded:  true,
    }
}

func main() {
    ctx := context.Background()
    
    response := getData(ctx)
    
    fmt.Printf("Data: %s\n", response.Data)
    fmt.Printf("From cache: %v\n", response.FromCache)
    fmt.Printf("Degraded: %v\n", response.Degraded)
}
```

---

## ⚠️ Частые ошибки

### 1. Слишком короткий таймаут

```go
// ❌ ПЛОХО — не учитываем сетевые задержки
ctx, _ := context.WithTimeout(ctx, 100*time.Millisecond)

// ✅ ХОРОШО — реалистичный таймаут
ctx, _ := context.WithTimeout(ctx, 5*time.Second)
```

### 2. Не проверяем ctx.Err()

```go
// ❌ ПЛОХО — непонятно, что случилось
select {
case <-ctx.Done():
    return errors.New("failed")
}

// ✅ ХОРОШО — возвращаем реальную причину
select {
case <-ctx.Done():
    return ctx.Err()  // Canceled или DeadlineExceeded
}
```

### 3. Блокировка без проверки контекста

```go
// ❌ ПЛОХО — может зависнуть навсегда
result := <-channel

// ✅ ХОРОШО — проверяем контекст
select {
case result := <-channel:
    // обрабатываем
case <-ctx.Done():
    return ctx.Err()
}
```

---

## 🏋️ Практические задания

### Задание 1: Timeout Wrapper

Создайте универсальную обёртку для выполнения функций с таймаутом.

**Начальный код:**
```go
package main

import (
    "fmt"
    "time"
)

// TODO: Создайте функцию согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Result: result, Error: <nil>
Result: , Error: operation timed out
```

**Баллы:** 15

### Задание 2: Retry с отменой

Реализуйте функцию retry с поддержкой отмены.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Создайте функцию согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Attempt 1 failed: temporary error
Attempt 2 failed: temporary error
Result: <nil>
---
Attempt 1 failed: always fails
Attempt 2 failed: always fails
Result: context deadline exceeded
```

**Баллы:** 15

### Задание 3: Parallel с лимитом

Выполните задачи параллельно с лимитом горутин и общим таймаутом.

**Начальный код:**
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// TODO: Запустите горутину

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 15

### Задание 4: Deadline check

Проверяйте оставшееся время перед выполнением операции.

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
Remaining time: 499ms
Work completed
Result 1: <nil>
---
Result 2: insufficient time for operation
```

**Баллы:** 10

### Задание 5: First success

Выполните несколько запросов параллельно, вернув первый успешный.

**Начальный код:**
```go
package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Result: server2, Error: <nil>, Time: 100ms
```

**Баллы:** 15

---

## 🔗 Полезные ссылки

- [Context Package](https://pkg.go.dev/context)
- [Timeouts in Go](https://blog.golang.org/context)
