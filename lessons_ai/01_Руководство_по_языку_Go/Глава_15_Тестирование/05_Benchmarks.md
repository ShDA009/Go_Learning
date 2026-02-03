# Benchmarks (бенчмарки)

---

## 💡 Ключевые идеи

1. **Benchmark** — измерение производительности кода
2. **Функция Benchmark*** — соглашение именования для бенчмарков
3. **b.N** — количество итераций, определяемое фреймворком
4. **go test -bench** — команда для запуска бенчмарков
5. **Memory allocation** — отслеживание выделения памяти
6. **Сравнение** — benchstat для сравнения результатов

### Зачем нужны бенчмарки?

| Ситуация | Решение |
|----------|---------|
| Какой алгоритм быстрее? | Написать бенчмарк |
| Где bottleneck? | Профилирование + бенчмарк |
| Регрессия производительности | CI с бенчмарками |
| Оптимизация помогла? | До/после бенчмарк |

---

## 📖 Теория

### Что такое бенчмарк?

**Бенчмарк** — это тест, который измеряет, насколько быстро работает ваш код. Если обычные тесты отвечают на вопрос "Работает ли код правильно?", то бенчмарки отвечают на вопрос "Насколько быстро работает код?".

### Зачем измерять производительность?

1. **Сравнение алгоритмов** — какой способ быстрее?
2. **Поиск узких мест** — где код тормозит?
3. **Контроль регрессий** — не стал ли код медленнее после изменений?
4. **Оптимизация** — помогла ли оптимизация?

### Как работают бенчмарки в Go?

Go сам определяет, сколько раз нужно выполнить код, чтобы получить точные измерения. Это число хранится в `b.N`:

```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(1, 2)  // выполняется b.N раз
    }
}
```

Go начинает с небольшого `b.N` и увеличивает его, пока результаты не станут стабильными (обычно 1 секунда измерений).

### Понимание результатов бенчмарка

```
BenchmarkAdd-8     1000000000     0.318 ns/op     0 B/op     0 allocs/op
```

Разберём каждое число:
- **BenchmarkAdd-8** — название и количество CPU (GOMAXPROCS)
- **1000000000** — сколько раз выполнился код
- **0.318 ns/op** — время на одну операцию (наносекунды)
- **0 B/op** — байт памяти на операцию
- **0 allocs/op** — аллокаций памяти на операцию

### Что такое ns/op?

**ns** — наносекунда, одна миллиардная секунды. Для понимания масштаба:
- 1 ns — доступ к кэшу CPU
- 100 ns — доступ к оперативной памяти
- 10 000 ns (10 μs) — простой системный вызов
- 1 000 000 ns (1 ms) — сетевой запрос в локальной сети

### Почему важны аллокации?

Выделение памяти (`allocs/op`) — одна из главных причин медленного кода в Go. Каждая аллокация:
1. Требует времени на выделение
2. Создаёт работу для сборщика мусора (GC)
3. Может вызвать паузы в программе

**Правило:** Меньше аллокаций = быстрее код.

### Типичные ошибки в бенчмарках

1. **Не сбросить таймер после setup:**
```go
func BenchmarkProcess(b *testing.B) {
    data := loadTestData()  // долгая операция
    b.ResetTimer()          // ВАЖНО: сбросить таймер
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}
```

2. **Компилятор оптимизирует "бесполезный" код:**
```go
// ❌ Плохо — компилятор может удалить вызов
for i := 0; i < b.N; i++ {
    Add(1, 2)  // результат не используется
}

// ✅ Хорошо — сохраняем результат
var result int
for i := 0; i < b.N; i++ {
    result = Add(1, 2)
}
_ = result
```

### Когда запускать бенчмарки?

- После оптимизаций — проверить, что стало лучше
- При выборе алгоритма — сравнить варианты
- В CI/CD — отслеживать регрессии производительности

---

## 📋 Синтаксис

### Структура бенчмарка

```go
func BenchmarkFunctionName(b *testing.B) {
    // Setup (не учитывается в измерениях)
    
    for i := 0; i < b.N; i++ {
        // Измеряемый код
        Function()
    }
}
```

### Основные методы *testing.B

```go
b.N                  // количество итераций
b.ResetTimer()       // сбросить таймер (после setup)
b.StopTimer()        // приостановить измерение
b.StartTimer()       // возобновить измерение
b.ReportAllocs()     // включить отчёт по памяти
b.SetBytes(n)        // установить размер обрабатываемых данных
b.RunParallel()      // параллельный бенчмарк
```

