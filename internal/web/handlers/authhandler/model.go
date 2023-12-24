package authhandler

import "net/http"

type AuthenticationRestAPIProvider interface {
	Register(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}
