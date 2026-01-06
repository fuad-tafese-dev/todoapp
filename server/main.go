package main

import (
	"log"
	"net"
	"os"

	pb "todoapp/proto/pb"

	"google.golang.org/grpc"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("usage: server [IP_ADDR] (e.g., server :50051)")
	}

	addr := args[0]
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	defer func() {
		if err := lis.Close(); err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
	}()

	// Create gRPC server
	s := grpc.NewServer()
	
	// Register our service
	pb.RegisterTodoServiceServer(s, &server{
		store: New(),
	})

	log.Printf("Todo server listening at %s\n", addr)
	
	// Start serving
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}