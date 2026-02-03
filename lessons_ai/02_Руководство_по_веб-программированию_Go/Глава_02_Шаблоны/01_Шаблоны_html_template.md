# Шаблоны (html/template)

## 💡 Ключевые идеи

1. **html/template** — безопасные HTML шаблоны
2. **Контекст** — данные, передаваемые в шаблон (.)
3. **ParseFiles** — загрузка шаблонов из файлов
4. **Execute** — рендеринг шаблона с данными
5. **Автоэкранирование** — защита от XSS

---

## 📋 Синтаксис

### Создание и рендеринг

```go
// Из строки
tmpl := template.Must(template.New("name").Parse(text))

// Из файла
tmpl, err := template.ParseFiles("template.html")

// Рендеринг
tmpl.Execute(w, data)
```

### Синтаксис шаблонов

```go
{{ . }}            // контекст (все данные)
{{ .Field }}       // поле структуры
{{ .Method }}      // вызов метода
{{/* comment */}}  // комментарий
```

---

## 📖 Теория

### Что такое шаблоны?

**Шаблоны** — это HTML с плейсхолдерами для данных:

```html
<h1>Привет, {{.Name}}!</h1>
<p>Возраст: {{.Age}}</p>
```

При рендеринге `{{.Name}}` заменяется на реальное значение.

### html/template vs text/template

| html/template | text/template |
|---------------|---------------|
| Автоэкранирование HTML | Без экранирования |
| Для веб-страниц | Для email, конфигов |
| Защита от XSS | Нет защиты |

```go
// Безопасно для HTML
import "html/template"

// Для текстовых файлов
import "text/template"
```

### Точка (.) — контекст данных

```go
data := map[string]string{
    "Title": "Hello",
    "Body":  "World",
}

tmpl.Execute(w, data)
```

В шаблоне `.` — это переданные данные:
```html
<h1>{{.Title}}</h1>
<p>{{.Body}}</p>
```

### Создание шаблонов

**Из строки:**
```go
tmpl := template.Must(template.New("hello").Parse(`
    <h1>Hello, {{.Name}}!</h1>
`))
```

**Из файла:**
```go
tmpl, err := template.ParseFiles("templates/index.html")
```

**Несколько файлов:**
```go
tmpl, err := template.ParseGlob("templates/*.html")
```

### Рендеринг

```go
// В io.Writer (http.ResponseWriter)
tmpl.Execute(w, data)

// Конкретный шаблон по имени
tmpl.ExecuteTemplate(w, "index.html", data)
```

### Автоэкранирование

```go
data := map[string]string{
    "Content": "<script>alert('XSS')</script>",
}
```

`html/template` автоматически экранирует:
```html
&lt;script&gt;alert('XSS')&lt;/script&gt;
```

### Отключение экранирования (осторожно!)

```go
import "html/template"

data := map[string]interface{}{
    "HTML": template.HTML("<b>Bold</b>"),
}
```

```html
{{.HTML}}  <!-- <b>Bold</b> без экранирования -->
```

⚠️ Используйте только для доверенного контента!

---

## 💻 Примеры кода

### Простой шаблон

```go
package main

import (
    "html/template"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl := template.Must(template.New("hello").Parse(`
<!DOCTYPE html>
<html>
<body>
    <h1>Hello, {{ . }}!</h1>
</body>
</html>
        `))
        
        tmpl.Execute(w, "World")
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### Шаблон со структурой

```go
package main

import (
    "html/template"
    "net/http"
)

type PageData struct {
    Title   string
    Message string
    Year    int
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            Title:   "Welcome",
            Message: "Hello from Go templates!",
            Year:    2025,
        }
        
        tmpl := template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    <h1>{{ .Title }}</h1>
    <p>{{ .Message }}</p>
    <footer>&copy; {{ .Year }}</footer>
</body>
</html>
        `))
        
        tmpl.Execute(w, data)
    })
    
    http.ListenAndServe(":8080", nil)
}
```

### Загрузка из файла

```go
package main

import (
    "html/template"
    "net/http"
)

type User struct {
    Name  string
    Email string
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        user := User{Name: "John", Email: "john@example.com"}
        
        tmpl, err := template.ParseFiles("templates/user.html")
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        
        tmpl.Execute(w, user)
    })
    
    http.ListenAndServe(":8080", nil)
}
```

**templates/user.html:**
```html
<!DOCTYPE html>
<html>
<head>
    <title>User Profile</title>
</head>
<body>
    <h1>{{ .Name }}</h1>
    <p>Email: {{ .Email }}</p>
</body>
</html>
```

### Несколько шаблонов

