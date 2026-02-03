package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	"golearning/internal/db"
)

func main() {
	dbPath := flag.String("db", "./data.db", "Путь к файлу базы данных SQLite")
	flag.Parse()

	database, err := db.Open(*dbPath)
	if err != nil {
		log.Fatalf("Ошибка открытия БД: %v", err)
	}
	defer database.Close()

	// На всякий случай убеждаемся, что схема актуальна (некоторые таблицы могут быть добавлены миграциями).
	if err := db.Migrate(database); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	demoModuleSlugs := []string{"osnovy", "tipy-dannyh", "upravlenie"}
	demoLessonSlugs := []string{"vvedenie", "peremennye", "tipy-dannyh", "operatory", "uslovnye-konstruktsii"}

	beforeModules := countIn(database, "modules", "slug", demoModuleSlugs)
	beforeLessons := countIn(database, "lessons", "slug", demoLessonSlugs)

	// Удаляем демо-модули: каскадно удалятся уроки/секции/задания/прогресс/заметки.
	deletedModules, err := deleteIn(database, "modules", "slug", demoModuleSlugs)
	if err != nil {
		log.Fatalf("Ошибка удаления демо-модулей: %v", err)
	}

	// На случай, если демо-уроки были созданы без ожидаемых модулей (или модули уже удалены),
	// дополнительно удаляем их по slug.
	deletedLessons, err := deleteIn(database, "lessons", "slug", demoLessonSlugs)
	if err != nil {
		log.Fatalf("Ошибка удаления демо-уроков: %v", err)
	}

	afterModules := countIn(database, "modules", "slug", demoModuleSlugs)
	afterLessons := countIn(database, "lessons", "slug", demoLessonSlugs)

	fmt.Println("✅ Демо-контент очищен")
	fmt.Printf("- modules: было %d, удалено %d, осталось %d\n", beforeModules, deletedModules, afterModules)
	fmt.Printf("- lessons: было %d, удалено %d, осталось %d\n", beforeLessons, deletedLessons, afterLessons)
}

func countIn(dbx *sql.DB, table, col string, values []string) int64 {
	if len(values) == 0 {
		return 0
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(values)), ",")
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s IN (%s)", table, col, placeholders)
	args := make([]any, 0, len(values))
	for _, v := range values {
		args = append(args, v)
	}
	var n int64
	_ = dbx.QueryRow(q, args...).Scan(&n) // best-effort
	return n
}

func deleteIn(dbx *sql.DB, table, col string, values []string) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(values)), ",")
	q := fmt.Sprintf("DELETE FROM %s WHERE %s IN (%s)", table, col, placeholders)
	args := make([]any, 0, len(values))
	for _, v := range values {
		args = append(args, v)
	}

	res, err := dbx.Exec(q, args...)
	if err != nil {
		return 0, err
	}
	ra, _ := res.RowsAffected()
	return ra, nil
}
