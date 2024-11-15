package bookservice

import (
	mock "Book/internal/service/bookService/mocks"
	booksauthors "Book/pkg/models/booksAuthors"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_CreateBook(t *testing.T) {
	type fields struct {
		repository BookRepository
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
				rep := mock.NewMockBookRepository(ctrl) // инициализирую пакет заглушки
				rep.EXPECT().CreateBook(gomock.Any(), booksauthors.Book{
					Book_title:  "Whale",
					Author_name: "Aron",
					Price:       300,
				}).Return(nil) // создаю метод заглушку(мок)репозитория
				field.repository = rep // определяю репозиторий
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
				rep := mock.NewMockBookRepository(ctrl)
				rep.EXPECT().CreateBook(gomock.Any(), booksauthors.Book{
					Book_title:  "Whale",
					Author_name: "Aron",
					Price:       300,
				}).Return(errors.New("error"))
				field.repository = rep
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &bookServiceStruct{
				bookRep: tt.fields.repository,
			}

			err := s.CreateBook(context.Background(), tt.args)
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
		repository      BookRepository
		redisRepository bookRedisRepository
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
			name:   "test get success empty book_tittle search id", // test case
			fields: fields{},
			args: booksauthors.Book{
				Book_id: 1,
			},
			wantErr: false,
			want: []booksauthors.Book{
				{
					Book_id:     1,
					Book_title:  "pipka",
					Author_name: "Ignat",
				},
			},
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockBookRepository(ctrl)
				rep.EXPECT().ReadBook(gomock.Any(), booksauthors.Book{
					Book_id: 1,
				}).Return([]booksauthors.Book{
					{
						Book_id:     1,
						Book_title:  "pipka",
						Author_name: "Ignat",
					}}, nil)
				field.repository = rep
			},
		},
		{
			name:   "test get success empty book_tittle search Author_name ", // test case
			fields: fields{},
			args: booksauthors.Book{
				Author_name: "Petr",
			},
			wantErr: false,
			want: []booksauthors.Book{
				{
					Book_id:     2,
					Book_title:  "Ribka",
					Author_name: "Petr",
				},
				{
					Book_id:     5,
					Book_title:  "Utka",
					Author_name: "Petr",
				},
			},
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockBookRepository(ctrl)
				rep.EXPECT().ReadBook(gomock.Any(), booksauthors.Book{
					Author_name: "Petr",
				}).Return([]booksauthors.Book{
					{
						Book_id:     2,
						Book_title:  "Ribka",
						Author_name: "Petr",
					},
					{
						Book_id:     5,
						Book_title:  "Utka",
						Author_name: "Petr",
					}}, nil)
				field.repository = rep
			},
		},
		{
			name:   "test get success with redis", // test case
			fields: fields{},
			args: booksauthors.Book{
				Book_title: "pipka",
			},
			wantErr: false,
			want: []booksauthors.Book{
				{
					Book_title:  "pipka",
					Author_name: "Jonny",
					Price:       100,
				},
			},
			setup: func(field *fields, ctrl *gomock.Controller) {
				red := mock.NewMockbookRedisRepository(ctrl)
				red.EXPECT().GetBook("pipka").Return(&booksauthors.Book{
					Book_title:  "pipka",
					Author_name: "Jonny",
					Price:       100}, nil)
				field.redisRepository = red
			},
		},
		{
			name:   "test get success with redis empty", // test case
			fields: fields{},
			args: booksauthors.Book{
				Book_title: "pipka",
			},
			wantErr: false,
			want: []booksauthors.Book{
				{
					Book_title:  "pipka",
					Author_name: "Jonny",
					Price:       100,
				},
			},
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockBookRepository(ctrl)
				red := mock.NewMockbookRedisRepository(ctrl)

				red.EXPECT().GetBook("pipka").Return(nil, nil)

				rep.EXPECT().ReadBook(gomock.Any(), booksauthors.Book{
					Book_title: "pipka",
				}).Return([]booksauthors.Book{
					{
						Book_title:  "pipka",
						Author_name: "Jonny",
						Price:       100,
					}}, nil)

				red.EXPECT().SetBook("pipka", booksauthors.Book{
					Book_title:  "pipka",
					Author_name: "Jonny",
					Price:       100,
				}).Return(nil)
				field.repository = rep
				field.redisRepository = red
			},
		},
		{
			name:   "test get error", // test case
			fields: fields{},
			args: booksauthors.Book{
				Book_id: 1,
			},
			wantErr: true,
			want:    nil,
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockBookRepository(ctrl)
				rep.EXPECT().ReadBook(gomock.Any(), booksauthors.Book{
					Book_id: 1,
				}).Return(nil, errors.New("test error"))
				field.repository = rep
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &bookServiceStruct{
				bookRep:      tt.fields.repository,
				bookRedisRep: tt.fields.redisRepository,
			}

			books, err := s.ReadBook(context.Background(), tt.args)
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
		repository      BookRepository
		redisRepository bookRedisRepository
	}
	tests := []struct {
		name    string
		fields  fields
		argsOne int
		argsTwo booksauthors.Book
		wantErr bool
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
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockBookRepository(ctrl)
				red := mock.NewMockbookRedisRepository(ctrl)

				rep.EXPECT().UpdateBook(gomock.Any(), 1, booksauthors.Book{
					Authors_old_name: "Ignat",
					Author_name:      "Petr",
				}).Return(nil)

				rep.EXPECT().ReadBook(gomock.Any(), booksauthors.Book{
					Book_id: 1,
				}).Return([]booksauthors.Book{
					{
						Book_title:  "Olega",
						Author_name: "Petr",
						Price:       300,
					}}, nil)

				red.EXPECT().SetBook("Olega", booksauthors.Book{
					Book_title:  "Olega",
					Author_name: "Petr",
					Price:       300,
				}).Return(nil)
				field.repository = rep
				field.redisRepository = red
			},
		},
		{
			name:    "test update error", // test case
			fields:  fields{},
			argsOne: 1,
			argsTwo: booksauthors.Book{
				Authors_old_name: "Ignat",
				Author_name:      "Petr",
			},
			wantErr: true,
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockBookRepository(ctrl)
				rep.EXPECT().UpdateBook(gomock.Any(), 1, booksauthors.Book{
					Authors_old_name: "Ignat",
					Author_name:      "Petr",
				}).Return(errors.New("test error"))
				field.repository = rep
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &bookServiceStruct{
				bookRep:      tt.fields.repository,
				bookRedisRep: tt.fields.redisRepository,
			}
			err := s.UpdateBook(context.Background(), tt.argsOne, tt.argsTwo)
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
		repository      BookRepository
		redisRepository bookRedisRepository
	}
	tests := []struct {
		name    string
		fields  fields
		args    int
		wantErr bool
		setup   func(field *fields, ctrl *gomock.Controller)
	}{
		{
			name:    "test delete success", // test case
			fields:  fields{},
			args:    1,
			wantErr: false,
			setup: func(field *fields, ctrl *gomock.Controller) {
				rep := mock.NewMockBookRepository(ctrl)
				red := mock.NewMockbookRedisRepository(ctrl)

				rep.EXPECT().ReadBook(gomock.Any(), booksauthors.Book{
					Book_id: 1,
				}).Return([]booksauthors.Book{
					{
						Book_title:  "Olega",
						Author_name: "Petr",
						Price:       300,
					}}, nil)

				rep.EXPECT().DeleteBook(gomock.Any(), 1).Return(nil)

				red.EXPECT().DelBook("Olega").Return(nil)

				field.repository = rep
				field.redisRepository = red
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctr := gomock.NewController(t)

			tt.setup(&tt.fields, ctr)

			s := &bookServiceStruct{
				bookRep:      tt.fields.repository,
				bookRedisRep: tt.fields.redisRepository,
			}

			err := s.DeleteBook(context.Background(), tt.args)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

		})
	}
}
