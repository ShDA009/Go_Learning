-- Модули (разделы курса)
CREATE TABLE IF NOT EXISTS modules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 0
);

-- Уроки
CREATE TABLE IF NOT EXISTS lessons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    module_id INTEGER NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    order_index INTEGER NOT NULL DEFAULT 0,
    source_url TEXT,
    body_md TEXT NOT NULL DEFAULT '',
    reading_time_min INTEGER NOT NULL DEFAULT 5,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_lessons_module ON lessons(module_id);
CREATE INDEX IF NOT EXISTS idx_lessons_slug ON lessons(slug);

-- Секции урока (overview, syntax, examples, pitfalls, extra)
CREATE TABLE IF NOT EXISTS lesson_sections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK(kind IN ('overview', 'syntax', 'examples', 'pitfalls', 'extra')),
    title TEXT NOT NULL,
    body_md TEXT NOT NULL DEFAULT '',
    order_index INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_lesson_sections_lesson ON lesson_sections(lesson_id);

-- Практические задания
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    prompt_md TEXT NOT NULL,
    starter_code TEXT NOT NULL DEFAULT '',
    tests_go TEXT NOT NULL DEFAULT '',
    points INTEGER NOT NULL DEFAULT 10,
    order_index INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_tasks_lesson ON tasks(lesson_id);

-- Прогресс пользователя (single-user)
CREATE TABLE IF NOT EXISTS progress (
    lesson_id INTEGER PRIMARY KEY REFERENCES lessons(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'new' CHECK(status IN ('new', 'reading', 'done')),
    practice_done INTEGER NOT NULL DEFAULT 0,
    points_earned INTEGER NOT NULL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Заметки пользователя
CREATE TABLE IF NOT EXISTS notes (
    lesson_id INTEGER PRIMARY KEY REFERENCES lessons(id) ON DELETE CASCADE,
    note_md TEXT NOT NULL DEFAULT '',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- История отправок решений
CREATE TABLE IF NOT EXISTS submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'success', 'error', 'timeout')),
    stdout TEXT,
    stderr TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_submissions_task ON submissions(task_id);

-- Полнотекстовый поиск по урокам
CREATE VIRTUAL TABLE IF NOT EXISTS lessons_fts USING fts5(
    title, 
    body_md,
    content='lessons',
    content_rowid='id'
);

-- Триггеры для синхронизации FTS
CREATE TRIGGER IF NOT EXISTS lessons_ai AFTER INSERT ON lessons BEGIN
    INSERT INTO lessons_fts(rowid, title, body_md) VALUES (new.id, new.title, new.body_md);
END;

CREATE TRIGGER IF NOT EXISTS lessons_ad AFTER DELETE ON lessons BEGIN
    INSERT INTO lessons_fts(lessons_fts, rowid, title, body_md) VALUES('delete', old.id, old.title, old.body_md);
END;

CREATE TRIGGER IF NOT EXISTS lessons_au AFTER UPDATE ON lessons BEGIN
    INSERT INTO lessons_fts(lessons_fts, rowid, title, body_md) VALUES('delete', old.id, old.title, old.body_md);
    INSERT INTO lessons_fts(rowid, title, body_md) VALUES (new.id, new.title, new.body_md);
END;