### Команды запуска

```bash
go test -bench=.                    # все бенчмарки
go test -bench=BenchmarkName        # конкретный бенчмарк
go test -bench=. -benchmem          # с отчётом по памяти
go test -bench=. -benchtime=5s      # длительность теста
go test -bench=. -count=5           # количество прогонов
go test -bench=. -cpu=1,2,4         # разное количество CPU
```

---

## 💻 Примеры кода

### Пример 1: Простой бенчмарк

```go
package math

import "testing"

func Sum(numbers []int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

func BenchmarkSum(b *testing.B) {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    for i := 0; i < b.N; i++ {
        Sum(numbers)
    }
}

// Запуск: go test -bench=BenchmarkSum
// Результат:
// BenchmarkSum-8    100000000    10.5 ns/op
```

### Пример 2: Бенчмарк с setup

```go
package search

import (
    "math/rand"
    "testing"
)

func LinearSearch(slice []int, target int) int {
    for i, v := range slice {
        if v == target {
            return i
        }
    }
    return -1
}

func BenchmarkLinearSearch(b *testing.B) {
    // Setup — не включается в измерение
    slice := make([]int, 10000)
    for i := range slice {
        slice[i] = rand.Intn(10000)
    }
    target := slice[len(slice)/2]  // ищем элемент в середине
    
    b.ResetTimer()  // сбрасываем таймер после setup
    
    for i := 0; i < b.N; i++ {
        LinearSearch(slice, target)
    }
}
```

### Пример 3: Сравнение алгоритмов

```go
package concat

import (
    "bytes"
    "strings"
    "testing"
)

// Конкатенация через +
func ConcatPlus(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s
    }
    return result
}

// Конкатенация через strings.Builder
func ConcatBuilder(strs []string) string {
    var builder strings.Builder
    for _, s := range strs {
        builder.WriteString(s)
    }
    return builder.String()
}

// Конкатенация через bytes.Buffer
func ConcatBuffer(strs []string) string {
    var buffer bytes.Buffer
    for _, s := range strs {
        buffer.WriteString(s)
    }
    return buffer.String()
}

// Конкатенация через strings.Join
func ConcatJoin(strs []string) string {
    return strings.Join(strs, "")
}

// Подготовка данных
func makeStrings(n int) []string {
    strs := make([]string, n)
    for i := range strs {
        strs[i] = "hello"
    }
    return strs
}

func BenchmarkConcatPlus(b *testing.B) {
    strs := makeStrings(100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ConcatPlus(strs)
    }
}

func BenchmarkConcatBuilder(b *testing.B) {
    strs := makeStrings(100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ConcatBuilder(strs)
    }
}

func BenchmarkConcatBuffer(b *testing.B) {
    strs := makeStrings(100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ConcatBuffer(strs)
    }
}

func BenchmarkConcatJoin(b *testing.B) {
    strs := makeStrings(100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ConcatJoin(strs)
    }
}

// Результаты (примерные):
// BenchmarkConcatPlus-8      10000    150000 ns/op   53000 B/op   99 allocs/op
// BenchmarkConcatBuilder-8  200000      8000 ns/op    1024 B/op    8 allocs/op
// BenchmarkConcatBuffer-8   150000      9000 ns/op    1024 B/op    5 allocs/op
// BenchmarkConcatJoin-8     300000      5000 ns/op     512 B/op    1 allocs/op
```

### Пример 4: Бенчмарк с разными размерами данных

```go
package sort

import (
    "math/rand"
    "sort"
    "testing"
)

func BenchmarkSort(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            // Создаём данные
            original := make([]int, size)
            for i := range original {
                original[i] = rand.Intn(size)
            }
            
            b.ResetTimer()
            
            for i := 0; i < b.N; i++ {
                // Копируем, чтобы каждая итерация сортировала несортированный массив
                b.StopTimer()
                data := make([]int, len(original))
                copy(data, original)
                b.StartTimer()
                
                sort.Ints(data)
            }
        })
    }
}

// Результаты:
// BenchmarkSort/size=10-8       5000000    250 ns/op
// BenchmarkSort/size=100-8       500000   3000 ns/op
// BenchmarkSort/size=1000-8       30000  45000 ns/op
// BenchmarkSort/size=10000-8       2000 600000 ns/op
```

