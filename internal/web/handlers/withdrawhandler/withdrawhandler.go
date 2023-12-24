package withdrawhandler

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/lammer90/gofermart/internal/dto/withdraw"
	"github.com/lammer90/gofermart/internal/services/withdrawservice"
	"net/http"
)

type withdrawHandler struct {
	withdrawService withdrawservice.WithdrawService
	cookieStore     *sessions.CookieStore
}

func New(withdrawService withdrawservice.WithdrawService, cookieStore *sessions.CookieStore) WithdrawRestAPIProvider {
	return withdrawHandler{
		withdrawService: withdrawService,
		cookieStore:     cookieStore,
	}
}

func (w withdrawHandler) Save(res http.ResponseWriter, req *http.Request) {
	login := getLogin(req, w.cookieStore)

	var request withdraw.WithdrawRequest
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&request)
	if err != nil || request.Order == "" || request.Sum == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = w.withdrawService.Save(request.Order, login, request.Sum)
	if err != nil && err == withdrawservice.ErrNotValidLuhnSum {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err != nil && err == withdrawservice.ErrNotEnoughMoney {
		res.WriteHeader(http.StatusPaymentRequired)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (w withdrawHandler) FindAll(res http.ResponseWriter, req *http.Request) {
	login := getLogin(req, w.cookieStore)

	withdraws, err := w.withdrawService.FindAll(login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(withdraws) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(res)
	if err := enc.Encode(withdraws); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (w withdrawHandler) FindBalance(res http.ResponseWriter, req *http.Request) {
	login := getLogin(req, w.cookieStore)

	bal, err := w.withdrawService.FindBalance(login)
	if err != nil || bal == nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(res)
	if err := enc.Encode(bal); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getLogin(req *http.Request, cookieStore *sessions.CookieStore) string {
	session, _ := cookieStore.Get(req, "Authorization")
	return session.Values["login"].(string)
}
