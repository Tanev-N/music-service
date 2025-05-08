-- Добавляем столбец description в таблицу playlists
ALTER TABLE playlists ADD COLUMN description VARCHAR(500);

-- Обновляем существующие записи, устанавливая пустую строку для description
UPDATE playlists SET description = '';

-- Создаем индекс для поиска по названию плейлиста
CREATE INDEX IF NOT EXISTS idx_playlists_name ON playlists (name); 