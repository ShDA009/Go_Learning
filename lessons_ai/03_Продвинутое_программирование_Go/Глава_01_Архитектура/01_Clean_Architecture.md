# Clean Architecture в Go

---

## 💡 Ключевые идеи

1. **Clean Architecture** — независимость бизнес-логики от фреймворков и БД
2. **Dependency Rule** — зависимости направлены внутрь (к бизнес-логике)
3. **Слои** — Entities, Use Cases, Adapters, Infrastructure
4. **Инверсия зависимостей** — высокоуровневые модули не зависят от низкоуровневых
5. **Тестируемость** — бизнес-логика тестируется без внешних зависимостей

### Слои Clean Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Infrastructure                      │
│  (Frameworks, DB, External APIs, UI)                │
├─────────────────────────────────────────────────────┤
│                    Adapters                          │
│  (Controllers, Presenters, Gateways)                │
├─────────────────────────────────────────────────────┤
│                   Use Cases                          │
│  (Application Business Rules)                        │
├─────────────────────────────────────────────────────┤
│                    Entities                          │
│  (Enterprise Business Rules)                         │
└─────────────────────────────────────────────────────┘
```

---

## 📖 Теория

### Что такое Clean Architecture?

**Clean Architecture** — это подход к организации кода, предложенный Робертом Мартином (Uncle Bob). Главная идея: **бизнес-логика должна быть независимой** от фреймворков, баз данных, UI и любых внешних систем.

### Проблема, которую решает Clean Architecture

Представьте типичный проект:
```go
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)
    
    // Валидация прямо в handler
    if user.Email == "" {
        http.Error(w, "email required", 400)
        return
    }
    
    // Прямой доступ к БД
    db.Query("INSERT INTO users...")
    
    // Отправка email
    smtp.Send(user.Email, "Welcome!")
    
    json.NewEncoder(w).Encode(user)
}
```

Проблемы:
- **Бизнес-логика смешана с HTTP** — нельзя вызвать из CLI
- **Прямая зависимость от БД** — нельзя заменить на другую БД
- **Сложно тестировать** — нужен реальный HTTP и БД
- **Код нельзя переиспользовать**

### Dependency Rule (Правило зависимостей)

Ключевое правило Clean Architecture: **зависимости направлены внутрь**.

```
Infrastructure → Adapters → Use Cases → Entities
              ←────────────────────────────────
              Код знает только о внутренних слоях
```

- **Entities** — ничего не знают о внешнем мире
- **Use Cases** — знают только об Entities
- **Adapters** — знают о Use Cases и Entities
- **Infrastructure** — знает обо всём

### Четыре слоя Clean Architecture

**1. Entities (Сущности)**
Это бизнес-объекты вашего домена. Они не знают ни о чём, кроме себя:
```go
type User struct {
    ID    int64
    Email string
    Name  string
}

func (u *User) IsActive() bool {
    return u.Email != ""
}
```

**2. Use Cases (Сценарии использования)**
Бизнес-логика приложения. Определяют, ЧТО может делать система:
```go
type UserService struct {
    repo UserRepository  // интерфейс, не конкретная реализация!
}

func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    // Бизнес-правила здесь
    if !isValidEmail(input.Email) {
        return nil, ErrInvalidEmail
    }
    return s.repo.Create(ctx, input)
}
```

**3. Adapters (Адаптеры)**
Преобразуют данные между форматами. Сюда относятся HTTP-хэндлеры, репозитории:
```go
// HTTP Adapter
func (h *Handler) CreateUser(c *gin.Context) {
    var input CreateUserInput
    c.BindJSON(&input)
    user, err := h.userService.CreateUser(c, input)
    c.JSON(200, user)
}

// Repository Adapter
type PostgresUserRepo struct {
    db *sql.DB
}
```

**4. Infrastructure (Инфраструктура)**
Внешние системы: БД, фреймворки, HTTP-серверы, очереди сообщений.

### Инверсия зависимостей

Use Cases не должны зависеть от конкретной БД. Как это работает?

```go
// В слое Use Cases определяем ИНТЕРФЕЙС
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id int64) (*User, error)
}

// В слое Adapters — РЕАЛИЗАЦИЯ
type PostgresUserRepo struct {
    db *sql.DB
}

