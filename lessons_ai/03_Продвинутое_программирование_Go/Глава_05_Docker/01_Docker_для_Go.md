# Docker для Go приложений

---

## 💡 Ключевые идеи

1. **Multi-stage builds** — минимальные образы
2. **Scratch/Distroless** — безопасные базовые образы
3. **Docker Compose** — оркестрация локального окружения
4. **Best Practices** — кэширование, безопасность

---

## 📖 Теория

### Зачем Docker для Go?

Go компилирует код в один бинарный файл без зависимостей. Казалось бы, зачем Docker? Причины:

1. **Единообразное окружение** — "работает на моей машине" → работает везде
2. **Зависимости** — БД, Redis, очереди запускаются одной командой
3. **Изоляция** — приложение работает в своём контейнере
4. **Масштабирование** — легко запустить несколько экземпляров
5. **CI/CD** — стандартный способ доставки в Kubernetes

### Multi-stage builds

Классическая проблема: образ с компилятором Go весит ~1GB, а нам нужен только бинарник.

**Решение: Multi-stage build**
```dockerfile
# Stage 1: Сборка (большой образ с компилятором)
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Stage 2: Запуск (минимальный образ)
FROM alpine:latest
COPY --from=builder /app/main /main
CMD ["/main"]
```

Результат:
- Сборочный образ: ~1GB
- Финальный образ: ~15MB

### Scratch vs Alpine vs Distroless

**1. scratch** — пустой образ (0 байт):
```dockerfile
FROM scratch
COPY main /main
CMD ["/main"]
```
- ✅ Минимальный размер (~10MB)
- ❌ Нет shell, нельзя отлаживать
- ❌ Нет сертификатов (нужно копировать)

**2. alpine** — минимальный Linux (~5MB):
```dockerfile
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY main /main
CMD ["/main"]
```
- ✅ Есть shell и утилиты
- ✅ Легко отлаживать
- ⚠️ Использует musl вместо glibc

**3. distroless** — от Google:
```dockerfile
FROM gcr.io/distroless/static-debian12
COPY main /main
CMD ["/main"]
```
- ✅ Безопаснее alpine (меньше attack surface)
- ✅ Есть сертификаты
- ❌ Нет shell

### CGO_ENABLED=0

Go может использовать C-библиотеки через CGO. Для статического бинарника нужно отключить:

```dockerfile
RUN CGO_ENABLED=0 go build -o main .
```

Если CGO включён, бинарник будет динамически линкован с glibc и не запустится в scratch/alpine.

### Флаги сборки

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=1.0.0" \
    -o main ./cmd/server
```

- `CGO_ENABLED=0` — статическая линковка
- `GOOS=linux GOARCH=amd64` — целевая платформа
- `-ldflags="-w -s"` — убрать отладочную информацию (меньше размер)
- `-X main.Version=1.0.0` — встроить версию в бинарник

### Docker Compose для разработки

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db
      - redis
    environment:
      - DATABASE_URL=postgres://...
  
  db:
    image: postgres:15
    volumes:
      - postgres_data:/var/lib/postgresql/data
  
  redis:
    image: redis:7-alpine
```

Одна команда запускает всё окружение:
```bash
docker compose up -d
```

### Best Practices

1. **Кэшируйте зависимости:**
```dockerfile
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

2. **Не запускайте от root:**
```dockerfile
USER nobody:nobody
```

3. **Используйте .dockerignore:**
```
.git
*.md
Dockerfile
docker-compose.yml
```

4. **Healthcheck:**
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -q --spider http://localhost:8080/health || exit 1
```

---

## 💻 Примеры кода

### Пример 1: Простой Dockerfile

```dockerfile
# Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Копируем go.mod и go.sum для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем
RUN go build -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]
```

### Пример 2: Multi-stage Build (рекомендуется)

```dockerfile
# Dockerfile
# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для CGO (если нужно)
RUN apk add --no-cache gcc musl-dev

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код
COPY . .

# Собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=1.0.0" \
    -o /app/server \
    ./cmd/server

# Stage 2: Runtime
FROM alpine:3.19

# Добавляем CA сертификаты для HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Создаём непривилегированного пользователя
RUN adduser -D -g '' appuser

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/server .

# Копируем конфиги/миграции если нужно
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

# Меняем владельца
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./server"]
```

### Пример 3: Минимальный образ со Scratch

