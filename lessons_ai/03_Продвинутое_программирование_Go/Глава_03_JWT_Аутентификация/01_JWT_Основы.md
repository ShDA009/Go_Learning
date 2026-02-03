# JWT Аутентификация

---

## 💡 Ключевые идеи

1. **JWT** — JSON Web Token для stateless аутентификации
2. **Структура** — Header.Payload.Signature
3. **Access Token** — короткоживущий токен доступа
4. **Refresh Token** — долгоживущий токен для обновления
5. **Claims** — данные в payload (user ID, роли, expiration)

### Структура JWT

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.     <- Header (base64)
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4ifQ. <- Payload (base64)
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c   <- Signature
```

---

## 📖 Теория

### Что такое JWT?

**JWT (JSON Web Token)** — это открытый стандарт (RFC 7519) для безопасной передачи информации между сторонами в виде JSON-объекта. Эта информация может быть проверена и является доверенной, так как имеет цифровую подпись.

### Зачем нужен JWT?

В классической сессионной аутентификации:
1. Пользователь логинится
2. Сервер создаёт сессию и сохраняет в БД
3. Сервер отправляет session_id в cookie
4. При каждом запросе сервер проверяет session_id в БД

**Проблемы:**
- Нагрузка на БД при каждом запросе
- Сложности с масштабированием (нужна общая БД сессий)
- Не подходит для микросервисов

**JWT-аутентификация:**
1. Пользователь логинится
2. Сервер создаёт JWT с данными пользователя
3. Клиент хранит JWT (localStorage, cookie)
4. При каждом запросе клиент отправляет JWT
5. Сервер **проверяет подпись** без обращения к БД

### Структура JWT

JWT состоит из трёх частей, разделённых точками:

```
xxxxx.yyyyy.zzzzz
  ↑      ↑      ↑
Header  Payload  Signature
```

**1. Header (Заголовок):**
```json
{
  "alg": "HS256",  // алгоритм подписи
  "typ": "JWT"     // тип токена
}
```

**2. Payload (Полезная нагрузка):**
```json
{
  "sub": "1234567890",    // subject (ID пользователя)
  "name": "John Doe",     // пользовательские данные
  "iat": 1516239022,      // issued at (время создания)
  "exp": 1516242622       // expiration (время истечения)
}
```

**3. Signature (Подпись):**
```
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret
)
```

### Стандартные Claims

| Claim | Название | Описание |
|-------|----------|----------|
| iss | Issuer | Кто выдал токен |
| sub | Subject | Кому выдан (ID пользователя) |
| aud | Audience | Для кого предназначен |
| exp | Expiration | Когда истекает |
| iat | Issued At | Когда создан |
| nbf | Not Before | Не использовать до |
| jti | JWT ID | Уникальный ID токена |

### Access Token vs Refresh Token

**Access Token:**
- Короткий срок жизни (5-30 минут)
- Отправляется с каждым запросом
- Содержит данные пользователя (ID, роль)
- При компрометации ущерб ограничен временем жизни

**Refresh Token:**
- Долгий срок жизни (дни, недели)
- Используется только для получения нового Access Token
- Хранится безопаснее (httpOnly cookie)
- Можно отозвать (revoke) на сервере

**Поток:**
```
1. Login → Access Token (15 мин) + Refresh Token (7 дней)
2. Запросы с Access Token
3. Access Token истёк → POST /refresh с Refresh Token
4. Получили новый Access Token
```

### Алгоритмы подписи

**Симметричные (один ключ):**
- `HS256`, `HS384`, `HS512` — HMAC с SHA
- Один секретный ключ для подписи и проверки
- Проще, подходит для монолита

**Асимметричные (пара ключей):**
- `RS256`, `RS384`, `RS512` — RSA
- `ES256`, `ES384`, `ES512` — ECDSA
- Приватный ключ для подписи, публичный для проверки
- Подходит для микросервисов (можно раздать публичный ключ)

### Безопасность JWT

**✅ Правильно:**
- Короткий срок жизни Access Token
- HTTPS всегда
- Секретный ключ достаточной длины (минимум 256 бит)
- Валидация всех claims (exp, iss, aud)

**❌ Неправильно:**
- Хранить чувствительные данные в payload (пароли!)
- Очень долгий срок жизни токена
- Слабый секретный ключ
- Передача токена в URL

### JWT — это НЕ шифрование!

Payload в JWT закодирован в Base64, но **НЕ зашифрован**. Любой может прочитать содержимое:

```bash
echo "eyJzdWIiOiIxMjM0NTY3ODkwIn0" | base64 -d
# {"sub":"1234567890"}
```

Подпись гарантирует только то, что данные не изменялись. Не храните секреты в JWT!

---

## 💻 Примеры кода

### Полная реализация JWT аутентификации

```go
package main

