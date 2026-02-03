# Введение в Gin Framework

---

## 💡 Ключевые идеи

1. **Gin** — самый популярный веб-фреймворк для Go
2. **Высокая производительность** — использует httprouter
3. **Middleware** — гибкая система промежуточных обработчиков
4. **JSON binding** — автоматическая сериализация/десериализация
5. **Валидация** — встроенная валидация запросов

### Сравнение фреймворков

| Фреймворк | Производительность | Функциональность | Популярность |
|-----------|-------------------|------------------|--------------|
| Gin | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | #1 |
| Echo | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | #2 |
| Fiber | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | #3 |
| Chi | ⭐⭐⭐⭐ | ⭐⭐⭐ | #4 |

---

## 📖 Теория

### Что такое Gin?

**Gin** — это высокопроизводительный веб-фреймворк для Go, написанный на основе httprouter. Он предоставляет удобный API для создания RESTful сервисов с минимальным количеством кода.

### Почему Gin так популярен?

1. **Скорость** — до 40 раз быстрее стандартного net/http благодаря radix tree роутингу
2. **Простота** — интуитивный API, легко начать
3. **Middleware** — богатая экосистема готовых middleware
4. **JSON** — встроенная сериализация/десериализация
5. **Валидация** — автоматическая валидация входных данных

### net/http vs Gin

**Стандартный net/http:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Парсинг JSON вручную
    var data MyStruct
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        w.WriteHeader(400)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    
    // Получение параметров вручную
    id := r.URL.Query().Get("id")
    
    // Ответ
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

**С Gin:**
```go
func handler(c *gin.Context) {
    var data MyStruct
    if err := c.ShouldBindJSON(&data); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    id := c.Query("id")
    
    c.JSON(200, response)
}
```

Gin убирает boilerplate и делает код чище.

### Основные концепции Gin

**1. gin.Engine** — основной объект приложения:
```go
r := gin.Default()  // с Logger и Recovery
r := gin.New()      // без middleware
```

**2. gin.Context** — контекст запроса:
```go
func handler(c *gin.Context) {
    c.JSON(200, data)     // Ответ JSON
    c.String(200, "OK")   // Ответ строкой
    c.HTML(200, "index.tmpl", data)  // HTML
    
    c.Query("name")       // ?name=value
    c.Param("id")         // /users/:id
    c.PostForm("field")   // form field
    
    c.BindJSON(&data)     // Парсинг JSON
    c.Set("key", value)   // Сохранить в контексте
    c.Get("key")          // Получить из контекста
}
```

**3. Роутинг:**
```go
r.GET("/users", handler)
r.POST("/users", handler)
r.PUT("/users/:id", handler)
r.DELETE("/users/:id", handler)

// Группы
api := r.Group("/api")
{
    api.GET("/users", getUsers)
    api.POST("/users", createUser)
}
```

**4. Middleware:**
```go
// Глобальный middleware
r.Use(gin.Logger())

// Для группы
authorized := r.Group("/admin")
authorized.Use(AuthRequired())

// Для конкретного route
r.GET("/protected", AuthRequired(), handler)
```

### Binding и валидация

Gin использует теги struct для биндинга и валидации:

```go
type CreateUserInput struct {
    Name     string `json:"name" binding:"required,min=2,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"gte=0,lte=150"`
    Password string `json:"password" binding:"required,min=8"`
}

func createUser(c *gin.Context) {
    var input CreateUserInput
    
    // Автоматическая валидация
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // input валиден!
}
```

### Когда использовать Gin?

✅ **Используйте когда:**
- Создаёте REST API
- Важна производительность
- Нужна валидация входных данных
- Хотите быстро начать

❌ **Рассмотрите альтернативы когда:**
- Нужна максимальная совместимость со стандартной библиотекой (Chi)
- Хотите fasthttp вместо net/http (Fiber)

---

## 📋 Синтаксис

### Установка

```bash
go get -u github.com/gin-gonic/gin
```

### Базовая структура

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()  // с Logger и Recovery middleware
    
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    r.Run(":8080")
}
```

---

## 💻 Примеры кода

### Пример 1: Полный REST API

```go
package main

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=150"`
}

