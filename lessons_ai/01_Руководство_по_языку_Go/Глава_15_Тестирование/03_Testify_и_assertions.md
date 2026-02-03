# Testify и assertions

---

## 💡 Ключевые идеи

1. **testify** — популярная библиотека для тестирования в Go
2. **assert** — мягкие проверки (тест продолжается после провала)
3. **require** — строгие проверки (тест останавливается при провале)
4. **suite** — организация тестов в наборы с setup/teardown
5. **mock** — создание моков для интерфейсов

### Сравнение подходов

| Подход | Код | Информативность |
|--------|-----|-----------------|
| Стандартный | `if got != want { t.Errorf(...) }` | Нужно писать сообщение |
| testify/assert | `assert.Equal(t, want, got)` | Автоматическое сообщение |

---

## 📖 Теория

### Проблема стандартного тестирования

В стандартной библиотеке Go нет встроенных assertion-функций. Вы должны сами писать проверки:

```go
if result != expected {
    t.Errorf("got %d; want %d", result, expected)
}
```

Это работает, но:
- **Много шаблонного кода** — каждая проверка требует 3 строки
- **Нужно писать сообщения** — легко забыть или написать неинформативно
- **Нет готовых проверок** — для срезов, map, ошибок нужен свой код

### Что такое Testify?

**Testify** — это самая популярная библиотека для тестирования в Go. Она предоставляет:
- **assert** — функции для проверки условий
- **require** — то же, но останавливает тест при ошибке
- **mock** — инструменты для создания моков
- **suite** — организация тестов в наборы

### Assert vs Require: в чём разница?

Это ключевое различие, которое нужно понимать:

**assert** — "мягкая" проверка:
```go
assert.Equal(t, 5, result)  // если false — логирует ошибку
// тест ПРОДОЛЖАЕТСЯ
```

**require** — "строгая" проверка:
```go
require.Equal(t, 5, result)  // если false — логирует ошибку
// тест ОСТАНАВЛИВАЕТСЯ
```

### Когда что использовать?

**Используйте require когда:**
- Дальнейшие проверки зависят от этой
- Провал делает остальные проверки бессмысленными

```go
user, err := GetUser(1)
require.NoError(t, err)     // если ошибка — дальше нет смысла
require.NotNil(t, user)     // если nil — будет panic
assert.Equal(t, "John", user.Name)  // теперь безопасно
```

**Используйте assert когда:**
- Хотите увидеть ВСЕ ошибки теста сразу
- Проверки независимы друг от друга

### Информативные сообщения об ошибках

Одно из главных преимуществ testify — понятные сообщения:

```
--- FAIL: TestAdd (0.00s)
    add_test.go:15: 
        	Error Trace:	add_test.go:15
        	Error:      	Not equal: 
        	            	expected: 5
        	            	actual  : 4
        	Test:       	TestAdd
```

Без testify вам пришлось бы писать это сообщение вручную!

### Полезные функции testify

```go
// Сравнение
assert.Equal(t, expected, actual)      // ==
assert.NotEqual(t, a, b)               // !=
assert.Same(t, ptr1, ptr2)             // один и тот же указатель

// Коллекции
assert.Contains(t, "hello", "ell")     // строка содержит подстроку
assert.Contains(t, []int{1,2,3}, 2)    // срез содержит элемент
assert.Len(t, slice, 5)                // длина равна 5
assert.Empty(t, slice)                 // пустой

// Ошибки
assert.NoError(t, err)                 // err == nil
assert.Error(t, err)                   // err != nil
assert.ErrorIs(t, err, ErrNotFound)    // errors.Is()

// Паника
assert.Panics(t, func() { panic("!") })
```

---

## 📋 Синтаксис

### Установка

```bash
go get github.com/stretchr/testify
```

### Импорт

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

### Основные функции assert/require

