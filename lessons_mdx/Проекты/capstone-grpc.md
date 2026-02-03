## Цель

Сделать прод‑подобный gRPC сервис с безопасностью (TLS/mTLS), interceptors, дедлайнами, наблюдаемостью (метрики/логи/трейсы). Опционально — grpc-gateway + OpenAPI для REST клиентов.

## Домен: Users/Accounts сервис

Минимальные RPC:
- CreateUser
- GetUser
- ListUsers (пагинация)
- UpdateUser
- DeleteUser

Опционально:
- Server streaming (например, ListUsersStream)

## Функциональные требования

### Контракт
- Определите proto с версиями (например, users.v1)
- Добавьте валидируемые поля (минимум на уровне приложения)

### Deadlines и отмена
- Клиент всегда выставляет deadline
- Сервер уважает ctx.Done и корректно прерывает работу
- (опционально) timeout interceptor задаёт дефолт, если клиент не выставил

### Interceptors
Минимум:
- logging (slog)
- recovery (panic → gRPC error)
- tracing/metrics (OpenTelemetry + Prometheus)
- auth (если добавите токен в metadata)

## Безопасность

- TLS обязателен
- mTLS для межсервисных вызовов (dev — self-signed CA)
- Никаких insecure режимов в CI/staging

## Наблюдаемость

### Логи (slog)
- JSON логи
- trace_id/span_id (после tracing)

### Трейсы (OpenTelemetry)
- server spans на каждый RPC
- client spans на исходящие gRPC/HTTP/DB вызовы (если есть)

### Метрики (Prometheus)
- счётчики вызовов по методу/статусу
- latency histogram по методу

## Тестирование

- Unit тесты бизнес‑логики
- gRPC тесты через bufconn/in-memory transport
- В CI: go test -race -cover ./...

## Docker Compose

Поднимите окружение для локальной проверки:
- сервис
- otel-collector + jaeger
- prometheus (+ grafana опционально)

## Definition of Done

- TLS/mTLS работает
- Есть interceptors и дедлайны
- Трейсы видны в Jaeger/Tempo, логи коррелируются по trace_id
- Есть CI, линт, тесты с -race

## Рекомендуемые уроки (ссылки)

- gRPC: [Введение в gRPC](/lessons/vvedenie-v-grpc-1)
- gRPC security: [gRPC TLS и mTLS](/lessons/grpc-tls-i-mtls-1)
- gRPC REST: [gRPC-Gateway и OpenAPI](/lessons/grpc-gateway-i-openapi-2)
- Логи: [Структурное логирование: log/slog](/lessons/strukturnoe-logirovanie-logslog-1)
- Трейсы: [OpenTelemetry Tracing в Go](/lessons/opentelemetry-tracing-v-go-1)
- Трейсы (инструментация): [OpenTelemetry: Gin и gRPC](/lessons/opentelemetry-gin-i-grpc-2)
- Метрики: [Prometheus Metrics](/lessons/prometheus-metrics-1)
- Docker: [Docker для Go приложений](/lessons/docker-dlya-go-prilozheniy-1)
- CI: [GitHub Actions для Go](/lessons/github-actions-dlya-go-1)
- Context: [Context Best Practices](/lessons/context-best-practices-4)

