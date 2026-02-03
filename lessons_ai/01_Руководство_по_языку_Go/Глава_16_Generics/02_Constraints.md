# Constraints (Ограничения типов)

---

## 💡 Ключевые идеи

1. **Constraint** — интерфейс, ограничивающий допустимые типы
2. **any** — любой тип (псевдоним `interface{}`)
3. **comparable** — типы, поддерживающие == и !=
4. **cmp.Ordered** — типы с операторами <, >, <=, >=
5. **Union types** — объединение типов через `|`
6. **~T** — приближённый тип (включая пользовательские на базе T)

### Встроенные ограничения

| Constraint | Описание |
|------------|----------|
| `any` | Любой тип |
| `comparable` | == и != |
| `cmp.Ordered` | Числа и строки |

---

## 📖 Теория

### Что такое Constraint?

**Constraint** (ограничение) — это правило, которое определяет, какие типы можно использовать с generic-функцией. Без ограничений Go не знает, что можно делать с типом.

### Зачем нужны ограничения?

Представьте функцию:
```go
func Max[T any](a, b T) T {
    if a > b {  // ОШИБКА: не все типы поддерживают >
        return a
    }
    return b
}
```

Тип `any` означает "любой тип", но не все типы можно сравнивать с помощью `>`. Например, нельзя написать `struct{} > struct{}`.

Чтобы использовать `>`, нужно **ограничить** типы теми, которые это поддерживают:

```go
func Max[T cmp.Ordered](a, b T) T {  // только упорядоченные типы
    if a > b {  // OK!
        return a
    }
    return b
}
```

### Встроенные ограничения

**1. any — любой тип**
```go
func Print[T any](value T) {
    fmt.Println(value)  // Println принимает любой тип
}
```

**2. comparable — сравниваемые типы (== и !=)**
```go
func Contains[T comparable](slice []T, item T) bool {
    for _, v := range slice {
        if v == item {  // нужен comparable для ==
            return true
        }
    }
    return false
}
```

Comparable включает: числа, строки, булевы, указатели, каналы, структуры (если все поля comparable).
НЕ включает: срезы, map, функции.

**3. cmp.Ordered — упорядоченные типы (<, >, <=, >=)**
```go
import "cmp"

func Sort[T cmp.Ordered](slice []T) {
    // можно использовать <, >, <=, >=
}
```

Ordered включает: все числа (int, float...), строки.

### Создание собственных ограничений

Constraint — это просто интерфейс. Но для generics в нём можно использовать специальный синтаксис:

**Union types (объединение типов):**
```go
type Number interface {
    int | int64 | float64
}

func Sum[T Number](values []T) T {
    var sum T
    for _, v := range values {
        sum += v
    }
    return sum
}
```

Теперь `Sum` принимает только int, int64 или float64.

### Оператор ~ (приближённый тип)

Проблема: пользовательские типы на базе int не подходят под `int`:

```go
type MyInt int
var x MyInt = 5

Sum([]MyInt{1, 2, 3})  // ОШИБКА: MyInt != int
```

Решение — оператор `~` означает "этот тип и все типы на его основе":

```go
type Number interface {
    ~int | ~int64 | ~float64
}

Sum([]MyInt{1, 2, 3})  // OK! MyInt основан на int
```

### Методы в constraint

Constraint может требовать наличие методов:

```go
type Stringer interface {
    String() string
}

func PrintAll[T Stringer](items []T) {
    for _, item := range items {
        fmt.Println(item.String())
    }
}
```

### Комбинирование ограничений

Можно объединять типы и методы:

```go
type StringableInt interface {
    ~int | ~int64
    String() string
}
```

Это означает: тип должен быть основан на int/int64 И иметь метод String().

### Правило: минимальные ограничения

Используйте самое слабое ограничение, которое работает:
- `any` — если не нужны никакие операции над типом
- `comparable` — если нужно сравнивать на равенство
- `cmp.Ordered` — если нужно сравнивать по порядку

---

## 📋 Синтаксис

### Использование constraint

```go
func Function[T constraint](param T) T { ... }
```

### Определение собственного constraint

```go
type MyConstraint interface {
    Method() string
}
```

### Union types

```go
type Number interface {
    int | int64 | float64
}
```

### Приближённый тип (~)

```go
type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}
```

---

## 💻 Примеры кода

### Пример 1: Встроенные constraints

```go
package main

import (
    "cmp"
    "fmt"
)

// any — любой тип
func Print[T any](value T) {
    fmt.Printf("Value: %v (type: %T)\n", value, value)
}

// comparable — можно сравнивать
func Equal[T comparable](a, b T) bool {
    return a == b
}

// cmp.Ordered — можно сортировать
func Sort[T cmp.Ordered](slice []T) {
    for i := 0; i < len(slice)-1; i++ {
        for j := i + 1; j < len(slice); j++ {
            if slice[j] < slice[i] {
                slice[i], slice[j] = slice[j], slice[i]
            }
        }
    }
}

func main() {
    // any
    Print(42)
    Print("hello")
    Print([]int{1, 2, 3})
    
    // comparable
    fmt.Println(Equal(5, 5))       // true
    fmt.Println(Equal("a", "b"))   // false
    
    // cmp.Ordered
    nums := []int{5, 2, 8, 1, 9}
    Sort(nums)
    fmt.Println(nums)  // [1 2 5 8 9]
    
    words := []string{"banana", "apple", "cherry"}
    Sort(words)
    fmt.Println(words)  // [apple banana cherry]
}
```

