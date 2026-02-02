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

## 📝 Практика

### Задача 1: Time server
Сервер, возвращающий текущее время.

### Задача 2: Calculator server
Сервер-калькулятор (принимает выражение, возвращает результат).

### Задача 3: Chat room
Простой чат-сервер для нескольких клиентов.

### Задача 4: File server
Сервер для передачи файлов.

### Задача 5: Key-value store
Простое key-value хранилище по сети.

### Задача 6: Command server
Сервер, выполняющий команды.

### Задача 7: Proxy server
Простой TCP прокси.

### Задача 8: Load balancer
Балансировщик нагрузки между серверами.
