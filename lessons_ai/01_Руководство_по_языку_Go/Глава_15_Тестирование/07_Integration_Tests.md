# Integration Tests (Интеграционные тесты)

---

## 💡 Ключевые идеи

1. **Integration tests** — тестирование взаимодействия компонентов системы
2. **Test fixtures** — подготовка тестовых данных и окружения
3. **TestMain** — настройка и очистка перед/после всех тестов
4. **Build tags** — разделение unit и integration тестов
5. **Docker** — изолированное окружение для интеграционных тестов
6. **Test containers** — автоматический запуск зависимостей

### Пирамида тестирования

```
        ╱╲
       ╱  ╲         E2E Tests (мало)
      ╱────╲
     ╱      ╲       Integration Tests (средне)
    ╱────────╲
   ╱          ╲     Unit Tests (много)
  ╱────────────╲
```

---

## 📖 Теория

### Что такое интеграционные тесты?

**Unit-тесты** проверяют отдельные функции изолированно, используя моки для зависимостей. **Интеграционные тесты** проверяют, как несколько компонентов работают **вместе**.

Представьте автомобиль:
- **Unit-тест:** проверяем, что двигатель заводится
- **Интеграционный тест:** проверяем, что двигатель + коробка передач + колёса вместе позволяют машине ехать

### Зачем нужны интеграционные тесты?

Unit-тесты с моками не гарантируют, что система работает целиком. Мок базы данных всегда возвращает то, что вы ему сказали — но реальная база может вести себя иначе.

**Типичные проблемы, которые находят интеграционные тесты:**
- SQL-запрос с синтаксической ошибкой
- Неправильная конфигурация подключения
- Проблемы с кодировкой данных
- Race conditions при параллельном доступе

### Пирамида тестирования

```
        E2E (мало)
    Интеграционные (средне)
      Unit (много)
```

- **Unit-тесты** (70%) — быстрые, много, тестируют логику
- **Интеграционные** (20%) — средне, проверяют связи между компонентами
- **E2E** (10%) — медленные, проверяют систему от и до

### Build Tags — разделение тестов

В Go можно пометить файлы тегами и запускать их отдельно:

```go
//go:build integration

package mypackage
```

Такой файл не будет компилироваться при обычном `go test`, только с флагом:
```bash
go test -tags=integration ./...
```

Это удобно, потому что:
- Unit-тесты запускаются быстро (каждый коммит)
- Интеграционные — медленнее (перед merge, в CI)

### TestMain — setup и teardown

`TestMain` — специальная функция, которая выполняется **один раз** перед всеми тестами пакета:

```go
func TestMain(m *testing.M) {
    // 1. Setup — запускаем БД, создаём таблицы
    db := setupDatabase()
    
    // 2. Запуск всех тестов
    code := m.Run()
    
    // 3. Teardown — очистка
    db.Close()
    
    // 4. Выход с кодом результата
    os.Exit(code)
}
```

### Test Containers — автоматические зависимости

Библиотека `testcontainers-go` позволяет автоматически запускать Docker-контейнеры для тестов:

```go
container, _ := postgres.Run(ctx, "postgres:15")
connString, _ := container.ConnectionString(ctx)
// используем реальный PostgreSQL
```

**Преимущества:**
- Каждый тест получает чистую базу
- Не нужно устанавливать зависимости локально
- Работает в CI без настройки

### Изоляция тестов

Каждый тест должен быть **независимым**. Способы изоляции:

1. **Транзакции с откатом:**
```go
tx, _ := db.Begin()
defer tx.Rollback()  // всё откатится
// тест работает с tx
```

2. **Очистка перед тестом:**
```go
func TestSomething(t *testing.T) {
    db.Exec("DELETE FROM users")  // чистим таблицу
    // тест
}
```

3. **Отдельная БД для каждого теста** (медленно, но надёжно)

---

## 📋 Синтаксис

### TestMain

```go
func TestMain(m *testing.M) {
    // Setup
    setup()
    
    // Run tests
    code := m.Run()
    
    // Teardown
    teardown()
    
    os.Exit(code)
}
```

### Build Tags

```go
//go:build integration
// +build integration

package mypackage
```

### Запуск с тегами

```bash
go test -tags=integration ./...
go test -tags=integration -v ./...
```

---

## 💻 Примеры кода

### Пример 1: Интеграционный тест с БД

