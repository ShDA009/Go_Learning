# Введение в Generics (Go 1.18+)

---

## 💡 Ключевые идеи

1. **Generics** — параметризованные типы, позволяющие писать универсальный код
2. **Type Parameters** — параметры типа в квадратных скобках `[T any]`
3. **Constraints** — ограничения на типы (интерфейсы)
4. **any** — псевдоним для `interface{}`, любой тип
5. **comparable** — встроенное ограничение для сравниваемых типов
6. **Type Inference** — автоматический вывод типов

### До и после Generics

```go
// ДО: дублирование кода
func MaxInt(a, b int) int { ... }
func MaxFloat(a, b float64) float64 { ... }
func MaxString(a, b string) string { ... }

// ПОСЛЕ: один универсальный код
func Max[T cmp.Ordered](a, b T) T { ... }
```

---

## 📖 Теория

### Что такое Generics?

**Generics** (обобщённое программирование) — это способ писать код, который работает с **разными типами** без дублирования. Вместо того чтобы писать отдельную функцию для каждого типа, вы пишете одну универсальную функцию.

### Проблема до Generics

До Go 1.18 (февраль 2022), если вам нужна была функция поиска максимума, приходилось писать несколько версий:

```go
func MaxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func MaxFloat64(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
}

func MaxString(a, b string) string {
    if a > b {
        return a
    }
    return b
}
```

Это называется **дублирование кода** — одна и та же логика повторяется трижды. Если нужно изменить алгоритм — менять придётся везде.

### Альтернатива: interface{}

До Generics можно было использовать `interface{}` (пустой интерфейс):

```go
func Max(a, b interface{}) interface{} {
    // но как сравнить a > b?
    // нужны type assertions и switch
}
```

Проблемы этого подхода:
- **Нет проверки типов** на этапе компиляции
- **Потеря информации о типе** — результат тоже `interface{}`
- **Runtime ошибки** — паника при неверном приведении типов
- **Неудобно использовать** — нужно приводить типы вручную

### Решение: Generics

С Generics вы пишете одну функцию с **параметром типа**:

```go
func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// Использование
fmt.Println(Max(1, 2))        // int, результат: 2
fmt.Println(Max(1.5, 2.5))    // float64, результат: 2.5
fmt.Println(Max("a", "b"))    // string, результат: "b"
```

### Как читать синтаксис Generics?

```go
func Max[T cmp.Ordered](a, b T) T
       ↑  ↑             ↑      ↑
       │  │             │      └─ возвращаемый тип (тот же T)
       │  │             └──────── параметры функции типа T
       │  └────────────────────── constraint (ограничение): какие типы допустимы
       └───────────────────────── T — параметр типа (можно назвать как угодно)
```

**Квадратные скобки `[T ...]`** — это место, где объявляются параметры типа. Это отличает их от обычных параметров в круглых скобках.

### Type Inference (вывод типов)

Go умеет автоматически определять тип по аргументам:

```go
Max(1, 2)         // Go понимает: T = int
Max[int](1, 2)    // явное указание типа (не обязательно)
```

Вывод типов делает код чище — не нужно указывать типы, если они очевидны.

### Когда использовать Generics?

✅ **Используйте когда:**
- Функция работает одинаково для разных типов (Min, Max, Sort, Contains)
- Структура данных может хранить разные типы (Stack, Queue, Set)
- Хотите избежать дублирования кода

❌ **Не используйте когда:**
- Функция специфична для одного типа
- Логика различается для разных типов (лучше интерфейсы)
- Код становится сложнее для понимания

### Generics vs Interfaces

| Generics | Interfaces |
|----------|------------|
| Типы известны на этапе компиляции | Типы определяются в runtime |
| Более производительны | Небольшой overhead |
| Для работы с **данными** | Для работы с **поведением** |

---

## 📋 Синтаксис

### Обобщённая функция

```go
func FunctionName[T constraint](param T) T {
    return param
}
```

### Обобщённый тип

