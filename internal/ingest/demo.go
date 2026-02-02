package ingest

import (
	"context"
	"log"

	"golearning/internal/content"
)

// DemoData содержит демонстрационные уроки для тестирования.
type DemoData struct {
	repo *content.Repository
}

// NewDemoData создаёт новый генератор демо-данных.
func NewDemoData(repo *content.Repository) *DemoData {
	return &DemoData{repo: repo}
}

// Seed заполняет БД демонстрационными уроками.
func (d *DemoData) Seed(ctx context.Context) error {
	log.Println("Создание демонстрационных данных...")

	// Создаём модули
	modules := []content.Module{
		{Slug: "osnovy", Title: "Основы Go", OrderIndex: 0},
		{Slug: "tipy-dannyh", Title: "Типы данных", OrderIndex: 1},
		{Slug: "upravlenie", Title: "Управляющие конструкции", OrderIndex: 2},
	}

	for i := range modules {
		if err := d.repo.CreateModule(&modules[i]); err != nil {
			return err
		}
	}

	// Урок 1: Введение в Go
	lesson1 := createLesson1()
	lesson1.Lesson.ModuleID = modules[0].ID
	if err := d.saveLesson(lesson1); err != nil {
		return err
	}

	// Урок 2: Переменные
	lesson2 := createLesson2()
	lesson2.Lesson.ModuleID = modules[0].ID
	if err := d.saveLesson(lesson2); err != nil {
		return err
	}

	// Урок 3: Типы данных
	lesson3 := createLesson3()
	lesson3.Lesson.ModuleID = modules[1].ID
	if err := d.saveLesson(lesson3); err != nil {
		return err
	}

	// Урок 4: Операторы
	lesson4 := createLesson4()
	lesson4.Lesson.ModuleID = modules[1].ID
	if err := d.saveLesson(lesson4); err != nil {
		return err
	}

	// Урок 5: Условные конструкции
	lesson5 := createLesson5()
	lesson5.Lesson.ModuleID = modules[2].ID
	if err := d.saveLesson(lesson5); err != nil {
		return err
	}

	log.Println("Демонстрационные данные созданы!")
	return nil
}

type lessonData struct {
	Lesson   content.Lesson
	Sections []content.Section
	Tasks    []content.Task
}

func (d *DemoData) saveLesson(data lessonData) error {
	if err := d.repo.CreateLesson(&data.Lesson); err != nil {
		return err
	}
	log.Printf("  Урок: %s (ID=%d)", data.Lesson.Title, data.Lesson.ID)

	d.repo.DeleteSectionsByLessonID(data.Lesson.ID)
	d.repo.DeleteTasksByLessonID(data.Lesson.ID)

	for i := range data.Sections {
		data.Sections[i].LessonID = data.Lesson.ID
		if err := d.repo.CreateSection(&data.Sections[i]); err != nil {
			log.Printf("    Ошибка секции: %v", err)
		}
	}

	for i := range data.Tasks {
		data.Tasks[i].LessonID = data.Lesson.ID
		if err := d.repo.CreateTask(&data.Tasks[i]); err != nil {
			log.Printf("    Ошибка задания: %v", err)
		}
	}

	return nil
}

// ============================================================
// УРОК 1: Введение в Go
// ============================================================

func createLesson1() lessonData {
	return lessonData{
		Lesson: content.Lesson{
			Slug:           "vvedenie",
			Title:          "Введение в Go",
			OrderIndex:     1,
			SourceURL:      "https://metanit.com/go/tutorial/1.1.php",
			ReadingTimeMin: 12,
			BodyMD:         lesson1Body,
		},
		Sections: []content.Section{
			{Kind: content.SectionOverview, Title: "Ключевые идеи", OrderIndex: 0, BodyMD: lesson1Overview},
			{Kind: content.SectionSyntax, Title: "Синтаксис", OrderIndex: 1, BodyMD: lesson1Syntax},
			{Kind: content.SectionExamples, Title: "Примеры кода", OrderIndex: 2, BodyMD: lesson1Examples},
			{Kind: content.SectionPitfalls, Title: "Частые ошибки", OrderIndex: 3, BodyMD: lesson1Pitfalls},
		},
		Tasks: []content.Task{
			{
				Title:       "Hello, World!",
				OrderIndex:  0,
				Points:      10,
				PromptMD:    lesson1Task1Prompt,
				StarterCode: lesson1Task1Starter,
				TestsGo:     lesson1Task1Tests,
			},
			{
				Title:       "Приветствие пользователя",
				OrderIndex:  1,
				Points:      15,
				PromptMD:    lesson1Task2Prompt,
				StarterCode: lesson1Task2Starter,
				TestsGo:     lesson1Task2Tests,
			},
		},
	}
}

const lesson1Body = `# Введение в язык Go

Go (часто называемый Golang) — это компилируемый, статически типизированный язык программирования, разработанный в Google. Язык был создан в 2007 году, а первая стабильная версия вышла в 2012 году.

## Создатели языка

Go был создан тремя выдающимися инженерами:
- **Роберт Гризмер** — работал над Java HotSpot VM
- **Роб Пайк** — создатель Plan 9 и UTF-8  
- **Кен Томпсон** — создатель Unix и языка B

## Почему Go?

Go создавался для решения проблем, с которыми Google сталкивался при разработке масштабных систем:
- Медленная компиляция C++
- Сложность управления зависимостями
- Трудности с написанием параллельного кода
- Низкая читаемость кодовой базы`

const lesson1Overview = `## Что такое Go?

Go — это современный язык программирования, который сочетает в себе:

| Характеристика | Описание |
|----------------|----------|
| **Простота** | Минималистичный синтаксис, всего 25 ключевых слов |
| **Производительность** | Компилируется в нативный код, сравним с C/C++ |
| **Параллелизм** | Встроенная поддержка горутин и каналов |
| **Безопасность** | Сборка мусора, проверка границ массивов |
| **Быстрая компиляция** | Даже большие проекты компилируются за секунды |

### Где используется Go?

Go широко применяется в:

1. **Бэкенд-разработке** — веб-серверы, REST API, микросервисы
2. **Облачной инфраструктуре** — Docker, Kubernetes, Terraform написаны на Go
3. **Системном программировании** — утилиты командной строки, сетевые инструменты
4. **Распределённых системах** — базы данных (CockroachDB, InfluxDB)

### Философия Go

> «Меньше — это больше» (Less is more)

Go намеренно избегает сложных концепций:
- Нет классического ООП с наследованием
- Нет исключений (только явная обработка ошибок)
- Нет дженериков (добавлены только в Go 1.18)
- Нет перегрузки функций`

