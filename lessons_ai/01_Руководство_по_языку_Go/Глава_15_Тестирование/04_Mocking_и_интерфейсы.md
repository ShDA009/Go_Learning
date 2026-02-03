# Mocking и интерфейсы

---

## 💡 Ключевые идеи

1. **Mock** — имитация зависимости для изолированного тестирования
2. **Интерфейсы** — ключ к тестируемому коду в Go
3. **Dependency Injection** — передача зависимостей через параметры
4. **testify/mock** — библиотека для создания моков
5. **Принцип** — тестируем логику, а не внешние зависимости (БД, API, файлы)

### Зачем нужны моки?

| Без моков | С моками |
|-----------|----------|
| Тесты зависят от БД | Тесты независимы |
| Тесты медленные | Тесты быстрые |
| Сложно проверить edge cases | Легко симулировать любой сценарий |
| Нестабильные тесты | Предсказуемые результаты |

---

## 📖 Теория

### Что такое Mock?

**Mock** (мок) — это "фальшивый" объект, который имитирует поведение реального объекта. Это как дублёр актёра в кино — выглядит похоже, но это не настоящий актёр.

### Зачем нужны моки?

Представьте, что вы тестируете сервис заказов, который:
1. Получает пользователя из базы данных
2. Проверяет баланс
3. Списывает деньги
4. Отправляет email

Без моков вам нужна реальная база данных и почтовый сервер. Это:
- **Медленно** — подключение к БД занимает время
- **Нестабильно** — что если БД недоступна?
- **Сложно** — как проверить отправку email?
- **Опасно** — тест может случайно удалить реальные данные

С моками вы заменяете БД и почту на "заглушки", которые возвращают нужные данные.

### Интерфейсы — ключ к тестируемости

В Go моки работают через **интерфейсы**. Если ваш код зависит от конкретного типа — его сложно замокать:

```go
// ❌ Плохо — зависимость от конкретного типа
type Service struct {
    db *sql.DB  // конкретный тип
}
```

Если код зависит от интерфейса — легко подменить реализацию:

```go
// ✅ Хорошо — зависимость от интерфейса
type UserRepository interface {
    GetByID(id int) (*User, error)
}

type Service struct {
    repo UserRepository  // интерфейс
}
```

### Dependency Injection (DI)

**Dependency Injection** — это паттерн, при котором зависимости передаются объекту извне, а не создаются внутри.

```go
// ❌ Плохо — создаёт зависимость внутри
func NewService() *Service {
    db := sql.Open("postgres", "...")  // создаёт сам
    return &Service{db: db}
}

// ✅ Хорошо — получает зависимость извне
func NewService(repo UserRepository) *Service {
    return &Service{repo: repo}  // получает готовую
}
```

Теперь в тестах можно передать мок:

```go
func TestService(t *testing.T) {
    mockRepo := &MockUserRepository{}  // мок
    service := NewService(mockRepo)    // передаём мок
    // тестируем
}
```

### Виды моков

1. **Ручные моки** — вы сами пишете структуру с нужным поведением
2. **Моки с функциями** — структура с функциями, которые можно переопределить
3. **testify/mock** — библиотека с удобными проверками вызовов

### Что проверять в тестах с моками?

1. **Что метод вызван** — мок позволяет проверить, что `Save()` был вызван
2. **С какими аргументами** — можно проверить, какие данные передали
3. **Сколько раз вызван** — один раз, несколько или ни разу
4. **Порядок вызовов** — в какой последовательности вызывались методы

### Принцип: тестируем логику, не зависимости

Моки позволяют тестировать **бизнес-логику** вашего кода, а не работу базы данных или HTTP-клиента. База данных уже протестирована её разработчиками — вам нужно проверить, что ВАШ код правильно её использует.

---

## 📋 Синтаксис

### Определение интерфейса

```go
type Repository interface {
    GetByID(id int) (*User, error)
    Save(user *User) error
    Delete(id int) error
}
```

### Ручной мок

```go
type MockRepository struct {
    GetByIDFunc func(id int) (*User, error)
    SaveFunc    func(user *User) error
}

func (m *MockRepository) GetByID(id int) (*User, error) {
    return m.GetByIDFunc(id)
}
```

### testify/mock

```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetByID(id int) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

// В тесте:
mockRepo := new(MockRepository)
mockRepo.On("GetByID", 1).Return(&User{Name: "John"}, nil)
```

