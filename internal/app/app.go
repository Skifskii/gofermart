package app

import (
	"fmt"
	"gophermart/internal/handler/api/user/register"
	"gophermart/internal/repository/inmem"
	"gophermart/internal/service/auth"
	"net/http"

	"github.com/go-chi/chi"
)

func Run() error {
	// Конфиг
	// TODO:

	// Логгер
	// TODO:

	// Репозиторий
	repo := inmem.New()

	// Сервис авторизации
	authService := auth.New(repo, "supersecret") // TODO: вынести секретный ключ в конфиг

	// HTTP сервер
	router := chi.NewRouter()

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", register.New(authService))
	})

	address := "localhost:8080"
	fmt.Printf("Starting server at %s\n", address)
	http.ListenAndServe(address, router)

	return nil
}
