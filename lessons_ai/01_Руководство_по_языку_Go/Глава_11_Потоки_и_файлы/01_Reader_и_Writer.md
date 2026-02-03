# Операции ввода-вывода: Reader и Writer

## 💡 Ключевые идеи

1. **io.Reader** — интерфейс для чтения данных из источника
2. **io.Writer** — интерфейс для записи данных в приёмник
3. **Поток данных** — последовательность байтов `[]byte`
4. **Абстракция** — единый API для файлов, сети, памяти, etc.
5. **io.EOF** — маркер конца данных (End Of File)
6. **Композиция** — интерфейсы комбинируются (ReadWriter, ReadCloser)

---

## 📋 Синтаксис

### Интерфейс io.Reader

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

- `p` — буфер для записи прочитанных данных
- `n` — количество прочитанных байт
- `err` — ошибка или io.EOF при конце данных

### Интерфейс io.Writer

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

- `p` — данные для записи
- `n` — количество записанных байт
- `err` — ошибка записи

---

## 📖 Теория

### Самые важные интерфейсы Go

`io.Reader` и `io.Writer` — **фундаментальные абстракции** в Go. Практически всё, что связано с вводом-выводом, работает через них:
- Файлы
- Сетевые соединения
- HTTP запросы/ответы
- Сжатие/шифрование
- Буферы в памяти

### Зачем нужна абстракция?

Один код работает с любым источником данных:

```go
func processData(r io.Reader) {
    // Работает с файлом, сетью, памятью — чем угодно!
    data, _ := io.ReadAll(r)
    process(data)
}

// Из файла
processData(file)

// Из сети
processData(response.Body)

// Из строки
processData(strings.NewReader("hello"))
```

### Как работает Read?

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

**Вы** создаёте буфер, **Reader** заполняет его:

```go
buf := make([]byte, 1024)  // ваш буфер
n, err := reader.Read(buf)  // reader заполняет
data := buf[:n]             // реальные данные
```

### io.EOF — это НЕ ошибка!

`io.EOF` означает "данные закончились". Это нормальное завершение:

```go
for {
    n, err := reader.Read(buf)
    if n > 0 {
        process(buf[:n])  // сначала обработать данные!
    }
    if err == io.EOF {
        break  // нормальный конец
    }
    if err != nil {
        return err  // реальная ошибка
    }
}
```

### Как работает Write?

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

**Вы** даёте данные, **Writer** записывает их:

```go
data := []byte("hello")
n, err := writer.Write(data)
if n < len(data) {
    // Записано меньше чем нужно!
}
```

### Композиция интерфейсов

```go
type ReadWriter interface {
    Reader
    Writer
}

type ReadCloser interface {
    Reader
    Closer
}

type WriteCloser interface {
    Writer
    Closer
}
```

### Полезные типы и функции

```go
// Читатели
strings.NewReader("text")      // строка как Reader
bytes.NewReader([]byte{...})   // байты как Reader

// Писатели
var buf bytes.Buffer           // буфер в памяти
buf.Write(data)
buf.String()                   // получить строку

// Утилиты
io.Copy(dst, src)              // скопировать всё
io.ReadAll(r)                  // прочитать всё
io.WriteString(w, "text")      // записать строку
```

---

## 💻 Примеры кода

### Реализация io.Reader

```go
package main

import (
    "fmt"
    "io"
)

// Кастомный Reader — считывает только цифры
type DigitReader struct {
    data string
    pos  int
}

func (r *DigitReader) Read(p []byte) (int, error) {
    if r.pos >= len(r.data) {
        return 0, io.EOF
    }
    
    n := 0
    for r.pos < len(r.data) && n < len(p) {
        c := r.data[r.pos]
        r.pos++
        if c >= '0' && c <= '9' {
            p[n] = c
            n++
        }
    }
    
    return n, nil
}

func main() {
    reader := &DigitReader{data: "abc123def456ghi"}
    
    buf := make([]byte, 10)
    for {
        n, err := reader.Read(buf)
        if err == io.EOF {
            break
        }
        fmt.Printf("Read: %s\n", string(buf[:n]))
    }
}
// Output: Read: 123456
```