### Пример 2: Собственный constraint с методом

```go
package main

import "fmt"

// Constraint — тип должен иметь метод String()
type Stringer interface {
    String() string
}

// Обобщённая функция, требующая Stringer
func PrintAll[T Stringer](items []T) {
    for _, item := range items {
        fmt.Println(item.String())
    }
}

// Реализация
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}

type Product struct {
    Name  string
    Price float64
}

func (p Product) String() string {
    return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
}

func main() {
    people := []Person{
        {Name: "Alice", Age: 30},
        {Name: "Bob", Age: 25},
    }
    PrintAll(people)
    
    products := []Product{
        {Name: "Phone", Price: 999.99},
        {Name: "Laptop", Price: 1499.99},
    }
    PrintAll(products)
}
```

### Пример 3: Union types (объединение типов)

```go
package main

import "fmt"

// Constraint: только числовые типы
type Number interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64 |
    float32 | float64
}

func Sum[T Number](numbers []T) T {
    var sum T
    for _, n := range numbers {
        sum += n
    }
    return sum
}

func Average[T Number](numbers []T) float64 {
    if len(numbers) == 0 {
        return 0
    }
    sum := Sum(numbers)
    return float64(sum) / float64(len(numbers))
}

func main() {
    ints := []int{1, 2, 3, 4, 5}
    fmt.Println("Sum:", Sum(ints))        // 15
    fmt.Println("Avg:", Average(ints))    // 3
    
    floats := []float64{1.5, 2.5, 3.5}
    fmt.Println("Sum:", Sum(floats))      // 7.5
    fmt.Println("Avg:", Average(floats))  // 2.5
    
    bytes := []uint8{10, 20, 30}
    fmt.Println("Sum:", Sum(bytes))       // 60
}
```

### Пример 4: Приближённый тип (~)

```go
package main

import "fmt"

// Без ~ — только точные типы
type ExactInt interface {
    int | int64
}

// С ~ — включая пользовательские типы на основе int
type ApproxInt interface {
    ~int | ~int64
}

// Пользовательский тип
type UserID int
type OrderID int64

func DoubleExact[T ExactInt](v T) T {
    return v * 2
}

func DoubleApprox[T ApproxInt](v T) T {
    return v * 2
}

func main() {
    var x int = 5
    fmt.Println(DoubleExact(x))   // OK: 10
    fmt.Println(DoubleApprox(x))  // OK: 10
    
    var userID UserID = 100
    // fmt.Println(DoubleExact(userID))  // ❌ ОШИБКА: UserID не int
    fmt.Println(DoubleApprox(userID))    // ✅ OK: 200 (UserID основан на int)
    
    var orderID OrderID = 1000
    fmt.Println(DoubleApprox(orderID))   // ✅ OK: 2000
}
```

### Пример 5: Комбинирование constraints

```go
package main

import "fmt"

// Constraint с несколькими требованиями
type OrderedStringer interface {
    ~int | ~string
    String() string  // должен иметь метод String()
}

// ИЛИ использование вложения интерфейсов
type Ordered interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
    ~float32 | ~float64 | ~string
}

type Printable interface {
    String() string
}

// Тип должен быть и Ordered, и Printable
type OrderedPrintable interface {
    Ordered
    Printable
}

// Использование в методе
type Score int

func (s Score) String() string {
    return fmt.Sprintf("Score: %d", s)
}

func PrintSorted[T OrderedPrintable](items []T) {
    // Сортируем
    for i := 0; i < len(items)-1; i++ {
        for j := i + 1; j < len(items); j++ {
            if items[j] < items[i] {
                items[i], items[j] = items[j], items[i]
            }
        }
    }
    
    // Печатаем
    for _, item := range items {
        fmt.Println(item.String())
    }
}

func main() {
    scores := []Score{95, 87, 92, 78}
    PrintSorted(scores)
}
```

### Пример 6: Constraint для map ключей

```go
package main

import "fmt"

// Map ключ должен быть comparable
type MapKey interface {
    comparable
}

// Обобщённый кэш
type Cache[K comparable, V any] struct {
    data map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{
        data: make(map[K]V),
    }
}

func (c *Cache[K, V]) Set(key K, value V) {
    c.data[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    v, ok := c.data[key]
    return v, ok
}

// Функция подсчёта вхождений
func CountOccurrences[T comparable](items []T) map[T]int {
    counts := make(map[T]int)
    for _, item := range items {
        counts[item]++
    }
    return counts
}

func main() {
    // Кэш с string ключами
    strCache := NewCache[string, int]()
    strCache.Set("one", 1)
    strCache.Set("two", 2)
    
    if v, ok := strCache.Get("one"); ok {
        fmt.Println("one =", v)
    }
    
    // Кэш с int ключами
    intCache := NewCache[int, string]()
    intCache.Set(1, "one")
    intCache.Set(2, "two")
    
    // Подсчёт вхождений
    words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
    counts := CountOccurrences(words)
    fmt.Println(counts)  // map[apple:3 banana:2 cherry:1]
}
```

