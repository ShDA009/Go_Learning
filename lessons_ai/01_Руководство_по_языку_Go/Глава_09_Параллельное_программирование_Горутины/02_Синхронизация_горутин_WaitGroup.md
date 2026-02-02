# Синхронизация горутин и WaitGroup

## 💡 Ключевые идеи

1. **sync.WaitGroup** — механизм ожидания завершения группы горутин
2. **Внутренний счётчик** — отслеживает количество активных горутин
3. **Add(n)** — увеличивает счётчик на n (регистрирует горутины)
4. **Done()** — уменьшает счётчик на 1 (сигнал о завершении)
5. **Wait()** — блокирует выполнение пока счётчик не станет 0
6. **Передача по указателю** — WaitGroup должен передаваться как `*sync.WaitGroup`

---

## 📋 Синтаксис

### Объявление и использование

```go
var wg sync.WaitGroup

wg.Add(n)    // добавить n горутин в группу
wg.Done()    // уменьшить счётчик на 1 (обычно defer wg.Done())
wg.Wait()    // ждать пока счётчик не станет 0
```

### Передача в функцию

```go
func worker(wg *sync.WaitGroup) {  // указатель!
    defer wg.Done()
    // работа
}
```

---

## 💻 Примеры кода

### Базовый пример WaitGroup

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done()  // гарантированно вызовется при выходе из функции
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1)          // регистрируем горутину
        go worker(i, &wg)  // передаём указатель
    }
    
    wg.Wait()  // ждём завершения всех горутин
    fmt.Println("All workers completed!")
}
```

### WaitGroup с анонимными функциями

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup
    
    tasks := []string{"task1", "task2", "task3"}
    
    for _, task := range tasks {
        wg.Add(1)
        go func(t string) {
            defer wg.Done()
            
            fmt.Printf("Processing %s...\n", t)
            time.Sleep(500 * time.Millisecond)
            fmt.Printf("Completed %s\n", t)
        }(task)
    }
    
    wg.Wait()
    fmt.Println("All tasks done!")
}
```

### Сбор результатов от горутин

```go
package main

import (
    "fmt"
    "sync"
)

func factorial(n int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    result := 1
    for i := 2; i <= n; i++ {
        result *= i
    }
    results <- result
}

func main() {
    var wg sync.WaitGroup
    results := make(chan int, 5)  // буферизированный канал
    
    numbers := []int{3, 5, 7, 10, 12}
    
    for _, n := range numbers {
        wg.Add(1)
        go factorial(n, results, &wg)
    }
    
    // Закрываем канал после завершения всех горутин
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Читаем результаты
    for result := range results {
        fmt.Println("Factorial:", result)
    }
}
```

### Паттерн: Add перед запуском горутины

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    
    // ✅ ПРАВИЛЬНО: Add вызывается ДО запуска горутины
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Println("Worker", id)
        }(i)
    }
    
    wg.Wait()
    
    // ❌ НЕПРАВИЛЬНО: Add внутри горутины
    // for i := 0; i < 3; i++ {
    //     go func(id int) {
    //         wg.Add(1)  // гонка! Wait может выполниться раньше
    //         defer wg.Done()
    //         fmt.Println("Worker", id)
    //     }(i)
    // }
}
```

### Параллельная обработка с ограничением

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func processItem(id int, sem chan struct{}, wg *sync.WaitGroup) {
    defer wg.Done()
    
    sem <- struct{}{}        // занимаем слот
    defer func() { <-sem }() // освобождаем слот
    
    fmt.Printf("Processing item %d\n", id)
    time.Sleep(time.Second)
    fmt.Printf("Done item %d\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    // Семафор: максимум 3 одновременных горутины
    semaphore := make(chan struct{}, 3)
    
    for i := 1; i <= 10; i++ {
        wg.Add(1)
        go processItem(i, semaphore, &wg)
    }
    
    wg.Wait()
    fmt.Println("All items processed!")
}
```

### Вложенные WaitGroups

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func task(id int, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Printf("Task %d running\n", id)
    time.Sleep(100 * time.Millisecond)
}

func batch(batchID int, taskCount int, mainWg *sync.WaitGroup) {
    defer mainWg.Done()
    
    var batchWg sync.WaitGroup
    
    fmt.Printf("Batch %d: starting %d tasks\n", batchID, taskCount)
    
    for i := 1; i <= taskCount; i++ {
        batchWg.Add(1)
        go task(i, &batchWg)
    }
    
    batchWg.Wait()
    fmt.Printf("Batch %d: completed\n", batchID)
}