```dockerfile
# Dockerfile.scratch
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Статическая сборка без CGO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o server ./cmd/server

# Используем пустой образ
FROM scratch

# Копируем CA сертификаты
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копируем бинарник
COPY --from=builder /app/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
```

### Пример 4: Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://postgres:postgres@db:5432/app?sslmode=disable
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=your-secret-key
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_started
    restart: unless-stopped
    networks:
      - app-network

  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: app
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./migrations/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    networks:
      - app-network

  migrate:
    image: migrate/migrate
    volumes:
      - ./migrations:/migrations
    command: ["-path", "/migrations", "-database", "postgres://postgres:postgres@db:5432/app?sslmode=disable", "up"]
    depends_on:
      db:
        condition: service_healthy
    networks:
      - app-network

volumes:
  postgres-data:
  redis-data:

networks:
  app-network:
    driver: bridge
```

### Пример 5: .dockerignore

```
# .dockerignore
# Git
.git
.gitignore

# IDE
.idea
.vscode
*.swp

# Build artifacts
bin/
dist/
*.exe

# Test files
*_test.go
**/*_test.go
coverage.out

# Documentation
*.md
docs/

# Development files
docker-compose.yml
docker-compose.*.yml
Makefile

# Secrets
.env
.env.*
*.pem
*.key
```

### Пример 6: Makefile для Docker

```makefile
# Makefile
.PHONY: build run test docker-build docker-run docker-push

APP_NAME := myapp
VERSION := $(shell git describe --tags --always --dirty)
DOCKER_REPO := myregistry.com/myapp

# Local development
build:
	go build -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v -race -cover ./...

# Docker
docker-build:
	docker build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .

docker-run:
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

docker-push:
	docker tag $(APP_NAME):$(VERSION) $(DOCKER_REPO):$(VERSION)
	docker tag $(APP_NAME):latest $(DOCKER_REPO):latest
	docker push $(DOCKER_REPO):$(VERSION)
	docker push $(DOCKER_REPO):latest

# Docker Compose
up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f app

# Clean
clean:
	rm -rf bin/
	docker-compose down -v
```

### Пример 7: Health Check в приложении

```go
package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "net/http"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
}

type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Services  map[string]string `json:"services"`
}

func (h *HealthChecker) Check(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    status := HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Services:  make(map[string]string),
    }
    
    // Check database
    if err := h.db.PingContext(ctx); err != nil {
        status.Services["database"] = "unhealthy: " + err.Error()
        status.Status = "unhealthy"
    } else {
        status.Services["database"] = "healthy"
    }
    
    // Check Redis
    if err := h.redis.Ping(ctx).Err(); err != nil {
        status.Services["redis"] = "unhealthy: " + err.Error()
        status.Status = "unhealthy"
    } else {
        status.Services["redis"] = "healthy"
    }
    
    statusCode := http.StatusOK
    if status.Status != "healthy" {
        statusCode = http.StatusServiceUnavailable
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(status)
}
```

---

## 📝 Best Practices

### 1. Минимизируйте размер образа
- Используйте multi-stage builds
- Используйте alpine или scratch
- Удаляйте ненужные файлы

### 2. Кэшируйте зависимости
- Копируйте go.mod/go.sum отдельно
- Используйте `go mod download`

### 3. Безопасность
- Не запускайте от root
- Не храните секреты в образе
- Используйте .dockerignore

### 4. Health Checks
- Добавьте endpoint /health
- Настройте HEALTHCHECK в Dockerfile

---

## 🏋️ Практические задания

### Задание 1: Multi-stage Dockerfile

Создайте оптимизированный Dockerfile для вашего приложения.

**Начальный код:**
```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Run stage
FROM scratch
COPY --from=builder /app/main /main
EXPOSE 8080
ENTRYPOINT ["/main"]
```

**Баллы:** 15

### Задание 2: Docker Compose

Настройте полное окружение с PostgreSQL и Redis.

**Начальный код:**
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/app
      - REDIS_URL=redis:6379
    depends_on:
      - db
      - redis

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: app
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine

volumes:
  pgdata:
```

**Баллы:** 15

### Задание 3: Development Setup с Air

Настройте hot reload для разработки.

**Начальный код:**
```yaml
# docker-compose.dev.yml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
    ports:
      - "8080:8080"
```

```dockerfile
# Dockerfile.dev
FROM golang:1.22
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [Docker Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)
