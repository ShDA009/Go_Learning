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

## 📝 Практика

### Задача 1: URL checker
Проверьте доступность списка URL.

### Задача 2: API client
Клиент для публичного API (погода, курсы валют).

### Задача 3: Web scraper
Простой парсер веб-страниц.

### Задача 4: File uploader
Загрузка файла на сервер.

### Задача 5: Parallel downloader
Параллельное скачивание файлов.

### Задача 6: RSS reader
Чтение и парсинг RSS лент.

### Задача 7: REST client
Клиент для REST API с CRUD операциями.

### Задача 8: Webhook sender
Отправка webhook уведомлений.
