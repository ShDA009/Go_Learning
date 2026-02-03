# Table-Driven Tests (Табличные тесты)

---

## 💡 Ключевые идеи

1. **Table-driven tests** — идиоматический подход к тестированию в Go
2. **Таблица тестов** — срез структур с входными данными и ожидаемыми результатами
3. **Преимущества** — легко добавлять новые случаи, код DRY, читаемость
4. **Субтесты** — `t.Run()` для именованных подтестов
5. **Параллельные тесты** — `t.Parallel()` для ускорения выполнения

### Структура табличного теста

```
┌─────────────────────────────────────────────────────────┐
│  tests := []struct{                                     │
│      name     string    // имя тест-кейса               │
│      input    T         // входные данные               │
│      expected R         // ожидаемый результат          │
│  }{                                                     │
│      {"case1", input1, expected1},                      │
│      {"case2", input2, expected2},                      │
│      ...                                                │
│  }                                                      │
│                                                         │
│  for _, tt := range tests {                             │
│      t.Run(tt.name, func(t *testing.T) {                │
│          result := Function(tt.input)                   │
│          if result != tt.expected { ... }               │
│      })                                                 │
│  }                                                      │
└─────────────────────────────────────────────────────────┘
```

---

## 📖 Теория

### Проблема обычных тестов

Представьте, что вы тестируете функцию сложения. Вам нужно проверить:
- Сложение двух положительных чисел
- Сложение с нулём
- Сложение отрицательных чисел
- Сложение больших чисел

При обычном подходе вы напишете 4 отдельных теста с почти одинаковым кодом:

```go
func TestAddPositive(t *testing.T) { ... }
func TestAddZero(t *testing.T) { ... }
func TestAddNegative(t *testing.T) { ... }
func TestAddLarge(t *testing.T) { ... }
```

Это приводит к **дублированию кода** — если изменится сигнатура функции, придётся менять все 4 теста.

### Решение: Table-Driven Tests

**Table-driven tests** (табличные тесты) — это подход, при котором все тестовые случаи описываются в виде **таблицы** (среза структур), а затем обрабатываются в одном цикле.

Это как Excel-таблица:
| Название | Вход A | Вход B | Ожидаемый результат |
|----------|--------|--------|---------------------|
| positive | 2 | 3 | 5 |
| zero | 0 | 5 | 5 |
| negative | -1 | -1 | -2 |

### Почему это идиоматичный подход в Go?

Go-сообщество любит табличные тесты по нескольким причинам:

1. **DRY (Don't Repeat Yourself)** — логика проверки пишется один раз
2. **Легко добавлять случаи** — просто добавьте строку в таблицу
3. **Читаемость** — все тестовые случаи видны в одном месте
4. **Отладка** — `t.Run()` позволяет запускать конкретный тест-кейс

### Что такое субтесты (t.Run)?

`t.Run()` создаёт **именованный подтест**. Это даёт важные преимущества:

1. **Понятный вывод** — при падении теста видно, какой именно случай упал
2. **Выборочный запуск** — можно запустить только один подтест: `go test -run TestAdd/negative`
3. **Параллельность** — подтесты можно запускать параллельно

### Параллельные тесты

Вызов `t.Parallel()` говорит Go: "Этот тест можно запускать одновременно с другими".

⚠️ **Важно:** При параллельном выполнении переменная цикла должна быть скопирована:

```go
for _, tt := range tests {
    tt := tt  // ОБЯЗАТЕЛЬНО! Создаём локальную копию
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        // используем tt
    })
}
```

Без этой строки все горутины будут использовать одну и ту же переменную, и тесты будут работать неправильно.

### Когда использовать табличные тесты?

✅ **Используйте когда:**
- Тестируете функцию с разными входными данными
- Есть много edge cases
- Логика проверки одинаковая

❌ **Не используйте когда:**
- Каждый тест требует уникальной настройки
- Тестов всего 1-2
- Логика проверки сильно отличается

---

## 📋 Синтаксис

### Базовый табличный тест

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
    }{
        {"описание случая", input, expected},
        // ...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Function(tt.input)
            if result != tt.expected {
                t.Errorf("got %v; want %v", result, tt.expected)
            }
        })
    }
}
```

### Субтесты (t.Run)

```go
t.Run("имя подтеста", func(t *testing.T) {
    // код теста
})
```

### Параллельные тесты

```go
t.Run(tt.name, func(t *testing.T) {
    t.Parallel()  // этот субтест выполняется параллельно
    // ...
})
```

---

## 💻 Примеры кода

### Пример 1: Простой табличный тест

```go
package math

