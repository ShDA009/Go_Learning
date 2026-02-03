# HTTP запросы (net/http)

## 💡 Ключевые идеи

1. **http.Get** — простой GET запрос
2. **http.Post** — POST запрос с телом
3. **http.Response** — структура ответа
4. **resp.Body** — тело ответа (io.ReadCloser)
5. **Обязательное закрытие** — defer resp.Body.Close()

---

## 📋 Синтаксис

### Функции для запросов

```go
// GET запрос
http.Get(url string) (*http.Response, error)

// HEAD запрос
http.Head(url string) (*http.Response, error)

// POST запрос
http.Post(url, contentType string, body io.Reader) (*http.Response, error)

// POST форма
http.PostForm(url string, data url.Values) (*http.Response, error)
```

### Структура Response

```go
type Response struct {
    Status     string      // "200 OK"
    StatusCode int         // 200
    Header     http.Header // заголовки
    Body       io.ReadCloser
    // ...
}
```

---

## 📖 Теория

### HTTP — протокол веба

**HTTP** (HyperText Transfer Protocol) — основа всего веба. Go отлично поддерживает HTTP из коробки.

### Простейший GET запрос

```go
resp, err := http.Get("https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()  // ОБЯЗАТЕЛЬНО!

body, err := io.ReadAll(resp.Body)
fmt.Println(string(body))
```

### ВАЖНО: всегда закрывайте Body!

```go
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()  // Даже если Body не читаем!
```

Незакрытый Body = утечка соединений!

### Проверка статуса

```go
resp, _ := http.Get(url)
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("bad status: %d", resp.StatusCode)
}
```

### HTTP методы

| Функция | Метод | Использование |
|---------|-------|---------------|
| `http.Get` | GET | Получить данные |
| `http.Post` | POST | Отправить данные |
| `http.PostForm` | POST | Отправить форму |
| `http.Head` | HEAD | Только заголовки |

### POST с JSON

```go
data := map[string]string{"name": "John", "age": "30"}
jsonData, _ := json.Marshal(data)

resp, err := http.Post(
    "https://api.example.com/users",
    "application/json",
    bytes.NewBuffer(jsonData),
)
defer resp.Body.Close()
```

### POST форма

```go
formData := url.Values{
    "username": {"john"},
    "password": {"secret"},
}

resp, err := http.PostForm("https://example.com/login", formData)
defer resp.Body.Close()
```

### Чтение заголовков ответа

```go
// Конкретный заголовок
contentType := resp.Header.Get("Content-Type")

// Все заголовки
for key, values := range resp.Header {
    fmt.Printf("%s: %v\n", key, values)
}
```

### Обработка JSON ответа

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

resp, _ := http.Get("https://api.example.com/user/1")
defer resp.Body.Close()

var user User
json.NewDecoder(resp.Body).Decode(&user)
fmt.Printf("User: %+v\n", user)
```

---

## 📋 Синтаксис

### Структура Response

```go
type Response struct {
    Status     string      // "200 OK"
    StatusCode int         // 200
    Header     http.Header // заголовки
    Body       io.ReadCloser
    // ...
}
```

---

## 💻 Примеры кода

### Простой GET запрос

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    resp, err := http.Get("https://httpbin.org/get")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    // Читаем тело ответа
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Read error:", err)
        return
    }
    
    fmt.Println("Status:", resp.Status)
    fmt.Println("Body:", string(body))
}
```

### Проверка статуса ответа

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    resp, err := http.Get("https://httpbin.org/status/404")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    // Проверяем статус
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Bad status: %d %s\n", resp.StatusCode, resp.Status)
        return
    }
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Чтение заголовков ответа

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    resp, err := http.Get("https://httpbin.org/get")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    // Все заголовки
    fmt.Println("=== Headers ===")
    for key, values := range resp.Header {
        for _, value := range values {
            fmt.Printf("%s: %s\n", key, value)
        }
    }
    
    // Конкретный заголовок
    contentType := resp.Header.Get("Content-Type")
    fmt.Println("\nContent-Type:", contentType)
}
```

### POST запрос с JSON

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type RequestData struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    data := RequestData{
        Name:  "John",
        Email: "john@example.com",
    }
    
    // Кодируем в JSON
    jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Println("JSON error:", err)
        return
    }
    
    // Отправляем POST
    resp, err := http.Post(
        "https://httpbin.org/post",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        fmt.Println("Request error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### POST форма

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "net/url"
)

func main() {
    // Данные формы
    formData := url.Values{
        "username": {"john"},
        "password": {"secret123"},
    }
    
    resp, err := http.PostForm("https://httpbin.org/post", formData)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Копирование ответа в stdout

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func main() {
    resp, err := http.Get("https://httpbin.org/html")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    // Копируем напрямую в stdout
    io.Copy(os.Stdout, resp.Body)
}
```

