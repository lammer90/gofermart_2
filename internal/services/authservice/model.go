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

var ErrNotAuthorized = errors.New("user not authorized")

var ErrUserAlreadyExist = errors.New("user already exist")

var ErrUserDidntFind = errors.New("user didn't find")
