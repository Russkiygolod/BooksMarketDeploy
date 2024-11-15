package authservice

import (
	"context"
)

//go:generate mockgen -source=authService.go -destination=mocks/auth_mock.go -package=mock
//export PATH=$PATH:$(go env GOPATH)/bin
type AuthRepository interface {
	Registration(ctx context.Context, email, password string) error // регистрация пользователя
	Login(ctx context.Context, email, password string) (int, error) // проверка логина и пароля для генерации токена
	ChekAuth(ctx context.Context, email string) error               // проверка email из полезной нагрузки токена
}

type AuthServiceStruct struct {
	authRep AuthRepository
}

func NewAuthService(a AuthRepository) AuthServiceStruct {
	var user AuthServiceStruct
	user.authRep = a
	return user
}

func (a AuthServiceStruct) Registration(ctx context.Context, email, password string) error {
	return a.authRep.Registration(ctx, email, password)
}

func (a AuthServiceStruct) Login(ctx context.Context, email, password string) (int, error) {
	return a.authRep.Login(ctx, email, password)
}

func (a AuthServiceStruct) ChekAuth(ctx context.Context, email string) error {
	return a.authRep.ChekAuth(ctx, email)
}
