# GitHub Actions для Go

---

## 💡 Ключевые идеи

1. **CI** — автоматическая проверка кода при каждом push/PR
2. **CD** — автоматический деплой
3. **Matrix builds** — тестирование на разных версиях Go
4. **Caching** — ускорение сборки
5. **Secrets** — безопасное хранение credentials

---

## 📖 Теория

### Что такое CI/CD?

**CI (Continuous Integration)** — автоматическая проверка кода при каждом изменении:
- Запуск тестов
- Линтинг
- Сборка
- Проверка безопасности

**CD (Continuous Deployment/Delivery)** — автоматическая доставка кода:
- Сборка Docker-образа
- Публикация в registry
- Деплой на сервер/Kubernetes

### Зачем нужен CI/CD?

Без CI/CD:
1. Разработчик пушит код
2. Забывает запустить тесты
3. Код ломает production
4. Все грустные

С CI/CD:
1. Разработчик пушит код
2. Автоматически запускаются тесты
3. Если тесты падают — merge блокируется
4. Если тесты проходят — автоматический деплой

### GitHub Actions

**GitHub Actions** — встроенный CI/CD в GitHub. Преимущества:
- Бесплатно для публичных репозиториев
- 2000 минут/месяц для приватных
- Интеграция с GitHub (PR, Issues, Releases)
- Большой marketplace готовых actions

### Основные концепции

**1. Workflow** — файл `.github/workflows/*.yml`:
```yaml
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: go test ./...
```

**2. Events (on)** — что запускает workflow:
```yaml
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  release:
    types: [published]
  schedule:
    - cron: '0 0 * * *'  # каждый день в полночь
```

**3. Jobs** — параллельные задачи:
```yaml
jobs:
  test:
    runs-on: ubuntu-latest
  lint:
    runs-on: ubuntu-latest
  deploy:
    needs: [test, lint]  # ждёт завершения test и lint
```

**4. Steps** — последовательные шаги внутри job:
```yaml
steps:
  - uses: actions/checkout@v4      # готовый action
  - run: go test ./...             # shell-команда
  - name: Custom step              # шаг с названием
    run: |
      echo "Multi-line"
      echo "commands"
```

### Matrix builds

Тестирование на разных версиях:
```yaml
strategy:
  matrix:
    go-version: ['1.21', '1.22']
    os: [ubuntu-latest, macos-latest]

steps:
  - uses: actions/setup-go@v5
    with:
      go-version: ${{ matrix.go-version }}
```

Создаст 4 job'а: 1.21+ubuntu, 1.21+macos, 1.22+ubuntu, 1.22+macos.

### Caching

Ускорение за счёт кэширования зависимостей:
```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.22'
    cache: true  # автоматическое кэширование go mod

# Или вручную:
- uses: actions/cache@v4
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

### Secrets

Для credentials используйте GitHub Secrets:
```yaml
steps:
  - name: Login to Docker Hub
    run: docker login -u ${{ secrets.DOCKER_USERNAME }} -p ${{ secrets.DOCKER_PASSWORD }}
```

Secrets добавляются в Settings → Secrets and variables → Actions.

### Типичный CI для Go

1. **Checkout** — получить код
2. **Setup Go** — установить Go
3. **Download deps** — скачать зависимости
4. **Lint** — проверить стиль (golangci-lint)
5. **Test** — запустить тесты
6. **Build** — собрать бинарник
7. **Upload artifacts** — сохранить результаты

---

## 💻 Примеры workflows

### Пример 1: Базовый CI

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true
      
      - name: Download dependencies
        run: go mod download
      
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: coverage.out

  lint:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  build:
    runs-on: ubuntu-latest
    needs: [test, lint]
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true
      
      - name: Build
        run: go build -v ./...
```

### Пример 2: CI с матрицей версий и базой данных

```yaml
# .github/workflows/ci-full.yml
name: CI Full

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    
    strategy:
      matrix:
        go-version: ['1.21', '1.22']
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: testdb
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      
      redis:
        image: redis:7
        ports:
          - 6379:6379
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go ${{ matrix.go-version }}
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
          cache: true
      
      - name: Run migrations
        run: |
          go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
          migrate -path ./migrations -database "postgres://test:test@localhost:5432/testdb?sslmode=disable" up
      
      - name: Run tests
        env:
          DATABASE_URL: postgres://test:test@localhost:5432/testdb?sslmode=disable
          REDIS_URL: redis://localhost:6379
        run: |
          go test -v -race -coverprofile=coverage.out ./...
      
      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | tr -d '%')
          echo "Coverage: $COVERAGE%"
          if (( $(echo "$COVERAGE < 70" | bc -l) )); then
            echo "Coverage is below 70%"
            exit 1
          fi
```

### Пример 3: CD с Docker

```yaml
# .github/workflows/cd.yml
name: CD

on:
  push:
    tags:
      - 'v*'

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Docker Buildx
        uses: docker/setup-buildx-action@v3
      
      - name: Login to Container Registry
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      
      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha
      
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
  
  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    
    steps:
      - name: Deploy to server
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          script: |
            cd /opt/myapp
            docker compose pull
            docker compose up -d
            docker system prune -f
```

### Пример 4: Release с GoReleaser

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true
      
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

```yaml
# .goreleaser.yml
project_name: myapp

before:
  hooks:
    - go mod tidy

builds:
  - main: ./cmd/server
    binary: myapp
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.Version={{.Version}}
      - -X main.Commit={{.Commit}}

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
```

### Пример 5: golangci-lint конфигурация

```yaml
# .golangci.yml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unconvert
    - bodyclose
    - noctx
    - gosec
    - prealloc

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  
  govet:
    check-shadowing: true
  
  gofmt:
    simplify: true
  
  gosec:
    excludes:
      - G104 # Audit errors not checked

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gosec
        - errcheck
```

---

## 🏋️ Практические задания

### Задание 1: CI Pipeline с тестами

Настройте CI с тестами и линтером.

**Начальный код:**
```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v4
      with:
        files: ./coverage.out

  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.22'
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v4
```

**Баллы:** 15

### Задание 2: Docker Build и Push

Создайте workflow для сборки Docker образа.

**Начальный код:**
```yaml
# .github/workflows/docker.yml
name: Docker

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Login to Docker Hub
      uses: docker/login-action@v3
      with:
        username: ${{ secrets.DOCKER_USERNAME }}
        password: ${{ secrets.DOCKER_PASSWORD }}
    
    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        push: true
        tags: |
          username/app:${{ github.ref_name }}
          username/app:latest
```

**Баллы:** 15

### Задание 3: Release с GoReleaser

Настройте автоматические релизы.

**Начальный код:**
```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags: ['v*']

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
    
    - uses: actions/setup-go@v5
      with:
        go-version: '1.22'
    
    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v5
      with:
        args: release --clean
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

```yaml
# .goreleaser.yml
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: 'checksums.txt'
```

**Баллы:** 10

---

## 🔗 Полезные ссылки

- [GitHub Actions](https://docs.github.com/en/actions)
- [GoReleaser](https://goreleaser.com/)
- [golangci-lint](https://golangci-lint.run/)
