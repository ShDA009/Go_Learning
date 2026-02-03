## Цель

Сделать прод‑подобный REST сервис на Gin с Postgres: регистрация/логин, CRUD доменных сущностей, транзакции, наблюдаемость, CI, контейнеризация, нагрузочное тестирование и профилирование.

## Домен: “сервис заказов”

Минимальные сущности:
- User
- Product
- Order
- OrderItem

## Функциональные требования (минимум)

### Аутентификация и авторизация
- Регистрация и логин
- JWT Access Token + Refresh Token
- Обновление токена (refresh)
- Logout (revocation refresh токена — через БД или Redis)
- Роли: user/admin (admin управляет продуктами)

### Продукты
- Создать/обновить/удалить продукт (admin)
- Список продуктов (пагинация + фильтр по имени)
- Получить продукт по ID

### Заказы
- Создать заказ (user)
- Получить заказ по ID (владелец или admin)
- Список заказов пользователя (пагинация)
- Транзакционность: создание заказа + списание остатков (если вы добавите склад/stock)

## Нефункциональные требования

- Конфигурация через переменные окружения (env) + разумные defaults
- Таймауты HTTP сервера и внешних вызовов
- Graceful shutdown (SIGTERM)
- Единый формат ошибок (HTTP status + код + сообщение)

## Наблюдаемость

### Метрики (Prometheus)
- HTTP: requests_total, duration_seconds (Histogram), errors_total
- Бизнес: orders_created_total, active_users (если есть)

### Логи (slog)
- JSON логи
- request_id обязателен
- trace_id/span_id (после подключения tracing)

### Трейсинг (OpenTelemetry)
- Сквозной trace для запроса: HTTP → DB → внешние вызовы (если есть)
- Экспорт в Jaeger/Tempo через OTLP

## Тестирование

- Unit тесты сервиса (use case слой)
- Интеграционные тесты с Postgres (поднимаем контейнер в CI)
- В CI запускать go test с -race

## CI/CD (GitHub Actions)

Минимум:
- golangci-lint
- go test -race -cover ./...
- build (и опционально docker build)

## Docker Compose (локальное окружение)

Поднимите как минимум:
- app
- postgres
- (опционально) migrate job
- prometheus + grafana
- otel-collector + jaeger

## Нагрузка и профилирование

Сделайте прогон hey/wrk/k6 и снимите:
- CPU профиль (pprof)
- Heap профиль (pprof)
- (опционально) trace (go tool trace)

## Definition of Done

- Есть README (как запустить, как прогнать тесты, как снять профили)
- API покрыто интеграционными тестами
- Метрики/логи/трейсы работают локально через docker-compose
- Есть отчёт “до/после” по одной оптимизации (latency/CPU/allocs)

## Рекомендуемые уроки (ссылки)

- Архитектура: [Clean Architecture в Go](/lessons/clean-architecture-v-go-1)
- Web: [Введение в Gin Framework](/lessons/vvedenie-v-gin-framework-1)
- Auth: [JWT Аутентификация](/lessons/jwt-autentifikatsiya-1)
- ORM: [Введение в GORM](/lessons/vvedenie-v-gorm-1)
- Docker: [Docker для Go приложений](/lessons/docker-dlya-go-prilozheniy-1)
- CI: [GitHub Actions для Go](/lessons/github-actions-dlya-go-1)
- Метрики: [Prometheus Metrics](/lessons/prometheus-metrics-1)
- Логи: [Структурное логирование: log/slog](/lessons/strukturnoe-logirovanie-logslog-1)
- Трейсы: [OpenTelemetry Tracing в Go](/lessons/opentelemetry-tracing-v-go-1)
- Трейсы (инструментация): [OpenTelemetry: Gin и gRPC](/lessons/opentelemetry-gin-i-grpc-2)
- Прод‑подготовка: [Production-ready HTTP сервер](/lessons/production-ready-http-server-1)
- Redis: [Redis: кэширование и TTL](/lessons/redis-keshirovanie-i-ttl-1)
- Миграции: [Миграции БД в Go](/lessons/migratsii-bd-v-go-1)
- SQL в проде: [database/sql в проде: пул, транзакции, EXPLAIN](/lessons/databasesql-v-prode-pul-tranzaktsii-explain-2)
- Тесты: [Integration Tests (Интеграционные тесты)](/lessons/integration-tests-integratsionnye-testy-7)
- Concurrency: [Гонка данных и её обнаружение (Data Race)](/lessons/gonka-dannyh-i-eyo-obnaruzhenie-data-race-9)
- Context: [Context Best Practices](/lessons/context-best-practices-4)
- Perf: [Профилирование с pprof](/lessons/profilirovanie-s-pprof-1)
- Perf: [Трассировка выполнения: go tool trace](/lessons/trassirovka-vypolneniya-go-tool-trace-2)
- Perf: [Нагрузочное тестирование и сбор профилей](/lessons/nagruzochnoe-testirovanie-i-sbor-profiley-3)

