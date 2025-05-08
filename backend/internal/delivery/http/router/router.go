package router

import (
	"music-service/internal/delivery/http/handlers"
	"music-service/internal/delivery/http/middleware"
	"music-service/internal/usecases/interfaces"

	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

func NewRouter(
	userUseCase interfaces.UserUseCase,
	trackUseCase interfaces.TrackUseCase,
	albumUseCase interfaces.AlbumUseCase,
	genreUseCase interfaces.GenreUseCase,
	allowedTypes []string,
	maxFileSizeMB int,
) *Router {
	r := mux.NewRouter()
	router := &Router{
		router: r,
	}

	r.Use(middleware.CORS)

	userHandler := handlers.NewUserHandler(userUseCase)
	trackHandler := handlers.NewTrackHandler(trackUseCase, allowedTypes, maxFileSizeMB)
	albumHandler := handlers.NewAlbumHandler(albumUseCase)
	genreHandler := handlers.NewGenreHandler(genreUseCase)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/users", userHandler.RegisterUser).Methods("POST", "OPTIONS")
	v1.HandleFunc("/users/auth", userHandler.AuthenticateUser).Methods("POST", "OPTIONS")
	v1.HandleFunc("/users/{id}", userHandler.GetUserProfile).Methods("GET", "OPTIONS")
	v1.HandleFunc("/users/{id}/permissions", userHandler.UpdateUserPermissions).Methods("PATCH", "OPTIONS")
	v1.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE", "OPTIONS")
	v1.HandleFunc("/users/logout", userHandler.LogoutUser).Methods("POST", "OPTIONS")

	v1.HandleFunc("/tracks", trackHandler.SearchTracks).Methods("GET", "OPTIONS")
	v1.HandleFunc("/tracks", trackHandler.UploadTrack).Methods("POST", "OPTIONS")
	v1.HandleFunc("/tracks/{id}", trackHandler.GetTrackDetails).Methods("GET", "OPTIONS")
	v1.HandleFunc("/tracks/{id}/stream", trackHandler.ServeTrackFile).Methods("GET", "OPTIONS")
	v1.HandleFunc("/tracks/{id}", trackHandler.DeleteTrack).Methods("DELETE", "OPTIONS")

	v1.HandleFunc("/albums", albumHandler.ListAllAlbums).Methods("GET", "OPTIONS")
	v1.HandleFunc("/albums", albumHandler.CreateAlbum).Methods("POST", "OPTIONS")
	v1.HandleFunc("/albums/{id}", albumHandler.GetAlbumDetails).Methods("GET", "OPTIONS")
	v1.HandleFunc("/albums/{id}", albumHandler.UpdateAlbum).Methods("PUT", "OPTIONS")
	v1.HandleFunc("/albums/{id}", albumHandler.DeleteAlbum).Methods("DELETE", "OPTIONS")
	v1.HandleFunc("/albums/{id}/tracks", albumHandler.AddTrackToAlbum).Methods("POST", "OPTIONS")
	v1.HandleFunc("/albums/{id}/tracks/{track_id}", albumHandler.RemoveTrackFromAlbum).Methods("DELETE", "OPTIONS")

	v1.HandleFunc("/genres", genreHandler.ListAllGenres).Methods("GET", "OPTIONS")
	v1.HandleFunc("/genres", genreHandler.CreateGenre).Methods("POST", "OPTIONS")
	v1.HandleFunc("/genres/tracks/{id}", genreHandler.GetGenresByTrack).Methods("GET", "OPTIONS")
	v1.HandleFunc("/genres/tracks/{id}/genres", genreHandler.AssignGenreToTrack).Methods("POST", "OPTIONS")
	v1.HandleFunc("/genres/tracks/{trackId}/genres/{genreId}", genreHandler.RemoveGenreFromTrack).Methods("DELETE", "OPTIONS")

	return router
}

func (r *Router) GetRouter() *mux.Router {
	return r.router
}
