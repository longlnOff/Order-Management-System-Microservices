package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/longln/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/longln/common/api"
)



var (
	httpAddress = common.EnvString("HTTP_ADDRESS", ":8080")
	ordersServiceAddress = "localhost:3000"
)

func main() {

	conn, err := grpc.NewClient(ordersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to order service at address", ordersServiceAddress, " error: ", err)
	}
	defer conn.Close()
	log.Println("connected to order service at address", ordersServiceAddress)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("starting server at %s", httpAddress)
	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		log.Fatal("failed to start server", err)
	}
}