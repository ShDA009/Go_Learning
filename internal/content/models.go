package content

import "time"

// SectionKind — тип секции урока.
type SectionKind string

const (
	SectionOverview SectionKind = "overview"
	SectionSyntax   SectionKind = "syntax"
	SectionExamples SectionKind = "examples"
	SectionPitfalls SectionKind = "pitfalls"
	SectionExtra    SectionKind = "extra"
)

// Module — раздел курса (например, "Основы", "Функции", "Структуры").
type Module struct {
	ID         int64
	Slug       string
	Title      string
	OrderIndex int
}

// Lesson — урок в модуле.
type Lesson struct {
	ID             int64
	ModuleID       int64
	Slug           string
	Title          string
	OrderIndex     int
	SourceURL      string
	BodyMD         string
	ReadingTimeMin int
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Связанные данные (заполняются при необходимости)
	Module   *Module
	Sections []Section
	Tasks    []Task
}

// Section — секция урока (overview, syntax, examples и т.д.).
type Section struct {
	ID         int64
	LessonID   int64
	Kind       SectionKind
	Title      string
	BodyMD     string
	OrderIndex int
}

// Task — практическое задание.
type Task struct {
	ID               int64
	LessonID         int64
	Title            string
	PromptMD         string
	StarterCode      string
	TestsGo          string
	ExpectedOutput   string // Ожидаемый вывод программы
	RequiredPatterns string // Паттерны, которые должны быть в коде (разделённые |)
	Points           int
	OrderIndex       int
}

// StructuredLesson — структурированный урок после обработки rewriter.
type StructuredLesson struct {
	Title          string
	BodyMD         string
	ReadingTimeMin int
	Sections       []Section
	Tasks          []Task
}

// SearchResult — результат поиска.
type SearchResult struct {
	LessonID int64
	Slug     string
	Title    string
	Snippet  string
	Rank     float64
}