func main() {
    var mainWg sync.WaitGroup
    
    for batch := 1; batch <= 3; batch++ {
        mainWg.Add(1)
        go batch(batch, 4, &mainWg)
    }
    
    mainWg.Wait()
    fmt.Println("All batches completed!")
}
```

### Практический пример: параллельный веб-скрапер

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Result struct {
    URL     string
    Status  string
    Latency time.Duration
}

func fetchURL(url string, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    
    start := time.Now()
    
    // Имитация HTTP-запроса
    time.Sleep(time.Duration(100+len(url)*10) * time.Millisecond)
    
    results <- Result{
        URL:     url,
        Status:  "200 OK",
        Latency: time.Since(start),
    }
}

func main() {
    urls := []string{
        "https://google.com",
        "https://github.com",
        "https://golang.org",
        "https://amazon.com",
        "https://microsoft.com",
    }
    
    var wg sync.WaitGroup
    results := make(chan Result, len(urls))
    
    start := time.Now()
    
    for _, url := range urls {
        wg.Add(1)
        go fetchURL(url, results, &wg)
    }
    
    // Закрываем канал когда все горутины завершены
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Собираем результаты
    fmt.Println("Results:")
    for result := range results {
        fmt.Printf("  %s - %s (%v)\n", result.URL, result.Status, result.Latency)
    }
    
    fmt.Printf("\nTotal time: %v\n", time.Since(start))
}
```

---

## ⚠️ Частые ошибки

### 1. Передача WaitGroup по значению

```go
// ❌ НЕПРАВИЛЬНО — передаётся копия
func worker(wg sync.WaitGroup) {
    defer wg.Done()  // изменяет копию!
}

// ✅ ПРАВИЛЬНО — передаётся указатель
func worker(wg *sync.WaitGroup) {
    defer wg.Done()
}
```

### 2. Add после запуска горутины

```go
// ❌ НЕПРАВИЛЬНО — гонка данных
go func() {
    wg.Add(1)  // может выполниться после wg.Wait()
    defer wg.Done()
}()
wg.Wait()

// ✅ ПРАВИЛЬНО — Add до запуска
wg.Add(1)
go func() {
    defer wg.Done()
}()
wg.Wait()
```

### 3. Отрицательный счётчик

```go
var wg sync.WaitGroup

wg.Add(1)
wg.Done()
wg.Done()  // ❌ PANIC: negative WaitGroup counter

// ✅ Количество Done должно равняться Add
wg.Add(2)
wg.Done()
wg.Done()  // OK
```

### 4. Забыли Done()

```go
func worker(wg *sync.WaitGroup) {
    // ❌ Забыли wg.Done() — программа зависнет
    fmt.Println("Working...")
}

func worker(wg *sync.WaitGroup) {
    defer wg.Done()  // ✅ Всегда используйте defer
    fmt.Println("Working...")
}
```

### 5. Повторное использование WaitGroup

```go
var wg sync.WaitGroup

// Первая группа
wg.Add(2)
go func() { wg.Done() }()
go func() { wg.Done() }()
wg.Wait()

// ⚠️ Можно использовать повторно ПОСЛЕ Wait()
wg.Add(1)
go func() { wg.Done() }()
wg.Wait()
```

---

## 📝 Практика

### Задача 1: Параллельное вычисление
Вычислите сумму квадратов чисел от 1 до N, разделив работу между несколькими горутинами.

### Задача 2: Пул воркеров
Создайте пул из 3 горутин-воркеров, которые обрабатывают 10 задач из канала.

### Задача 3: Параллельный поиск файлов
Имитируйте параллельный поиск файлов в нескольких директориях с использованием WaitGroup.

### Задача 4: Таймаут
Реализуйте WaitGroup с таймаутом — ждать максимум N секунд.

### Задача 5: Прогресс-бар
Выводите прогресс выполнения: "Completed 3/10" при завершении каждой горутины.

### Задача 6: Сбор ошибок
Соберите все ошибки от горутин в срез, используя мьютекс для безопасности.

### Задача 7: Pipeline
Создайте конвейер из 3 стадий обработки, где каждая стадия — отдельная горутина.

### Задача 8: Rate limiter
Реализуйте ограничитель запросов: максимум N запросов в секунду.
