# Пакет strings: Операции со строками

## 💡 Ключевые идеи

1. **strings** — стандартный пакет для работы со строками UTF-8
2. **Неизменяемость** — все функции возвращают новую строку, не изменяя исходную
3. **Unicode-совместимость** — функции корректно работают с UTF-8
4. **Производительность** — strings.Builder для эффективной конкатенации

---

## 📖 Теория

### Пакет strings — ваш основной инструмент

Пакет `strings` содержит десятки функций для работы со строками. Все они работают с UTF-8 и учитывают Unicode.

### Главный принцип: строки неизменяемы

Все функции **возвращают новую строку**, не изменяя исходную:
```go
s := "hello"
upper := strings.ToUpper(s)
fmt.Println(s)     // "hello" — без изменений
fmt.Println(upper) // "HELLO" — новая строка
```

### Категории функций

**Поиск:**
```go
strings.Contains("hello", "ell")    // true
strings.HasPrefix("hello", "he")    // true
strings.HasSuffix("hello", "lo")    // true
strings.Index("hello", "l")         // 2
strings.Count("banana", "a")        // 3
```

**Регистр:**
```go
strings.ToUpper("hello")  // "HELLO"
strings.ToLower("HELLO")  // "hello"
strings.Title("hello go") // "Hello Go" (устарела!)
```

**Разбиение и объединение:**
```go
strings.Split("a,b,c", ",")     // ["a", "b", "c"]
strings.Join([]string{"a","b"}, "-") // "a-b"
strings.Fields("  a   b  ")     // ["a", "b"]
```

**Замена и очистка:**
```go
strings.Replace("aaa", "a", "b", 2)  // "bba"
strings.ReplaceAll("aaa", "a", "b")  // "bbb"
strings.TrimSpace("  hello  ")       // "hello"
strings.Trim("!!!hello!!!", "!")     // "hello"
```

### Проблема конкатенации

```go
// ПЛОХО — O(n²), создаёт много временных строк
result := ""
for i := 0; i < 1000; i++ {
    result += "x"
}

// ХОРОШО — O(n), один буфер
var b strings.Builder
for i := 0; i < 1000; i++ {
    b.WriteString("x")
}
result := b.String()
```

### strings.Builder — эффективная сборка

```go
var b strings.Builder
b.WriteString("Hello")
b.WriteByte(' ')
b.WriteRune('🌍')
s := b.String() // "Hello 🌍"
```

**Преимущества:**
- Минимум аллокаций
- Можно предвыделить размер: `b.Grow(100)`

### Полезные функции для валидации

```go
// Проверка email (упрощённо)
func isEmail(s string) bool {
    return strings.Contains(s, "@") && 
           strings.Contains(s, ".")
}

// Нормализация ввода
func normalize(s string) string {
    s = strings.TrimSpace(s)
    s = strings.ToLower(s)
    return s
}
```

### Пакет strconv — конвертация типов

```go
import "strconv"

// Строка → число
n, _ := strconv.Atoi("42")      // int
f, _ := strconv.ParseFloat("3.14", 64)

// Число → строка  
s := strconv.Itoa(42)           // "42"
s := strconv.FormatFloat(3.14, 'f', 2, 64) // "3.14"
```

---

## 📋 Синтаксис

### Основные функции пакета strings

```go
import "strings"

// Регистр
strings.ToUpper(s)           // в верхний регистр
strings.ToLower(s)           // в нижний регистр
strings.Title(s)             // первые буквы слов заглавные

// Поиск
strings.Contains(s, substr)  // содержит подстроку?
strings.HasPrefix(s, prefix) // начинается с?
strings.HasSuffix(s, suffix) // заканчивается на?
strings.Index(s, substr)     // индекс первого вхождения
strings.LastIndex(s, substr) // индекс последнего вхождения
strings.Count(s, substr)     // количество вхождений

// Разбиение и объединение
strings.Split(s, sep)        // строка → срез
strings.Join(slice, sep)     // срез → строка
strings.Fields(s)            // разбить по пробелам

// Замена и удаление
strings.Replace(s, old, new, n)   // заменить n вхождений
strings.ReplaceAll(s, old, new)   // заменить все
strings.Trim(s, cutset)           // удалить символы с обеих сторон
strings.TrimSpace(s)              // удалить пробелы
strings.TrimPrefix(s, prefix)     // удалить префикс
strings.TrimSuffix(s, suffix)     // удалить суффикс
```

---

## 💻 Примеры кода

