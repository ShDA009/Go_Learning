package ingest

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"golearning/internal/content"
)

// Pipeline — конвейер импорта контента.
type Pipeline struct {
	crawler  *Crawler
	parser   *Parser
	rewriter Rewriter
	repo     *content.Repository
}

// NewPipeline создаёт новый pipeline.
func NewPipeline(crawler *Crawler, parser *Parser, rewriter Rewriter, repo *content.Repository) *Pipeline {
	return &Pipeline{
		crawler:  crawler,
		parser:   parser,
		rewriter: rewriter,
		repo:     repo,
	}
}

// Run запускает импорт контента.
func (p *Pipeline) Run(ctx context.Context, limit int) error {
	log.Println("Получение оглавления...")

	toc, err := p.crawler.FetchTOC(ctx)
	if err != nil {
		return fmt.Errorf("fetch TOC: %w", err)
	}

	log.Printf("Найдено %d уроков", len(toc))

	if limit > 0 && limit < len(toc) {
		toc = toc[:limit]
		log.Printf("Ограничение: импортируем только %d уроков", limit)
	}

	// Группируем по модулям
	modules := p.groupByModules(toc)

	for _, mod := range modules {
		// Создаём или обновляем модуль
		if err := p.repo.CreateModule(mod.Module); err != nil {
			return fmt.Errorf("create module %s: %w", mod.Module.Slug, err)
		}
		log.Printf("Модуль: %s (ID=%d)", mod.Module.Title, mod.Module.ID)

		for _, entry := range mod.Entries {
			if err := p.processLesson(ctx, entry, mod.Module.ID); err != nil {
				log.Printf("Ошибка обработки урока %s: %v", entry.URL, err)
				continue
			}

			// Пауза между запросами
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(500 * time.Millisecond):
			}
		}
	}

	log.Println("Импорт завершён!")
	return nil
}

// ModuleGroup — группа уроков в модуле.
type ModuleGroup struct {
	Module  *content.Module
	Entries []TOCEntry
}

// groupByModules группирует уроки по модулям.
func (p *Pipeline) groupByModules(entries []TOCEntry) []ModuleGroup {
	moduleMap := make(map[string]*ModuleGroup)
	var order []string

	for _, entry := range entries {
		moduleSlug := entry.ModuleSlug
		if moduleSlug == "" {
			moduleSlug = "osnovy" // По умолчанию
		}

		if _, exists := moduleMap[moduleSlug]; !exists {
			moduleTitle := p.moduleSlugToTitle(moduleSlug)
			moduleMap[moduleSlug] = &ModuleGroup{
				Module: &content.Module{
					Slug:       moduleSlug,
					Title:      moduleTitle,
					OrderIndex: len(order),
				},
			}
			order = append(order, moduleSlug)
		}

		moduleMap[moduleSlug].Entries = append(moduleMap[moduleSlug].Entries, entry)
	}

	// Возвращаем в порядке добавления
	result := make([]ModuleGroup, 0, len(order))
	for _, slug := range order {
		result = append(result, *moduleMap[slug])
	}

	return result
}

// moduleSlugToTitle преобразует slug модуля в заголовок.
func (p *Pipeline) moduleSlugToTitle(slug string) string {
	titles := map[string]string{
		"osnovy":           "Основы Go",
		"osnovy-yazyka":    "Основы языка",
		"peremennye":       "Переменные и типы данных",
		"operatory":        "Операторы",
		"uslovnye":         "Условные конструкции",
		"tsikly":           "Циклы",
		"funktsii":         "Функции",
		"massivy":          "Массивы и срезы",
		"map":              "Отображения (map)",
		"struktury":        "Структуры",
		"interfeysy":       "Интерфейсы",
		"obrabotka-oshibok": "Обработка ошибок",
		"goroutiny":        "Горутины и каналы",
		"pakety":           "Пакеты и модули",
		"rabota-s-faylami": "Работа с файлами",
		"":                 "Основы Go",
	}

	if title, ok := titles[slug]; ok {
		return title
	}

	// Убираем дефисы и делаем первую букву заглавной
	title := strings.ReplaceAll(slug, "-", " ")
	if len(title) > 0 {
		title = strings.ToUpper(title[:1]) + title[1:]
	}
	return title
}

// processLesson обрабатывает один урок.
func (p *Pipeline) processLesson(ctx context.Context, entry TOCEntry, moduleID int64) error {
	log.Printf("  Загрузка: %s", entry.Title)

	// Скачиваем страницу
	html, err := p.crawler.FetchPage(ctx, entry.URL)
	if err != nil {
		return fmt.Errorf("fetch page: %w", err)
	}

	// Парсим HTML
	parsed, err := p.parser.Parse(html)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	// Если заголовок пустой, берём из оглавления
	if parsed.Title == "" {
		parsed.Title = entry.Title
	}

	// Преобразуем в структурированный урок
	structured, err := p.rewriter.Rewrite(ctx, parsed, entry)
	if err != nil {
		return fmt.Errorf("rewrite: %w", err)
	}

	// Генерируем slug
	slug := slugify(parsed.Title)

	// Сохраняем урок
	lesson := &content.Lesson{
		ModuleID:       moduleID,
		Slug:           slug,
		Title:          structured.Title,
		OrderIndex:     entry.OrderIndex,
		SourceURL:      entry.URL,
		BodyMD:         structured.BodyMD,
		ReadingTimeMin: structured.ReadingTimeMin,
	}

	if err := p.repo.CreateLesson(lesson); err != nil {
		return fmt.Errorf("create lesson: %w", err)
	}

	log.Printf("    -> Урок сохранён: %s (ID=%d)", lesson.Slug, lesson.ID)

	// Удаляем старые секции и задания
	p.repo.DeleteSectionsByLessonID(lesson.ID)
	p.repo.DeleteTasksByLessonID(lesson.ID)

	// Сохраняем секции
	for i := range structured.Sections {
		structured.Sections[i].LessonID = lesson.ID
		if err := p.repo.CreateSection(&structured.Sections[i]); err != nil {
			log.Printf("    Ошибка сохранения секции: %v", err)
		}
	}
	log.Printf("    -> Секций: %d", len(structured.Sections))

	// Сохраняем задания
	for i := range structured.Tasks {
		structured.Tasks[i].LessonID = lesson.ID
		if err := p.repo.CreateTask(&structured.Tasks[i]); err != nil {
			log.Printf("    Ошибка сохранения задания: %v", err)
		}
	}
	log.Printf("    -> Заданий: %d", len(structured.Tasks))

	return nil
}