```go
//go:build integration

package repository

import (
    "database/sql"
    "os"
    "testing"
    
    _ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
    // Setup: подключаемся к тестовой БД
    var err error
    dsn := os.Getenv("TEST_DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://test:test@localhost:5432/testdb?sslmode=disable"
    }
    
    testDB, err = sql.Open("postgres", dsn)
    if err != nil {
        panic("failed to connect to test database: " + err.Error())
    }
    
    // Создаём таблицы
    setupSchema()
    
    // Запускаем тесты
    code := m.Run()
    
    // Teardown: очищаем БД
    teardownSchema()
    testDB.Close()
    
    os.Exit(code)
}

func setupSchema() {
    testDB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
}

func teardownSchema() {
    testDB.Exec("DROP TABLE IF EXISTS users")
}

// Очистка данных между тестами
func cleanupUsers(t *testing.T) {
    _, err := testDB.Exec("DELETE FROM users")
    if err != nil {
        t.Fatalf("failed to cleanup users: %v", err)
    }
}

// Тесты
func TestUserRepository_Create(t *testing.T) {
    cleanupUsers(t)
    
    repo := NewUserRepository(testDB)
    
    user := &User{Name: "John", Email: "john@example.com"}
    err := repo.Create(user)
    
    if err != nil {
        t.Fatalf("failed to create user: %v", err)
    }
    
    if user.ID == 0 {
        t.Error("expected user ID to be set")
    }
}

func TestUserRepository_GetByID(t *testing.T) {
    cleanupUsers(t)
    
    repo := NewUserRepository(testDB)
    
    // Создаём пользователя
    user := &User{Name: "John", Email: "john@example.com"}
    repo.Create(user)
    
    // Получаем по ID
    found, err := repo.GetByID(user.ID)
    
    if err != nil {
        t.Fatalf("failed to get user: %v", err)
    }
    
    if found.Name != "John" {
        t.Errorf("expected name John, got %s", found.Name)
    }
}

func TestUserRepository_GetByEmail(t *testing.T) {
    cleanupUsers(t)
    
    repo := NewUserRepository(testDB)
    
    // Создаём пользователя
    user := &User{Name: "John", Email: "john@example.com"}
    repo.Create(user)
    
    // Получаем по email
    found, err := repo.GetByEmail("john@example.com")
    
    if err != nil {
        t.Fatalf("failed to get user: %v", err)
    }
    
    if found.ID != user.ID {
        t.Errorf("expected ID %d, got %d", user.ID, found.ID)
    }
}
```

### Пример 2: Тест с SQLite в памяти

```go
package repository

import (
    "database/sql"
    "testing"
    
    _ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    
    // SQLite в памяти — быстро и изолированно
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open db: %v", err)
    }
    
    // Создаём схему
    _, err = db.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL
        )
    `)
    if err != nil {
        t.Fatalf("failed to create schema: %v", err)
    }
    
    return db
}

func TestUserCRUD(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewUserRepository(db)
    
    t.Run("Create", func(t *testing.T) {
        user := &User{Name: "Alice", Email: "alice@example.com"}
        err := repo.Create(user)
        
        if err != nil {
            t.Fatalf("create failed: %v", err)
        }
        if user.ID == 0 {
            t.Error("expected ID to be set")
        }
    })
    
    t.Run("Read", func(t *testing.T) {
        // Создаём
        user := &User{Name: "Bob", Email: "bob@example.com"}
        repo.Create(user)
        
        // Читаем
        found, err := repo.GetByID(user.ID)
        if err != nil {
            t.Fatalf("read failed: %v", err)
        }
        if found.Name != "Bob" {
            t.Errorf("expected Bob, got %s", found.Name)
        }
    })
    
    t.Run("Update", func(t *testing.T) {
        user := &User{Name: "Charlie", Email: "charlie@example.com"}
        repo.Create(user)
        
        user.Name = "Charles"
        err := repo.Update(user)
        
        if err != nil {
            t.Fatalf("update failed: %v", err)
        }
        
        found, _ := repo.GetByID(user.ID)
        if found.Name != "Charles" {
            t.Errorf("expected Charles, got %s", found.Name)
        }
    })
    
    t.Run("Delete", func(t *testing.T) {
        user := &User{Name: "David", Email: "david@example.com"}
        repo.Create(user)
        
        err := repo.Delete(user.ID)
        if err != nil {
            t.Fatalf("delete failed: %v", err)
        }
        
        _, err = repo.GetByID(user.ID)
        if err == nil {
            t.Error("expected error, got nil")
        }
    })
}
```

### Пример 3: Интеграционный тест API

```go
//go:build integration

package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

var testServer *httptest.Server
var testApp *Application