### Преобразование регистра

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello World"
    
    fmt.Println(strings.ToUpper(s))  // HELLO WORLD
    fmt.Println(strings.ToLower(s))  // hello world
    fmt.Println(strings.Title(s))    // Hello World
    
    // Работает с UTF-8
    rus := "привет мир"
    fmt.Println(strings.ToUpper(rus))  // ПРИВЕТ МИР
}
```

### Поиск подстроки

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello, World! Hello, Go!"
    
    // Содержит подстроку?
    fmt.Println(strings.Contains(s, "World"))  // true
    fmt.Println(strings.Contains(s, "world"))  // false (регистрозависимо)
    
    // Содержит любой из символов?
    fmt.Println(strings.ContainsAny(s, "xyz"))  // false
    fmt.Println(strings.ContainsAny(s, "aeo"))  // true
    
    // Начинается/заканчивается на?
    fmt.Println(strings.HasPrefix(s, "Hello"))  // true
    fmt.Println(strings.HasSuffix(s, "Go!"))    // true
    
    // Количество вхождений
    fmt.Println(strings.Count(s, "Hello"))  // 2
    fmt.Println(strings.Count(s, "o"))      // 4
}
```

### Поиск индекса

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello, World!"
    
    // Первое вхождение
    fmt.Println(strings.Index(s, "o"))      // 4
    fmt.Println(strings.Index(s, "World"))  // 7
    fmt.Println(strings.Index(s, "xyz"))    // -1 (не найдено)
    
    // Последнее вхождение
    fmt.Println(strings.LastIndex(s, "o"))  // 8
    
    // Первое вхождение любого символа
    fmt.Println(strings.IndexAny(s, "aeiou"))      // 1 (e)
    fmt.Println(strings.LastIndexAny(s, "aeiou"))  // 8 (o)
}
```

### Разбиение строки (Split)

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    // Split — по разделителю
    s1 := "apple,banana,cherry"
    parts := strings.Split(s1, ",")
    fmt.Println(parts)  // [apple banana cherry]
    
    // SplitN — ограничение количества частей
    s2 := "a:b:c:d:e"
    parts2 := strings.SplitN(s2, ":", 3)
    fmt.Println(parts2)  // [a b c:d:e]
    
    // Fields — по пробелам (любому количеству)
    s3 := "  Hello   World   Go  "
    words := strings.Fields(s3)
    fmt.Println(words)  // [Hello World Go]
    
    // FieldsFunc — по условию
    s4 := "a,b;c:d"
    parts4 := strings.FieldsFunc(s4, func(r rune) bool {
        return r == ',' || r == ';' || r == ':'
    })
    fmt.Println(parts4)  // [a b c d]
}
```

### Объединение строк (Join)

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    words := []string{"Hello", "World", "Go"}
    
    // С пробелом
    s1 := strings.Join(words, " ")
    fmt.Println(s1)  // Hello World Go
    
    // С запятой
    s2 := strings.Join(words, ", ")
    fmt.Println(s2)  // Hello, World, Go
    
    // Без разделителя
    s3 := strings.Join(words, "")
    fmt.Println(s3)  // HelloWorldGo
}
```

### Замена подстрок

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello World, Hello Go!"
    
    // Заменить первые n вхождений
    s1 := strings.Replace(s, "Hello", "Hi", 1)
    fmt.Println(s1)  // Hi World, Hello Go!
    
    // Заменить все вхождения
    s2 := strings.ReplaceAll(s, "Hello", "Hi")
    fmt.Println(s2)  // Hi World, Hi Go!
    
    // Заменить -1 = все (как ReplaceAll)
    s3 := strings.Replace(s, "Hello", "Hi", -1)
    fmt.Println(s3)  // Hi World, Hi Go!
}
```

### Удаление символов (Trim)

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    // TrimSpace — удаляет пробелы с обеих сторон
    s1 := "   Hello World   "
    fmt.Printf("[%s]\n", strings.TrimSpace(s1))  // [Hello World]
    
    // Trim — удаляет указанные символы
    s2 := "###Hello###"
    fmt.Println(strings.Trim(s2, "#"))  // Hello
    
    // TrimLeft / TrimRight
    s3 := "...Hello..."
    fmt.Println(strings.TrimLeft(s3, "."))   // Hello...
    fmt.Println(strings.TrimRight(s3, "."))  // ...Hello
    
    // TrimPrefix / TrimSuffix — удаляет подстроку
    s4 := "Hello World"
    fmt.Println(strings.TrimPrefix(s4, "Hello "))  // World
    fmt.Println(strings.TrimSuffix(s4, " World"))  // Hello
}
```

### Удаление по условию (TrimFunc)

```go
package main

import (
    "fmt"
    "strings"
    "unicode"
)

