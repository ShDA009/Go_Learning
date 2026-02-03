# Coverage и Best Practices

---

## 💡 Ключевые идеи

1. **Code Coverage** — процент кода, выполняемого тестами
2. **go test -cover** — встроенный инструмент измерения покрытия
3. **coverprofile** — файл с детальным отчётом
4. **go tool cover** — визуализация покрытия
5. **Best Practices** — правила написания качественных тестов

### Виды покрытия

| Тип | Описание |
|-----|----------|
| Statement | Какие строки выполнены |
| Branch | Какие ветви (if/else) пройдены |
| Function | Какие функции вызваны |
| Line | Покрытие по строкам |

---

## 📖 Теория

### Что такое покрытие кода (Code Coverage)?

**Покрытие кода** — это метрика, показывающая, какой процент вашего кода выполняется при запуске тестов. Если покрытие 80% — значит 80% строк кода были выполнены хотя бы один раз.

### Зачем измерять покрытие?

1. **Находить непротестированный код** — если функция не покрыта, там могут быть баги
2. **Оценивать качество тестов** — низкое покрытие = мало тестов
3. **Требования проекта** — многие команды устанавливают минимум (например, 80%)

### Какое покрытие считается хорошим?

| Покрытие | Оценка |
|----------|--------|
| < 50% | Плохо — много непротестированного кода |
| 50-70% | Средне — основные пути покрыты |
| 70-80% | Хорошо — большинство сценариев протестировано |
| > 80% | Отлично — но 100% не всегда нужно |

### 100% покрытие — это цель?

**Нет!** 100% покрытие не означает отсутствие багов. Покрытие показывает, что код **выполнялся**, но не что он **работает правильно**:

```go
func Add(a, b int) int {
    return a * b  // баг: умножение вместо сложения
}

func TestAdd(t *testing.T) {
    Add(2, 2)  // покрытие 100%, но баг не найден!
}
```

### На что обращать внимание?

1. **Критический код** — бизнес-логика, обработка платежей, авторизация
2. **Edge cases** — граничные случаи, обработка ошибок
3. **Сложные условия** — все ветви if/else/switch

### Best Practices: как писать хорошие тесты

**1. Один тест — одна проверка**
```go
// ❌ Плохо
func TestUser(t *testing.T) {
    // создание, обновление, удаление в одном тесте
}

// ✅ Хорошо
func TestUserCreate(t *testing.T) { ... }
func TestUserUpdate(t *testing.T) { ... }
func TestUserDelete(t *testing.T) { ... }
```

**2. Понятные имена тестов**
```go
// ❌ Плохо
func TestAdd1(t *testing.T) { ... }

// ✅ Хорошо
func TestAdd_PositiveNumbers(t *testing.T) { ... }
func TestAdd_NegativeNumbers(t *testing.T) { ... }
func TestAdd_WithZero(t *testing.T) { ... }
```

**3. Arrange-Act-Assert (AAA)**
```go
func TestCalculateTotal(t *testing.T) {
    // Arrange — подготовка
    items := []Item{{Price: 100}, {Price: 200}}
    
    // Act — действие
    total := CalculateTotal(items)
    
    // Assert — проверка
    assert.Equal(t, 300, total)
}
```

**4. Тестируйте поведение, не реализацию**
```go
// ❌ Плохо — привязка к реализации
assert.Equal(t, 3, len(cache.items))

// ✅ Хорошо — проверка поведения
item, found := cache.Get("key")
assert.True(t, found)
assert.Equal(t, expectedItem, item)
```

**5. Используйте table-driven tests**
```go
tests := []struct{
    name string
    input int
    want int
}{
    {"positive", 5, 25},
    {"zero", 0, 0},
    {"negative", -3, 9},
}
```

**6. Изолируйте тесты**
- Каждый тест должен работать независимо
- Не полагайтесь на порядок выполнения
- Очищайте состояние после теста

---

## 📋 Синтаксис

### Команды покрытия

```bash
# Показать процент покрытия
go test -cover ./...

# Сохранить профиль
go test -coverprofile=coverage.out ./...

# Покрытие по пакетам
go test -cover -coverpkg=./... ./...

# HTML отчёт
go tool cover -html=coverage.out -o coverage.html

# Покрытие функций
go tool cover -func=coverage.out

# Покрытие в процентах
go test -cover -covermode=count ./...
```

