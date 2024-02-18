package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lammer90/gofermart/internal/web/handlers/authhandler"
	"github.com/lammer90/gofermart/internal/web/handlers/orderhandler"
	"github.com/lammer90/gofermart/internal/web/handlers/withdrawhandler"
)

func Start(
	servAddress string,
	authProvider authhandler.AuthenticationRestAPIProvider,
	orderProvider orderhandler.OrderRestAPIProvider,
	withdrawProvider withdrawhandler.WithdrawRestAPIProvider,
	middlewares ...func(next http.Handler) http.Handler) {

	router := chi.NewRouter()
	for _, f := range middlewares {
		router.Use(f)
	}
	router.Post("/api/user/register", authProvider.Register)
	router.Post("/api/user/login", authProvider.Login)
	router.Post("/api/user/orders", orderProvider.Save)
	router.Get("/api/user/orders", orderProvider.FindAll)
	router.Get("/api/user/balance", withdrawProvider.FindBalance)
	router.Post("/api/user/balance/withdraw", withdrawProvider.Save)
	router.Get("/api/user/withdrawals", withdrawProvider.FindAll)

	http.ListenAndServe(servAddress, router)
}
