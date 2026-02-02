# Стандартные потоки ввода-вывода и io.Copy

## 💡 Ключевые идеи

1. **os.Stdin** — стандартный ввод (клавиатура)
2. **os.Stdout** — стандартный вывод (консоль)
3. **os.Stderr** — стандартный поток ошибок
4. **io.Copy** — копирует данные из Reader в Writer
5. **io.CopyN** — копирует N байт
6. **Pipe** — связывает Reader и Writer

---

## 📋 Синтаксис

### Стандартные потоки

```go
os.Stdin   // *os.File, io.Reader
os.Stdout  // *os.File, io.Writer
os.Stderr  // *os.File, io.Writer
```

### io.Copy

```go
n, err := io.Copy(dst Writer, src Reader)
n, err := io.CopyN(dst Writer, src Reader, n int64)
n, err := io.CopyBuffer(dst, src, buf []byte)
```

---

## 💻 Примеры кода

### Вывод в os.Stdout

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // fmt.Println использует os.Stdout
    fmt.Println("Hello, World!")
    
    // Прямая запись в os.Stdout
    os.Stdout.WriteString("Direct write\n")
    
    // Используя fmt.Fprintln
    fmt.Fprintln(os.Stdout, "Using Fprintln")
    
    // Запись байтов
    os.Stdout.Write([]byte("Bytes\n"))
}
```

### Вывод ошибок в os.Stderr

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Обычный вывод
    fmt.Println("Normal output to stdout")
    
    // Вывод ошибок
    fmt.Fprintln(os.Stderr, "Error message to stderr")
    
    // При перенаправлении можно разделить:
    // go run main.go > output.txt 2> errors.txt
}
```

### io.Copy: файл в консоль

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    file, err := os.Open("data.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        return
    }
    defer file.Close()
    
    // Копируем содержимое файла в консоль
    n, err := io.Copy(os.Stdout, file)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Copy error:", err)
        return
    }
    
    fmt.Fprintf(os.Stderr, "\n--- Copied %d bytes ---\n", n)
}
```

### io.Copy: stdin в файл

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    file, err := os.Create("input.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        return
    }
    defer file.Close()
    
    fmt.Println("Enter text (Ctrl+D to finish):")
    
    // Копируем ввод пользователя в файл
    n, err := io.Copy(file, os.Stdin)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error:", err)
        return
    }
    
    fmt.Printf("Saved %d bytes to file\n", n)
}
```

### io.CopyN: копирование N байт

```go
package main

import (
    "fmt"
    "io"
    "os"
    "strings"
)

func main() {
    reader := strings.NewReader("Hello, World! This is a test.")
    
    // Копируем только первые 13 байт
    n, err := io.CopyN(os.Stdout, reader, 13)
    fmt.Println()  // новая строка
    
    if err != nil && err != io.EOF {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Printf("Copied %d bytes\n", n)
}
// Output:
// Hello, World!
// Copied 13 bytes
```

### io.Copy с кастомным Reader

```go
package main

import (
    "fmt"
    "io"
    "os"
)

// Reader, который преобразует в uppercase
type UpperReader struct {
    src io.Reader
}

func (r UpperReader) Read(p []byte) (int, error) {
    n, err := r.src.Read(p)
    for i := 0; i < n; i++ {
        if p[i] >= 'a' && p[i] <= 'z' {
            p[i] -= 32  // a-z -> A-Z
        }
    }
    return n, err
}

func main() {
    file, err := os.Open("data.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()
    
    upperReader := UpperReader{src: file}
    io.Copy(os.Stdout, upperReader)
}
```

### io.Pipe: связь Reader-Writer

```go
package main

import (
    "fmt"
    "io"
)

func main() {
    // Создаём pipe
    reader, writer := io.Pipe()
    
    // Горутина-писатель
    go func() {
        defer writer.Close()
        
        for i := 1; i <= 5; i++ {
            fmt.Fprintf(writer, "Line %d\n", i)
        }
    }()
    
    // Читатель в main
    buf := make([]byte, 100)
    for {
        n, err := reader.Read(buf)
        if err == io.EOF {
            break
        }
        fmt.Print(string(buf[:n]))
    }
}
```

### TeeReader: чтение с дублированием

```go
package main

import (
    "fmt"
    "io"
    "os"
    "strings"
)

func main() {
    input := strings.NewReader("Hello, TeeReader!")
    
    // Создаём файл для копии
    logFile, _ := os.Create("copy.txt")
    defer logFile.Close()
    
    // TeeReader: данные читаются из input, копируются в logFile
    tee := io.TeeReader(input, logFile)
    
    // Читаем через tee
    data, _ := io.ReadAll(tee)
    fmt.Println("Read:", string(data))
    
    // copy.txt теперь содержит те же данные
}
```

### MultiWriter: запись в несколько мест

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    // Создаём файл для лога
    logFile, _ := os.Create("app.log")
    defer logFile.Close()
    
    // MultiWriter пишет и в файл, и в консоль
    multiWriter := io.MultiWriter(os.Stdout, logFile)
    
    fmt.Fprintln(multiWriter, "Starting application...")
    fmt.Fprintln(multiWriter, "Processing...")
    fmt.Fprintln(multiWriter, "Done!")
    
    // Вывод появится и в консоли, и в файле
}
```

### Практический пример: прогресс копирования

```go
package main

import (
    "fmt"
    "io"
    "os"
    "strings"
)

type ProgressWriter struct {
    total   int64
    current int64
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
    n := len(p)
    pw.current += int64(n)
    
    percent := float64(pw.current) / float64(pw.total) * 100
    fmt.Printf("\rProgress: %.1f%%", percent)
    
    return n, nil
}

func main() {
    // Имитация большого файла
    data := strings.Repeat("Hello World! ", 1000)
    reader := strings.NewReader(data)
    
    // Writer с прогрессом
    progress := &ProgressWriter{total: int64(len(data))}
    
    // Копируем с отслеживанием
    multiWriter := io.MultiWriter(io.Discard, progress)
    io.Copy(multiWriter, reader)
    
    fmt.Println("\nDone!")
}
```

---

## ⚠️ Частые ошибки

### 1. Путаница stdout/stderr

```go
// ❌ Ошибки в stdout
fmt.Println("Error: something went wrong")

// ✅ Ошибки в stderr
fmt.Fprintln(os.Stderr, "Error: something went wrong")
```

### 2. Не проверяем ошибку io.Copy

```go
// ❌ Игнорируем ошибку
io.Copy(dst, src)

// ✅ Проверяем
n, err := io.Copy(dst, src)
if err != nil {
    return fmt.Errorf("copy failed: %w", err)
}
```

### 3. Использование io.Copy для маленьких данных

```go
// ❌ Избыточно для маленьких данных
io.Copy(file, strings.NewReader("Hi"))

// ✅ Проще напрямую
file.WriteString("Hi")
```

---

## 📝 Практика

### Задача 1: cat
Реализуйте команду cat: вывод файлов в stdout.

### Задача 2: tee
Реализуйте команду tee: stdin → stdout + file.

### Задача 3: head/tail
Реализуйте head -n (первые N строк) и tail -n (последние N строк).

### Задача 4: wc
Реализуйте команду wc: подсчёт строк, слов, байт.

### Задача 5: Прогресс-бар
Создайте Writer с отображением прогресса.

### Задача 6: Throttle writer
Ограничьте скорость записи (N байт в секунду).

### Задача 7: Redirect
Перенаправьте stdout в файл программно.

### Задача 8: Broadcast
Создайте Writer, дублирующий данные в N других Writer'ов.