### Режимы покрытия

```bash
-covermode=set    # был ли код выполнен (по умолчанию)
-covermode=count  # сколько раз выполнен
-covermode=atomic # атомарный подсчёт (для параллельных тестов)
```

---

## 💻 Примеры кода

### Пример 1: Базовое измерение покрытия

```go
package calc

// calculator.go
func Add(a, b int) int {
    return a + b
}

func Subtract(a, b int) int {
    return a - b
}

func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func Multiply(a, b int) int {
    return a * b
}
```

```go
// calculator_test.go
package calc

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", result)
    }
}

func TestSubtract(t *testing.T) {
    result := Subtract(5, 3)
    if result != 2 {
        t.Errorf("Subtract(5, 3) = %d; want 2", result)
    }
}

func TestDivide(t *testing.T) {
    result, err := Divide(10, 2)
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if result != 5 {
        t.Errorf("Divide(10, 2) = %d; want 5", result)
    }
}

// Multiply не тестируется — покрытие будет неполным
```

```bash
$ go test -cover
PASS
coverage: 75.0% of statements
ok      calc    0.005s

$ go test -coverprofile=coverage.out
$ go tool cover -func=coverage.out
calc/calculator.go:4:   Add             100.0%
calc/calculator.go:8:   Subtract        100.0%
calc/calculator.go:12:  Divide          50.0%   # ветка с ошибкой не покрыта
calc/calculator.go:19:  Multiply        0.0%    # не тестируется
total:                  (statements)    75.0%
```

### Пример 2: Полное покрытие с ветвями

```go
package calc

import (
    "errors"
    "testing"
)

func TestDivideSuccess(t *testing.T) {
    result, err := Divide(10, 2)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != 5 {
        t.Errorf("got %d; want 5", result)
    }
}

func TestDivideByZero(t *testing.T) {
    _, err := Divide(10, 0)
    if err == nil {
        t.Error("expected error for division by zero")
    }
}

func TestMultiply(t *testing.T) {
    tests := []struct {
        a, b, want int
    }{
        {2, 3, 6},
        {0, 5, 0},
        {-2, 3, -6},
    }
    
    for _, tt := range tests {
        result := Multiply(tt.a, tt.b)
        if result != tt.want {
            t.Errorf("Multiply(%d, %d) = %d; want %d", 
                tt.a, tt.b, result, tt.want)
        }
    }
}
```

### Пример 3: Покрытие сложной логики

```go
package user

type User struct {
    Name   string
    Age    int
    Email  string
    Active bool
}

func ValidateUser(u *User) []string {
    var errors []string
    
    if u == nil {
        return []string{"user is nil"}
    }
    
    if u.Name == "" {
        errors = append(errors, "name is required")
    } else if len(u.Name) < 2 {
        errors = append(errors, "name too short")
    } else if len(u.Name) > 100 {
        errors = append(errors, "name too long")
    }
    
    if u.Age < 0 {
        errors = append(errors, "age cannot be negative")
    } else if u.Age > 150 {
        errors = append(errors, "age is unrealistic")
    }
    
    if u.Email == "" {
        errors = append(errors, "email is required")
    }
    
    return errors
}
```