func TestMain(m *testing.M) {
    // Setup
    db := setupTestDatabase()
    testApp = NewApplication(db)
    testServer = httptest.NewServer(testApp.Router())
    
    // Run
    code := m.Run()
    
    // Teardown
    testServer.Close()
    db.Close()
    
    os.Exit(code)
}

func TestAPI_CreateAndGetUser(t *testing.T) {
    // Очищаем данные
    testApp.DB.Exec("DELETE FROM users")
    
    // 1. Создаём пользователя
    createBody := map[string]string{
        "name":  "John",
        "email": "john@example.com",
    }
    bodyBytes, _ := json.Marshal(createBody)
    
    resp, err := http.Post(
        testServer.URL+"/api/users",
        "application/json",
        bytes.NewBuffer(bodyBytes),
    )
    if err != nil {
        t.Fatalf("create request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("expected 201, got %d", resp.StatusCode)
    }
    
    var createdUser User
    json.NewDecoder(resp.Body).Decode(&createdUser)
    
    if createdUser.ID == 0 {
        t.Error("expected user ID")
    }
    
    // 2. Получаем пользователя
    resp, err = http.Get(testServer.URL + "/api/users/" + fmt.Sprint(createdUser.ID))
    if err != nil {
        t.Fatalf("get request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected 200, got %d", resp.StatusCode)
    }
    
    var foundUser User
    json.NewDecoder(resp.Body).Decode(&foundUser)
    
    if foundUser.Name != "John" {
        t.Errorf("expected John, got %s", foundUser.Name)
    }
}

func TestAPI_UserNotFound(t *testing.T) {
    resp, err := http.Get(testServer.URL + "/api/users/99999")
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusNotFound {
        t.Errorf("expected 404, got %d", resp.StatusCode)
    }
}
```

### Пример 4: Тест с Redis

```go
//go:build integration

package cache

import (
    "context"
    "os"
    "testing"
    "time"
    
    "github.com/go-redis/redis/v8"
)

var testRedis *redis.Client
var ctx = context.Background()

func TestMain(m *testing.M) {
    // Setup Redis
    redisAddr := os.Getenv("TEST_REDIS_URL")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }
    
    testRedis = redis.NewClient(&redis.Options{
        Addr: redisAddr,
        DB:   15,  // используем отдельную БД для тестов
    })
    
    // Проверяем подключение
    if err := testRedis.Ping(ctx).Err(); err != nil {
        panic("failed to connect to Redis: " + err.Error())
    }
    
    code := m.Run()
    
    // Очищаем и закрываем
    testRedis.FlushDB(ctx)
    testRedis.Close()
    
    os.Exit(code)
}

func cleanupRedis(t *testing.T) {
    t.Helper()
    testRedis.FlushDB(ctx)
}

func TestCache_SetGet(t *testing.T) {
    cleanupRedis(t)
    
    cache := NewCache(testRedis)
    
    // Set
    err := cache.Set(ctx, "key1", "value1", time.Minute)
    if err != nil {
        t.Fatalf("set failed: %v", err)
    }
    
    // Get
    value, err := cache.Get(ctx, "key1")
    if err != nil {
        t.Fatalf("get failed: %v", err)
    }
    
    if value != "value1" {
        t.Errorf("expected value1, got %s", value)
    }
}

func TestCache_Expiration(t *testing.T) {
    cleanupRedis(t)
    
    cache := NewCache(testRedis)
    
    // Set с коротким TTL
    err := cache.Set(ctx, "temp", "data", 100*time.Millisecond)
    if err != nil {
        t.Fatalf("set failed: %v", err)
    }
    
    // Сразу доступно
    _, err = cache.Get(ctx, "temp")
    if err != nil {
        t.Error("expected value to exist immediately")
    }
    
    // Ждём истечения
    time.Sleep(150 * time.Millisecond)
    
    // Должно исчезнуть
    _, err = cache.Get(ctx, "temp")
    if err == nil {
        t.Error("expected key to be expired")
    }
}
```

### Пример 5: Testcontainers

```go
//go:build integration

package repository

import (
    "context"
    "database/sql"
    "testing"
    
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    _ "github.com/lib/pq"
)

func setupPostgresContainer(t *testing.T) (*sql.DB, func()) {
    ctx := context.Background()
    
    // Создаём контейнер
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
            "POSTGRES_DB":       "testdb",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections").
            WithOccurrence(2).
            WithStartupTimeout(60 * time.Second),
    }
    
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("failed to start container: %v", err)
    }
    
    // Получаем хост и порт
    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "5432")
    
    // Подключаемся
    dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        t.Fatalf("failed to connect: %v", err)
    }
    
    // Функция очистки
    cleanup := func() {
        db.Close()
        container.Terminate(ctx)
    }
    
    return db, cleanup
}

