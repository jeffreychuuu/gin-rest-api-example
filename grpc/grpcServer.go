package grpc

import (
	"context"
	"gin-rest-api-example/models"
	pb "gin-rest-api-example/proto"
	"gin-rest-api-example/services"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func StartGrpcServer() {
	// Create gRPC Server
	const host = "localhost"
	const port = "5000"
	lis, err := net.Listen("tcp", host + ":" + port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		s := grpc.NewServer()
		log.Printf("gRPC server is running in port: %s.", port)

		pb.RegisterBookServiceServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (s server) CreateBook(c context.Context, input *pb.CreateBookInput) (*pb.Result, error) {
	request := models.CreateBookInput{Title: input.Title, Author: input.Author, IsEnable: input.IsEnable}
	response := services.CreateBook(request)

	return &pb.Result{Code: int32(response.Code), Message: response.Message,Data: response.Data.(string)}, nil
}