---

## 💻 Примеры кода

### Пример 1: Проблема без интерфейсов

```go
package user

import "database/sql"

// ❌ ПЛОХО — жёсткая зависимость от БД
type UserService struct {
    db *sql.DB
}

func (s *UserService) GetUser(id int) (*User, error) {
    var user User
    err := s.db.QueryRow("SELECT * FROM users WHERE id = ?", id).
        Scan(&user.ID, &user.Name, &user.Email)
    return &user, err
}

// Как протестировать без реальной БД? 🤔
```

### Пример 2: Решение с интерфейсами

```go
package user

// ✅ ХОРОШО — зависимость от интерфейса
type UserRepository interface {
    GetByID(id int) (*User, error)
    Save(user *User) error
    Delete(id int) error
}

type UserService struct {
    repo UserRepository  // интерфейс вместо конкретной реализации
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, ErrInvalidID
    }
    return s.repo.GetByID(id)
}

func (s *UserService) CreateUser(name, email string) (*User, error) {
    user := &User{Name: name, Email: email}
    err := s.repo.Save(user)
    return user, err
}
```

### Пример 3: Ручной мок

```go
package user

import (
    "errors"
    "testing"
)

// Мок репозитория
type MockUserRepository struct {
    users      map[int]*User
    saveError  error
    getByIDErr error
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users: make(map[int]*User),
    }
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
    if m.getByIDErr != nil {
        return nil, m.getByIDErr
    }
    user, ok := m.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}

func (m *MockUserRepository) Save(user *User) error {
    if m.saveError != nil {
        return m.saveError
    }
    if user.ID == 0 {
        user.ID = len(m.users) + 1
    }
    m.users[user.ID] = user
    return nil
}

func (m *MockUserRepository) Delete(id int) error {
    delete(m.users, id)
    return nil
}

// Тесты
func TestGetUser(t *testing.T) {
    // Arrange
    mockRepo := NewMockUserRepository()
    mockRepo.users[1] = &User{ID: 1, Name: "John", Email: "john@example.com"}
    
    service := NewUserService(mockRepo)
    
    // Act
    user, err := service.GetUser(1)
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.Name != "John" {
        t.Errorf("expected name John, got %s", user.Name)
    }
}

func TestGetUserNotFound(t *testing.T) {
    mockRepo := NewMockUserRepository()
    service := NewUserService(mockRepo)
    
    _, err := service.GetUser(999)
    
    if !errors.Is(err, ErrNotFound) {
        t.Errorf("expected ErrNotFound, got %v", err)
    }
}

func TestGetUserInvalidID(t *testing.T) {
    mockRepo := NewMockUserRepository()
    service := NewUserService(mockRepo)
    
    _, err := service.GetUser(-1)
    
    if !errors.Is(err, ErrInvalidID) {
        t.Errorf("expected ErrInvalidID, got %v", err)
    }
}
```

### Пример 4: Мок с функциями

```go
package payment

import (
    "testing"
)

type PaymentGateway interface {
    Charge(amount float64, cardToken string) (transactionID string, err error)
    Refund(transactionID string) error
}

type PaymentService struct {
    gateway PaymentGateway
}

func (s *PaymentService) ProcessPayment(amount float64, cardToken string) (string, error) {
    if amount <= 0 {
        return "", errors.New("invalid amount")
    }
    return s.gateway.Charge(amount, cardToken)
}

// Мок с функциями — гибкий подход
type MockPaymentGateway struct {
    ChargeFunc func(amount float64, cardToken string) (string, error)
    RefundFunc func(transactionID string) error
}

func (m *MockPaymentGateway) Charge(amount float64, cardToken string) (string, error) {
    if m.ChargeFunc != nil {
        return m.ChargeFunc(amount, cardToken)
    }
    return "tx_123", nil
}

func (m *MockPaymentGateway) Refund(transactionID string) error {
    if m.RefundFunc != nil {
        return m.RefundFunc(transactionID)
    }
    return nil
}

// Тесты
func TestProcessPaymentSuccess(t *testing.T) {
    mock := &MockPaymentGateway{
        ChargeFunc: func(amount float64, cardToken string) (string, error) {
            // Можем проверить аргументы
            if amount != 100.0 {
                t.Errorf("expected amount 100, got %f", amount)
            }
            return "tx_success", nil
        },
    }
    
    service := &PaymentService{gateway: mock}
    txID, err := service.ProcessPayment(100.0, "card_token")
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if txID != "tx_success" {
        t.Errorf("expected tx_success, got %s", txID)
    }
}

func TestProcessPaymentGatewayError(t *testing.T) {
    mock := &MockPaymentGateway{
        ChargeFunc: func(amount float64, cardToken string) (string, error) {
            return "", errors.New("gateway error")
        },
    }
    
    service := &PaymentService{gateway: mock}
    _, err := service.ProcessPayment(100.0, "card_token")
    
    if err == nil {
        t.Error("expected error, got nil")
    }
}
```

