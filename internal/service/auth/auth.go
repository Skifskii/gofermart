package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	repo      Repository
	secretKey string
}

type Repository interface {
	AddUser(login, password string) error
	AuthenticateUser(login, password string) error
}

type Claims struct {
	jwt.RegisteredClaims
	UserLogin string
}

func New(repo Repository, secretKey string) *AuthService {
	return &AuthService{repo: repo, secretKey: secretKey}
}

func (a *AuthService) CreateUser(login, password string) (jwt string, err error) {
	if err := a.repo.AddUser(login, password); err != nil {
		return "", err
	}

	return a.buildJWTString(login)
}

func (a *AuthService) AuthenticateUser(login, password string) (jwt string, err error) {
	if err := a.repo.AuthenticateUser(login, password); err != nil {
		return "", err
	}

	return a.buildJWTString(login)
}

func (a *AuthService) buildJWTString(login string) (jwtString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{UserLogin: login})

	return token.SignedString([]byte(a.secretKey))
}
