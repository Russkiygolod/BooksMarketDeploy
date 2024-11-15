package bookrepository

import (
	booksauthors "Book/pkg/models/booksAuthors"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
)

// создаем логгер
var fileError, _ = os.OpenFile("log_BookRepository._error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
var logError = log.New(fileError, "ERROR:", log.LstdFlags|log.Lshortfile)

// Хранилище данных.
type Store struct {
	cl *pgxpool.Pool
	m  sync.Mutex
}

// Конструктор объекта хранилища.
func New(conf *pgxpool.Config) *Store {
	var postgres Store
	var err error
	postgres.cl, err = pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		log.Fatalf("Failed to init DB conf - %v", err)
	}
	return &postgres
}

// добавляет книгу и автора
func (s *Store) CreateBook(ctx context.Context, Book booksauthors.Book) error {
	s.m.Lock()
	defer s.m.Unlock()
	var authorID int
	authorID = 0
	if Book.Book_title == "" || Book.Author_name == "" || Book.Price == 0 {
		return errors.New("empty fields are not allowed")
	}
	tx, err := s.cl.Begin(ctx)
	if err != nil {
		logError.Println(err)
		return err
	}
	//добаляем информацию о книге
	rows, err := tx.Query(
		ctx,
		`
			INSERT INTO books (title, price)
			VALUES ($1, $2) RETURNING id; 
		`,
		Book.Book_title,
		Book.Price,
	)
	if err != nil {
		tx.Rollback(ctx)
		logError.Println(err)
		return err
	}
	var bookID int
	for rows.Next() {
		err = rows.Scan(
			&Book.Book_id,
		)
		if err != nil {
			tx.Rollback(ctx)
			logError.Println(err)
			return err
		}
		bookID = Book.Book_id
	}
	// забираем id автора по имени, если автор есть в БД
	rows, err = tx.Query(
		ctx,
		`
			SELECT id
			FROM authors
			WHERE name = $1;
	   				`,
		Book.Author_name,
	)
	if err != nil {
		tx.Rollback(ctx)
		logError.Println(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(
			&Book.Author_id,
		)
		if err != nil {
			tx.Rollback(ctx)
			logError.Println(err)
			return err
		}
		authorID = Book.Author_id
	}
	// если автора нет в БД создаем запись в БД и забираем id автора по имени
	if authorID == 0 {
		rows, err = tx.Query(
			ctx,
			`
				INSERT INTO authors (name) 
				VALUES ($1) 
				RETURNING id;
						   `,
			Book.Author_name,
		)
		if err != nil {
			tx.Rollback(ctx)
			logError.Println(err)
			return err
		}

		for rows.Next() {
			err = rows.Scan(
				&Book.Author_id,
			)
			if err != nil {
				tx.Rollback(ctx)
				logError.Println(err)
				return err
			}
		}
		authorID = Book.Author_id
	}
	// доавляем в таблицу id автора и id книги
	_, err = tx.Exec(

		ctx,
		`
			INSERT INTO authors_books (books_id, authors_id)
			VALUES ($1,$2);
				`,
		bookID,
		authorID,
	)
	if err != nil {
		tx.Rollback(ctx)
		logError.Println(err)
		return err
	}
	tx.Commit(ctx)
	return nil
}

// чтение книги
func (s *Store) ReadBook(ctx context.Context, B booksauthors.Book) ([]booksauthors.Book, error) {
	var booksAuthors []booksauthors.Book
	if B.Book_id != 0 {
		// забираем информацию по id книги
		rows, err := s.cl.Query(
			ctx,
			`
			SELECT books.id, books.title, books.price, authors.name 
			FROM books
			JOIN authors_books as ab ON books.id = ab.books_id
			JOIN authors ON ab.authors_id = authors.id
			WHERE books.id = $1
	   				`,
			B.Book_id,
		)
		if err != nil {
			logError.Println(err)
			return nil, err
		}

		for rows.Next() {
			err = rows.Scan(
				&B.Book_id,
				&B.Book_title,
				&B.Price,
				&B.Author_name,
			)
			if err != nil {
				logError.Println(err)
				return nil, err
			}
			booksAuthors = append(booksAuthors, B)
		}
	} else if B.Book_title != "" {
		// забираем информацию по названию книги
		rows, err := s.cl.Query(
			ctx,
			`
			SELECT books.id, books.title, books.price, authors.name 
			FROM books
			JOIN authors_books as ab ON books.id = ab.books_id
			JOIN authors ON ab.authors_id = authors.id
			WHERE books.title = $1
	   				`,
			B.Book_title,
		)
		if err != nil {
			logError.Println(err)
			return nil, err
		}

		for rows.Next() {
			err = rows.Scan(
				&B.Book_id,
				&B.Book_title,
				&B.Price,
				&B.Author_name,
			)
			if err != nil {
				logError.Println(err)
				return nil, err
			}
			booksAuthors = append(booksAuthors, B)
		}
	} else if B.Author_id != 0 {
		// забираем информацию по id автора
		rows, err := s.cl.Query(
			ctx,
			`
			SELECT books.id, books.title, books.price, authors.name 
			FROM books
			JOIN authors_books as ab ON books.id = ab.books_id
			JOIN authors ON ab.authors_id = authors.id
			WHERE authors.id = $1
	   				`,
			B.Author_id,
		)
		if err != nil {
			logError.Println(err)
			return nil, err
		}
		for rows.Next() {
			err = rows.Scan(
				&B.Book_id,
				&B.Book_title,
				&B.Price,
				&B.Author_name,
			)
			if err != nil {
				logError.Println(err)
				return nil, err
			}
			booksAuthors = append(booksAuthors, B)
		}
	} else if B.Author_name != "" {
		// забираем информацию по имени автора
		rows, err := s.cl.Query(
			ctx,
			`
			SELECT books.id, books.title, books.price, authors.name 
			FROM books
			JOIN authors_books as ab ON books.id = ab.books_id
			JOIN authors ON ab.authors_id = authors.id
			WHERE authors.name = $1
			
	   				`,
			B.Author_name,
		)
		if err != nil {
			logError.Println(err)
			return nil, err
		}

		for rows.Next() {
			err = rows.Scan(
				&B.Book_id,
				&B.Book_title,
				&B.Price,
				&B.Author_name,
			)
			if err != nil {
				logError.Println(err)
				return nil, err
			}
			booksAuthors = append(booksAuthors, B)
		}
	} else {
		// забираем все книги
		rows, err := s.cl.Query(
			ctx,
			`
			SELECT books.id, books.title, books.price, authors.name
			FROM books
			JOIN authors_books as ab ON books.id = ab.books_id
			JOIN authors ON ab.authors_id = authors.id
	   				`,
		)
		if err != nil {
			logError.Println(err)
			return nil, err
		}
		for rows.Next() {
			err = rows.Scan(
				&B.Book_id,
				&B.Book_title,
				&B.Price,
				&B.Author_name,
			)
			if err != nil {
				logError.Println(err)
				return nil, err
			}
			booksAuthors = append(booksAuthors, B)
		}
	}
	return booksAuthors, nil
}

// изменение книги
func (s *Store) UpdateBook(ctx context.Context, bookID int, B booksauthors.Book) error {
	s.m.Lock()
	defer s.m.Unlock()
	var authorID int
	var authorOldID int
	// если меняется имя автора
	if B.Author_name != "" {
		tx, err := s.cl.Begin(ctx)
		if err != nil {
			logError.Println(err)
			return err
		}
		// находим id старого автора по имени
		rows, err := tx.Query(
			ctx,
			`
			SELECT id
			FROM authors
			WHERE name = $1
		`,
			B.Authors_old_name,
		)
		if err != nil {
			tx.Rollback(ctx)
			logError.Println(err)
			return err
		}
		for rows.Next() {
			err = rows.Scan(
				&B.Author_id,
			)
			if err != nil {
				tx.Rollback(ctx)
				logError.Println(err)
				return err
			}
			authorOldID = B.Author_id
		}
		// находим id нового автора по имени если он есть в БД
		rows, err = tx.Query(
			ctx,
			`
			SELECT id
			FROM authors
			WHERE name = $1;
	   				`,
			B.Author_name,
		)
		if err != nil {
			tx.Rollback(ctx)
			logError.Println(err)
			return err
		}
		for rows.Next() {
			err = rows.Scan(
				&B.Author_id,
			)
			if err != nil {
				tx.Rollback(ctx)
				logError.Println(err)
				return err
			}
			authorID = B.Author_id
		}
		// если автора нет в БД добавляем автора и забираем id нового автора по имени
		if authorID == 0 {
			rows, err = tx.Query(
				ctx,
				`
				INSERT INTO authors(name)
				VALUES ($1)
				RETURNING id
						   `,
				B.Author_name,
			)
			if err != nil {
				tx.Rollback(ctx)
				logError.Println(err)
				return err
			}
			for rows.Next() {
				err = rows.Scan(
					&B.Author_id,
				)
				if err != nil {
					tx.Rollback(ctx)
					logError.Println(err)
					return err
				}
				authorID = B.Author_id
			}
		}
		// доавляем в таблицу id нового автора и id книги
		_, err = tx.Exec(
			ctx,
			`
			UPDATE authors_books 
			SET authors_id = $1
			WHERE books_id = $2 AND authors_id = $3;`,
			authorID,
			bookID,
			authorOldID,
		)
		if err != nil {
			tx.Rollback(ctx)
			logError.Println(err)
			return err
		}
		err = tx.Commit(ctx)
		if err != nil {
			logError.Println(err)
			return err
		}
	}
	// если меняется название
	if B.Book_title != "" {
		_, err := s.cl.Exec(
			ctx,
			`
			UPDATE books
			SET title = $2
			WHERE id = $1;`,
			bookID,
			B.Book_title,
		)
		if err != nil {
			logError.Println(err)
			return err
		}
	}
	// если меняется цена
	if B.Price != 0 {
		_, err := s.cl.Exec(
			ctx,
			`
			UPDATE books
			SET price = $2
			WHERE id = $1;`,
			bookID,
			B.Price,
		)
		if err != nil {
			logError.Println(err)
			return err
		}
	}
	return nil
}

// удаляет книгу
func (s *Store) DeleteBook(ctx context.Context, id int) error {
	s.m.Lock()
	defer s.m.Unlock()
	var search bool
	tx, err := s.cl.Begin(ctx)
	if err != nil {
		logError.Println(err)
		return err
	}
	// Проверяем, что пользователь книга, которую мы хотим удалить существует
	rows, err := tx.Query(
		ctx,
		`
		SELECT EXISTS(SELECT 1 FROM books WHERE id = $1 LIMIT 1);
		;`, id)
	if err != nil {
		tx.Rollback(ctx)
		logError.Println(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(
			&search,
		)
		if err != nil {
			tx.Rollback(ctx)
			logError.Println(err)
			return err
		}
		if search == false {
			tx.Rollback(ctx)
			return fmt.Errorf("the book does not exist")
		}
	}
	// удаляет запись из таблицы authors_books
	_, err = tx.Exec(
		ctx, `DELETE FROM authors_books WHERE books_id = $1 ;`, id)
	if err != nil {
		tx.Rollback(ctx)
		logError.Println(err)
		return err
	}
	// удаляет запись из таблицы books
	_, err = tx.Exec(
		ctx, `DELETE FROM books WHERE id = $1 ;`, id)
	if err != nil {
		tx.Rollback(ctx)
		logError.Println(err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		logError.Println(err)
		return err
	}
	return nil
}
