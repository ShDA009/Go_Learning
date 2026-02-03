# TCP Сервер (net.Listen)

## 💡 Ключевые идеи

1. **net.Listen** — создание слушателя для входящих соединений
2. **Accept()** — приём входящего подключения
3. **Горутины** — параллельная обработка клиентов
4. **net.Listener** — интерфейс для приёма соединений
5. **Graceful shutdown** — корректное завершение сервера

---

## 📋 Синтаксис

### Функция Listen

```go
func net.Listen(network, address string) (net.Listener, error)
```

### Интерфейс Listener

```go
type Listener interface {
    Accept() (Conn, error)  // принять соединение
    Close() error           // закрыть слушатель
    Addr() Addr             // адрес слушателя
}
```

### Поддерживаемые протоколы

```go
"tcp", "tcp4", "tcp6"    // TCP соединения
"unix", "unixpacket"     // Unix сокеты
```

---

## 📖 Теория

### Клиент vs Сервер

**Клиент** — инициирует соединение (`net.Dial`)
**Сервер** — ждёт и принимает соединения (`net.Listen` + `Accept`)

### Как работает TCP сервер?

1. **Listen** — начинаем слушать порт
2. **Accept** — принимаем подключение
3. **Handle** — обрабатываем клиента
4. **Повторяем** с пункта 2

```go
listener, err := net.Listen("tcp", ":8080")
if err != nil {
    log.Fatal(err)
}
defer listener.Close()

for {
    conn, err := listener.Accept()  // блокируется до подключения
    if err != nil {
        continue
    }
    go handleConnection(conn)  // обработка в горутине
}
```

### Порт и адрес

```go
// Только порт — слушать на всех интерфейсах
net.Listen("tcp", ":8080")

// Конкретный IP
net.Listen("tcp", "127.0.0.1:8080")  // только localhost
net.Listen("tcp", "0.0.0.0:8080")    // все интерфейсы

// Порт 0 — система выберет свободный
listener, _ := net.Listen("tcp", ":0")
fmt.Println(listener.Addr())  // выведет реальный порт
```

### Почему горутины?

Без горутин сервер обрабатывает **одного** клиента за раз:

```go
// ПЛОХО — последовательная обработка
for {
    conn, _ := listener.Accept()
    handleConnection(conn)  // другие клиенты ждут!
}

// ХОРОШО — параллельная обработка
for {
    conn, _ := listener.Accept()
    go handleConnection(conn)  // каждый клиент в своей горутине
}
```

### Обработка клиента

```go
func handleConnection(conn net.Conn) {
    defer conn.Close()  // обязательно закрыть!
    
    // Читаем запрос
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        return
    }
    request := buf[:n]
    
    // Отправляем ответ
    conn.Write([]byte("Hello from server!"))
}
```

### Graceful Shutdown

Корректное завершение сервера:

```go
func main() {
    listener, _ := net.Listen("tcp", ":8080")
    
    // Ловим сигнал завершения
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigCh
        listener.Close()  // Accept вернёт ошибку
    }()
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            return  // сервер остановлен
        }
        go handle(conn)
    }
}
```

---

### Форматы адресов

```go
":8080"              // все интерфейсы, порт 8080
"localhost:8080"     // только localhost
"192.168.1.1:8080"   // конкретный IP
"0.0.0.0:8080"       // все IPv4 интерфейсы
"[::]:8080"          // все IPv6 интерфейсы
```

---

## 💻 Примеры кода

### Простейший TCP сервер

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Создаём слушатель на порту 4545
    listener, err := net.Listen("tcp", ":4545")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("Server started on :4545")
    
    for {
        // Принимаем входящее соединение
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Accept error:", err)
            continue
        }
        
        // Отправляем приветствие и закрываем
        conn.Write([]byte("Hello from server!\n"))
        conn.Close()
    }
}
```

### Сервер с горутинами

```go
package main

import (
    "fmt"
    "net"
    "time"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    clientAddr := conn.RemoteAddr().String()
    fmt.Println("Client connected:", clientAddr)
    
    // Симулируем обработку
    time.Sleep(2 * time.Second)
    
    message := fmt.Sprintf("Hello, %s! Time: %s\n", 
        clientAddr, time.Now().Format(time.RFC3339))
    conn.Write([]byte(message))
    
    fmt.Println("Client disconnected:", clientAddr)
}

func main() {
    listener, err := net.Listen("tcp", ":4545")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("Server started on :4545")
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Accept error:", err)
            continue
        }
        
        // Обрабатываем каждое соединение в горутине
        go handleConnection(conn)
    }
}
```

### Echo сервер

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

func handleClient(conn net.Conn) {
    defer conn.Close()
    
    reader := bufio.NewReader(conn)
    
    for {
        // Читаем строку от клиента
        message, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Client disconnected")
            return
        }
        
        message = strings.TrimSpace(message)
        fmt.Printf("Received: %s\n", message)
        
        // Отправляем эхо обратно
        response := fmt.Sprintf("Echo: %s\n", message)
        conn.Write([]byte(response))
        
        // Команда выхода
        if message == "quit" {
            return
        }
    }
}

func main() {
    listener, err := net.Listen("tcp", ":4545")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("Echo server on :4545")
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleClient(conn)
    }
}
```

### Сервер-словарь (пример из оригинала)