### Пример 5: testify/mock

```go
package order

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type Order struct {
    ID     int
    UserID int
    Amount float64
    Status string
}

type OrderRepository interface {
    Create(order *Order) error
    GetByID(id int) (*Order, error)
    UpdateStatus(id int, status string) error
}

// Мок с testify
type MockOrderRepository struct {
    mock.Mock
}

func (m *MockOrderRepository) Create(order *Order) error {
    args := m.Called(order)
    return args.Error(0)
}

func (m *MockOrderRepository) GetByID(id int) (*Order, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateStatus(id int, status string) error {
    args := m.Called(id, status)
    return args.Error(0)
}

// Сервис
type OrderService struct {
    repo OrderRepository
}

func (s *OrderService) PlaceOrder(userID int, amount float64) (*Order, error) {
    order := &Order{
        UserID: userID,
        Amount: amount,
        Status: "pending",
    }
    err := s.repo.Create(order)
    return order, err
}

func (s *OrderService) CompleteOrder(orderID int) error {
    order, err := s.repo.GetByID(orderID)
    if err != nil {
        return err
    }
    if order.Status != "pending" {
        return errors.New("order cannot be completed")
    }
    return s.repo.UpdateStatus(orderID, "completed")
}

// Тесты
func TestPlaceOrder(t *testing.T) {
    mockRepo := new(MockOrderRepository)
    
    // Настраиваем ожидание
    mockRepo.On("Create", mock.AnythingOfType("*order.Order")).Return(nil)
    
    service := &OrderService{repo: mockRepo}
    order, err := service.PlaceOrder(1, 99.99)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, order.UserID)
    assert.Equal(t, 99.99, order.Amount)
    assert.Equal(t, "pending", order.Status)
    
    // Проверяем, что метод был вызван
    mockRepo.AssertExpectations(t)
}

func TestCompleteOrder(t *testing.T) {
    mockRepo := new(MockOrderRepository)
    
    // Настраиваем цепочку вызовов
    mockRepo.On("GetByID", 1).Return(&Order{
        ID:     1,
        Status: "pending",
    }, nil)
    mockRepo.On("UpdateStatus", 1, "completed").Return(nil)
    
    service := &OrderService{repo: mockRepo}
    err := service.CompleteOrder(1)
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestCompleteOrderAlreadyCompleted(t *testing.T) {
    mockRepo := new(MockOrderRepository)
    
    mockRepo.On("GetByID", 1).Return(&Order{
        ID:     1,
        Status: "completed",  // уже завершён
    }, nil)
    
    service := &OrderService{repo: mockRepo}
    err := service.CompleteOrder(1)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "cannot be completed")
    
    // UpdateStatus НЕ должен быть вызван
    mockRepo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
}
```

### Пример 6: Мок HTTP клиента

```go
package api

import (
    "encoding/json"
    "io"
    "net/http"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type APIClient struct {
    client  HTTPClient
    baseURL string
}

func (c *APIClient) GetUser(id int) (*User, error) {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%d", c.baseURL, id), nil)
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
    
    var user User
    err = json.NewDecoder(resp.Body).Decode(&user)
    return &user, err
}

// Мок HTTP клиента
type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}

// Тест
func TestGetUser(t *testing.T) {
    mockClient := &MockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            // Проверяем запрос
            assert.Equal(t, "GET", req.Method)
            assert.Contains(t, req.URL.Path, "/users/1")
            
            // Возвращаем мок-ответ
            body := `{"id": 1, "name": "John", "email": "john@example.com"}`
            return &http.Response{
                StatusCode: 200,
                Body:       io.NopCloser(strings.NewReader(body)),
            }, nil
        },
    }
    
    client := &APIClient{
        client:  mockClient,
        baseURL: "https://api.example.com",
    }
    
    user, err := client.GetUser(1)
    
    require.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}
```