```go
// Равенство
assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)

// Булевы значения
assert.True(t, value)
assert.False(t, value)

// Nil проверки
assert.Nil(t, value)
assert.NotNil(t, value)

// Ошибки
assert.NoError(t, err)
assert.Error(t, err)
assert.ErrorIs(t, err, targetErr)

// Содержание
assert.Contains(t, collection, element)
assert.NotContains(t, collection, element)

// Длина
assert.Len(t, collection, length)
assert.Empty(t, collection)
assert.NotEmpty(t, collection)

// Паника
assert.Panics(t, func() { /* код */ })
assert.NotPanics(t, func() { /* код */ })

// Типы
assert.IsType(t, expectedType, actual)

// С кастомным сообщением
assert.Equal(t, expected, actual, "custom message")
assert.Equalf(t, expected, actual, "value should be %d", expected)
```

---

## 💻 Примеры кода

### Пример 1: Базовое использование assert

```go
package math

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func Add(a, b int) int {
    return a + b
}

func TestAddWithAssert(t *testing.T) {
    // Вместо:
    // if result != 5 { t.Errorf("...") }
    
    // Используем assert:
    result := Add(2, 3)
    assert.Equal(t, 5, result)
    
    // С кастомным сообщением
    assert.Equal(t, 5, result, "2 + 3 should equal 5")
    
    // Несколько проверок
    assert.Equal(t, 0, Add(0, 0))
    assert.Equal(t, -2, Add(-5, 3))
    assert.Equal(t, 10, Add(7, 3))
}
```

### Пример 2: Разница между assert и require

```go
package user

import (
    "errors"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type User struct {
    Name string
    Age  int
}

func GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid id")
    }
    return &User{Name: "John", Age: 30}, nil
}

func TestGetUserWithAssert(t *testing.T) {
    user, err := GetUser(1)
    
    // assert — тест продолжится даже при ошибке
    assert.NoError(t, err)
    assert.NotNil(t, user)
    // Если user == nil, следующая строка вызовет panic!
    assert.Equal(t, "John", user.Name)
}

func TestGetUserWithRequire(t *testing.T) {
    user, err := GetUser(1)
    
    // require — тест остановится при первой ошибке
    require.NoError(t, err)
    require.NotNil(t, user)
    // Безопасно — сюда дойдём только если user != nil
    assert.Equal(t, "John", user.Name)
    assert.Equal(t, 30, user.Age)
}
```

### Пример 3: Проверка коллекций

```go
package collection

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func GetNumbers() []int {
    return []int{1, 2, 3, 4, 5}
}

func GetEmptySlice() []int {
    return []int{}
}

func GetMap() map[string]int {
    return map[string]int{
        "one":   1,
        "two":   2,
        "three": 3,
    }
}

func TestCollections(t *testing.T) {
    // Срезы
    nums := GetNumbers()
    assert.Len(t, nums, 5)
    assert.Contains(t, nums, 3)
    assert.NotContains(t, nums, 10)
    assert.NotEmpty(t, nums)
    
    // Пустой срез
    empty := GetEmptySlice()
    assert.Empty(t, empty)
    assert.Len(t, empty, 0)
    
    // Карты
    m := GetMap()
    assert.Len(t, m, 3)
    assert.Contains(t, m, "one")
    assert.NotContains(t, m, "four")
    
    // Элементы в порядке (для срезов)
    assert.ElementsMatch(t, []int{5, 4, 3, 2, 1}, nums)
}
```

### Пример 4: Проверка ошибок

```go
package validator

import (
    "errors"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

var (
    ErrEmpty    = errors.New("value is empty")
    ErrTooShort = errors.New("value is too short")
)

func Validate(s string) error {
    if s == "" {
        return ErrEmpty
    }
    if len(s) < 3 {
        return fmt.Errorf("validation failed: %w", ErrTooShort)
    }
    return nil
}

func TestValidate(t *testing.T) {
    // Нет ошибки
    err := Validate("hello")
    assert.NoError(t, err)
    
    // Есть ошибка
    err = Validate("")
    assert.Error(t, err)
    
    // Конкретный тип ошибки
    err = Validate("")
    assert.ErrorIs(t, err, ErrEmpty)
    
    // Обёрнутая ошибка
    err = Validate("ab")
    assert.ErrorIs(t, err, ErrTooShort)
    
    // Сообщение ошибки содержит текст
    assert.ErrorContains(t, err, "validation failed")
}
```