const lesson1Syntax = `## Структура программы на Go

Каждая программа на Go состоит из **пакетов** (packages). Исполняемая программа всегда начинается с пакета ` + "`main`" + ` и функции ` + "`main()`" + `.

### Минимальная программа

` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

### Разбор структуры

| Элемент | Описание |
|---------|----------|
| ` + "`package main`" + ` | Объявление пакета. ` + "`main`" + ` — специальный пакет для исполняемых программ |
| ` + "`import \"fmt\"`" + ` | Импорт пакета ` + "`fmt`" + ` для форматированного ввода-вывода |
| ` + "`func main()`" + ` | Главная функция — точка входа в программу |
| ` + "`fmt.Println()`" + ` | Вывод строки с переводом строки в конце |

### Правила именования

- **Пакеты**: нижний регистр, одно слово (` + "`fmt`" + `, ` + "`http`" + `, ` + "`json`" + `)
- **Экспортируемые имена**: начинаются с заглавной буквы (` + "`Println`" + `, ` + "`Reader`" + `)
- **Приватные имена**: начинаются со строчной буквы (` + "`internalFunc`" + `)

### Множественный импорт

` + "```go" + `
import (
    "fmt"
    "math"
    "strings"
)
` + "```" + `

### Комментарии

` + "```go" + `
// Однострочный комментарий

/*
   Многострочный
   комментарий
*/
` + "```" + ``

const lesson1Examples = `## Примеры программ

### Пример 1: Вывод текста

` + "```go" + `
package main

import "fmt"

func main() {
    // Println добавляет перевод строки
    fmt.Println("Первая строка")
    fmt.Println("Вторая строка")
    
    // Print не добавляет перевод строки
    fmt.Print("Это ")
    fmt.Print("одна строка\n")
}
` + "```" + `

**Вывод:**
` + "```" + `
Первая строка
Вторая строка
Это одна строка
` + "```" + `

### Пример 2: Форматированный вывод

` + "```go" + `
package main

import "fmt"

func main() {
    name := "Go"
    version := 1.22
    year := 2024
    
    // Printf — форматированный вывод
    fmt.Printf("Язык: %s\n", name)
    fmt.Printf("Версия: %.2f\n", version)
    fmt.Printf("Год: %d\n", year)
    
    // Sprintf возвращает строку вместо вывода
    message := fmt.Sprintf("%s версии %.1f", name, version)
    fmt.Println(message)
}
` + "```" + `

**Вывод:**
` + "```" + `
Язык: Go
Версия: 1.22
Год: 2024
Go версии 1.2
` + "```" + `

### Пример 3: Использование нескольких пакетов

` + "```go" + `
package main

import (
    "fmt"
    "math"
    "strings"
)

func main() {
    // Пакет math — математические функции
    fmt.Printf("Квадратный корень из 16: %.0f\n", math.Sqrt(16))
    fmt.Printf("Число Пи: %.4f\n", math.Pi)
    
    // Пакет strings — работа со строками
    text := "Hello, Go!"
    fmt.Printf("Верхний регистр: %s\n", strings.ToUpper(text))
    fmt.Printf("Содержит 'Go': %t\n", strings.Contains(text, "Go"))
}
` + "```" + `

**Вывод:**
` + "```" + `
Квадратный корень из 16: 4
Число Пи: 3.1416
Верхний регистр: HELLO, GO!
Содержит 'Go': true
` + "```" + ``

const lesson1Pitfalls = `## Частые ошибки начинающих

### 1. Неиспользуемые импорты

Go **не позволяет** импортировать пакеты, которые не используются:

` + "```go" + `
package main

import (
    "fmt"
    "math"  // ОШИБКА: imported and not used: "math"
)

func main() {
    fmt.Println("Hello")
}
` + "```" + `

**Решение:** удалите неиспользуемый импорт или используйте ` + "`_`" + ` для подавления:

` + "```go" + `
import (
    "fmt"
    _ "math"  // Импорт только для побочных эффектов
)
` + "```" + `

### 2. Неиспользуемые переменные

Объявленные переменные **обязательно** должны использоваться:

` + "```go" + `
func main() {
    x := 10  // ОШИБКА: x declared but not used
    fmt.Println("Hello")
}
` + "```" + `

### 3. Открывающая скобка на новой строке

В Go **обязательно** размещать ` + "`{`" + ` на той же строке:

` + "```go" + `
// ОШИБКА: syntax error
func main()
{
    fmt.Println("Hello")
}

// ПРАВИЛЬНО
func main() {
    fmt.Println("Hello")
}
` + "```" + `

### 4. Забыли package main

` + "```go" + `
// ОШИБКА: программа не запустится без package main
import "fmt"

func main() {
    fmt.Println("Hello")
}
` + "```" + `

### 5. Неправильный регистр при импорте

` + "```go" + `
import "Fmt"  // ОШИБКА: пакеты пишутся в нижнем регистре

import "fmt"  // ПРАВИЛЬНО
` + "```" + `

### Советы

- Используйте ` + "`go fmt`" + ` для автоматического форматирования кода
- Настройте IDE на автоудаление неиспользуемых импортов
- Читайте сообщения об ошибках — они очень информативны в Go`

const lesson1Task1Prompt = `### Задание: Hello, World!

Напишите программу, которая выводит на экран приветствие **"Hello, Go!"**.

**Требования:**
- Используйте пакет ` + "`fmt`" + `
- Вывод должен быть **ровно** ` + "`Hello, Go!`" + ` (с восклицательным знаком)
- Используйте функцию ` + "`fmt.Println`" + `

**Ожидаемый вывод:**
` + "```" + `
Hello, Go!
` + "```" + ``

const lesson1Task1Starter = `package main

import "fmt"

func main() {
	// Напишите код для вывода "Hello, Go!"
	
}
`

const lesson1Task1Tests = `package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Программа завершилась с ошибкой: %v", err)
	}
	
	output := strings.TrimSpace(out.String())
	expected := "Hello, Go!"
	
	if output != expected {
		t.Errorf("Ожидалось %q, получено %q", expected, output)
	}
}
`

const lesson1Task2Prompt = `### Задание: Приветствие

Напишите программу, которая выводит три строки:
1. Ваше имя (например, "Gopher")
2. Год изучения Go (например, "2024")
3. Фразу "Учу Go!"

**Ожидаемый формат вывода:**
` + "```" + `
Gopher
2024
Учу Go!
` + "```" + `

**Подсказка:** используйте три вызова ` + "`fmt.Println()`" + ``

const lesson1Task2Starter = `package main

import "fmt"

func main() {
	// Выведите своё имя
	
	// Выведите год
	
	// Выведите "Учу Go!"
	
}
`