var users = []User{
    {ID: 1, Name: "Alice", Email: "alice@example.com", Age: 30},
    {ID: 2, Name: "Bob", Email: "bob@example.com", Age: 25},
}

func main() {
    r := gin.Default()
    
    // Routes
    api := r.Group("/api/v1")
    {
        api.GET("/users", getUsers)
        api.GET("/users/:id", getUser)
        api.POST("/users", createUser)
        api.PUT("/users/:id", updateUser)
        api.DELETE("/users/:id", deleteUser)
    }
    
    r.Run(":8080")
}

func getUsers(c *gin.Context) {
    c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    
    for _, user := range users {
        if user.ID == id {
            c.JSON(http.StatusOK, user)
            return
        }
    }
    
    c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}

func createUser(c *gin.Context) {
    var user User
    
    // Binding с валидацией
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user.ID = len(users) + 1
    users = append(users, user)
    
    c.JSON(http.StatusCreated, user)
}

func updateUser(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    
    var input User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    for i, user := range users {
        if user.ID == id {
            input.ID = id
            users[i] = input
            c.JSON(http.StatusOK, input)
            return
        }
    }
    
    c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}

func deleteUser(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    
    for i, user := range users {
        if user.ID == id {
            users = append(users[:i], users[i+1:]...)
            c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
            return
        }
    }
    
    c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}
```

### Пример 2: Middleware

```go
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
)

// Custom Logger Middleware
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        // Выполняем следующий handler
        c.Next()
        
        // После выполнения
        duration := time.Since(start)
        status := c.Writer.Status()
        
        log.Printf("[%s] %s %d %v", c.Request.Method, path, status, duration)
    }
}

// Auth Middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "authorization header required",
            })
            return
        }
        
        // Валидация токена (упрощённо)
        if token != "Bearer secret-token" {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "invalid token",
            })
            return
        }
        
        // Сохраняем данные пользователя в context
        c.Set("userID", 123)
        c.Set("role", "admin")
        
        c.Next()
    }
}

// Rate Limiter Middleware (упрощённый)
func RateLimiter(rps int) gin.HandlerFunc {
    limiter := make(chan struct{}, rps)
    
    // Наполняем канал
    go func() {
        ticker := time.NewTicker(time.Second / time.Duration(rps))
        for range ticker.C {
            select {
            case limiter <- struct{}{}:
            default:
            }
        }
    }()
    
    return func(c *gin.Context) {
        select {
        case <-limiter:
            c.Next()
        default:
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
            })
        }
    }
}

func main() {
    r := gin.New()  // без стандартных middleware
    
    // Глобальные middleware
    r.Use(gin.Recovery())
    r.Use(Logger())
    
    // Публичные endpoints
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // Защищённые endpoints
    protected := r.Group("/api")
    protected.Use(AuthMiddleware())
    {
        protected.GET("/profile", func(c *gin.Context) {
            userID := c.GetInt("userID")
            role := c.GetString("role")
            
            c.JSON(200, gin.H{
                "userID": userID,
                "role":   role,
            })
        })
    }
    
    r.Run(":8080")
}
```

### Пример 3: Binding и валидация

```go
package main

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/validator/v10"
)

// Custom validator
func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    return len(username) >= 3 && len(username) <= 20
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,username"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Age      int    `json:"age" binding:"gte=18,lte=120"`
}

type QueryParams struct {
    Page     int    `form:"page" binding:"gte=1"`
    PageSize int    `form:"page_size" binding:"gte=1,lte=100"`
    Sort     string `form:"sort" binding:"omitempty,oneof=asc desc"`
}

