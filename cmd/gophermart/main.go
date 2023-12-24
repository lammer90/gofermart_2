package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lammer90/gofermart/internal/config"
	"github.com/lammer90/gofermart/internal/logger"
	"github.com/lammer90/gofermart/internal/repository/balance"
	"github.com/lammer90/gofermart/internal/repository/orderstorage"
	"github.com/lammer90/gofermart/internal/repository/userstorage"
	"github.com/lammer90/gofermart/internal/repository/withdrawstorage"
	"github.com/lammer90/gofermart/internal/services/accrualservice"
	"github.com/lammer90/gofermart/internal/services/authservice"
	"github.com/lammer90/gofermart/internal/services/orderservice"
	"github.com/lammer90/gofermart/internal/services/withdrawservice"
	"github.com/lammer90/gofermart/internal/web/handlers/authhandler"
	"github.com/lammer90/gofermart/internal/web/handlers/orderhandler"
	"github.com/lammer90/gofermart/internal/web/handlers/withdrawhandler"
	"github.com/lammer90/gofermart/internal/web/middleware/authfilter"
	"net/http"
)

func main() {
	config.InitConfig()
	logger.InitLogger("info")

	db := InitDB("pgx", config.DataSource)
	defer db.Close()

	cookieStore := buildSession()

	balRep := balance.New(db)

	authSrv := authservice.New(userstorage.New(db), balRep, config.PrivateKey)
	authMdl := authfilter.New(authSrv, cookieStore)
	authHdl := authhandler.New(authSrv)

	orderSrv := orderservice.New(orderstorage.New(db), balRep, db)
	orderHdl := orderhandler.New(orderSrv, cookieStore)

	withSrv := withdrawservice.New(withdrawstorage.New(db), balRep, db)
	withHdl := withdrawhandler.New(withSrv, cookieStore)

	go accrualservice.New(orderSrv, config.AccrualAddress).Start()

	http.ListenAndServe(config.ServAddress, shortenerRouter(authHdl, orderHdl, withHdl, authMdl))
}

func shortenerRouter(
	authProvider authhandler.AuthenticationRestAPIProvider,
	orderProvider orderhandler.OrderRestAPIProvider,
	withdrawProvider withdrawhandler.WithdrawRestAPIProvider,
	middlewares ...func(next http.Handler) http.Handler) chi.Router {

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
	return router
}

func InitDB(driverName, dataSource string) *sql.DB {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		panic(err)
	}
	return db
}

func buildSession() *sessions.CookieStore {
	key := []byte("abc123")
	return sessions.NewCookieStore(key)
}