const lesson1Task2Tests = `package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestGreeting(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Ошибка выполнения: %v", err)
	}
	
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	
	if len(lines) < 3 {
		t.Errorf("Ожидалось минимум 3 строки, получено %d", len(lines))
		return
	}
	
	if !strings.Contains(out.String(), "Учу Go!") {
		t.Errorf("Вывод должен содержать 'Учу Go!'")
	}
}
`

// ============================================================
// УРОК 2: Переменные
// ============================================================

func createLesson2() lessonData {
	return lessonData{
		Lesson: content.Lesson{
			Slug:           "peremennye",
			Title:          "Переменные и константы",
			OrderIndex:     2,
			SourceURL:      "https://metanit.com/go/tutorial/2.2.php",
			ReadingTimeMin: 15,
			BodyMD:         lesson2Body,
		},
		Sections: []content.Section{
			{Kind: content.SectionOverview, Title: "Ключевые идеи", OrderIndex: 0, BodyMD: lesson2Overview},
			{Kind: content.SectionSyntax, Title: "Синтаксис", OrderIndex: 1, BodyMD: lesson2Syntax},
			{Kind: content.SectionExamples, Title: "Примеры кода", OrderIndex: 2, BodyMD: lesson2Examples},
			{Kind: content.SectionPitfalls, Title: "Частые ошибки", OrderIndex: 3, BodyMD: lesson2Pitfalls},
		},
		Tasks: []content.Task{
			{
				Title:       "Объявление переменных",
				OrderIndex:  0,
				Points:      10,
				PromptMD:    lesson2Task1Prompt,
				StarterCode: lesson2Task1Starter,
				TestsGo:     lesson2Task1Tests,
			},
			{
				Title:       "Обмен значениями",
				OrderIndex:  1,
				Points:      15,
				PromptMD:    lesson2Task2Prompt,
				StarterCode: lesson2Task2Starter,
				TestsGo:     lesson2Task2Tests,
			},
		},
	}
}

const lesson2Body = `# Переменные и константы в Go

Переменные — это именованные области памяти для хранения данных. В Go переменные **строго типизированы**: каждая переменная имеет определённый тип, который нельзя изменить.`

const lesson2Overview = `## Что такое переменная?

**Переменная** — это именованный контейнер для хранения значения определённого типа.

### Характеристики переменных в Go

| Свойство | Описание |
|----------|----------|
| **Статическая типизация** | Тип переменной определяется при объявлении и не может измениться |
| **Нулевые значения** | Неинициализированные переменные получают значение по умолчанию |
| **Область видимости** | Переменные видны только в блоке, где объявлены |
| **Обязательное использование** | Объявленная переменная должна быть использована |

### Нулевые значения (zero values)

Если переменная объявлена без инициализации, она получает **нулевое значение**:

| Тип | Нулевое значение |
|-----|------------------|
| ` + "`int`" + `, ` + "`float64`" + ` | ` + "`0`" + ` |
| ` + "`string`" + ` | ` + "`\"\"`" + ` (пустая строка) |
| ` + "`bool`" + ` | ` + "`false`" + ` |
| указатели, slice, map, channel, function | ` + "`nil`" + ` |

### Константы

**Константы** — неизменяемые значения, известные на этапе компиляции:

` + "```go" + `
const Pi = 3.14159
const AppName = "GoLearning"
const MaxRetries = 3
` + "```" + ``

const lesson2Syntax = `## Способы объявления переменных

### 1. Полное объявление с var

` + "```go" + `
var name string = "Alice"
var age int = 25
var isActive bool = true
var price float64 = 99.99
` + "```" + `

### 2. Объявление без начального значения

` + "```go" + `
var count int      // count = 0
var message string // message = ""
var flag bool      // flag = false
` + "```" + `

### 3. Вывод типа (type inference)

Компилятор сам определяет тип по значению:

` + "```go" + `
var name = "Alice"   // string
var age = 25         // int
var price = 99.99    // float64
var active = true    // bool
` + "```" + `

### 4. Краткое объявление := (short declaration)

Самый распространённый способ внутри функций:

` + "```go" + `
name := "Alice"     // string
age := 25           // int
price := 99.99      // float64
active := true      // bool
` + "```" + `

> ⚠️ **Важно:** ` + "`:=`" + ` работает только внутри функций!

### 5. Множественное объявление

` + "```go" + `
// Несколько переменных одного типа
var x, y, z int = 1, 2, 3

// Несколько переменных разных типов
var (
    name    string  = "Alice"
    age     int     = 25
    balance float64 = 1000.50
)

// Краткое множественное объявление
a, b, c := 1, "hello", true
` + "```" + `

### 6. Константы

` + "```go" + `
// Одиночные константы
const Pi = 3.14159
const MaxSize = 100

// Блок констант
const (
    StatusOK    = 200
    StatusError = 500
    AppVersion  = "1.0.0"
)

// Константы с iota (автоинкремент)
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
    Wednesday      // 3
)
` + "```" + ``

const lesson2Examples = `## Практические примеры

### Пример 1: Работа с переменными

` + "```go" + `
package main

import "fmt"

func main() {
    // Объявление и инициализация
    name := "Иван"
    age := 30
    salary := 75000.50
    
    // Вывод значений
    fmt.Println("Имя:", name)
    fmt.Println("Возраст:", age)
    fmt.Printf("Зарплата: %.2f руб.\n", salary)
    
    // Изменение значений
    age = 31
    salary = salary * 1.1 // повышение на 10%
    
    fmt.Println("\nПосле изменений:")
    fmt.Println("Возраст:", age)
    fmt.Printf("Зарплата: %.2f руб.\n", salary)
}
` + "```" + `

**Вывод:**
` + "```" + `
Имя: Иван
Возраст: 30
Зарплата: 75000.50 руб.

После изменений:
Возраст: 31
Зарплата: 82500.55 руб.
` + "```" + `

### Пример 2: Множественное присваивание

` + "```go" + `
package main

import "fmt"

func main() {
    // Одновременное присваивание
    a, b := 10, 20
    fmt.Printf("a = %d, b = %d\n", a, b)
    
    // Обмен значениями (Go позволяет сделать это элегантно!)
    a, b = b, a
    fmt.Printf("После обмена: a = %d, b = %d\n", a, b)
    
    // Множественное присваивание с функцией
    x, y, z := getCoordinates()
    fmt.Printf("Координаты: (%d, %d, %d)\n", x, y, z)
}

func getCoordinates() (int, int, int) {
    return 10, 20, 30
}
` + "```" + `

**Вывод:**
` + "```" + `
a = 10, b = 20
После обмена: a = 20, b = 10
Координаты: (10, 20, 30)
` + "```" + `

### Пример 3: Константы и iota

` + "```go" + `
package main

import "fmt"

const (
    KB = 1 << (10 * iota) // 1 << 0 = 1
    MB                     // 1 << 10 = 1024
    GB                     // 1 << 20 = 1048576
    TB                     // 1 << 30 = 1073741824
)

func main() {
    fileSize := 2.5 * GB
    
    fmt.Printf("1 KB = %d байт\n", KB)
    fmt.Printf("1 MB = %d байт\n", MB)
    fmt.Printf("1 GB = %d байт\n", GB)
    fmt.Printf("\nРазмер файла: %.0f байт\n", fileSize)
    fmt.Printf("Или: %.2f GB\n", fileSize/float64(GB))
}
` + "```" + `

### Пример 4: Область видимости

` + "```go" + `
package main

import "fmt"

var globalVar = "Я глобальная" // Видна во всём пакете

func main() {
    localVar := "Я локальная" // Видна только в main()
    
    fmt.Println(globalVar)
    fmt.Println(localVar)
    
    if true {
        blockVar := "Я в блоке if" // Видна только в этом if
        fmt.Println(blockVar)
    }
    // fmt.Println(blockVar) // ОШИБКА: blockVar не видна здесь
}
` + "```" + ``

