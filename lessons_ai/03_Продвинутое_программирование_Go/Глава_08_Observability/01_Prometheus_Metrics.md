# Prometheus Metrics

---

## 💡 Ключевые идеи

1. **Prometheus** — система мониторинга и алертинга
2. **Метрики** — Counter, Gauge, Histogram, Summary
3. **Labels** — теги для фильтрации метрик
4. **Scraping** — Prometheus забирает метрики с endpoint'а
5. **Grafana** — визуализация метрик

### Типы метрик

| Тип | Описание | Пример |
|-----|----------|--------|
| Counter | Только увеличивается | Запросы, ошибки |
| Gauge | Может увеличиваться/уменьшаться | Температура, goroutines |
| Histogram | Распределение значений | Latency, размер ответа |
| Summary | Квантили (percentiles) | P50, P95, P99 |

---

## 📖 Теория

### Что такое Observability?

**Observability (наблюдаемость)** — способность понять состояние системы по её внешним выходам. Три столпа observability:

1. **Metrics** — числовые показатели (RPS, latency, errors)
2. **Logs** — текстовые записи событий
3. **Traces** — путь запроса через систему

В этом уроке мы фокусируемся на **метриках**.

### Зачем нужны метрики?

Без метрик вы узнаёте о проблемах, когда:
- Пользователи жалуются
- Сервис упал полностью
- Деньги уже потеряны

С метриками вы видите:
- Сервис замедляется → исправить до того, как упадёт
- Ошибки растут → найти причину до массовых жалоб
- Ресурсы заканчиваются → масштабировать заранее

### Что такое Prometheus?

**Prometheus** — это система мониторинга с открытым исходным кодом от Cloud Native Computing Foundation (CNCF). Особенности:

- **Pull-модель** — Prometheus сам забирает метрики с ваших сервисов
- **Time-series DB** — хранит метрики с временными метками
- **PromQL** — мощный язык запросов
- **Alerting** — отправка уведомлений при проблемах

### Как работает Prometheus?

```
┌─────────────┐     scrape      ┌─────────────┐
│  Ваш сервис │ ←────────────── │  Prometheus │
│  /metrics   │                  │             │
└─────────────┘                  └─────────────┘
                                       │
                                       │ query
                                       ↓
                                 ┌─────────────┐
                                 │   Grafana   │
                                 └─────────────┘
```

1. Ваш сервис отдаёт метрики на `/metrics`
2. Prometheus периодически (каждые 15-30 сек) забирает их
3. Grafana визуализирует данные из Prometheus

### Четыре типа метрик

**1. Counter (счётчик):**
```go
requestsTotal.Inc()  // +1
```
- Только увеличивается (или сбрасывается при рестарте)
- Примеры: запросы, ошибки, отправленные сообщения

**2. Gauge (датчик):**
```go
activeConnections.Set(42)
activeConnections.Inc()
activeConnections.Dec()
```
- Может увеличиваться и уменьшаться
- Примеры: текущие соединения, температура, память

**3. Histogram (гистограмма):**
```go
requestDuration.Observe(0.25)  // 250ms
```
- Распределение значений по бакетам
- Примеры: время ответа, размер запроса
- Позволяет вычислить перцентили (P50, P95, P99)

**4. Summary:**
- Похож на Histogram, но считает квантили на стороне клиента
- Обычно используют Histogram (гибче)

### Labels (метки)

Метки позволяют фильтровать и группировать метрики:

```go
requestsTotal.WithLabelValues("GET", "/api/users", "200").Inc()
requestsTotal.WithLabelValues("POST", "/api/users", "500").Inc()
```

Теперь можно строить графики:
- Все запросы: `sum(http_requests_total)`
- Только ошибки: `sum(http_requests_total{status=~"5.."})`
- По endpoint'ам: `sum(http_requests_total) by (path)`

### Важные метрики для веб-сервиса (RED)

**R**ate — количество запросов в секунду:
```promql
rate(http_requests_total[5m])
```

**E**rrors — процент ошибок:
```promql
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

**D**uration — время ответа (P99):
```promql
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
```

### Best Practices

1. **Не создавайте метрику для каждого пользователя** — это взорвёт кардинальность
2. **Используйте стандартные имена** — `http_requests_total`, не `myapp_requests`
3. **Добавляйте unit в имя** — `_seconds`, `_bytes`, `_total`
4. **Не используйте много labels** — каждая комбинация = отдельный time series

---

## 💻 Примеры кода

### Пример 1: Базовые метрики

```go
package main

import (
    "net/http"
    "time"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Определение метрик
var (
    // Counter — общее количество запросов
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    // Histogram — время ответа
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
        },
        []string{"method", "path"},
    )
    
    // Gauge — количество активных соединений
    activeConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )
    
    // Counter — бизнес-метрика
    ordersCreated = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "orders_created_total",
            Help: "Total number of orders created",
        },
    )
    
    // Gauge — размер очереди
    queueSize = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "queue_size",
            Help: "Current size of the queue",
        },
        []string{"queue_name"},
    )
)