### Пример 7: Мок времени

```go
package scheduler

import (
    "testing"
    "time"
)

// Интерфейс для работы со временем
type Clock interface {
    Now() time.Time
}

// Реальная реализация
type RealClock struct{}

func (RealClock) Now() time.Time {
    return time.Now()
}

// Мок
type MockClock struct {
    CurrentTime time.Time
}

func (m MockClock) Now() time.Time {
    return m.CurrentTime
}

// Сервис
type Scheduler struct {
    clock Clock
}

func (s *Scheduler) IsWorkingHours() bool {
    hour := s.clock.Now().Hour()
    return hour >= 9 && hour < 18
}

// Тесты
func TestIsWorkingHours(t *testing.T) {
    tests := []struct {
        name     string
        time     time.Time
        expected bool
    }{
        {"morning 9am", time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), true},
        {"afternoon 2pm", time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC), true},
        {"evening 6pm", time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC), false},
        {"night 11pm", time.Date(2024, 1, 1, 23, 0, 0, 0, time.UTC), false},
        {"early morning 5am", time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC), false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            scheduler := &Scheduler{
                clock: MockClock{CurrentTime: tt.time},
            }
            
            result := scheduler.IsWorkingHours()
            
            if result != tt.expected {
                t.Errorf("at %v: got %v, want %v", tt.time, result, tt.expected)
            }
        })
    }
}
```

---

## ⚠️ Частые ошибки

### 1. Мокирование конкретных типов вместо интерфейсов

```go
// ❌ НЕВОЗМОЖНО замокать
type Service struct {
    db *sql.DB
}

// ✅ МОЖНО замокать
type Service struct {
    repo Repository  // интерфейс
}
```

### 2. Слишком большие интерфейсы

```go
// ❌ ПЛОХО — сложно мокать
type Repository interface {
    GetUser(id int) (*User, error)
    SaveUser(user *User) error
    DeleteUser(id int) error
    GetOrder(id int) (*Order, error)
    SaveOrder(order *Order) error
    // ... 20 методов
}

// ✅ ХОРОШО — маленькие интерфейсы
type UserRepository interface {
    GetByID(id int) (*User, error)
    Save(user *User) error
}

type OrderRepository interface {
    GetByID(id int) (*Order, error)
    Save(order *Order) error
}
```

### 3. Забыли проверить вызовы

```go
// ❌ ПЛОХО — не проверяем, был ли вызван метод
func TestSomething(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("Save", mock.Anything).Return(nil)
    
    service.DoSomething()
    // Забыли: mockRepo.AssertExpectations(t)
}

// ✅ ХОРОШО
func TestSomething(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("Save", mock.Anything).Return(nil)
    
    service.DoSomething()
    
    mockRepo.AssertExpectations(t)  // проверяем!
}
```

---

## 🏋️ Практические задания

### Задание 1: Мок Email сервиса

Создайте интерфейс `EmailSender` и мок для тестирования сервиса уведомлений.

**Начальный код:**
```go
// email.go
package notification

type EmailSender interface {
    Send(to, subject, body string) error
}

type NotificationService struct {
    emailer EmailSender
}

func NewNotificationService(e EmailSender) *NotificationService {
    return &NotificationService{emailer: e}
}

func (s *NotificationService) NotifyUser(email, message string) error {
    return s.emailer.Send(email, "Notification", message)
}
```

```go
// email_test.go
package notification

import "testing"

// Создайте MockEmailSender
type MockEmailSender struct {
    SendFunc func(to, subject, body string) error
    Calls    []EmailCall
}

type EmailCall struct {
    To, Subject, Body string
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    m.Calls = append(m.Calls, EmailCall{to, subject, body})
    if m.SendFunc != nil {
        return m.SendFunc(to, subject, body)
    }
    return nil
}

func TestNotifyUser(t *testing.T) {
    mock := &MockEmailSender{}
    service := NewNotificationService(mock)
    
    err := service.NotifyUser("user@example.com", "Hello!")
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    // Проверьте, что Send был вызван с правильными аргументами
    if len(mock.Calls) != 1 {
        t.Errorf("expected 1 call, got %d", len(mock.Calls))
    }
    // Добавьте проверки аргументов
}
```

