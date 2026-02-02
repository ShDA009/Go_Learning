# Пакет strings: Операции со строками

## 💡 Ключевые идеи

1. **strings** — стандартный пакет для работы со строками UTF-8
2. **Неизменяемость** — все функции возвращают новую строку, не изменяя исходную
3. **Unicode-совместимость** — функции корректно работают с UTF-8
4. **Производительность** — strings.Builder для эффективной конкатенации

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

## 📝 Практика

### Задача 1: Регистр
Напишите функцию `SwapCase(s string) string`, которая меняет регистр каждой буквы на противоположный.

### Задача 2: Подсчёт слов
Напишите функцию `WordCount(s string) int`, которая считает количество слов в строке.

### Задача 3: Анаграмма
Напишите функцию `IsAnagram(s1, s2 string) bool`, которая проверяет, являются ли две строки анаграммами.

### Задача 4: Центрирование
Напишите функцию `Center(s string, width int) string`, которая центрирует строку в заданной ширине.

### Задача 5: Маскирование
Напишите функцию `MaskEmail(email string) string`, которая скрывает часть email: "john@example.com" → "j***@example.com".

### Задача 6: CamelCase → snake_case
Напишите функцию `ToSnakeCase(s string) string`, которая преобразует CamelCase в snake_case.

### Задача 7: Сжатие пробелов
Напишите функцию `CompressSpaces(s string) string`, которая заменяет множественные пробелы на один.

### Задача 8: Валидация
Напишите функцию `IsValidUsername(s string) bool`, которая проверяет: только буквы, цифры и _, длина 3-20 символов.
