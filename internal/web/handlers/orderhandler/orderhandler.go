package orderhandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/lammer90/gofermart/internal/services/orderservice"
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
	if err != nil && err == orderservice.ErrNotValidLuhnSum {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err != nil && err == orderservice.ErrOrderNumberAlreadyExistThisUser {
		res.WriteHeader(http.StatusOK)
		return
	}
	if err != nil && err == orderservice.ErrOrderNumberAlreadyExistAnotherUser {
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
