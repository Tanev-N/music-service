package ui

import (
	"bufio"
	"fmt"
	"music-service/internal/models"
	"music-service/internal/usecases/interfaces"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ConsoleApp struct {
	userUseCase     interfaces.UserUseCase
	trackUseCase    interfaces.TrackUseCase
	albumUseCase    interfaces.AlbumUseCase
	playlistUseCase interfaces.PlaylistUseCase
	genreUseCase    interfaces.GenreUseCase
	historyUseCase  interfaces.HistoryUseCase

	currentUser    *models.User
	currentSession *models.Session
}

func NewConsoleApp(
	userUseCase interfaces.UserUseCase,
	trackUseCase interfaces.TrackUseCase,
	albumUseCase interfaces.AlbumUseCase,
	playlistUseCase interfaces.PlaylistUseCase,
	genreUseCase interfaces.GenreUseCase,
	historyUseCase interfaces.HistoryUseCase,
) *ConsoleApp {
	return &ConsoleApp{
		userUseCase:     userUseCase,
		trackUseCase:    trackUseCase,
		albumUseCase:    albumUseCase,
		playlistUseCase: playlistUseCase,
		genreUseCase:    genreUseCase,
		historyUseCase:  historyUseCase,
	}
}

func (app *ConsoleApp) StartInteractive() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Музыкальный сервис - Технологический UI ===")
	fmt.Println("Введите 'help' для получения списка команд или 'exit' для выхода")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		command := scanner.Text()
		if command == "exit" {
			break
		}

		output, err := app.ExecuteCommand(command)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
		} else {
			fmt.Println(output)
		}
	}
}

func (app *ConsoleApp) ExecuteCommand(command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", nil
	}

	switch parts[0] {
	case "help":
		return app.showHelp(), nil
	case "register":
		return app.register(parts[1:])
	case "login":
		return app.login(parts[1:])
	case "profile":
		return app.getProfile()
	case "update-permissions":
		return app.updatePermissions(parts[1:])
	case "create-genre":
		return app.createGenre(parts[1:])
	case "list-genres":
		return app.listGenres()
	case "create-album":
		return app.createAlbum(parts[1:])
	case "album-details":
		return app.getAlbumDetails(parts[1:])
	case "add-track-to-album":
		return app.addTrackToAlbum(parts[1:])
	case "remove-track-from-album":
		return app.removeTrackFromAlbum(parts[1:])
	case "create-playlist":
		return app.createPlaylist(parts[1:])
	case "add-to-playlist":
		return app.addTrackToPlaylist(parts[1:])
	case "remove-from-playlist":
		return app.removeTrackFromPlaylist(parts[1:])
	case "get-playlist":
		return app.getPlaylist(parts[1:])
	case "play-track":
		return app.playTrack(parts[1:])
	case "delete-track":
		return app.deleteTrack(parts[1:])
	case "get-history":
		return app.getHistory()
	default:
		return "", fmt.Errorf("неизвестная команда: %s", parts[0])
	}
}

func (app *ConsoleApp) showHelp() string {
	return `
Доступные команды:
  Управление пользователями:
    register <login> <password>       - Регистрация нового пользователя
    login <login> <password>          - Вход в систему
    profile                           - Просмотр профиля текущего пользователя
    update-permissions <userID> <role> - Обновление прав пользователя (admin или user)

  Управление жанрами:
    create-genre <name>               - Создание нового жанра
    list-genres                       - Список всех жанров

  Управление альбомами:
    create-album <title> <artist> <coverURL> - Создание нового альбома
    album-details <albumID>           - Получение деталей альбома
    add-track-to-album <albumID> <trackID> - Добавление трека в альбом
    remove-track-from-album <albumID> <trackID> - Удаление трека из альбома

  Управление плейлистами:
    create-playlist <name> <description> - Создание нового плейлиста
    add-to-playlist <playlistID> <trackID> - Добавление трека в плейлист
    remove-from-playlist <playlistID> <trackID> - Удаление трека из плейлиста
    get-playlist <playlistID>          - Получение плейлиста с треками

  Управление треками:
    play-track <trackID>              - Воспроизведение трека (запись в историю)
    delete-track <trackID>            - Удаление трека
    
  Управление историей:
    get-history                       - Получение истории прослушивания
`
}

