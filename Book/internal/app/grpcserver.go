package app

import (
	"Book/internal/books"
	authservice "Book/internal/service/authService"
	bookservice "Book/internal/service/bookService"
	booksauthors "Book/pkg/models/booksAuthors"
	"context"
)

//go:generate mockgen -source=grpcserver.go -destination=mocks/books_mock.go -package=mock

type AuthService interface {
	Registration(ctx context.Context, email, password string) error // регистрация пользователя
	Login(ctx context.Context, email, password string) (int, error) // проверка логина и пароля для генерации токена
	ChekAuth(ctx context.Context, email string) error               // проверка email из полезной нагрузки токена
}

type BookService interface {
	CreateBook(context.Context, booksauthors.Book) error                      // создать книгу
	ReadBook(context.Context, booksauthors.Book) ([]booksauthors.Book, error) // прочитать книгу
	UpdateBook(context.Context, int, booksauthors.Book) error                 // изменить книгу
	DeleteBook(context.Context, int) error                                    // удалить книгу
}

type GRPCserver struct {
	authService AuthService
	bookService BookService
}

func New(a authservice.AuthRepository, b bookservice.BookRepository) *GRPCserver {
	return &GRPCserver{authService: a, bookService: b}
}

func (g *GRPCserver) GetBooks(ctx context.Context, reg *books.GetBooksReq) (*books.GetBooksResp, error) {
	b, err := g.bookService.ReadBook(ctx, booksauthors.Book{
		Book_id:     int(*reg.Id),
		Book_title:  *reg.Tittle,
		Author_id:   int(*reg.AuthorID),
		Author_name: *reg.AuthorName})
	if err != nil {
		return nil, err
	}
	res := books.GetBooksResp{}
	res.Books = make([]*books.Book, 0, len(b))
	for _, book := range b {
		res.Books = append(res.Books, &books.Book{
			BookID:     uint64(book.Book_id),
			BookTitle:  book.Book_title,
			Price:      uint64(book.Price),
			AuthorName: book.Author_name,
		})
	}
	return &res, nil
}
func (g *GRPCserver) PostBooks(ctx context.Context, reg *books.PostBooksReq) (*books.PostBooksResp, error) {
	err := g.bookService.CreateBook(ctx, booksauthors.Book{
		Book_title:  *reg.Tittle,
		Author_name: *reg.AuthorName,
		Price:       uint(reg.Price),
	})
	if err != nil {
		return nil, err
	}
	return &books.PostBooksResp{}, nil
}

func (g *GRPCserver) PatchBooks(ctx context.Context, reg *books.PatchBooksReq) (*books.PatchBooksResp, error) {
	err := g.bookService.UpdateBook(ctx, int(*reg.Id), booksauthors.Book{
		Book_title:       *reg.Tittle,
		Author_name:      *reg.AuthorName,
		Authors_old_name: *reg.AuthorsOldName,
		Price:            uint(reg.Price),
	})
	if err != nil {
		return nil, err
	}
	return &books.PatchBooksResp{}, nil
}

func (g *GRPCserver) DelBooks(ctx context.Context, reg *books.DelBooksReq) (*books.DelBooksResp, error) {
	err := g.bookService.DeleteBook(ctx, int(reg.BookID))
	if err != nil {
		return nil, err
	}
	return &books.DelBooksResp{}, nil
}

func (g *GRPCserver) PostRegistration(ctx context.Context, reg *books.PostRegistrationReq) (*books.PostRegistrationResp, error) {
	err := g.authService.Registration(ctx, reg.Email, reg.Password)
	if err != nil {
		return nil, err
	}
	return &books.PostRegistrationResp{}, nil
}

func (g *GRPCserver) PostLogin(ctx context.Context, reg *books.PostLoginReq) (*books.PostLoginResp, error) {
	id, err := g.authService.Login(ctx, reg.Email, reg.Password)
	if err != nil {
		return nil, err
	}
	return &books.PostLoginResp{Id: uint64(id)}, nil
}

func (g *GRPCserver) PostChekAuth(ctx context.Context, reg *books.PostChekAuthReq) (*books.PostChekAuthResp, error) {
	err := g.authService.ChekAuth(ctx, reg.Email)
	if err != nil {
		return nil, err
	}
	return &books.PostChekAuthResp{}, nil
}