func (r *PostgresUserRepo) Create(ctx context.Context, user *User) error {
    _, err := r.db.ExecContext(ctx, "INSERT INTO users...")
    return err
}
```

Теперь Use Case зависит от интерфейса (абстракции), а не от Postgres.

### Преимущества Clean Architecture

| Проблема | Решение |
|----------|---------|
| Сложно тестировать | Моки для интерфейсов |
| Привязка к фреймворку | Фреймворк только в Infrastructure |
| Смена БД = переписать всё | Только новый адаптер |
| Код спагетти | Чёткое разделение ответственности |

### Когда НЕ нужна Clean Architecture?

- Маленькие проекты (< 5000 строк)
- Прототипы и MVP
- Скрипты и утилиты
- Проекты с одним разработчиком и без перспективы роста

Clean Architecture требует больше кода. Используйте её, когда преимущества перевешивают затраты.

---

## 📋 Структура проекта

```
project/
├── cmd/
│   └── api/
│       └── main.go              # Точка входа
├── internal/
│   ├── domain/                  # Entities
│   │   ├── user.go
│   │   ├── order.go
│   │   └── errors.go
│   ├── usecase/                 # Use Cases
│   │   ├── user/
│   │   │   ├── interface.go
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   └── order/
│   │       └── service.go
│   ├── adapter/                 # Adapters
│   │   ├── http/
│   │   │   ├── handler.go
│   │   │   └── middleware.go
│   │   └── repository/
│   │       ├── postgres/
│   │       │   └── user.go
│   │       └── memory/
│   │           └── user.go
│   └── infrastructure/          # Infrastructure
│       ├── database/
│       │   └── postgres.go
│       └── config/
│           └── config.go
├── pkg/                         # Shared packages
│   └── logger/
│       └── logger.go
├── go.mod
└── go.sum
```

---

## 💻 Примеры кода

### Domain Layer (Entities)

```go
// internal/domain/user.go
package domain

import (
    "errors"
    "time"
)

// Domain errors
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidEmail     = errors.New("invalid email")
)

// User — доменная сущность
type User struct {
    ID        int64
    Email     string
    Name      string
    Password  string // hashed
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Validate — бизнес-валидация
func (u *User) Validate() error {
    if u.Email == "" {
        return ErrInvalidEmail
    }
    if u.Name == "" {
        return errors.New("name is required")
    }
    return nil
}

// CanBeDeleted — бизнес-правило
func (u *User) CanBeDeleted() bool {
    // Пример: нельзя удалить пользователя младше 24 часов
    return time.Since(u.CreatedAt) > 24*time.Hour
}
```

### Use Case Layer

```go
// internal/usecase/user/interface.go
package user

import (
    "context"
    "myapp/internal/domain"
)

// Repository — интерфейс репозитория (определяется use case)
type Repository interface {
    Create(ctx context.Context, user *domain.User) error
    GetByID(ctx context.Context, id int64) (*domain.User, error)
    GetByEmail(ctx context.Context, email string) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id int64) error
}

// PasswordHasher — интерфейс для хэширования паролей
type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(hash, password string) bool
}

// EmailSender — интерфейс для отправки email
type EmailSender interface {
    SendWelcome(ctx context.Context, email, name string) error
}
```

```go
// internal/usecase/user/service.go
package user

import (
    "context"
    "myapp/internal/domain"
)

// Service — бизнес-логика пользователей
type Service struct {
    repo     Repository
    hasher   PasswordHasher
    emailer  EmailSender
}

func NewService(repo Repository, hasher PasswordHasher, emailer EmailSender) *Service {
    return &Service{
        repo:    repo,
        hasher:  hasher,
        emailer: emailer,
    }
}

// CreateUserInput — входные данные для создания пользователя
type CreateUserInput struct {
    Email    string
    Name     string
    Password string
}

// CreateUser — use case создания пользователя
func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*domain.User, error) {
    // Проверяем, существует ли пользователь
    existing, err := s.repo.GetByEmail(ctx, input.Email)
    if err != nil && err != domain.ErrUserNotFound {
        return nil, err
    }
    if existing != nil {
        return nil, domain.ErrUserAlreadyExists
    }
    
    // Хэшируем пароль
    hashedPassword, err := s.hasher.Hash(input.Password)
    if err != nil {
        return nil, err
    }
    
    // Создаём пользователя
    user := &domain.User{
        Email:    input.Email,
        Name:     input.Name,
        Password: hashedPassword,
    }
    
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // Отправляем приветственное письмо (async)
    go s.emailer.SendWelcome(context.Background(), user.Email, user.Name)
    
    return user, nil
}