```go
// Тесты для 100% покрытия
func TestValidateUser(t *testing.T) {
    tests := []struct {
        name       string
        user       *User
        wantErrors []string
    }{
        {
            name:       "nil user",
            user:       nil,
            wantErrors: []string{"user is nil"},
        },
        {
            name:       "valid user",
            user:       &User{Name: "John", Age: 30, Email: "john@example.com"},
            wantErrors: nil,
        },
        {
            name:       "empty name",
            user:       &User{Name: "", Age: 30, Email: "john@example.com"},
            wantErrors: []string{"name is required"},
        },
        {
            name:       "short name",
            user:       &User{Name: "J", Age: 30, Email: "john@example.com"},
            wantErrors: []string{"name too short"},
        },
        {
            name:       "long name",
            user:       &User{Name: strings.Repeat("a", 101), Age: 30, Email: "john@example.com"},
            wantErrors: []string{"name too long"},
        },
        {
            name:       "negative age",
            user:       &User{Name: "John", Age: -1, Email: "john@example.com"},
            wantErrors: []string{"age cannot be negative"},
        },
        {
            name:       "unrealistic age",
            user:       &User{Name: "John", Age: 200, Email: "john@example.com"},
            wantErrors: []string{"age is unrealistic"},
        },
        {
            name:       "empty email",
            user:       &User{Name: "John", Age: 30, Email: ""},
            wantErrors: []string{"email is required"},
        },
        {
            name: "multiple errors",
            user: &User{Name: "", Age: -1, Email: ""},
            wantErrors: []string{
                "name is required",
                "age cannot be negative",
                "email is required",
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            errors := ValidateUser(tt.user)
            
            if len(errors) != len(tt.wantErrors) {
                t.Errorf("got %d errors; want %d", len(errors), len(tt.wantErrors))
                return
            }
            
            for i, err := range errors {
                if err != tt.wantErrors[i] {
                    t.Errorf("error[%d] = %q; want %q", i, err, tt.wantErrors[i])
                }
            }
        })
    }
}
```

### Пример 4: HTML отчёт о покрытии

```bash
# Генерация профиля
go test -coverprofile=coverage.out ./...

# Создание HTML отчёта
go tool cover -html=coverage.out -o coverage.html

# Открытие в браузере
open coverage.html  # macOS
xdg-open coverage.html  # Linux
start coverage.html  # Windows
```

### Пример 5: Покрытие нескольких пакетов

```bash
# Покрытие всех пакетов
go test -coverprofile=coverage.out -coverpkg=./... ./...

# Исключение тестовых пакетов
go test -coverprofile=coverage.out -coverpkg=./internal/...,./pkg/... ./...
```

### Пример 6: CI/CD интеграция

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run tests with coverage
      run: go test -coverprofile=coverage.out -covermode=atomic ./...
    
    - name: Check coverage threshold
      run: |
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | tr -d '%')
        echo "Coverage: $COVERAGE%"
        if (( $(echo "$COVERAGE < 80" | bc -l) )); then
          echo "Coverage is below 80%"
          exit 1
        fi
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: coverage.out
```

---

## 📋 Best Practices

### 1. Именование тестов

```go
// ✅ ХОРОШО — описательные имена
func TestUserService_CreateUser_Success(t *testing.T) {}
func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {}
func TestUserService_CreateUser_InvalidInput(t *testing.T) {}

// ❌ ПЛОХО — непонятные имена
func TestUser1(t *testing.T) {}
func TestUserOK(t *testing.T) {}
```

### 2. Структура теста (AAA)

```go
func TestCalculator_Add(t *testing.T) {
    // Arrange — подготовка
    calc := NewCalculator()
    a, b := 2, 3
    expected := 5
    
    // Act — действие
    result := calc.Add(a, b)
    
    // Assert — проверка
    if result != expected {
        t.Errorf("Add(%d, %d) = %d; want %d", a, b, result, expected)
    }
}
```

### 3. Один тест — одна проверка

```go
// ✅ ХОРОШО — каждый тест проверяет одно
func TestUser_Name_Empty(t *testing.T) {
    user := User{Name: ""}
    err := user.Validate()
    assert.ErrorContains(t, err, "name")
}

func TestUser_Email_Invalid(t *testing.T) {
    user := User{Name: "John", Email: "invalid"}
    err := user.Validate()
    assert.ErrorContains(t, err, "email")
}

// ❌ ПЛОХО — много проверок в одном тесте
func TestUser_Validate(t *testing.T) {
    // проверка имени
    // проверка email
    // проверка возраста
    // ... 20 проверок
}
```

### 4. Независимость тестов

```go
// ✅ ХОРОШО — тест не зависит от других
func TestCreateUser(t *testing.T) {
    db := setupCleanDB(t)
    defer db.Close()
    
    user := &User{Name: "John"}
    err := CreateUser(db, user)
    
    assert.NoError(t, err)
}

// ❌ ПЛОХО — тест зависит от состояния
var globalUser *User

func TestCreateUser(t *testing.T) {
    globalUser = &User{Name: "John"}
    CreateUser(globalUser)
}

