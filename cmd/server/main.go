package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golearning/internal/content"
	"golearning/internal/db"
	"golearning/internal/practice"
	"golearning/internal/progress"
	"golearning/internal/web"
)

func main() {
	// Флаги командной строки
	dbPath := flag.String("db", "./data.db", "Путь к файлу базы данных SQLite")
	addr := flag.String("addr", ":8080", "Адрес для прослушивания")
	flag.Parse()

	log.Printf("Go Learning — Веб-сервер")
	log.Printf("База данных: %s", *dbPath)
	log.Printf("Адрес: %s", *addr)

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

	// Создаём репозитории
	contentRepo := content.NewRepository(database)
	progressRepo := progress.NewRepository(database)

	// Создаём runner и checker
	runner := practice.NewLocalRunner()
	checker := practice.NewChecker(runner, contentRepo, progressRepo)

	// Создаём HTTP-сервер
	server, err := web.NewServer(contentRepo, progressRepo, checker)
	if err != nil {
		log.Fatalf("Ошибка создания сервера: %v", err)
	}

	httpServer := &http.Server{
		Addr:         *addr,
		Handler:      server.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Получен сигнал завершения, останавливаем сервер...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("Ошибка остановки сервера: %v", err)
		}

		close(done)
	}()

	log.Printf("Сервер запущен: http://localhost%s", *addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	<-done
	log.Println("Сервер остановлен")
}