### Реализация io.Writer

```go
package main

import (
    "fmt"
    "strings"
)

// Кастомный Writer — собирает данные в uppercase
type UpperWriter struct {
    builder strings.Builder
}

func (w *UpperWriter) Write(p []byte) (int, error) {
    upper := strings.ToUpper(string(p))
    w.builder.WriteString(upper)
    return len(p), nil
}

func (w *UpperWriter) String() string {
    return w.builder.String()
}

func main() {
    writer := &UpperWriter{}
    
    writer.Write([]byte("hello "))
    writer.Write([]byte("world!"))
    
    fmt.Println(writer.String())  // HELLO WORLD!
}
```

### Использование стандартных Reader'ов

```go
package main

import (
    "fmt"
    "io"
    "strings"
)

func main() {
    // strings.Reader реализует io.Reader
    reader := strings.NewReader("Hello, Go!")
    
    buf := make([]byte, 5)
    for {
        n, err := reader.Read(buf)
        if err == io.EOF {
            break
        }
        fmt.Printf("Read %d bytes: %s\n", n, string(buf[:n]))
    }
}
```

### Комбинированные интерфейсы

```go
package main

import (
    "fmt"
    "io"
)

// io.ReadWriter = io.Reader + io.Writer
// io.ReadCloser = io.Reader + io.Closer
// io.WriteCloser = io.Writer + io.Closer
// io.ReadWriteCloser = io.Reader + io.Writer + io.Closer

type Buffer struct {
    data []byte
    pos  int
}

func (b *Buffer) Read(p []byte) (int, error) {
    if b.pos >= len(b.data) {
        return 0, io.EOF
    }
    n := copy(p, b.data[b.pos:])
    b.pos += n
    return n, nil
}

func (b *Buffer) Write(p []byte) (int, error) {
    b.data = append(b.data, p...)
    return len(p), nil
}

func main() {
    // Buffer реализует io.ReadWriter
    var rw io.ReadWriter = &Buffer{}
    
    rw.Write([]byte("Hello"))
    
    buf := make([]byte, 10)
    n, _ := rw.Read(buf)
    fmt.Println(string(buf[:n]))  // Hello
}
```

### io.Copy — копирование между потоками

```go
package main

import (
    "fmt"
    "io"
    "os"
    "strings"
)

func main() {
    // Копирование из Reader в Writer
    reader := strings.NewReader("Hello from Reader!")
    
    // os.Stdout реализует io.Writer
    n, err := io.Copy(os.Stdout, reader)
    fmt.Println()
    
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Printf("Copied %d bytes\n", n)
}
```

### io.ReadAll — чтение всех данных

```go
package main

import (
    "fmt"
    "io"
    "strings"
)

func main() {
    reader := strings.NewReader("Read all at once!")
    
    // io.ReadAll читает все данные до EOF
    data, err := io.ReadAll(reader)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println(string(data))
}
```

### Практический пример: подсчёт слов

```go
package main

import (
    "fmt"
    "io"
    "strings"
    "unicode"
)

type WordCounter struct {
    count int
    inWord bool
}

func (wc *WordCounter) Write(p []byte) (int, error) {
    for _, b := range p {
        isLetter := unicode.IsLetter(rune(b))
        if isLetter && !wc.inWord {
            wc.count++
        }
        wc.inWord = isLetter
    }
    return len(p), nil
}

func (wc *WordCounter) Count() int {
    return wc.count
}

func main() {
    text := "Hello, World! This is a test."
    reader := strings.NewReader(text)
    counter := &WordCounter{}
    
    io.Copy(counter, reader)
    
    fmt.Printf("Text: %s\n", text)
    fmt.Printf("Word count: %d\n", counter.Count())
}
```

---

## ⚠️ Частые ошибки

### 1. Игнорирование возвращаемого n