import (
    "errors"
    "net/http"
    "strings"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

var (
    accessSecret  = []byte("access-secret-key")
    refreshSecret = []byte("refresh-secret-key")
    
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken       = errors.New("invalid token")
    ErrTokenExpired       = errors.New("token expired")
)

// User model
type User struct {
    ID       int64  `json:"id"`
    Email    string `json:"email"`
    Password string `json:"-"`
    Role     string `json:"role"`
}

// Claims для access token
type AccessClaims struct {
    UserID int64  `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// Claims для refresh token
type RefreshClaims struct {
    UserID int64 `json:"user_id"`
    jwt.RegisteredClaims
}

// TokenPair — пара токенов
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
}

// Auth Service
type AuthService struct {
    users map[string]*User // в реальности — база данных
}

func NewAuthService() *AuthService {
    // Создаём тестового пользователя
    hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    
    return &AuthService{
        users: map[string]*User{
            "user@example.com": {
                ID:       1,
                Email:    "user@example.com",
                Password: string(hash),
                Role:     "user",
            },
        },
    }
}

// Login — аутентификация и генерация токенов
func (s *AuthService) Login(email, password string) (*TokenPair, error) {
    user, ok := s.users[email]
    if !ok {
        return nil, ErrInvalidCredentials
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, ErrInvalidCredentials
    }
    
    return s.generateTokenPair(user)
}

// RefreshTokens — обновление токенов
func (s *AuthService) RefreshTokens(refreshToken string) (*TokenPair, error) {
    claims, err := s.ValidateRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }
    
    // Находим пользователя
    for _, user := range s.users {
        if user.ID == claims.UserID {
            return s.generateTokenPair(user)
        }
    }
    
    return nil, ErrInvalidToken
}

// ValidateAccessToken — валидация access token
func (s *AuthService) ValidateAccessToken(tokenString string) (*AccessClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return accessSecret, nil
    })
    
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrTokenExpired
        }
        return nil, ErrInvalidToken
    }
    
    if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, ErrInvalidToken
}

// ValidateRefreshToken — валидация refresh token
func (s *AuthService) ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
        return refreshSecret, nil
    })
    
    if err != nil {
        return nil, ErrInvalidToken
    }
    
    if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, ErrInvalidToken
}

func (s *AuthService) generateTokenPair(user *User) (*TokenPair, error) {
    // Access token (15 минут)
    accessExpiration := time.Now().Add(15 * time.Minute)
    accessClaims := AccessClaims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(accessExpiration),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   user.Email,
        },
    }
    
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    accessString, err := accessToken.SignedString(accessSecret)
    if err != nil {
        return nil, err
    }
    
    // Refresh token (7 дней)
    refreshExpiration := time.Now().Add(7 * 24 * time.Hour)
    refreshClaims := RefreshClaims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(refreshExpiration),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshString, err := refreshToken.SignedString(refreshSecret)
    if err != nil {
        return nil, err
    }
    
    return &TokenPair{
        AccessToken:  accessString,
        RefreshToken: refreshString,
        ExpiresIn:    int64(accessExpiration.Sub(time.Now()).Seconds()),
    }, nil
}

// Auth Middleware
func AuthMiddleware(authService *AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "authorization header required",
            })
            return
        }
        
        // Формат: "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "invalid authorization format",
            })
            return
        }
        
        claims, err := authService.ValidateAccessToken(parts[1])
        if err != nil {
            status := http.StatusUnauthorized
            message := "invalid token"
            
            if errors.Is(err, ErrTokenExpired) {
                message = "token expired"
            }
            
            c.AbortWithStatusJSON(status, gin.H{"error": message})
            return
        }
        
        // Сохраняем claims в context
        c.Set("userID", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)
        
        c.Next()
    }
}

// Role Middleware
func RoleMiddleware(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("role")
        
        for _, role := range roles {
            if userRole == role {
                c.Next()
                return
            }
        }
        
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "error": "insufficient permissions",
        })
    }
}

// Handlers
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

func main() {
    authService := NewAuthService()
    
    r := gin.Default()
    
    // Public routes
    r.POST("/auth/login", func(c *gin.Context) {
        var req LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        tokens, err := authService.Login(req.Email, req.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, tokens)
    })
    
    r.POST("/auth/refresh", func(c *gin.Context) {
        var req RefreshRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        tokens, err := authService.RefreshTokens(req.RefreshToken)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, tokens)
    })
    
    // Protected routes
    protected := r.Group("/api")
    protected.Use(AuthMiddleware(authService))
    {
        protected.GET("/profile", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "user_id": c.GetInt64("userID"),
                "email":   c.GetString("email"),
                "role":    c.GetString("role"),
            })
        })
        
        // Admin only
        admin := protected.Group("/admin")
        admin.Use(RoleMiddleware("admin"))
        {
            admin.GET("/stats", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{"stats": "admin data"})
            })
        }
    }
    
    r.Run(":8080")
}
```

---

## 🏋️ Практические задания

### Задание 1: Генерация и проверка JWT

Реализуйте базовую JWT аутентификацию.

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

### Задание 2: Refresh Token

Реализуйте систему с access и refresh токенами.

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

### Задание 3: Token Blacklist

Реализуйте blacklist для отозванных токенов.

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

### Задание 4: Login/Logout Handlers

Создайте полный flow аутентификации.

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

### Задание 5: Role-Based Access

Реализуйте middleware для проверки ролей.

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

- [JWT.io](https://jwt.io/)
- [golang-jwt](https://github.com/golang-jwt/jwt)