```go
type TypeName[T any] struct {
    value T
}
```

### Вызов с явным типом

```go
result := FunctionName[int](42)
```

### Вызов с выводом типа

```go
result := FunctionName(42)  // T = int выводится автоматически
```

### Несколько параметров типа

```go
func Pair[K, V any](key K, value V) (K, V) {
    return key, value
}
```

---

## 💻 Примеры кода

### Пример 1: Простая обобщённая функция

```go
package main

import "fmt"

// Обобщённая функция идентичности
func Identity[T any](value T) T {
    return value
}

func main() {
    // Явное указание типа
    intResult := Identity[int](42)
    fmt.Println(intResult)  // 42
    
    // Вывод типа (type inference)
    strResult := Identity("Hello")
    fmt.Println(strResult)  // Hello
    
    floatResult := Identity(3.14)
    fmt.Println(floatResult)  // 3.14
}
```

### Пример 2: Функция поиска минимума/максимума

```go
package main

import (
    "cmp"
    "fmt"
)

// Используем constraint cmp.Ordered
func Min[T cmp.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

func main() {
    fmt.Println(Min(5, 3))          // 3
    fmt.Println(Min(5.5, 3.3))      // 3.3
    fmt.Println(Min("apple", "banana"))  // apple
    
    fmt.Println(Max(10, 20))        // 20
    fmt.Println(Max("z", "a"))      // z
}
```

### Пример 3: Обобщённая работа со срезами

```go
package main

import "fmt"

// Фильтрация среза
func Filter[T any](slice []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Преобразование среза
func Map[T, R any](slice []T, transform func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = transform(v)
    }
    return result
}

// Свёртка среза
func Reduce[T, R any](slice []T, initial R, reducer func(R, T) R) R {
    result := initial
    for _, v := range slice {
        result = reducer(result, v)
    }
    return result
}

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Фильтруем чётные
    evens := Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println("Evens:", evens)  // [2 4 6 8 10]
    
    // Удваиваем
    doubled := Map(numbers, func(n int) int {
        return n * 2
    })
    fmt.Println("Doubled:", doubled)  // [2 4 6 8 10 12 14 16 18 20]
    
    // Конвертируем в строки
    strings := Map(numbers, func(n int) string {
        return fmt.Sprintf("#%d", n)
    })
    fmt.Println("Strings:", strings)  // [#1 #2 #3 ...]
    
    // Сумма
    sum := Reduce(numbers, 0, func(acc, n int) int {
        return acc + n
    })
    fmt.Println("Sum:", sum)  // 55
}
```

### Пример 4: Поиск в срезе

```go
package main

import "fmt"

// Содержит ли срез элемент (требуется comparable)
func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {
            return true
        }
    }
    return false
}

// Индекс элемента
func IndexOf[T comparable](slice []T, target T) int {
    for i, v := range slice {
        if v == target {
            return i
        }
    }
    return -1
}

// Удаление дубликатов
func Unique[T comparable](slice []T) []T {
    seen := make(map[T]bool)
    result := make([]T, 0)
    
    for _, v := range slice {
        if !seen[v] {
            seen[v] = true
            result = append(result, v)
        }
    }
    return result
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    fmt.Println(Contains(numbers, 3))   // true
    fmt.Println(Contains(numbers, 10))  // false
    
    fmt.Println(IndexOf(numbers, 4))    // 3
    fmt.Println(IndexOf(numbers, 10))   // -1
    
    words := []string{"a", "b", "a", "c", "b", "d"}
    fmt.Println(Unique(words))  // [a b c d]
}
```

### Пример 5: Обобщённый Stack