import "testing"

func Abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}

func TestAbs(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected int
    }{
        {"positive", 5, 5},
        {"negative", -5, 5},
        {"zero", 0, 0},
        {"large negative", -1000000, 1000000},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Abs(tt.input)
            if result != tt.expected {
                t.Errorf("Abs(%d) = %d; want %d", tt.input, result, tt.expected)
            }
        })
    }
}
```

### Пример 2: Тест функции с несколькими параметрами

```go
package math

import "testing"

func Max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func TestMax(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"first greater", 10, 5, 10},
        {"second greater", 5, 10, 10},
        {"equal", 7, 7, 7},
        {"negative numbers", -5, -10, -5},
        {"mixed signs", -5, 5, 5},
        {"zero and positive", 0, 5, 5},
        {"zero and negative", 0, -5, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Max(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Max(%d, %d) = %d; want %d", 
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Пример 3: Тест с ошибками

```go
package converter

import (
    "errors"
    "strconv"
    "testing"
)

var ErrNegative = errors.New("negative numbers not allowed")

func ParsePositiveInt(s string) (int, error) {
    n, err := strconv.Atoi(s)
    if err != nil {
        return 0, err
    }
    if n < 0 {
        return 0, ErrNegative
    }
    return n, nil
}

func TestParsePositiveInt(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expected    int
        expectError bool
        errorType   error
    }{
        {"valid positive", "42", 42, false, nil},
        {"zero", "0", 0, false, nil},
        {"negative", "-5", 0, true, ErrNegative},
        {"invalid string", "abc", 0, true, nil},
        {"empty string", "", 0, true, nil},
        {"large number", "999999", 999999, false, nil},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParsePositiveInt(tt.input)
            
            // Проверяем наличие ошибки
            if tt.expectError {
                if err == nil {
                    t.Errorf("expected error, got nil")
                    return
                }
                // Проверяем тип ошибки, если указан
                if tt.errorType != nil && !errors.Is(err, tt.errorType) {
                    t.Errorf("expected error %v, got %v", tt.errorType, err)
                }
                return
            }
            
            // Ошибки не ожидалось
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            
            if result != tt.expected {
                t.Errorf("got %d; want %d", result, tt.expected)
            }
        })
    }
}
```

### Пример 4: Тест строковых функций

```go
package stringutil

import (
    "strings"
    "testing"
)

func Capitalize(s string) string {
    if s == "" {
        return ""
    }
    return strings.ToUpper(s[:1]) + s[1:]
}

func TestCapitalize(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"lowercase word", "hello", "Hello"},
        {"uppercase word", "HELLO", "HELLO"},
        {"mixed case", "hELLO", "HELLO"},
        {"single char", "a", "A"},
        {"empty string", "", ""},
        {"with spaces", "hello world", "Hello world"},
        {"number start", "123abc", "123abc"},
        {"unicode", "привет", "Привет"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Capitalize(tt.input)
            if result != tt.expected {
                t.Errorf("Capitalize(%q) = %q; want %q", 
                    tt.input, result, tt.expected)
            }
        })
    }
}
```

### Пример 5: Параллельные табличные тесты

```go
package math

import (
    "testing"
    "time"
)

func SlowSquare(n int) int {
    time.Sleep(100 * time.Millisecond)  // имитация медленной операции
    return n * n
}

func TestSlowSquareParallel(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected int
    }{
        {"zero", 0, 0},
        {"positive", 5, 25},
        {"negative", -3, 9},
        {"one", 1, 1},
        {"large", 100, 10000},
    }
    
    for _, tt := range tests {
        tt := tt  // важно! создаём локальную копию для горутины
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // запускаем параллельно
            
            result := SlowSquare(tt.input)
            if result != tt.expected {
                t.Errorf("SlowSquare(%d) = %d; want %d", 
                    tt.input, result, tt.expected)
            }
        })
    }
}
```

### Пример 6: Тест с map входных данных

```go
package validator

import "testing"

func IsValidUsername(username string) bool {
    if len(username) < 3 || len(username) > 20 {
        return false
    }
    for _, r := range username {
        if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
             (r >= '0' && r <= '9') || r == '_') {
            return false
        }
    }
    return true
}

func TestIsValidUsername(t *testing.T) {
    tests := map[string]struct {
        input    string
        expected bool
    }{
        "valid lowercase":    {"john", true},
        "valid with numbers": {"john123", true},
        "valid with underscore": {"john_doe", true},
        "too short":          {"ab", false},
        "too long":           {"abcdefghijklmnopqrstuvwxyz", false},
        "empty":              {"", false},
        "with space":         {"john doe", false},
        "with special char":  {"john@doe", false},
        "minimum length":     {"abc", true},
        "maximum length":     {"abcdefghijklmnopqrst", true},
    }
    
    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            result := IsValidUsername(tt.input)
            if result != tt.expected {
                t.Errorf("IsValidUsername(%q) = %v; want %v", 
                    tt.input, result, tt.expected)
            }
        })
    }
}
```

### Пример 7: Комплексный тест структуры

```go
package user

import (
    "testing"
    "time"
)

type User struct {
    Name      string
    Email     string
    Age       int
    CreatedAt time.Time
}

func NewUser(name, email string, age int) (*User, error) {
    if name == "" {
        return nil, errors.New("name is required")
    }
    if email == "" {
        return nil, errors.New("email is required")
    }
    if age < 0 || age > 150 {
        return nil, errors.New("invalid age")
    }
    
    return &User{
        Name:      name,
        Email:     email,
        Age:       age,
        CreatedAt: time.Now(),
    }, nil
}

func TestNewUser(t *testing.T) {
    tests := []struct {
        name        string
        userName    string
        email       string
        age         int
        expectError bool
        checkUser   func(*User) bool
    }{
        {
            name:        "valid user",
            userName:    "John",
            email:       "john@example.com",
            age:         30,
            expectError: false,
            checkUser: func(u *User) bool {
                return u.Name == "John" && 
                       u.Email == "john@example.com" && 
                       u.Age == 30
            },
        },
        {
            name:        "empty name",
            userName:    "",
            email:       "john@example.com",
            age:         30,
            expectError: true,
        },
        {
            name:        "empty email",
            userName:    "John",
            email:       "",
            age:         30,
            expectError: true,
        },
        {
            name:        "negative age",
            userName:    "John",
            email:       "john@example.com",
            age:         -1,
            expectError: true,
        },
        {
            name:        "age too high",
            userName:    "John",
            email:       "john@example.com",
            age:         200,
            expectError: true,
        },
        {
            name:        "zero age",
            userName:    "Baby",
            email:       "baby@example.com",
            age:         0,
            expectError: false,
            checkUser: func(u *User) bool {
                return u.Age == 0
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := NewUser(tt.userName, tt.email, tt.age)
            
            if tt.expectError {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }
            
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            
            if tt.checkUser != nil && !tt.checkUser(user) {
                t.Errorf("user validation failed: %+v", user)
            }
        })
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Забыли скопировать переменную для параллельных тестов

```go
// ❌ НЕПРАВИЛЬНО — все горутины используют одну переменную
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        // tt может измениться к моменту выполнения!
        result := Function(tt.input)
    })
}

