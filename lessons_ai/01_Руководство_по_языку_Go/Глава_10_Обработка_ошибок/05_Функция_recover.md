# Восстановление после ошибки и функция recover

## 💡 Ключевые идеи

1. **recover()** — перехватывает panic и возвращает переданное значение
2. **Только в defer** — recover работает ТОЛЬКО внутри deferred функции
3. **Возвращает nil** — если panic не было или вызван не в defer
4. **Продолжение выполнения** — после recover программа продолжает работу
5. **Аналог try-catch** — механизм похож на исключения, но идиоматически отличается
6. **Применение** — серверы, обработчики, изоляция критических участков

---

## 📋 Синтаксис

### Базовый recover

```go
defer func() {
    if r := recover(); r != nil {
        // Обработка panic
        fmt.Println("Recovered:", r)
    }
}()
```

### recover вне defer не работает

```go
if r := recover(); r != nil {  // Всегда nil — не в defer!
    fmt.Println("This won't work")
}
```

---

## 💻 Примеры кода

### Базовый recover

```go
package main

import "fmt"

func mayPanic() {
    panic("something went wrong!")
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from:", r)
        }
    }()
    
    fmt.Println("Before panic")
    mayPanic()
    fmt.Println("After panic")  // Не выполнится
}
// Output:
// Before panic
// Recovered from: something went wrong!
```

### Программа продолжает после recover

```go
package main

import "fmt"

func riskyOperation() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Caught panic:", r)
        }
    }()
    
    panic("boom!")
}

func main() {
    fmt.Println("Start")
    
    riskyOperation()  // Panic перехвачен внутри
    
    fmt.Println("Continue after risky operation")
    fmt.Println("End")
}
// Output:
// Start
// Caught panic: boom!
// Continue after risky operation
// End
```

### Паттерн try-catch

```go
package main

import "fmt"

func tryCatch(try func(), catch func(interface{})) {
    defer func() {
        if r := recover(); r != nil {
            catch(r)
        }
    }()
    
    try()
}

func main() {
    tryCatch(
        func() {
            fmt.Println("Doing risky stuff...")
            panic("error occurred!")
        },
        func(err interface{}) {
            fmt.Println("Caught:", err)
        },
    )
    
    fmt.Println("Program continues")
}
```

### Конвертация panic в error

```go
package main

import (
    "fmt"
)

func safeCall(f func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    
    f()
    return nil
}

func riskyFunction() {
    panic("something bad happened")
}

func main() {
    err := safeCall(riskyFunction)
    if err != nil {
        fmt.Println("Error:", err)
    }
    
    err = safeCall(func() {
        fmt.Println("Safe operation")
    })
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Success!")
    }
}
// Output:
// Error: panic: something bad happened
// Safe operation
// Success!
```

### Recover с именованным возвращаемым значением

```go
package main

import "fmt"

func divide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered: %v", r)
            result = 0
        }
    }()
    
    if b == 0 {
        panic("division by zero")
    }
    
    return a / b, nil
}

func main() {
    result, err := divide(10, 2)
    fmt.Printf("10 / 2 = %d, err = %v\n", result, err)
    
    result, err = divide(10, 0)
    fmt.Printf("10 / 0 = %d, err = %v\n", result, err)
}
// Output:
// 10 / 2 = 5, err = <nil>
// 10 / 0 = 0, err = recovered: division by zero
```

### Recover в HTTP handler

```go
package main

import (
    "fmt"
    "net/http"
)

func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                fmt.Printf("Panic in handler: %v\n", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
    panic("handler crashed!")
}

func okHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "OK")
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/panic", panicHandler)
    mux.HandleFunc("/ok", okHandler)
    
    handler := recoveryMiddleware(mux)
    
    // Для демонстрации
    fmt.Println("Server would start on :8080")
    // http.ListenAndServe(":8080", handler)
    _ = handler
}
```

### Recover в горутине

```go
package main

import (
    "fmt"
    "sync"
)

func safeGo(f func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Goroutine panicked:", r)
            }
        }()
        
        f()
    }()
}

func main() {
    var wg sync.WaitGroup
    
    wg.Add(3)
    
    safeGo(func() {
        defer wg.Done()
        fmt.Println("Goroutine 1: OK")
    })
    
    safeGo(func() {
        defer wg.Done()
        panic("Goroutine 2 crashed!")
    })
    
    safeGo(func() {
        defer wg.Done()
        fmt.Println("Goroutine 3: OK")
    })
    
    wg.Wait()
    fmt.Println("All goroutines finished")
}
```

### Типизированный recover