func main() {
    s := "123Hello456"
    
    // Удаляем цифры с обеих сторон
    result := strings.TrimFunc(s, func(r rune) bool {
        return unicode.IsDigit(r)
    })
    fmt.Println(result)  // Hello
    
    // Только слева
    s2 := "   Hello   "
    result2 := strings.TrimLeftFunc(s2, unicode.IsSpace)
    fmt.Printf("[%s]\n", result2)  // [Hello   ]
}
```

### Повторение строки

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    // Repeat — повторяет строку n раз
    s := strings.Repeat("Go! ", 3)
    fmt.Println(s)  // Go! Go! Go! 
    
    // Полезно для создания разделителей
    line := strings.Repeat("-", 40)
    fmt.Println(line)  // ----------------------------------------
}
```

### strings.Builder (эффективная конкатенация)

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    var builder strings.Builder
    
    // Добавляем строки
    builder.WriteString("Hello")
    builder.WriteString(" ")
    builder.WriteString("World")
    
    // Добавляем руну
    builder.WriteRune('!')
    
    // Добавляем байт
    builder.WriteByte('\n')
    
    // Получаем результат
    result := builder.String()
    fmt.Print(result)  // Hello World!
    
    // Сброс для повторного использования
    builder.Reset()
}
```

### strings.Reader (чтение строки как io.Reader)

```go
package main

import (
    "fmt"
    "io"
    "strings"
)

func main() {
    s := "Hello, World!"
    reader := strings.NewReader(s)
    
    // Размер
    fmt.Println("Size:", reader.Size())  // 13
    
    // Чтение
    buf := make([]byte, 5)
    n, _ := reader.Read(buf)
    fmt.Println("Read:", string(buf[:n]))  // Hello
    
    // Чтение всего
    reader.Reset(s)
    data, _ := io.ReadAll(reader)
    fmt.Println("All:", string(data))  // Hello, World!
}
```

### Сравнение строк

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    s1 := "Hello"
    s2 := "hello"
    s3 := "Hello"
    
    // Обычное сравнение (регистрозависимое)
    fmt.Println(s1 == s2)  // false
    fmt.Println(s1 == s3)  // true
    
    // Без учёта регистра
    fmt.Println(strings.EqualFold(s1, s2))  // true
    
    // Compare: -1 если a < b, 0 если a == b, 1 если a > b
    fmt.Println(strings.Compare("abc", "abd"))  // -1
    fmt.Println(strings.Compare("abc", "abc"))  // 0
    fmt.Println(strings.Compare("abd", "abc"))  // 1
}
```

### Практический пример: парсинг CSV

```go
package main

import (
    "fmt"
    "strings"
)

func parseCSV(line string) []string {
    // Разбиваем по запятой
    fields := strings.Split(line, ",")
    
    // Убираем пробелы
    for i, f := range fields {
        fields[i] = strings.TrimSpace(f)
    }
    
    return fields
}

func main() {
    line := "John, Doe, john@example.com, 30"
    fields := parseCSV(line)
    
    fmt.Println("Name:", fields[0], fields[1])
    fmt.Println("Email:", fields[2])
    fmt.Println("Age:", fields[3])
}
```

### Практический пример: форматирование текста

```go
package main

import (
    "fmt"
    "strings"
)

func toSlug(title string) string {
    // В нижний регистр
    slug := strings.ToLower(title)
    
    // Заменяем пробелы на дефисы
    slug = strings.ReplaceAll(slug, " ", "-")
    
    // Удаляем спецсимволы (упрощённо)
    slug = strings.ReplaceAll(slug, "!", "")
    slug = strings.ReplaceAll(slug, "?", "")
    
    return slug
}

func main() {
    title := "Hello World! How Are You?"
    slug := toSlug(title)
    fmt.Println(slug)  // hello-world-how-are-you
}
```

---

## ⚠️ Частые ошибки

### 1. Забыли про регистрозависимость

```go
s := "Hello World"

// ❌ НЕПРАВИЛЬНО — не найдёт
fmt.Println(strings.Contains(s, "hello"))  // false

// ✅ ПРАВИЛЬНО — приводим к одному регистру
fmt.Println(strings.Contains(strings.ToLower(s), "hello"))  // true

// Или используем EqualFold для сравнения
fmt.Println(strings.EqualFold("Hello", "hello"))  // true
```

### 2. Модификация исходной строки

```go
s := "Hello"

// ❌ НЕПРАВИЛЬНО — s не изменится
strings.ToUpper(s)
fmt.Println(s)  // Hello

// ✅ ПРАВИЛЬНО — присваиваем результат
s = strings.ToUpper(s)
fmt.Println(s)  // HELLO
```

