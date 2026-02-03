-- Добавляем тип секции 'theory' в lesson_sections
-- SQLite не поддерживает ALTER CHECK CONSTRAINT, поэтому пересоздаём таблицу

-- 1. Создаём новую таблицу с обновлённым constraint
CREATE TABLE lesson_sections_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK(kind IN ('overview', 'theory', 'syntax', 'examples', 'pitfalls', 'extra')),
    title TEXT NOT NULL,
    body_md TEXT NOT NULL DEFAULT '',
    order_index INTEGER NOT NULL DEFAULT 0
);

-- 2. Копируем данные
INSERT INTO lesson_sections_new SELECT * FROM lesson_sections;

-- 3. Удаляем старую таблицу
DROP TABLE lesson_sections;

-- 4. Переименовываем новую
ALTER TABLE lesson_sections_new RENAME TO lesson_sections;

-- 5. Восстанавливаем индекс
CREATE INDEX IF NOT EXISTS idx_lesson_sections_lesson ON lesson_sections(lesson_id);