// ✅ ПРАВИЛЬНО — создаём локальную копию
for _, tt := range tests {
    tt := tt  // создаём копию!
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        result := Function(tt.input)
    })
}
```

### 2. Нечитаемые имена тестов

```go
// ❌ ПЛОХО
{"test1", 5, 25},
{"test2", -3, 9},

// ✅ ХОРОШО — описательные имена
{"positive number", 5, 25},
{"negative number", -3, 9},
```

### 3. Слишком много логики в таблице

```go
// ❌ ПЛОХО — сложно читать
tests := []struct {
    name     string
    input    string
    setup    func()
    validate func(result) bool
    cleanup  func()
    // ... много полей
}

// ✅ ХОРОШО — выносим сложную логику в отдельные тесты
```

### 4. Не проверяем все случаи

```go
// ❌ Забыли граничные случаи
tests := []struct{...}{
    {"positive", 5, 25},
}

// ✅ Покрываем граничные случаи
tests := []struct{...}{
    {"positive", 5, 25},
    {"negative", -5, 25},
    {"zero", 0, 0},
    {"one", 1, 1},
    {"max int", math.MaxInt32, /* ... */},
}
```

---

## 🏋️ Практические задания

### Задание 1: FizzBuzz с табличными тестами

Напишите функцию `FizzBuzz(n int) string` и полный табличный тест.

**Начальный код:**
```go
package main

import (
    "fmt"
    "testing"
)

// TODO: Напишите тест

func main() {
    // Ваш код здесь
    
}
```

```go
// fizzbuzz_test.go
package fizzbuzz

import "testing"

