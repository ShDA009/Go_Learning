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

## 📝 Практика

### Задача 1: Rot13 Reader
Создайте Reader, который декодирует ROT13 на лету.

### Задача 2: Counting Writer
Создайте Writer, который считает записанные байты.

### Задача 3: Tee Reader
Создайте Reader, который дублирует данные в Writer.

### Задача 4: Limiting Reader
Создайте Reader, ограничивающий количество читаемых байт.

### Задача 5: Multi Reader
Создайте Reader, объединяющий несколько Reader'ов последовательно.

### Задача 6: Pipe
Используйте io.Pipe для связи Reader и Writer.

### Задача 7: Progress Writer
Создайте Writer, показывающий прогресс записи.

### Задача 8: Checksum Writer
Создайте Writer, вычисляющий контрольную сумму данных.
