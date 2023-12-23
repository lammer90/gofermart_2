package authhandler

import "net/http"

type AuthenticationRestApiProvider interface {
	Register(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}
