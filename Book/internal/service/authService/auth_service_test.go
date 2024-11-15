package authservice

import (
	mock "Book/internal/service/authService/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Registration(t *testing.T) {
	type fields struct {
		repository AuthRepository
	}

	tests := []struct {
		name     string
		fields   fields
		email    string
		password string
		wantErr  bool
		setup    func(field *fields, ctrl *gomock.Controller)
	}{
		{
			name:     "test get success",
			fields:   fields{},
			email:    "test@yandex.ru",
			password: "12345",
			wantErr:  false,
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockAuthRepository(ctrl)                                        // инициализирую пакет заглушки
				rep.EXPECT().Registration(gomock.Any(), "test@yandex.ru", "12345").Return(nil) // создаю метод заглушку(мок)репозитория
				field.repository = rep                                                         // определяю репозиторий
			},
		},
		{
			name:     "test get error",
			fields:   fields{},
			email:    "test@yandex.ru",
			password: "12345",
			wantErr:  true,
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockAuthRepository(ctrl)
				rep.EXPECT().Registration(gomock.Any(), "test@yandex.ru", "12345").Return(errors.New("test error"))
				field.repository = rep
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &AuthServiceStruct{
				authRep: tt.fields.repository,
			}

			err := s.Registration(context.Background(), tt.email, tt.password)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	type fields struct {
		repository AuthRepository
	}

	tests := []struct {
		name     string
		fields   fields
		email    string
		password string
		want     int
		wantErr  bool
		setup    func(field *fields, ctrl *gomock.Controller)
	}{
		{
			name:     "test get success",
			fields:   fields{},
			email:    "test@yandex.ru",
			password: "12345",
			want:     3,
			wantErr:  false,
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockAuthRepository(ctrl)
				rep.EXPECT().Login(gomock.Any(), "test@yandex.ru", "12345").Return(3, nil)
				field.repository = rep
			},
		},
		{
			name:     "test get error",
			fields:   fields{},
			email:    "test@yandex.ru",
			password: "12345",
			want:     0,
			wantErr:  true,
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockAuthRepository(ctrl)
				rep.EXPECT().Login(gomock.Any(), "test@yandex.ru", "12345").Return(0, errors.New("test error"))
				field.repository = rep
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &AuthServiceStruct{
				authRep: tt.fields.repository,
			}

			token, err := s.Login(context.Background(), tt.email, tt.password)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Equal(t, tt.want, token)
			assert.Nil(t, err)
		})
	}
}

func TestAuthService_ChekAuth(t *testing.T) {
	type fields struct {
		repository AuthRepository
	}

	tests := []struct {
		name    string
		fields  fields
		args    string
		wantErr bool
		setup   func(field *fields, ctrl *gomock.Controller)
	}{
		{
			name:    "test get success",
			fields:  fields{},
			args:    "test@yandex.ru",
			wantErr: false,
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockAuthRepository(ctrl)
				rep.EXPECT().ChekAuth(gomock.Any(), "test@yandex.ru").Return(nil)
				field.repository = rep
			},
		},

		{
			name:    "test get error",
			fields:  fields{},
			args:    "test@yandex.ru",
			wantErr: true,
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockAuthRepository(ctrl)
				rep.EXPECT().ChekAuth(gomock.Any(), "test@yandex.ru").Return(errors.New("test error"))
				field.repository = rep
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &AuthServiceStruct{
				authRep: tt.fields.repository,
			}

			err := s.ChekAuth(context.Background(), tt.args)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
