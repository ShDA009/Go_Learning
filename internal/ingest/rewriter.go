package ingest

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"golearning/internal/content"
)

// Rewriter преобразует сырой контент в структурированный урок.
type Rewriter interface {
	Rewrite(ctx context.Context, parsed *ParsedContent, meta TOCEntry) (*content.StructuredLesson, error)
}

// LocalRuleBasedRewriter — реализация на основе правил (без LLM).
type LocalRuleBasedRewriter struct{}

// NewLocalRewriter создаёт новый локальный rewriter.
func NewLocalRewriter() *LocalRuleBasedRewriter {
	return &LocalRuleBasedRewriter{}
}

// Rewrite преобразует распарсенный контент в структурированный урок.
func (r *LocalRuleBasedRewriter) Rewrite(ctx context.Context, parsed *ParsedContent, meta TOCEntry) (*content.StructuredLesson, error) {
	lesson := &content.StructuredLesson{
		Title:          parsed.Title,
		ReadingTimeMin: r.estimateReadingTime(parsed),
	}

	// Формируем полный текст урока в Markdown
	var bodyParts []string

	// Заголовок
	bodyParts = append(bodyParts, "# "+parsed.Title)
	bodyParts = append(bodyParts, "")

	// Обзор (первые 2-3 параграфа)
	overview := r.extractOverview(parsed)
	if overview != "" {
		bodyParts = append(bodyParts, "## Обзор")
		bodyParts = append(bodyParts, "")
		bodyParts = append(bodyParts, overview)
		bodyParts = append(bodyParts, "")

		lesson.Sections = append(lesson.Sections, content.Section{
			Kind:       content.SectionOverview,
			Title:      "Ключевые идеи",
			BodyMD:     overview,
			OrderIndex: 0,
		})
	}

	// Синтаксис
	syntax := r.extractSyntax(parsed)
	if syntax != "" {
		bodyParts = append(bodyParts, "## Синтаксис")
		bodyParts = append(bodyParts, "")
		bodyParts = append(bodyParts, syntax)
		bodyParts = append(bodyParts, "")

		lesson.Sections = append(lesson.Sections, content.Section{
			Kind:       content.SectionSyntax,
			Title:      "Синтаксис",
			BodyMD:     syntax,
			OrderIndex: 1,
		})
	}

	// Примеры кода
	examples := r.extractExamples(parsed)
	if examples != "" {
		bodyParts = append(bodyParts, "## Примеры")
		bodyParts = append(bodyParts, "")
		bodyParts = append(bodyParts, examples)
		bodyParts = append(bodyParts, "")

		lesson.Sections = append(lesson.Sections, content.Section{
			Kind:       content.SectionExamples,
			Title:      "Примеры кода",
			BodyMD:     examples,
			OrderIndex: 2,
		})
	}

	// Частые ошибки
	pitfalls := r.extractPitfalls(parsed)
	if pitfalls != "" {
		bodyParts = append(bodyParts, "## Частые ошибки")
		bodyParts = append(bodyParts, "")
		bodyParts = append(bodyParts, pitfalls)
		bodyParts = append(bodyParts, "")

		lesson.Sections = append(lesson.Sections, content.Section{
			Kind:       content.SectionPitfalls,
			Title:      "Частые ошибки",
			BodyMD:     pitfalls,
			OrderIndex: 3,
		})
	}

	// Дополнительная информация (оставшиеся параграфы)
	extra := r.extractExtra(parsed)
	if extra != "" {
		bodyParts = append(bodyParts, "## Дополнительно")
		bodyParts = append(bodyParts, "")
		bodyParts = append(bodyParts, extra)

		lesson.Sections = append(lesson.Sections, content.Section{
			Kind:       content.SectionExtra,
			Title:      "Дополнительно",
			BodyMD:     extra,
			OrderIndex: 4,
		})
	}

	lesson.BodyMD = strings.Join(bodyParts, "\n")

	// Генерируем задания
	lesson.Tasks = r.generateTasks(parsed, meta)

	return lesson, nil
}

// estimateReadingTime оценивает время чтения в минутах.
func (r *LocalRuleBasedRewriter) estimateReadingTime(parsed *ParsedContent) int {
	wordCount := 0
	for _, p := range parsed.Paragraphs {
		wordCount += len(strings.Fields(p))
	}
	// ~200 слов в минуту для технического текста
	minutes := wordCount / 200
	if minutes < 3 {
		minutes = 3
	}
	if minutes > 30 {
		minutes = 30
	}
	return minutes
}