```go
package main

import (
    "errors"
    "fmt"
)

// Обобщённый стек
type Stack[T any] struct {
    items []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, error) {
    var zero T
    if len(s.items) == 0 {
        return zero, errors.New("stack is empty")
    }
    
    index := len(s.items) - 1
    item := s.items[index]
    s.items = s.items[:index]
    return item, nil
}

func (s *Stack[T]) Peek() (T, error) {
    var zero T
    if len(s.items) == 0 {
        return zero, errors.New("stack is empty")
    }
    return s.items[len(s.items)-1], nil
}

func (s *Stack[T]) IsEmpty() bool {
    return len(s.items) == 0
}

func (s *Stack[T]) Size() int {
    return len(s.items)
}

func main() {
    // Стек целых чисел
    intStack := NewStack[int]()
    intStack.Push(1)
    intStack.Push(2)
    intStack.Push(3)
    
    fmt.Println("Size:", intStack.Size())  // 3
    
    for !intStack.IsEmpty() {
        val, _ := intStack.Pop()
        fmt.Println(val)  // 3, 2, 1
    }
    
    // Стек строк
    strStack := NewStack[string]()
    strStack.Push("Hello")
    strStack.Push("World")
    
    top, _ := strStack.Peek()
    fmt.Println("Top:", top)  // World
}
```

### Пример 6: Обобщённая Map (словарь)

```go
package main

import "fmt"

// Обобщённая карта с методами
type Dictionary[K comparable, V any] struct {
    data map[K]V
}

func NewDictionary[K comparable, V any]() *Dictionary[K, V] {
    return &Dictionary[K, V]{data: make(map[K]V)}
}

func (d *Dictionary[K, V]) Set(key K, value V) {
    d.data[key] = value
}

func (d *Dictionary[K, V]) Get(key K) (V, bool) {
    value, ok := d.data[key]
    return value, ok
}

func (d *Dictionary[K, V]) GetOrDefault(key K, defaultValue V) V {
    if value, ok := d.data[key]; ok {
        return value
    }
    return defaultValue
}

func (d *Dictionary[K, V]) Delete(key K) {
    delete(d.data, key)
}

func (d *Dictionary[K, V]) Keys() []K {
    keys := make([]K, 0, len(d.data))
    for k := range d.data {
        keys = append(keys, k)
    }
    return keys
}

func (d *Dictionary[K, V]) Values() []V {
    values := make([]V, 0, len(d.data))
    for _, v := range d.data {
        values = append(values, v)
    }
    return values
}

func main() {
    // string -> int
    scores := NewDictionary[string, int]()
    scores.Set("Alice", 100)
    scores.Set("Bob", 85)
    scores.Set("Charlie", 92)
    
    if score, ok := scores.Get("Alice"); ok {
        fmt.Println("Alice's score:", score)  // 100
    }
    
    fmt.Println("Dave's score:", scores.GetOrDefault("Dave", 0))  // 0
    
    fmt.Println("Keys:", scores.Keys())  // [Alice Bob Charlie]
    
    // int -> string
    names := NewDictionary[int, string]()
    names.Set(1, "One")
    names.Set(2, "Two")
    
    fmt.Println("Values:", names.Values())  // [One Two]
}
```

### Пример 7: Несколько параметров типа

```go
package main

import "fmt"

// Пара значений разных типов
type Pair[K, V any] struct {
    Key   K
    Value V
}

func NewPair[K, V any](key K, value V) Pair[K, V] {
    return Pair[K, V]{Key: key, Value: value}
}

// Функция обмена
func Swap[T any](a, b *T) {
    *a, *b = *b, *a
}

// Zip — объединение двух срезов
func Zip[T, U any](a []T, b []U) []Pair[T, U] {
    minLen := len(a)
    if len(b) < minLen {
        minLen = len(b)
    }
    
    result := make([]Pair[T, U], minLen)
    for i := 0; i < minLen; i++ {
        result[i] = NewPair(a[i], b[i])
    }
    return result
}

func main() {
    // Пара
    p := NewPair("name", 42)
    fmt.Printf("Key: %v, Value: %v\n", p.Key, p.Value)
    
    // Swap
    x, y := 10, 20
    Swap(&x, &y)
    fmt.Printf("x=%d, y=%d\n", x, y)  // x=20, y=10
    
    // Zip
    names := []string{"Alice", "Bob", "Charlie"}
    ages := []int{25, 30, 35}
    
    pairs := Zip(names, ages)
    for _, p := range pairs {
        fmt.Printf("%s is %d years old\n", p.Key, p.Value)
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Забыли constraint

```go
// ❌ ОШИБКА — нельзя сравнивать any
func Max[T any](a, b T) T {
    if a > b {  // invalid operation: a > b
        return a
    }
    return b
}

