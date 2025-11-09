package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/mattn/go-grpc-hello/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}

func (s *server) StreamHello(req *pb.HelloRequest, stream grpc.ServerStreamingServer[pb.HelloReply]) error {
	greetings := []string{"Hello", "Hi", "Hey"}
	for _, greeting := range greetings {
		reply := &pb.HelloReply{Message: greeting + ", " + req.Name + "!"}
		if err := stream.Send(reply); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