```go
package main

import (
    "html/template"
    "net/http"
)

var templates *template.Template

func init() {
    templates = template.Must(template.ParseGlob("templates/*.html"))
}

type PageData struct {
    Title string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    data := PageData{Title: "Home"}
    templates.ExecuteTemplate(w, "home.html", data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    data := PageData{Title: "About"}
    templates.ExecuteTemplate(w, "about.html", data)
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/about", aboutHandler)
    
    http.ListenAndServe(":8080", nil)
}
```

### Базовый layout (наследование)

```go
package main

import (
    "html/template"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Порядок важен: сначала layout, потом content
        tmpl := template.Must(template.ParseFiles(
            "templates/layout.html",
            "templates/home.html",
        ))
        
        data := map[string]string{
            "Title": "Home Page",
        }
        
        tmpl.ExecuteTemplate(w, "layout", data)
    })
    
    http.ListenAndServe(":8080", nil)
}
```

**templates/layout.html:**
```html
{{ define "layout" }}
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    <nav>Navigation here</nav>
    <main>
        {{ template "content" . }}
    </main>
    <footer>Footer here</footer>
</body>
</html>
{{ end }}
```

**templates/home.html:**
```html
{{ define "content" }}
<h1>Welcome!</h1>
<p>This is the home page content.</p>
{{ end }}
```

### Передача функций в шаблон

```go
package main

import (
    "html/template"
    "net/http"
    "strings"
    "time"
)

func main() {
    funcMap := template.FuncMap{
        "upper":  strings.ToUpper,
        "lower":  strings.ToLower,
        "now":    time.Now,
        "format": func(t time.Time) string {
            return t.Format("02.01.2006 15:04")
        },
    }
    
    tmpl := template.Must(
        template.New("page").Funcs(funcMap).Parse(`
<!DOCTYPE html>
<html>
<body>
    <h1>{{ upper .Name }}</h1>
    <p>Current time: {{ now | format }}</p>
</body>
</html>
        `))
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := map[string]string{"Name": "hello world"}
        tmpl.Execute(w, data)
    })
    
    http.ListenAndServe(":8080", nil)
}
```

---

## ⚠️ Частые ошибки

### 1. Неэкспортируемые поля

```go
type User struct {
    name string  // ❌ Маленькая буква — недоступно в шаблоне
    Name string  // ✅ Большая буква — доступно
}
```

### 2. Не проверяют ошибку парсинга

```go
// ❌ Паника при ошибке в шаблоне
tmpl, _ := template.ParseFiles("broken.html")

// ✅ Используйте Must или проверяйте ошибку
tmpl := template.Must(template.ParseFiles("template.html"))
```

### 3. XSS при text/template

```go
// ❌ text/template не экранирует HTML
import "text/template"

// ✅ Используйте html/template для веба
import "html/template"
```

---

## 🏋️ Практические задания

### Задание 1: Unit тесты

Напишите unit тест для handler.

**Ожидаемый результат:**
```
Unit: httptest.NewRecorder(), httptest.NewRequest()
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Unit: httptest.NewRecorder(), httptest.NewRequest()")
}
```

**Ожидаемый вывод:**
```
Unit: httptest.NewRecorder(), httptest.NewRequest()
```

**Баллы:** 15

### Задание 2: Моки

Используйте моки для зависимостей.

**Ожидаемый результат:**
```
Моки: интерфейсы + testify/mock
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Моки: интерфейсы + testify/mock")
}
```

**Ожидаемый вывод:**
```
Моки: интерфейсы + testify/mock
```

**Баллы:** 15

### Задание 3: Integration тесты

Напишите интеграционный тест.

**Ожидаемый результат:**
```
Integration: реальная БД + httptest.Server
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Integration: реальная БД + httptest.Server")
}
```

**Ожидаемый вывод:**
```
Integration: реальная БД + httptest.Server
```

**Баллы:** 20

### Задание 4: TestMain

Настройте TestMain.

**Ожидаемый результат:**
```
TestMain: func TestMain(m *testing.M) { setup(); m.Run(); teardown() }
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("TestMain: func TestMain(m *testing.M) { setup(); m.Run(); teardown() }")
}
```

**Ожидаемый вывод:**
```
TestMain: func TestMain(m *testing.M) { setup(); m.Run(); teardown() }
```

**Баллы:** 15

### Задание 5: Test coverage

Измерьте покрытие тестами.

**Ожидаемый результат:**
```
Coverage: go test -cover ./...
```

**Начальный код:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Coverage: go test -cover ./...")
}
```

**Ожидаемый вывод:**
```
Coverage: go test -cover ./...
```

**Баллы:** 10