func TestGetUser(t *testing.T) {
    // зависит от TestCreateUser
    user := GetUser(globalUser.ID)
}
```

### 5. Тестовые данные

```go
// ✅ ХОРОШО — данные в тесте
func TestParseDate(t *testing.T) {
    input := "2024-01-15"
    expected := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
    
    result, _ := ParseDate(input)
    
    assert.Equal(t, expected, result)
}

// ✅ ХОРОШО — для сложных данных используйте testdata/
// testdata/users.json
func TestImportUsers(t *testing.T) {
    data, _ := os.ReadFile("testdata/users.json")
    // ...
}
```

### 6. Используйте t.Helper()

```go
// ✅ ХОРОШО
func assertEqual(t *testing.T, got, want int) {
    t.Helper()  // ошибка покажет строку вызова, а не эту строку
    if got != want {
        t.Errorf("got %d; want %d", got, want)
    }
}

func TestAdd(t *testing.T) {
    assertEqual(t, Add(2, 2), 4)  // ошибка укажет на эту строку
}
```

### 7. Параллельные тесты

```go
func TestParallel(t *testing.T) {
    tests := []struct {
        name string
        input int
    }{
        {"case1", 1},
        {"case2", 2},
        {"case3", 3},
    }
    
    for _, tt := range tests {
        tt := tt  // важно для параллельных тестов!
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // тест
        })
    }
}
```

### 8. Очистка ресурсов

```go
func TestWithCleanup(t *testing.T) {
    // Setup
    file, err := os.CreateTemp("", "test")
    if err != nil {
        t.Fatal(err)
    }
    
    // Cleanup — выполнится даже при падении теста
    t.Cleanup(func() {
        os.Remove(file.Name())
    })
    
    // Test
    // ...
}
```

### 9. Skip условных тестов

```go
func TestRequiresDocker(t *testing.T) {
    if os.Getenv("DOCKER_HOST") == "" {
        t.Skip("Docker not available")
    }
    // ...
}

func TestSlow(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping slow test in short mode")
    }
    // ...
}
```

### 10. Логирование в тестах

```go
func TestWithLogs(t *testing.T) {
    // Логи видны только при -v или провале теста
    t.Log("starting test")
    
    result := DoSomething()
    t.Logf("result: %v", result)
    
    if result != expected {
        t.Errorf("got %v; want %v", result, expected)
    }
}
```

---

## ⚠️ Анти-паттерны

### 1. Тестирование приватных деталей

```go
// ❌ ПЛОХО — тестируем внутреннюю реализацию
func TestInternalCache(t *testing.T) {
    service := NewService()
    // проверяем внутренний кэш
    assert.Len(t, service.cache, 0)
}

// ✅ ХОРОШО — тестируем поведение
func TestService_CachesResults(t *testing.T) {
    service := NewService()
    
    // Первый вызов
    result1 := service.GetData("key")
    
    // Второй вызов должен быть из кэша (быстрее)
    start := time.Now()
    result2 := service.GetData("key")
    duration := time.Since(start)
    
    assert.Equal(t, result1, result2)
    assert.Less(t, duration, time.Millisecond)
}
```

### 2. Хрупкие тесты

```go
// ❌ ПЛОХО — зависит от форматирования
func TestOutput(t *testing.T) {
    output := FormatUser(user)
    expected := "Name: John\nAge: 30\nEmail: john@example.com"
    assert.Equal(t, expected, output)
}

// ✅ ХОРОШО — проверяем содержание
func TestOutput(t *testing.T) {
    output := FormatUser(user)
    assert.Contains(t, output, "John")
    assert.Contains(t, output, "30")
    assert.Contains(t, output, "john@example.com")
}
```

### 3. Игнорирование ошибок

```go
// ❌ ПЛОХО
func TestSomething(t *testing.T) {
    result, _ := DoSomething()  // игнорируем ошибку!
    assert.Equal(t, expected, result)
}

