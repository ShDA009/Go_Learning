# SQLite

## 💡 Ключевые идеи

1. **go-sqlite3** — CGO драйвер для SQLite
2. **Файловая БД** — база данных в одном файле
3. **Без сервера** — не требует установки СУБД
4. **Плейсхолдеры** — `$1, $2` или `?`
5. **Идеально для** — тестов, прототипов, встраиваемых приложений

---

## 📋 Синтаксис

### Установка драйвера

```bash
go get -u github.com/mattn/go-sqlite3
```

> ⚠️ Требует CGO и компилятор C (gcc)

### Открытие базы данных

```go
// Файловая база
db, err := sql.Open("sqlite3", "database.db")

// In-memory база (для тестов)
db, err := sql.Open("sqlite3", ":memory:")

// С параметрами
db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=rwc")
```

### Параметры подключения

| Параметр | Описание |
|----------|----------|
| `cache=shared` | Разделяемый кэш |
| `mode=rwc` | Read-write-create |
| `_busy_timeout=5000` | Таймаут блокировки (мс) |
| `_journal=WAL` | Write-Ahead Logging |
| `_foreign_keys=on` | Включить внешние ключи |

---

## 💻 Примеры кода

### Создание базы и таблицы

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Создаём/открываем файл базы данных
    db, err := sql.Open("sqlite3", "store.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Создаём таблицу
    createTable := `
        CREATE TABLE IF NOT EXISTS products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            model TEXT NOT NULL,
            company TEXT NOT NULL,
            price INTEGER NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `
    
    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Database and table created!")
}
```

### INSERT — добавление данных

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, _ := sql.Open("sqlite3", "store.db")
    defer db.Close()
    
    // Вставка с параметрами ($1, $2 или ?)
    result, err := db.Exec(
        "INSERT INTO products (model, company, price) VALUES ($1, $2, $3)",
        "iPhone X", "Apple", 72000,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // SQLite поддерживает LastInsertId
    id, _ := result.LastInsertId()
    fmt.Println("Inserted ID:", id)
    
    rows, _ := result.RowsAffected()
    fmt.Println("Rows affected:", rows)
}
```

### Множественная вставка

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

type Product struct {
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("sqlite3", "store.db")
    defer db.Close()
    
    products := []Product{
        {"Pixel 7", "Google", 65000},
        {"Galaxy S23", "Samsung", 70000},
        {"Mi 13", "Xiaomi", 45000},
    }
    
    // Prepared statement
    stmt, err := db.Prepare("INSERT INTO products (model, company, price) VALUES (?, ?, ?)")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    
    for _, p := range products {
        result, err := stmt.Exec(p.Model, p.Company, p.Price)
        if err != nil {
            log.Println("Insert error:", err)
            continue
        }
        id, _ := result.LastInsertId()
        fmt.Printf("Inserted %s with ID %d\n", p.Model, id)
    }
}
```

### SELECT — получение данных

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("sqlite3", "store.db")
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

### SELECT с условием

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, _ := sql.Open("sqlite3", "store.db")
    defer db.Close()
    
    // Товары дороже 60000
    rows, err := db.Query(
        "SELECT id, model, price FROM products WHERE price > ?",
        60000,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    fmt.Println("Products over 60000:")
    for rows.Next() {
        var id, price int
        var model string
        rows.Scan(&id, &model, &price)
        fmt.Printf("  %d: %s - %d\n", id, model, price)
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
    
    _ "github.com/mattn/go-sqlite3"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("sqlite3", "store.db")
    defer db.Close()
    
    var p Product
    
    err := db.QueryRow(
        "SELECT id, model, company, price FROM products WHERE id = ?",
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

### UPDATE — обновление

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, _ := sql.Open("sqlite3", "store.db")
    defer db.Close()
    
    result, err := db.Exec(
        "UPDATE products SET price = ? WHERE id = ?",
        69000, 1,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    rows, _ := result.RowsAffected()
    fmt.Printf("Updated %d row(s)\n", rows)
}
```