func main() {
    r := gin.Default()
    
    // Регистрируем custom validator
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("username", validateUsername)
    }
    
    // JSON binding
    r.POST("/register", func(c *gin.Context) {
        var req RegisterRequest
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "validation failed",
                "details": err.Error(),
            })
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message":  "user registered",
            "username": req.Username,
        })
    })
    
    // Query binding
    r.GET("/users", func(c *gin.Context) {
        var query QueryParams
        query.Page = 1      // default
        query.PageSize = 10 // default
        
        if err := c.ShouldBindQuery(&query); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "page":      query.Page,
            "page_size": query.PageSize,
            "sort":      query.Sort,
        })
    })
    
    // URI binding
    r.GET("/users/:id", func(c *gin.Context) {
        type URI struct {
            ID int `uri:"id" binding:"required,gte=1"`
        }
        
        var uri URI
        if err := c.ShouldBindUri(&uri); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{"id": uri.ID})
    })
    
    r.Run(":8080")
}
```

### Пример 4: Группировка и версионирование

```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    // API v1
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", getV1Users)
        v1.POST("/users", createV1User)
        
        // Nested group
        admin := v1.Group("/admin")
        admin.Use(AdminAuthMiddleware())
        {
            admin.GET("/stats", getStats)
            admin.DELETE("/users/:id", deleteUser)
        }
    }
    
    // API v2
    v2 := r.Group("/api/v2")
    {
        v2.GET("/users", getV2Users)  // новая версия
    }
    
    r.Run(":8080")
}
```

### Пример 5: Обработка ошибок

```go
package main

import (
    "errors"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

// Custom errors
var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrBadRequest   = errors.New("bad request")
)

// APIError — структура ошибки API
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

// Error handler middleware
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        // Проверяем ошибки после выполнения handler
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            var status int
            var message string
            
            switch {
            case errors.Is(err, ErrNotFound):
                status = http.StatusNotFound
                message = "Resource not found"
            case errors.Is(err, ErrUnauthorized):
                status = http.StatusUnauthorized
                message = "Unauthorized"
            case errors.Is(err, ErrForbidden):
                status = http.StatusForbidden
                message = "Forbidden"
            case errors.Is(err, ErrBadRequest):
                status = http.StatusBadRequest
                message = err.Error()
            default:
                status = http.StatusInternalServerError
                message = "Internal server error"
            }
            
            c.JSON(status, APIError{
                Code:    status,
                Message: message,
            })
        }
    }
}

func main() {
    r := gin.Default()
    r.Use(ErrorHandler())
    
    r.GET("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        
        if id == "0" {
            c.Error(ErrNotFound)
            return
        }
        
        c.JSON(http.StatusOK, gin.H{"id": id})
    })
    
    r.Run(":8080")
}
```

### Пример 6: Загрузка файлов

```go
package main

import (
    "fmt"
    "net/http"
    "path/filepath"
    
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    // Лимит размера файла
    r.MaxMultipartMemory = 8 << 20 // 8 MB
    
    // Один файл
    r.POST("/upload", func(c *gin.Context) {
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        // Валидация типа файла
        ext := filepath.Ext(file.Filename)
        if ext != ".jpg" && ext != ".png" && ext != ".gif" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "only images allowed"})
            return
        }
        
        // Сохранение
        dst := fmt.Sprintf("./uploads/%s", file.Filename)
        if err := c.SaveUploadedFile(file, dst); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message":  "file uploaded",
            "filename": file.Filename,
            "size":     file.Size,
        })
    })
    
    // Множественные файлы
    r.POST("/uploads", func(c *gin.Context) {
        form, err := c.MultipartForm()
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        files := form.File["files"]
        
        for _, file := range files {
            dst := fmt.Sprintf("./uploads/%s", file.Filename)
            c.SaveUploadedFile(file, dst)
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("%d files uploaded", len(files)),
        })
    })
    
    r.Run(":8080")
}
```

---

## ⚠️ Частые ошибки

### 1. Не обрабатывают ошибки binding

```go
// ❌ ПЛОХО
c.ShouldBindJSON(&req)  // игнорируем ошибку

// ✅ ХОРОШО
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### 2. Забывают return после ответа

```go
// ❌ ПЛОХО
if err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    // продолжает выполнение!
}

// ✅ ХОРОШО
if err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

---

## 🏋️ Практические задания

### Задание 1: CRUD API для Product

Создайте полный CRUD API.

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

### Задание 2: Custom Middleware

Напишите middleware для логирования и аутентификации.

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

### Задание 3: Custom Validator

Создайте custom validator для телефонных номеров.

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

### Задание 4: Error Handling

Создайте централизованную обработку ошибок.

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

### Задание 5: Группировка и версионирование API

Создайте структурированный API с версионированием.

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

---

## 🔗 Полезные ссылки

- [Gin Documentation](https://gin-gonic.com/docs/)
- [Gin GitHub](https://github.com/gin-gonic/gin)
- [Validator v10](https://github.com/go-playground/validator)
