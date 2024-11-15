package authrepository

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// создаем логгер
var fileError, _ = os.OpenFile("log_AuthRepository_error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
var logError = log.New(fileError, "ERROR:", log.LstdFlags|log.Lshortfile)

// Хранилище данных.
type Store struct {
	cl *pgxpool.Pool
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

// регистрация пользователя
func (s *Store) Registration(ctx context.Context, email, password string) error {
	var err error
	var search bool
	tx, err := s.cl.Begin(ctx)
	if err != nil {
		logError.Println(err)
		return err
	}
	// Проверяем, что пользователь с таким email еще не зарегистрирован
	rows, err := tx.Query(
		ctx,
		`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 LIMIT 1);
		;`, email)
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
		if search == true {
			tx.Rollback(ctx)
			return fmt.Errorf("user alredy exist")
		}
	}

	// регистрируем нового пользователя
	_, err = tx.Exec(
		ctx, `
		INSERT INTO users (email, password)
		VALUES($1, $2)
		 ;`, email, password)
	if err != nil {
		logError.Println(err)
		return fmt.Errorf("user registration error")
	}
	tx.Commit(ctx)
	return nil
}

// проверка логина и пароля для генерации токена
func (s *Store) Login(ctx context.Context, email, password string) (int, error) {
	var err error
	var id int
	rows, err := s.cl.Query(
		ctx,
		`
		SELECT id
		FROM users
		WHERE email = $1 and password = $2
		 ;`, email, password)
	if err != nil {
		return 0, fmt.Errorf("login or password error")
	}
	for rows.Next() {
		err = rows.Scan(
			&id,
		)
		if err != nil {
			return 0, fmt.Errorf("login or password error")
		}
	}
	return id, nil
}

func (s *Store) ChekAuth(ctx context.Context, email string) error {
	var emailSearch string
	rows, err := s.cl.Query(
		ctx, `SELECT email FROM users WHERE email = $1`, email)
	if err != nil {
		logError.Println(err)
		return err
	}
	for rows.Next() {
		err = rows.Scan(
			&emailSearch,
		)
	}
	if err != nil {
		logError.Println(err)
		return err
	}
	if emailSearch == "" {
		return fmt.Errorf("incorrect password")
	}
	return nil
}