func TestFizzBuzz(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected string
    }{
        // Добавьте тестовые случаи
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Реализуйте проверку
        })
    }
}
```

**Ожидаемый вывод:**
```
=== RUN   TestFizzBuzz
=== RUN   TestFizzBuzz/fizzbuzz_15
=== RUN   TestFizzBuzz/fizz_3
=== RUN   TestFizzBuzz/buzz_5
=== RUN   TestFizzBuzz/number_7
--- PASS: TestFizzBuzz (0.00s)
PASS
```

**Баллы:** 15

### Задание 2: Валидация пароля

Создайте функцию валидации пароля с табличными тестами.

**Начальный код:**
```go
package main

import (
    "fmt"
    "testing"
)

// TODO: Напишите тест

func main() {
    // Ваш код здесь
    
}
```

```go
// password_test.go
package password

import "testing"

func TestValidate(t *testing.T) {
    tests := []struct {
        name       string
        password   string
        wantErrors int // количество ошибок
    }{
        {"valid", "Password123", 0},
        {"too_short", "Pass1", 1},
        {"no_digit", "Password", 1},
        {"no_upper", "password123", 1},
        {"all_wrong", "abc", 3},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            errors := Validate(tt.password)
            if len(errors) != tt.wantErrors {
                t.Errorf("Validate(%q) returned %d errors; want %d",
                    tt.password, len(errors), tt.wantErrors)
            }
        })
    }
}
```

**Ожидаемый вывод:**
```
=== RUN   TestValidate
=== RUN   TestValidate/valid
=== RUN   TestValidate/too_short
=== RUN   TestValidate/no_digit
=== RUN   TestValidate/no_upper
=== RUN   TestValidate/all_wrong
--- PASS: TestValidate (0.00s)
PASS
```

**Баллы:** 15

### Задание 3: Параллельные субтесты

Добавьте параллельное выполнение в табличный тест.

**Начальный код:**
```go
// converter.go
package converter

func CelsiusToFahrenheit(c float64) float64 {
    return c*9/5 + 32
}
```

```go
// converter_test.go
package converter

import "testing"

func TestCelsiusToFahrenheit(t *testing.T) {
    tests := []struct {
        name    string
        celsius float64
        want    float64
    }{
        {"freezing", 0, 32},
        {"boiling", 100, 212},
        {"body_temp", 37, 98.6},
        {"negative", -40, -40},
    }
    
    for _, tt := range tests {
        tt := tt // capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // добавьте это
            got := CelsiusToFahrenheit(tt.celsius)
            if got != tt.want {
                t.Errorf("CelsiusToFahrenheit(%v) = %v; want %v",
                    tt.celsius, got, tt.want)
            }
        })
    }
}
```

**Ожидаемый вывод:**
```
=== RUN   TestCelsiusToFahrenheit
=== RUN   TestCelsiusToFahrenheit/freezing
=== PAUSE TestCelsiusToFahrenheit/freezing
=== RUN   TestCelsiusToFahrenheit/boiling
=== PAUSE TestCelsiusToFahrenheit/boiling
=== CONT  TestCelsiusToFahrenheit/freezing
=== CONT  TestCelsiusToFahrenheit/boiling
--- PASS: TestCelsiusToFahrenheit (0.00s)
PASS
```

**Баллы:** 10

### Задание 4: Тест с map

Перепишите тест FizzBuzz используя map вместо slice.

**Начальный код:**
```go
func TestFizzBuzzMap(t *testing.T) {
    tests := map[string]struct {
        input    int
        expected string
    }{
        // Используйте имя теста как ключ map
    }
    
    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            // Реализуйте
        })
    }
}
```

**Баллы:** 10

### Задание 5: Генерация тестовых данных

Создайте функцию-генератор тестовых данных для факториала.

**Начальный код:**
```go
// factorial.go
package factorial

func Factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * Factorial(n-1)
}
```

```go
// factorial_test.go
package factorial

import "testing"

func generateTestCases() []struct {
    n    int
    want int
} {
    // Сгенерируйте случаи: 0!, 1!, 2!, ..., 10!
    return nil
}

func TestFactorial(t *testing.T) {
    for _, tt := range generateTestCases() {
        t.Run(fmt.Sprintf("n=%d", tt.n), func(t *testing.T) {
            if got := Factorial(tt.n); got != tt.want {
                t.Errorf("Factorial(%d) = %d; want %d", tt.n, got, tt.want)
            }
        })
    }
}
```

**Ожидаемый вывод:**
```
=== RUN   TestFactorial
=== RUN   TestFactorial/n=0
=== RUN   TestFactorial/n=1
=== RUN   TestFactorial/n=5
=== RUN   TestFactorial/n=10
--- PASS: TestFactorial (0.00s)
PASS
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Go Blog — Using Subtests](https://go.dev/blog/subtests)
- [Go Wiki — Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Dave Cheney — Writing Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
