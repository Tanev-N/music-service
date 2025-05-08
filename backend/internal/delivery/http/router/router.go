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
) *Router {
	r := mux.NewRouter()
	router := &Router{
		router: r,
	}

	r.Use(middleware.CORS)

	userHandler := handlers.NewUserHandler(userUseCase)

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/users", userHandler.RegisterUser).Methods("POST", "OPTIONS")
	v1.HandleFunc("/users/auth", userHandler.AuthenticateUser).Methods("POST", "OPTIONS")
	v1.HandleFunc("/users/{id}", userHandler.GetUserProfile).Methods("GET", "OPTIONS")
	v1.HandleFunc("/users/{id}/permissions", userHandler.UpdateUserPermissions).Methods("PATCH", "OPTIONS")
	v1.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE", "OPTIONS")
	v1.HandleFunc("/users/logout", userHandler.LogoutUser).Methods("POST", "OPTIONS")

	return router
}

func (r *Router) GetRouter() *mux.Router {
	return r.router
}