```go
// ❌ НЕПРАВИЛЬНО — использовать весь буфер
n, _ := reader.Read(buf)
fmt.Println(string(buf))  // Может содержать мусор!

// ✅ ПРАВИЛЬНО — использовать только прочитанные байты
n, _ := reader.Read(buf)
fmt.Println(string(buf[:n]))
```

### 2. Неправильная обработка EOF

```go
// ❌ EOF — это не ошибка в обычном смысле
n, err := reader.Read(buf)
if err != nil {
    log.Fatal(err)  // Завершится при нормальном EOF!
}

// ✅ Проверяйте EOF отдельно
n, err := reader.Read(buf)
if err == io.EOF {
    // Нормальное завершение
    break
}
if err != nil {
    log.Fatal(err)
}
```

### 3. Чтение может вернуть меньше данных

```go
// ❌ Ожидаем, что Read заполнит весь буфер
buf := make([]byte, 1000)
reader.Read(buf)  // Может прочитать меньше 1000 байт!

// ✅ Используйте io.ReadFull для гарантированного чтения
buf := make([]byte, 1000)
n, err := io.ReadFull(reader, buf)
```

---

## 🏋️ Практические задания

### Задание 1: http.Request структура

Изучите поля структуры Request.

**Ожидаемый результат:**
```
Request: Method, URL, Header, Body
```

**Критерии приёмки:**
- r.Method — HTTP метод
- r.URL — URL запроса
- r.Header — заголовки
- r.Body — тело запроса

**Подсказки:**
- r.URL.Path — путь
- r.URL.Query() — параметры

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Request: Method, URL, Header, Body")
}
```

**Ожидаемый вывод:**
```
Request: Method, URL, Header, Body
```

**Баллы:** 10

### Задание 2: Query параметры

Получите параметры из URL.

**Ожидаемый результат:**
```
Параметры: ?name=value -> r.URL.Query().Get("name")
```

**Критерии приёмки:**
- r.URL.Query() возвращает url.Values
- Get("key") возвращает первое значение
- Пустая строка если нет

**Подсказки:**
- values := r.URL.Query()
- values["key"] — все значения (срез)

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Параметры: ?name=value -> r.URL.Query().Get(\"name\")")
}
```

**Ожидаемый вывод:**
```
Параметры: ?name=value -> r.URL.Query().Get("name")
```

**Баллы:** 10

### Задание 3: Заголовки запроса

Прочитайте заголовки запроса.

**Ожидаемый результат:**
```
Заголовки: r.Header.Get("Content-Type")
```

**Критерии приёмки:**
- r.Header — http.Header (map)
- Get нечувствителен к регистру

**Подсказки:**
- Content-Type, Authorization, User-Agent
- r.Header["Key"] — все значения

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Заголовки: r.Header.Get(\"Content-Type\")")
}
```

**Ожидаемый вывод:**
```
Заголовки: r.Header.Get("Content-Type")
```

**Баллы:** 10

### Задание 4: Чтение тела запроса

Прочитайте тело POST запроса.

**Ожидаемый результат:**
```
Тело: io.ReadAll(r.Body)
```

**Критерии приёмки:**
- r.Body — io.ReadCloser
- io.ReadAll читает полностью
- defer r.Body.Close()

**Подсказки:**
- body, err := io.ReadAll(r.Body)
- Тело можно прочитать только раз

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Тело: io.ReadAll(r.Body)")
}
```

**Ожидаемый вывод:**
```
Тело: io.ReadAll(r.Body)
```

**Баллы:** 15

### Задание 5: Парсинг JSON тела

Десериализуйте JSON из тела запроса.

**Ожидаемый результат:**
```
JSON: json.NewDecoder(r.Body).Decode(&data)
```

**Критерии приёмки:**
- json.NewDecoder для потокового чтения
- Decode в указатель на структуру

**Подсказки:**
- var user User
- json.NewDecoder(r.Body).Decode(&user)

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("JSON: json.NewDecoder(r.Body).Decode(&data)")
}
```

**Ожидаемый вывод:**
```
JSON: json.NewDecoder(r.Body).Decode(&data)
```

**Баллы:** 15
