# Каналы и конструкция select

## 💡 Ключевые идеи

1. **select** — мультиплексор для нескольких операций с каналами
2. **Неблокирующий выбор** — выполняется первый готовый case
3. **default** — выполняется если все каналы заблокированы
4. **Случайный выбор** — при нескольких готовых case выбирается случайный
5. **Таймауты** — реализуются через `time.After`
6. **Отмена** — через закрытие канала в case

---

## 📋 Синтаксис

### Базовый select

```go
select {
case val := <-ch1:
    // получили из ch1
case ch2 <- value:
    // отправили в ch2
case val, ok := <-ch3:
    // получили с проверкой закрытия
default:
    // если все case заблокированы
}
```

### Бесконечный select

```go
for {
    select {
    case <-done:
        return
    case val := <-ch:
        // обработка
    }
}
```

---

## 💻 Примеры кода

### Базовый select

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(100 * time.Millisecond)
        ch1 <- "from channel 1"
    }()
    
    go func() {
        time.Sleep(200 * time.Millisecond)
        ch2 <- "from channel 2"
    }()
    
    // Ждём первый готовый
    select {
    case msg := <-ch1:
        fmt.Println("Received:", msg)
    case msg := <-ch2:
        fmt.Println("Received:", msg)
    }
}
```

### Обработка всех каналов

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    ch3 := make(chan int)
    
    go func() { time.Sleep(50 * time.Millisecond); ch1 <- 1 }()
    go func() { time.Sleep(100 * time.Millisecond); ch2 <- 2 }()
    go func() { time.Sleep(150 * time.Millisecond); ch3 <- 3 }()
    
    // Получаем все три значения
    for i := 0; i < 3; i++ {
        select {
        case v := <-ch1:
            fmt.Println("ch1:", v)
        case v := <-ch2:
            fmt.Println("ch2:", v)
        case v := <-ch3:
            fmt.Println("ch3:", v)
        }
    }
}
```

### Неблокирующий select с default

```go
package main

import "fmt"

func main() {
    ch := make(chan int)
    
    // Неблокирующее чтение
    select {
    case val := <-ch:
        fmt.Println("Received:", val)
    default:
        fmt.Println("No data available")
    }
    
    // Неблокирующая отправка
    select {
    case ch <- 42:
        fmt.Println("Sent successfully")
    default:
        fmt.Println("Cannot send (no receiver)")
    }
}
```

### Таймаут с time.After

```go
package main

import (
    "fmt"
    "time"
)

func slowOperation() <-chan string {
    ch := make(chan string)
    go func() {
        time.Sleep(3 * time.Second)
        ch <- "result"
    }()
    return ch
}

func main() {
    result := slowOperation()
    
    select {
    case data := <-result:
        fmt.Println("Got result:", data)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout! Operation took too long")
    }
}
```

### Отмена операции

```go
package main

import (
    "fmt"
    "time"
)

func worker(done <-chan struct{}) <-chan int {
    result := make(chan int)
    
    go func() {
        defer close(result)
        
        for i := 1; ; i++ {
            select {
            case <-done:
                fmt.Println("Worker: cancelled")
                return
            case result <- i:
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()
    
    return result
}

func main() {
    done := make(chan struct{})
    results := worker(done)
    
    // Получаем 3 значения
    for i := 0; i < 3; i++ {
        fmt.Println("Received:", <-results)
    }
    
    // Отменяем
    close(done)
    
    time.Sleep(time.Second)
}
```

### Периодические задачи с select

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ticker := time.NewTicker(500 * time.Millisecond)
    done := make(chan struct{})
    
    go func() {
        time.Sleep(2 * time.Second)
        close(done)
    }()
    
    for {
        select {
        case <-done:
            ticker.Stop()
            fmt.Println("Done!")
            return
        case t := <-ticker.C:
            fmt.Println("Tick at", t.Format("15:04:05.000"))
        }
    }
}
```

### Приоритетный select (псевдо-приоритет)

```go
package main

import "fmt"

func main() {
    high := make(chan string, 10)
    low := make(chan string, 10)
    
    // Заполняем каналы
    for i := 0; i < 5; i++ {
        high <- fmt.Sprintf("high-%d", i)
        low <- fmt.Sprintf("low-%d", i)
    }
    
    // Приоритетное чтение
    for i := 0; i < 10; i++ {
        select {
        case msg := <-high:
            fmt.Println("Priority:", msg)
        default:
            select {
            case msg := <-high:
                fmt.Println("Priority:", msg)
            case msg := <-low:
                fmt.Println("Normal:", msg)
            }
        }
    }
}
```

### Fan-In с select

```go
package main