func (app *ConsoleApp) register(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("использование: register <login> <password>")
	}

	login, password := args[0], args[1]
	user, err := app.userUseCase.Register(login, password)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Пользователь %s успешно зарегистрирован (ID: %s)", user.Login, user.ID), nil
}

func (app *ConsoleApp) login(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("использование: login <login> <password>")
	}

	login, password := args[0], args[1]
	user, session, err := app.userUseCase.Authenticate(login, password)
	if err != nil {
		return "", err
	}

	app.currentUser = user
	app.currentSession = session

	return fmt.Sprintf("Успешный вход пользователя %s (ID: %s)", user.Login, user.ID), nil
}

func (app *ConsoleApp) getProfile() (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	user, err := app.userUseCase.GetUserProfile(app.currentUser.ID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Профиль пользователя:\nID: %s\nЛогин: %s\nПрава: %s",
		user.ID, user.Login, user.Permission), nil
}

func (app *ConsoleApp) updatePermissions(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if app.currentUser.Permission != models.AdminPermission {
		return "", fmt.Errorf("требуются права администратора")
	}

	if len(args) < 2 {
		return "", fmt.Errorf("использование: update-permissions <userID> <role>")
	}

	userID, err := uuid.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID: %v", err)
	}

	permission := models.Permission(args[1])
	if permission != models.AdminPermission && permission != models.UserPermission {
		return "", fmt.Errorf("недопустимые права: %s (разрешены: admin, user)", permission)
	}

	err = app.userUseCase.UpdatePermissions(userID, permission)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Права пользователя %s обновлены на %s", userID, permission), nil
}

func (app *ConsoleApp) createGenre(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 1 {
		return "", fmt.Errorf("использование: create-genre <name>")
	}

	name := args[0]
	genre, err := app.genreUseCase.CreateGenre(name)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Жанр '%s' успешно создан (ID: %s)", genre.Name, genre.ID), nil
}

func (app *ConsoleApp) listGenres() (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	genres, err := app.genreUseCase.ListAllGenres()
	if err != nil {
		return "", err
	}

	result := "Список жанров:\n"
	for i, genre := range genres {
		result += fmt.Sprintf("%d. %s (ID: %s)\n", i+1, genre.Name, genre.ID)
	}

	return result, nil
}

func (app *ConsoleApp) createAlbum(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 3 {
		return "", fmt.Errorf("использование: create-album <title> <artist> <coverURL>")
	}

	title, _, coverURL := args[0], args[1], args[2]
	releaseDate := time.Now() // Упрощение для примера

	album, err := app.albumUseCase.CreateAlbum(title, releaseDate, coverURL)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Альбом '%s' успешно создан (ID: %s)", album.Title, album.ID), nil
}

func (app *ConsoleApp) getAlbumDetails(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 1 {
		return "", fmt.Errorf("использование: album-details <albumID>")
	}

	albumID, err := uuid.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID: %v", err)
	}

	album, tracks, err := app.albumUseCase.GetAlbumDetails(albumID)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Альбом: %s\nОбложка: %s\n\nТреки:\n",
		album.Title, album.CoverURL)

	for i, track := range tracks {
		result += fmt.Sprintf("%d. %s (ID: %s, Длительность: %d сек)\n",
			i+1, track.Title, track.ID, track.Duration)
	}

	return result, nil
}

func (app *ConsoleApp) addTrackToAlbum(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 2 {
		return "", fmt.Errorf("использование: add-track-to-album <albumID> <trackID>")
	}

	albumID, err := uuid.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID альбома: %v", err)
	}

	trackID, err := uuid.Parse(args[1])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID трека: %v", err)
	}

	err = app.albumUseCase.AddTrackToAlbum(albumID, trackID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Трек %s успешно добавлен в альбом %s", trackID, albumID), nil
}

