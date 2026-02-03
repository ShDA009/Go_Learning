-- Курсы (руководства верхнего уровня)
CREATE TABLE IF NOT EXISTS courses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    icon TEXT NOT NULL DEFAULT '📚',
    order_index INTEGER NOT NULL DEFAULT 0
);

-- Добавляем связь модулей с курсами
ALTER TABLE modules ADD COLUMN course_id INTEGER REFERENCES courses(id);

-- Создаём индекс для быстрого поиска модулей по курсу
CREATE INDEX IF NOT EXISTS idx_modules_course ON modules(course_id);
