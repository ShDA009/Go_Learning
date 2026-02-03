# MySQL

## 💡 Ключевые идеи

1. **go-sql-driver/mysql** — популярный драйвер для MySQL
2. **DSN формат** — user:password@tcp(host:port)/dbname
3. **Плейсхолдеры** — знаки `?` для параметров
4. **LastInsertId** — поддерживается для AUTO_INCREMENT
5. **Charset** — важно указать utf8mb4 для Unicode

---

## 📋 Синтаксис

### Установка драйвера

```bash
go get -u github.com/go-sql-driver/mysql
```

### Формат DSN (Data Source Name)

```go
// Базовый формат
"user:password@tcp(host:port)/dbname"

// С параметрами
"user:password@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=true"

// Локальный сокет (Unix)
"user:password@unix(/var/run/mysqld/mysqld.sock)/dbname"
```

### Параметры подключения

| Параметр | Описание |
|----------|----------|
| `charset` | Кодировка (utf8mb4) |
| `parseTime` | Парсить TIME/DATE в time.Time |
| `loc` | Временная зона (Local, UTC) |
| `timeout` | Таймаут подключения |
| `readTimeout` | Таймаут чтения |
| `writeTimeout` | Таймаут записи |

---

## 💻 Примеры кода

### Подключение к MySQL

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // DSN с полезными параметрами
    dsn := "root:password@tcp(localhost:3306)/productdb?charset=utf8mb4&parseTime=true"
    
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Open error:", err)
    }
    defer db.Close()
    
    // Проверяем подключение
    if err = db.Ping(); err != nil {
        log.Fatal("Ping error:", err)
    }
    
    fmt.Println("Connected to MySQL!")
}
```

### Создание таблицы

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/productdb")
    defer db.Close()
    
    createTable := `
        CREATE TABLE IF NOT EXISTS products (
            id INT AUTO_INCREMENT PRIMARY KEY,
            model VARCHAR(100) NOT NULL,
            company VARCHAR(100) NOT NULL,
            price INT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
    `
    
    _, err := db.Exec(createTable)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Table created!")
}
```

### INSERT — добавление данных

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/productdb")
    defer db.Close()
    
    // Добавление с параметрами
    result, err := db.Exec(
        "INSERT INTO products (model, company, price) VALUES (?, ?, ?)",
        "iPhone X", "Apple", 72000,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // ID добавленной записи
    id, _ := result.LastInsertId()
    fmt.Println("Inserted ID:", id)
    
    // Количество затронутых строк
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
    
    _ "github.com/go-sql-driver/mysql"
)

type Product struct {
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/productdb")
    defer db.Close()
    
    products := []Product{
        {"Pixel 7", "Google", 65000},
        {"Galaxy S23", "Samsung", 70000},
        {"Xperia 5", "Sony", 55000},
    }
    
    // Prepared statement для эффективности
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

### SELECT — получение всех записей

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/productdb")
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
    
    if err = rows.Err(); err != nil {
        log.Fatal(err)
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
    
    _ "github.com/go-sql-driver/mysql"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/productdb")
    defer db.Close()
    
    // Товары дороже 60000
    rows, err := db.Query(
        "SELECT id, model, company, price FROM products WHERE price > ?",
        60000,
    )
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    fmt.Println("Products over 60000:")
    for rows.Next() {
        var p Product
        rows.Scan(&p.ID, &p.Model, &p.Company, &p.Price)
        fmt.Printf("  %s - %d руб.\n", p.Model, p.Price)
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
    
    _ "github.com/go-sql-driver/mysql"
)

type Product struct {
    ID      int
    Model   string
    Company string
    Price   int
}

func getProductByID(db *sql.DB, id int) (*Product, error) {
    var p Product
    
    err := db.QueryRow(
        "SELECT id, model, company, price FROM products WHERE id = ?",
        id,
    ).Scan(&p.ID, &p.Model, &p.Company, &p.Price)
    
    if err == sql.ErrNoRows {
        return nil, nil  // не найдено
    }
    if err != nil {
        return nil, err
    }
    
    return &p, nil
}

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/productdb")
    defer db.Close()
    
    product, err := getProductByID(db, 1)
    if err != nil {
        log.Fatal(err)
    }
    
    if product == nil {
        fmt.Println("Product not found")
    } else {
        fmt.Printf("Found: %+v\n", product)
    }
}
```

### UPDATE — обновление

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/productdb")
    defer db.Close()
    
    // Обновляем цену для id=1
    result, err := db.Exec(
        "UPDATE products SET price = ? WHERE id = ?",
        69000, 1,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        fmt.Println("No rows updated (product not found?)")
    } else {
        fmt.Printf("Updated %d row(s)\n", rows)
    }
}
```

### DELETE — удаление

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/productdb")
    defer db.Close()
    
    result, err := db.Exec("DELETE FROM products WHERE id = ?", 5)
    if err != nil {
        log.Fatal(err)
    }
    
    rows, _ := result.RowsAffected()
    fmt.Printf("Deleted %d row(s)\n", rows)
}
```

### Транзакции

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/go-sql-driver/mysql"
)

func transferMoney(db *sql.DB, fromID, toID int, amount int) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    // Списание
    _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("debit failed: %w", err)
    }
    
    // Зачисление
    _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("credit failed: %w", err)
    }
    
    return tx.Commit()
}

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/bank")
    defer db.Close()
    
    err := transferMoney(db, 1, 2, 1000)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Transfer completed!")
}
```

### Работа с NULL значениями

```go
package main

import (
    "database/sql"
    "fmt"
    
    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID    int
    Name  string
    Email sql.NullString  // может быть NULL
    Age   sql.NullInt64   // может быть NULL
}

func main() {
    db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/mydb")
    defer db.Close()
    
    var user User
    db.QueryRow("SELECT id, name, email, age FROM users WHERE id = ?", 1).
        Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    
    fmt.Println("Name:", user.Name)
    
    if user.Email.Valid {
        fmt.Println("Email:", user.Email.String)
    } else {
        fmt.Println("Email: not set")
    }
    
    if user.Age.Valid {
        fmt.Println("Age:", user.Age.Int64)
    } else {
        fmt.Println("Age: not set")
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Неверный формат DSN

```go
// ❌ Неверно
"root:password@localhost/dbname"

// ✅ Правильно (tcp протокол обязателен)
"root:password@tcp(localhost:3306)/dbname"
```

### 2. Проблемы с кодировкой

```go
// ❌ Не указана кодировка — проблемы с кириллицей
"root:password@tcp(localhost:3306)/dbname"

// ✅ Указана utf8mb4
"root:password@tcp(localhost:3306)/dbname?charset=utf8mb4"
```

### 3. Проблемы с датами

```go
// ❌ time.Time не парсится
var createdAt time.Time
rows.Scan(&createdAt)  // error!

// ✅ Добавьте parseTime=true в DSN
dsn := "user:pass@tcp(host)/db?parseTime=true"
```

---

## 🏋️ Практические задания

### Задание 1: JWT структура

Понимание JWT токена.

**Ожидаемый результат:**
```
JWT: header.payload.signature
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("JWT: header.payload.signature")
}
```

**Ожидаемый вывод:**
```
JWT: header.payload.signature
```

**Баллы:** 10

### Задание 2: Создание JWT

Сгенерируйте JWT токен.

**Ожидаемый результат:**
```
Создание: jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Создание: jwt.NewWithClaims(jwt.SigningMethodHS256, claims)")
}
```

**Ожидаемый вывод:**
```
Создание: jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
```

**Баллы:** 15

### Задание 3: Подпись JWT

Подпишите токен.

**Ожидаемый результат:**
```
Подпись: tokenString, err := token.SignedString(secretKey)
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Подпись: tokenString, err := token.SignedString(secretKey)")
}
```

**Ожидаемый вывод:**
```
Подпись: tokenString, err := token.SignedString(secretKey)
```

**Баллы:** 15

### Задание 4: Валидация JWT

Проверьте токен.

**Ожидаемый результат:**
```
Валидация: jwt.Parse(tokenString, keyFunc)
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Валидация: jwt.Parse(tokenString, keyFunc)")
}
```

**Ожидаемый вывод:**
```
Валидация: jwt.Parse(tokenString, keyFunc)
```

**Баллы:** 15

### Задание 5: Claims

Работа с claims.

**Ожидаемый результат:**
```
Claims: jwt.MapClaims{"user_id": 1, "exp": time}
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Claims: jwt.MapClaims{\"user_id\": 1, \"exp\": time}")
}
```

**Ожидаемый вывод:**
```
Claims: jwt.MapClaims{"user_id": 1, "exp": time}
```

**Баллы:** 15
