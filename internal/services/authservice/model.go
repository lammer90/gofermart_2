package authservice

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

type AuthenticationService interface {
	CheckAuthentication(token string) (string, error)
	ToRegisterUser(login, password string) (string, error)
	ToLoginUser(login, password string) (string, error)
}

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

var NotAuthorized = errors.New("user not authorized")

var UserAlreadyExist = errors.New("user already exist")

var UserDidntFind = errors.New("user didn't find")