// GetUser — use case получения пользователя
func (s *Service) GetUser(ctx context.Context, id int64) (*domain.User, error) {
    return s.repo.GetByID(ctx, id)
}

// DeleteUser — use case удаления пользователя
func (s *Service) DeleteUser(ctx context.Context, id int64) error {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    // Бизнес-правило
    if !user.CanBeDeleted() {
        return errors.New("user cannot be deleted yet")
    }
    
    return s.repo.Delete(ctx, id)
}
```

### Adapter Layer — Repository

```go
// internal/adapter/repository/postgres/user.go
package postgres

import (
    "context"
    "database/sql"
    "myapp/internal/domain"
    "myapp/internal/usecase/user"
)

// Проверяем, что реализуем интерфейс
var _ user.Repository = (*UserRepository)(nil)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
    query := `
        INSERT INTO users (email, name, password, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `
    
    return r.db.QueryRowContext(ctx, query, u.Email, u.Name, u.Password).
        Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
    query := `SELECT id, email, name, password, created_at, updated_at FROM users WHERE id = $1`
    
    var u domain.User
    err := r.db.QueryRowContext(ctx, query, id).
        Scan(&u.ID, &u.Email, &u.Name, &u.Password, &u.CreatedAt, &u.UpdatedAt)
    
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    
    return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `SELECT id, email, name, password, created_at, updated_at FROM users WHERE email = $1`
    
    var u domain.User
    err := r.db.QueryRowContext(ctx, query, email).
        Scan(&u.ID, &u.Email, &u.Name, &u.Password, &u.CreatedAt, &u.UpdatedAt)
    
    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, err
    }
    
    return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
    query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3`
    _, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.ID)
    return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}
```

### Adapter Layer — HTTP Handler

```go
// internal/adapter/http/user_handler.go
package http

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "myapp/internal/usecase/user"
)

type UserHandler struct {
    userService *user.Service
}

func NewUserHandler(userService *user.Service) *UserHandler {
    return &UserHandler{userService: userService}
}

// CreateUserRequest — DTO для запроса
type CreateUserRequest struct {
    Email    string `json:"email"`
    Name     string `json:"name"`
    Password string `json:"password"`
}

// UserResponse — DTO для ответа
type UserResponse struct {
    ID    int64  `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    // Преобразуем DTO в input для use case
    input := user.CreateUserInput{
        Email:    req.Email,
        Name:     req.Name,
        Password: req.Password,
    }
    
    // Вызываем use case
    u, err := h.userService.CreateUser(r.Context(), input)
    if err != nil {
        handleError(w, err)
        return
    }
    
    // Преобразуем domain entity в DTO
    response := UserResponse{
        ID:    u.ID,
        Email: u.Email,
        Name:  u.Name,
    }
    
    respondJSON(w, http.StatusCreated, response)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid id")
        return
    }
    
    u, err := h.userService.GetUser(r.Context(), id)
    if err != nil {
        handleError(w, err)
        return
    }
    
    response := UserResponse{
        ID:    u.ID,
        Email: u.Email,
        Name:  u.Name,
    }
    
    respondJSON(w, http.StatusOK, response)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
    respondJSON(w, status, map[string]string{"error": message})
}

func handleError(w http.ResponseWriter, err error) {
    switch err {
    case domain.ErrUserNotFound:
        respondError(w, http.StatusNotFound, err.Error())
    case domain.ErrUserAlreadyExists:
        respondError(w, http.StatusConflict, err.Error())
    case domain.ErrInvalidEmail:
        respondError(w, http.StatusBadRequest, err.Error())
    default:
        respondError(w, http.StatusInternalServerError, "internal error")
    }
}
```

### Infrastructure — Dependency Injection

```go
// cmd/api/main.go
package main

import (
    "database/sql"
    "log"
    "net/http"
    
    "myapp/internal/adapter/http"
    "myapp/internal/adapter/repository/postgres"
    "myapp/internal/infrastructure/config"
    "myapp/internal/usecase/user"
    "myapp/pkg/hasher"
    "myapp/pkg/mailer"
    
    _ "github.com/lib/pq"
)

