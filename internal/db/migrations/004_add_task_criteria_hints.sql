-- Добавляем поля criteria и hints в tasks
ALTER TABLE tasks ADD COLUMN criteria TEXT NOT NULL DEFAULT '';
ALTER TABLE tasks ADD COLUMN hints TEXT NOT NULL DEFAULT '';