func TestWithTestcontainers(t *testing.T) {
    db, cleanup := setupPostgresContainer(t)
    defer cleanup()
    
    // Создаём схему
    db.Exec(`CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT)`)
    
    // Тестируем
    repo := NewUserRepository(db)
    
    user := &User{Name: "Test"}
    err := repo.Create(user)
    
    if err != nil {
        t.Fatalf("create failed: %v", err)
    }
    
    if user.ID == 0 {
        t.Error("expected ID")
    }
}
```

### Пример 6: Параллельные интеграционные тесты

```go
//go:build integration

package service

import (
    "testing"
)

func TestParallelIntegration(t *testing.T) {
    // Общий setup
    db := setupTestDB(t)
    defer db.Close()
    
    // Параллельные тесты с изолированными данными
    t.Run("CreateUser", func(t *testing.T) {
        t.Parallel()
        
        // Каждый тест использует уникальные данные
        user := &User{
            Name:  "User_CreateUser",
            Email: "create_user@example.com",
        }
        // ... test
    })
    
    t.Run("UpdateUser", func(t *testing.T) {
        t.Parallel()
        
        user := &User{
            Name:  "User_UpdateUser",
            Email: "update_user@example.com",
        }
        // ... test
    })
    
    t.Run("DeleteUser", func(t *testing.T) {
        t.Parallel()
        
        user := &User{
            Name:  "User_DeleteUser",
            Email: "delete_user@example.com",
        }
        // ... test
    })
}
```

### Пример 7: Тест с фикстурами

```go
//go:build integration

package repository

import (
    "encoding/json"
    "os"
    "testing"
)

type TestFixtures struct {
    Users    []User    `json:"users"`
    Products []Product `json:"products"`
}

func loadFixtures(t *testing.T, db *sql.DB) {
    t.Helper()
    
    // Читаем файл с фикстурами
    data, err := os.ReadFile("testdata/fixtures.json")
    if err != nil {
        t.Fatalf("failed to read fixtures: %v", err)
    }
    
    var fixtures TestFixtures
    if err := json.Unmarshal(data, &fixtures); err != nil {
        t.Fatalf("failed to parse fixtures: %v", err)
    }
    
    // Загружаем пользователей
    for _, user := range fixtures.Users {
        db.Exec(
            "INSERT INTO users (id, name, email) VALUES (?, ?, ?)",
            user.ID, user.Name, user.Email,
        )
    }
    
    // Загружаем продукты
    for _, product := range fixtures.Products {
        db.Exec(
            "INSERT INTO products (id, name, price) VALUES (?, ?, ?)",
            product.ID, product.Name, product.Price,
        )
    }
}

