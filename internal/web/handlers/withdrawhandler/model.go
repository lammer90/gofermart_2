package withdrawhandler

import "net/http"

type WithdrawRestApiProvider interface {
	Save(http.ResponseWriter, *http.Request)
	FindAll(http.ResponseWriter, *http.Request)
	FindBalance(http.ResponseWriter, *http.Request)
}
