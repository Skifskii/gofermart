package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	repo      UserAdder
	secretKey string
}

type UserAdder interface {
	AddUser(login, password string) error
}

type Claims struct {
	jwt.RegisteredClaims
	UserLogin string
}

func New(repo UserAdder, secretKey string) *AuthService {
	return &AuthService{repo: repo, secretKey: secretKey}
}

func (a *AuthService) CreateUser(login, password string) (jwt string, err error) {
	if err := a.repo.AddUser(login, password); err != nil {
		return "", err
	}

	return a.buildJWTString(login)
}

func (a *AuthService) buildJWTString(login string) (jwtString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{UserLogin: login})

	return token.SignedString([]byte(a.secretKey))
}
