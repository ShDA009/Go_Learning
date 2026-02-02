package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golearning/internal/content"
	"golearning/internal/db"
	"golearning/internal/ingest"
)

func main() {
	// Флаги командной строки
	dbPath := flag.String("db", "./data.db", "Путь к файлу базы данных SQLite")
	limit := flag.Int("limit", 0, "Ограничение количества уроков (0 = без ограничения)")
	baseURL := flag.String("url", "https://metanit.com/go/tutorial", "Базовый URL для импорта")
	demo := flag.Bool("demo", false, "Использовать демонстрационные данные вместо загрузки")
	dir := flag.String("dir", "", "Директория с Markdown файлами уроков (например: ./lessons_ai)")
	flag.Parse()

	log.Printf("Go Learning — Импорт контента")
	log.Printf("База данных: %s", *dbPath)

	// Контекст с обработкой сигналов
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Получен сигнал завершения, останавливаем...")
		cancel()
	}()

	// Открываем базу данных
	database, err := db.Open(*dbPath)
	if err != nil {
		log.Fatalf("Ошибка открытия БД: %v", err)
	}
	defer database.Close()

	// Применяем миграции
	if err := db.Migrate(database); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	repo := content.NewRepository(database)

	// Выбираем режим импорта
	switch {
	case *dir != "":
		// Импорт из директории с Markdown файлами
		log.Printf("Режим: импорт из директории %s", *dir)
		importer := ingest.NewMarkdownImporter(repo, *dir)
		if err := importer.Import(ctx); err != nil {
			log.Fatalf("Ошибка импорта: %v", err)
		}

	case *demo:
		// Демонстрационные данные
		log.Println("Режим: демонстрационные данные")
		demoData := ingest.NewDemoData(repo)
		if err := demoData.Seed(ctx); err != nil {
			log.Fatalf("Ошибка создания демо-данных: %v", err)
		}

	default:
		// Импорт с сайта
		log.Printf("Источник: %s", *baseURL)

		// Создаём компоненты pipeline
		crawler := ingest.NewCrawler(*baseURL)
		parser := ingest.NewParser()
		rewriter := ingest.NewLocalRewriter()

		// Создаём и запускаем pipeline
		pipeline := ingest.NewPipeline(crawler, parser, rewriter, repo)

		if err := pipeline.Run(ctx, *limit); err != nil {
			if ctx.Err() != nil {
				log.Println("Импорт прерван пользователем")
				os.Exit(0)
			}

			log.Printf("Ошибка загрузки с сайта: %v", err)
			log.Println("Переключаемся на демонстрационные данные...")

			demoData := ingest.NewDemoData(repo)
			if err := demoData.Seed(ctx); err != nil {
				log.Fatalf("Ошибка создания демо-данных: %v", err)
			}
		}
	}

	log.Println("Импорт успешно завершён!")
}
