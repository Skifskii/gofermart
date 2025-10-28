package app

import (
	"fmt"
	"gophermart/internal/handler/api/user/balance"
	"gophermart/internal/handler/api/user/login"
	"gophermart/internal/handler/api/user/orders"

	"gophermart/internal/handler/api/user/register"
	"gophermart/internal/middleware/authmw"
	"gophermart/internal/repository/inmem"
	"gophermart/internal/repository/postgres"
	"gophermart/internal/service/auth"
	bm "gophermart/internal/service/balance"
	om "gophermart/internal/service/orders"
	"gophermart/internal/service/orders/loyalsys"
	"net/http"

	"github.com/go-chi/chi"
)

func Run() error {
	// Конфиг
	// TODO:

	// Логгер
	// TODO:

	// Репозиторий
	dsn := "postgresql://gouser:gopassword@localhost:5432/gophermartdb?sslmode=disable"
	repo, _ := postgres.New(dsn)
	inmemRepo := inmem.New()

	// Сервисы:
	// - сервис авторизации
	authService := auth.New(repo, "supersecret") // TODO: вынести секретный ключ в конфиг

	// - сервис управления балансом
	balanceManager := bm.New(inmemRepo)

	// - сервис управления заказами
	loyaltySystem := loyalsys.New()
	ordersManager := om.New(inmemRepo, loyaltySystem)

	// HTTP сервер
	router := chi.NewRouter()

	router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", register.New(authService))
		r.Post("/login", login.New(authService))

		r.Group(func(r chi.Router) {
			r.Use(authmw.AuthMiddleware(authService))

			r.Get("/balance", balance.New(balanceManager))
			r.Post("/orders", orders.NewPost(ordersManager))
			r.Get("/orders", orders.NewGet(ordersManager))
		})
	})

	address := "localhost:8080"
	fmt.Printf("Starting server at %s\n", address)
	http.ListenAndServe(address, router)

	return nil
}