func (app *ConsoleApp) removeTrackFromAlbum(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 2 {
		return "", fmt.Errorf("использование: remove-track-from-album <albumID> <trackID>")
	}

	albumID, err := uuid.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID альбома: %v", err)
	}

	trackID, err := uuid.Parse(args[1])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID трека: %v", err)
	}

	err = app.albumUseCase.RemoveTrackFromAlbum(albumID, trackID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Трек %s успешно удален из альбома %s", trackID, albumID), nil
}

func (app *ConsoleApp) createPlaylist(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 2 {
		return "", fmt.Errorf("использование: create-playlist <name> <description>")
	}

	name, description := args[0], args[1]
	playlist, err := app.playlistUseCase.CreatePlaylist(app.currentUser.ID, name, description)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Плейлист '%s' успешно создан (ID: %s)", playlist.Name, playlist.ID), nil
}

func (app *ConsoleApp) addTrackToPlaylist(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 2 {
		return "", fmt.Errorf("использование: add-to-playlist <playlistID> <trackID>")
	}

	playlistID, err := uuid.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID плейлиста: %v", err)
	}

	trackID, err := uuid.Parse(args[1])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID трека: %v", err)
	}

	err = app.playlistUseCase.AddTrackToPlaylist(playlistID, trackID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Трек %s успешно добавлен в плейлист %s", trackID, playlistID), nil
}

func (app *ConsoleApp) removeTrackFromPlaylist(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 2 {
		return "", fmt.Errorf("использование: remove-from-playlist <playlistID> <trackID>")
	}

	playlistID, err := uuid.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID плейлиста: %v", err)
	}

	trackID, err := uuid.Parse(args[1])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID трека: %v", err)
	}

	err = app.playlistUseCase.RemoveTrackFromPlaylist(playlistID, trackID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Трек %s успешно удален из плейлиста %s", trackID, playlistID), nil
}

func (app *ConsoleApp) getPlaylist(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 1 {
		return "", fmt.Errorf("использование: get-playlist <playlistID>")
	}

	playlistID, err := uuid.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID: %v", err)
	}

	playlistTracks, err := app.playlistUseCase.GetPlaylistTracks(playlistID)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Плейлист: %s\n\nТреки:\n", playlistID)

	for i, track := range playlistTracks {
		result += fmt.Sprintf("%d. %s (ID: %s, Исполнитель: %s)\n",
			i+1, track.Title, track.ID, track.ArtistName)
	}

	return result, nil
}

func (app *ConsoleApp) playTrack(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 1 {
		return "", fmt.Errorf("использование: play-track <trackID>")
	}

	trackID, err := uuid.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID: %v", err)
	}

	err = app.trackUseCase.PlayTrack(app.currentUser.ID, trackID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Трек %s воспроизводится", trackID), nil
}

func (app *ConsoleApp) getHistory() (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	history, err := app.historyUseCase.GetUserHistory(app.currentUser.ID)
	if err != nil {
		return "", err
	}

	result := "История прослушивания:\n"
	for i, item := range history {
		result += fmt.Sprintf("%d. Трек: %s, Дата: %s\n",
			i+1, item.TrackID, item.ListenedAt.Format("2006-01-02 15:04:05"))
	}

	return result, nil
}

func (app *ConsoleApp) deleteTrack(args []string) (string, error) {
	if app.currentUser == nil {
		return "", fmt.Errorf("необходимо войти в систему")
	}

	if len(args) < 1 {
		return "", fmt.Errorf("использование: delete-track <trackID>")
	}

	trackID, err := uuid.Parse(args[0])
	if err != nil {
		return "", fmt.Errorf("неверный формат ID: %v", err)
	}

	_, err = app.trackUseCase.GetTrackDetails(trackID)
	if err != nil {
		return "", fmt.Errorf("трек не найден: %v", err)
	}

	err = app.trackUseCase.DeleteTrack(trackID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Трек с ID %s успешно удален", trackID), nil
}
