package main

import (
	"api-gw/internal/api"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
}

func main() {
	conn, err := grpc.NewClient(os.Getenv("HOST")+":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	appi := api.New(conn)
	log.Fatal(http.ListenAndServe(":80", appi.Router()))
}