const lesson2Pitfalls = `## Частые ошибки

### 1. Повторное объявление с :=

` + "```go" + `
name := "Alice"
name := "Bob"  // ОШИБКА: no new variables on left side

// Правильно: используйте = для присваивания
name = "Bob"
` + "```" + `

### 2. := вне функции

` + "```go" + `
package main

name := "Alice"  // ОШИБКА: non-declaration statement outside function body

// Правильно: используйте var на уровне пакета
var name = "Alice"
` + "```" + `

### 3. Неиспользуемая переменная

` + "```go" + `
func main() {
    x := 10  // ОШИБКА: x declared but not used
    fmt.Println("Hello")
}

// Решение 1: использовать переменную
func main() {
    x := 10
    fmt.Println(x)
}

// Решение 2: использовать _ для игнорирования
func main() {
    x, _ := someFunction()
    fmt.Println(x)
}
` + "```" + `

### 4. Затенение переменных (shadowing)

` + "```go" + `
package main

import "fmt"

var x = 10 // глобальная

func main() {
    fmt.Println(x) // 10 (глобальная)
    
    x := 20 // ВНИМАНИЕ: создаётся НОВАЯ локальная переменная!
    fmt.Println(x) // 20 (локальная)
    
    if true {
        x := 30 // Ещё одна НОВАЯ переменная в блоке if
        fmt.Println(x) // 30
    }
    
    fmt.Println(x) // 20 (локальная из main)
}
` + "```" + `

### 5. Изменение константы

` + "```go" + `
const Pi = 3.14
Pi = 3.14159  // ОШИБКА: cannot assign to Pi

// Константы нельзя изменять!
` + "```" + `

### Советы

- Предпочитайте ` + "`:=`" + ` внутри функций для краткости
- Используйте ` + "`var`" + ` когда нужно явно указать тип
- Группируйте связанные переменные в блоки ` + "`var (...)`" + `
- Давайте переменным понятные имена`

const lesson2Task1Prompt = `### Задание: Объявление переменных

Объявите три переменные:
- ` + "`name`" + ` типа ` + "`string`" + ` со значением ` + "`\"Gopher\"`" + `
- ` + "`age`" + ` типа ` + "`int`" + ` со значением ` + "`15`" + `
- ` + "`pi`" + ` типа ` + "`float64`" + ` со значением ` + "`3.14`" + `

Выведите их через пробел.

**Ожидаемый вывод:**
` + "```" + `
Gopher 15 3.14
` + "```" + `

**Подсказка:** используйте ` + "`fmt.Println(name, age, pi)`" + ``

const lesson2Task1Starter = `package main

import "fmt"

func main() {
	// Объявите переменные name, age, pi
	
	
	// Выведите их через пробел
	fmt.Println(name, age, pi)
}
`

const lesson2Task1Tests = `package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestVariables(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Ошибка: %v", err)
	}
	
	output := strings.TrimSpace(out.String())
	if output != "Gopher 15 3.14" {
		t.Errorf("Ожидалось 'Gopher 15 3.14', получено '%s'", output)
	}
}
`

const lesson2Task2Prompt = `### Задание: Обмен значениями

Даны две переменные ` + "`a = 5`" + ` и ` + "`b = 10`" + `.

Поменяйте их значения местами **без использования третьей переменной** (Go позволяет сделать это в одну строку!).

**Ожидаемый вывод:**
` + "```" + `
a = 10, b = 5
` + "```" + `

**Подсказка:** в Go можно писать ` + "`a, b = b, a`" + ``

const lesson2Task2Starter = `package main

import "fmt"

func main() {
	a := 5
	b := 10
	
	// Поменяйте значения a и b местами
	
	
	fmt.Printf("a = %d, b = %d\n", a, b)
}
`

const lesson2Task2Tests = `package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestSwap(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Ошибка: %v", err)
	}
	
	output := strings.TrimSpace(out.String())
	if output != "a = 10, b = 5" {
		t.Errorf("Ожидалось 'a = 10, b = 5', получено '%s'", output)
	}
}
`

// ============================================================
// УРОК 3: Типы данных
// ============================================================

func createLesson3() lessonData {
	return lessonData{
		Lesson: content.Lesson{
			Slug:           "tipy-dannyh",
			Title:          "Базовые типы данных",
			OrderIndex:     3,
			SourceURL:      "https://metanit.com/go/tutorial/2.3.php",
			ReadingTimeMin: 18,
			BodyMD:         lesson3Body,
		},
		Sections: []content.Section{
			{Kind: content.SectionOverview, Title: "Ключевые идеи", OrderIndex: 0, BodyMD: lesson3Overview},
			{Kind: content.SectionSyntax, Title: "Синтаксис", OrderIndex: 1, BodyMD: lesson3Syntax},
			{Kind: content.SectionExamples, Title: "Примеры кода", OrderIndex: 2, BodyMD: lesson3Examples},
			{Kind: content.SectionPitfalls, Title: "Частые ошибки", OrderIndex: 3, BodyMD: lesson3Pitfalls},
		},
		Tasks: []content.Task{
			{
				Title:       "Работа с числами",
				OrderIndex:  0,
				Points:      15,
				PromptMD:    lesson3Task1Prompt,
				StarterCode: lesson3Task1Starter,
				TestsGo:     lesson3Task1Tests,
			},
		},
	}
}

const lesson3Body = `# Базовые типы данных в Go

Go — **статически типизированный** язык. Это означает, что тип каждой переменной известен на этапе компиляции и не может измениться во время выполнения программы.`

