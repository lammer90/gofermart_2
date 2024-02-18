package authhandler

import (
	"encoding/json"
	"net/http"

	"github.com/lammer90/gofermart/internal/dto/auth"
	"github.com/lammer90/gofermart/internal/services/authservice"
)

type authenticationHandler struct {
	authenticationService authservice.AuthenticationService
}

func New(authenticationService authservice.AuthenticationService) AuthenticationRestAPIProvider {
	return authenticationHandler{
		authenticationService: authenticationService,
	}
}

func (a authenticationHandler) Register(res http.ResponseWriter, req *http.Request) {
	var request auth.AuthRequest
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&request)
	if err != nil || request.Login == "" || request.Password == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.authenticationService.ToRegisterUser(request.Login, request.Password)
	if err != nil && err == authservice.ErrUserAlreadyExist {
		res.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  "Authorization",
		Value: token,
		Path:  "/",
	})
	res.WriteHeader(http.StatusOK)
}

func (a authenticationHandler) Login(res http.ResponseWriter, req *http.Request) {
	var request auth.AuthRequest
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&request)
	if err != nil || request.Login == "" || request.Password == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.authenticationService.ToLoginUser(request.Login, request.Password)
	if err != nil && (err == authservice.ErrUserDidntFind || err == authservice.ErrNotAuthorized) {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name:  "Authorization",
		Value: token,
		Path:  "/",
	})
	res.WriteHeader(http.StatusOK)
}
