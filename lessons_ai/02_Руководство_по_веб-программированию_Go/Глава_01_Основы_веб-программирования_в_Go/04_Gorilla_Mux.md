# Маршрутизация и gorilla/mux

## 💡 Ключевые идеи

1. **gorilla/mux** — мощный маршрутизатор для Go
2. **Параметры URL** — `/users/{id}`
3. **Регулярные выражения** — `/users/{id:[0-9]+}`
4. **HTTP методы** — `.Methods("GET", "POST")`
5. **Middleware** — промежуточные обработчики

---

## 📋 Синтаксис

### Установка

```bash
go get -u github.com/gorilla/mux
```

### Создание роутера

```go
router := mux.NewRouter()
router.HandleFunc(pattern, handler)
```

### Параметры URL

```go
// {name} — параметр
"/users/{id}"

// {name:regex} — с регулярным выражением
"/users/{id:[0-9]+}"

// Получение параметра
vars := mux.Vars(r)
id := vars["id"]
```

---

## 💻 Примеры кода

### Базовое использование

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Home Page")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "About Page")
}

func main() {
    router := mux.NewRouter()
    
    router.HandleFunc("/", homeHandler)
    router.HandleFunc("/about", aboutHandler)
    
    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", router)
}
```

### Параметры в URL

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    
    fmt.Fprintf(w, "User ID: %s", id)
}

func main() {
    router := mux.NewRouter()
    
    router.HandleFunc("/users/{id}", userHandler)
    
    http.ListenAndServe(":8080", router)
}
```

### Регулярные выражения в параметрах

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
)

func productHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    
    category := vars["category"]
    id := vars["id"]
    
    fmt.Fprintf(w, "Category: %s, Product ID: %s", category, id)
}

func main() {
    router := mux.NewRouter()
    
    // id должен быть числом
    router.HandleFunc("/products/{category}/{id:[0-9]+}", productHandler)
    
    http.ListenAndServe(":8080", router)
}
```

### Ограничение по HTTP методам

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "GET: List all users")
}

func createUser(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "POST: Create user")
}

func getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fmt.Fprintf(w, "GET: User %s", vars["id"])
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fmt.Fprintf(w, "PUT: Update user %s", vars["id"])
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fmt.Fprintf(w, "DELETE: Delete user %s", vars["id"])
}

func main() {
    router := mux.NewRouter()
    
    router.HandleFunc("/users", getUsers).Methods("GET")
    router.HandleFunc("/users", createUser).Methods("POST")
    router.HandleFunc("/users/{id:[0-9]+}", getUser).Methods("GET")
    router.HandleFunc("/users/{id:[0-9]+}", updateUser).Methods("PUT")
    router.HandleFunc("/users/{id:[0-9]+}", deleteUser).Methods("DELETE")
    
    http.ListenAndServe(":8080", router)
}
```

### Subrouter (подроутеры)

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()
    
    // API v1
    apiV1 := router.PathPrefix("/api/v1").Subrouter()
    apiV1.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "API v1: Users")
    })
    apiV1.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "API v1: Products")
    })
    
    // API v2
    apiV2 := router.PathPrefix("/api/v2").Subrouter()
    apiV2.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "API v2: Users")
    })
    
    http.ListenAndServe(":8080", router)
}
```

### Middleware

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main() {
    router := mux.NewRouter()
    
    // Глобальный middleware
    router.Use(loggingMiddleware)
    
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Home")
    })
    
    // Защищённый роут
    api := router.PathPrefix("/api").Subrouter()
    api.Use(authMiddleware)
    api.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Protected data")
    })
    
    http.ListenAndServe(":8080", router)
}
```

### Query параметры

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
    // Query параметры из URL
    query := r.URL.Query().Get("q")
    page := r.URL.Query().Get("page")
    
    if page == "" {
        page = "1"
    }
    
    fmt.Fprintf(w, "Search: %s, Page: %s", query, page)
}

func main() {
    router := mux.NewRouter()
    
    // URL: /search?q=golang&page=2
    router.HandleFunc("/search", searchHandler).
        Queries("q", "{q}").
        Methods("GET")
    
    http.ListenAndServe(":8080", router)
}
```

### Статические файлы с mux

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()
    
    // Статические файлы
    router.PathPrefix("/static/").Handler(
        http.StripPrefix("/static/", 
            http.FileServer(http.Dir("static"))))
    
    // Динамические маршруты
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Home")
    })
    
    router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        fmt.Fprintf(w, "User: %s", vars["id"])
    })
    
    http.ListenAndServe(":8080", router)
}
```

### Полный REST API пример

```go
package main

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/gorilla/mux"
)

type Product struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price int    `json:"price"`
}

var products = []Product{
    {ID: 1, Name: "iPhone", Price: 999},
    {ID: 2, Name: "MacBook", Price: 1999},
}

func getProducts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := strconv.Atoi(vars["id"])
    
    for _, p := range products {
        if p.ID == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(p)
            return
        }
    }
    
    http.NotFound(w, r)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
    var p Product
    json.NewDecoder(r.Body).Decode(&p)
    
    p.ID = len(products) + 1
    products = append(products, p)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(p)
}

func main() {
    router := mux.NewRouter()
    
    router.HandleFunc("/products", getProducts).Methods("GET")
    router.HandleFunc("/products/{id:[0-9]+}", getProduct).Methods("GET")
    router.HandleFunc("/products", createProduct).Methods("POST")
    
    http.ListenAndServe(":8080", router)
}
```

---

## ⚠️ Частые ошибки

### 1. Забыли передать router в ListenAndServe

```go
// ❌ Используется default mux
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)

// ✅ Передаём gorilla router
router := mux.NewRouter()
router.HandleFunc("/", handler)
http.ListenAndServe(":8080", router)
```

### 2. Порядок маршрутов

```go
// ❌ Более общий маршрут перехватывает
router.HandleFunc("/users/{id}", specificHandler)
router.HandleFunc("/users/new", newHandler)  // не сработает!

// ✅ Более специфичные маршруты первыми
router.HandleFunc("/users/new", newHandler)
router.HandleFunc("/users/{id}", specificHandler)
```

---

## 📝 Практика

### Задача 1: REST API
Полный CRUD API с gorilla/mux.

### Задача 2: Nested routes
Вложенные маршруты с subrouters.

### Задача 3: Auth middleware
Middleware для авторизации.

### Задача 4: Rate limiting
Ограничение запросов.

### Задача 5: CORS
Middleware для CORS.

### Задача 6: Validation
Валидация параметров URL.

### Задача 7: Error handling
Централизованная обработка ошибок.

### Задача 8: Documentation
Автодокументация API.