const lesson3Overview = `## Категории типов данных

### Числовые типы

| Тип | Размер | Диапазон значений |
|-----|--------|-------------------|
| ` + "`int8`" + ` | 8 бит | -128 до 127 |
| ` + "`int16`" + ` | 16 бит | -32768 до 32767 |
| ` + "`int32`" + ` | 32 бита | -2147483648 до 2147483647 |
| ` + "`int64`" + ` | 64 бита | -9223372036854775808 до 9223372036854775807 |
| ` + "`int`" + ` | 32/64 бита | Зависит от платформы |

| Тип | Размер | Диапазон значений |
|-----|--------|-------------------|
| ` + "`uint8`" + ` (` + "`byte`" + `) | 8 бит | 0 до 255 |
| ` + "`uint16`" + ` | 16 бит | 0 до 65535 |
| ` + "`uint32`" + ` | 32 бита | 0 до 4294967295 |
| ` + "`uint64`" + ` | 64 бита | 0 до 18446744073709551615 |
| ` + "`uint`" + ` | 32/64 бита | Зависит от платформы |

### Числа с плавающей точкой

| Тип | Размер | Точность |
|-----|--------|----------|
| ` + "`float32`" + ` | 32 бита | ~7 значащих цифр |
| ` + "`float64`" + ` | 64 бита | ~15 значащих цифр |

### Другие базовые типы

| Тип | Описание |
|-----|----------|
| ` + "`bool`" + ` | Логический тип: ` + "`true`" + ` или ` + "`false`" + ` |
| ` + "`string`" + ` | Строка (последовательность байтов UTF-8) |
| ` + "`byte`" + ` | Псевдоним для ` + "`uint8`" + ` |
| ` + "`rune`" + ` | Псевдоним для ` + "`int32`" + ` (символ Unicode) |
| ` + "`complex64`" + `, ` + "`complex128`" + ` | Комплексные числа |`

const lesson3Syntax = `## Объявление типов данных

### Целые числа

` + "```go" + `
var a int = 42           // int (размер зависит от платформы)
var b int8 = 127         // от -128 до 127
var c int16 = 32000      // от -32768 до 32767
var d int32 = 2000000000 // от -2147483648 до 2147483647
var e int64 = 9223372036854775807 // очень большие числа

var f uint = 42          // беззнаковый int
var g uint8 = 255        // от 0 до 255 (byte)
var h uint64 = 18446744073709551615 // максимальное uint64
` + "```" + `

### Числа с плавающей точкой

` + "```go" + `
var pi float64 = 3.14159265358979  // рекомендуется по умолчанию
var e float32 = 2.71828            // меньше памяти, меньше точность

// Научная нотация
var avogadro = 6.022e23  // 6.022 × 10²³
var planck = 6.626e-34   // 6.626 × 10⁻³⁴
` + "```" + `

### Строки

` + "```go" + `
var greeting string = "Привет, мир!"
var multiline = ` + "`" + `Это
многострочная
строка` + "`" + `

// Строки неизменяемы!
s := "hello"
// s[0] = 'H'  // ОШИБКА: cannot assign to s[0]
s = "Hello"   // Можно только заменить всю строку
` + "```" + `

### Логический тип

` + "```go" + `
var isActive bool = true
var isAdmin bool = false
var isValid = (10 > 5)   // true
` + "```" + `

### Преобразование типов

Go **не выполняет** неявное преобразование типов:

` + "```go" + `
var x int = 10
var y float64 = float64(x)  // явное преобразование
var z int = int(y)          // явное преобразование

var a int32 = 100
var b int64 = int64(a)      // даже между int32 и int64

// Строки и числа
num := 42
str := strconv.Itoa(num)    // "42"
num2, _ := strconv.Atoi(str) // 42
` + "```" + ``

const lesson3Examples = `## Практические примеры

### Пример 1: Работа с числами

` + "```go" + `
package main

import (
    "fmt"
    "math"
)

func main() {
    // Целые числа
    var a int = 10
    var b int = 3
    
    fmt.Println("Сложение:", a + b)      // 13
    fmt.Println("Вычитание:", a - b)     // 7
    fmt.Println("Умножение:", a * b)     // 30
    fmt.Println("Деление:", a / b)       // 3 (целочисленное!)
    fmt.Println("Остаток:", a % b)       // 1
    
    // Числа с плавающей точкой
    var x float64 = 10.0
    var y float64 = 3.0
    
    fmt.Printf("Точное деление: %.2f\n", x / y)  // 3.33
    fmt.Printf("Степень: %.2f\n", math.Pow(x, 2)) // 100.00
    fmt.Printf("Корень: %.2f\n", math.Sqrt(x))    // 3.16
}
` + "```" + `

### Пример 2: Работа со строками

` + "```go" + `
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello, Go!"
    
    // Длина строки (в байтах)
    fmt.Println("Длина:", len(s))  // 10
    
    // Доступ к символу (возвращает байт)
    fmt.Printf("Первый байт: %c\n", s[0])  // H
    
    // Срез строки
    fmt.Println("Срез [0:5]:", s[0:5])  // Hello
    
    // Конкатенация
    s2 := s + " Welcome!"
    fmt.Println(s2)  // Hello, Go! Welcome!
    
    // Методы пакета strings
    fmt.Println("Верхний регистр:", strings.ToUpper(s))
    fmt.Println("Замена:", strings.Replace(s, "Go", "World", 1))
    fmt.Println("Разделение:", strings.Split(s, ", "))
}
` + "```" + `

### Пример 3: Работа с Unicode (rune)

` + "```go" + `
package main

import "fmt"

func main() {
    // Строка с русскими символами
    s := "Привет"
    
    // len() возвращает количество БАЙТОВ
    fmt.Println("Байтов:", len(s))  // 12 (каждая буква = 2 байта)
    
    // Для подсчёта символов используем []rune
    runes := []rune(s)
    fmt.Println("Символов:", len(runes))  // 6
    
    // Итерация по символам
    for i, r := range s {
        fmt.Printf("Позиция %d: %c (код %d)\n", i, r, r)
    }
}
` + "```" + ``

