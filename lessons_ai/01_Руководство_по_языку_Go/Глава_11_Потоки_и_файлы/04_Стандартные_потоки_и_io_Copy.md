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

## 📖 Теория

### Три стандартных потока

Каждая программа получает три "файла" при запуске:

| Поток | Переменная | Назначение |
|-------|------------|------------|
| stdin | `os.Stdin` | Ввод (клавиатура, pipe) |
| stdout | `os.Stdout` | Вывод (консоль) |
| stderr | `os.Stderr` | Ошибки (консоль) |

### stdin — откуда читать данные

```go
// Чтение с клавиатуры
scanner := bufio.NewScanner(os.Stdin)
scanner.Scan()
input := scanner.Text()

// Работает и с pipe!
// echo "hello" | ./myprogram
```

### stdout vs stderr

```go
fmt.Println("Normal output")        // → stdout
fmt.Fprintln(os.Stderr, "Error!")   // → stderr
```

**Зачем разделять?**
```bash
./program > output.txt  # stdout в файл, stderr на экране
./program 2> errors.txt # stdout на экране, stderr в файл
```

### io.Copy — универсальный копировщик

```go
// Копирует ВСЁ из Reader в Writer
n, err := io.Copy(dst, src)
```

**Примеры:**
```go
// Файл в файл
io.Copy(dstFile, srcFile)

// HTTP ответ в файл
io.Copy(file, response.Body)

// Файл на консоль
io.Copy(os.Stdout, file)

// Сеть в файл
io.Copy(file, conn)
```

### io.CopyN — копировать N байт

```go
// Первые 1024 байта
io.CopyN(dst, src, 1024)
```

### io.Pipe — связь между горутинами

```go
r, w := io.Pipe()

go func() {
    defer w.Close()
    w.Write([]byte("data"))
}()

data, _ := io.ReadAll(r)  // "data"
```

### Утилиты io

```go
io.ReadAll(r)           // прочитать всё
io.ReadFull(r, buf)     // заполнить буфер полностью
io.WriteString(w, s)    // записать строку
io.Discard              // /dev/null — игнорирует данные
io.LimitReader(r, n)    // ограничить чтение
io.TeeReader(r, w)      // читать и копировать
io.MultiWriter(w1, w2)  // писать в несколько
io.MultiReader(r1, r2)  // читать последовательно
```

### Пример: прогресс скачивания

```go
type ProgressWriter struct {
    Total   int64
    Written int64
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
    pw.Written += int64(len(p))
    fmt.Printf("\r%.2f%%", float64(pw.Written)/float64(pw.Total)*100)
    return len(p), nil
}

// Использование
io.Copy(io.MultiWriter(file, &progress), response.Body)
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

## 🏋️ Практические задания

### Задание 1: html/template базовый

Создайте и выполните шаблон.

**Ожидаемый результат:**
```
Шаблон: template.New().Parse()
```

**Критерии приёмки:**
- template.New("name")
- Parse для парсинга
- Execute для выполнения

**Подсказки:**
- t, _ := template.New("t").Parse("Hello {{.}}")
- t.Execute(w, "World")

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Шаблон: template.New().Parse()")
}
```

**Ожидаемый вывод:**
```
Шаблон: template.New().Parse()
```

**Баллы:** 10

### Задание 2: Передача данных в шаблон

Передайте структуру в шаблон.

**Ожидаемый результат:**
```
Данные: {{.Field}} в шаблоне
```

**Критерии приёмки:**
- {{.FieldName}} для доступа к полям
- Структура передаётся в Execute

**Подсказки:**
- type Data struct { Name string }
- {{.Name}} в шаблоне

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Данные: {{.Field}} в шаблоне")
}
```

**Ожидаемый вывод:**
```
Данные: {{.Field}} в шаблоне
```

**Баллы:** 10

### Задание 3: Условия в шаблоне

Используйте if в шаблоне.

**Ожидаемый результат:**
```
Условие: {{if .Value}}...{{else}}...{{end}}
```

**Критерии приёмки:**
- {{if .Condition}}...{{end}}
- {{else}} для альтернативы

**Подсказки:**
- if проверяет истинность
- Пустая строка, 0, nil — ложь

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Условие: {{if .Value}}...{{else}}...{{end}}")
}
```

**Ожидаемый вывод:**
```
Условие: {{if .Value}}...{{else}}...{{end}}
```

**Баллы:** 15

### Задание 4: Циклы в шаблоне

Используйте range для итерации.

**Ожидаемый результат:**
```
Цикл: {{range .Items}}...{{end}}
```

**Критерии приёмки:**
- {{range .Slice}}...{{end}}
- . внутри range — текущий элемент

**Подсказки:**
- {{range .Items}}<li>{{.}}</li>{{end}}
- {{range $i, $v := .Items}}...

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Цикл: {{range .Items}}...{{end}}")
}
```

**Ожидаемый вывод:**
```
Цикл: {{range .Items}}...{{end}}
```

**Баллы:** 15

### Задание 5: Вложенные шаблоны

Используйте define и template.

**Ожидаемый результат:**
```
Вложение: {{define "name"}}...{{end}}, {{template "name" .}}
```

**Критерии приёмки:**
- {{define "name"}} определяет блок
- {{template "name" .}} вставляет

**Подсказки:**
- Для layout и компонентов
- Передавайте данные через .

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Вложение: {{define \"name\"}}...{{end}}, {{template \"name\" .}}")
}
```

**Ожидаемый вывод:**
```
Вложение: {{define "name"}}...{{end}}, {{template "name" .}}
```

**Баллы:** 20