**Ожидаемый вывод:**
```
=== RUN   TestNotifyUser
--- PASS: TestNotifyUser (0.00s)
PASS
```

**Баллы:** 15

### Задание 2: Мок кэша

Создайте интерфейс `Cache` и протестируйте сервис с кэшированием.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Определите и реализуйте интерфейс

func main() {
    // Ваш код здесь
    
}
```

```go
// cache_test.go
package user

import "testing"

type MockCache struct {
    data map[string]interface{}
}

func NewMockCache() *MockCache {
    return &MockCache{data: make(map[string]interface{})}
}

func (m *MockCache) Get(key string) (interface{}, bool) {
    v, ok := m.data[key]
    return v, ok
}

func (m *MockCache) Set(key string, value interface{}) {
    m.data[key] = value
}

type MockRepository struct {
    users map[int]*User
    calls int
}

func (m *MockRepository) GetByID(id int) (*User, error) {
    m.calls++
    if user, ok := m.users[id]; ok {
        return user, nil
    }
    return nil, errors.New("not found")
}

func TestUserServiceCaching(t *testing.T) {
    cache := NewMockCache()
    repo := &MockRepository{
        users: map[int]*User{1: {ID: 1, Name: "John"}},
    }
    service := &UserService{cache: cache, repo: repo}
    
    // Первый вызов — должен обратиться к repo
    user1, _ := service.GetUser(1)
    
    // Второй вызов — должен взять из кэша
    user2, _ := service.GetUser(1)
    
    // Проверьте, что repo был вызван только один раз
    if repo.calls != 1 {
        t.Errorf("expected 1 repo call, got %d", repo.calls)
    }
    
    // Проверьте, что оба результата одинаковы
    if user1.Name != user2.Name {
        t.Error("users should be equal")
    }
}
```

**Ожидаемый вывод:**
```
=== RUN   TestUserServiceCaching
--- PASS: TestUserServiceCaching (0.00s)
PASS
```

**Баллы:** 15

### Задание 3: testify/mock

Перепишите мок с использованием testify/mock.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Реализуйте решение согласно заданию

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
=== RUN   TestWithTestifyMock
--- PASS: TestWithTestifyMock (0.00s)
PASS
```

**Баллы:** 15

### Задание 4: Мок времени

Создайте интерфейс `Clock` и протестируйте функцию, зависящую от времени.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Создайте функцию согласно заданию
// TODO: Определите и реализуйте интерфейс

func main() {
    // Ваш код здесь
    
}
```

```go
// scheduler_test.go
package scheduler

import (
    "testing"
    "time"
)

type MockClock struct {
    Time time.Time
}

func (m MockClock) Now() time.Time { return m.Time }

func TestScheduler(t *testing.T) {
    tests := []struct {
        name      string
        time      time.Time
        isWeekend bool
        isWorking bool
    }{
        {
            name:      "monday_9am",
            time:      time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), // Monday
            isWeekend: false,
            isWorking: true,
        },
        // Добавьте ещё тестовые случаи
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            s := &Scheduler{clock: MockClock{Time: tt.time}}
            
            if s.IsWeekend() != tt.isWeekend {
                t.Errorf("IsWeekend() = %v, want %v", s.IsWeekend(), tt.isWeekend)
            }
            if s.IsWorkingHours() != tt.isWorking {
                t.Errorf("IsWorkingHours() = %v, want %v", s.IsWorkingHours(), tt.isWorking)
            }
        })
    }
}
```

**Ожидаемый вывод:**
```
=== RUN   TestScheduler
=== RUN   TestScheduler/monday_9am
=== RUN   TestScheduler/saturday_noon
=== RUN   TestScheduler/friday_8pm
--- PASS: TestScheduler (0.00s)
PASS
```

**Баллы:** 10

### Задание 5: Проверка вызовов с аргументами

Убедитесь, что метод был вызван с конкретными аргументами определённое количество раз.

**Начальный код:**
```go
package main

import "fmt"

// TODO: Добавьте методы

func main() {
    // Ваш код здесь
    
}
```

**Ожидаемый вывод:**
```
=== RUN   TestServiceLogging
--- PASS: TestServiceLogging (0.00s)
PASS
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Testify Mock](https://pkg.go.dev/github.com/stretchr/testify/mock)
- [Go Interface Mocking](https://go.dev/doc/effective_go#interface_and_methods)
- [GoMock](https://github.com/golang/mock)
