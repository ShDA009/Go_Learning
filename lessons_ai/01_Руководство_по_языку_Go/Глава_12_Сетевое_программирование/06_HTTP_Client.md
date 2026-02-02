# HTTP Client (http.Client)

## 💡 Ключевые идеи

1. **http.Client** — настраиваемый HTTP клиент
2. **Timeout** — общий таймаут запроса
3. **http.Request** — кастомизация запроса
4. **Заголовки** — добавление своих headers
5. **client.Do** — выполнение кастомного запроса

---

## 📋 Синтаксис

### Создание клиента

```go
// Клиент по умолчанию
client := &http.Client{}

// Клиент с настройками
client := &http.Client{
    Timeout: 10 * time.Second,
    // Transport, Jar, CheckRedirect...
}
```

### Создание запроса

```go
req, err := http.NewRequest(method, url, body)
req, err := http.NewRequestWithContext(ctx, method, url, body)
```

### Методы Client

```go
client.Get(url)              // GET
client.Post(url, type, body) // POST
client.PostForm(url, data)   // POST form
client.Head(url)             // HEAD
client.Do(req)               // кастомный запрос
```

---

## 💻 Примеры кода

### Клиент с таймаутом

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "time"
)

func main() {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    resp, err := client.Get("https://httpbin.org/delay/2")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Кастомные заголовки

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    client := &http.Client{}
    
    req, err := http.NewRequest("GET", "https://httpbin.org/headers", nil)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    // Добавляем заголовки
    req.Header.Set("User-Agent", "MyApp/1.0")
    req.Header.Set("Accept", "application/json")
    req.Header.Set("X-Custom-Header", "custom-value")
    
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Basic Authentication

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    client := &http.Client{}
    
    req, err := http.NewRequest("GET", "https://httpbin.org/basic-auth/user/pass", nil)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    // Basic Auth
    req.SetBasicAuth("user", "pass")
    
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Println("Status:", resp.Status)
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Bearer Token

```go
package main

import (
    "fmt"
    "io"
    "net/http"
)

func main() {
    client := &http.Client{}
    token := "my-secret-token"
    
    req, err := http.NewRequest("GET", "https://httpbin.org/bearer", nil)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    // Bearer Token
    req.Header.Set("Authorization", "Bearer "+token)
    
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### PUT/DELETE запросы

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

func main() {
    client := &http.Client{}
    
    // PUT запрос
    data := map[string]string{"name": "John", "email": "john@example.com"}
    jsonData, _ := json.Marshal(data)
    
    req, err := http.NewRequest("PUT", "https://httpbin.org/put", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Println("PUT Status:", resp.Status)
    
    // DELETE запрос
    req, _ = http.NewRequest("DELETE", "https://httpbin.org/delete", nil)
    resp, _ = client.Do(req)
    defer resp.Body.Close()
    
    fmt.Println("DELETE Status:", resp.Status)
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Запрос с контекстом

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
)

func main() {
    // Контекст с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    req, err := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/2", nil)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Управление редиректами

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Отключение редиректов
    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    }
    
    resp, err := client.Get("https://httpbin.org/redirect/3")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Println("Status:", resp.Status)
    fmt.Println("Location:", resp.Header.Get("Location"))
    
    // Ограничение числа редиректов
    client2 := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 2 {
                return fmt.Errorf("too many redirects")
            }
            return nil
        },
    }
    
    resp2, _ := client2.Get("https://httpbin.org/redirect/5")
    if resp2 != nil {
        fmt.Println("Status2:", resp2.Status)
    }
}
```

### Работа с cookies

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "net/http/cookiejar"
)

func main() {
    // Создаём cookie jar
    jar, _ := cookiejar.New(nil)
    
    client := &http.Client{
        Jar: jar,
    }
    
    // Первый запрос — сервер устанавливает cookie
    resp, _ := client.Get("https://httpbin.org/cookies/set/session/abc123")
    resp.Body.Close()
    
    // Второй запрос — cookie отправляется автоматически
    resp, _ = client.Get("https://httpbin.org/cookies")
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Кастомный Transport

```go
package main

import (
    "fmt"
    "io"
    "net/http"
    "time"
)

func main() {
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  false,
    }
    
    client := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
    
    resp, err := client.Get("https://httpbin.org/get")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

### Реюзабельный API клиент

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type APIClient struct {
    BaseURL    string
    HTTPClient *http.Client
    Token      string
}

func NewAPIClient(baseURL, token string) *APIClient {
    return &APIClient{
        BaseURL: baseURL,
        Token:   token,
        HTTPClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *APIClient) doRequest(method, path string, body interface{}) (*http.Response, error) {
    var bodyReader io.Reader
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        bodyReader = bytes.NewBuffer(jsonData)
    }
    
    req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    
    return c.HTTPClient.Do(req)
}

func (c *APIClient) Get(path string) (*http.Response, error) {
    return c.doRequest("GET", path, nil)
}

func (c *APIClient) Post(path string, body interface{}) (*http.Response, error) {
    return c.doRequest("POST", path, body)
}

func main() {
    client := NewAPIClient("https://httpbin.org", "my-token")
    
    resp, err := client.Get("/get")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
}
```

---

## ⚠️ Частые ошибки

### 1. Создание клиента в цикле

```go
// ❌ Создаёт новый клиент каждую итерацию
for _, url := range urls {
    client := &http.Client{}
    client.Get(url)
}

// ✅ Переиспользуйте клиент
client := &http.Client{}
for _, url := range urls {
    client.Get(url)
}
```

### 2. Без таймаута

```go
// ❌ Может зависнуть навсегда
client := &http.Client{}

// ✅ Всегда устанавливайте таймаут
client := &http.Client{
    Timeout: 30 * time.Second,
}
```

### 3. Игнорирование Body при ошибках

```go
// ❌ Body не закрывается
resp, err := client.Get(url)
if resp.StatusCode != 200 {
    return errors.New("bad status")
}

// ✅ Всегда закрывайте Body
resp, err := client.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()

if resp.StatusCode != 200 {
    return errors.New("bad status")
}
```

---

## 📝 Практика

### Задача 1: REST API wrapper
Создайте обёртку для REST API с CRUD методами.

### Задача 2: OAuth client
Клиент с OAuth авторизацией.

### Задача 3: Retry middleware
Middleware для повторных попыток.

### Задача 4: Rate limited client
Клиент с ограничением запросов.

### Задача 5: Circuit breaker
Паттерн circuit breaker для HTTP.

### Задача 6: Logging transport
Transport с логированием запросов.

### Задача 7: Caching client
Клиент с кэшированием ответов.

### Задача 8: Metrics collector
Сбор метрик HTTP запросов.