### Скачивание файла

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
)

func downloadFile(url, filepath string) error {
    // Создаём файл
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()
    
    // Скачиваем
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Проверяем статус
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %s", resp.Status)
    }
    
    // Копируем в файл
    _, err = io.Copy(out, resp.Body)
    return err
}

func main() {
    url := "https://golang.org/doc/gopher/frontpage.png"
    filepath := "gopher.png"
    
    err := downloadFile(url, filepath)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println("Downloaded:", filepath)
}
```

### Парсинг JSON ответа

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type HTTPBinResponse struct {
    Origin string `json:"origin"`
    URL    string `json:"url"`
}

func main() {
    resp, err := http.Get("https://httpbin.org/get")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    var result HTTPBinResponse
    
    // Декодируем JSON напрямую из Body
    decoder := json.NewDecoder(resp.Body)
    if err := decoder.Decode(&result); err != nil {
        fmt.Println("JSON error:", err)
        return
    }
    
    fmt.Println("Origin:", result.Origin)
    fmt.Println("URL:", result.URL)
}
```

---

## ⚠️ Частые ошибки

### 1. Не закрывают Body

```go
// ❌ Утечка ресурсов!
resp, _ := http.Get(url)
body, _ := io.ReadAll(resp.Body)
// Body не закрыт!

// ✅ Всегда закрывайте Body
resp, err := http.Get(url)
if err != nil {
    return
}
defer resp.Body.Close()
```

### 2. Не проверяют статус

```go
// ❌ Код ошибки игнорируется
resp, _ := http.Get(url)
body, _ := io.ReadAll(resp.Body)
// body может содержать сообщение об ошибке!

// ✅ Проверяйте StatusCode
if resp.StatusCode != http.StatusOK {
    fmt.Println("Error:", resp.Status)
    return
}
```

### 3. Путают err и StatusCode

```go
// ❌ 404 — это не err!
resp, err := http.Get("https://example.com/notfound")
if err != nil {
    // НЕ сработает для 404!
}

// ✅ err — только для сетевых ошибок
// StatusCode — для HTTP ошибок
if err != nil {
    // Сетевая ошибка (DNS, таймаут, и т.д.)
}
if resp.StatusCode >= 400 {
    // HTTP ошибка (404, 500, и т.д.)
}
```

### 4. Используют http.Get для всего

```go
// ❌ http.Get не позволяет настроить заголовки
http.Get(url)

// ✅ Для кастомных заголовков используйте Client
```

---

## 🏋️ Практические задания

### Задание 1: QueryRow

Получите одну строку.

**Ожидаемый результат:**
```
QueryRow: db.QueryRow("SELECT...").Scan(&val)
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("QueryRow: db.QueryRow(\"SELECT...\").Scan(&val)")
}
```

**Ожидаемый вывод:**
```
QueryRow: db.QueryRow("SELECT...").Scan(&val)
```

**Баллы:** 10

### Задание 2: sql.NullString

Обработайте NULL значения.

**Ожидаемый результат:**
```
NULL: var name sql.NullString; if name.Valid {...}
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("NULL: var name sql.NullString; if name.Valid {...}")
}
```

**Ожидаемый вывод:**
```
NULL: var name sql.NullString; if name.Valid {...}
```

**Баллы:** 15

### Задание 3: Транзакции

Используйте транзакции.

**Ожидаемый результат:**
```
Tx: tx, _ := db.Begin(); tx.Commit() или tx.Rollback()
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Tx: tx, _ := db.Begin(); tx.Commit() или tx.Rollback()")
}
```

**Ожидаемый вывод:**
```
Tx: tx, _ := db.Begin(); tx.Commit() или tx.Rollback()
```

**Баллы:** 20

### Задание 4: Context в запросах

Используйте context для timeout.

**Ожидаемый результат:**
```
Context: db.QueryContext(ctx, "SELECT...")
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Context: db.QueryContext(ctx, \"SELECT...\")")
}
```

**Ожидаемый вывод:**
```
Context: db.QueryContext(ctx, "SELECT...")
```

**Баллы:** 15

### Задание 5: Connection pool

Настройте пул соединений.

**Ожидаемый результат:**
```
Pool: db.SetMaxOpenConns(25); db.SetMaxIdleConns(5)
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Pool: db.SetMaxOpenConns(25); db.SetMaxIdleConns(5)")
}
```

**Ожидаемый вывод:**
```
Pool: db.SetMaxOpenConns(25); db.SetMaxIdleConns(5)
```

**Баллы:** 15
