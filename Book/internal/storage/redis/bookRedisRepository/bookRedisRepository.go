package bookredisrepository

import (
	booksauthors "Book/pkg/models/booksAuthors"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

// создаем логгер
var fileError, _ = os.OpenFile("log_BookRepository._error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
var logError = log.New(fileError, "ERROR:", log.LstdFlags|log.Lshortfile)

type Client struct {
	client *redis.Client
}

func New(addr string) Client {
	var C Client
	opt, err := redis.ParseURL(addr)
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	C.client = redis.NewClient(opt)
	return C
}
func (C *Client) SetBook(key string, book booksauthors.Book) error {
	b, err := json.Marshal(book)
	if err != nil {
		logError.Println(err)
		return err
	}
	err = C.client.Set(key, string(b), 0).Err()
	if err != nil {
		logError.Println(err)
		return err
	}
	return nil
}
func (C *Client) GetBook(key string) (*booksauthors.Book, error) {
	var book booksauthors.Book
	bookByte, err := C.client.Get(key).Result()
	if err != nil {
		return nil, fmt.Errorf("this key does not exist")
	}

	err = json.Unmarshal([]byte(bookByte), &book)
	if err != nil {
		logError.Println(err)
		return nil, err
	}
	return &book, nil
}
func (C *Client) DelBook(key string) error {
	err := C.client.Del(key)
	if err != nil {
		return fmt.Errorf("error del")
	}
	return nil
}
func (C *Client) Close() {
	C.client.Close()
}