### Пример 5: Проверка паники

```go
package safety

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func MustGetIndex(slice []int, index int) int {
    if index < 0 || index >= len(slice) {
        panic("index out of range")
    }
    return slice[index]
}

func SafeGetIndex(slice []int, index int) (int, bool) {
    if index < 0 || index >= len(slice) {
        return 0, false
    }
    return slice[index], true
}

func TestPanic(t *testing.T) {
    nums := []int{1, 2, 3}
    
    // Проверяем, что паника происходит
    assert.Panics(t, func() {
        MustGetIndex(nums, 10)
    })
    
    // Проверяем, что паника НЕ происходит
    assert.NotPanics(t, func() {
        MustGetIndex(nums, 1)
    })
    
    // Проверяем значение паники
    assert.PanicsWithValue(t, "index out of range", func() {
        MustGetIndex(nums, -1)
    })
}
```

### Пример 6: Табличные тесты с testify

```go
package converter

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func CelsiusToFahrenheit(c float64) float64 {
    return c*9/5 + 32
}

func TestCelsiusToFahrenheit(t *testing.T) {
    tests := []struct {
        name     string
        celsius  float64
        expected float64
    }{
        {"freezing point", 0, 32},
        {"boiling point", 100, 212},
        {"body temperature", 37, 98.6},
        {"negative", -40, -40},
        {"room temperature", 20, 68},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CelsiusToFahrenheit(tt.celsius)
            assert.InDelta(t, tt.expected, result, 0.01, 
                "conversion of %.1f°C", tt.celsius)
        })
    }
}
```

### Пример 7: Проверка структур

```go
package user

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
)

type User struct {
    ID        int
    Name      string
    Email     string
    CreatedAt time.Time
}

func CreateUser(name, email string) *User {
    return &User{
        ID:        1,
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
    }
}

func TestCreateUser(t *testing.T) {
    user := CreateUser("John", "john@example.com")
    
    // Проверка отдельных полей
    assert.Equal(t, 1, user.ID)
    assert.Equal(t, "John", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
    assert.False(t, user.CreatedAt.IsZero())
    
    // Проверка типа
    assert.IsType(t, &User{}, user)
    
    // Проверка что время недавнее
    assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
}

func TestUserEquality(t *testing.T) {
    user1 := User{ID: 1, Name: "John", Email: "john@example.com"}
    user2 := User{ID: 1, Name: "John", Email: "john@example.com"}
    user3 := User{ID: 2, Name: "Jane", Email: "jane@example.com"}
    
    // Сравнение структур
    assert.Equal(t, user1, user2)
    assert.NotEqual(t, user1, user3)
    
    // Глубокое сравнение (для вложенных структур)
    assert.EqualValues(t, user1, user2)
}
```

### Пример 8: Comparison helpers

```go
package numbers

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestComparisons(t *testing.T) {
    // Числовые сравнения
    assert.Greater(t, 10, 5)
    assert.GreaterOrEqual(t, 10, 10)
    assert.Less(t, 5, 10)
    assert.LessOrEqual(t, 5, 5)
    
    // Положительное/отрицательное
    assert.Positive(t, 5)
    assert.Negative(t, -5)
    assert.Zero(t, 0)
    assert.NotZero(t, 5)
    
    // Строки
    assert.Contains(t, "hello world", "world")
    assert.NotContains(t, "hello", "world")
    assert.Regexp(t, `^\d{3}-\d{2}-\d{4}$`, "123-45-6789")
    
    // JSON (проверка структурного равенства)
    json1 := `{"name": "John", "age": 30}`
    json2 := `{"age": 30, "name": "John"}`
    assert.JSONEq(t, json1, json2)
}
```

