// https://blog.lelonek.me/a-brief-introduction-to-grpc-in-go-e66e596fe244
// https://github.com/KamilLelonek/grpc-with-go

package main

import (
	"context"
	"log"
	"net"

	pb "github.com/ConduitVC/grpc2/gravatar"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const port = ":50051"

type gravatarService struct{}

func (s *gravatarService) Generate(ctx context.Context, in *pb.GravatarRequest) (*pb.GravatarResponse, error) {
	log.Printf("Received email %v with size %v", in.Email, in.Size)
	return &pb.GravatarResponse{Url: gravatar(in.Email, in.Size)}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterGravatarServiceServer(server, &gravatarService{})
	if err := server.Serve(lis); err != nil {
		log.Fatal(errors.Wrap(err, "Failed to start server!"))
	}
}