import (
    "fmt"
    "time"
)

func producer(name string, delay time.Duration) <-chan string {
    ch := make(chan string)
    go func() {
        for i := 1; i <= 3; i++ {
            time.Sleep(delay)
            ch <- fmt.Sprintf("%s-%d", name, i)
        }
        close(ch)
    }()
    return ch
}

func fanIn(ch1, ch2 <-chan string) <-chan string {
    out := make(chan string)
    
    go func() {
        var v1, v2 string
        var ok1, ok2 bool = true, true
        
        for ok1 || ok2 {
            select {
            case v1, ok1 = <-ch1:
                if ok1 {
                    out <- v1
                }
            case v2, ok2 = <-ch2:
                if ok2 {
                    out <- v2
                }
            }
        }
        close(out)
    }()
    
    return out
}

func main() {
    ch1 := producer("A", 200*time.Millisecond)
    ch2 := producer("B", 300*time.Millisecond)
    
    merged := fanIn(ch1, ch2)
    
    for msg := range merged {
        fmt.Println("Received:", msg)
    }
}
```

### Практический пример: HTTP-like timeout

```go
package main

import (
    "errors"
    "fmt"
    "time"
)

func fetchData(url string) <-chan string {
    result := make(chan string)
    go func() {
        time.Sleep(time.Duration(len(url)*50) * time.Millisecond)
        result <- fmt.Sprintf("Data from %s", url)
    }()
    return result
}

func fetchWithTimeout(url string, timeout time.Duration) (string, error) {
    select {
    case data := <-fetchData(url):
        return data, nil
    case <-time.After(timeout):
        return "", errors.New("request timeout")
    }
}

func main() {
    urls := []string{
        "short.com",
        "medium-length-url.com",
        "very-very-long-website-address.com",
    }
    
    for _, url := range urls {
        result, err := fetchWithTimeout(url, 500*time.Millisecond)
        if err != nil {
            fmt.Printf("%s: %v\n", url, err)
        } else {
            fmt.Printf("%s: %s\n", url, result)
        }
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Пустой select блокируется навсегда

```go
// ❌ DEADLOCK
select {}

// ✅ Используйте хотя бы один case или default
select {
case <-done:
    return
}
```

### 2. Утечка time.After в цикле

```go
// ❌ Каждая итерация создаёт новый таймер
for {
    select {
    case <-ch:
    case <-time.After(time.Second):  // утечка памяти!
    }
}

// ✅ Используйте time.NewTimer
timer := time.NewTimer(time.Second)
for {
    select {
    case <-ch:
        if !timer.Stop() {
            <-timer.C
        }
        timer.Reset(time.Second)
    case <-timer.C:
        timer.Reset(time.Second)
    }
}
```

### 3. Неверное понимание случайности

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)

ch1 <- 1
ch2 <- 2

// Оба case готовы — выбор СЛУЧАЙНЫЙ, не последовательный!
select {
case v := <-ch1:
    fmt.Println(v)  // может быть 1 или 2
case v := <-ch2:
    fmt.Println(v)
}
```

### 4. Забыли break в for-select

```go
// ❌ break выходит только из select, не из for!
for {
    select {
    case <-done:
        break  // НЕ выходит из for!
    }
}

// ✅ Используйте return или label
outer:
for {
    select {
    case <-done:
        break outer  // выходит из for
    }
}
```

---

## 📝 Практика

### Задача 1: First response
Отправьте запрос на 3 "сервера" и верните первый ответ.

### Задача 2: Rate limiter
Реализуйте rate limiter: максимум N операций в секунду.

### Задача 3: Timeout retry
Реализуйте операцию с таймаутом и повторными попытками.

### Задача 4: Heartbeat monitor
Следите за heartbeat сигналами, сообщайте если нет сигнала N секунд.

### Задача 5: Multiplexer
Объедините 3 канала в один с приоритетом.

### Задача 6: Debounce
Реализуйте debounce: выполнять действие только после паузы в событиях.

### Задача 7: Graceful shutdown
Реализуйте завершение с таймаутом: сначала мягкое, потом принудительное.

### Задача 8: Circuit breaker
Реализуйте circuit breaker с состояниями: closed, open, half-open.