### Пример 9: Eventually и Eventually (async)

```go
package async

import (
    "sync/atomic"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
)

var counter int64

func IncrementAsync() {
    go func() {
        time.Sleep(50 * time.Millisecond)
        atomic.AddInt64(&counter, 1)
    }()
}

func TestEventually(t *testing.T) {
    counter = 0
    IncrementAsync()
    
    // Eventually проверяет условие с повторами
    assert.Eventually(t, func() bool {
        return atomic.LoadInt64(&counter) == 1
    }, time.Second, 10*time.Millisecond)
}

func TestNever(t *testing.T) {
    value := 0
    
    // Never проверяет, что условие никогда не станет true
    assert.Never(t, func() bool {
        return value == 1
    }, 100*time.Millisecond, 10*time.Millisecond)
}
```

---

## ⚠️ Частые ошибки

### 1. Перепутан порядок аргументов

```go
// ❌ НЕПРАВИЛЬНО — порядок: expected, actual
result := Add(2, 3)
assert.Equal(t, result, 5)  // ошибка покажет "expected: result, actual: 5"

// ✅ ПРАВИЛЬНО
assert.Equal(t, 5, result)  // expected: 5, actual: result
```

### 2. Использование assert когда нужен require

```go
// ❌ ПЛОХО — может вызвать panic
func TestUser(t *testing.T) {
    user, err := GetUser(1)
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)  // panic если user == nil
}

// ✅ ХОРОШО
func TestUser(t *testing.T) {
    user, err := GetUser(1)
    require.NoError(t, err)      // остановит тест при ошибке
    require.NotNil(t, user)
    assert.Equal(t, "John", user.Name)  // безопасно
}
```

### 3. Сравнение float без погрешности

```go
// ❌ ПЛОХО — может не сработать из-за floating point
assert.Equal(t, 0.1+0.2, 0.3)

// ✅ ХОРОШО — используем InDelta
assert.InDelta(t, 0.3, 0.1+0.2, 0.0001)

// Или InEpsilon (относительная погрешность)
assert.InEpsilon(t, 0.3, 0.1+0.2, 0.01)
```

### 4. Забыли передать t

```go
// ❌ НЕ СКОМПИЛИРУЕТСЯ
assert.Equal(5, result)

// ✅ ПРАВИЛЬНО
assert.Equal(t, 5, result)
```

---

## 🏋️ Практические задания

### Задание 1: Переписать тесты на testify

Перепишите тесты из предыдущих уроков, используя testify/assert.

**Начальный код:**
```go
// math.go
package math

func Add(a, b int) int { return a + b }
func Sub(a, b int) int { return a - b }
func Mul(a, b int) int { return a * b }
```

```go
// math_test.go
package math

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestMath(t *testing.T) {
    // Используйте assert.Equal, assert.NotEqual
    // Добавьте осмысленные сообщения третьим аргументом
    
    assert.Equal(t, 5, Add(2, 3), "2 + 3 should equal 5")
    // Добавьте тесты для Sub и Mul
}
```

**Ожидаемый вывод:**
```
=== RUN   TestMath
--- PASS: TestMath (0.00s)
PASS
```

**Баллы:** 10

### Задание 2: Тестирование слайсов

Напишите функции работы со слайсами и протестируйте их с testify.

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
// slices_test.go
package slices

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
    input := []int{1, 2, 2, 3, 3, 3}
    result := Unique(input)
    
    assert.Len(t, result, 3)
    assert.ElementsMatch(t, []int{1, 2, 3}, result)
}

func TestContains(t *testing.T) {
    s := []int{1, 2, 3, 4, 5}
    
    assert.True(t, Contains(s, 3))
    assert.False(t, Contains(s, 10))
}

