# Введение в GORM

---

## 💡 Ключевые идеи

1. **GORM** — самый популярный ORM для Go
2. **Auto Migration** — автоматическое создание таблиц
3. **CRUD** — Create, Read, Update, Delete операции
4. **Associations** — связи между моделями
5. **Hooks** — callback'и для событий

---

## 📖 Теория

### Что такое ORM?

**ORM (Object-Relational Mapping)** — это техника, которая позволяет работать с базой данных через объекты вашего языка программирования, а не через SQL-запросы.

```go
// Без ORM
db.Query("SELECT * FROM users WHERE id = ?", 1)

// С ORM (GORM)
var user User
db.First(&user, 1)
```

### Зачем нужен GORM?

**Проблемы работы с database/sql напрямую:**
- Много boilerplate кода для сканирования результатов
- Ручное написание всех SQL-запросов
- Нет автоматических миграций
- Сложно работать со связями между таблицами

**GORM решает эти проблемы:**
- Автоматическое маппирование строк на структуры
- Генерация SQL из методов Go
- Auto Migration — создание таблиц по структурам
- Простая работа со связями (один-к-одному, один-ко-многим)

### Основные концепции GORM

**1. Модель** — структура Go, представляющая таблицу:
```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    
    Name  string `gorm:"size:100"`
    Email string `gorm:"uniqueIndex"`
}
```

**2. gorm.Model** — встраиваемая структура с базовыми полями:
```go
type User struct {
    gorm.Model  // ID, CreatedAt, UpdatedAt, DeletedAt
    Name string
}
```

**3. Теги** — настройка полей через теги:
```go
type Product struct {
    Code  string `gorm:"primaryKey;size:10"`
    Price uint   `gorm:"default:100"`
    Name  string `gorm:"not null;index"`
}
```

### CRUD операции

```go
// Create
user := User{Name: "John", Email: "john@example.com"}
db.Create(&user)  // user.ID теперь заполнен

// Read
var user User
db.First(&user, 1)                    // по ID
db.First(&user, "email = ?", "john@example.com")
db.Find(&users)                       // все записи
db.Where("age > ?", 18).Find(&users)  // с условием

// Update
db.Model(&user).Update("Name", "Jane")
db.Model(&user).Updates(User{Name: "Jane", Age: 30})
db.Save(&user)  // сохраняет все поля

// Delete
db.Delete(&user, 1)  // Soft Delete (если есть DeletedAt)
db.Unscoped().Delete(&user)  // Hard Delete
```

### Связи между моделями

**Один-к-одному:**
```go
type User struct {
    gorm.Model
    Profile Profile  // has one
}

type Profile struct {
    ID     uint
    UserID uint  // foreign key
    Bio    string
}
```

**Один-ко-многим:**
```go
type User struct {
    gorm.Model
    Orders []Order  // has many
}

type Order struct {
    ID     uint
    UserID uint  // foreign key
    Amount float64
}
```

**Многие-ко-многим:**
```go
type User struct {
    gorm.Model
    Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    ID   uint
    Name string
}
```

### Preload — загрузка связей

```go
// Без Preload — N+1 проблема
var users []User
db.Find(&users)
for _, user := range users {
    db.Find(&user.Orders, "user_id = ?", user.ID)  // N запросов!
}

// С Preload — 2 запроса
db.Preload("Orders").Find(&users)
```

### GORM vs database/sql

| | database/sql | GORM |
|---|---|---|
| Производительность | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| Простота | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| Контроль | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| Миграции | Нужен sqlx/migrate | Встроенные |
| Связи | Вручную | Автоматически |

Используйте GORM для быстрой разработки, database/sql — для критичных к производительности частей.

---

## 📋 Установка

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u gorm.io/driver/sqlite
go get -u gorm.io/driver/mysql
```

---

## 💻 Примеры кода

### Пример 1: Подключение и модели

```go
package main