### Пример 5: Бенчмарк с отслеживанием памяти

```go
package memory

import "testing"

type User struct {
    ID    int
    Name  string
    Email string
}

func CreateUserSlice(n int) []*User {
    users := make([]*User, n)
    for i := 0; i < n; i++ {
        users[i] = &User{
            ID:    i,
            Name:  fmt.Sprintf("User%d", i),
            Email: fmt.Sprintf("user%d@example.com", i),
        }
    }
    return users
}

func CreateUserValue(n int) []User {
    users := make([]User, n)
    for i := 0; i < n; i++ {
        users[i] = User{
            ID:    i,
            Name:  fmt.Sprintf("User%d", i),
            Email: fmt.Sprintf("user%d@example.com", i),
        }
    }
    return users
}

func BenchmarkCreateUserSlice(b *testing.B) {
    b.ReportAllocs()  // включаем отчёт по памяти
    
    for i := 0; i < b.N; i++ {
        CreateUserSlice(1000)
    }
}

func BenchmarkCreateUserValue(b *testing.B) {
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        CreateUserValue(1000)
    }
}

// Результаты:
// BenchmarkCreateUserSlice-8   10000   100000 ns/op   88000 B/op   1001 allocs/op
// BenchmarkCreateUserValue-8   20000    60000 ns/op   48000 B/op      1 allocs/op
```

### Пример 6: Параллельный бенчмарк

```go
package counter

import (
    "sync"
    "sync/atomic"
    "testing"
)

// Счётчик с мьютексом
type MutexCounter struct {
    mu    sync.Mutex
    value int64
}

func (c *MutexCounter) Inc() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}

// Атомарный счётчик
type AtomicCounter struct {
    value int64
}

func (c *AtomicCounter) Inc() {
    atomic.AddInt64(&c.value, 1)
}

func BenchmarkMutexCounterParallel(b *testing.B) {
    counter := &MutexCounter{}
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counter.Inc()
        }
    })
}

func BenchmarkAtomicCounterParallel(b *testing.B) {
    counter := &AtomicCounter{}
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            counter.Inc()
        }
    })
}

// Результаты:
// BenchmarkMutexCounterParallel-8    20000000    80 ns/op
// BenchmarkAtomicCounterParallel-8  100000000    15 ns/op
```

### Пример 7: Бенчмарк с SetBytes

```go
package io

import (
    "bytes"
    "testing"
)

func BenchmarkBufferWrite(b *testing.B) {
    data := []byte("hello world")
    b.SetBytes(int64(len(data)))  // устанавливаем размер данных
    
    var buf bytes.Buffer
    
    for i := 0; i < b.N; i++ {
        buf.Reset()
        buf.Write(data)
    }
}

// Результат покажет пропускную способность:
// BenchmarkBufferWrite-8   50000000   30 ns/op   366.67 MB/s
```

### Пример 8: Сравнение map vs slice

```go
package lookup

import "testing"

func BenchmarkMapLookup(b *testing.B) {
    m := make(map[int]string)
    for i := 0; i < 1000; i++ {
        m[i] = fmt.Sprintf("value%d", i)
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _ = m[500]  // lookup в середине
    }
}

func BenchmarkSliceLookup(b *testing.B) {
    s := make([]string, 1000)
    for i := 0; i < 1000; i++ {
        s[i] = fmt.Sprintf("value%d", i)
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _ = s[500]  // lookup по индексу
    }
}

func BenchmarkSliceLinearSearch(b *testing.B) {
    type item struct {
        key   int
        value string
    }
    s := make([]item, 1000)
    for i := 0; i < 1000; i++ {
        s[i] = item{key: i, value: fmt.Sprintf("value%d", i)}
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        for _, item := range s {
            if item.key == 500 {
                _ = item.value
                break
            }
        }
    }
}

// Результаты:
// BenchmarkMapLookup-8           100000000   15 ns/op
// BenchmarkSliceLookup-8        1000000000    1 ns/op
// BenchmarkSliceLinearSearch-8    5000000  300 ns/op
```

---

## ⚠️ Частые ошибки

### 1. Компилятор оптимизирует код

```go
// ❌ ПЛОХО — результат не используется, компилятор может убрать вызов
func BenchmarkBad(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Calculate()  // может быть оптимизировано
    }
}

// ✅ ХОРОШО — используем результат
var result int

func BenchmarkGood(b *testing.B) {
    var r int
    for i := 0; i < b.N; i++ {
        r = Calculate()
    }
    result = r  // предотвращаем оптимизацию
}
```