```go
package main

import "fmt"

type ValidationError struct {
    Field   string
    Message string
}

type SystemError struct {
    Code    int
    Message string
}

func process(input int) {
    if input < 0 {
        panic(ValidationError{Field: "input", Message: "must be positive"})
    }
    if input > 100 {
        panic(SystemError{Code: 500, Message: "overflow"})
    }
    fmt.Println("Processing:", input)
}

func safeProcess(input int) error {
    defer func() {
        if r := recover(); r != nil {
            switch e := r.(type) {
            case ValidationError:
                fmt.Printf("Validation error on %s: %s\n", e.Field, e.Message)
            case SystemError:
                fmt.Printf("System error [%d]: %s\n", e.Code, e.Message)
            default:
                fmt.Printf("Unknown panic: %v\n", r)
            }
        }
    }()
    
    process(input)
    return nil
}

func main() {
    safeProcess(50)   // OK
    safeProcess(-1)   // ValidationError
    safeProcess(150)  // SystemError
}
```

### Практический пример: Worker Pool с recover

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID   int
    Data interface{}
}

type Result struct {
    JobID int
    Data  interface{}
    Error error
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        result := processJobSafe(job)
        results <- result
    }
}

func processJobSafe(job Job) (result Result) {
    result.JobID = job.ID
    
    defer func() {
        if r := recover(); r != nil {
            result.Error = fmt.Errorf("job %d panicked: %v", job.ID, r)
        }
    }()
    
    // Имитация обработки
    if job.ID == 3 {
        panic("job 3 always fails!")
    }
    
    result.Data = fmt.Sprintf("processed-%d", job.ID)
    return result
}

func main() {
    jobs := make(chan Job, 5)
    results := make(chan Result, 5)
    
    var wg sync.WaitGroup
    
    // Запускаем 2 воркера
    for w := 1; w <= 2; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }
    
    // Отправляем задачи
    for j := 1; j <= 5; j++ {
        jobs <- Job{ID: j, Data: fmt.Sprintf("data-%d", j)}
    }
    close(jobs)
    
    // Закрываем results когда воркеры завершены
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Собираем результаты
    for r := range results {
        if r.Error != nil {
            fmt.Printf("Job %d FAILED: %v\n", r.JobID, r.Error)
        } else {
            fmt.Printf("Job %d OK: %v\n", r.JobID, r.Data)
        }
    }
}
```

---

## ⚠️ Частые ошибки

### 1. recover вне defer

```go
// ❌ НЕ РАБОТАЕТ
func wrong() {
    r := recover()  // Всегда nil!
    if r != nil {
        fmt.Println(r)
    }
    panic("error")
}

// ✅ ПРАВИЛЬНО
func right() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println(r)
        }
    }()
    panic("error")
}
```

### 2. recover в другой горутине

```go
// ❌ НЕ РАБОТАЕТ — recover в main, panic в горутине
func wrong() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Caught:", r)  // Никогда не выполнится!
        }
    }()
    
    go func() {
        panic("panic in goroutine")
    }()
    
    time.Sleep(time.Second)
}

// ✅ ПРАВИЛЬНО — recover в той же горутине
func right() {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Caught:", r)
            }
        }()
        panic("panic in goroutine")
    }()
    
    time.Sleep(time.Second)
}
```

### 3. Злоупотребление recover

```go
// ❌ ПЛОХО — скрывает все ошибки
func badPractice() {
    defer func() { recover() }()  // Съедает все panic без обработки!
    // код...
}

// ✅ ХОРОШО — логируем и обрабатываем
func goodPractice() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
            // Возможно, re-panic или преобразование в error
        }
    }()
    // код...
}
```

### 4. re-panic без контекста

```go
// ❌ Теряем информацию
defer func() {
    if r := recover(); r != nil {
        panic(r)  // Теряем stack trace
    }
}()

// ✅ Добавляем контекст
defer func() {
    if r := recover(); r != nil {
        panic(fmt.Sprintf("in myFunction: %v", r))
    }
}()
```

---

## 📝 Практика

### Задача 1: Safe wrapper
Создайте функцию `Safe(f func()) error`, конвертирующую panic в error.

### Задача 2: Retry с recover
Реализуйте retry, который повторяет вызов при panic.

### Задача 3: HTTP middleware
Создайте middleware, перехватывающий panic и возвращающий 500.

### Задача 4: Graceful degradation
Реализуйте функцию, которая при panic возвращает значение по умолчанию.

### Задача 5: Типизированный recover
Обработайте разные типы panic по-разному.

### Задача 6: Worker pool
Создайте пул воркеров с recover для каждой задачи.

### Задача 7: Цепочка recover
Продемонстрируйте, как panic проходит через несколько уровней defer.

### Задача 8: Logging panic
Создайте функцию для логирования panic с полным stack trace.