// extractOverview извлекает обзор (первые параграфы с определениями).
func (r *LocalRuleBasedRewriter) extractOverview(parsed *ParsedContent) string {
	var overview []string
	definitionPatterns := []string{
		"представляет", "является", "позволяет", "используется",
		"это", "служит", "предназначен", "определяет",
	}

	count := 0
	for _, p := range parsed.Paragraphs {
		if count >= 3 {
			break
		}
		lower := strings.ToLower(p)

		// Первые параграфы или параграфы с определениями
		isDefinition := false
		for _, pattern := range definitionPatterns {
			if strings.Contains(lower, pattern) {
				isDefinition = true
				break
			}
		}

		if count < 2 || isDefinition {
			overview = append(overview, p)
			count++
		}
	}

	return strings.Join(overview, "\n\n")
}

// extractSyntax извлекает информацию о синтаксисе.
func (r *LocalRuleBasedRewriter) extractSyntax(parsed *ParsedContent) string {
	var syntax []string

	// Ищем параграфы с описанием синтаксиса
	syntaxPatterns := []string{
		"синтаксис", "имеет вид", "записывается", "объявляется",
		"определяется", "формат", "структура", "шаблон",
	}

	for _, p := range parsed.Paragraphs {
		lower := strings.ToLower(p)
		for _, pattern := range syntaxPatterns {
			if strings.Contains(lower, pattern) {
				syntax = append(syntax, p)
				break
			}
		}
	}

	// Добавляем первый блок кода как пример синтаксиса
	if len(parsed.CodeBlocks) > 0 {
		cb := parsed.CodeBlocks[0]
		syntax = append(syntax, "")
		syntax = append(syntax, "```"+cb.Language)
		syntax = append(syntax, cb.Code)
		syntax = append(syntax, "```")
	}

	// Добавляем списки (часто содержат ключевые слова/операторы)
	for _, list := range parsed.Lists {
		if len(syntax) < 5 { // Ограничиваем размер
			syntax = append(syntax, "")
			syntax = append(syntax, list)
		}
	}

	return strings.Join(syntax, "\n")
}

// extractExamples извлекает примеры кода с пояснениями.
func (r *LocalRuleBasedRewriter) extractExamples(parsed *ParsedContent) string {
	var examples []string

	for i, cb := range parsed.CodeBlocks {
		if i >= 4 { // Максимум 4 примера
			break
		}

		// Добавляем заголовок примера
		examples = append(examples, fmt.Sprintf("### Пример %d", i+1))
		examples = append(examples, "")

		// Пытаемся найти пояснение к коду
		if i < len(parsed.Paragraphs) {
			// Берём ближайший параграф
			for _, p := range parsed.Paragraphs {
				lower := strings.ToLower(p)
				if strings.Contains(lower, "пример") ||
					strings.Contains(lower, "рассмотрим") ||
					strings.Contains(lower, "следующий") {
					examples = append(examples, p)
					examples = append(examples, "")
					break
				}
			}
		}

		// Добавляем код
		examples = append(examples, "```"+cb.Language)
		examples = append(examples, cb.Code)
		examples = append(examples, "```")
		examples = append(examples, "")
	}

	return strings.Join(examples, "\n")
}

// extractPitfalls извлекает информацию о частых ошибках.
func (r *LocalRuleBasedRewriter) extractPitfalls(parsed *ParsedContent) string {
	var pitfalls []string

	pitfallPatterns := []string{
		"ошибка", "ошибки", "нельзя", "важно", "следует помнить",
		"внимание", "осторожно", "не рекомендуется", "избегайте",
		"проблема", "неправильно", "ограничение", "исключение",
	}

	for _, p := range parsed.Paragraphs {
		lower := strings.ToLower(p)
		for _, pattern := range pitfallPatterns {
			if strings.Contains(lower, pattern) {
				pitfalls = append(pitfalls, "- "+p)
				break
			}
		}
	}

	if len(pitfalls) == 0 {
		// Генерируем общие советы
		pitfalls = append(pitfalls, "- Внимательно следите за типами данных")
		pitfalls = append(pitfalls, "- Не забывайте про обработку ошибок")
		pitfalls = append(pitfalls, "- Используйте понятные имена переменных")
	}

	return strings.Join(pitfalls, "\n")
}

// extractExtra извлекает дополнительную информацию.
func (r *LocalRuleBasedRewriter) extractExtra(parsed *ParsedContent) string {
	var extra []string

	// Собираем оставшиеся параграфы, которые не попали в другие секции
	usedPatterns := []string{
		"представляет", "является", "синтаксис", "ошибка", "нельзя", "важно",
	}

	for _, p := range parsed.Paragraphs[min(3, len(parsed.Paragraphs)):] {
		lower := strings.ToLower(p)
		isUsed := false
		for _, pattern := range usedPatterns {
			if strings.Contains(lower, pattern) {
				isUsed = true
				break
			}
		}
		if !isUsed && len(extra) < 5 {
			extra = append(extra, p)
		}
	}

	return strings.Join(extra, "\n\n")
}