### 2. Не сбросили таймер после setup

```go
// ❌ ПЛОХО — setup включается в измерение
func BenchmarkBad(b *testing.B) {
    data := expensiveSetup()  // долгий setup
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}

// ✅ ХОРОШО
func BenchmarkGood(b *testing.B) {
    data := expensiveSetup()
    b.ResetTimer()  // сбрасываем таймер!
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}
```

### 3. Аллокация внутри цикла

```go
// ❌ ПЛОХО — измеряем аллокацию, а не функцию
func BenchmarkBad(b *testing.B) {
    for i := 0; i < b.N; i++ {
        data := make([]int, 1000)  // аллокация каждый раз
        Process(data)
    }
}

// ✅ ХОРОШО — если аллокация не часть измерения
func BenchmarkGood(b *testing.B) {
    data := make([]int, 1000)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Process(data)
    }
}
```

### 4. Неправильное использование b.N

```go
// ❌ НЕПРАВИЛЬНО
func BenchmarkBad(b *testing.B) {
    for i := 0; i < 1000; i++ {  // фиксированное количество!
        Function()
    }
}

// ✅ ПРАВИЛЬНО
func BenchmarkGood(b *testing.B) {
    for i := 0; i < b.N; i++ {  // используем b.N
        Function()
    }
}
```

---

## 🏋️ Практические задания

### Задание 1: Бенчмарк сортировки

Сравните производительность пузырьковой сортировки и стандартной.

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
// sort_test.go
package sort

import (
    "math/rand"
    "testing"
)

func generateSlice(n int) []int {
    s := make([]int, n)
    for i := range s {
        s[i] = rand.Intn(1000)
    }
    return s
}

var result []int

func BenchmarkBubbleSort100(b *testing.B) {
    data := generateSlice(100)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        result = BubbleSort(data)
    }
}

func BenchmarkStandardSort100(b *testing.B) {
    data := generateSlice(100)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        result = StandardSort(data)
    }
}

// Добавьте бенчмарки для 1000 и 10000 элементов
```

**Запуск:**
```bash
go test -bench=. -benchmem
```

**Ожидаемый вывод:**
```
BenchmarkBubbleSort100-8         50000     25000 ns/op     896 B/op    1 allocs/op
BenchmarkStandardSort100-8      500000      3000 ns/op     896 B/op    1 allocs/op
```

**Баллы:** 15

### Задание 2: Конкатенация строк

Сравните разные способы конкатенации строк.

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

```go
// strings_test.go
package strings

import "testing"

func generateStrings(n int) []string {
    result := make([]string, n)
    for i := range result {
        result[i] = "hello"
    }
    return result
}

func BenchmarkConcatPlus(b *testing.B) {
    strs := generateStrings(100)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ConcatPlus(strs)
    }
}

// Добавьте бенчмарки для остальных функций
```

**Баллы:** 10

### Задание 3: Map vs Slice поиск

Сравните поиск в map и в slice.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 10

### Задание 4: Бенчмарк с параллелизмом

Используйте b.RunParallel для параллельного бенчмарка.

**Начальный код:**
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Запуск:**
```bash
go test -bench=Hash -cpu=1,2,4,8
```

**Баллы:** 10

### Задание 5: Оптимизация аллокаций

Оптимизируйте функцию на основе данных бенчмарка.

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
// alloc_test.go
package alloc

import "testing"

func BenchmarkProcessSlow(b *testing.B) {
    data := make([]int, 1000)
    for i := range data {
        data[i] = i - 500
    }
    
    b.ResetTimer()
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        ProcessSlow(data)
    }
}

func BenchmarkProcessFast(b *testing.B) {
    data := make([]int, 1000)
    for i := range data {
        data[i] = i - 500
    }
    
    b.ResetTimer()
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        ProcessFast(data)
    }
}
```

**Ожидаемый вывод:**
```
BenchmarkProcessSlow-8     50000    30000 ns/op   41000 B/op   12 allocs/op
BenchmarkProcessFast-8    100000    10000 ns/op    8000 B/op    1 allocs/op
```

**Баллы:** 15

---

## 🔗 Полезные ссылки

- [Go Testing — Benchmarks](https://pkg.go.dev/testing#hdr-Benchmarks)
- [Benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat)
- [High Performance Go Workshop](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html)
