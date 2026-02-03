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

---

## 📖 Теория

### Что такое сокеты и TCP?

**Сокет** — это конечная точка сетевого соединения. Представьте телефон: чтобы позвонить, нужен номер (адрес) и связь (соединение).

**TCP** (Transmission Control Protocol) — надёжный протокол:
- Гарантирует доставку данных
- Сохраняет порядок данных
- Устанавливает соединение перед передачей

### net.Dial — установка соединения

```go
conn, err := net.Dial("tcp", "google.com:80")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()
```

Это аналог "позвонить по телефону":
1. `"tcp"` — какой протокол использовать
2. `"google.com:80"` — куда звоним (адрес:порт)
3. `conn` — установленное соединение (можно говорить)

### net.Conn — соединение

`net.Conn` реализует:
- `io.Reader` — можно читать данные
- `io.Writer` — можно писать данные
- `io.Closer` — нужно закрывать

```go
// Отправить данные
conn.Write([]byte("Hello, server!"))

// Получить ответ
buf := make([]byte, 1024)
n, err := conn.Read(buf)
response := buf[:n]
```

### Пример: HTTP запрос вручную

```go
conn, _ := net.Dial("tcp", "example.com:80")
defer conn.Close()

// Отправляем HTTP запрос
request := "GET / HTTP/1.0\r\nHost: example.com\r\n\r\n"
conn.Write([]byte(request))

// Читаем ответ
response, _ := io.ReadAll(conn)
fmt.Println(string(response))
```

### TCP vs UDP

| TCP | UDP |
|-----|-----|
| Надёжный | Ненадёжный |
| Устанавливает соединение | Без соединения |
| Медленнее | Быстрее |
| HTTP, SSH, FTP | DNS, игры, видео |

```go
// TCP
conn, _ := net.Dial("tcp", "server:80")

// UDP  
conn, _ := net.Dial("udp", "dns:53")
```

### Важные моменты

1. **Всегда закрывайте соединение:**
   ```go
   defer conn.Close()
   ```

2. **Обрабатывайте ошибки:**
   ```go
   if err != nil {
       // сервер недоступен, сеть упала, etc.
   }
   ```

3. **Чтение может быть частичным:**
   ```go
   // Данные приходят порциями!
   for {
       n, err := conn.Read(buf)
       if err == io.EOF {
           break
       }
       process(buf[:n])
   }
   ```

---

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

## 🏋️ Практические задания

### Задание 1: gorilla/sessions

Используйте сессии.

**Ожидаемый результат:**
```
Сессии: store.Get(r, "session-name")
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Сессии: store.Get(r, \"session-name\")")
}
```

**Ожидаемый вывод:**
```
Сессии: store.Get(r, "session-name")
```

**Баллы:** 15

### Задание 2: Session store

Создайте хранилище сессий.

**Ожидаемый результат:**
```
Store: sessions.NewCookieStore([]byte("secret"))
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Store: sessions.NewCookieStore([]byte(\"secret\"))")
}
```

**Ожидаемый вывод:**
```
Store: sessions.NewCookieStore([]byte("secret"))
```

**Баллы:** 10

### Задание 3: Запись в сессию

Сохраните данные в сессию.

**Ожидаемый результат:**
```
Запись: session.Values["key"] = value; session.Save(r, w)
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Запись: session.Values[\"key\"] = value; session.Save(r, w)")
}
```

**Ожидаемый вывод:**
```
Запись: session.Values["key"] = value; session.Save(r, w)
```

**Баллы:** 15

### Задание 4: Чтение из сессии

Получите данные из сессии.

**Ожидаемый результат:**
```
Чтение: val, ok := session.Values["key"]
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Чтение: val, ok := session.Values[\"key\"]")
}
```

**Ожидаемый вывод:**
```
Чтение: val, ok := session.Values["key"]
```

**Баллы:** 10

### Задание 5: Flash сообщения

Используйте flash сообщения.

**Ожидаемый результат:**
```
Flash: session.AddFlash("message")
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Flash: session.AddFlash(\"message\")")
}
```

**Ожидаемый вывод:**
```
Flash: session.AddFlash("message")
```

**Баллы:** 15