func main() {
    // Load config
    cfg := config.Load()
    
    // Infrastructure
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Adapters
    userRepo := postgres.NewUserRepository(db)
    passwordHasher := hasher.NewBcryptHasher()
    emailSender := mailer.NewSMTPMailer(cfg.SMTPConfig)
    
    // Use Cases
    userService := user.NewService(userRepo, passwordHasher, emailSender)
    
    // HTTP Handlers
    userHandler := http.NewUserHandler(userService)
    
    // Routes
    mux := http.NewServeMux()
    mux.HandleFunc("POST /users", userHandler.Create)
    mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
    
    // Start server
    log.Printf("Server starting on %s", cfg.ServerAddr)
    if err := http.ListenAndServe(cfg.ServerAddr, mux); err != nil {
        log.Fatal(err)
    }
}
```

### Тестирование Use Case

```go
// internal/usecase/user/service_test.go
package user

import (
    "context"
    "testing"
    
    "myapp/internal/domain"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock Repository
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *domain.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

// ... other methods

// Mock Hasher
type MockHasher struct {
    mock.Mock
}

func (m *MockHasher) Hash(password string) (string, error) {
    args := m.Called(password)
    return args.String(0), args.Error(1)
}

func (m *MockHasher) Compare(hash, password string) bool {
    args := m.Called(hash, password)
    return args.Bool(0)
}

// Mock Emailer
type MockEmailer struct {
    mock.Mock
}

func (m *MockEmailer) SendWelcome(ctx context.Context, email, name string) error {
    args := m.Called(ctx, email, name)
    return args.Error(0)
}

// Tests
func TestService_CreateUser(t *testing.T) {
    ctx := context.Background()
    
    t.Run("success", func(t *testing.T) {
        repo := new(MockRepository)
        hasher := new(MockHasher)
        emailer := new(MockEmailer)
        
        service := NewService(repo, hasher, emailer)
        
        // Setup expectations
        repo.On("GetByEmail", ctx, "test@example.com").Return(nil, domain.ErrUserNotFound)
        hasher.On("Hash", "password123").Return("hashed_password", nil)
        repo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
        emailer.On("SendWelcome", mock.Anything, "test@example.com", "John").Return(nil)
        
        // Execute
        input := CreateUserInput{
            Email:    "test@example.com",
            Name:     "John",
            Password: "password123",
        }
        
        user, err := service.CreateUser(ctx, input)
        
        // Assert
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, "test@example.com", user.Email)
        
        repo.AssertExpectations(t)
        hasher.AssertExpectations(t)
    })
    
    t.Run("user already exists", func(t *testing.T) {
        repo := new(MockRepository)
        hasher := new(MockHasher)
        emailer := new(MockEmailer)
        
        service := NewService(repo, hasher, emailer)
        
        existingUser := &domain.User{ID: 1, Email: "test@example.com"}
        repo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)
        
        input := CreateUserInput{
            Email:    "test@example.com",
            Name:     "John",
            Password: "password123",
        }
        
        _, err := service.CreateUser(ctx, input)
        
        assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
    })
}
```

---

## ⚠️ Частые ошибки

### 1. Нарушение Dependency Rule

```go
// ❌ ПЛОХО — domain зависит от infrastructure
package domain

import "database/sql"

type User struct {
    db *sql.DB  // зависимость от инфраструктуры!
}

// ✅ ХОРОШО — domain не знает о БД
package domain

type User struct {
    ID    int64
    Email string
}
```

### 2. Бизнес-логика в handlers

```go
// ❌ ПЛОХО
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // Бизнес-логика в handler!
    if existingUser != nil {
        // ...
    }
    hashedPassword := bcrypt.Hash(password)
    // ...
}

// ✅ ХОРОШО — handler только преобразует данные
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    user, err := h.service.CreateUser(r.Context(), req.ToInput())
    if err != nil { ... }
    
    respondJSON(w, http.StatusCreated, UserResponse{}.FromDomain(user))
}
```

---

## 🏋️ Практические задания

### Задание 1: Domain Entity

Создайте доменную сущность с валидацией.

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

### Задание 2: Repository Interface и Implementation

Создайте интерфейс репозитория и in-memory реализацию.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Определите и реализуйте интерфейс

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 15

### Задание 3: Use Case

Реализуйте use case для создания заказа.

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

### Задание 4: HTTP Handler

Создайте HTTP handler для API.

**Начальный код:**
```go
package main

import (
    "fmt"
    "net/http"
)

// TODO: Реализуйте HTTP обработчик

func main() {
    // Ваш код здесь
    
}
```

**Баллы:** 15

### Задание 5: Dependency Injection

Соберите все компоненты вместе.

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

- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Clean Architecture Example](https://github.com/bxcodec/go-clean-arch)
