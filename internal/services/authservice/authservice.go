package authservice

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lammer90/gofermart/internal/repository/balance"
	"github.com/lammer90/gofermart/internal/repository/userstorage"
)

type authenticationServiceImpl struct {
	userRepository    userstorage.UserRepository
	balanceRepository balance.BalanceRepository
	privateKey        string
}

func New(userRepository userstorage.UserRepository, balanceRepository balance.BalanceRepository, privateKey string) AuthenticationService {
	return &authenticationServiceImpl{userRepository: userRepository, balanceRepository: balanceRepository, privateKey: privateKey}
}

func (a *authenticationServiceImpl) CheckAuthentication(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(a.privateKey), nil
		})
	if err != nil {
		return "", ErrNotAuthorized
	}
	if !token.Valid || claims.Login == "" {
		return "", ErrNotAuthorized
	}
	return claims.Login, nil
}

func (a *authenticationServiceImpl) ToRegisterUser(login, password string) (token string, err error) {
	existHash, err := a.userRepository.Find(login)
	if err != nil {
		return "", err
	}
	if existHash != "" {
		return "", ErrUserAlreadyExist
	}

	newHash := buildHash(login, password)
	err = a.balanceRepository.CreateBalance(login)
	if err != nil {
		return "", err
	}
	err = a.userRepository.Save(login, newHash)
	if err != nil {
		return "", err
	}

	token, err = buildJWTString(login, a.privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *authenticationServiceImpl) ToLoginUser(login, password string) (token string, err error) {
	existHash, err := a.userRepository.Find(login)
	if err != nil {
		return "", err
	}
	if existHash == "" {
		return "", ErrUserDidntFind
	}

	sentHash := buildHash(login, password)
	if existHash != sentHash {
		return "", ErrNotAuthorized
	}

	token, err = buildJWTString(login, a.privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func buildHash(login string, password string) string {
	src := []byte(login + ":" + password)
	newHashByte := sha256.Sum256(src)
	return base64.StdEncoding.EncodeToString(newHashByte[:])
}

func buildJWTString(login, privateKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		Login:            login,
	})

	tokenString, err := token.SignedString([]byte(privateKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
