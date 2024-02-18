package main

import (
	"context"
	"database/sql"
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
	"github.com/lammer90/gofermart/internal/web/app"
	"github.com/lammer90/gofermart/internal/web/handlers/authhandler"
	"github.com/lammer90/gofermart/internal/web/handlers/orderhandler"
	"github.com/lammer90/gofermart/internal/web/handlers/withdrawhandler"
	"github.com/lammer90/gofermart/internal/web/middleware/authfilter"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go accrualservice.New(orderSrv, config.AccrualAddress).Start(ctx)

	app.Start(config.ServAddress, authHdl, orderHdl, withHdl, authMdl)
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