### Пример 7: Signed и Unsigned constraints

```go
package main

import (
    "fmt"
    "golang.org/x/exp/constraints"
)

// Из пакета constraints (golang.org/x/exp/constraints)
// type Signed interface {
//     ~int | ~int8 | ~int16 | ~int32 | ~int64
// }
// type Unsigned interface {
//     ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
// }
// type Integer interface {
//     Signed | Unsigned
// }
// type Float interface {
//     ~float32 | ~float64
// }

// Абсолютное значение для знаковых чисел
func Abs[T constraints.Signed](n T) T {
    if n < 0 {
        return -n
    }
    return n
}

// Проверка на чётность для целых
func IsEven[T constraints.Integer](n T) bool {
    return n%2 == 0
}

// Округление для float
func Round[T constraints.Float](n T) T {
    if n < 0 {
        return T(int(n - 0.5))
    }
    return T(int(n + 0.5))
}

func main() {
    fmt.Println(Abs(-5))     // 5
    fmt.Println(Abs(int8(-10)))  // 10
    
    fmt.Println(IsEven(4))   // true
    fmt.Println(IsEven(uint(7))) // false
    
    fmt.Println(Round(3.7))  // 4
    fmt.Println(Round(-2.3)) // -2
}
```

### Пример 8: Собственные сложные constraints

```go
package main

import "fmt"

// Constraint для сериализуемых типов
type Serializable interface {
    Serialize() []byte
    Deserialize([]byte) error
}

// Constraint для клонируемых типов
type Cloneable[T any] interface {
    Clone() T
}

// Пример реализации
type Config struct {
    Name  string
    Value int
}

func (c Config) Clone() Config {
    return Config{
        Name:  c.Name,
        Value: c.Value,
    }
}

// Обобщённая функция клонирования
func CloneAll[T Cloneable[T]](items []T) []T {
    result := make([]T, len(items))
    for i, item := range items {
        result[i] = item.Clone()
    }
    return result
}

func main() {
    configs := []Config{
        {Name: "debug", Value: 1},
        {Name: "timeout", Value: 30},
    }
    
    cloned := CloneAll(configs)
    
    // Изменяем оригинал
    configs[0].Value = 999
    
    // Клон не изменился
    fmt.Println("Original:", configs[0])  // {debug 999}
    fmt.Println("Cloned:", cloned[0])     // {debug 1}
}
```

---

## ⚠️ Частые ошибки

### 1. Слишком строгий constraint

```go
// ❌ СЛИШКОМ СТРОГО — только int
type OnlyInt interface {
    int
}

// ✅ ЛУЧШЕ — все целые числа
type Integer interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}
```

### 2. Забыли ~ для пользовательских типов

```go
type MyInt int

// ❌ НЕ РАБОТАЕТ с MyInt
type Strict interface { int }

// ✅ РАБОТАЕТ с MyInt
type Flexible interface { ~int }
```

### 3. Constraint в неправильном месте

```go
// ❌ НЕПРАВИЛЬНО — constraint должен быть интерфейсом
func Bad[T int | string](v T) {}

// ✅ ПРАВИЛЬНО — определяем интерфейс
type IntOrString interface { int | string }
func Good[T IntOrString](v T) {}
```

---

## 🏋️ Практические задания

### Задание 1: Numeric constraint

Создайте constraint для всех числовых типов и функцию суммирования.

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
15
6.6
60
```

**Баллы:** 10

### Задание 2: Comparable + Stringer

Создайте constraint для типов, которые можно сравнивать и преобразовывать в строку.

**Начальный код:**
```go
package main

import (
    "fmt"
    "strings"
)

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Found: Bob (25)
Not found
```

**Баллы:** 10

### Задание 3: Entity constraint

Создайте интерфейс Entity для работы с сущностями БД.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Определите и реализуйте интерфейс

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
User 1: &{ID:1 Name:Alice}
User 2: &{ID:2 Name:Bob}
Found: &{ID:1 Name:Alice}
```

**Баллы:** 15

### Задание 4: Validator constraint

Создайте систему валидации с обобщёнными типами.

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
item 1: invalid email: missing @
item 0: password too short
```

**Баллы:** 15

### Задание 5: Ordered с custom types

Используйте ~ для поддержки пользовательских типов.

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
10
3.14
36.6
200
100
0
50
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Go Constraints Package](https://pkg.go.dev/golang.org/x/exp/constraints)
- [Type Parameters Design](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)
- [Generics Tutorial](https://go.dev/doc/tutorial/generics)