const lesson3Pitfalls = `## Частые ошибки

### 1. Целочисленное деление

` + "```go" + `
a := 7
b := 2
result := a / b  // 3, а не 3.5!

// Для точного деления преобразуйте в float
result := float64(a) / float64(b)  // 3.5
` + "```" + `

### 2. Переполнение типа

` + "```go" + `
var x int8 = 127
x = x + 1  // Переполнение! x станет -128

var y uint8 = 255
y = y + 1  // Переполнение! y станет 0
` + "```" + `

### 3. Неявное преобразование типов

` + "```go" + `
var a int32 = 100
var b int64 = a  // ОШИБКА: cannot use a (type int32) as type int64

// Правильно:
var b int64 = int64(a)
` + "```" + `

### 4. Сравнение float

` + "```go" + `
a := 0.1 + 0.2
b := 0.3

if a == b {  // Может быть false из-за погрешности!
    fmt.Println("Равны")
}

// Правильно: сравнивайте с погрешностью
const epsilon = 1e-9
if math.Abs(a - b) < epsilon {
    fmt.Println("Примерно равны")
}
` + "```" + `

### 5. Изменение символа в строке

` + "```go" + `
s := "hello"
s[0] = 'H'  // ОШИБКА: строки неизменяемы

// Правильно: создайте новую строку
s = "H" + s[1:]  // "Hello"

// Или через []byte
b := []byte(s)
b[0] = 'H'
s = string(b)
` + "```" + ``

const lesson3Task1Prompt = `### Задание: Калькулятор

Вычислите и выведите:
- Сумму чисел 17 и 25
- Разность чисел 100 и 37
- Произведение чисел 6 и 7

**Ожидаемый вывод:**
` + "```" + `
Сумма: 42
Разность: 63
Произведение: 42
` + "```" + ``

const lesson3Task1Starter = `package main

import "fmt"

func main() {
	// Вычислите и выведите результаты
	
}
`

const lesson3Task1Tests = `package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestCalc(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Ошибка: %v", err)
	}
	
	output := out.String()
	if !strings.Contains(output, "42") || !strings.Contains(output, "63") {
		t.Errorf("Неверный результат: %s", output)
	}
}
`

// ============================================================
// УРОК 4: Операторы
// ============================================================

func createLesson4() lessonData {
	return lessonData{
		Lesson: content.Lesson{
			Slug:           "operatory",
			Title:          "Арифметические и логические операторы",
			OrderIndex:     4,
			SourceURL:      "https://metanit.com/go/tutorial/2.5.php",
			ReadingTimeMin: 12,
			BodyMD:         lesson4Body,
		},
		Sections: []content.Section{
			{Kind: content.SectionOverview, Title: "Ключевые идеи", OrderIndex: 0, BodyMD: lesson4Overview},
			{Kind: content.SectionSyntax, Title: "Синтаксис", OrderIndex: 1, BodyMD: lesson4Syntax},
			{Kind: content.SectionExamples, Title: "Примеры кода", OrderIndex: 2, BodyMD: lesson4Examples},
			{Kind: content.SectionPitfalls, Title: "Частые ошибки", OrderIndex: 3, BodyMD: lesson4Pitfalls},
		},
		Tasks: []content.Task{
			{
				Title:       "Площадь круга",
				OrderIndex:  0,
				Points:      15,
				PromptMD:    lesson4Task1Prompt,
				StarterCode: lesson4Task1Starter,
				TestsGo:     lesson4Task1Tests,
			},
		},
	}
}

const lesson4Body = `# Операторы в Go

Операторы — это специальные символы, которые выполняют операции над операндами.`

const lesson4Overview = `## Типы операторов

### Арифметические операторы

| Оператор | Описание | Пример |
|----------|----------|--------|
| ` + "`+`" + ` | Сложение | ` + "`5 + 3 = 8`" + ` |
| ` + "`-`" + ` | Вычитание | ` + "`5 - 3 = 2`" + ` |
| ` + "`*`" + ` | Умножение | ` + "`5 * 3 = 15`" + ` |
| ` + "`/`" + ` | Деление | ` + "`5 / 3 = 1`" + ` (целочисленное) |
| ` + "`%`" + ` | Остаток от деления | ` + "`5 % 3 = 2`" + ` |

### Операторы сравнения

| Оператор | Описание | Пример |
|----------|----------|--------|
| ` + "`==`" + ` | Равно | ` + "`5 == 5`" + ` → ` + "`true`" + ` |
| ` + "`!=`" + ` | Не равно | ` + "`5 != 3`" + ` → ` + "`true`" + ` |
| ` + "`<`" + ` | Меньше | ` + "`3 < 5`" + ` → ` + "`true`" + ` |
| ` + "`>`" + ` | Больше | ` + "`5 > 3`" + ` → ` + "`true`" + ` |
| ` + "`<=`" + ` | Меньше или равно | ` + "`3 <= 3`" + ` → ` + "`true`" + ` |
| ` + "`>=`" + ` | Больше или равно | ` + "`5 >= 5`" + ` → ` + "`true`" + ` |

### Логические операторы

| Оператор | Описание | Пример |
|----------|----------|--------|
| ` + "`&&`" + ` | Логическое И | ` + "`true && false`" + ` → ` + "`false`" + ` |
| ` + "`||`" + ` | Логическое ИЛИ | ` + "`true || false`" + ` → ` + "`true`" + ` |
| ` + "`!`" + ` | Логическое НЕ | ` + "`!true`" + ` → ` + "`false`" + ` |`

const lesson4Syntax = `## Операторы присваивания

### Базовые

` + "```go" + `
x := 10      // Инициализация
x = 20       // Присваивание
` + "```" + `

### Составные операторы присваивания

` + "```go" + `
x := 10

x += 5   // x = x + 5 → 15
x -= 3   // x = x - 3 → 12
x *= 2   // x = x * 2 → 24
x /= 4   // x = x / 4 → 6
x %= 4   // x = x % 4 → 2
` + "```" + `

### Инкремент и декремент

` + "```go" + `
x := 5

x++  // x = 6 (только постфиксная форма)
x--  // x = 5

// В Go НЕТ префиксных ++x и --x
// ++x  // ОШИБКА

// Инкремент/декремент — это ОПЕРАТОРЫ, не выражения
// y := x++  // ОШИБКА: нельзя использовать в выражениях
` + "```" + `

### Битовые операторы

` + "```go" + `
a := 5  // 0101 в двоичной
b := 3  // 0011 в двоичной

fmt.Println(a & b)   // 1  (AND:  0001)
fmt.Println(a | b)   // 7  (OR:   0111)
fmt.Println(a ^ b)   // 6  (XOR:  0110)
fmt.Println(a << 1)  // 10 (сдвиг влево: 1010)
fmt.Println(a >> 1)  // 2  (сдвиг вправо: 0010)
` + "```" + ``

