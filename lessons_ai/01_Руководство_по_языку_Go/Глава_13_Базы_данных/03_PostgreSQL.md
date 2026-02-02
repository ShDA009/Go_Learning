# PostgreSQL

## 💡 Ключевые идеи

1. **lib/pq** — популярный драйвер для PostgreSQL
2. **Плейсхолдеры** — `$1, $2, $3` вместо `?`
3. **RETURNING** — получение данных после INSERT/UPDATE
4. **LastInsertId** — НЕ поддерживается (используйте RETURNING)
5. **SSL** — настройка через sslmode

---

## 📋 Синтаксис

### Установка драйвера

```bash
go get -u github.com/lib/pq
```

### Формат строки подключения

```go
// Key=value формат
"user=postgres password=secret dbname=mydb sslmode=disable"

// URL формат
"postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

### Параметры подключения

| Параметр | Описание |
|----------|----------|
| `user` | Имя пользователя |
| `password` | Пароль |
| `host` | Хост (по умолчанию localhost) |
| `port` | Порт (по умолчанию 5432) |
| `dbname` | Имя базы данных |
| `sslmode` | disable, require, verify-full |
| `connect_timeout` | Таймаут подключения в секундах |

---

## 💻 Примеры кода

### Подключение к PostgreSQL

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

func main() {
    connStr := "user=postgres password=secret dbname=productdb sslmode=disable"
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Open error:", err)
    }
    defer db.Close()
    
    // Проверяем подключение
    if err = db.Ping(); err != nil {
        log.Fatal("Ping error:", err)
    }
    
    fmt.Println("Connected to PostgreSQL!")
}
```

### Альтернативный формат URL

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

func main() {
    // URL формат
    connStr := "postgres://postgres:secret@localhost:5432/productdb?sslmode=disable"
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    fmt.Println("Connected!")
}
```

### Создание таблицы

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=productdb sslmode=disable")
    defer db.Close()
    
    createTable := `
        CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            model VARCHAR(100) NOT NULL,
            company VARCHAR(100) NOT NULL,
            price INTEGER NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `
    
    _, err := db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Table created!")
}
```

### INSERT с RETURNING (получение ID)

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=productdb sslmode=disable")
    defer db.Close()
    
    var id int
    
    // ВАЖНО: используем QueryRow + RETURNING вместо Exec
    err := db.QueryRow(
        "INSERT INTO products (model, company, price) VALUES ($1, $2, $3) RETURNING id",
        "iPhone X", "Apple", 72000,
    ).Scan(&id)
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Inserted ID:", id)
}
```

### INSERT без получения ID

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=productdb sslmode=disable")
    defer db.Close()
    
    result, err := db.Exec(
        "INSERT INTO products (model, company, price) VALUES ($1, $2, $3)",
        "Pixel 7", "Google", 65000,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // LastInsertId НЕ работает в PostgreSQL!
    // id, _ := result.LastInsertId()  // вернёт 0 и ошибку
    
    // RowsAffected работает
    rows, _ := result.RowsAffected()
    fmt.Printf("Rows affected: %d\n", rows)
}
```

### SELECT — получение данных

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=productdb sslmode=disable")
    defer db.Close()
    
    rows, err := db.Query("SELECT id, model, company, price FROM products")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    var products []Product
    
    for rows.Next() {
        var p Product
        err := rows.Scan(&p.ID, &p.Model, &p.Company, &p.Price)
        if err != nil {
            log.Println("Scan error:", err)
            continue
        }
        products = append(products, p)
    }
    
    for _, p := range products {
        fmt.Printf("%d: %s (%s) - %d руб.\n", p.ID, p.Model, p.Company, p.Price)
    }
}
```

### SELECT с параметрами

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=productdb sslmode=disable")
    defer db.Close()
    
    // Параметры: $1, $2, $3...
    rows, err := db.Query(
        "SELECT id, model, company, price FROM products WHERE price > $1 AND company = $2",
        50000, "Apple",
    )
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    fmt.Println("Apple products over 50000:")
    for rows.Next() {
        var p Product
        rows.Scan(&p.ID, &p.Model, &p.Company, &p.Price)
        fmt.Printf("  %s - %d\n", p.Model, p.Price)
    }
}
```

### QueryRow — одна запись

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=productdb sslmode=disable")
    defer db.Close()
    
    var p Product
    
    err := db.QueryRow(
        "SELECT id, model, company, price FROM products WHERE id = $1",
        1,
    ).Scan(&p.ID, &p.Model, &p.Company, &p.Price)
    
    if err == sql.ErrNoRows {
        fmt.Println("Product not found")
        return
    }
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found: %+v\n", p)
}
```

### UPDATE с RETURNING

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=productdb sslmode=disable")
    defer db.Close()
    
    var p Product
    
    // UPDATE с RETURNING возвращает обновлённые данные
    err := db.QueryRow(
        "UPDATE products SET price = $1 WHERE id = $2 RETURNING id, model, company, price",
        69000, 1,
    ).Scan(&p.ID, &p.Model, &p.Company, &p.Price)
    
    if err == sql.ErrNoRows {
        fmt.Println("Product not found")
        return
    }
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Updated: %+v\n", p)
}
```