### 3. Split с пустой строкой

```go
s := "Hello"

// Осторожно: Split("", "") вернёт срез рун!
parts := strings.Split(s, "")
fmt.Println(parts)  // [H e l l o]
```

### 4. Неэффективная конкатенация в цикле

```go
// ❌ НЕПРАВИЛЬНО — неэффективно
result := ""
for i := 0; i < 1000; i++ {
    result += "x"  // создаёт новую строку каждый раз!
}

// ✅ ПРАВИЛЬНО — используем Builder
var builder strings.Builder
for i := 0; i < 1000; i++ {
    builder.WriteString("x")
}
result := builder.String()
```

### 5. TrimPrefix не удаляет, если нет совпадения

```go
s := "Hello World"

// Не вызовет ошибку, но и не изменит строку
result := strings.TrimPrefix(s, "Hi")
fmt.Println(result)  // Hello World (без изменений)

// Проверяйте HasPrefix, если это важно
```

---

## 🏋️ Практические задания

### Задание 1: interface{} как any

Используйте пустой интерфейс для хранения любого типа.

**Ожидаемый результат:**
```
42 (int)
Hello (string)
true (bool)
```

**Критерии приёмки:**
- Переменная interface{} принимает разные типы
- %T показывает тип
- Три разных типа

**Подсказки:**
- var x interface{} = 42
- fmt.Printf("%v (%T)\n", x, x)

**Начальный код:**
```go
package main

import "fmt"

func main() {
    var x interface{}
    x = 42
    fmt.Printf("%v (%T)\n", x, x)
    x = "Hello"
    fmt.Printf("%v (%T)\n", x, x)
    x = true
    fmt.Printf("%v (%T)\n", x, x)
}
```

**Ожидаемый вывод:**
```
42 (int)
Hello (string)
true (bool)
```

**Баллы:** 10

### Задание 2: Срез interface{}

Создайте срез с разными типами.

**Ожидаемый результат:**
```
Элементы: [1 hello true 3.14]
```

**Критерии приёмки:**
- []interface{} или []any
- Содержит int, string, bool, float64
- Вывод всех элементов

**Подсказки:**
- items := []interface{}{1, "hello", true, 3.14}
- any — алиас для interface{}

**Начальный код:**
```go
package main

import "fmt"

func main() {
    items := []interface{}{1, "hello", true, 3.14}
    fmt.Println("Элементы:", items)
}
```

**Ожидаемый вывод:**
```
Элементы: [1 hello true 3.14]
```

**Баллы:** 10

### Задание 3: Type assertion

Извлеките конкретный тип из interface{}.

**Ожидаемый результат:**
```
Значение: 42
Тип int: true
```

**Критерии приёмки:**
- v, ok := x.(int)
- Проверка ok перед использованием
- Успешное извлечение

**Подсказки:**
- Паника если тип неверный без ok
- Всегда используйте форму с ok

**Начальный код:**
```go
package main

import "fmt"

func main() {
    var x interface{} = 42
    v, ok := x.(int)
    fmt.Println("Значение:", v)
    fmt.Println("Тип int:", ok)
}
```

**Ожидаемый вывод:**
```
Значение: 42
Тип int: true
```

**Баллы:** 15

### Задание 4: Type switch

Определите тип через type switch.

**Ожидаемый результат:**
```
42 — это int
hello — это string
3.14 — это другой тип
```

**Критерии приёмки:**
- switch v := x.(type)
- Три case: int, string, default
- Корректное определение типа

**Подсказки:**
- case int: ...
- default для остальных типов

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
42 — это int
hello — это string
3.14 — это другой тип
```

**Баллы:** 20

### Задание 5: Функция с any параметром

Создайте функцию printType, которая определяет и выводит тип.

**Ожидаемый результат:**
```
Тип: int, Значение: 100
Тип: string, Значение: Go
Тип: []int, Значение: [1 2 3]
```

**Критерии приёмки:**
- Функция принимает any (interface{})
- Выводит тип через %T
- Работает для срезов тоже

**Подсказки:**
- func printType(x any) {...}
- %T — формат типа

**Начальный код:**
```go
package main

import "fmt"

func printType(x any) {
    fmt.Printf("Тип: %T, Значение: %v\n", x, x)
}

func main() {
    printType(100)
    printType("Go")
    printType([]int{1, 2, 3})
}
```

**Ожидаемый вывод:**
```
Тип: int, Значение: 100
Тип: string, Значение: Go
Тип: []int, Значение: [1 2 3]
```

**Баллы:** 20