const lesson4Examples = `## Практические примеры

### Пример 1: Арифметические операции

` + "```go" + `
package main

import "fmt"

func main() {
    a, b := 17, 5
    
    fmt.Printf("%d + %d = %d\n", a, b, a+b)
    fmt.Printf("%d - %d = %d\n", a, b, a-b)
    fmt.Printf("%d * %d = %d\n", a, b, a*b)
    fmt.Printf("%d / %d = %d\n", a, b, a/b)
    fmt.Printf("%d %% %d = %d\n", a, b, a%b)
}
` + "```" + `

**Вывод:**
` + "```" + `
17 + 5 = 22
17 - 5 = 12
17 * 5 = 85
17 / 5 = 3
17 % 5 = 2
` + "```" + `

### Пример 2: Проверка чётности

` + "```go" + `
package main

import "fmt"

func main() {
    for i := 1; i <= 10; i++ {
        if i % 2 == 0 {
            fmt.Printf("%d — чётное\n", i)
        } else {
            fmt.Printf("%d — нечётное\n", i)
        }
    }
}
` + "```" + `

### Пример 3: Логические операторы

` + "```go" + `
package main

import "fmt"

func main() {
    age := 25
    hasLicense := true
    
    // Проверка нескольких условий
    canDrive := age >= 18 && hasLicense
    fmt.Printf("Может водить: %t\n", canDrive)
    
    // Логическое ИЛИ
    isWeekend := false
    isHoliday := true
    dayOff := isWeekend || isHoliday
    fmt.Printf("Выходной: %t\n", dayOff)
    
    // Отрицание
    isWorking := !dayOff
    fmt.Printf("Рабочий день: %t\n", isWorking)
}
` + "```" + ``

const lesson4Pitfalls = `## Частые ошибки

### 1. Целочисленное деление

` + "```go" + `
result := 5 / 2  // 2, а не 2.5!

// Для точного результата:
result := 5.0 / 2.0  // 2.5
// или
result := float64(5) / float64(2)  // 2.5
` + "```" + `

### 2. Приоритет операторов

` + "```go" + `
result := 2 + 3 * 4  // 14, а не 20
// Умножение имеет более высокий приоритет

result := (2 + 3) * 4  // 20 (используйте скобки)
` + "```" + `

### 3. Короткое замыкание

` + "```go" + `
// && прекращает вычисление, если первый операнд false
// || прекращает вычисление, если первый операнд true

func expensive() bool {
    fmt.Println("Вызвана!")
    return true
}

false && expensive()  // expensive() НЕ будет вызвана
true || expensive()   // expensive() НЕ будет вызвана
` + "```" + `

### 4. Использование = вместо ==

` + "```go" + `
x := 5

// ОШИБКА: это присваивание, не сравнение
if x = 10 {  // syntax error
    
}

// ПРАВИЛЬНО: используйте ==
if x == 10 {
    
}
` + "```" + ``

const lesson4Task1Prompt = `### Задание: Площадь круга

Вычислите площадь круга с радиусом ` + "`5`" + `.

Используйте формулу: **S = π × r²**

Значение π возьмите как ` + "`3.14159`" + `.

**Ожидаемый вывод (с точностью до 2 знаков):**
` + "```" + `
Площадь круга: 78.54
` + "```" + `

**Подсказка:** используйте ` + "`fmt.Printf(\"Площадь круга: %.2f\\n\", area)`" + ``

const lesson4Task1Starter = `package main

import "fmt"

func main() {
	radius := 5.0
	pi := 3.14159
	
	// Вычислите площадь круга
	
	
	fmt.Printf("Площадь круга: %.2f\n", area)
}
`

const lesson4Task1Tests = `package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestCircleArea(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Ошибка: %v", err)
	}
	
	output := out.String()
	if !strings.Contains(output, "78.5") {
		t.Errorf("Неверный результат: %s", output)
	}
}
`

// ============================================================
// УРОК 5: Условные конструкции
// ============================================================

func createLesson5() lessonData {
	return lessonData{
		Lesson: content.Lesson{
			Slug:           "uslovnye-konstruktsii",
			Title:          "Условные конструкции if/else и switch",
			OrderIndex:     5,
			SourceURL:      "https://metanit.com/go/tutorial/2.6.php",
			ReadingTimeMin: 15,
			BodyMD:         lesson5Body,
		},
		Sections: []content.Section{
			{Kind: content.SectionOverview, Title: "Ключевые идеи", OrderIndex: 0, BodyMD: lesson5Overview},
			{Kind: content.SectionSyntax, Title: "Синтаксис", OrderIndex: 1, BodyMD: lesson5Syntax},
			{Kind: content.SectionExamples, Title: "Примеры кода", OrderIndex: 2, BodyMD: lesson5Examples},
			{Kind: content.SectionPitfalls, Title: "Частые ошибки", OrderIndex: 3, BodyMD: lesson5Pitfalls},
		},
		Tasks: []content.Task{
			{
				Title:       "Определение знака числа",
				OrderIndex:  0,
				Points:      15,
				PromptMD:    lesson5Task1Prompt,
				StarterCode: lesson5Task1Starter,
				TestsGo:     lesson5Task1Tests,
			},
			{
				Title:       "Максимум из трёх",
				OrderIndex:  1,
				Points:      20,
				PromptMD:    lesson5Task2Prompt,
				StarterCode: lesson5Task2Starter,
				TestsGo:     lesson5Task2Tests,
			},
		},
	}
}

const lesson5Body = `# Условные конструкции в Go

Условные конструкции позволяют выполнять разный код в зависимости от условий. В Go есть два основных инструмента: ` + "`if/else`" + ` и ` + "`switch`" + `.`

const lesson5Overview = `## Управление потоком выполнения

### if/else

Позволяет выполнить код, если условие истинно:

` + "```go" + `
if условие {
    // выполняется если true
} else {
    // выполняется если false
}
` + "```" + `

### switch

Позволяет выбрать один из нескольких вариантов:

` + "```go" + `
switch значение {
case вариант1:
    // код для варианта 1
case вариант2:
    // код для варианта 2
default:
    // код по умолчанию
}
` + "```" + `

### Особенности Go

| Особенность | Описание |
|-------------|----------|
| Без скобок | Условие пишется **без** круглых скобок |
| Обязательные ` + "`{}`" + ` | Фигурные скобки обязательны даже для одной строки |
| ` + "`{`" + ` на той же строке | Открывающая скобка **должна** быть на той же строке |
| Нет ` + "`break`" + ` в switch | Каждый ` + "`case`" + ` автоматически завершается |
| if с инициализацией | Можно объявить переменную в условии |`