// ✅ ХОРОШО
func TestSomething(t *testing.T) {
    result, err := DoSomething()
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

---

## 🏋️ Практические задания

### Задание 1: Достигните 80% покрытия

Напишите тесты для достижения 80% покрытия кода.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

```go
// calculator_test.go
package calculator

import (
    "testing"
)

// Напишите тесты для всех функций
// Запустите: go test -cover

func TestAdd(t *testing.T) {
    tests := []struct {
        a, b, want float64
    }{
        {1, 2, 3},
        {-1, 1, 0},
        {0, 0, 0},
    }
    
    for _, tt := range tests {
        got := Add(tt.a, tt.b)
        if got != tt.want {
            t.Errorf("Add(%v, %v) = %v; want %v", tt.a, tt.b, got, tt.want)
        }
    }
}

// Добавьте тесты для Subtract, Multiply, Divide, Power, Sqrt
```

**Запуск:**
```bash
go test -cover
# Или с детальным отчётом:
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Ожидаемый вывод:**
```
PASS
coverage: 85.0% of statements
```

**Баллы:** 15

### Задание 2: HTML отчёт покрытия

Сгенерируйте HTML отчёт и найдите непокрытые ветки.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

```bash
# Команды для генерации отчёта
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
# Откройте coverage.html в браузере
```

**Баллы:** 10

### Задание 3: Рефакторинг тестов по best practices

Перепишите плохие тесты согласно best practices.

**Начальный код (плохой):**
```go
// bad_test.go
package user

import "testing"

var globalUser *User

func TestCreate(t *testing.T) {
    globalUser = &User{Name: "John", Email: "john@example.com"}
    err := CreateUser(globalUser)
    if err != nil {
        t.Error(err)
    }
}

func TestGet(t *testing.T) {
    // Зависит от TestCreate!
    user, _ := GetUser(globalUser.ID)
    if user.Name != "John" {
        t.Error("wrong name")
    }
}

func TestUpdate(t *testing.T) {
    // Игнорируем ошибку
    UpdateUser(globalUser.ID, "Jane", "jane@example.com")
    user, _ := GetUser(globalUser.ID)
    if user.Email != "jane@example.com" {
        t.Error("update failed")
    }
}
```

**Перепишите согласно best practices:**
```go
// good_test.go
package user

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (*UserService, func()) {
    t.Helper()
    
    db := setupTestDB(t)
    service := NewUserService(db)
    
    cleanup := func() {
        db.Exec("DELETE FROM users")
        db.Close()
    }
    
    return service, cleanup
}

func TestUserService_Create(t *testing.T) {
    service, cleanup := setupTest(t)
    defer cleanup()
    
    // Arrange
    input := CreateUserInput{
        Name:  "John",
        Email: "john@example.com",
    }
    
    // Act
    user, err := service.Create(input)
    
    // Assert
    require.NoError(t, err)
    assert.NotZero(t, user.ID)
    assert.Equal(t, "John", user.Name)
}

// Каждый тест независим и создаёт свои данные
func TestUserService_Get(t *testing.T) {
    service, cleanup := setupTest(t)
    defer cleanup()
    
    // Setup - создаём данные для этого теста
    created, err := service.Create(CreateUserInput{
        Name:  "Jane",
        Email: "jane@example.com",
    })
    require.NoError(t, err)
    
    // Act
    user, err := service.GetByID(created.ID)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, "Jane", user.Name)
}
```

**Баллы:** 15

### Задание 4: GitHub Actions для тестов

Настройте CI с проверкой тестов и покрытия.

**Начальный код:**
```yaml
# .github/workflows/test.yml
name: Test

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Check coverage
      run: |
        coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        echo "Coverage: $coverage%"
        if (( $(echo "$coverage < 80" | bc -l) )); then
          echo "Coverage is below 80%"
          exit 1
        fi
    
    - name: Upload coverage
      uses: codecov/codecov-action@v4
      with:
        files: ./coverage.out
```

**Баллы:** 10

### Задание 5: t.Cleanup и t.Helper

Используйте t.Cleanup и t.Helper для улучшения тестов.

**Начальный код:**
```go
package main

import (
    "errors"
    "fmt"
    "os"
    "sync"
    "time"
)

// TODO: Создайте функцию согласно заданию
// TODO: Запустите горутину
// TODO: Используйте defer

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Go Cover](https://go.dev/blog/cover)
- [Testing Techniques](https://go.dev/doc/code#Testing)
- [Codecov](https://codecov.io/)
- [go-mutesting](https://github.com/zimmski/go-mutesting)
