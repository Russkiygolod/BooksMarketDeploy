package app

import (
	mock "Book/internal/app/mocks"
	booksauthors "Book/pkg/models/booksAuthors"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestServer_Registration(t *testing.T) {
	type fields struct {
		aService AuthService
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
				aserv := mock.NewMockAuthService(ctrl)                                           // инициализирую пакет заглушки
				aserv.EXPECT().Registration(gomock.Any(), "test@yandex.ru", "12345").Return(nil) // создаю метод заглушку(мок)репозитория
				field.aService = aserv                                                           // определяю репозиторий
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &GRPCserver{
				authService: tt.fields.aService,
			}

			err := s.authService.Registration(context.Background(), tt.email, tt.password)
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
		aService AuthService
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
				aserv := mock.NewMockAuthService(ctrl)                                       // инициализирую пакет заглушки
				aserv.EXPECT().Login(gomock.Any(), "test@yandex.ru", "12345").Return(3, nil) // создаю метод заглушку(мок)репозитория
				field.aService = aserv                                                       // определяю репозиторий
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &GRPCserver{
				authService: tt.fields.aService,
			}

			token, err := s.authService.Login(context.Background(), tt.email, tt.password)
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
		aService AuthService
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
				aserv := mock.NewMockAuthService(ctrl)                              // инициализирую пакет заглушки
				aserv.EXPECT().ChekAuth(gomock.Any(), "test@yandex.ru").Return(nil) // создаю метод заглушку(мок)репозитория
				field.aService = aserv                                              // определяю репозиторий
			},
		},

		{
			name:    "test get error",
			fields:  fields{},
			args:    "test@yandex.ru",
			wantErr: true,
			setup: func(field *fields, ctrl *gomock.Controller) {
				aserv := mock.NewMockAuthService(ctrl)                                                   // инициализирую пакет заглушки
				aserv.EXPECT().ChekAuth(gomock.Any(), "test@yandex.ru").Return(errors.New("test error")) // создаю метод заглушку(мок)репозитория
				field.aService = aserv                                                                   // определяю репозиторий
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &GRPCserver{
				authService: tt.fields.aService,
			}

			err := s.authService.ChekAuth(context.Background(), tt.args)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestAuthService_CreateBook(t *testing.T) {
	type fields struct {
		bService BookService
	}
	tests := []struct {
		name    string
		fields  fields
		args    booksauthors.Book
		wantErr bool
		setup   func(field *fields, ctrl *gomock.Controller)
	}{
		{
			name:   "test create success",
			fields: fields{},
			args: booksauthors.Book{
				Book_title:  "Whale",
				Author_name: "Aron",
				Price:       300,
			},
			wantErr: false,
			setup: func(field *fields, ctrl *gomock.Controller) {
				bserv := mock.NewMockBookService(ctrl) // инициализирую пакет заглушки
				bserv.EXPECT().CreateBook(gomock.Any(), booksauthors.Book{
					Book_title:  "Whale",
					Author_name: "Aron",
					Price:       300,
				}).Return(nil) // создаю метод заглушку(мок)репозитория
				field.bService = bserv // определяю репозиторий
			},
		},
		{
			name:   "test get error",
			fields: fields{},
			args: booksauthors.Book{
				Book_title:  "Whale",
				Author_name: "Aron",
				Price:       300,
			},
			wantErr: true,
			setup: func(field *fields, ctrl *gomock.Controller) {
				bserv := mock.NewMockBookService(ctrl) // инициализирую пакет заглушки
				bserv.EXPECT().CreateBook(gomock.Any(), booksauthors.Book{
					Book_title:  "Whale",
					Author_name: "Aron",
					Price:       300,
				}).Return(errors.New("error")) // создаю метод заглушку(мок)репозитория
				field.bService = bserv // определяю репозиторий
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &GRPCserver{
				bookService: tt.fields.bService,
			}

			err := s.bookService.CreateBook(context.Background(), tt.args)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestAuthService_ReadBook(t *testing.T) {
	type fields struct {
		bService BookService
	}
	tests := []struct {
		name    string
		fields  fields
		args    booksauthors.Book
		wantErr bool
		want    []booksauthors.Book
		setup   func(field *fields, ctrl *gomock.Controller)
	}{
		{
			name:   "test get success empty book_tittle", // test case
			fields: fields{},
			args: booksauthors.Book{
				Book_id: 1,
			},
			wantErr: false,
			want: []booksauthors.Book{
				{
					Book_id:    1,
					Book_title: "pipka",
				},
			},
			setup: func(field *fields, ctrl *gomock.Controller) {
				bserv := mock.NewMockBookService(ctrl)
				bserv.EXPECT().ReadBook(gomock.Any(), booksauthors.Book{
					Book_id: 1,
				}).Return([]booksauthors.Book{
					{
						Book_id:    1,
						Book_title: "pipka",
					}}, nil)
				field.bService = bserv
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &GRPCserver{
				bookService: tt.fields.bService,
			}

			books, err := s.bookService.ReadBook(context.Background(), tt.args)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Len(t, tt.want, len(books))
			for i := range tt.want {
				assert.Equal(t, tt.want[i], books[i])
			}
		})
	}
}

func TestAuthService_UpdateBook(t *testing.T) {
	type fields struct {
		bService BookService
	}
	tests := []struct {
		name    string
		fields  fields
		argsOne int
		argsTwo booksauthors.Book
		wantErr bool
		want    booksauthors.Book
		setup   func(field *fields, ctrl *gomock.Controller)
	}{
		{
			name:    "test update success", // test case
			fields:  fields{},
			argsOne: 1,
			argsTwo: booksauthors.Book{
				Authors_old_name: "Ignat",
				Author_name:      "Petr",
			},
			wantErr: false,
			want: booksauthors.Book{
				Book_title:  "Olega",
				Author_name: "Petr",
				Price:       300,
			},

			setup: func(field *fields, ctrl *gomock.Controller) {
				bserv := mock.NewMockBookService(ctrl)
				bserv.EXPECT().UpdateBook(gomock.Any(), 1, booksauthors.Book{
					Authors_old_name: "Ignat",
					Author_name:      "Petr",
				}).Return(nil)
				field.bService = bserv
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)
			tt.setup(&tt.fields, ctr)
			s := &GRPCserver{
				bookService: tt.fields.bService,
			}
			err := s.bookService.UpdateBook(context.Background(), tt.argsOne, tt.argsTwo)
			// проверить функцию readbook
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

		})
	}
}

func TestAuthService_DeleteBook(t *testing.T) {
	type fields struct {
		bService BookService
	}
	tests := []struct {
		name    string
		fields  fields
		args    int
		wantErr bool
		want    booksauthors.Book
		setup   func(field *fields, ctrl *gomock.Controller)
	}{
		{
			name:    "test delete success", // test case
			fields:  fields{},
			args:    1,
			wantErr: false,
			setup: func(field *fields, ctrl *gomock.Controller) {
				bserv := mock.NewMockBookService(ctrl)
				bserv.EXPECT().DeleteBook(gomock.Any(), 1).Return(nil)
				field.bService = bserv
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &GRPCserver{
				bookService: tt.fields.bService,
			}

			err := s.bookService.DeleteBook(context.Background(), tt.args)
			// проверить функцию readbook
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

		})
	}
}
