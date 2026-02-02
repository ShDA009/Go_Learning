-- Добавляем поле для ожидаемого вывода задания
ALTER TABLE tasks ADD COLUMN expected_output TEXT NOT NULL DEFAULT '';

-- Добавляем поле для проверки паттернов в коде (опционально)
ALTER TABLE tasks ADD COLUMN required_patterns TEXT NOT NULL DEFAULT '';