// generateTasks генерирует практические задания.
func (r *LocalRuleBasedRewriter) generateTasks(parsed *ParsedContent, meta TOCEntry) []content.Task {
	tasks := []content.Task{}

	slug := slugify(parsed.Title)

	// Задание 1: Разогрев — простая задача на синтаксис
	tasks = append(tasks, content.Task{
		Title:    "Разогрев: базовый синтаксис",
		PromptMD: r.generateWarmupPrompt(parsed),
		StarterCode: `package main

import "fmt"

func main() {
	// Ваш код здесь
	
}`,
		TestsGo:    r.generateWarmupTests(parsed, slug),
		Points:     10,
		OrderIndex: 0,
	})

	// Задание 2: На понимание
	if len(parsed.CodeBlocks) > 0 {
		tasks = append(tasks, content.Task{
			Title:    "На понимание: применяем концепцию",
			PromptMD: r.generateUnderstandingPrompt(parsed),
			StarterCode: `package main

import "fmt"

func main() {
	// Напишите решение здесь
	
}`,
			TestsGo:    r.generateUnderstandingTests(parsed, slug),
			Points:     15,
			OrderIndex: 1,
		})
	}

	// Задание 3: На ошибку
	tasks = append(tasks, content.Task{
		Title:    "Найди ошибку",
		PromptMD: r.generateDebugPrompt(parsed),
		StarterCode: `package main

import "fmt"

func main() {
	// Исправьте код ниже
	
}`,
		TestsGo:    r.generateDebugTests(parsed, slug),
		Points:     20,
		OrderIndex: 2,
	})

	return tasks
}

// generateWarmupPrompt генерирует промпт для разогревочного задания.
func (r *LocalRuleBasedRewriter) generateWarmupPrompt(parsed *ParsedContent) string {
	// Базовый шаблон
	return fmt.Sprintf(`Напишите программу, которая выводит на экран текст "Hello, Go!".

**Требования:**
- Используйте функцию %sfmt.Println%s
- Программа должна вывести ровно одну строку

**Пример вывода:**
%s
Hello, Go!
%s`, "`", "`", "```", "```")
}

// generateWarmupTests генерирует тесты для разогревочного задания.
func (r *LocalRuleBasedRewriter) generateWarmupTests(parsed *ParsedContent, slug string) string {
	return `package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
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
}`
}

// generateUnderstandingPrompt генерирует промпт для задания на понимание.
func (r *LocalRuleBasedRewriter) generateUnderstandingPrompt(parsed *ParsedContent) string {
	return `Напишите функцию, которая принимает число и возвращает его квадрат.

**Требования:**
- Создайте функцию ` + "`square(n int) int`" + `
- Функция должна возвращать n * n
- В main() вызовите функцию с числом 5 и выведите результат

**Пример вывода:**
` + "```" + `
25
` + "```"
}

// generateUnderstandingTests генерирует тесты для задания на понимание.
func (r *LocalRuleBasedRewriter) generateUnderstandingTests(parsed *ParsedContent, slug string) string {
	return `package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestSquare(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Программа завершилась с ошибкой: %v", err)
	}
	
	output := strings.TrimSpace(out.String())
	if output != "25" {
		t.Errorf("Ожидалось 25, получено %q", output)
	}
}`
}

// generateDebugPrompt генерирует промпт для задания на поиск ошибки.
func (r *LocalRuleBasedRewriter) generateDebugPrompt(parsed *ParsedContent) string {
	return `В коде ниже есть ошибка. Найдите и исправьте её.

` + "```go" + `
package main

import "fmt"

func main() {
	var x int = 10
	var y int = 0
	fmt.Println(x / y)  // Проблема здесь
}
` + "```" + `

**Задание:**
- Исправьте код так, чтобы он корректно обрабатывал деление на ноль
- Программа должна выводить сообщение "Деление на ноль!" если делитель равен 0
- Если делитель не равен 0, выводите результат деления`
}

// generateDebugTests генерирует тесты для задания на поиск ошибки.
func (r *LocalRuleBasedRewriter) generateDebugTests(parsed *ParsedContent, slug string) string {
	return `package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestDivision(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	if err := cmd.Run(); err != nil {
		t.Fatalf("Программа завершилась с ошибкой: %v", err)
	}
	
	output := strings.TrimSpace(out.String())
	if !strings.Contains(output, "Деление на ноль") && !strings.Contains(output, "деление на ноль") {
		t.Errorf("Программа должна обрабатывать деление на ноль")
	}
}`
}

// slugify преобразует строку в slug.
func slugify(s string) string {
	// Транслитерация русских букв
	translit := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "",
		'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo",
		'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M",
		'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U",
		'Ф': "F", 'Х': "H", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch", 'Ъ': "",
		'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
	}

	var result strings.Builder
	for _, r := range s {
		if tr, ok := translit[r]; ok {
			result.WriteString(tr)
		} else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}

	// Приводим к нижнему регистру и убираем лишние дефисы
	slug := strings.ToLower(result.String())
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	if slug == "" {
		slug = "lesson"
	}

	return slug
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
