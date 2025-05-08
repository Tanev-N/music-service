-- Создаём функцию для генерации пути хранения трека на основе его ID
CREATE OR REPLACE FUNCTION track_file_path(track_id UUID) 
RETURNS TEXT AS $$
BEGIN
  RETURN CONCAT(
    SUBSTRING(track_id::text, 1, 2), '/', 
    SUBSTRING(track_id::text, 3, 2), '/', 
    track_id::text, '.mp3'
  );
END;
$$ LANGUAGE plpgsql;

-- Создаём триггерную функцию для автоматического обновления пути файла при вставке/обновлении
CREATE OR REPLACE FUNCTION update_track_file_path() 
RETURNS TRIGGER AS $$
BEGIN
  -- Используем функцию трека для генерации корректного пути
  NEW.file_path := track_file_path(NEW.id);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаём триггер для новых записей
DROP TRIGGER IF EXISTS track_file_path_trigger ON tracks;
CREATE TRIGGER track_file_path_trigger
BEFORE INSERT ON tracks
FOR EACH ROW
WHEN (NEW.file_path IS NULL OR NEW.file_path = '')
EXECUTE FUNCTION update_track_file_path();

-- Обновляем существующие файловые пути
UPDATE tracks
SET file_path = track_file_path(id)
WHERE file_path NOT LIKE '%/%/%'; 