// client/main.go
package main

import (
	"context"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/mattn/go-grpc-hello/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var st bool
	var addr string
	flag.BoolVar(&st, "stream", false, "replys in stream")
	flag.StringVar(&addr, "addr", "127.0.0.1:50051", "server address")
	flag.Parse()

	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("failed to get system cert pool: %v", err)
	}
	creds := credentials.NewClientTLSFromCert(caCertPool, "")

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	arg := "World"
	if flag.NArg() > 0 {
		arg = flag.Arg(0)
	}

	if st {
		reply, err := c.StreamHello(ctx, &pb.HelloRequest{Name: arg})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		for {
			r, err := reply.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("could not greet: %v", err)
			}
			fmt.Printf("Greeting: %q\n", r.GetMessage())
		}
	} else {
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: arg})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		fmt.Printf("Greeting: %q\n", r.GetMessage())
	}
}
