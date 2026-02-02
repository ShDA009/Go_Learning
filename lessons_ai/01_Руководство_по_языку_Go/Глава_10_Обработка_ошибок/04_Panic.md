# Оператор panic

## 💡 Ключевые идеи

1. **panic** — аварийное завершение программы при критической ошибке
2. **Раскрутка стека** — при panic выполняются все defer-функции
3. **Аргумент** — panic принимает любое значение (обычно строку или error)
4. **Использование** — только для невосстановимых ошибок (ошибки программиста)
5. **Не для управления потоком** — НЕ замена обычной обработки ошибок
6. **recover** — позволяет перехватить panic

---

## 📋 Синтаксис

### Вызов panic

```go
panic("error message")      // строка
panic(errors.New("error"))  // error
panic(123)                  // любое значение
```

### Panic в стандартной библиотеке

```go
// Многие встроенные операции вызывают panic
arr := []int{1, 2, 3}
_ = arr[10]  // panic: runtime error: index out of range

var m map[string]int
m["key"] = 1  // panic: assignment to entry in nil map
```

---

## 💻 Примеры кода

### Базовый пример panic

```go
package main

import "fmt"

func main() {
    fmt.Println("Start")
    
    panic("Something went terribly wrong!")
    
    fmt.Println("This will never print")
}
// Output:
// Start
// panic: Something went terribly wrong!
// 
// goroutine 1 [running]:
// main.main()
//     /main.go:7 +0x...
```

### Panic с defer

```go
package main

import "fmt"

func cleanup() {
    fmt.Println("Cleanup executed!")
}

func main() {
    defer cleanup()  // Выполнится даже при panic!
    
    fmt.Println("Start")
    panic("Critical error!")
    fmt.Println("Never reached")
}
// Output:
// Start
// Cleanup executed!
// panic: Critical error!
```

### Panic в функции

```go
package main

import "fmt"

func divide(a, b int) int {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}

func main() {
    defer fmt.Println("Main cleanup")
    
    fmt.Println("10 / 2 =", divide(10, 2))
    fmt.Println("10 / 0 =", divide(10, 0))  // panic!
    fmt.Println("Never reached")
}
```

### Порядок выполнения defer при panic

```go
package main

import "fmt"

func level3() {
    defer fmt.Println("level3 defer")
    fmt.Println("level3 start")
    panic("error in level3")
}

func level2() {
    defer fmt.Println("level2 defer")
    fmt.Println("level2 start")
    level3()
    fmt.Println("level2 end")  // не выполнится
}

func level1() {
    defer fmt.Println("level1 defer")
    fmt.Println("level1 start")
    level2()
    fmt.Println("level1 end")  // не выполнится
}

func main() {
    defer fmt.Println("main defer")
    fmt.Println("main start")
    level1()
    fmt.Println("main end")  // не выполнится
}
// Output:
// main start
// level1 start
// level2 start
// level3 start
// level3 defer
// level2 defer
// level1 defer
// main defer
// panic: error in level3
```

### Когда использовать panic

```go
package main

import "fmt"

// ✅ Хорошее использование: инициализация, которая обязана успешно завершиться
func mustParseConfig(path string) Config {
    config, err := parseConfig(path)
    if err != nil {
        panic(fmt.Sprintf("failed to parse config %s: %v", path, err))
    }
    return config
}

// ✅ Хорошее использование: невозможная ситуация
func processStatus(status int) {
    switch status {
    case 0:
        // OK
    case 1:
        // Warning
    case 2:
        // Error
    default:
        panic(fmt.Sprintf("unknown status: %d", status))
    }
}

// ❌ Плохое использование: обычная ошибка
func badExample(x int) int {
    if x < 0 {
        panic("negative number")  // Лучше вернуть error!
    }
    return x * 2
}

type Config struct{}

func parseConfig(path string) (Config, error) {
    return Config{}, nil
}

func main() {
    fmt.Println("Example")
}
```

### Panic с кастомным типом

```go
package main

import "fmt"

type CriticalError struct {
    Code    int
    Message string
}

func (e CriticalError) Error() string {
    return fmt.Sprintf("CRITICAL [%d]: %s", e.Code, e.Message)
}

func criticalOperation() {
    panic(CriticalError{
        Code:    500,
        Message: "Database connection lost",
    })
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            if err, ok := r.(CriticalError); ok {
                fmt.Printf("Caught critical error: %v\n", err)
            } else {
                fmt.Println("Caught unknown panic:", r)
            }
        }
    }()
    
    criticalOperation()
}
```

