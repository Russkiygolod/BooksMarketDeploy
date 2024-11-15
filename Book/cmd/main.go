package main

import (
	"Book/internal/app"
	"Book/internal/books"
	authservice "Book/internal/service/authService"
	bookservice "Book/internal/service/bookService"
	authrepository "Book/internal/storage/postgreSQL/AuthRepository"
	bookrepository "Book/internal/storage/postgreSQL/BookRepository"
	bookredisrepository "Book/internal/storage/redis/bookRedisRepository"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	//PostgreSQL
	// Получаем переменные окружения
	userPostgres, passwordPostgres, dbNamePostgres := os.Getenv("DB_USER_P"), os.Getenv("DB_PASSWORD_P"), os.Getenv("DB_NAME_P")
	//userRedis, passwordRedis, dbNameRedis := os.Getenv("DB_USER_R"), os.Getenv("DB_PASSWORD_R"), os.Getenv("DB_NAME_R")
	url := "postgres://" + userPostgres + ":" + passwordPostgres + "@" + os.Getenv("HOST_P") + ":" + "5432" + "/" + dbNamePostgres
	conn, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Failed to init DB conf - %v", err)
	}
	//Redis
	//addr := "redis://" + userRedis + passwordRedis + ":@localhost:6379/" + dbNameRedis
	addr := "redis://:@" + os.Getenv("HOST_R") + ":6379/"
	Redis := bookredisrepository.New(addr)
	defer Redis.Close()
	bookRepository := bookrepository.New(conn)
	authRepository := authrepository.New(conn)
	authservice := authservice.NewAuthService(authRepository)
	bookservice := bookservice.NewBookService(bookRepository, &Redis)
	s := grpc.NewServer()
	srv := app.New(authservice, bookservice)
	books.RegisterBooksServer(s, srv)
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	err = s.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
}