```go
package main

import (
    "fmt"
    "net"
)

var dictionary = map[string]string{
    "red":    "красный",
    "green":  "зелёный",
    "blue":   "синий",
    "yellow": "жёлтый",
    "black":  "чёрный",
    "white":  "белый",
}

func handleClient(conn net.Conn) {
    defer conn.Close()
    
    buffer := make([]byte, 1024)
    
    for {
        n, err := conn.Read(buffer)
        if n == 0 || err != nil {
            fmt.Println("Client disconnected")
            return
        }
        
        word := string(buffer[:n])
        
        // Ищем перевод
        translation, found := dictionary[word]
        if !found {
            translation = "не найдено"
        }
        
        fmt.Printf("%s -> %s\n", word, translation)
        conn.Write([]byte(translation))
    }
}

func main() {
    listener, err := net.Listen("tcp", ":4545")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("Dictionary server on :4545")
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleClient(conn)
    }
}
```

### Сервер с лимитом подключений

```go
package main

import (
    "fmt"
    "net"
    "sync"
)

const maxConnections = 10

func main() {
    listener, err := net.Listen("tcp", ":4545")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer listener.Close()
    
    // Семафор для ограничения соединений
    semaphore := make(chan struct{}, maxConnections)
    var wg sync.WaitGroup
    
    fmt.Printf("Server started (max %d connections)\n", maxConnections)
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        
        // Занимаем слот
        semaphore <- struct{}{}
        wg.Add(1)
        
        go func(c net.Conn) {
            defer func() {
                <-semaphore  // освобождаем слот
                wg.Done()
                c.Close()
            }()
            
            // Обработка клиента
            c.Write([]byte("Hello!\n"))
        }(conn)
    }
}
```

### Graceful shutdown

```go
package main

import (
    "context"
    "fmt"
    "net"
    "os"
    "os/signal"
    "sync"
    "syscall"
)

func main() {
    listener, err := net.Listen("tcp", ":4545")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    ctx, cancel := context.WithCancel(context.Background())
    var wg sync.WaitGroup
    
    // Обработка сигналов
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        fmt.Println("\nShutting down...")
        cancel()
        listener.Close()
    }()
    
    fmt.Println("Server started on :4545")
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            select {
            case <-ctx.Done():
                wg.Wait()
                fmt.Println("Server stopped")
                return
            default:
                continue
            }
        }
        
        wg.Add(1)
        go func(c net.Conn) {
            defer wg.Done()
            defer c.Close()
            
            c.Write([]byte("Hello!\n"))
        }(conn)
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Не закрывают соединение клиента

```go
// ❌ Утечка соединений
func handle(conn net.Conn) {
    conn.Write([]byte("Hello"))
    // соединение не закрыто!
}

// ✅ Используйте defer
func handle(conn net.Conn) {
    defer conn.Close()
    conn.Write([]byte("Hello"))
}
```

### 2. Блокирующий Accept в main

```go
// ❌ Обработка по одному
for {
    conn, _ := listener.Accept()
    handleClient(conn)  // блокирует следующий Accept
}

// ✅ Горутины для параллельности
for {
    conn, _ := listener.Accept()
    go handleClient(conn)
}
```

### 3. Игнорирование ошибки Accept

```go
// ❌ Паника на nil conn
conn, _ := listener.Accept()
conn.Write(...)  // panic if err != nil

// ✅ Проверяйте ошибку
conn, err := listener.Accept()
if err != nil {
    continue
}
```

### 4. Порт уже занят

```go
// При повторном запуске сервера
// Error: listen tcp :4545: bind: address already in use

// Решение: проверьте, не запущен ли уже сервер
// Или используйте другой порт
```

---

## 🏋️ Практические задания

### Задание 1: gorilla/mux роутер

Создайте роутер с mux.

**Ожидаемый результат:**
```
Mux: r := mux.NewRouter()
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Mux: r := mux.NewRouter()")
}
```

**Ожидаемый вывод:**
```
Mux: r := mux.NewRouter()
```

**Баллы:** 10

### Задание 2: Параметры пути

Извлеките параметры из URL.

**Ожидаемый результат:**
```
Параметры: r.HandleFunc("/users/{id}", ...); mux.Vars(r)["id"]
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Параметры: r.HandleFunc(\"/users/{id}\", ...); mux.Vars(r)[\"id\"]")
}
```

**Ожидаемый вывод:**
```
Параметры: r.HandleFunc("/users/{id}", ...); mux.Vars(r)["id"]
```

**Баллы:** 15

### Задание 3: Методы HTTP

Ограничьте маршрут по методу.

**Ожидаемый результат:**
```
Метод: r.HandleFunc(...).Methods("GET", "POST")
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Метод: r.HandleFunc(...).Methods(\"GET\", \"POST\")")
}
```

**Ожидаемый вывод:**
```
Метод: r.HandleFunc(...).Methods("GET", "POST")
```

**Баллы:** 10

### Задание 4: Подроутер

Создайте подроутер.

**Ожидаемый результат:**
```
Подроутер: api := r.PathPrefix("/api").Subrouter()
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Подроутер: api := r.PathPrefix(\"/api\").Subrouter()")
}
```

**Ожидаемый вывод:**
```
Подроутер: api := r.PathPrefix("/api").Subrouter()
```

**Баллы:** 15

### Задание 5: Middleware в mux

Добавьте middleware.

**Ожидаемый результат:**
```
Middleware: r.Use(loggingMiddleware)
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Middleware: r.Use(loggingMiddleware)")
}
```

**Ожидаемый вывод:**
```
Middleware: r.Use(loggingMiddleware)
```

**Баллы:** 15