func TestWithFixtures(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    loadFixtures(t, db)
    
    repo := NewUserRepository(db)
    
    // Тест с предзагруженными данными
    user, err := repo.GetByID(1)
    if err != nil {
        t.Fatalf("failed to get user: %v", err)
    }
    
    // Проверяем данные из фикстур
    if user.Name != "John Doe" {
        t.Errorf("expected John Doe, got %s", user.Name)
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Тесты зависят от порядка выполнения

```go
// ❌ ПЛОХО — тест зависит от предыдущего
func TestCreate(t *testing.T) {
    repo.Create(&User{Name: "John"})
}

func TestGet(t *testing.T) {
    user, _ := repo.GetByName("John")  // зависит от TestCreate
}

// ✅ ХОРОШО — каждый тест независим
func TestCreate(t *testing.T) {
    cleanup()
    repo.Create(&User{Name: "John"})
}

func TestGet(t *testing.T) {
    cleanup()
    repo.Create(&User{Name: "Jane"})  // свои данные
    user, _ := repo.GetByName("Jane")
}
```

### 2. Не очищают данные между тестами

```go
// ❌ ПЛОХО
func TestA(t *testing.T) {
    db.Exec("INSERT INTO users ...")
}

func TestB(t *testing.T) {
    // данные от TestA могут влиять
}

// ✅ ХОРОШО
func TestA(t *testing.T) {
    db.Exec("DELETE FROM users")  // очистка
    db.Exec("INSERT INTO users ...")
}
```

### 3. Хардкод адресов

```go
// ❌ ПЛОХО
db, _ := sql.Open("postgres", "localhost:5432/testdb")

// ✅ ХОРОШО — используем переменные окружения
dsn := os.Getenv("TEST_DATABASE_URL")
if dsn == "" {
    dsn = "localhost:5432/testdb"  // default для локальной разработки
}
```

---

## 🏋️ Практические задания

### Задание 1: Интеграционный тест с SQLite

Напишите интеграционные тесты для UserRepository с SQLite.

**Начальный код:**
```go
package main

import (
    "errors"
    "fmt"
)

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

```go
// repository_test.go
//go:build integration

package repository

import (
    "database/sql"
    "os"
    "testing"
    
    _ "github.com/mattn/go-sqlite3"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
    // Setup
    var err error
    testDB, err = sql.Open("sqlite3", ":memory:")
    if err != nil {
        panic(err)
    }
    
    // Создаём таблицу
    testDB.Exec(`
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT NOT NULL
        )
    `)
    
    code := m.Run()
    
    // Teardown
    testDB.Close()
    os.Exit(code)
}

func cleanup(t *testing.T) {
    t.Helper()
    testDB.Exec("DELETE FROM users")
}

func TestUserRepository_Create(t *testing.T) {
    cleanup(t)
    
    repo := NewUserRepository(testDB)
    user := &User{Name: "John", Email: "john@example.com"}
    
    err := repo.Create(user)
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.ID == 0 {
        t.Error("expected user ID to be set")
    }
}

// Добавьте тесты для GetByID, Update, Delete
```

**Запуск:**
```bash
go test -tags=integration -v
```

**Ожидаемый вывод:**
```
=== RUN   TestUserRepository_Create
--- PASS: TestUserRepository_Create (0.00s)
=== RUN   TestUserRepository_GetByID
--- PASS: TestUserRepository_GetByID (0.00s)
PASS
```

**Баллы:** 15

### Задание 2: E2E тест REST API

Напишите end-to-end тест для REST API с базой данных.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 15

### Задание 3: TestMain с очисткой

Создайте структуру тестов с TestMain и автоматической очисткой.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Определите структуру

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 10

### Задание 4: Тест с фикстурами

Создайте систему загрузки тестовых данных из JSON.

**Начальный код:**
```go
// testdata/fixtures.json
{
    "users": [
        {"id": 1, "name": "John", "email": "john@example.com"},
        {"id": 2, "name": "Jane", "email": "jane@example.com"}
    ],
    "products": [
        {"id": 1, "name": "Laptop", "price": 999.99},
        {"id": 2, "name": "Phone", "price": 499.99}
    ]
}
```

```go
// fixtures_test.go
//go:build integration

package repository

import (
    "encoding/json"
    "os"
    "testing"
)

type Fixtures struct {
    Users    []User    `json:"users"`
    Products []Product `json:"products"`
}

func loadFixtures(t *testing.T, db *sql.DB) Fixtures {
    t.Helper()
    
    data, err := os.ReadFile("testdata/fixtures.json")
    if err != nil {
        t.Fatalf("failed to load fixtures: %v", err)
    }
    
    var fixtures Fixtures
    json.Unmarshal(data, &fixtures)
    
    // Загружаем в БД
    for _, u := range fixtures.Users {
        db.Exec("INSERT INTO users (id, name, email) VALUES (?, ?, ?)",
            u.ID, u.Name, u.Email)
    }
    
    return fixtures
}

func TestWithFixtures(t *testing.T) {
    cleanup := withCleanDB(t)
    defer cleanup()
    
    fixtures := loadFixtures(t, testDB)
    
    repo := NewUserRepository(testDB)
    
    // Проверяем загруженные данные
    user, _ := repo.GetByID(1)
    
    if user.Name != fixtures.Users[0].Name {
        t.Errorf("expected %s, got %s", fixtures.Users[0].Name, user.Name)
    }
}
```

**Баллы:** 10

### Задание 5: Параллельные интеграционные тесты

Напишите параллельные интеграционные тесты с изоляцией данных.

**Начальный код:**
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
=== RUN   TestParallelUsers
=== RUN   TestParallelUsers/Alice
=== PAUSE TestParallelUsers/Alice
=== RUN   TestParallelUsers/Bob
=== PAUSE TestParallelUsers/Bob
=== RUN   TestParallelUsers/Charlie
=== PAUSE TestParallelUsers/Charlie
=== CONT  TestParallelUsers/Alice
=== CONT  TestParallelUsers/Bob
=== CONT  TestParallelUsers/Charlie
--- PASS: TestParallelUsers (0.05s)
PASS
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Go Testing — TestMain](https://pkg.go.dev/testing#hdr-Main)
- [Testcontainers Go](https://golang.testcontainers.org/)
- [Go Build Tags](https://pkg.go.dev/go/build#hdr-Build_Constraints)
