# Отправка запросов (net.Dial)

## 💡 Ключевые идеи

1. **net.Dial** — функция для создания сетевых соединений
2. **Протоколы** — поддержка TCP, UDP, IP, Unix сокетов
3. **net.Conn** — интерфейс соединения (Reader + Writer)
4. **Адресация** — IP:порт или домен:порт
5. **Двунаправленность** — можно читать и писать в одно соединение

---

## 📋 Синтаксис

### Функция Dial

```go
func net.Dial(network, address string) (net.Conn, error)
```

### Поддерживаемые протоколы

| Протокол | Описание |
|----------|----------|
| `tcp`, `tcp4`, `tcp6` | TCP соединение (IPv4/IPv6) |
| `udp`, `udp4`, `udp6` | UDP соединение (IPv4/IPv6) |
| `ip`, `ip4`, `ip6` | IP соединение |
| `unix`, `unixgram`, `unixpacket` | Unix сокеты |

### Форматы адресов

```go
"127.0.0.1:80"           // IPv4 с портом
"localhost:8080"         // домен с портом
"[::1]:80"               // IPv6 с портом
"example.com:443"        // внешний домен
```

### Интерфейс net.Conn

```go
type Conn interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
    LocalAddr() Addr
    RemoteAddr() Addr
    SetDeadline(t time.Time) error
    SetReadDeadline(t time.Time) error
    SetWriteDeadline(t time.Time) error
}
```

---

## 💻 Примеры кода

### Простой TCP запрос

```go
package main

import (
    "fmt"
    "io"
    "net"
    "os"
)

func main() {
    // Подключаемся к серверу
    conn, err := net.Dial("tcp", "golang.org:80")
    if err != nil {
        fmt.Println("Connection error:", err)
        return
    }
    defer conn.Close()
    
    // Формируем HTTP запрос вручную
    request := "GET / HTTP/1.1\r\n" +
        "Host: golang.org\r\n" +
        "Connection: close\r\n\r\n"
    
    // Отправляем запрос
    _, err = conn.Write([]byte(request))
    if err != nil {
        fmt.Println("Write error:", err)
        return
    }
    
    // Читаем ответ и выводим на консоль
    io.Copy(os.Stdout, conn)
}
```

### Чтение с буфером

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    conn, err := net.Dial("tcp", "example.com:80")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()
    
    request := "GET / HTTP/1.1\r\nHost: example.com\r\nConnection: close\r\n\r\n"
    conn.Write([]byte(request))
    
    // Читаем порциями
    buffer := make([]byte, 1024)
    for {
        n, err := conn.Read(buffer)
        if n > 0 {
            fmt.Print(string(buffer[:n]))
        }
        if err != nil {
            break  // EOF или ошибка
        }
    }
}
```

### Отправка и получение данных

```go
package main

import (
    "bufio"
    "fmt"
    "net"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:4545")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()
    
    // Отправляем сообщение
    message := "Hello, Server!"
    conn.Write([]byte(message))
    
    // Получаем ответ
    reader := bufio.NewReader(conn)
    response, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Read error:", err)
        return
    }
    
    fmt.Println("Response:", response)
}
```

### UDP соединение

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // UDP соединение (без установки "рукопожатия")
    conn, err := net.Dial("udp", "127.0.0.1:5000")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()
    
    // Отправка данных
    message := []byte("UDP message")
    _, err = conn.Write(message)
    if err != nil {
        fmt.Println("Write error:", err)
        return
    }
    
    // Получение ответа (если сервер отвечает)
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Read error:", err)
        return
    }
    
    fmt.Println("Response:", string(buffer[:n]))
}
```

### Получение информации о соединении

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    conn, err := net.Dial("tcp", "google.com:80")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()
    
    // Локальный адрес (наш)
    fmt.Println("Local:", conn.LocalAddr().String())
    fmt.Println("Local Network:", conn.LocalAddr().Network())
    
    // Удалённый адрес (сервер)
    fmt.Println("Remote:", conn.RemoteAddr().String())
    fmt.Println("Remote Network:", conn.RemoteAddr().Network())
}
```

### Dial с контекстом и отменой

```go
package main

import (
    "context"
    "fmt"
    "net"
    "time"
)

func main() {
    // Контекст с таймаутом
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Dialer с настройками
    dialer := net.Dialer{
        Timeout:   3 * time.Second,
        KeepAlive: 30 * time.Second,
    }
    
    conn, err := dialer.DialContext(ctx, "tcp", "example.com:80")
    if err != nil {
        fmt.Println("Connection error:", err)
        return
    }
    defer conn.Close()
    
    fmt.Println("Connected to:", conn.RemoteAddr())
}
```

### Резолвинг DNS

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Получить IP адреса домена
    ips, err := net.LookupIP("google.com")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Println("IP addresses for google.com:")
    for _, ip := range ips {
        fmt.Println(" ", ip)
    }
    
    // Получить записи MX
    mxRecords, err := net.LookupMX("google.com")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Println("\nMX records:")
    for _, mx := range mxRecords {
        fmt.Printf("  %s (priority %d)\n", mx.Host, mx.Pref)
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Забыли закрыть соединение

```go
// ❌ Утечка ресурсов
conn, _ := net.Dial("tcp", "example.com:80")
conn.Write([]byte("data"))
// соединение не закрыто!

// ✅ Всегда закрывайте соединение
conn, err := net.Dial("tcp", "example.com:80")
if err != nil {
    return
}
defer conn.Close()
```

### 2. Не проверяют ошибку Dial

```go
// ❌ Паника при использовании nil conn
conn, _ := net.Dial("tcp", "nonexistent:80")
conn.Write([]byte("data"))  // panic!

// ✅ Всегда проверяйте ошибку
conn, err := net.Dial("tcp", "example.com:80")
if err != nil {
    fmt.Println("Failed to connect:", err)
    return
}
```

### 3. Неверный формат IPv6

```go
// ❌ Неверный формат
conn, _ := net.Dial("tcp", "::1:80")

// ✅ IPv6 адрес в квадратных скобках
conn, _ := net.Dial("tcp", "[::1]:80")
```

### 4. Блокировка на Read без таймаута

```go
// ❌ Может заблокироваться навсегда
n, err := conn.Read(buffer)

// ✅ Установите таймаут
conn.SetReadDeadline(time.Now().Add(10 * time.Second))
n, err := conn.Read(buffer)
```

---

## 📝 Практика

### Задача 1: Ping
Проверьте доступность хоста через TCP соединение.

### Задача 2: Port scanner
Просканируйте открытые порты на хосте.

### Задача 3: HTTP client
Реализуйте простой HTTP клиент без net/http.

### Задача 4: DNS lookup
Создайте утилиту для DNS запросов.

### Задача 5: Echo client
Клиент для echo-сервера.

### Задача 6: File downloader
Скачайте файл по TCP.

### Задача 7: Health checker
Проверка доступности нескольких серверов.

### Задача 8: Load tester
Простой инструмент нагрузочного тестирования.