### DELETE

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=productdb sslmode=disable")
    defer db.Close()
    
    result, err := db.Exec("DELETE FROM products WHERE id = $1", 5)
    if err != nil {
        log.Fatal(err)
    }
    
    rows, _ := result.RowsAffected()
    fmt.Printf("Deleted %d row(s)\n", rows)
}
```

### UPSERT (INSERT ON CONFLICT)

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=productdb sslmode=disable")
    defer db.Close()
    
    // Вставить или обновить если существует
    var id int
    err := db.QueryRow(`
        INSERT INTO products (model, company, price) 
        VALUES ($1, $2, $3)
        ON CONFLICT (model) 
        DO UPDATE SET price = EXCLUDED.price
        RETURNING id
    `, "iPhone X", "Apple", 75000).Scan(&id)
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Upserted ID:", id)
}
```

### Работа с массивами PostgreSQL

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=mydb sslmode=disable")
    defer db.Close()
    
    // Создаём таблицу с массивом
    db.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id SERIAL PRIMARY KEY,
            title TEXT,
            tags TEXT[]
        )
    `)
    
    // Вставляем массив
    tags := []string{"go", "programming", "tutorial"}
    _, err := db.Exec(
        "INSERT INTO posts (title, tags) VALUES ($1, $2)",
        "Learning Go", pq.Array(tags),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Читаем массив
    var title string
    var readTags []string
    
    err = db.QueryRow("SELECT title, tags FROM posts WHERE id = 1").
        Scan(&title, pq.Array(&readTags))
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Title:", title)
    fmt.Println("Tags:", readTags)
}
```

### JSONB в PostgreSQL

```go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    
    _ "github.com/lib/pq"
)

type Metadata struct {
    Color  string `json:"color"`
    Size   string `json:"size"`
    Weight int    `json:"weight"`
}

func main() {
    db, _ := sql.Open("postgres", "user=postgres password=secret dbname=mydb sslmode=disable")
    defer db.Close()
    
    // Таблица с JSONB
    db.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            id SERIAL PRIMARY KEY,
            name TEXT,
            metadata JSONB
        )
    `)
    
    // Вставляем JSON
    meta := Metadata{Color: "red", Size: "M", Weight: 100}
    metaJSON, _ := json.Marshal(meta)
    
    db.Exec("INSERT INTO items (name, metadata) VALUES ($1, $2)", "T-Shirt", metaJSON)
    
    // Читаем JSON
    var name string
    var metaBytes []byte
    
    db.QueryRow("SELECT name, metadata FROM items WHERE id = 1").
        Scan(&name, &metaBytes)
    
    var readMeta Metadata
    json.Unmarshal(metaBytes, &readMeta)
    
    fmt.Println("Name:", name)
    fmt.Printf("Metadata: %+v\n", readMeta)
}
```

---

## ⚠️ Частые ошибки

### 1. Использование ? вместо $n

```go
// ❌ MySQL синтаксис
db.Query("SELECT * FROM products WHERE id = ?", 1)

// ✅ PostgreSQL синтаксис
db.Query("SELECT * FROM products WHERE id = $1", 1)
```

### 2. Попытка использовать LastInsertId

```go
// ❌ Не работает в PostgreSQL
result, _ := db.Exec("INSERT INTO products ...")
id, _ := result.LastInsertId()  // всегда 0!

// ✅ Используйте RETURNING
var id int
db.QueryRow("INSERT INTO products ... RETURNING id").Scan(&id)
```

### 3. SSL ошибки

```go
// ❌ SSL включён по умолчанию
"user=postgres password=secret dbname=mydb"
// Error: SSL is not enabled on the server

// ✅ Отключите SSL для локальной разработки
"user=postgres password=secret dbname=mydb sslmode=disable"
```

---

## 📝 Практика

### Задача 1: Product CRUD
Полный CRUD с использованием RETURNING.

### Задача 2: Full-text search
Полнотекстовый поиск с tsvector.

### Задача 3: JSONB queries
Работа с JSONB полями.

### Задача 4: Array operations
Операции с массивами PostgreSQL.

### Задача 5: CTE queries
Использование WITH (Common Table Expressions).

### Задача 6: Window functions
Оконные функции (RANK, ROW_NUMBER).

### Задача 7: Partitioning
Работа с партиционированными таблицами.

### Задача 8: Listen/Notify
Реализация pub/sub через LISTEN/NOTIFY.
