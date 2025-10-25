package app

import (
	"fmt"
	"gophermart/internal/handler/api/user/balance"
	"gophermart/internal/handler/api/user/login"
	"gophermart/internal/handler/api/user/register"
	"gophermart/internal/repository/inmem"
	"gophermart/internal/service/auth"
	bm "gophermart/internal/service/balance"
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

	// Сервисы:
	// - сервис авторизации
	authService := auth.New(repo, "supersecret") // TODO: вынести секретный ключ в конфиг

	// - сервис управления балансом
	balanceManager := bm.New(repo)

	// HTTP сервер
	router := chi.NewRouter()

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", register.New(authService))
		r.Post("/login", login.New(authService))

		r.Get("/balance", balance.New(balanceManager))
	})

	address := "localhost:8080"
	fmt.Printf("Starting server at %s\n", address)
	http.ListenAndServe(address, router)

	return nil
}