func TestReverse(t *testing.T) {
    input := []int{1, 2, 3}
    result := Reverse(input)
    
    assert.Equal(t, []int{3, 2, 1}, result)
    // Проверьте, что оригинал не изменился
    assert.Equal(t, []int{1, 2, 3}, input, "original should not be modified")
}
```

**Ожидаемый вывод:**
```
=== RUN   TestUnique
--- PASS: TestUnique (0.00s)
=== RUN   TestContains
--- PASS: TestContains (0.00s)
=== RUN   TestReverse
--- PASS: TestReverse (0.00s)
PASS
```

**Баллы:** 15

### Задание 3: require vs assert

Напишите тест, где важен порядок проверок, используя require.

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
// user_test.go
package user

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
    t.Run("valid user", func(t *testing.T) {
        user, err := GetUser(1)
        
        // Используйте require для критичных проверок
        require.NoError(t, err)
        require.NotNil(t, user)
        
        // После require безопасно использовать assert
        assert.Equal(t, 1, user.ID)
        assert.NotEmpty(t, user.Name)
        assert.Contains(t, user.Email, "@")
    })
    
    t.Run("invalid id", func(t *testing.T) {
        user, err := GetUser(-1)
        
        require.Error(t, err)
        assert.Nil(t, user)
        assert.Contains(t, err.Error(), "invalid")
    })
}
```

**Ожидаемый вывод:**
```
=== RUN   TestGetUser
=== RUN   TestGetUser/valid_user
=== RUN   TestGetUser/invalid_id
--- PASS: TestGetUser (0.00s)
PASS
```

**Баллы:** 15

### Задание 4: Тест с Eventually

Протестируйте асинхронную функцию с помощью Eventually.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Создайте функцию согласно заданию

func main() {
    // Ваш код здесь
    
}
```

```go
// async_test.go
package async

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

func TestCounterAsync(t *testing.T) {
    c := &Counter{}
    
    c.IncrementAsync()
    c.IncrementAsync()
    c.IncrementAsync()
    
    // Используйте Eventually для ожидания результата
    assert.Eventually(t, func() bool {
        return c.Value() == 3
    }, time.Second, 10*time.Millisecond, "counter should reach 3")
}
```

**Ожидаемый вывод:**
```
=== RUN   TestCounterAsync
--- PASS: TestCounterAsync (0.35s)
PASS
```

**Баллы:** 10

### Задание 5: Тест ошибок с ErrorIs

Создайте иерархию ошибок и протестируйте их.

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
// errors_test.go
package errors

import (
    "errors"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestGetResource(t *testing.T) {
    t.Run("unauthorized", func(t *testing.T) {
        _, err := GetResource(1, "")
        assert.ErrorIs(t, err, ErrUnauthorized)
    })
    
    t.Run("validation", func(t *testing.T) {
        _, err := GetResource(-1, "token")
        assert.ErrorIs(t, err, ErrValidation)
        assert.ErrorContains(t, err, "invalid id")
    })
    
    t.Run("not found", func(t *testing.T) {
        _, err := GetResource(999, "token")
        assert.ErrorIs(t, err, ErrNotFound)
    })
    
    t.Run("success", func(t *testing.T) {
        result, err := GetResource(1, "token")
        assert.NoError(t, err)
        assert.Equal(t, "resource-1", result)
    })
}
```

**Ожидаемый вывод:**
```
=== RUN   TestGetResource
=== RUN   TestGetResource/unauthorized
=== RUN   TestGetResource/validation
=== RUN   TestGetResource/not_found
=== RUN   TestGetResource/success
--- PASS: TestGetResource (0.00s)
PASS
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Testify GitHub](https://github.com/stretchr/testify)
- [Testify Assert Package](https://pkg.go.dev/github.com/stretchr/testify/assert)
- [Testify Require Package](https://pkg.go.dev/github.com/stretchr/testify/require)
