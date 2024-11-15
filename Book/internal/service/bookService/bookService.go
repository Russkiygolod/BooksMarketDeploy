package bookservice

import (
	booksauthors "Book/pkg/models/booksAuthors"
	"context"
	"fmt"
)

//go:generate mockgen -source=bookService.go -destination=mocks/book_mock.go -package=mock

type BookRepository interface {
	CreateBook(context.Context, booksauthors.Book) error                      // создать книгу
	ReadBook(context.Context, booksauthors.Book) ([]booksauthors.Book, error) // прочитать книгу
	UpdateBook(context.Context, int, booksauthors.Book) error                 // изменить книгу
	DeleteBook(context.Context, int) error                                    // удалить книгу
}

type bookRedisRepository interface {
	SetBook(key string, book booksauthors.Book) error
	GetBook(key string) (*booksauthors.Book, error)
	DelBook(key string) error
}

type bookServiceStruct struct {
	bookRep      BookRepository
	bookRedisRep bookRedisRepository
}

func NewBookService(b BookRepository, bb bookRedisRepository) bookServiceStruct {
	var book bookServiceStruct
	book.bookRep = b
	book.bookRedisRep = bb
	return book
}

func (b bookServiceStruct) CreateBook(ctx context.Context, B booksauthors.Book) error {
	return b.bookRep.CreateBook(ctx, B)
}
func (b bookServiceStruct) ReadBook(ctx context.Context, B booksauthors.Book) ([]booksauthors.Book, error) {
	// все книги всегда постгрес
	if B.Book_title == "" {
		return b.bookRep.ReadBook(ctx, B)
	}
	var bookPostgres []booksauthors.Book
	// сначала идем в редис,
	book, err := b.bookRedisRep.GetBook(B.Book_title)
	// если там нет идем в постгрес
	if book == nil {
		bookPostgres, err = b.bookRep.ReadBook(ctx, B)
		if err != nil {
			return nil, err
		}
		// добавляем в редис
		for _, bookRedis := range bookPostgres {
			err := b.bookRedisRep.SetBook(bookRedis.Book_title, bookRedis)
			if err != nil {
				return nil, err
			}
		}
	} else if err != nil {
		return nil, err
	} else {
		bookPostgres = append(bookPostgres, *book)
	}
	// и отдаем в апи гв
	return bookPostgres, nil
}
func (b bookServiceStruct) UpdateBook(ctx context.Context, id int, B booksauthors.Book) error {
	// идем в постгрес делаем изменение
	err := b.bookRep.UpdateBook(ctx, id, B)
	if err != nil {
		return err
	}
	// забираем данные по измененной книге
	var book booksauthors.Book
	book.Book_id = id
	bookPostgres, err := b.bookRep.ReadBook(ctx, book)
	if err != nil {
		return err
	}
	// обновляем информацию в редис
	for _, bookRedis := range bookPostgres {
		fmt.Println(bookRedis)
		b.bookRedisRep.SetBook(bookRedis.Book_title, bookRedis)
	}
	return nil
}
func (b bookServiceStruct) DeleteBook(ctx context.Context, id int) error {
	// забираем данные по удаленной книги
	var book booksauthors.Book
	book.Book_id = id
	bookPostgres, err := b.bookRep.ReadBook(ctx, book)
	if err != nil {
		return err
	}
	// идем в постгрес удаляем книгу
	err = b.bookRep.DeleteBook(ctx, id)
	if err != nil {
		return err
	}
	// удаляем книгу так же из редис
	for _, bookRedis := range bookPostgres {
		b.bookRedisRep.DelBook(bookRedis.Book_title)
	}
	return nil
}
