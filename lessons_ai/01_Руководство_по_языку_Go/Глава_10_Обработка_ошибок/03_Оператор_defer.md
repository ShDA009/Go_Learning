# Оператор defer и отложенное выполнение

## 💡 Ключевые идеи

1. **defer** — откладывает выполнение функции до выхода из текущей функции
2. **LIFO порядок** — несколько defer выполняются в обратном порядке (стек)
3. **Гарантия выполнения** — defer срабатывает даже при panic
4. **Оценка аргументов** — аргументы вычисляются сразу, не при выполнении
5. **Основное применение** — освобождение ресурсов, unlock мьютексов, закрытие файлов
6. **Производительность** — минимальный overhead в современном Go

---

## 📋 Синтаксис

### Базовый defer

```go
defer function()          // вызов функции
defer fmt.Println("end")  // встроенные функции
defer file.Close()        // методы
```

### Defer с анонимной функцией

```go
defer func() {
    // код, который выполнится при выходе
}()
```

### Defer с параметрами

```go
x := 10
defer fmt.Println(x)  // выведет 10, не текущее значение x при выполнении!
```

---

## 💻 Примеры кода

### Базовый пример

```go
package main

import "fmt"

func main() {
    fmt.Println("Start")
    
    defer fmt.Println("This will print last")
    
    fmt.Println("Middle")
    fmt.Println("End")
}
// Output:
// Start
// Middle
// End
// This will print last
```

### Порядок выполнения (LIFO)

```go
package main

import "fmt"

func main() {
    defer fmt.Println("First defer - executes last")
    defer fmt.Println("Second defer - executes second")
    defer fmt.Println("Third defer - executes first")
    
    fmt.Println("Main function body")
}
// Output:
// Main function body
// Third defer - executes first
// Second defer - executes second
// First defer - executes last
```

### Закрытие файла

```go
package main

import (
    "fmt"
    "os"
)

func readFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()  // ✅ Закроется при любом выходе из функции
    
    // Читаем файл...
    buf := make([]byte, 100)
    n, err := file.Read(buf)
    if err != nil {
        return err  // file.Close() всё равно вызовется
    }
    
    fmt.Printf("Read %d bytes\n", n)
    return nil
}

func main() {
    err := readFile("test.txt")
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Unlock мьютекса

```go
package main

import (
    "fmt"
    "sync"
)

var (
    counter int
    mutex   sync.Mutex
)

func increment() {
    mutex.Lock()
    defer mutex.Unlock()  // ✅ Гарантированно разблокируется
    
    counter++
    
    if counter > 100 {
        return  // mutex.Unlock() всё равно вызовется
    }
    
    // Какая-то логика...
}

func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            increment()
        }()
    }
    
    wg.Wait()
    fmt.Println("Counter:", counter)
}
```

### Аргументы вычисляются сразу!

```go
package main

import "fmt"

func main() {
    x := 10
    
    defer fmt.Println("Deferred x:", x)  // x = 10 вычисляется СЕЙЧАС
    
    x = 20
    x = 30
    
    fmt.Println("Current x:", x)
}
// Output:
// Current x: 30
// Deferred x: 10  (не 30!)
```

### Захват переменной через замыкание

```go
package main

import "fmt"

func main() {
    x := 10
    
    // Анонимная функция захватывает переменную
    defer func() {
        fmt.Println("Deferred x:", x)  // использует текущее значение x
    }()
    
    x = 20
    x = 30
    
    fmt.Println("Current x:", x)
}
// Output:
// Current x: 30
// Deferred x: 30  (замыкание видит текущее значение)
```

### Defer в цикле

```go
package main

import "fmt"

func main() {
    // ⚠️ Все defer накапливаются!
    for i := 0; i < 5; i++ {
        defer fmt.Println("Iteration:", i)
    }
    fmt.Println("Loop done")
}
// Output:
// Loop done
// Iteration: 4
// Iteration: 3
// Iteration: 2
// Iteration: 1
// Iteration: 0
```

### Defer для измерения времени

```go
package main

