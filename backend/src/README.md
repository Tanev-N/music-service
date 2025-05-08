# Технологический UI для тестирования музыкального сервиса

Это консольное приложение предназначено для системного тестирования всех Use Case из ТЗ музыкального сервиса.

## Запуск приложения

```bash
go run src/main.go
```

## Доступные команды

### Управление пользователями
- `register <login> <password>` - Регистрация нового пользователя
- `login <login> <password>` - Вход в систему
- `profile` - Просмотр профиля текущего пользователя
- `update-permissions <userID> <role>` - Обновление прав пользователя (admin или user)

### Управление жанрами
- `create-genre <name>` - Создание нового жанра
- `list-genres` - Список всех жанров

### Управление альбомами
- `create-album <title> <artist> <coverURL>` - Создание нового альбома
- `album-details <albumID>` - Получение деталей альбома

### Управление плейлистами
- `create-playlist <name> <description>` - Создание нового плейлиста
- `add-to-playlist <playlistID> <trackID>` - Добавление трека в плейлист
- `get-playlist <playlistID>` - Получение плейлиста с треками

### Управление треками
- `play-track <trackID>` - Воспроизведение трека (запись в историю)

### Управление историей
- `get-history` - Получение истории прослушивания

## Примеры использования

```bash
# Регистрация и вход
register user1 password123
login user1 password123

# Создание жанра
create-genre "Rock"

# Создание альбома
create-album "Greatest Hits" "Artist Name" "http://example.com/cover.jpg"

# Создание плейлиста
create-playlist "My Favorites" "Мои любимые треки"

# Все доступные команды можно посмотреть с помощью
help
``` 