const lesson5Syntax = `## Конструкция if/else

### Простой if

` + "```go" + `
age := 18

if age >= 18 {
    fmt.Println("Совершеннолетний")
}
` + "```" + `

### if-else

` + "```go" + `
age := 16

if age >= 18 {
    fmt.Println("Совершеннолетний")
} else {
    fmt.Println("Несовершеннолетний")
}
` + "```" + `

### if-else if-else

` + "```go" + `
score := 75

if score >= 90 {
    fmt.Println("Отлично")
} else if score >= 70 {
    fmt.Println("Хорошо")
} else if score >= 50 {
    fmt.Println("Удовлетворительно")
} else {
    fmt.Println("Неудовлетворительно")
}
` + "```" + `

### if с инициализацией

` + "```go" + `
// Переменная err видна только внутри if
if err := doSomething(); err != nil {
    fmt.Println("Ошибка:", err)
    return
}
// err здесь НЕ доступна
` + "```" + `

## Конструкция switch

### Базовый switch

` + "```go" + `
day := "Monday"

switch day {
case "Monday":
    fmt.Println("Понедельник")
case "Tuesday":
    fmt.Println("Вторник")
case "Wednesday":
    fmt.Println("Среда")
default:
    fmt.Println("Другой день")
}
` + "```" + `

### switch с несколькими значениями

` + "```go" + `
day := "Saturday"

switch day {
case "Saturday", "Sunday":
    fmt.Println("Выходной!")
default:
    fmt.Println("Рабочий день")
}
` + "```" + `

### switch без выражения (замена if-else)

` + "```go" + `
score := 85

switch {
case score >= 90:
    fmt.Println("A")
case score >= 80:
    fmt.Println("B")
case score >= 70:
    fmt.Println("C")
default:
    fmt.Println("F")
}
` + "```" + `

### fallthrough

` + "```go" + `
x := 1

switch x {
case 1:
    fmt.Println("Один")
    fallthrough  // переход к следующему case
case 2:
    fmt.Println("Два")
}
// Вывод:
// Один
// Два
` + "```" + ``

const lesson5Examples = `## Практические примеры

### Пример 1: Проверка возраста

` + "```go" + `
package main

import "fmt"

func main() {
    age := 25
    
    if age < 0 {
        fmt.Println("Некорректный возраст")
    } else if age < 7 {
        fmt.Println("Дошкольник")
    } else if age < 18 {
        fmt.Println("Школьник")
    } else if age < 23 {
        fmt.Println("Студент")
    } else if age < 65 {
        fmt.Println("Работающий")
    } else {
        fmt.Println("Пенсионер")
    }
}
` + "```" + `

### Пример 2: Калькулятор на switch

` + "```go" + `
package main

import "fmt"

func main() {
    a, b := 10.0, 3.0
    op := "/"
    
    var result float64
    
    switch op {
    case "+":
        result = a + b
    case "-":
        result = a - b
    case "*":
        result = a * b
    case "/":
        if b != 0 {
            result = a / b
        } else {
            fmt.Println("Ошибка: деление на ноль")
            return
        }
    default:
        fmt.Println("Неизвестная операция")
        return
    }
    
    fmt.Printf("%.2f %s %.2f = %.2f\n", a, op, b, result)
}
` + "```" + `

### Пример 3: Проверка типа (type switch)

` + "```go" + `
package main

import "fmt"

func checkType(x interface{}) {
    switch v := x.(type) {
    case int:
        fmt.Printf("int: %d\n", v)
    case string:
        fmt.Printf("string: %s\n", v)
    case bool:
        fmt.Printf("bool: %t\n", v)
    default:
        fmt.Printf("неизвестный тип: %T\n", v)
    }
}

func main() {
    checkType(42)
    checkType("hello")
    checkType(true)
    checkType(3.14)
}
` + "```" + ``

const lesson5Pitfalls = `## Частые ошибки

### 1. else на новой строке

` + "```go" + `
// ОШИБКА: syntax error
if x > 5 {
    fmt.Println("больше")
}
else {  // else должен быть на той же строке
    fmt.Println("меньше или равно")
}

// ПРАВИЛЬНО:
if x > 5 {
    fmt.Println("больше")
} else {
    fmt.Println("меньше или равно")
}
` + "```" + `

### 2. Скобки вокруг условия

` + "```go" + `
// Не рекомендуется (работает, но не идиоматично)
if (x > 5) {
    
}

// ПРАВИЛЬНО:
if x > 5 {
    
}
` + "```" + `

### 3. Присваивание вместо сравнения

` + "```go" + `
// ОШИБКА: это присваивание
if x = 5 {
    
}

// ПРАВИЛЬНО: используйте ==
if x == 5 {
    
}
` + "```" + `

### 4. Забытый break в switch (не нужен!)

` + "```go" + `
// В Go break НЕ нужен — каждый case завершается автоматически
switch day {
case "Monday":
    fmt.Println("Понедельник")
    // break не нужен!
case "Tuesday":
    fmt.Println("Вторник")
}
` + "```" + `

### 5. Область видимости переменной в if

` + "```go" + `
if x := getValue(); x > 0 {
    fmt.Println(x)  // x доступна
}
fmt.Println(x)  // ОШИБКА: x не определена здесь

// Если нужна переменная вне if:
x := getValue()
if x > 0 {
    fmt.Println(x)
}
fmt.Println(x)  // OK
` + "```" + ``

const lesson5Task1Prompt = `### Задание: Знак числа

Напишите программу, которая определяет знак числа ` + "`-7`" + `.

**Требования:**
- Если число положительное — выведите "положительное"
- Если число отрицательное — выведите "отрицательное"
- Если число равно нулю — выведите "ноль"

**Ожидаемый вывод для числа -7:**
` + "```" + `
Число -7 отрицательное
` + "```" + ``

const lesson5Task1Starter = `package main

import "fmt"

func main() {
	num := -7
	
	// Определите знак числа и выведите результат
	
}
`

const lesson5Task1Tests = `package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestSign(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Ошибка: %v", err)
	}
	
	output := strings.ToLower(out.String())
	if !strings.Contains(output, "отрицательн") {
		t.Errorf("Ожидалось 'отрицательное', получено: %s", out.String())
	}
}
`

const lesson5Task2Prompt = `### Задание: Максимум из трёх

Найдите максимальное значение среди трёх чисел: ` + "`15`" + `, ` + "`42`" + `, ` + "`8`" + `.

**Ожидаемый вывод:**
` + "```" + `
Максимум: 42
` + "```" + `

**Подсказка:** используйте вложенные if или несколько сравнений`

const lesson5Task2Starter = `package main

import "fmt"

func main() {
	a, b, c := 15, 42, 8
	
	// Найдите максимальное значение
	var max int
	
	
	fmt.Println("Максимум:", max)
}
`

const lesson5Task2Tests = `package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestMax(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Ошибка: %v", err)
	}
	
	output := out.String()
	if !strings.Contains(output, "42") {
		t.Errorf("Ожидалось '42', получено: %s", output)
	}
}
`