import (
    "fmt"
    "time"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// Models
type User struct {
    ID        uint           `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete
    
    Email    string `gorm:"uniqueIndex;size:255"`
    Name     string `gorm:"size:100"`
    Age      int
    Active   bool   `gorm:"default:true"`
    
    // Associations
    Profile  Profile
    Orders   []Order
}

type Profile struct {
    ID     uint
    UserID uint   `gorm:"uniqueIndex"`
    Bio    string `gorm:"type:text"`
    Avatar string
}

type Order struct {
    ID        uint
    UserID    uint
    Total     float64
    Status    string `gorm:"default:'pending'"`
    CreatedAt time.Time
    
    Items []OrderItem
}

type OrderItem struct {
    ID        uint
    OrderID   uint
    ProductID uint
    Quantity  int
    Price     float64
}

func main() {
    // Подключение к PostgreSQL
    dsn := "host=localhost user=postgres password=postgres dbname=testdb port=5432"
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info), // логируем SQL
    })
    if err != nil {
        panic("failed to connect database")
    }
    
    // Авто-миграция
    db.AutoMigrate(&User{}, &Profile{}, &Order{}, &OrderItem{})
    
    fmt.Println("Database connected and migrated!")
}
```

### Пример 2: CRUD операции

```go
package main

import "gorm.io/gorm"

func CRUDExamples(db *gorm.DB) {
    // CREATE
    user := User{
        Email: "john@example.com",
        Name:  "John Doe",
        Age:   30,
    }
    result := db.Create(&user)
    fmt.Printf("Created user ID: %d, Rows affected: %d\n", user.ID, result.RowsAffected)
    
    // CREATE multiple
    users := []User{
        {Email: "alice@example.com", Name: "Alice", Age: 25},
        {Email: "bob@example.com", Name: "Bob", Age: 35},
    }
    db.Create(&users)
    
    // READ - Get by ID
    var foundUser User
    db.First(&foundUser, 1) // SELECT * FROM users WHERE id = 1
    
    // READ - Get by condition
    db.First(&foundUser, "email = ?", "john@example.com")
    
    // READ - Get all
    var allUsers []User
    db.Find(&allUsers)
    
    // READ - With conditions
    var activeUsers []User
    db.Where("active = ? AND age > ?", true, 18).Find(&activeUsers)
    
    // READ - Select specific fields
    var names []string
    db.Model(&User{}).Pluck("name", &names)
    
    // UPDATE - Save
    foundUser.Name = "John Smith"
    db.Save(&foundUser)
    
    // UPDATE - Update single field
    db.Model(&foundUser).Update("age", 31)
    
    // UPDATE - Update multiple fields
    db.Model(&foundUser).Updates(User{Name: "John Updated", Age: 32})
    
    // UPDATE - Update with map
    db.Model(&foundUser).Updates(map[string]interface{}{
        "name": "John Map",
        "age":  33,
    })
    
    // UPDATE - Batch update
    db.Model(&User{}).Where("age < ?", 18).Update("active", false)
    
    // DELETE - Soft delete (если есть DeletedAt)
    db.Delete(&foundUser)
    
    // DELETE - Permanent
    db.Unscoped().Delete(&foundUser)
    
    // DELETE - By condition
    db.Where("email = ?", "test@example.com").Delete(&User{})
}
```

### Пример 3: Ассоциации

```go
package main

import "gorm.io/gorm"

func AssociationsExamples(db *gorm.DB) {
    // Создание с ассоциацией
    user := User{
        Name:  "Jane",
        Email: "jane@example.com",
        Profile: Profile{
            Bio:    "Software Developer",
            Avatar: "https://example.com/avatar.jpg",
        },
        Orders: []Order{
            {
                Total:  99.99,
                Status: "completed",
                Items: []OrderItem{
                    {ProductID: 1, Quantity: 2, Price: 49.99},
                },
            },
        },
    }
    db.Create(&user)
    
    // Preload - загрузка связанных данных
    var userWithProfile User
    db.Preload("Profile").First(&userWithProfile, user.ID)
    
    // Preload multiple
    var userFull User
    db.Preload("Profile").
       Preload("Orders").
       Preload("Orders.Items").
       First(&userFull, user.ID)
    
    // Preload with conditions
    var userWithCompletedOrders User
    db.Preload("Orders", "status = ?", "completed").
       First(&userWithCompletedOrders, user.ID)
    
    // Joins
    var usersWithProfiles []User
    db.Joins("Profile").Find(&usersWithProfiles)
    
    // Association methods
    var orders []Order
    db.Model(&user).Association("Orders").Find(&orders)
    
    // Add association
    newOrder := Order{Total: 50.00}
    db.Model(&user).Association("Orders").Append(&newOrder)
    
    // Remove association
    db.Model(&user).Association("Orders").Delete(&newOrder)
    
    // Replace associations
    db.Model(&user).Association("Orders").Replace(&[]Order{
        {Total: 100.00},
        {Total: 200.00},
    })
    
    // Count associations
    count := db.Model(&user).Association("Orders").Count()
    fmt.Println("Orders count:", count)
}
```

### Пример 4: Запросы

```go
package main

import "gorm.io/gorm"

func QueryExamples(db *gorm.DB) {
    var users []User
    
    // Where
    db.Where("name = ?", "John").Find(&users)
    db.Where("name IN ?", []string{"John", "Jane"}).Find(&users)
    db.Where("age BETWEEN ? AND ?", 18, 30).Find(&users)
    db.Where("name LIKE ?", "%John%").Find(&users)
    
    // Or
    db.Where("name = ?", "John").Or("name = ?", "Jane").Find(&users)
    
    // Not
    db.Not("name = ?", "John").Find(&users)
    
    // Order
    db.Order("age desc, name").Find(&users)
    
    // Limit & Offset
    db.Limit(10).Offset(20).Find(&users)
    
    // Group & Having
    type Result struct {
        Status string
        Count  int
    }
    var results []Result
    db.Model(&Order{}).
       Select("status, count(*) as count").
       Group("status").
       Having("count(*) > ?", 5).
       Scan(&results)
    
    // Distinct
    var emails []string
    db.Model(&User{}).Distinct("email").Pluck("email", &emails)
    
    // Subquery
    subQuery := db.Model(&Order{}).Select("user_id").Where("total > ?", 100)
    db.Where("id IN (?)", subQuery).Find(&users)
    
    // Raw SQL
    db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)
    
    // Exec for non-select
    db.Exec("UPDATE users SET active = ? WHERE age < ?", false, 18)
}
```

### Пример 5: Транзакции

```go
package main

import (
    "errors"
    "gorm.io/gorm"
)

func TransactionExample(db *gorm.DB) error {
    // Автоматическая транзакция
    return db.Transaction(func(tx *gorm.DB) error {
        // Создаём пользователя
        user := User{Name: "Transaction User", Email: "tx@example.com"}
        if err := tx.Create(&user).Error; err != nil {
            return err // rollback
        }
        
        // Создаём заказ
        order := Order{UserID: user.ID, Total: 99.99}
        if err := tx.Create(&order).Error; err != nil {
            return err // rollback
        }
        
        // Проверка бизнес-логики
        if order.Total > 1000 {
            return errors.New("order total too high") // rollback
        }
        
        return nil // commit
    })
}

func ManualTransactionExample(db *gorm.DB) error {
    tx := db.Begin()
    
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    if err := tx.Create(&User{Name: "Test"}).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    if err := tx.Create(&Order{Total: 50}).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit().Error
}
```

### Пример 6: Hooks

```go
package main

import (
    "time"
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
)

type UserWithHooks struct {
    gorm.Model
    Email        string
    Password     string
    PasswordHash string
    LastLoginAt  *time.Time
}

// BeforeCreate hook
func (u *UserWithHooks) BeforeCreate(tx *gorm.DB) error {
    // Хэшируем пароль перед созданием
    if u.Password != "" {
        hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.PasswordHash = string(hash)
        u.Password = "" // не сохраняем plain text
    }
    return nil
}

// AfterCreate hook
func (u *UserWithHooks) AfterCreate(tx *gorm.DB) error {
    // Отправляем email после создания
    // sendWelcomeEmail(u.Email)
    return nil
}

// BeforeUpdate hook
func (u *UserWithHooks) BeforeUpdate(tx *gorm.DB) error {
    // Хэшируем пароль если он изменился
    if u.Password != "" {
        hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        u.PasswordHash = string(hash)
        u.Password = ""
    }
    return nil
}

// BeforeDelete hook
func (u *UserWithHooks) BeforeDelete(tx *gorm.DB) error {
    // Проверяем можно ли удалить
    var orderCount int64
    tx.Model(&Order{}).Where("user_id = ?", u.ID).Count(&orderCount)
    if orderCount > 0 {
        return errors.New("cannot delete user with orders")
    }
    return nil
}
```

### Пример 7: Scopes

```go
package main

import "gorm.io/gorm"

// Scopes — переиспользуемые условия
func Active(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func Adult(db *gorm.DB) *gorm.DB {
    return db.Where("age >= ?", 18)
}

func OrderBy(field string, desc bool) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        order := field
        if desc {
            order += " DESC"
        }
        return db.Order(order)
    }
}

func Paginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

func ScopesExample(db *gorm.DB) {
    var users []User
    
    // Использование scopes
    db.Scopes(Active, Adult).Find(&users)
    
    // Комбинирование
    db.Scopes(Active, Adult, OrderBy("name", false), Paginate(1, 10)).Find(&users)
}
```

---

## ⚠️ Частые ошибки

### 1. N+1 Query Problem

```go
// ❌ ПЛОХО — N+1 запросов
var users []User
db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Orders").Find(&user.Orders)
}

// ✅ ХОРОШО — Preload
db.Preload("Orders").Find(&users)
```

### 2. Забыли проверить ошибку

```go
// ❌ ПЛОХО
db.First(&user, 1)

// ✅ ХОРОШО
if err := db.First(&user, 1).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        // не найдено
    }
}
```

---

## 🏋️ Практические задания

### Задание 1: Модели со связями

Создайте модели для блога с отношениями.

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

### Задание 2: CRUD операции

Реализуйте Repository с CRUD операциями.

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

### Задание 3: Scopes и фильтрация

Создайте переиспользуемые scopes.

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

### Задание 4: Транзакции

Реализуйте сложную операцию в транзакции.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 10

### Задание 5: Soft Delete и восстановление

Реализуйте soft delete с возможностью восстановления.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [GORM Documentation](https://gorm.io/docs/)
- [GORM GitHub](https://github.com/go-gorm/gorm)