### Runtime panic

```go
package main

import "fmt"

func main() {
    // Index out of range
    arr := []int{1, 2, 3}
    // _ = arr[10]  // panic: runtime error: index out of range [10] with length 3
    
    // Nil pointer dereference
    var ptr *int
    // _ = *ptr  // panic: runtime error: invalid memory address or nil pointer dereference
    
    // Nil map write
    var m map[string]int
    // m["key"] = 1  // panic: assignment to entry in nil map
    
    // Type assertion failure
    var i interface{} = "hello"
    // _ = i.(int)  // panic: interface conversion: interface {} is string, not int
    
    fmt.Println("Program ended safely")
    _ = arr
    _ = ptr
    _ = m
    _ = i
}
```

### Panic vs return error

```go
package main

import (
    "errors"
    "fmt"
)

// Вариант 1: возврат error (рекомендуется)
func divideWithError(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Вариант 2: panic (только для критических случаев)
func mustDivide(a, b int) int {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}

func main() {
    // Предпочтительный способ
    result, err := divideWithError(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", result)
    }
    
    // Panic способ — для случаев когда ошибка невозможна
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Panic:", r)
        }
    }()
    
    _ = mustDivide(10, 0)
}
```

### Практический пример: инициализация приложения

```go
package main

import (
    "fmt"
    "os"
)

type App struct {
    Config   Config
    Database *Database
}

type Config struct {
    DBHost string
    DBPort int
}

type Database struct{}

func mustLoadConfig() Config {
    host := os.Getenv("DB_HOST")
    if host == "" {
        panic("DB_HOST environment variable is required")
    }
    
    return Config{
        DBHost: host,
        DBPort: 5432,
    }
}

func mustConnectDB(config Config) *Database {
    // Имитация подключения
    if config.DBHost == "" {
        panic("cannot connect to database: empty host")
    }
    
    fmt.Println("Connected to database at", config.DBHost)
    return &Database{}
}

func NewApp() *App {
    // При инициализации panic допустим
    config := mustLoadConfig()
    db := mustConnectDB(config)
    
    return &App{
        Config:   config,
        Database: db,
    }
}

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Failed to start application:", r)
            os.Exit(1)
        }
    }()
    
    // Для теста
    os.Setenv("DB_HOST", "localhost")
    
    app := NewApp()
    fmt.Printf("App started with config: %+v\n", app.Config)
}
```

---

## ⚠️ Частые ошибки

### 1. Использование panic вместо error

```go
// ❌ ПЛОХО
func GetUser(id int) *User {
    if id <= 0 {
        panic("invalid id")  // Обычная ошибка валидации
    }
    return &User{}
}

// ✅ ХОРОШО
func GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid id")
    }
    return &User{}, nil
}
```

### 2. Panic без recover в библиотеке

```go
// ❌ Библиотечная функция не должна паниковать
func LibraryFunction() {
    panic("error")  // Плохо для пользователей библиотеки
}

// ✅ Возвращайте error
func LibraryFunction() error {
    return errors.New("error")
}
```

### 3. Игнорирование defer при panic

```go
// ❌ Ресурс не будет освобождён
func bad() {
    file, _ := os.Open("file.txt")
    // ...
    panic("error")
    file.Close()  // Никогда не выполнится
}

// ✅ Используйте defer
func good() {
    file, _ := os.Open("file.txt")
    defer file.Close()  // Выполнится даже при panic
    // ...
    panic("error")
}
```

---

## 📝 Практика

### Задача 1: Must функция
Создайте `MustParseInt(s string) int`, которая паникует при ошибке парсинга.

### Задача 2: Валидация конфига
Напишите функцию инициализации, которая паникует при невалидном конфиге.

### Задача 3: Assert
Создайте функцию `Assert(condition bool, message string)`.

### Задача 4: Stack trace
Напишите программу, демонстрирующую раскрутку стека при panic.

### Задача 5: Типизированный panic
Создайте кастомный тип для panic с кодом ошибки.

### Задача 6: Защита от panic
Напишите функцию-обёртку, которая перехватывает panic из переданной функции.

### Задача 7: Safe call
Создайте `SafeCall(f func()) error`, преобразующую panic в error.

### Задача 8: Graceful shutdown
Реализуйте graceful shutdown при panic с логированием.