import (
    "fmt"
    "time"
)

func measureTime(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}

func slowOperation() {
    defer measureTime("slowOperation")()
    
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Working...")
}

func main() {
    slowOperation()
}
// Output:
// Working...
// slowOperation took 100.xxxms
```

### Defer и именованные возвращаемые значения

```go
package main

import "fmt"

func triple(x int) (result int) {
    defer func() {
        result *= 3  // Можем изменить возвращаемое значение!
    }()
    
    result = x
    return  // result = x, потом defer изменит его на x * 3
}

func main() {
    fmt.Println(triple(10))  // 30
}
```

### Восстановление после panic

```go
package main

import "fmt"

func mayPanic() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from:", r)
        }
    }()
    
    panic("something went wrong!")
    
    fmt.Println("This won't print")
}

func main() {
    mayPanic()
    fmt.Println("Program continues...")
}
// Output:
// Recovered from: something went wrong!
// Program continues...
```

### Практический пример: транзакция БД

```go
package main

import "fmt"

type Transaction struct {
    committed bool
}

func (t *Transaction) Commit() {
    t.committed = true
    fmt.Println("Transaction committed")
}

func (t *Transaction) Rollback() {
    if !t.committed {
        fmt.Println("Transaction rolled back")
    }
}

func performDatabaseOperation() error {
    tx := &Transaction{}
    defer tx.Rollback()  // Откат, если не было commit
    
    // Операция 1
    fmt.Println("Operation 1")
    
    // Операция 2 - ошибка!
    if true {  // имитация ошибки
        return fmt.Errorf("operation 2 failed")
    }
    
    tx.Commit()  // Не дойдём сюда
    return nil
}

func main() {
    err := performDatabaseOperation()
    if err != nil {
        fmt.Println("Error:", err)
    }
}
// Output:
// Operation 1
// Transaction rolled back
// Error: operation 2 failed
```

---

## ⚠️ Частые ошибки

### 1. Defer в цикле для ресурсов

```go
// ❌ ПЛОХО — все файлы закроются только в конце функции
func processFiles(paths []string) {
    for _, path := range paths {
        f, _ := os.Open(path)
        defer f.Close()  // накапливаются!
        // процессинг...
    }
}

// ✅ ХОРОШО — используйте отдельную функцию
func processFile(path string) {
    f, _ := os.Open(path)
    defer f.Close()
    // процессинг...
}

func processFiles(paths []string) {
    for _, path := range paths {
        processFile(path)
    }
}
```

### 2. Ошибка при вызове defer

```go
// ❌ ПЛОХО — ошибка открытия не проверена
f, _ := os.Open("file.txt")
defer f.Close()  // panic если f == nil

// ✅ ХОРОШО — defer после проверки ошибки
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close()
```

### 3. Непонимание момента вычисления аргументов

```go
// ❌ Ожидаем 10, получаем 0
i := 0
defer fmt.Println(i)  // i = 0 запоминается здесь!
i = 10

// ✅ Используйте замыкание
i := 0
defer func() { fmt.Println(i) }()
i = 10  // выведет 10
```

### 4. Возврат значения из defer

```go
// defer НЕ возвращает значения!
func example() int {
    defer func() int {
        return 42  // это значение никуда не идёт
    }()
    return 0
}
```

---

## 📝 Практика

### Задача 1: Таймер функции
Создайте функцию для измерения времени выполнения через defer.

### Задача 2: Trace функций
Реализуйте трассировку входа/выхода из функции.

### Задача 3: Cleanup
Напишите функцию создания временного файла с автоматическим удалением.

### Задача 4: Lock/Unlock
Создайте обёртку для мьютекса с автоматическим unlock.

### Задача 5: Транзакция
Реализуйте паттерн транзакции с автоматическим rollback при ошибке.

### Задача 6: Порядок defer
Напишите программу, демонстрирующую LIFO порядок.

### Задача 7: Счётчик вызовов
Используйте defer для подсчёта вызовов функции.

### Задача 8: HTTP handler
Создайте middleware для логирования запросов с использованием defer.