func main() {
    // Endpoint для метрик
    http.Handle("/metrics", promhttp.Handler())
    
    // Бизнес endpoints
    http.HandleFunc("/api/orders", metricsMiddleware(ordersHandler))
    http.HandleFunc("/api/users", metricsMiddleware(usersHandler))
    
    http.ListenAndServe(":8080", nil)
}

// Middleware для сбора метрик
func metricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        activeConnections.Inc()
        defer activeConnections.Dec()
        
        // Wrap ResponseWriter для получения статуса
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        next(wrapped, r)
        
        duration := time.Since(start).Seconds()
        
        httpRequestsTotal.WithLabelValues(
            r.Method,
            r.URL.Path,
            http.StatusText(wrapped.statusCode),
        ).Inc()
        
        httpRequestDuration.WithLabelValues(
            r.Method,
            r.URL.Path,
        ).Observe(duration)
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
    // Бизнес-логика
    ordersCreated.Inc()
    w.Write([]byte("Order created"))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Users list"))
}
```

### Пример 2: Gin с Prometheus

```go
package main

import (
    "strconv"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

func init() {
    prometheus.MustRegister(httpRequests, httpDuration)
}

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.FullPath()
        if path == "" {
            path = "not_found"
        }
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(c.Writer.Status())
        
        httpRequests.WithLabelValues(c.Request.Method, path, status).Inc()
        httpDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
    }
}

func main() {
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(PrometheusMiddleware())
    
    // Metrics endpoint
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))
    
    // API endpoints
    r.GET("/api/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"users": []string{"Alice", "Bob"}})
    })
    
    r.Run(":8080")
}
```

### Пример 3: Custom Collector

```go
package main

import (
    "runtime"
    
    "github.com/prometheus/client_golang/prometheus"
)

// Custom collector для runtime метрик
type RuntimeCollector struct {
    goroutines *prometheus.Desc
    heapAlloc  *prometheus.Desc
    heapSys    *prometheus.Desc
    gcPause    *prometheus.Desc
}

func NewRuntimeCollector() *RuntimeCollector {
    return &RuntimeCollector{
        goroutines: prometheus.NewDesc(
            "go_goroutines_count",
            "Number of goroutines",
            nil, nil,
        ),
        heapAlloc: prometheus.NewDesc(
            "go_heap_alloc_bytes",
            "Heap allocation in bytes",
            nil, nil,
        ),
        heapSys: prometheus.NewDesc(
            "go_heap_sys_bytes",
            "Heap system bytes",
            nil, nil,
        ),
        gcPause: prometheus.NewDesc(
            "go_gc_pause_total_seconds",
            "Total GC pause time",
            nil, nil,
        ),
    }
}

func (c *RuntimeCollector) Describe(ch chan<- *prometheus.Desc) {
    ch <- c.goroutines
    ch <- c.heapAlloc
    ch <- c.heapSys
    ch <- c.gcPause
}

func (c *RuntimeCollector) Collect(ch chan<- prometheus.Metric) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    ch <- prometheus.MustNewConstMetric(
        c.goroutines,
        prometheus.GaugeValue,
        float64(runtime.NumGoroutine()),
    )
    
    ch <- prometheus.MustNewConstMetric(
        c.heapAlloc,
        prometheus.GaugeValue,
        float64(m.HeapAlloc),
    )
    
    ch <- prometheus.MustNewConstMetric(
        c.heapSys,
        prometheus.GaugeValue,
        float64(m.HeapSys),
    )
    
    ch <- prometheus.MustNewConstMetric(
        c.gcPause,
        prometheus.CounterValue,
        float64(m.PauseTotalNs)/1e9,
    )
}

func main() {
    prometheus.MustRegister(NewRuntimeCollector())
    // ...
}
```

### Пример 4: prometheus.yml

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

rule_files:
  - "alerts.yml"

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'myapp'
    static_configs:
      - targets: ['app:8080']
    metrics_path: '/metrics'

  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']
```

### Пример 5: Алерты

```yaml
# alerts.yml
groups:
  - name: myapp
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate"
          description: "Error rate is {{ $value | humanizePercentage }}"

      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency"
          description: "P95 latency is {{ $value }}s"

      - alert: TooManyGoroutines
        expr: go_goroutines_count > 10000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Too many goroutines"
```

### Пример 6: Docker Compose с мониторингом

```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - ./alerts.yml:/etc/prometheus/alerts.yml
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning

  alertmanager:
    image: prom/alertmanager:latest
    ports:
      - "9093:9093"
    volumes:
      - ./alertmanager.yml:/etc/alertmanager/alertmanager.yml

volumes:
  prometheus-data:
  grafana-data:
```

---

## 🏋️ Практические задания

### Задание 1: Базовые метрики

Добавьте Counter и Gauge метрики.

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

### Задание 2: Histogram для latency

Измеряйте время ответа с помощью Histogram.

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

### Задание 3: Бизнес-метрики

Добавьте метрики для бизнес-операций.

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

### Задание 4: Middleware для Gin

Создайте Prometheus middleware для Gin.

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

### Задание 5: Custom Collector

Создайте custom collector для системных метрик.

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

---

## 🔗 Полезные ссылки

- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
