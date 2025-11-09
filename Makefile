all: go-grpc-hello client/client

go-grpc-hello: pb/greeter.pb.go main.go
	go build -o go-grpc-hello main.go

client/client: client/main.go pb/greeter.pb.go
	go build -o client/client client/main.go

pb/greeter.pb.go: greeter.proto
	protoc --go_out=. --go-grpc_out=. greeter.proto