### DELETE — удаление

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, _ := sql.Open("sqlite3", "store.db")
    defer db.Close()
    
    result, err := db.Exec("DELETE FROM products WHERE id = ?", 5)
    if err != nil {
        log.Fatal(err)
    }
    
    rows, _ := result.RowsAffected()
    fmt.Printf("Deleted %d row(s)\n", rows)
}
```

### In-memory база для тестов

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // База в памяти (удаляется при закрытии)
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Создаём таблицу
    db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY,
            name TEXT
        )
    `)
    
    // Добавляем данные
    db.Exec("INSERT INTO users (name) VALUES (?)", "Alice")
    db.Exec("INSERT INTO users (name) VALUES (?)", "Bob")
    
    // Читаем
    rows, _ := db.Query("SELECT id, name FROM users")
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        rows.Scan(&id, &name)
        fmt.Printf("%d: %s\n", id, name)
    }
}
```

### Транзакции

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, _ := sql.Open("sqlite3", "store.db")
    defer db.Close()
    
    // Начинаем транзакцию
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }
    
    // Выполняем операции
    _, err = tx.Exec("INSERT INTO products (model, company, price) VALUES (?, ?, ?)",
        "Test Product", "Test", 1000)
    if err != nil {
        tx.Rollback()
        log.Fatal("Insert failed:", err)
    }
    
    _, err = tx.Exec("UPDATE products SET price = price + 100 WHERE company = ?", "Test")
    if err != nil {
        tx.Rollback()
        log.Fatal("Update failed:", err)
    }
    
    // Фиксируем
    err = tx.Commit()
    if err != nil {
        log.Fatal("Commit failed:", err)
    }
    
    fmt.Println("Transaction completed!")
}
```

### WAL режим для производительности

```go
package main

import (
    "database/sql"
    "fmt"
    
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Включаем WAL режим через параметры
    db, _ := sql.Open("sqlite3", "store.db?_journal=WAL&_busy_timeout=5000")
    defer db.Close()
    
    // Или через PRAGMA
    db.Exec("PRAGMA journal_mode=WAL")
    db.Exec("PRAGMA busy_timeout=5000")
    
    // Проверяем режим
    var journalMode string
    db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
    fmt.Println("Journal mode:", journalMode)
}
```

### Внешние ключи

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Включаем внешние ключи
    db, _ := sql.Open("sqlite3", "store.db?_foreign_keys=on")
    defer db.Close()
    
    // Создаём таблицы с внешними ключами
    db.Exec(`
        CREATE TABLE IF NOT EXISTS categories (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL
        )
    `)
    
    db.Exec(`
        CREATE TABLE IF NOT EXISTS items (
            id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            category_id INTEGER,
            FOREIGN KEY (category_id) REFERENCES categories(id)
        )
    `)
    
    // Добавляем категорию
    db.Exec("INSERT INTO categories (name) VALUES (?)", "Electronics")
    
    // Добавляем товар с существующей категорией — ОК
    _, err := db.Exec("INSERT INTO items (name, category_id) VALUES (?, ?)", "Phone", 1)
    if err != nil {
        log.Println("Error:", err)
    } else {
        fmt.Println("Item added successfully")
    }
    
    // Попытка добавить с несуществующей категорией — ошибка
    _, err = db.Exec("INSERT INTO items (name, category_id) VALUES (?, ?)", "Tablet", 999)
    if err != nil {
        fmt.Println("Foreign key constraint:", err)
    }
}
```

---

## ⚠️ Частые ошибки

### 1. CGO не включён

```bash
# Ошибка: Binary was compiled with 'CGO_ENABLED=0'
# Решение: включите CGO
CGO_ENABLED=1 go build
```

### 2. Database is locked

```go
// ❌ Множественные параллельные записи
// Error: database is locked

// ✅ Используйте WAL и busy_timeout
db, _ := sql.Open("sqlite3", "db.sqlite?_journal=WAL&_busy_timeout=5000")
```

### 3. Путь к файлу

```go
// ❌ Относительный путь может не работать
db, _ := sql.Open("sqlite3", "data/store.db")

// ✅ Используйте абсолютный путь
db, _ := sql.Open("sqlite3", "/app/data/store.db")
```

---

## 📝 Практика

### Задача 1: Todo app
Создайте приложение для задач с SQLite.

### Задача 2: Test fixtures
In-memory база для unit-тестов.

### Задача 3: Config store
Хранение конфигурации в SQLite.

### Задача 4: Cache layer
Кэширование данных в SQLite.

### Задача 5: Backup/Restore
Резервное копирование базы.

### Задача 6: FTS search
Full-Text Search в SQLite.

### Задача 7: Migrations
Система миграций для SQLite.

### Задача 8: Analytics
Хранение и анализ логов.
