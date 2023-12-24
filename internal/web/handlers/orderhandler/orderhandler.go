package orderhandler

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/lammer90/gofermart/internal/services/orderservice"
	"io"
	"net/http"
)

type orderHandler struct {
	orderService orderservice.OrderService
	cookieStore  *sessions.CookieStore
}

func New(orderService orderservice.OrderService, cookieStore *sessions.CookieStore) OrderRestAPIProvider {
	return orderHandler{
		orderService: orderService,
		cookieStore:  cookieStore,
	}
}

func (o orderHandler) Save(res http.ResponseWriter, req *http.Request) {
	login := getLogin(req, o.cookieStore)
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	number := string(body[:])

	err = o.orderService.Save(number, login)
	if err != nil && err == orderservice.NotValidLuhnSum {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err != nil && err == orderservice.OrderNumberAlreadyExistThisUser {
		res.WriteHeader(http.StatusOK)
		return
	}
	if err != nil && err == orderservice.OrderNumberAlreadyExistAnotherUser {
		res.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

func (o orderHandler) FindAll(res http.ResponseWriter, req *http.Request) {
	login := getLogin(req, o.cookieStore)

	orders, err := o.orderService.FindAll(login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(res)
	if err := enc.Encode(orders); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getLogin(req *http.Request, cookieStore *sessions.CookieStore) string {
	session, _ := cookieStore.Get(req, "Authorization")
	return session.Values["login"].(string)
}
