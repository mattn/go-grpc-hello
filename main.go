package main

import (
	"context"
	"log"
	"net"

	pb "github.com/mattn/go-grpc-hello/pb" // proto生成ディレクトリを適宜変更

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051") // Cloudflare例のポート
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