// ✅ ПРАВИЛЬНО — используем constraint
func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

### 2. Нельзя использовать == с any

```go
// ❌ ОШИБКА
func Contains[T any](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {  // invalid operation
            return true
        }
    }
    return false
}

// ✅ ПРАВИЛЬНО — используем comparable
func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {
            return true
        }
    }
    return false
}
```

### 3. Нельзя использовать оператор приведения типа напрямую

```go
// ❌ ОШИБКА
func Convert[T any](value interface{}) T {
    return value.(T)  // может panic
}

// ✅ ЛУЧШЕ — с проверкой
func Convert[T any](value interface{}) (T, bool) {
    result, ok := value.(T)
    return result, ok
}
```

### 4. Zero value

```go
// ❌ ПЛОХО — nil не работает для всех типов
func GetOrNil[T any](slice []T, index int) T {
    if index < 0 || index >= len(slice) {
        return nil  // cannot use nil as T value
    }
    return slice[index]
}

// ✅ ПРАВИЛЬНО — zero value
func GetOrZero[T any](slice []T, index int) T {
    var zero T
    if index < 0 || index >= len(slice) {
        return zero
    }
    return slice[index]
}
```

---

## 🏋️ Практические задания

### Задание 1: Generic Reverse

Напишите обобщённую функцию для разворота среза.

**Начальный код:**
```go
package main

import "fmt"

func Reverse[T any](slice []T) []T {
    // Реализуйте функцию
    // Не изменяйте оригинальный срез!
}

func main() {
    ints := []int{1, 2, 3, 4, 5}
    strings := []string{"a", "b", "c"}
    
    fmt.Println(Reverse(ints))    // [5 4 3 2 1]
    fmt.Println(Reverse(strings)) // [c b a]
    fmt.Println(ints)             // [1 2 3 4 5] — не изменился
}
```

**Ожидаемый вывод:**
```
[5 4 3 2 1]
[c b a]
[1 2 3 4 5]
```

**Баллы:** 10

### Задание 2: Generic Max

Напишите функцию поиска максимума с использованием cmp.Ordered.

**Начальный код:**
```go
package main

import (
    "cmp"
    "fmt"
)

func Max[T cmp.Ordered](slice []T) (T, bool) {
    // Верните максимальный элемент и true
    // Если срез пустой — zero value и false
}

func main() {
    fmt.Println(Max([]int{3, 1, 4, 1, 5}))       // 5 true
    fmt.Println(Max([]string{"cat", "dog", "apple"})) // dog true
    fmt.Println(Max([]float64{}))               // 0 false
}
```

**Ожидаемый вывод:**
```
5 true
dog true
0 false
```

**Баллы:** 10

### Задание 3: Generic Map функция

Реализуйте функцию Map для трансформации среза.

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
[2 4 6 8 10]
[#1 #2 #3 #4 #5]
```

**Баллы:** 10

### Задание 4: Generic Filter

Реализуйте функцию фильтрации среза.

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
[2 4 6 8 10]
[6 7 8 9 10]
[apple banana grape]
```

**Баллы:** 10

### Задание 5: Generic Stack

Создайте обобщённую структуру стека.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Определите структуру

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
Size: 3
Peek: 3
Pop: 3
Pop: 2
Pop: 1
Empty pop ok: false
```

**Баллы:** 15

---

## 🔗 Полезные ссылки

- [Go Generics Tutorial](https://go.dev/doc/tutorial/generics)
- [Type Parameters Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)
- [cmp Package](https://pkg.go.dev/cmp)
