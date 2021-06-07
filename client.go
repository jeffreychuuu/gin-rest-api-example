package main

import (
	"context"
	"gin-rest-api-example/models"
	pb "gin-rest-api-example/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBookServiceClient(conn)

	input := models.CreateBookInput{
		Title: "Test", Author: "Jeffrey", IsEnable: true,
	}
	createBook(c, input)
}

func createBook(c pb.BookServiceClient, input models.CreateBookInput) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.CreateBook(ctx, &pb.CreateBookInput{Author: input.Author, Title: input.Title, IsEnable: input.IsEnable})
	if err != nil {
		log.Fatalf("Could not createBook: %v", err)
	}
	log.Printf("gRPC response: %s", res.Data)
}
